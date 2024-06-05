package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"strconv"
)

// ScreenLockEnabled is a function that checks if the screen lock is enabled on the system and set to secure settings.
// The function reads the registry key for the screen lock settings and checks if the screen saver is active, requires a log-in, and has a timeout of 2 minutes or less.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): A mocker of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the screen lock is enabled and set to secure settings or an error, if one occurred.
func ScreenLockEnabled(registryKey mocking.RegistryKey) checks.Check {
	key, err := checks.OpenRegistryKey(registryKey, `Control Panel\Desktop`)

	if err != nil {
		return checks.NewCheckErrorf(checks.ScreenLockID, "error opening screen lock registry key", err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(key)

	// Read the values of ScreenSaveActive, ScreenSaverIsSecure, and ScreenSaveTimeOut
	// ScreenSaveActive indicates if the screen saver is enabled (1) or disabled (0)
	// ScreenSaverIsSecure indicates if the screen saver requires a password to unlock (1) or not (0)
	// ScreenSaveTimeOut indicates the time in minutes before the screen saver activates
	ssOn, _, err := key.GetStringValue("ScreenSaveActive")
	if err != nil {
		return checks.NewCheckErrorf(checks.ScreenLockID, "error reading ScreenSaveActive", err)
	}
	ssSecure, _, err := key.GetStringValue("ScreenSaverIsSecure")
	if err != nil {
		return checks.NewCheckErrorf(checks.ScreenLockID, "error reading ScreenSaverIsSecure", err)
	}
	ssInterval, _, err := key.GetStringValue("ScreenSaveTimeOut")
	if err != nil {
		return checks.NewCheckErrorf(checks.ScreenLockID, "error reading ScreenSaveTimeOut", err)
	}
	ssIntInterval, err := strconv.Atoi(ssInterval)
	if err != nil {
		return checks.NewCheckErrorf(checks.ScreenLockID, "error converting ScreenSaveTimeOut to int", err)
	}

	// Check that the screen saver is enabled, secure, and has a timeout of 2 minutes or less
	if ssOn == "1" && ssSecure == "1" && ssIntInterval <= 120 {
		return checks.NewCheckResult(checks.ScreenLockID, 0)
	}
	return checks.NewCheckResult(checks.ScreenLockID, 1)
}
