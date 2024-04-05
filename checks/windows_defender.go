package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// WindowsDefender is a function that checks the status of Windows Defender and its periodic scan feature on the system.
//
// Parameters:
//   - scanKey registrymock.RegistryKey: A registry key object for accessing the Windows Defender registry key.
//   - defenderKey registrymock.RegistryKey: A registry key object for accessing the Windows Defender Real-Time Protection registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Defender and its periodic scan feature are enabled or disabled.
//
// The function works by opening and reading the values of the Windows Defender and Real-Time Protection registry keys. Based on these values, it determines the status of Windows Defender and its periodic scan feature. The function returns a Check instance containing a string that describes the status of Windows Defender and its periodic scan feature.
func WindowsDefender(scanKey registrymock.RegistryKey, defenderKey registrymock.RegistryKey) Check {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := registrymock.OpenRegistryKey(scanKey, `SOFTWARE\Microsoft\Windows Defender`)
	if err != nil {
		return NewCheckErrorf("WindowsDefender", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(windowsDefenderKey)

	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := registrymock.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return NewCheckErrorf("WindowsDefender", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(realTimeKey)

	// Read the value of the registry keys
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	if err != nil {
		return NewCheckErrorf("WindowsDefender", "error reading value", err)
	}
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of these keys, we can determine if Windows Defender and the periodic scan
	// are enabled or disabled
	if err != nil {
		if antiVirusPeriodic == 1 {
			return NewCheckResult("WindowsDefender", "Windows real-time defender is enabled but the "+
				"windows periodic scan is disabled")
		}
		return NewCheckResult("WindowsDefender", "Windows real-time defender is enabled and also the "+
			"windows periodic scan is enabled")
	}
	if realTimeDefender == 1 {
		if antiVirusPeriodic == 1 {
			return NewCheckResult("WindowsDefender", "Windows real-time defender is disabled and also "+
				"the windows periodic scan is disabled")
		}
		return NewCheckResult("WindowsDefender", "Windows real-time defender is disabled but the "+
			"windows periodic scan is enabled")
	}
	return NewCheckResult("WindowsDefender", "No windows defender data found")
}
