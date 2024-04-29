package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

const nonPackaged = "NonPackaged"

// Permission is a function that checks if a user has granted a specific permission to an application.
//
// Parameters:
//   - permissionID (int): The ID of the permission check.
//   - permission (string): The specific permission to check.
//   - registryKey (mocking.RegistryKey): The registry key to use for the check.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the permission check. The Result field of the Check instance will contain a list of applications that have been granted the specified permission.
//
// This function opens the registry key for the given permission and retrieves the names of all sub-keys, which represent applications. It then iterates through these applications, checking if they have been granted the specified permission. If the permission value is "Allow", the application name is added to the results. The function also handles non-packaged applications separately. Finally, it removes any duplicate results before returning them.
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
		return NewCheckErrorf(permissionID, "error reading sub-key names", err)
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
		if appName == nonPackaged {
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
		if appName == nonPackaged {
			nonPackagedApplicationNames, err = nonPackagedAppNames(appKey)
			if err != nil {
				return NewCheckErrorf(permissionID, "error reading sub-key names", err)
			}
			results = append(results, nonPackagedApplicationNames...)
		} else {
			winApp := strings.Split(appName, "_")
			results = append(results, winApp[0])
		}
	}
	// Remove duplicate results
	filteredResults := RemoveDuplicateStr(results)
	prettyResults := []string{}
	for _, result := range filteredResults {
		cleanedFilename := strings.TrimSuffix(result, ".exe")
		cleanedFilename = strings.ReplaceAll(cleanedFilename, ".", " ")
		prettyResults = append(prettyResults, cleanedFilename)
	}
	if len(prettyResults) == 0 {
		return NewCheckResult(permissionID, 1)
	}
	return NewCheckResult(permissionID, 0, prettyResults...)
}

// appKeyName is a helper function that returns the appropriate registry key name for a given application name.
//
// Parameters:
//   - appName (string): The name of the application for which the registry key name is required.
//
// Returns:
//   - string: The appropriate registry key name for the given application name.
//
// This function is used to handle a special case where the application name is "NonPackaged". In such a case, it returns the string "NonPackaged" as the registry key name. For all other application names, it returns the application name itself as the registry key name. This function is used in the context of checking permissions for applications.
func appKeyName(appName string) string {
	if appName == nonPackaged {
		return nonPackaged
	}
	return appName
}

// nonPackagedAppNames is a helper function that retrieves the names of non-packaged applications from a given registry key.
//
// Parameters:
//   - appKey (mocking.RegistryKey): The registry key that contains the sub-keys representing non-packaged applications.
//
// Returns:
//   - []string: A slice of strings representing the names of non-packaged applications.
//   - error: An error object that describes the error, if any occurred during the operation.
//
// This function reads the names of all sub-keys from the provided registry key, which represent non-packaged applications. It then iterates through these names, splitting each one at the '#' character and appending the last segment to the results. This is done because the names of non-packaged applications are stored in the format 'path#applicationName'. The function returns the list of application names, or an error if one occurred during the operation.
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

// RemoveDuplicateStr is a utility function that eliminates duplicate string values from a given slice.
//
// Parameters:
//   - strSlice []string: The input slice from which duplicate string values need to be removed.
//
// Returns:
//   - []string: A new slice that contains the unique string values from the input slice. The order of the elements is preserved based on their first occurrence in the input slice.
func RemoveDuplicateStr(strSlice []string) []string {
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
