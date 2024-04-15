package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TestWindowsDefender is a function that tests the WindowsDefender function's behavior with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the WindowsDefender function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the registry key access for checking the status of Windows Defender and its periodic scan feature. Each test case checks if the WindowsDefender function correctly identifies the status of Windows Defender and its periodic scan feature based on the simulated registry key values. The function asserts that the returned Check instance contains the expected results.
func TestWindowsDefender(t *testing.T) {
	tests := []struct {
		name        string
		scanKey     mocking.RegistryKey
		defenderKey mocking.RegistryKey
		want        checks.Check
	}{
		{
			name: "Windows Defender disabled and periodic scan disabled",
			scanKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 1}, Err: nil}}},
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.WindowsDefenderID, 3,
				"Windows real-time defender is disabled and also the windows periodic scan is disabled"),
		},
		{
			name: "Windows Defender disabled and periodic scan enabled",
			scanKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 0}, Err: nil}}},
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.WindowsDefenderID, 2,
				"Windows real-time defender is disabled but the windows periodic scan is enabled"),
		},
		{
			name: "Unknown status",
			scanKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 0}, Err: nil}}},
			defenderKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 0}, Err: nil}}},
			want: checks.NewCheckResult(checks.WindowsDefenderID, 4, "No windows defender data found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.WindowsDefender(tt.scanKey, tt.defenderKey)
			require.Equal(t, tt.want, got)
		})
	}
}
