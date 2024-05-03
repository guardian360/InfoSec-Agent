package checks

import (
	"errors"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// ExternalDevices is a function that conducts a security assessment for any external devices connected to the system.
//
// Parameters:
//   - executorClass (commandmock.CommandExecutor): An instance of CommandExecutor that is utilized to execute commands at the system level.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the external devices check. If any external devices are detected, their names are included in the Result field of the Check object. If an error is encountered during the check, it is encapsulated in the Error and ErrorMSG fields of the Check object.
//
// The primary use of this function is to identify potential security threats associated with external devices that are connected to the system.
func ExternalDevices(executorClass mocking.CommandExecutor) Check {
	// All the classes you want to check with the Get-PnpDevice command
	// This list can easily be extended; we refer to the Microsoft documentation for the Get-PnpDevice command
	// (for example: Biometric, Printer, etc.)
	classesToCheck := []string{"Mouse", "Camera", "AudioEndpoint", "Keyboard", "Biometric"}
	output, err := CheckDeviceClasses(classesToCheck, executorClass)

	if err != nil {
		return NewCheckErrorf(ExternalDevicesID, "error checking device", err)
	}
	return NewCheckResult(ExternalDevicesID, 0, output...)
}

// CheckDeviceClasses is a function that runs the Get-PnpDevice command for device classes.
//
// Parameters:
//   - deviceClass (string): The specific device class to be checked using the Get-PnpDevice command.
//   - executorClass (commandmock.CommandExecutor): An instance of CommandExecutor that is responsible for executing system-level commands.
//
// Returns:
//   - ([]string): A list of devices that belong to the specified device class. Each string in the list represents a device name.
//   - (error): An error object that captures any error that occurred during the command execution. If no devices are found, an error is returned.
//
// The main purpose of this function is to identify devices of a specific class that are connected to the system. It runs the Get-PnpDevice command with the specified device class and parses the output to extract the device names. If no devices are found, the function returns an error.
func CheckDeviceClasses(deviceClasses []string, executorClass mocking.CommandExecutor) ([]string, error) {
	// Convert the device classes to a string
	classesString := strings.Join(deviceClasses, ",")
	// Run the Get-PnpDevice command with the given class
	command := "powershell"
	output, err := executorClass.Execute(command, "-Command", "Get-PnpDevice -Class", classesString, " "+
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
