package commandmock

import "os/exec"

// CommandExecutor is an interface for executing commands
type CommandExecutor interface {
	Execute(command string, args ...string) ([]byte, error)
}

// MockCommandExecutor is a mock implementation of CommandExecutor
type MockCommandExecutor struct {
	Output string
	Err    error
}

// Execute executes a mock command and returns the output
func (m *MockCommandExecutor) Execute(_ string, _ ...string) ([]byte, error) {
	return []byte(m.Output), m.Err
}

// RealCommandExecutor is a real implementation of CommandExecutor
type RealCommandExecutor struct {
}

// Execute executes an actual command and returns the output
func (r *RealCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}
