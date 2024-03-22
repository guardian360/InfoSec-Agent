package utils

import "os/exec"

type MockCommandExecutor struct {
	Output string
	Err    error
}

func (m *MockCommandExecutor) Execute() ([]byte, error) {
	return []byte(m.Output), m.Err
}

type CommandExecutor interface {
	Execute(command string, args ...string) ([]byte, error)
}

type RealCommandExecutor struct {
}

func (r *RealCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}
