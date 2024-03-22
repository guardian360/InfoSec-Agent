package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/RegistryKey"
	"golang.org/x/sys/windows/registry"
)

// RemoteDesktopCheck checks if Remote Desktop is enabled
//
// Parameters: _
//
// Returns: If Remote Desktop is enabled or not
func RemoteDesktopCheck() Check {
	// Open the registry key for Terminal Server settings
	key, err := RegistryKey.OpenRegistryKey(RegistryKey.NewRegistryKeyWrapper(registry.LOCAL_MACHINE), `System\CurrentControlSet\Control\Terminal Server`)
	//was registry.QUERY_VALUE, is now registry.READ
	if err != nil {
		return NewCheckErrorf("RemoteDesktop", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer RegistryKey.CloseRegistryKey(key)

	// Read the value of fDenyTSConnections, which contains the information if Remote Desktop is enabled or not
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return NewCheckErrorf("RemoteDesktop", "error reading fDenyTSConnections", err)
	} else {
		// Check if Remote Desktop is enabled or disabled based on the value of fDenyTSConnections
		if val == 0 {
			return NewCheckResult("RemoteDesktop", "Remote Desktop is enabled")
		} else {
			return NewCheckResult("RemoteDesktop", "Remote Desktop is disabled")
		}
	}

}
