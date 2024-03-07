package checks

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func RemoteDesktopCheck() {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Terminal Server`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()
	val, _, err := key.GetIntegerValue("fDenyTSConnections")
	if err != nil {
		fmt.Println("Error reading fDenyTSConnections:", err)
		return
	} else {
		if val == 0 {
			fmt.Println("Remote Desktop is enabled")
		} else {
			fmt.Println("Remote Desktop is disabled")
		}
	}

}
