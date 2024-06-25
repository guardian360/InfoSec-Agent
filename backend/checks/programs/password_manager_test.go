package programs_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/stretchr/testify/mock"
)

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
			want:         checks.NewCheckResult(checks.PasswordManagerID, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(mocking.MockProgramLister)
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
		{
			name:      "With Programs86",
			directory: "C:\\Program Files (x86)",
			want:      []string{"Program1", "Program2"},
		},
		{
			name:      "No Programs86",
			directory: "C:\\Program Files (x86)",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLister := new(mocking.MockProgramLister)
			mockLister.On("ListInstalledPrograms", mock.Anything).Return(tt.want, nil)

			result, err := mockLister.ListInstalledPrograms(tt.directory)
			require.Equal(t, tt.want, result)
			if err != nil {
				t.Errorf("Test %s failed. Expected no error, got %v", tt.name, err)
			}
		})
	}
}
