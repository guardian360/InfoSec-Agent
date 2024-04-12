package checks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TestBluetooth tests the Bluetooth function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestBluetooth(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
	}{
		{
			name: "No Devices found",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices"}}},
			want: checks.NewCheckResult(checks.BluetoothID, 0, "No Bluetooth devices found"),
		},
		{
			name: "Bluetooth devices found",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "4dbndas2", BinaryValues: map[string][]byte{"Name": []byte("Device1")}, Err: nil}},
					},
				}, Err: nil},
			want: checks.NewCheckResult(checks.BluetoothID, 1, "Device1"),
		},
		{
			name: "Error reading device name",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "FAFA", StringValues: map[string]string{"Name2": "fsdfs"}, Err: nil}},
					},
				}},
			want: checks.NewCheckResult(checks.BluetoothID, 1, "Error reading device name FAFA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.Bluetooth(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
