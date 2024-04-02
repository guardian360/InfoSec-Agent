package checks_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

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
			want:      checks.NewCheckResult("smb", "SMB1: enabled", "SMB2: enabled"),
		},
		{
			name:      "SMB1 enabled and SMB2 not enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckResult("smb", "SMB1: enabled", "SMB2: not enabled"),
		},
		{
			name:      "SMB1 not enabled and SMB2 enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue", Err: nil},
			want:      checks.NewCheckResult("smb", "SMB1: not enabled", "SMB2: enabled"),
		},
		{
			name:      "SMB1 and SMB2 not enabled",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckResult("smb", "SMB1: not enabled", "SMB2: not enabled"),
		},
		{
			name:      "command smb1 error",
			executor1: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("command smb1 error")},
			executor2: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			want:      checks.NewCheckError("smb", errors.New("command smb1 error")),
		},
		{
			name:      "command smb2 error",
			executor1: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse", Err: nil},
			executor2: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("command smb2 error")},
			want:      checks.NewCheckError("smb", errors.New("command smb2 error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.SmbCheck(tt.executor1, tt.executor2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SmbCheck() = %v, want %v", got, tt.want)
			}
			if got := checks.SmbCheck(tt.executor1, tt.executor2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SmbCheck() = %v, want %v", got, tt.want)
			}
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
			if strings.Contains(tt.name, "SMB1") {
				got, err := checks.SmbEnabled("SMB1", tt.executor)
				if (err != nil) != tt.wantErr {
					t.Errorf("SmbEnabled() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("SmbEnabled() got = %v, want %v", got, tt.want)
				}
			} else {
				got, err := checks.SmbEnabled("SMB2", tt.executor)
				if (err != nil) != tt.wantErr {
					t.Errorf("SmbEnabled() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("SmbEnabled() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
