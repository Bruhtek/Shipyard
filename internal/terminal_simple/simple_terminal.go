package terminal_simple

import (
	"os/exec"
	"strings"
)

// SimpleTerminal
// takes a command, runs it, returns the output (or errors)
// no support for output streaming or cancellations
type SimpleTerminal struct {
	Command string
}

func (t *SimpleTerminal) Run() (string, error) {
	splitCmd := strings.Split(t.Command, " ")
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RunSimpleCommand(command string) (string, error) {
	term := &SimpleTerminal{Command: command}
	return term.Run()
}
