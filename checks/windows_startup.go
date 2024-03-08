package checks

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func Startup() Check {
	// Open the relevant key so we can get the startup data entries out of them
	cuKey, err1 := openRegistryKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := openRegistryKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := openRegistryKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)

	if err1 != nil || err2 != nil || err3 != nil {
		return newCheckError("Startup", fmt.Errorf("error opening registry keys"))
	}

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

// Simple function to open registry keys and handle errors
func openRegistryKey(k registry.Key, path string) (registry.Key, error) {
	key, err := registry.OpenKey(k, path, registry.READ)

	if err != nil {
		return key, fmt.Errorf("error opening registry key: %w", err)
	}

	return key, nil
}

// Print the values of the entries inside the corresponding registry key
func findEntries(entries []string, key registry.Key) []string {
	elements := make([]string, 0)

	for _, element := range entries {
		val, _, _ := key.GetBinaryValue(element)

		// We check the binary values to make sure we only print the programs ENABLED to startup
		// For example: in registry it lists all startup programs, but also the ones that are disabled on startup.
		if val[4] == 0 && val[5] == 0 && val[6] == 0 {
			elements = append(elements, element)
		}
	}

	return elements
}
