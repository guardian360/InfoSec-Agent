package network_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

// TestSMBCheck is a function that tests the SmbCheck function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the SmbEnabled function with different scenarios. It uses a mock implementation of the CommandExecutor interface to simulate the behavior of the command execution for checking the status of SMB1 and SMB2 protocols. Each test case checks if the SmbEnabled function correctly identifies the status of the SMB protocols based on the simulated command output. The function asserts that the returned string and error match the expected results.
func TestSMBCheck(t *testing.T) {
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
			got := network.SmbCheck(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestSmbEnabled is a function that tests the SmbEnabled function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
func TestSmbEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor *mocking.MockCommandExecutor
		SMB1     string
		SMB2     string
		resultID int
	}{
		{
			name:     "SMB1 and SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "True True", Err: nil},
			SMB1:     "SMB1: enabled",
			SMB2:     "SMB2: enabled",
			resultID: 3,
		},
		{
			name:     "Only SMB1 enabled",
			executor: &mocking.MockCommandExecutor{Output: "True False", Err: nil},
			SMB1:     "SMB1: enabled",
			SMB2:     "SMB2: not enabled",
			resultID: 1,
		},
		{
			name:     "Only SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "False True", Err: nil},
			SMB1:     "SMB1: not enabled",
			SMB2:     "SMB2: enabled",
			resultID: 2,
		},
		{
			name:     "Neither SMB1 nor SMB2 enabled",
			executor: &mocking.MockCommandExecutor{Output: "False False", Err: nil},
			SMB1:     "SMB1: not enabled",
			SMB2:     "SMB2: not enabled",
			resultID: 0,
		},
		{
			name:     "command error",
			executor: &mocking.MockCommandExecutor{Output: "", Err: errors.New("command error")},
			SMB1:     "",
			SMB2:     "",
			resultID: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSMB1, gotSMB2, gotResultID, _ := network.SmbEnabled(tt.executor, 0)
			require.Equal(t, tt.SMB1, gotSMB1)
			require.Equal(t, tt.SMB2, gotSMB2)
			require.Equal(t, tt.resultID, gotResultID)
		})
	}
}
