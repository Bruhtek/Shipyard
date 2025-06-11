package terminals

import (
	"Shipyard/internal/utils"
	"bufio"
	"bytes"
	"context"
	"github.com/creack/pty"
	"github.com/rs/zerolog/log"
	"io"
	"os/exec"
)

type Runner struct {
	Command                 []string
	Ctx                     context.Context
	CancelFunc              context.CancelFunc
	OutputFn                func(string)
	OutputMetaFn            func(status utils.ActionStatus)
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
	cmd := exec.CommandContext(r.Ctx, r.Command[0], r.Command[1:]...)

	r.OutputMetaFn(utils.Pending)

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	r.OutputMetaFn(utils.Running)

	go streamOutput(f, r.OutputFn)

	if err := cmd.Wait(); err == nil {
		toBroadcast := "\r\n\n\nCommand finished\r\n"
		r.OutputFn(toBroadcast)

		r.OutputMetaFn(utils.Success)

		go r.DeleteFn()
	} else {
		toBroadcast := "\r\n\n\nCommand finished with error\r\n"
		r.OutputFn(toBroadcast)

		r.OutputMetaFn(utils.Failed)
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
