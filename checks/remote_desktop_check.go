package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// RemoteDesktopCheck is a function that checks if the Remote Desktop feature is enabled on the system.
//
// Parameters:
//   - registryKey (registrymock.RegistryKey): A mock of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the Remote Desktop feature is enabled or not.
//
// The function works by opening the registry key for Terminal Server settings. It then reads the value of 'fDenyTSConnections', which indicates whether Remote Desktop is enabled or not. If the value is 0, it means that Remote Desktop is enabled. Otherwise, it is disabled. The function returns a Check instance containing the result of the check.
func RemoteDesktopCheck(registryKey registrymock.RegistryKey) Check {
	// Open the registry key for Terminal Server settings
	key, err := registrymock.OpenRegistryKey(registryKey, `System\CurrentControlSet\Control\Terminal Server`)

	if err != nil {
		return NewCheckErrorf("RemoteDesktop", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Read the value of fDenyTSConnections, which contains the information if Remote Desktop is enabled or not
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return NewCheckErrorf("RemoteDesktop", "error reading fDenyTSConnections", err)
	}
	// Check if Remote Desktop is enabled or disabled based on the value of fDenyTSConnections
	if val == 0 {
		return NewCheckResult("RemoteDesktop", "Remote Desktop is enabled")
	}
	return NewCheckResult("RemoteDesktop", "Remote Desktop is disabled")
}
