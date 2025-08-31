package actions

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/rs/zerolog/log"
)

type Runner struct {
	Command                 []string
	nextCommands            [][]string
	Ctx                     context.Context
	CancelFunc              context.CancelFunc
	OutputFn                func(string)
	OutputMetaFn            func(status ActionStatus)
	DeleteFn                func()
	RemotelyMarkedForDelete bool
}

func (r *Runner) Run() {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if ok {
				log.Err(err).
					Strs("command", r.Command).
					Msg("[WS] Panic while running command:")
			} else {
				log.Error().
					Strs("command", r.Command).
					Msg("[WS] Panic while running command - unable to cast to error")
			}

			r.OutputFn("\r\n\nError running command\r\n")
		}
	}()
	allCommands := make([][]string, 0)
	allCommands = append(allCommands, make([]string, 0))
	for _, cmdArg := range r.Command {
		if cmdArg == "&&" {
			allCommands = append(allCommands, make([]string, 0))
		} else {
			currIndex := len(allCommands) - 1
			allCommands[currIndex] = append(allCommands[currIndex], cmdArg)
		}
	}
	r.Command = allCommands[0]
	if len(allCommands) > 1 {
		log.Debug().
			Strs("first", r.Command).
			Strs("next", allCommands[1]).
			Msg("Running an action with multiple commands")
		r.nextCommands = allCommands[1:]
	}

	cmd := exec.CommandContext(r.Ctx, r.Command[0], r.Command[1:]...)

	r.OutputMetaFn(Pending)
	commandString := strings.Join(r.Command, " ")
	r.OutputFn("$ " + commandString + "\r\n")

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	r.OutputMetaFn(Running)

	go streamOutput(f, r.OutputFn)

	if err := cmd.Wait(); err == nil {
		if len(r.nextCommands) > 0 {
			r.Command = r.nextCommands[0]
			if len(r.nextCommands) > 1 {
				r.nextCommands = r.nextCommands[1:]
			} else {
				r.nextCommands = make([][]string, 0)
			}

			toBroadcast := "\r\n\n\nProceeding to next command in chain\r\n\n"
			r.OutputFn(toBroadcast)

			log.Debug().
				Strs("next", r.Command).
				Msg("Proceeding to next command in chain")
			go r.Run()
			return
		}

		toBroadcast := "\r\n\n\nCommand finished\r\n"
		r.OutputFn(toBroadcast)

		r.OutputMetaFn(Success)

		go r.DeleteFn()
	} else {
		toBroadcast := "\r\n\n\nCommand finished with error\r\n"
		r.OutputFn(toBroadcast)

		r.OutputMetaFn(Failed)
	}
}

func streamOutput(reader io.Reader, outputFn func(string)) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitterFunc)

	for scanner.Scan() {
		text := scanner.Text()

		outputFn(text + "\r")
	}
}

func splitterFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	for i, b := range data {
		if b == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), bytes.TrimSpace(data), nil
	}

	return 0, nil, nil
}
