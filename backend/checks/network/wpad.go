package network

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// WPADEnabled checks if the WPAD service is enabled.
// The WPAD service is used to automatically configure proxy settings for a network.
// If the service is running, it is possible that an attacker could use it to redirect traffic.
//
// Parameters:
//   - key (mocking.RegistryKey): A mock registry key that allows access to the Windows registry for testing purposes.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the WPAD check, or an error if one occurred.
func WPADEnabled(key mocking.RegistryKey) checks.Check {
	resultID := 0
	output, err := mocking.OpenRegistryKey(key, `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Wpad`)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading key", err)
		return checks.NewCheckError(checks.WPADID, err)
	}
	disabled, _, intErr := output.GetIntegerValue("WpadOverride")
	if intErr != nil || disabled != 1 {
		logger.Log.ErrorWithErr("Error reading value", intErr)
		resultID++
	}

	return checks.NewCheckResult(checks.WPADID, resultID)
}
