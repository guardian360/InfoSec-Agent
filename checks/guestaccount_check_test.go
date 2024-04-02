package checks_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

func TestGuestAccount(t *testing.T) {
	tests := []struct {
		name                      string
		executorLocalGroup        commandmock.CommandExecutor
		executorLocalGroupMembers commandmock.CommandExecutor
		executorYesWord           commandmock.CommandExecutor
		executorNetUser           commandmock.CommandExecutor
		want                      checks.Check
	}{
		// TODO: Add test cases.
		{
			name:                      "wmiObjectError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "", Err: errors.New("Get-WmiObject error")},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckErrorf("Guest account", "error executing command Get-WmiObject", errors.New("Get-WmiObject error")),
		},
		{
			name:                      "guestLocalGroupNotFound",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckResult("Guest account", "Guest localgroup not found"),
		},
		{
			name:                      "netLocalGroupError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("net localgroup error")},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckErrorf("Guest account", "error executing command net localgroup", errors.New("net localgroup error")),
		},
		{
			name:                      "guestAccountNotFound",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "The command completed successfully.", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckResult("Guest account", "Guest account not found"),
		},
		{
			name:                      "YesWordError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "", Err: errors.New("net user yesWord error")},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want:                      checks.NewCheckErrorf("Guest account", "error executing command net user", errors.New("net user yesWord error")),
		},
		{
			name:                      "netUserError",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "", Err: errors.New("net user error")},
			want:                      checks.NewCheckErrorf("Guest account", "error executing command net user", errors.New("net user error")),
		},
		{
			name:                      "guestAccountFoundAndActive",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nyes", Err: nil},
			want:                      checks.NewCheckResult("Guest account", "Guest account is active"),
		},
		{
			name:                      "guestAccountFoundAndInactive",
			executorLocalGroup:        &commandmock.MockCommandExecutor{Output: "             S-1-5-32-546", Err: nil},
			executorLocalGroupMembers: &commandmock.MockCommandExecutor{Output: "-----\r\nguest", Err: nil},
			executorYesWord:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno yes", Err: nil},
			executorNetUser:           &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\n\r\n\r\nno", Err: nil},
			want:                      checks.NewCheckResult("Guest account", "Guest account is not active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.GuestAccount(
				tt.executorLocalGroup,
				tt.executorLocalGroupMembers,
				tt.executorYesWord,
				tt.executorNetUser,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GuestAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
