package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// Permission checks if the user has given permission to an application to access a certain capability
//
// Parameters: permission (string) represents the permission to check
//
// Returns: A list of applications that have the given permission
func Permission(permission string, registryKey registrymock.RegistryKey) Check {
	// Open the registry key for the given permission
	key, err := registrymock.OpenRegistryKey(registryKey,
		`Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission)
	if err != nil {
		return NewCheckErrorf(permission, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Get the names of all sub-keys (which represent applications)
	applicationNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return NewCheckErrorf(permission, "error reading subkey names", err)
	}

	var results []string

	// Iterate through the application names and append them to the results
	for _, appName := range applicationNames {
		// The registry key for packaged/non-packaged applications is different, so they get handled separately
		if appName == "NonPackaged" {
			appKey, err := registrymock.OpenRegistryKey(key, `NonPackaged`)
			if err != nil {
				return NewCheckErrorf(permission, "error opening registry key", err)
			}
			defer registrymock.CloseRegistryKey(appKey)
			nonPackagedApplicationNames, err := appKey.ReadSubKeyNames(-1)
			v, _, err := key.GetStringValue("Value")
			registrymock.CloseRegistryKey(appKey)

			// Check if the application has the specified permission
			if err == nil && v == "Allow" {
				for _, nonPackagedAppName := range nonPackagedApplicationNames {
					exeString := strings.Split(nonPackagedAppName, "#")
					results = append(results, exeString[len(exeString)-1])
				}
			}
		} else {
			appKey, err := registrymock.OpenRegistryKey(key, appName)
			v, _, err := appKey.GetStringValue("Value")
			registrymock.CloseRegistryKey(appKey)

			// Check if the application has the specified permission
			if err == nil && v == "Allow" {
				winApp := strings.Split(appName, "_")
				results = append(results, winApp[0])
			}
		}
	}
	// Remove duplicate results
	filteredResults := utils.RemoveDuplicateStr(results)
	return NewCheckResult(permission, filteredResults...)
}
