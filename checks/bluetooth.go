package checks

import (
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"strings"
)

// Bluetooth checks for bluetooth devices which are / have been connected to the system
//
// Parameters: _
//
// Returns: A list of bluetooth devices
func Bluetooth() Check {
	// Open the registry key for bluetooth devices
	key, err := utils.OpenRegistryKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`)
	if err != nil {
		return NewCheckErrorf("Bluetooth", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer utils.CloseRegistryKey(key)

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
		deviceKey, err := utils.OpenRegistryKey(key, deviceName)
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error opening device subkey %s", deviceName))
			continue
		}
		defer utils.CloseRegistryKey(deviceKey)

		// Get the device name
		deviceNameValue, _, err := deviceKey.GetBinaryValue("Name")
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error reading device name %s", deviceName))
		} else {
			result.Result = append(result.Result, string(deviceNameValue))
		}

	}

	// Check for currently connected bluetooth devices
	bt := "null"

	output, _ := exec.Command("ipconfig").Output()

	//re := regexp.MustCompile(":")
	lines := strings.Split(string(output), "\r\n")
	for _, i := range lines[len(lines)-3:] {
		// Within ipconfig, Media State represents the status of the bluetooth adapter
		if strings.Contains(i, "Media State") {
			bt = i
		}
	}
	if bt != "null" {
		btLines := strings.Split(bt, ": ")
		if btLines[1] == "Media disconnected" {
			result.Result = append(result.Result, "You have no devices connected via bluetooth.")
		} else {
			result.Result = append(result.Result, fmt.Sprintf("Device connected: %s", btLines[1]))
		}
	}

	return result
}
