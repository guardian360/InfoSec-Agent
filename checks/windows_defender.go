package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// WindowsDefender is a function that checks the status of Windows Defender and its periodic scan feature on the system.
//
// Parameters:
//   - scanKey mocking.RegistryKey: A registry key object for accessing the Windows Defender registry key.
//   - defenderKey mocking.RegistryKey: A registry key object for accessing the Windows Defender Real-Time Protection registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Defender and its periodic scan feature are enabled or disabled.
//
// The function works by opening and reading the values of the Windows Defender and Real-Time Protection registry keys. Based on these values, it determines the status of Windows Defender and its periodic scan feature. The function returns a Check instance containing a string that describes the status of Windows Defender and its periodic scan feature.
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
