package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestPasswordLength(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Password complexity is less than 15",
			executor: &mocking.MockCommandExecutor{
				Output: "\n\n\n:0\n",
			},
			want: checks.NewCheckResult(checks.PasswordComplexityID, 1),
		},
		{
			name: "Password complexity is 15",
			executor: &mocking.MockCommandExecutor{
				Output: "\n\n\n:15\n",
			},
			want: checks.NewCheckResult(checks.PasswordComplexityID, 0),
		},
		{
			name: "Error executing command",
			executor: &mocking.MockCommandExecutor{
				Err: errors.New("error"),
			},
			want: checks.NewCheckError(checks.PasswordComplexityID, errors.New("error")),
		},
		{
			name: "Command output is empty",
			executor: &mocking.MockCommandExecutor{
				Output: "",
			},
			want: checks.NewCheckError(checks.PasswordComplexityID, errors.New("command output does not have expected structure")),
		},
		{
			name: "Error parsing password length",
			executor: &mocking.MockCommandExecutor{
				Output: "\n\n\n:abc\n",
			},
			want: checks.NewCheckError(checks.PasswordComplexityID, &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.PasswordLength(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
