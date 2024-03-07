package checks

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// SecureBoot checks if windows secure boot is enabled or disabled
func SecureBoot() {

	windowsSecureBoot, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\SecureBoot\State`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer func(windowsSecureBoot registry.Key) {
		err := windowsSecureBoot.Close()
		if err != nil {
			return
		}
	}(windowsSecureBoot)
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		fmt.Println("Error reading registry key:", err)
		return
	}
	if secureBootStatus == 1 {
		fmt.Println("Secure boot is enabled")
	} else {
		fmt.Println("Secure boot is disabled")
	}
}
