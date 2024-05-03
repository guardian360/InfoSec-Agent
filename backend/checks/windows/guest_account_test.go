package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestGuestAccount is a test function that validates the behavior of the GuestAccount function.
//
// It executes a series of test cases, each with different inputs, to ensure that the function behaves as expected in various scenarios.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If a test case fails, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the GuestAccount function correctly identifies the status of the guest account on the Windows system and handles errors as expected.
func TestGuestAccount(t *testing.T) {
	tests := []struct {
		name                      string
		executorLocalGroup        mocking.CommandExecutor
		executorLocalGroupMembers mocking.CommandExecutor
		executorYesWord           mocking.CommandExecutor
		executorNetUser           mocking.CommandExecutor
		want                      checks.Check
	}{
		{
			name: "wmiObjectError",
			executorLocalGroup: &mocking.MockCommandExecutor{Output: "",
				Err: errors.New("Get-WmiObject error")},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &mocking.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command Get-WmiObject", errors.New("Get-WmiObject error")),
		},
		{
			name:                      "guestLocalGroupNotFound",
			executorLocalGroup:        &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &mocking.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 0, "Guest localgroup not found"),
		},
		{
			name:               "netLocalGroupError",
			executorLocalGroup: &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "",
				Err: errors.New("net localgroup error")},
			executorYesWord: &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser: &mocking.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net localgroup", errors.New("net localgroup error")),
		},
		{
			name:               "guestAccountNotFound",
			executorLocalGroup: &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "The command completed successfully.",
				Err: nil},
			executorYesWord: &mocking.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser: &mocking.MockCommandExecutor{Output: "", Err: nil},
			want:            checks.NewCheckResult(checks.GuestAccountID, 0, "Guest account not found"),
		},
		{
			name:                      "YesWordError",
			executorLocalGroup:        &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord: &mocking.MockCommandExecutor{Output: "",
				Err: errors.New("net user yesWord error")},
			executorNetUser: &mocking.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net user", errors.New("net user yesWord error")),
		},
		{
			name:                      "netUserError",
			executorLocalGroup:        &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser: &mocking.MockCommandExecutor{Output: "",
				Err: errors.New("net user error")},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net user", errors.New("net user error")),
		},
		{
			name:                      "guestAccountFoundAndActive",
			executorLocalGroup:        &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nyes", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 1, "Guest account is active"),
		},
		{
			name:                      "guestAccountFoundAndInactive",
			executorLocalGroup:        &mocking.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &mocking.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 2, "Guest account is not active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.GuestAccount(
				tt.executorLocalGroup,
				tt.executorLocalGroupMembers,
				tt.executorYesWord,
				tt.executorNetUser,
			)
			require.Equal(t, tt.want, got)
		})
	}
}
