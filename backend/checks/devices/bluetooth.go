// Package devices provides functions related to security/privacy checks of (external) devices
package devices

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// Bluetooth is a function that checks for Bluetooth devices which are currently connected or have been previously connected to the system.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//
// Returns:
//   - Check: A Check object that encapsulates the results of the Bluetooth check. The Check object includes a list of strings, where each string represents a Bluetooth device that is currently or was previously connected to the system. If an error occurs during the Bluetooth check, the Check object will encapsulate this error.
//
// This function first opens the registry key for Bluetooth devices. It then reads the names of all sub-keys, which represent Bluetooth devices. For each device, the function opens the device sub-key, retrieves the device name, and adds it to the results. If an error occurs at any point during this process, it is encapsulated in the Check object and returned.
func Bluetooth(registryKey mocking.RegistryKey) checks.Check {
	var err error
	var deviceKey mocking.RegistryKey
	var deviceNames []string
	var deviceNameValue []byte
	// Open the registry key for bluetooth devices
	key, err := checks.OpenRegistryKey(registryKey,
		`SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`)
	if err != nil {
		return checks.NewCheckError(checks.BluetoothID, err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(key)

	// Get the names of all sub keys (which represent bluetooth devices)
	deviceNames, err = key.ReadSubKeyNames(-1)
	if err != nil {
		return checks.NewCheckErrorf(checks.BluetoothID, "error reading sub key names", err)
	}

	if len(deviceNames) == 0 {
		return checks.NewCheckResult(checks.BluetoothID, 0)
	}

	result := checks.NewCheckResult(checks.BluetoothID, 1)
	// Open each device sub key within the registry
	for _, deviceName := range deviceNames {
		deviceKey, err = checks.OpenRegistryKey(key, deviceName)
		if err != nil {
			logger.Log.Error("Error opening device subkey " + deviceName)
			continue
		}

		defer checks.CloseRegistryKey(deviceKey)

		// Get the device name
		deviceNameValue, _, err = deviceKey.GetBinaryValue("Name")
		if err != nil {
			logger.Log.Error("Error reading device name " + deviceName)
		} else {
			result.Result = append(result.Result, string(deviceNameValue))
		}
	}
	return result
}
