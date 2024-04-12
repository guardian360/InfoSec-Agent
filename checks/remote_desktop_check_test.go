package checks_test

import (
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

// TestRemoteDesktopCheck tests the RemoteDesktopCheck function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRemoteDesktopCheck(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
	}{
		{
			name: "Remote Desktop enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteDesktopID, 0, "Remote Desktop is enabled"),
		},
		{
			name: "Remote Desktop disabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 1}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteDesktopID, 1, "Remote Desktop is disabled"),
		},
		{
			name: "Unknown status",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 3}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteDesktopID, 1, "Remote Desktop is disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.RemoteDesktopCheck(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestRegistryOutput ensures the registry output has the expected format
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRegistryOutputRemoteDesktop(t *testing.T) {
	path := "System\\CurrentControlSet\\Control\\Terminal Server"
	expectedValueName := "fDenyTSConnections"

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	require.NoError(t, err)

	defer func(key registry.Key) {
		err = key.Close()
		require.NoError(t, err)
	}(key)

	valueNames, err := key.ReadValueNames(-1)
	require.NoError(t, err)
	var found bool
	for _, subKeyName := range valueNames {
		if subKeyName == expectedValueName {
			found = true
			break
		}
	}
	require.True(t, found, "Value name %s not found", expectedValueName)

	val, valType, err := key.GetIntegerValue(expectedValueName)
	require.NoError(t, err)
	require.Equal(t, reflect.Uint32, reflect.TypeOf(valType).Kind())
	require.True(t, val == 0 || val == 1, "Unexpected value: %v, want 0 or 1", val)
}
