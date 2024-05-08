package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// SecureBoot is a function that checks if Windows Secure Boot is enabled on the system.
//
// Parameters:
//   - registryKey mocking.RegistryKey: A registry key object for accessing the Windows Secure Boot registry key.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether Windows Secure Boot is enabled or not.
//
// The function works by opening the Windows Secure Boot registry key and reading its 'UEFISecureBootEnabled' value. This value represents the status of Secure Boot. If the value is 1, Secure Boot is enabled. If the value is 0, Secure Boot is disabled. If the function encounters an error while accessing the registry key or reading the value, it returns a Check instance containing an error message. If the 'UEFISecureBootEnabled' value is not 1 or 0, the function returns a Check instance indicating that the Secure Boot status is unknown.
func SecureBoot(registryKey mocking.RegistryKey) checks.Check {
	// Get secure boot information from the registry
	windowsSecureBoot, err := checks.OpenRegistryKey(registryKey,
		`SYSTEM\CurrentControlSet\Control\SecureBoot\State`)
	if err != nil {
		return checks.NewCheckError(checks.SecureBootID, err)
	}
	defer checks.CloseRegistryKey(windowsSecureBoot)

	// Read the status of secure boot
	secureBootStatus, _, err := windowsSecureBoot.GetIntegerValue("UEFISecureBootEnabled")
	if err != nil {
		return checks.NewCheckError(checks.SecureBootID, err)
	}

	// Using the status, determine if secure boot is enabled or not
	if secureBootStatus == 1 {
		return checks.NewCheckResult(checks.SecureBootID, 1)
	}
	if secureBootStatus == 0 {
		return checks.NewCheckResult(checks.SecureBootID, 0)
	}
	return checks.NewCheckResult(checks.SecureBootID, 2)
}
