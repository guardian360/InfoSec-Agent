package checks

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func Defender() {
	// Open the Windows Defender registry key
	windowsDefenderKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer windowsDefenderKey.Close()
	realTimeKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Real-Time Protection`, registry.READ)

	// Read the value of the registry key
	antiVirusPeriodic, _, err := windowsDefenderKey.GetIntegerValue("DisableAntiVirus")
	realTimeDefender, _, err := realTimeKey.GetIntegerValue("DisableRealtimeMonitoring")
	if err != nil {
		if antiVirusPeriodic == 1 {
			fmt.Println("Windows real-time defender is enabled but the windows periodic scan is disabled")
		} else {
			fmt.Println("Windows real-time defender is enabled and also the windows periodic scan is enabled")
		}
	} else {
		if realTimeDefender == 1 {
			if antiVirusPeriodic == 1 {
				fmt.Println("Windows real-time defender is disabled and also the windows periodic scan is disabled")
			} else {
				fmt.Println("Windows real-time defender is disabled but the windows periodic scan is enabled")
			}
		} else {
			fmt.Println("No windows defender data found")
		}
	}
}
