package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// RemoteDesktopCheck is a function that checks if the Remote Desktop feature is enabled on the system.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): A mocker of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the Remote Desktop feature is enabled or not.
//
// The function works by opening the registry key for Terminal Server settings. It then reads the value of 'fDenyTSConnections', which indicates whether Remote Desktop is enabled or not. If the value is 0, it means that Remote Desktop is enabled. Otherwise, it is disabled. The function returns a Check instance containing the result of the check.
func RemoteDesktopCheck(registryKey mocking.RegistryKey) checks.Check {
	// Open the registry key for Terminal Server settings
	key, err := checks.OpenRegistryKey(registryKey, `System\CurrentControlSet\Control\Terminal Server`)

	if err != nil {
		return checks.NewCheckErrorf(checks.RemoteDesktopID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer checks.CloseRegistryKey(key)

	// Read the value of fDenyTSConnections, which contains the information if Remote Desktop is enabled or not
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return checks.NewCheckErrorf(checks.RemoteDesktopID, "error reading fDenyTSConnections", err)
	}
	// Check if Remote Desktop is enabled or disabled based on the value of fDenyTSConnections
	if val == 0 {
		return checks.NewCheckResult(checks.RemoteDesktopID, 0, "Remote Desktop is enabled")
	}
	return checks.NewCheckResult(checks.RemoteDesktopID, 1, "Remote Desktop is disabled")
}
