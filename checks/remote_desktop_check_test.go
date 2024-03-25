package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"reflect"
	"testing"
)

// TestRemoteDesktopCheck tests the RemoteDesktopCheck function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRemoteDesktopCheck(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "Remote Desktop enabled",
			key: &registrymock.MockRegistryKey{StringValue: "fDenyTSConnections",
				BinaryValue: nil, IntegerValue: 0, Err: nil},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is enabled"),
		},
		{
			name: "Remote Desktop disabled",
			key: &registrymock.MockRegistryKey{StringValue: "fDenyTSConnections",
				BinaryValue: nil, IntegerValue: 1, Err: nil},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is disabled"),
		},
		{
			name: "Unknown status",
			key: &registrymock.MockRegistryKey{StringValue: "fDenyTSConnections",
				BinaryValue: nil, IntegerValue: 3, Err: nil},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.RemoteDesktopCheck(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteDesktopCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
