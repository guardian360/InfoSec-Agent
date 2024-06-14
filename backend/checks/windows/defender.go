// Package windows provides functions related to security/privacy checks of windows settings
package windows

import (
	"errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// Defender is a function that checks the status of the Windows Defender.
//
// Parameters:
//   - defenderKey mocking.RegistryKey: A registry key object for accessing the Windows Defender Real-Time Protection registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Defender and its periodic scan feature are enabled or disabled.
//
// The function works by opening and reading the value of the Real-Time Protection registry key. Based on this value, it determines the status of the Windows Defender.
func Defender(defenderKey mocking.RegistryKey) checks.Check {
	// Open the Windows Defender real-time protection registry key
	realTimeKey, err := checks.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return checks.NewCheckErrorf(checks.WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(realTimeKey)

	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")

	// Based on the values of this key, determine the status of Windows Defender
	if err != nil || realTimeDefender == 0 {
		return checks.NewCheckResult(checks.WindowsDefenderID, 0)
	}
	if realTimeDefender == 1 {
		return checks.NewCheckResult(checks.WindowsDefenderID, 1)
	}
	return checks.NewCheckError(checks.WindowsDefenderID, errors.New("unexpected error occurred while checking Windows Defender status"))
}
