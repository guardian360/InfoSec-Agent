package mocking

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// TODO: Update documentation
// ProgramLister is an interface that defines a method for listing installed programs.
//
// The ListInstalledPrograms method takes a directory path as input and returns a slice of strings representing the names of installed programs, or an error if the operation fails.
//
// This interface is used in the PasswordManager function to abstract the operation of listing installed programs, allowing for different implementations that can be swapped out as needed. This is particularly useful for testing, as a mock implementation can be used to simulate different scenarios.
type ProgramLister interface {
	ListInstalledPrograms(directory string) ([]string, error)
}

// TODO: Update documentation
// RealProgramLister is a struct that implements the ProgramLister interface.
//
// It provides a real-world implementation of the ListInstalledPrograms method, which lists all installed programs in a given directory by reading the directory's contents and returning the names of all subdirectories, which represent installed programs.
//
// This struct is used in the PasswordManager function to list installed programs when checking for the presence of known password managers.
type RealProgramLister struct{}

// TODO: Update documentation
// MockProgramLister is a mock implementation of the ProgramLister interface used for testing.
// It provides a controlled environment to simulate the behavior of the real ProgramLister,
// allowing us to test how our code interacts with the ProgramLister interface.
// This is particularly useful for testing the PasswordManager function, as it allows us to
// simulate different scenarios of installed programs on a system.
type MockProgramLister struct {
	mock.Mock
}

// TODO: Update documentation
// ListInstalledPrograms is a method of the RealProgramLister struct that lists all installed programs in a given directory.
//
// Parameters:
//   - directory (string): The path of the directory to list the installed programs from.
//
// Returns:
//   - []string: A slice of strings representing the names of installed programs.
//   - error: An error object that describes the error, if any occurred.
//
// This method reads the contents of the specified directory and returns the names of all subdirectories, which represent installed programs. If an error occurs during the operation, it returns the error.
func (rpl RealProgramLister) ListInstalledPrograms(directory string) ([]string, error) {
	var programs []string
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			programs = append(programs, file.Name())
		}
	}
	return programs, nil
}

// TODO: Update documentation
// ListInstalledPrograms is a method of the MockProgramLister struct that simulates the behavior of the real ProgramLister's ListInstalledPrograms method.
//
// Parameters:
//   - directory (string): The path of the directory to list the installed programs from.
//
// Returns:
//   - []string: A slice of strings representing the names of installed programs.
//   - error: An error object that describes the error, if any occurred.
//
// This method is used in tests to control the output of the ListInstalledPrograms method, allowing us to simulate different scenarios of installed programs on a system. It returns the values provided when the method is mocked in the test.
func (m *MockProgramLister) ListInstalledPrograms(directory string) ([]string, error) {
	args := m.Called(directory)
	return args.Get(0).([]string), args.Error(1)
}
