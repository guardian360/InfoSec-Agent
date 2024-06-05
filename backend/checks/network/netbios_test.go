package network_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.SetupTests()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNetBIOSEnabled(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.CommandExecutor
		want     checks.Check
	}{
		{
			name: "NetBIOS enabled",
			executor: &mocking.MockCommandExecutor{
				Output: "NetBIOS over Tcpip: Enabled",
			},
			want: checks.NewCheckResult(checks.NetBIOSID, 1),
		},
		{
			name: "NetBIOS disabled",
			executor: &mocking.MockCommandExecutor{
				Output: "NetBIOS over Tcpip: Disabled",
			},
			want: checks.NewCheckResult(checks.NetBIOSID, 0),
		},
		{
			name: "Error executing command",
			executor: &mocking.MockCommandExecutor{
				Err: errors.New("error"),
			},
			want: checks.NewCheckError(checks.NetBIOSID, errors.New("error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network.NetBIOSEnabled(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
