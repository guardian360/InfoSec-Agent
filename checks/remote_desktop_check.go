package checks

import (
	"golang.org/x/sys/windows/registry"
)

func RemoteDesktopCheck() Check {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Terminal Server`, registry.QUERY_VALUE)
	if err != nil {
		return newCheckErrorf("RemoteDesktop", "error opening registry key", err)
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		return newCheckErrorf("RemoteDesktop", "error reading fDenyTSConnections", err)
	} else {
		if val == 0 {
			return newCheckResult("RemoteDesktop", "Remote Desktop is enabled")
		} else {
			return newCheckResult("RemoteDesktop", "Remote Desktop is disabled")
		}
	}

}
