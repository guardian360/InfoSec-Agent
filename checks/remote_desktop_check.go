package checks

import (
	"golang.org/x/sys/windows/registry"
	"log"
)

// RemoteDesktopCheck checks if Remote Desktop is enabled
//
// Parameters: _
//
// Returns: If Remote Desktop is enabled or not
func RemoteDesktopCheck() Check {
	// Open the registry key for Terminal Server settings
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Terminal Server`,
		registry.QUERY_VALUE)
	if err != nil {
		return NewCheckErrorf("RemoteDesktop", "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer func(key registry.Key) {
		err := key.Close()
		if err != nil {
			log.Printf("error closing registry key: %v", err)
		}
	}(key)

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
