package websocket

import (
	"bufio"
	"bytes"
	"context"
	"github.com/creack/pty"
	"io"
	"os/exec"
)

type Runner struct {
	Command []string
	TaskId  string
	Ctx     context.Context
}

func (r *Runner) Run() {
	defer func() {
		//if rec := recover(); rec != nil {
		//	println("[WS] Recovered from panic while running command '" +
		//		strings.Join(r.Command, " ") +
		//		"':")
		//
		//	println(rec)
		//
		//	ConnectionManager.Broadcast(r.TaskId, "error running command")
		//}
	}()
	cmd := exec.CommandContext(r.Ctx, r.Command[0], r.Command[1:]...)

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	go streamOutput(r.TaskId, f)

	if err := cmd.Wait(); err == nil {
		ConnectionManager.Broadcast(r.TaskId, "command finished")
	} else {
		ConnectionManager.Broadcast(r.TaskId, "command finished with error")
	}
}

func streamOutput(taskId string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitterFunc)

	for scanner.Scan() {
		text := scanner.Text()
		ConnectionManager.Broadcast(taskId, text)
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
