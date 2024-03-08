package checks

import (
	"os/exec"
	"strings"
)

// TODO: Improve formatting of output, Check more classes
func ExternalDevices() Check {
	// All the classes you want to check within the Get-PnpDevice command
	classesToCheck := [2]string{"Mouse", "Camera"}
	outputs := make([]string, 0)
	for _, s := range classesToCheck {
		output, err := checkDeviceClass(s)

		if err != nil {
			return newCheckErrorf("externaldevices", "error checking device "+s, err)
		}

		outputs = append(outputs, output...)
	}

	return newCheckResult("externaldevices", outputs...)
}

// Run the command for a specific class within the Get-PnpDevice and print its results
func checkDeviceClass(deviceClass string) ([]string, error) {
	output, err := exec.Command("powershell", "-Command", "Get-PnpDevice -Class", deviceClass, " | Where-Object -Property Status -eq 'OK' | Select-Object FriendlyName").Output()

	if err != nil {
		return nil, err
	}

	// Get all devices from the output
	devices := strings.Split(string(output), "\r\n")
	devices = devices[3 : len(devices)-3]

	// Trim all spaces in devices
	for i, device := range devices {
		devices[i] = strings.TrimSpace(device)
	}

	return devices, nil
}
