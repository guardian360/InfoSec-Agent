package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// WindowsDefender checks if the Windows Defender is enabled and if the periodic scan is enabled
//
// Parameters: _
//
// Returns: If Windows Defender and periodic scan are enabled/disabled
func WindowsDefender(scanKey mocking.RegistryKey, defenderKey mocking.RegistryKey) Check {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := mocking.OpenRegistryKey(scanKey, `SOFTWARE\Microsoft\Windows Defender`)
	if err != nil {
		return NewCheckErrorf(WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer mocking.CloseRegistryKey(windowsDefenderKey)

	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := mocking.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return NewCheckErrorf(WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer mocking.CloseRegistryKey(realTimeKey)

	// Read the value of the registry keys
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	if err != nil {
		return NewCheckErrorf(WindowsDefenderID, "error reading value", err)
	}
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of these keys, we can determine if Windows Defender and the periodic scan
	// are enabled or disabled
	if err != nil {
		if antiVirusPeriodic == 1 {
			return NewCheckResult(WindowsDefenderID, 1, "Windows real-time defender is enabled but the "+
				"windows periodic scan is disabled")
		}
		return NewCheckResult(WindowsDefenderID, 0, "Windows real-time defender is enabled and also the "+
			"windows periodic scan is enabled")
	}
	if realTimeDefender == 1 {
		if antiVirusPeriodic == 1 {
			return NewCheckResult(WindowsDefenderID, 3, "Windows real-time defender is disabled and also "+
				"the windows periodic scan is disabled")
		}
		return NewCheckResult(WindowsDefenderID, 2, "Windows real-time defender is disabled but the "+
			"windows periodic scan is enabled")
	}
	return NewCheckResult(WindowsDefenderID, 4, "No windows defender data found")
}
