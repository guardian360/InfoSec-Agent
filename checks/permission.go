package checks

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Permission(permission string) Check {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission, registry.READ)
	if err != nil {
		return newCheckErrorf(permission, "error opening registry key", err)
	}
	defer key.Close()

	// Get the names of all subkeys (which represent applications)
	applicationNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return newCheckErrorf(permission, "error reading subkey names", err)
	}

	var results []string

	// Iterate through the application names and print them
	for _, appName := range applicationNames {
		if appName == "NonPackaged" {
			key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\NonPackaged`, registry.READ)
			nonPackagedApplicationNames, err := key.ReadSubKeyNames(-1)
			v, vint, err := key.GetStringValue("Value")
			if vint == 1 && err == nil && v == "Allow" {
				for _, nonPackagedAppName := range nonPackagedApplicationNames {
					exeString := strings.Split(nonPackagedAppName, "#")
					results = append(results, exeString[len(exeString)-1])
				}
			}
		} else {
			key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`+permission+`\`+appName, registry.READ)
			v, vint, err := key.GetStringValue("Value")
			if vint == 1 && err == nil && v == "Allow" {
				winApp := strings.Split(appName, "_")
				results = append(results, winApp[0])
			}
		}
	}
	filteredResults := removeDuplicateStr(results)
	return newCheckResult(permission, filteredResults...)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
