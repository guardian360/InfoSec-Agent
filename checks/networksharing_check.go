// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"os/exec"
	"strings"
)

// NetworkSharing checks if network sharing is enabled or disabled
//
// Parameters: _
//
// Returns: If network sharing is enabled or not
func NetworkSharing() Check {
	// Execute a powershell command to get the network adapter binding status
	output, err := exec.Command("powershell", "Get-NetAdapterBinding | Where-Object "+
		"{$_.ComponentID -eq 'ms_server'} | Select-Object Enabled").Output()

	if err != nil {
		return newCheckErrorf("NetworkSharing",
			"error executing command Get-NetAdapterBinding", err)
	}

	outputString := strings.Split(string(output), "\r\n")
	counter := 0                   // Counter to keep track of the number of enabled network adapters
	total := len(outputString) - 6 // Expected number of enabled network adapters for network sharing to be enabled

	for _, line := range outputString[3 : len(outputString)-3] {
		// Check if the line contains "True" indicating network sharing is enabled for the adapter
		if strings.Contains(line, "True") {
			counter++
		}
	}

	// Check the status of network sharing based on the number of enabled network adapters
	if counter == total {
		return newCheckResult("NetworkSharing", "Network sharing is enabled")
	}
	if counter > 0 && counter < total {
		return newCheckResult("NetworkSharing", "Network sharing is partially enabled")
	}
	return newCheckResult("NetworkSharing", "Network sharing is disabled")
}
