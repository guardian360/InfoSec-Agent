package checks

import (
	"errors"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TODO: Improve formatting of output, check more classes

// ExternalDevices performs a security check for external devices connected to the system.
//
// Parameters:
//   - executorClass (commandmock.CommandExecutor): An instance of CommandExecutor used to execute system commands.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the external devices check. If any external devices are found, their names are included in the Result field of the Check instance. If an error occurs during the check, it is encapsulated in the Error and ErrorMSG fields of the Check instance.
//
// This function is primarily used to identify potential security risks associated with external devices connected to the system.
func ExternalDevices(executorClass commandmock.CommandExecutor) Check {
	// All the classes you want to check with the Get-PnpDevice command
	classesToCheck := [2]string{"Mouse", "Camera"}
	outputs := make([]string, 0)
	for _, s := range classesToCheck {
		output, err := CheckDeviceClass(s, executorClass)

		if err != nil {
			return NewCheckErrorf(ExternalDevicesID, "error checking device "+s, err)
		}

		outputs = append(outputs, output...)
	}

	return NewCheckResult(ExternalDevicesID, 0, outputs...)
}

// CheckDeviceClass executes the Get-PnpDevice command for a specific device class.
//
// Parameters:
//   - deviceClass (string): The device class to check with the Get-PnpDevice command.
//   - executorClass (commandmock.CommandExecutor): An instance of CommandExecutor used to execute system commands.
//
// Returns:
//   - ([]string): A list of devices belonging to the specified device class. Each string represents a device name.
//   - (error): An error object that captures any error that occurred during the execution of the command. If no devices are found, an error is returned.
//
// This function is primarily used to identify devices of a specific class connected to the system. It executes the Get-PnpDevice command with the specified device class and parses the output to extract the device names. If no devices are found, it returns an error.
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
