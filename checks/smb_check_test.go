package checks_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

func TestSmbCheck(t *testing.T) {
	tests := []struct {
		name      string
		executor1 *commandmock.MockCommandExecutor
		executor2 *commandmock.MockCommandExecutor
		want      checks.Check
	}{
		{
			name:      "SMB1 and SMB2 enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			want:      checks.NewCheckResult(checks.SmbID, 3, "SMB1: enabled", "SMB2: enabled"),
		},
		{
			name:      "SMB1 enabled and SMB2 not enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckResult(checks.SmbID, 1, "SMB1: enabled", "SMB2: not enabled"),
		},
		{
			name:      "SMB1 not enabled and SMB2 enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			want:      checks.NewCheckResult(checks.SmbID, 2, "SMB1: not enabled", "SMB2: enabled"),
		},
		{
			name:      "SMB1 and SMB2 not enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckResult(checks.SmbID, 0, "SMB1: not enabled", "SMB2: not enabled"),
		},
		{
			name:      "command smb1 error",
			executor1: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("command smb1 error")},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckError(checks.SmbID, errors.New("command smb1 error")),
		},
		{
			name:      "command smb2 error",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("command smb2 error")},
			want:      checks.NewCheckError(checks.SmbID, errors.New("command smb2 error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.SmbCheck(tt.executor1, tt.executor2)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSmbEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor *commandmock.MockCommandExecutor
		want     string
		wantErr  bool
	}{
		{
			name:     "SMB1 enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			want:     "SMB1: enabled",
			wantErr:  false,
		},
		{
			name:     "SMB1 not enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:     "SMB1: not enabled",
			wantErr:  false,
		},
		{
			name:     "SMB2 enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			want:     "SMB2: enabled",
			wantErr:  false,
		},
		{
			name:     "SMB2 not enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:     "SMB2: not enabled",
			wantErr:  false,
		},
		{
			name:     "command error",
			executor: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("command error")},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var smbVersion string
			if strings.Contains(tt.name, "SMB1") {
				smbVersion = "SMB1"
			} else {
				smbVersion = "SMB2"
			}
			got, _, err := checks.SmbEnabled(smbVersion, tt.executor, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("SmbEnabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SmbEnabled() got = %v, want %v", got, tt.want)
			}
		})
	}
}
