package devices_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestBluetooth is a unit test function for the Bluetooth function in the checks package.
//
// Parameters:
//   - t (*testing.T): A pointer to an instance of testing.T that provides methods for indicating test success or failure.
//
// This function defines a series of test cases to validate the functionality of the Bluetooth function. Each test case is represented as a struct that includes a name, a mock registry key, and the expected result of the Bluetooth function when called with the mock registry key.
//
// The function iterates over these test cases. For each test case, it calls the Bluetooth function with the provided mock registry key and compares the actual result to the expected result. If the actual result matches the expected result, the test case is considered to have passed; otherwise, it is considered to have failed.
//
// This function does not return a value. Instead, it uses the testing framework's functionality to indicate whether each test case passed or failed.
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
			got := devices.Bluetooth(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
