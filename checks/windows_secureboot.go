// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"golang.org/x/sys/windows/registry"
)

// SecureBoot checks if Windows secure boot is enabled
//
// Parameters: _
//
// Returns: If Windows secure boot is enabled or not
func SecureBoot() Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\SecureBoot\State`, registry.READ)

	if err != nil {
		return newCheckError("SecureBoot", err)
	}

	defer func(windowsSecureBoot registry.Key) {
		err := windowsSecureBoot.Close()
		if err != nil {
			return
		}
	}(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return newCheckError("SecureBoot", err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return newCheckResult("SecureBoot", "Secure boot is enabled")
	}

	return newCheckResult("SecureBoot", "Secure boot is disabled")
}
