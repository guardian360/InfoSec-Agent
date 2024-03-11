// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// Startup checks the registry for startup programs
//
// Parameters: _
//
// Returns: A list of start-up programs
func Startup() Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := openRegistryKey(registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := openRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := openRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)

	if err1 != nil || err2 != nil || err3 != nil {
		return newCheckError("Startup", fmt.Errorf("error opening registry keys"))
	}

	// Close the keys after we have received all relevant information
	defer cuKey.Close()
	defer lmKey.Close()
	defer lmKey2.Close()

	// Read the entries within the registry key
	cuValueNames, err1 := cuKey.ReadValueNames(0)
	lmValueNames, err2 := lmKey.ReadValueNames(0)
	lm2ValueNames, err3 := lmKey2.ReadValueNames(0)

	if err1 != nil || err2 != nil || err3 != nil {
		return newCheckError("Startup", fmt.Errorf("error reading value names"))
	}

	output := make([]string, 0)
	output = append(output, findEntries(cuValueNames, cuKey)...)
	output = append(output, findEntries(lmValueNames, lmKey)...)
	output = append(output, findEntries(lm2ValueNames, lmKey2)...)

	return newCheckResult("Startup", output...)
}

// openRegistryKey opens registry keys and handles associated errors
//
// Parameters: k (registry.Key) represents the registry key to open,
// path (string) represents the path to the registry key
//
// Returns: The opened registry key
func openRegistryKey(k registry.Key, path string) (registry.Key, error) {
	key, err := registry.OpenKey(k, path, registry.READ)

	if err != nil {
		return key, fmt.Errorf("error opening registry key: %w", err)
	}

	return key, nil
}

// findEntries returns the values of the entries inside the corresponding registry key
//
// Parameters: entries ([]string) represents the entries to check in the registry key,
// key (registry.Key) represents the registry key in which to look
//
// Returns: A slice of the values of the entries inside the registry key
func findEntries(entries []string, key registry.Key) []string {
	elements := make([]string, 0)

	for _, element := range entries {
		val, _, _ := key.GetBinaryValue(element)

		// Check the binary values to make sure we only return the programs that are ENABLED on startup
		// This is because the registry lists all programs that are related to the start-up,
		// including those that are disabled
		if val[4] == 0 && val[5] == 0 && val[6] == 0 {
			elements = append(elements, element)
		}
	}

	return elements
}
