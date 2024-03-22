package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/RegistryKey"
)

// SecureBoot checks if Windows secure boot is enabled
//
// Parameters: _
//
// Returns: If Windows secure boot is enabled or not
func SecureBoot(registryKey RegistryKey.RegistryKey) Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := RegistryKey.OpenRegistryKey(registryKey, `SYSTEM\CurrentControlSet\Control\SecureBoot\State`)
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}
	defer RegistryKey.CloseRegistryKey(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return NewCheckResult("SecureBoot", "Secure boot is enabled")
	}
	if secureBootStatus == 0 {
		return NewCheckResult("SecureBoot", "Secure boot is disabled")
	}
	return NewCheckResult("SecureBoot", "Secure boot status is unknown")
}
