package checks_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

func TestNetworkSharing(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Get-NetAdapterBinding command error",
			executor: &mocking.MockCommandExecutor{Output: "",
				Err: errors.New("error executing command Get-NetAdapterBinding")},
			want: checks.NewCheckErrorf(checks.NetworkSharingID,
				"error executing command Get-NetAdapterBinding",
				errors.New("error executing command Get-NetAdapterBinding")),
		},
		{
			name:     "Network sharing is enabled",
			executor: &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nTrue\r\nTrue\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 0, "Network sharing is enabled"),
		},
		{
			name:     "Network sharing is partially enabled",
			executor: &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 1, "Network sharing is partially enabled"),
		},
		{
			name:     "Network sharing is disabled",
			executor: &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 2, "Network sharing is disabled"),
		},
		{
			name:     "Network sharing status is unknown",
			executor: &mocking.MockCommandExecutor{Output: "\r\n\r\n\r\nHelloWorld\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 3, "Network sharing status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.NetworkSharing(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
