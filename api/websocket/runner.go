package websocket

import (
	"bufio"
	"bytes"
	"context"
	"github.com/creack/pty"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Runner struct {
	Command  []string
	ActionId string
	Ctx      context.Context
	Action   *Action
}

func (r *Runner) Run() {
	defer func() {
		if rec := recover(); rec != nil {
			println("[WS] Recovered from panic while running command '" +
				strings.Join(r.Command, " ") +
				"':")

			println(rec)

			ConnectionManager.Broadcast(r.ActionId, "\r\n\nError running command\r\n")
		}
	}()
	cmd := exec.CommandContext(r.Ctx, r.Command[0], r.Command[1:]...)

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	go streamOutput(r.ActionId, r.Action, f)

	if err := cmd.Wait(); err == nil {
		ConnectionManager.Broadcast(r.ActionId, "\r\n\n\nCommand finished\r\n")

		r.Action.Mutex.Lock()
		r.Action.Status = Success
		r.Action.FinishedAt = time.Now()
		r.Action.Mutex.Unlock()

		// TODO: URGENT - Change this
		go ActionManager.DeleteFinishedAction(r.Action, time.Hour*10)
	} else {
		ConnectionManager.Broadcast(r.ActionId, "\r\n\n\nCommand finished with error\r\n")

		r.Action.Mutex.Lock()
		r.Action.Status = Failed
		r.Action.FinishedAt = time.Now()
		r.Action.Mutex.Unlock()
	}

	r.Action.Mutex.Lock()
	// remove the last \r from the output
	if len(r.Action.Output) > 0 {
		r.Action.Output = r.Action.Output[:len(r.Action.Output)-1]
	}
	r.Action.Mutex.Unlock()
}

func streamOutput(actionId string, action *Action, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitterFunc)

	for scanner.Scan() {
		text := scanner.Text()

		action.Mutex.Lock()
		action.Output += text + "\r"
		action.Mutex.Unlock()

		ConnectionManager.Broadcast(actionId, text)
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
