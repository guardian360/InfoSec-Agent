package checks_test

import (
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// TestBluetooth tests the Bluetooth function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestBluetooth(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "No Devices found",
			key:  &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{}, Err: nil},
			want: checks.NewCheckResult("Bluetooth", "No Bluetooth devices found"),
		},
		{
			name: "One Device Connected",
			key: &registrymock.MockRegistryKey{Name: "ConnectedDevice1",
				BinaryValue: []byte{109, 97, 110, 111, 110}, SubKeys: []registrymock.MockRegistryKey{
					{
						BinaryValue: []byte{109, 97, 12, 111, 110, 0},
						SubKeys:     []registrymock.MockRegistryKey{},
					},
				},
				Err: nil,
			},
			want: checks.NewCheckResult("Bluetooth", "manon"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.Bluetooth(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bluetooth() = %v, want %v", got, tt.want)
			}
		})
	}
}
