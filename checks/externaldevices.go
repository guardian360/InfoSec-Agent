package checks

import (
	"errors"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TODO: Improve formatting of output, check more classes

// ExternalDevices checks for external devices connected to the system
//
// Parameters: _
//
// Returns: list of external devices
func ExternalDevices(executorClass commandmock.CommandExecutor) Check {
	// All the classes you want to check with the Get-PnpDevice command
	classesToCheck := [2]string{"Mouse", "Camera"}
	outputs := make([]string, 0)
	for _, s := range classesToCheck {
		output, err := CheckDeviceClass(s, executorClass)

		if err != nil {
			return NewCheckErrorf("externaldevices", "error checking device "+s, err)
		}

		outputs = append(outputs, output...)
	}

	return NewCheckResult("externaldevices", outputs...)
}

// CheckDeviceClass runs a specific class within the Get-PnpDevice command
//
// Parameters: deviceClass (string) representing the class to check with the Get-PnpDevice command
//
// Returns: list of devices of the given class
func CheckDeviceClass(deviceClass string, executorClass commandmock.CommandExecutor) ([]string, error) {
	// Run the Get-PnpDevice command with the given class
	command := "powershell"
	output, err := executorClass.Execute(command, "-Command", "Get-PnpDevice -Class", deviceClass, " "+
		"| Where-Object -Property Status -eq 'OK' | Select-Object FriendlyName")

	if err != nil {
		return nil, err
	}

	// Get all devices from the output
	devices := strings.Split(string(output), "\r\n")
	if len(devices) == 1 {
		return nil, errors.New("no devices found")
	}
	devices = devices[3 : len(devices)-3]

	// Trim all spaces in devices
	for i, device := range devices {
		devices[i] = strings.TrimSpace(device)
	}

	return devices, nil
}
