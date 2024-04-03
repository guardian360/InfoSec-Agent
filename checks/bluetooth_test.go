package checks_test

import (
	"github.com/stretchr/testify/require"
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
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices"}}},
			want: checks.NewCheckResult("Bluetooth", "No Bluetooth devices found"),
		},
		{
			name: "Bluetooth devices found",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "4dbndas2", BinaryValues: map[string][]byte{"Name": []byte("Device1")}, Err: nil}},
					},
				}, Err: nil},
			want: checks.NewCheckResult("Bluetooth", "Device1"),
		},
		{
			name: "Error reading device name",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "FAFA", StringValues: map[string]string{"Name2": "fsdfs"}, Err: nil}},
					},
				}},
			want: checks.NewCheckResult("Bluetooth", "Error reading device name FAFA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.Bluetooth(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
