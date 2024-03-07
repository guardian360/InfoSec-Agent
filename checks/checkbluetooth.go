package checks

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func bluetooth() {
	// Open the registry key for Bluetooth devices
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\BTHPORT\Parameters\Devices`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	// Get the names of all subkeys (which represent Bluetooth devices)
	deviceNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkey names:", err)
		return
	}

	if len(deviceNames) == 0 {
		fmt.Println("No Bluetooth devices found.")
		return
	}

	fmt.Println("Bluetooth devices:")
	for _, deviceName := range deviceNames {
		deviceKey, err := registry.OpenKey(key, deviceName, registry.READ)
		if err != nil {
			fmt.Println("Error opening device subkey:", err)
			continue
		}

		// Get the device name
		deviceNameValue, _, err := deviceKey.GetBinaryValue("Name")
		if err != nil {
			fmt.Println("Error reading device name:", err)
		} else {
			fmt.Println("-", string(deviceNameValue))
		}

		deviceKey.Close()
	}
	var bt string = "null"
	output, _ := exec.Command("ipconfig").Output()
	//re := regexp.MustCompile(":")
	lines := strings.Split(string(output), "\r\n")
	for _, i := range lines[len(lines)-3:] {
		if strings.Contains(i, "Media State") {
			bt = i
		}
	}
	if bt != "null" {
		btlines := strings.Split(bt, ": ")
		if btlines[1] == "Media disconnected" {
			fmt.Println("You have no devices connected via bluetooth.")
		} else {
			fmt.Println(btlines[1])
		}
	}
}
