// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"golang.org/x/sys/windows/registry"
)

// RemoteDesktopCheck checks if Remote Desktop is enabled
//
// Parameters: _
//
// Returns: If Remote Desktop is enabled or not
func RemoteDesktopCheck() Check {
	// Open the registry key for Terminal Server settings
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Terminal Server`, registry.QUERY_VALUE)
	if err != nil {
		return newCheckErrorf("RemoteDesktop", "error opening registry key", err)
	}
	defer key.Close()

	// Read the value of fDenyTSConnections, which contains the information if Remote Desktop is enabled or not
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return newCheckErrorf("RemoteDesktop", "error reading fDenyTSConnections", err)
	} else {
		// Check if Remote Desktop is enabled or disabled based on the value of fDenyTSConnections
		if val == 0 {
			return newCheckResult("RemoteDesktop", "Remote Desktop is enabled")
		} else {
			return newCheckResult("RemoteDesktop", "Remote Desktop is disabled")
		}
	}

}
