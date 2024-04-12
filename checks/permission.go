package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

const nonpackaged = "NonPackaged"

// Permission checks if the user has given permission to an application to access a certain capability
//
// Parameters: permission (string) represents the permission to check
//
// Returns: A list of applications that have the given permission
func Permission(permissionID int, permission string, registryKey mocking.RegistryKey) Check {
	var err error
	var appKey mocking.RegistryKey
	var nonPackagedApplicationNames []string
	// Open the registry key for the given permission
	key, err := mocking.OpenRegistryKey(registryKey,
		`Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission)
	if err != nil {
		return NewCheckErrorf(permissionID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer mocking.CloseRegistryKey(key)

	// Get the names of all sub-keys (which represent applications)
	applicationNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return NewCheckErrorf(permissionID, "error reading subkey names", err)
	}

	var results []string
	var val string
	// Iterate through the application names and append them to the results
	for _, appName := range applicationNames {
		appKey, err = mocking.OpenRegistryKey(key, appKeyName(appName))
		defer mocking.CloseRegistryKey(appKey)
		if err != nil {
			return NewCheckErrorf(permissionID, "error opening registry key", err)
		}
		if appName == nonpackaged {
			val, _, err = key.GetStringValue("Value")
		} else {
			val, _, err = appKey.GetStringValue("Value")
		}
		if err != nil {
			return NewCheckErrorf(permissionID, "error reading value", err)
		}
		// If the value is not "Allow", the application does not have permission
		if val != "Allow" {
			continue
		}
		if appName == nonpackaged {
			nonPackagedApplicationNames, err = nonPackagedAppNames(appKey)
			if err != nil {
				return NewCheckErrorf(permissionID, "error reading subkey names", err)
			}
			results = append(results, nonPackagedApplicationNames...)
		} else {
			winApp := strings.Split(appName, "_")
			results = append(results, winApp[0])
		}
	}
	// Remove duplicate results
	filteredResults := utils.RemoveDuplicateStr(results)
	return NewCheckResult(permissionID, 0, filteredResults...)
}

// appKeyName returns the appropriate key name for the given application name
func appKeyName(appName string) string {
	if appName == nonpackaged {
		return nonpackaged
	}
	return appName
}

// nonPackagedAppNames returns the names of non-packaged applications
func nonPackagedAppNames(appKey mocking.RegistryKey) ([]string, error) {
	nonPackagedApplicationNames, err := appKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, nonPackagedAppName := range nonPackagedApplicationNames {
		exeString := strings.Split(nonPackagedAppName, "#")
		results = append(results, exeString[len(exeString)-1])
	}
	return results, nil
}
