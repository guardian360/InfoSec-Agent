package checks_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
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
		executorLocalGroup        commandmock.CommandExecutor
		executorLocalGroupMembers commandmock.CommandExecutor
		executorYesWord           commandmock.CommandExecutor
		executorNetUser           commandmock.CommandExecutor
		want                      checks.Check
	}{
		{
			name: "wmiObjectError",
			executorLocalGroup: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("Get-WmiObject error")},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command Get-WmiObject", errors.New("Get-WmiObject error")),
		},
		{
			name:                      "guestLocalGroupNotFound",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 0, "Guest localgroup not found"),
		},
		{
			name:               "netLocalGroupError",
			executorLocalGroup: &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("net localgroup error")},
			executorYesWord: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net localgroup", errors.New("net localgroup error")),
		},
		{
			name:               "guestAccountNotFound",
			executorLocalGroup: &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "The command completed successfully.",
				Err: nil},
			executorYesWord: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:            checks.NewCheckResult(checks.GuestAccountID, 0, "Guest account not found"),
		},
		{
			name:                      "YesWordError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("net user yesWord error")},
			executorNetUser: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net user", errors.New("net user yesWord error")),
		},
		{
			name:                      "netUserError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("net user error")},
			want: checks.NewCheckErrorf(checks.GuestAccountID,
				"error executing command net user", errors.New("net user error")),
		},
		{
			name:                      "guestAccountFoundAndActive",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nyes", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 1, "Guest account is active"),
		},
		{
			name:                      "guestAccountFoundAndInactive",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno", Err: nil},
			want:                      checks.NewCheckResult(checks.GuestAccountID, 2, "Guest account is not active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.GuestAccount(
				tt.executorLocalGroup,
				tt.executorLocalGroupMembers,
				tt.executorYesWord,
				tt.executorNetUser,
			)
			require.Equal(t, tt.want, got)
		})
	}
}
