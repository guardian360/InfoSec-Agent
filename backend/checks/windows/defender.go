// Package windows provides functions related to security/privacy checks of windows settings
package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// Defender is a function that checks the status of Windows Defender and its periodic scan feature on the system.
//
// Parameters:
//   - scanKey mocking.RegistryKey: A registry key object for accessing the Windows Defender registry key.
//   - defenderKey mocking.RegistryKey: A registry key object for accessing the Windows Defender Real-Time Protection registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Defender and its periodic scan feature are enabled or disabled.
//
// The function works by opening and reading the values of the Windows Defender and Real-Time Protection registry keys. Based on these values, it determines the status of Windows Defender and its periodic scan feature. The function returns a Check instance containing a string that describes the status of Windows Defender and its periodic scan feature.
func Defender(scanKey mocking.RegistryKey, defenderKey mocking.RegistryKey) checks.Check {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := checks.OpenRegistryKey(scanKey, `SOFTWARE\Microsoft\Windows Defender`)
	if err != nil {
		return checks.NewCheckErrorf(checks.WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(windowsDefenderKey)

	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := checks.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return checks.NewCheckErrorf(checks.WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(realTimeKey)

	// Read the value of the registry keys
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	if err != nil {
		return checks.NewCheckErrorf(checks.WindowsDefenderID, "error reading value", err)
	}
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of these keys, we can determine if Windows Defender and the periodic scan
	// are enabled or disabled
	if err != nil || realTimeDefender == 0 {
		if antiVirusPeriodic == 1 {
			return checks.NewCheckResult(checks.WindowsDefenderID, 1)
		}
		return checks.NewCheckResult(checks.WindowsDefenderID, 0)
	}
	if realTimeDefender == 1 {
		if antiVirusPeriodic == 1 {
			return checks.NewCheckResult(checks.WindowsDefenderID, 3)
		}
		return checks.NewCheckResult(checks.WindowsDefenderID, 2)
	}
	return checks.NewCheckResult(checks.WindowsDefenderID, 4)
}
