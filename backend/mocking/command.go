// Package mocking contains different mocking implementations for various components of Windows.
//
// It contains mocking implementations for:
//   - command execution
//   - file reading and writing
//   - installed programs listing
//   - Windows registry access
//   - username retrieval
package mocking

import (
	"errors"
	"os/exec"
	"syscall"
)

// TODO: Update documentation
// CommandExecutor is an interface that defines a contract for executing system commands.
// It abstracts the details of command execution, allowing for different implementations
// that can either execute real system commands or simulate command execution for testing purposes.
type CommandExecutor interface {
	Execute(command string, args ...string) ([]byte, error)
}

// TODO: Update documentation
// MockCommandExecutor is a mock implementation of the CommandExecutor interface.
// It is used for testing purposes to simulate the behavior of a real command executor.
// This allows tests to control the output and error conditions of command execution,
// ensuring that the code under test can handle various scenarios correctly.
type MockCommandExecutor struct {
	Output string
	Err    error
}

// TODO: Update documentation
// Execute simulates the execution of a system command for testing purposes.
// This method is part of the MockCommandExecutor struct, which is a mock implementation of the CommandExecutor interface.
//
// Parameters:
//   - _: A string representing the system command to be executed. This parameter is ignored in the mock implementation.
//   - _: A variadic string slice representing the arguments to be passed to the command. This parameter is ignored in the mock implementation.
//
// Returns:
//   - A byte slice representing the predefined output of the simulated command execution.
//   - An error that will be non-nil if a predefined error condition is simulated.
//
// This method allows tests to control the outcomes of command execution, ensuring that the code under test can handle various scenarios correctly.
func (m *MockCommandExecutor) Execute(_ string, _ ...string) ([]byte, error) {
	if m.Output == "test1" {
		return nil, errors.New("test error")
	}
	return []byte(m.Output), m.Err
}

// TODO: Update documentation
// RealCommandExecutor is a struct that implements the CommandExecutor interface. It is responsible for executing actual system commands.
// The execution is performed using the os/exec package, which allows the commands to be run and their output to be captured.
// This struct provides a concrete implementation of the CommandExecutor interface, enabling interaction with the system's command line interface.
type RealCommandExecutor struct {
}

// TODO: Update documentation
// Execute runs a system command and returns its output.
// This method is part of the RealCommandExecutor struct, which is an implementation of the CommandExecutor interface.
//
// Parameters:
//   - command: A string representing the system command to be executed.
//   - args: A variadic string slice representing the arguments to be passed to the command.
//
// Returns:
//   - A byte slice representing the output of the executed command.
//   - An error that will be non-nil if the command execution fails.
//
// This method uses the os/exec package to execute the command, capturing and returning the output.
// It provides a mechanism for actual interaction with the system's command line interface.
func (r *RealCommandExecutor) Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Output()
}
