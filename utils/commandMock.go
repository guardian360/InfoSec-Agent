package utils

import "os/exec"

type CommandExecutor interface {
	Execute(command string, args ...string) ([]byte, error)
}

type MockCommandExecutor struct {
	Output string
	Err    error
}

func (m *MockCommandExecutor) Execute(_ string, _ ...string) ([]byte, error) {
	return []byte(m.Output), m.Err
}

type RealCommandExecutor struct {
}

func (r *RealCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}
