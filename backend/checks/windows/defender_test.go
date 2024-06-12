package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestWindowsDefender is a function that tests the Defender function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Defender function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the registry key access for checking the status of Windows Defender and its periodic scan feature. Each test case checks if the Defender function correctly identifies the status of Windows Defender and its periodic scan feature based on the simulated registry key values. The function asserts that the returned Check instance contains the expected results.
func TestWindowsDefender(t *testing.T) {
	tests := []struct {
		name        string
		defenderKey mocking.RegistryKey
		want        checks.Check
	}{
		{
			name: "Windows Defender disabled",
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.WindowsDefenderID, 1),
		},
		{
			name: "Windows Defender enabled",
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 0}, Err: nil}}},
			want: checks.NewCheckResult(checks.WindowsDefenderID, 0),
		},
		{
			name: "Unknown status",
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 2}, Err: nil}}},
			want: checks.NewCheckError(checks.WindowsDefenderID, errors.New("unexpected error occurred while checking Windows Defender status")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.Defender(tt.defenderKey)
			require.Equal(t, tt.want, got)
		})
	}
}
