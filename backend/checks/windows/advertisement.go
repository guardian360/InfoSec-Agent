// Package windows provides functions related to security/privacy checks of windows settings
package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// Advertisement is a function that checks if the Advertisement ID is set to be shared with apps.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): A mocker of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Advertisement ID is shared with apps or not.
//
// The function works by opening and reading the value of the AdvertisingInfo registry key.
// Based on this value, it determines if the Advertisement ID is shared with apps.
// The function returns a Check instance containing a string that describes the status of the Advertisement ID.
func Advertisement(registryKey mocking.RegistryKey) checks.Check {
	key, err := mocking.OpenRegistryKey(registryKey, `SOFTWARE\Microsoft\Windows\CurrentVersion\AdvertisingInfo`)
	if err != nil {
		return checks.NewCheckError(checks.AdvertisementID, err)
	}
	defer mocking.CloseRegistryKey(key)

	// Read the value of Enabled, which contains the information if Advertisement ID is shared with apps or not
	value, _, err := key.GetIntegerValue("Enabled")
	if err != nil {
		logger.Log.ErrorWithErr("Error reading Enabled value", err)
		return checks.NewCheckError(checks.AdvertisementID, err)
	}

	if value == 1 {
		return checks.NewCheckResult(checks.AdvertisementID, 1)
	}
	return checks.NewCheckResult(checks.AdvertisementID, 0)
}
