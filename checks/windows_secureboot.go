package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// SecureBoot checks if Windows secure boot is enabled
//
// Parameters: _
//
// Returns: If Windows secure boot is enabled or not
func SecureBoot(registryKey mocking.RegistryKey) Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := mocking.OpenRegistryKey(registryKey,
		`SYSTEM\CurrentControlSet\Control\SecureBoot\State`)
	if err != nil {
		return NewCheckError(SecureBootID, err)
	}
	defer mocking.CloseRegistryKey(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return NewCheckError(SecureBootID, err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return NewCheckResult(SecureBootID, 1, "Secure boot is enabled")
	}
	if secureBootStatus == 0 {
		return NewCheckResult(SecureBootID, 0, "Secure boot is disabled")
	}
	return NewCheckResult(SecureBootID, 2, "Secure boot status is unknown")
}
