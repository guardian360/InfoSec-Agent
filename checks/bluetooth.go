// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Bluetooth checks for bluetooth devices which are / have been connected to the system
//
// Parameters: _
//
// Returns: A list of bluetooth devices
func Bluetooth() Check {
	// Open the registry key for bluetooth devices
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`, registry.READ)
	if err != nil {
		return newCheckErrorf("Bluetooth", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer key.Close()

	// Get the names of all sub keys (which represent bluetooth devices)
	deviceNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return newCheckErrorf("Bluetooth", "error reading sub key names", err)
	}

	if len(deviceNames) == 0 {
		return newCheckResult("Bluetooth", "No Bluetooth devices found")
	}

	result := newCheckResult("Bluetooth")
	// Open each device sub key within the registry
	for _, deviceName := range deviceNames {
		deviceKey, err := registry.OpenKey(key, deviceName, registry.READ)
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error opening device subkey %s", deviceName))
			continue
		}

		// Get the device name
		deviceNameValue, _, err := deviceKey.GetBinaryValue("Name")
		if err != nil {
			result.Result = append(result.Result, fmt.Sprintf("Error reading device name %s", deviceName))
		} else {
			result.Result = append(result.Result, string(deviceNameValue))
		}

		deviceKey.Close()
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
