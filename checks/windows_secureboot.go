package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"golang.org/x/sys/windows/registry"
)

// SecureBoot checks if Windows secure boot is enabled
//
// Parameters: _
//
// Returns: If Windows secure boot is enabled or not
func SecureBoot() Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := utils.OpenRegistryKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\SecureBoot\State`)
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}
	defer utils.CloseRegistryKey(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return NewCheckResult("SecureBoot", "Secure boot is enabled")
	}

	return NewCheckResult("SecureBoot", "Secure boot is disabled")
}
