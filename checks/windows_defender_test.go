package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// TestWindowsDefender tests the WindowsDefender function with (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestWindowsDefender(t *testing.T) {
	tests := []struct {
		name        string
		scanKey     registrymock.RegistryKey
		defenderKey registrymock.RegistryKey
		want        checks.Check
	}{
		{
			name: "Windows Defender disabled and periodic scan disabled",
			scanKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 1}, Err: nil}}},
			defenderKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 1}, Err: nil}}},
			want: checks.NewCheckResult("WindowsDefender",
				"Windows real-time defender is disabled and also the windows periodic scan is disabled"),
		},
		{
			name: "Windows Defender disabled and periodic scan enabled",
			scanKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 0}, Err: nil}}},
			defenderKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 1}, Err: nil}}},
			want: checks.NewCheckResult("WindowsDefender",
				"Windows real-time defender is disabled but the windows periodic scan is enabled"),
		},
		{
			name: "Unknown status",
			scanKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender",
				IntegerValues: map[string]uint64{"DisableAntiVirus": 0}, Err: nil}}},
			defenderKey: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SOFTWARE\\Microsoft\\Windows Defender\\Real-Time Protection",
				IntegerValues: map[string]uint64{"DisableRealtimeMonitoring": 0}, Err: nil}}},
			want: checks.NewCheckResult("WindowsDefender", "No windows defender data found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.WindowsDefender(tt.scanKey, tt.defenderKey)
			require.Equal(t, tt.want, got)
		})
	}
}
