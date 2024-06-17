package devices_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

func TestMain(m *testing.M) {
	logger.SetupTests()

	exitCode := m.Run()

	os.Exit(exitCode)
}

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
		err  bool
	}{
		{
			name: "No Devices found",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices"}}},
			want: checks.NewCheckResult(checks.BluetoothID, 0),
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
			want: checks.NewCheckResult(checks.BluetoothID, 1),
		},
		{
			name: "Error opening registry key",
			key:  &mocking.MockRegistryKey{},
			err:  true,
		},
		{
			name: "Error reading sub key names",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "4dbndas2", StringValues: map[string]string{"test": "test"},
								BinaryValues: map[string][]byte{"Name": []byte("Device1")}, Err: nil}},
					},
				}, Err: nil},
			want: checks.NewCheckError(checks.BluetoothID, errors.New("error")),
			err:  true,
		},
		{
			name: "Error opening sub key",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SYSTEM\\CurrentControlSet\\Services\\BTHPORT\\Parameters\\Devices",
						SubKeys: []mocking.MockRegistryKey{
							{KeyName: "4dbndas2", IntegerValues: map[string]uint64{"test": 1},
								BinaryValues: map[string][]byte{"Name": []byte("Device1")}, Err: nil}},
					},
				}, Err: nil},
			want: checks.NewCheckResult(checks.BluetoothID, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := devices.Bluetooth(tt.key)
			if tt.err {
				require.Error(t, got.Error)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
