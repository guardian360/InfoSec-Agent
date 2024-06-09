package network_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWPADEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "WPAD enabled",
			executor: &mocking.MockCommandExecutor{
				Output: "STATE RUNNING",
			},
			want: checks.NewCheckResult(checks.WPADID, 1),
		},
		{
			name: "WPAD disabled",
			executor: &mocking.MockCommandExecutor{
				Output: "STATE STOPPED",
			},
			want: checks.NewCheckResult(checks.WPADID, 0),
		},
		{
			name: "Error executing command",
			executor: &mocking.MockCommandExecutor{
				Err: errors.New("error"),
			},
			want: checks.NewCheckError(checks.WPADID, errors.New("error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network.WPADEnabled(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
