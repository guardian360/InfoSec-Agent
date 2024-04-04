package checks_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

func TestNetworkSharing(t *testing.T) {
	tests := []struct {
		name     string
		executor commandmock.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Get-NetAdapterBinding command error",
			executor: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("error executing command Get-NetAdapterBinding")},
			want: checks.NewCheckErrorf(3,
				"error executing command Get-NetAdapterBinding",
				errors.New("error executing command Get-NetAdapterBinding")),
		},
		{
			name:     "Network sharing is enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nTrue\r\nTrue\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(3, 0, "Network sharing is enabled"),
		},
		{
			name:     "Network sharing is partially enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(3, 1, "Network sharing is partially enabled"),
		},
		{
			name:     "Network sharing is disabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(3, 2, "Network sharing is disabled"),
		},
		{
			name:     "Network sharing status is unknown",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nHelloWorld\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(3, 3, "Network sharing status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.NetworkSharing(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
