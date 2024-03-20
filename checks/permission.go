package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
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
	key, err := utils.OpenRegistryKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission)
	if err != nil {
		return NewCheckErrorf(permission, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer utils.CloseRegistryKey(key)

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
			key, err = utils.OpenRegistryKey(registry.CURRENT_USER,
				`Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\NonPackaged`)
			defer utils.CloseRegistryKey(key)
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
			key, err = utils.OpenRegistryKey(registry.CURRENT_USER,
				`Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\`+appName)
			defer utils.CloseRegistryKey(key)
			v, vint, err := key.GetStringValue("Value")

			// Check if the application has the specified permission
			if vint == 1 && err == nil && v == "Allow" {
				winApp := strings.Split(appName, "_")
				results = append(results, winApp[0])
			}
		}
	}
	// Remove duplicate results
	filteredResults := utils.RemoveDuplicateStr(results)
	return NewCheckResult(permission, filteredResults...)
}
