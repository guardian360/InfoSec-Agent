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
// The function works by opening and reading the values of the Windows Defender and Real-Time Protection registry keys. Based on these values, it determines the status of Windows Defender and its periodic scan feature. The function returns a Check instance containing a string that describes the status of Windows Defender and its periodic scan feature.
func Defender(defenderKey mocking.RegistryKey) checks.Check {
	// Open the Windows Defender real-time protection registry key, representing the periodic scan
	realTimeKey, err := mocking.OpenRegistryKey(defenderKey,
		`SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return checks.NewCheckErrorf(checks.WindowsDefenderID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer mocking.CloseRegistryKey(realTimeKey)

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
