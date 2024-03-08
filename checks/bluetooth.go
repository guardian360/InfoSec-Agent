package checks

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Bluetooth() Check {
	// Open the registry key for Bluetooth devices
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`, registry.READ)
	if err != nil {
		return newCheckErrorf("Bluetooth", "error opening registry key", err)
	}
	defer key.Close()

	// Get the names of all sub keys (which represent Bluetooth devices)
	deviceNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return newCheckErrorf("Bluetooth", "error reading sub key names", err)
	}

	if len(deviceNames) == 0 {
		return newCheckResult("Bluetooth", "No Bluetooth devices found")
	}

	fmt.Println("Bluetooth devices:")
	result := newCheckResult("Bluetooth")
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
	bt := "null"
	output, _ := exec.Command("ipconfig").Output()
	//re := regexp.MustCompile(":")
	lines := strings.Split(string(output), "\r\n")
	for _, i := range lines[len(lines)-3:] {
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
