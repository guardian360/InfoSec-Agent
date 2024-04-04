package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// WindowsDefender checks if the Windows Defender is enabled and if the periodic scan is enabled
//
// Parameters: _
//
// Returns: If Windows Defender and periodic scan are enabled/disabled
func WindowsDefender(scanKey registrymock.RegistryKey, defenderKey registrymock.RegistryKey) Check {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := registrymock.OpenRegistryKey(scanKey, `SOFTWARE\Microsoft\Windows Defender`)
	if err != nil {
		return NewCheckErrorf(14, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(windowsDefenderKey)

	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := registrymock.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return NewCheckErrorf(14, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(realTimeKey)

	// Read the value of the registry keys
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	if err != nil {
		return NewCheckErrorf(14, "error reading value", err)
	}
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of these keys, we can determine if Windows Defender and the periodic scan
	// are enabled or disabled
	if err != nil {
		if antiVirusPeriodic == 1 {
			return NewCheckResult(14, 1, "Windows real-time defender is enabled but the "+
				"windows periodic scan is disabled")
		}
		return NewCheckResult(14, 0, "Windows real-time defender is enabled and also the "+
			"windows periodic scan is enabled")
	}
	if realTimeDefender == 1 {
		if antiVirusPeriodic == 1 {
			return NewCheckResult(14, 3, "Windows real-time defender is disabled and also "+
				"the windows periodic scan is disabled")
		}
		return NewCheckResult(14, 2, "Windows real-time defender is disabled but the "+
			"windows periodic scan is enabled")
	}
	return NewCheckResult(14, 4, "No windows defender data found")
}
