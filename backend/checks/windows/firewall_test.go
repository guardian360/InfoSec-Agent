package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFirewallEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Firewall enabled",
			executor: &mocking.MockCommandExecutor{
				Output: "State ON\n State ON\n State ON\n",
			},
			want: checks.NewCheckResult(checks.FirewallID, 0),
		},
		{
			name: "Firewall disabled",
			executor: &mocking.MockCommandExecutor{
				Output: "State OFF\n State OFF\n State OFF\n",
			},
			want: checks.NewCheckResult(checks.FirewallID, 1),
		},
		{
			name: "Firewall partially enabled",
			executor: &mocking.MockCommandExecutor{
				Output: "State ON\n State OFF\n State ON\n",
			},
			want: checks.NewCheckResult(checks.FirewallID, 1),
		},
		{
			name: "Error executing command",
			executor: &mocking.MockCommandExecutor{
				Err: errors.New("error"),
			},
			want: checks.NewCheckError(checks.FirewallID, errors.New("error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.FirewallEnabled(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
