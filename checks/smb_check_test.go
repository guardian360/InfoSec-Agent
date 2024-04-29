package checks_test

import (
	"errors"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/stretchr/testify/require"
)

// TestCheckSMB is a function that tests the SmbEnabled function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the SmbEnabled function with different scenarios. It uses a mock implementation of the CommandExecutor interface to simulate the behavior of the command execution for checking the status of SMB1 and SMB2 protocols. Each test case checks if the SmbEnabled function correctly identifies the status of the SMB protocols based on the simulated command output. The function asserts that the returned string and error match the expected results.
func TestCheckSMB(t *testing.T) {
	tests := []struct {
		name     string
		executor *mocking.MockCommandExecutor
		want     checks.Check
		wantErr  bool
	}{
		{
			name:     "SMB1 and SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "True True", Err: nil},
			want:     checks.NewCheckResult(checks.SmbID, 3, "SMB1: enabled", "SMB2: enabled"),
			wantErr:  false,
		},
		{
			name:     "Only SMB1 enabled",
			executor: &mocking.MockCommandExecutor{Output: "True False", Err: nil},
			want:     checks.NewCheckResult(checks.SmbID, 1, "SMB1: enabled", "SMB2: not enabled"),
			wantErr:  false,
		},
		{
			name:     "Only SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "False True", Err: nil},
			want:     checks.NewCheckResult(checks.SmbID, 2, "SMB1: not enabled", "SMB2: enabled"),
			wantErr:  false,
		},
		{
			name:     "Neither SMB1 nor SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "False False", Err: nil},
			want:     checks.NewCheckResult(checks.SmbID, 0, "SMB1: not enabled", "SMB2: not enabled"),
			wantErr:  false,
		},
		{
			name:     "command error",
			executor: &mocking.MockCommandExecutor{Output: "", Err: errors.New("command error")},
			want:     checks.NewCheckError(checks.SmbID, errors.New("command error")),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.SmbCheck(tt.executor)
			if tt.wantErr {
				require.Equal(t, tt.want, got)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
