// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"os/exec"
	"strings"
	"syscall"
)

// TODO: Improve formatting of output, check more classes

// ExternalDevices checks for external devices connected to the system
//
// Parameters: _
//
// Returns: list of external devices
func ExternalDevices() Check {
	// All the classes you want to check with the Get-PnpDevice command
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

// checkDeviceClass runs a specific class within the Get-PnpDevice command
//
// Parameters: deviceClass (string) representing the class to check with the Get-PnpDevice command
//
// Returns: list of devices of the given class
func checkDeviceClass(deviceClass string) ([]string, error) {
	// Run the Get-PnpDevice command with the given class
	cmd := exec.Command("powershell", "-Command", "Get-PnpDevice -Class", deviceClass, " "+
		"| Where-Object -Property Status -eq 'OK' | Select-Object FriendlyName")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()

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
