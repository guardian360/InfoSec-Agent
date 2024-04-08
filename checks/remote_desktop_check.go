package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// RemoteDesktopCheck checks if Remote Desktop is enabled
//
// Parameters: registryKey (registrymock.RegistryKey) - A Windows registry mock
//
// Returns: If Remote Desktop is enabled or not
func RemoteDesktopCheck(registryKey registrymock.RegistryKey) Check {
	// Open the registry key for Terminal Server settings
	key, err := registrymock.OpenRegistryKey(registryKey, `System\CurrentControlSet\Control\Terminal Server`)

	if err != nil {
		return NewCheckErrorf(RemoteDesktopID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Read the value of fDenyTSConnections, which contains the information if Remote Desktop is enabled or not
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return NewCheckErrorf(RemoteDesktopID, "error reading fDenyTSConnections", err)
	}
	// Check if Remote Desktop is enabled or disabled based on the value of fDenyTSConnections
	if val == 0 {
		return NewCheckResult(RemoteDesktopID, 0, "Remote Desktop is enabled")
	}
	return NewCheckResult(RemoteDesktopID, 1, "Remote Desktop is disabled")
}
