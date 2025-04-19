package websocket

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
)

// TODO: Change all this to use some kind of terminal instead

type Runner struct {
	Command []string
	TaskId  string
	Ctx     context.Context
}

func (r *Runner) Run() {
	defer func() {
		if rec := recover(); rec != nil {
			println("[WS] Recovered from panic while running command '" +
				strings.Join(r.Command, " ") +
				"':")

			println(rec)

			ConnectionManager.Broadcast(r.TaskId, "error running command")
		}
	}()
	cmd := exec.CommandContext(r.Ctx, r.Command[0], r.Command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	go streamOutput(r.TaskId, stdout)
	go streamOutput(r.TaskId, stderr)

	if err := cmd.Wait(); err == nil {
		ConnectionManager.Broadcast(r.TaskId, "command finished")
	} else {
		ConnectionManager.Broadcast(r.TaskId, "command finished with error")
	}
}

func streamOutput(taskId string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		message := scanner.Text()
		ConnectionManager.Broadcast(taskId, message)
	}
}
