package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"reflect"
	"testing"
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
			name:        "Windows Defender disabled and periodic scan disabled",
			scanKey:     &registrymock.MockRegistryKey{StringValue: "DisableAntiVirus", BinaryValue: nil, IntegerValue: 1, Err: nil},
			defenderKey: &registrymock.MockRegistryKey{StringValue: "DisableRealtimeMonitoring", BinaryValue: nil, IntegerValue: 1, Err: nil},
			want:        checks.NewCheckResult("WindowsDefender", "Windows real-time defender is disabled and also the windows periodic scan is disabled"),
		},
		{
			name:        "Windows Defender disabled and periodic scan enabled",
			scanKey:     &registrymock.MockRegistryKey{StringValue: "DisableAntiVirus", BinaryValue: nil, IntegerValue: 0, Err: nil},
			defenderKey: &registrymock.MockRegistryKey{StringValue: "DisableRealtimeMonitoring", BinaryValue: nil, IntegerValue: 1, Err: nil},
			want:        checks.NewCheckResult("WindowsDefender", "Windows real-time defender is disabled but the windows periodic scan is enabled"),
		},
		{
			name:        "Unknown status",
			scanKey:     &registrymock.MockRegistryKey{StringValue: "DisableAntiVirus", BinaryValue: nil, IntegerValue: 0, Err: nil},
			defenderKey: &registrymock.MockRegistryKey{StringValue: "DisableRealtimeMonitoring", BinaryValue: nil, IntegerValue: 0, Err: nil},
			want:        checks.NewCheckResult("WindowsDefender", "No windows defender data found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.WindowsDefender(tt.scanKey, tt.defenderKey)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WindowsDefender() = %v, want %v", got, tt.want)
			}
		})
	}
}
