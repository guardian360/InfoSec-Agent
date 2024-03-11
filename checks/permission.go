// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Permission checks if the user has given permission to an application to access a certain capability
//
// Parameters: permission (string) represents the permission to check
//
// Returns: A list of applications that have the given permission
func Permission(permission string) Check {
	// Open the registry key for the given permission
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission, registry.READ)
	if err != nil {
		return newCheckErrorf(permission, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer key.Close()

	// Get the names of all sub-keys (which represent applications)
	applicationNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return newCheckErrorf(permission, "error reading subkey names", err)
	}

	var results []string

	// Iterate through the application names and append them to the results
	for _, appName := range applicationNames {
		// The registry key for packaged/non-packaged applications is different, so they get handled separately
		if appName == "NonPackaged" {
			key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\NonPackaged`, registry.READ)
			nonPackagedApplicationNames, err := key.ReadSubKeyNames(-1)
			v, vint, err := key.GetStringValue("Value")

			// Check if the application has the specified permission
			if vint == 1 && err == nil && v == "Allow" {
				for _, nonPackagedAppName := range nonPackagedApplicationNames {
					exeString := strings.Split(nonPackagedAppName, "#")
					results = append(results, exeString[len(exeString)-1])
				}
			}
		} else {
			key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\`+appName, registry.READ)
			v, vint, err := key.GetStringValue("Value")

			// Check if the application has the specified permission
			if vint == 1 && err == nil && v == "Allow" {
				winApp := strings.Split(appName, "_")
				results = append(results, winApp[0])
			}
		}
	}
	// Remove duplicate results
	filteredResults := removeDuplicateStr(results)
	return newCheckResult(permission, filteredResults...)
}

// removeDuplicateStr removes duplicate strings from a slice
//
// Parameters: strSlice (string slice) represents the slice to remove duplicates from
//
// Returns: A slice with the duplicates removed
func removeDuplicateStr(strSlice []string) []string {
	// Keep a map of found values, where true means the value has (already) been found
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			// If the value is found for the first time, append it to the list of results
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
