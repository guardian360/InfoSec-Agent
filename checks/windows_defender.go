// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"golang.org/x/sys/windows/registry"
)

// WindowsDefender checks if the Windows Defender is enabled and if the periodic scan is enabled
//
// Parameters: _
//
// Returns: If Windows Defender and periodic scan are enabled/disabled
func WindowsDefender() Check {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`,
		registry.READ)
	if err != nil {
		return newCheckErrorf("WindowsDefender", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer windowsDefenderKey.Close()

	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`, registry.READ)
	if err != nil {
		return newCheckErrorf("WindowsDefender", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer realTimeKey.Close()

	// Read the value of the registry keys
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of these keys, we can determine if Windows Defender and the periodic scan
	// are enabled or disabled
	if err != nil {
		if antiVirusPeriodic == 1 {
			return newCheckResult("WindowsDefender", "Windows real-time defender is enabled but the "+
				"windows periodic scan is disabled")
		} else {
			return newCheckResult("WindowsDefender", "Windows real-time defender is enabled and also the "+
				"windows periodic scan is enabled")
		}
	} else {
		if realTimeDefender == 1 {
			if antiVirusPeriodic == 1 {
				return newCheckResult("WindowsDefender", "Windows real-time defender is disabled and also "+
					"the windows periodic scan is disabled")
			} else {
				return newCheckResult("WindowsDefender", "Windows real-time defender is disabled but the "+
					"windows periodic scan is enabled")
			}
		} else {
			return newCheckResult("WindowsDefender", "No windows defender data found")
		}
	}
}
