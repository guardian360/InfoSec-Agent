package checks_test

import (
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

// TestRemoteDesktopCheck is a function that tests the RemoteDesktopCheck function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the RemoteDesktopCheck function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the Windows registry. Each test case checks if the RemoteDesktopCheck function correctly identifies the status of the Remote Desktop feature based on the simulated registry key values. The function asserts that the returned Check instance contains the expected results.
func TestRemoteDesktopCheck(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "Remote Desktop enabled",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteDesktopID, 0, "Remote Desktop is enabled"),
		},
		{
			name: "Remote Desktop disabled",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 1}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.RemoteDesktopID, 1, "Remote Desktop is disabled"),
		},
		{
			name: "Unknown status",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
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

// TestRegistryOutputRemoteDesktop is a function that verifies the format and values of the registry output related to the Remote Desktop settings.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function opens the registry key for Terminal Server settings and retrieves the value names. It checks if the expected value name 'fDenyTSConnections' is present among the retrieved value names. If the expected value name is not found, the function reports a test failure. The function then retrieves the integer value and its type for 'fDenyTSConnections'. It checks if the value type is uint32 and if the value is either 0 or 1, which represent the enabled or disabled status of the Remote Desktop feature. If the value type or value does not match the expected results, the function reports a test failure.
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
