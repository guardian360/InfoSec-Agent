package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCredentialGuardRunning(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Credential Guard is running",
			executor: &mocking.MockCommandExecutor{
				Output: "LsaIso.exe",
			},
			want: checks.NewCheckResult(checks.CredentialGuardID, 0),
		},
		{
			name: "Credential Guard is not running",
			executor: &mocking.MockCommandExecutor{
				Output: "",
			},
			want: checks.NewCheckResult(checks.CredentialGuardID, 1),
		},
		{
			name: "Error executing command",
			executor: &mocking.MockCommandExecutor{
				Err: errors.New("error"),
			},
			want: checks.NewCheckError(checks.CredentialGuardID, errors.New("error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.CredentialGuardRunning(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
