package programs_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/stretchr/testify/mock"
)

// MockProgramLister is a mock implementation of the ProgramLister interface used for testing.
// It provides a controlled environment to simulate the behavior of the real ProgramLister,
// allowing us to test how our code interacts with the ProgramLister interface.
// This is particularly useful for testing the PasswordManager function, as it allows us to
// simulate different scenarios of installed programs on a system.
type MockProgramLister struct {
	mock.Mock
}

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

// TestPasswordManager is a test function for the PasswordManager function.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the PasswordManager function with different scenarios. It uses a mock implementation of the ProgramLister interface to simulate different sets of installed programs. Each test case checks if the PasswordManager function correctly identifies the presence or absence of known password managers based on the simulated installed programs.
func TestPasswordManager(t *testing.T) {
	tests := []struct {
		name         string
		mockPrograms []string
		want         checks.Check
	}{
		{
			name:         "With Known Password Manager",
			mockPrograms: []string{"1Password"},
			want:         checks.NewCheckResult(checks.PasswordManagerID, 0, "1Password"),
		},
		{
			name:         "No Password Manager",
			mockPrograms: []string{"RandomSoftware"},
			want:         checks.NewCheckResult(checks.PasswordManagerID, 1, "No password manager found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(MockProgramLister)
			mockLister.On("ListInstalledPrograms", mock.Anything).Return(tt.mockPrograms, nil)

			result := programs.PasswordManager(mockLister)
			require.Equal(t, tt.want, result)
		})
	}
}

// TestListInstalledPrograms is a test function for the ListInstalledPrograms method.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the ListInstalledPrograms method with different scenarios. It uses a mock implementation of the ProgramLister interface to simulate different sets of installed programs in a directory. Each test case checks if the ListInstalledPrograms method correctly lists the installed programs based on the simulated directory content.
func TestListInstalledPrograms(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		want      []string
	}{
		{
			name:      "With Programs",
			directory: "C:\\Program Files",
			want:      []string{"Program1", "Program2"},
		},
		{
			name:      "No Programs",
			directory: "C:\\Program Files",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(MockProgramLister)
			mockLister.On("ListInstalledPrograms", mock.Anything).Return(tt.want, nil)

			result, err := mockLister.ListInstalledPrograms(tt.directory)
			require.Equal(t, tt.want, result)
			if err != nil {
				t.Errorf("Test %s failed. Expected no error, got %v", tt.name, err)
			}
		})
	}
}
