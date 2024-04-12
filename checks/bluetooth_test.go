package checks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// TestBluetooth is a function that validates the functionality of the Bluetooth function with both valid and invalid inputs.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework used to run the test cases.
//
// This function defines a series of test cases, each represented as a struct. Each test case includes a name, a mock registry key to simulate the system's registry, and the expected result of the Bluetooth function when it is called with the mock registry key. The function iterates over these test cases, and for each one, it calls the Bluetooth function with the provided mock registry key and compares the result to the expected result. If the actual result matches the expected result, the test case passes; otherwise, it fails.
//
// This function does not return a value. Instead, it uses the testing framework's functionality to indicate whether each test case passed or failed.
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
			want: checks.NewCheckResult(checks.BluetoothID, 0, "No Bluetooth devices found"),
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
			want: checks.NewCheckResult(checks.BluetoothID, 1, "Device1"),
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
