package checks_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TestUACCheck is a function that tests the UACCheck function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the UACCheck function with different scenarios. It uses a mock implementation of the CommandExecutor interface to simulate the behavior of the command execution for checking the User Account Control (UAC) level. Each test case checks if the UACCheck function correctly identifies the UAC level based on the simulated command output. The function asserts that the returned Check instance contains the expected results.
func TestUACCheck(t *testing.T) {
	tests := []struct {
		name        string
		executorUAC *commandmock.MockCommandExecutor
		want        checks.Check
	}{
		{
			name:        "UAC disabled",
			executorUAC: &commandmock.MockCommandExecutor{Output: "0", Err: nil},
			want:        checks.NewCheckResult("UAC", "UAC is disabled."),
		},
		{
			name:        "UAC enabled for apps and settings",
			executorUAC: &commandmock.MockCommandExecutor{Output: "2", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to your computer "+
				"and for changing your settings."),
		},
		{
			name:        "UAC enabled for apps but not for settings",
			executorUAC: &commandmock.MockCommandExecutor{Output: "5", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to "+
				"your computer."),
		},
		{
			name:        "unknown UAC level",
			executorUAC: &commandmock.MockCommandExecutor{Output: "3", Err: nil},
			want:        checks.NewCheckResult("UAC", "Unknown UAC level"),
		},
		{
			name:        "UAC error",
			executorUAC: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("error retrieving UAC")},
			want: checks.NewCheckErrorf("UAC", "error retrieving UAC",
				errors.New("error retrieving UAC")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.UACCheck(tt.executorUAC)
			require.Equal(t, tt.want, got)
		})
	}
}
