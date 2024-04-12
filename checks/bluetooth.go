package checks

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// Bluetooth is a function that checks for Bluetooth devices which are currently connected or have been previously connected to the system.
//
// Parameters:
//   - registryKey (registrymock.RegistryKey): The registry key used to access the system's registry.
//
// Returns:
//   - Check: A Check object that encapsulates the results of the Bluetooth check. The Check object includes a list of strings, where each string represents a Bluetooth device that is currently or was previously connected to the system. If an error occurs during the Bluetooth check, the Check object will encapsulate this error.
//
// This function first opens the registry key for Bluetooth devices. It then reads the names of all subkeys, which represent Bluetooth devices. For each device, the function opens the device subkey, retrieves the device name, and adds it to the results. If an error occurs at any point during this process, it is encapsulated in the Check object and returned.
func Bluetooth(registryKey registrymock.RegistryKey) Check {
	var err error
	var deviceKey registrymock.RegistryKey
	var deviceNames []string
	var deviceNameValue []byte
	// Open the registry key for bluetooth devices
	key, err := registrymock.OpenRegistryKey(registryKey,
		`SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`)
	if err != nil {
		return NewCheckError(BluetoothID, err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Get the names of all sub keys (which represent bluetooth devices)
	deviceNames, err = key.ReadSubKeyNames(-1)
	if err != nil {
		return NewCheckErrorf(BluetoothID, "error reading sub key names", err)
	}

	if len(deviceNames) == 0 {
		return NewCheckResult(BluetoothID, 0, "No Bluetooth devices found")
	}

	result := NewCheckResult(BluetoothID, 1)
	// Open each device sub key within the registry
	for _, deviceName := range deviceNames {
		deviceKey, err = registrymock.OpenRegistryKey(key, deviceName)
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error opening device subkey %s", deviceName))
			continue
		}
		defer registrymock.CloseRegistryKey(deviceKey)

		// Get the device name
		deviceNameValue, _, err = deviceKey.GetBinaryValue("Name")
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error reading device name %s", deviceName))
		} else {
			result.Result = append(result.Result, string(deviceNameValue))
		}
	}
	return result
}
