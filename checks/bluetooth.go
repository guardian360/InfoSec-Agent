package checks

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// Bluetooth checks for bluetooth devices which are / have been connected to the system
//
// Parameters: _
//
// Returns: A list of bluetooth devices
func Bluetooth(registryKey registrymock.RegistryKey) Check {
	// Open the registry key for bluetooth devices
	key, err := registrymock.OpenRegistryKey(registrymock.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`)
	if err != nil {
		return NewCheckErrorf("Bluetooth", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Get the names of all sub keys (which represent bluetooth devices)
	deviceNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return NewCheckErrorf("Bluetooth", "error reading sub key names", err)
	}

	if len(deviceNames) == 0 {
		return NewCheckResult("Bluetooth", "No Bluetooth devices found")
	}

	result := NewCheckResult("Bluetooth")
	// Open each device sub key within the registry
	for _, deviceName := range deviceNames {
		deviceKey, err := registrymock.OpenRegistryKey(key, deviceName)
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error opening device subkey %s", deviceName))
			continue
		}
		defer registrymock.CloseRegistryKey(deviceKey)

		// Get the device name
		deviceNameValue, _, err := deviceKey.GetBinaryValue("Name")
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error reading device name %s", deviceName))
		} else {
			result.Result = append(result.Result, string(deviceNameValue))
		}
	}
	return result
}
