package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// SecureBoot is a function that checks if Windows Secure Boot is enabled on the system.
//
// Parameters:
//   - registryKey registrymock.RegistryKey: A registry key object for accessing the Windows Secure Boot registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Secure Boot is enabled or not.
//
// The function works by opening the Windows Secure Boot registry key and reading its 'UEFISecureBootEnabled' value. This value represents the status of Secure Boot. If the value is 1, Secure Boot is enabled. If the value is 0, Secure Boot is disabled. If the function encounters an error while accessing the registry key or reading the value, it returns a Check instance containing an error message. If the 'UEFISecureBootEnabled' value is not 1 or 0, the function returns a Check instance indicating that the Secure Boot status is unknown.
func SecureBoot(registryKey registrymock.RegistryKey) Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := registrymock.OpenRegistryKey(registryKey,
		`SYSTEM\CurrentControlSet\Control\SecureBoot\State`)
	if err != nil {
		return NewCheckError("SecureBoot", err)
	}
	defer registrymock.CloseRegistryKey(windowsSecureBoot)

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
