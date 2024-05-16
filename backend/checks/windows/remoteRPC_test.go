package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAllowRemoteRPC(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
		err  bool
	}{
		{
			name: "Remote RPC enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"AllowRemoteRPC": 1}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteRPCID, 1),
		},
		{
			name: "Automatic log-in disabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"AllowRemoteRPC": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteRPCID, 0),
		},
		{
			name: "Error opening registry key",
			key:  &mocking.MockRegistryKey{},
			want: checks.NewCheckError(checks.RemoteRPCID, errors.New("error opening registry key: key not found")),
			err:  true,
		},
		{
			name: "Error reading AllowRemoteRPC value",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"AllowRemoteRPC2": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteRPCID, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.AllowRemoteRPC(tt.key)
			if tt.err {
				require.Error(t, got.Error)
			} else {
				require.Equal(t, tt.want.ResultID, got.ResultID)
			}
		})
	}
}
