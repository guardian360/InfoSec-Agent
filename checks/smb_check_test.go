package checks_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"testing"
)

func TestSmbEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor *utils.MockCommandExecutor
		want     string
		wantErr  bool
	}{
		{
			name:     "SMB1 enabled",
			executor: &utils.MockCommandExecutor{Output: "True", Err: nil},
			want:     "SMB1: enabled",
			wantErr:  false,
		},
		{
			name:     "SMB1 not enabled",
			executor: &utils.MockCommandExecutor{Output: "False", Err: nil},
			want:     "SMB1: not enabled",
			wantErr:  false,
		},
		{
			name:     "SMB2 enabled",
			executor: &utils.MockCommandExecutor{Output: "True", Err: nil},
			want:     "SMB2: enabled",
			wantErr:  false,
		},
		{
			name:     "SMB2 not enabled",
			executor: &utils.MockCommandExecutor{Output: "False", Err: nil},
			want:     "SMB2: not enabled",
			wantErr:  false,
		},
		{
			name:     "command error",
			executor: &utils.MockCommandExecutor{Output: "", Err: errors.New("command error")},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checks.SmbEnabled(tt.name[:3], tt.executor)
			if (err != nil) != tt.wantErr {
				t.Errorf("smbEnabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("smbEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
