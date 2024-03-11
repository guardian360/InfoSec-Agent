package checks

import (
	"golang.org/x/sys/windows/registry"
)

// SecureBoot checks if windows secure boot is enabled or disabled
func SecureBoot() Check {
	windowsSecureBoot, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\SecureBoot\State`, registry.READ)

	if err != nil {
		return newCheckError("SecureBoot", err)
	}
	defer func(windowsSecureBoot registry.Key) {
		err := windowsSecureBoot.Close()
		if err != nil {
			return
		}
	}(windowsSecureBoot)
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return newCheckError("SecureBoot", err)
	}

	// Output successful result
	if secureBootStatus == 1 {
		return newCheckResult("SecureBoot", "Secure boot is enabled")
	}

	return newCheckResult("SecureBoot", "Secure boot is disabled")
}
