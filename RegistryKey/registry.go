package RegistryKey

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
)

// OpenRegistryKey opens registry keys and handles associated errors
//
// Parameters: k (registry.Key) - the registry key to open,
//
// path (string) represents the path to the registry key
//
// Returns: The opened registry key
func OpenRegistryKey(k RegistryKey, path string) (RegistryKey, error) {
	key, err := k.OpenKey(path, registry.READ)

	if err != nil {
		return key, fmt.Errorf("error opening registry key: %w", err)
	}

	return key, nil
}

// CloseRegistryKey closes registry keys and handles associated errors
//
// Parameters: key (registry.Key) - the registry key to close
//
// Returns: _
func CloseRegistryKey(key RegistryKey) {
	err := key.Close()
	if err != nil {
		log.Printf("error closing registry key: %s", err)
	}
}

// FindEntries returns the values of the entries inside the corresponding registry key
//
// Parameters: entries ([]string) - the entries to check in the registry key,
//
// key (registry.Key) represents the registry key in which to look
//
// Returns: A slice of the values of the entries inside the registry key
func FindEntries(entries []string, key RegistryKey) []string {
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

// CheckKey checks the value of a certain element within a registry key
//
// Parameters: key (registry.Key) - the registry key to be checked,
//
// elem (string) - the element to be checked
//
// Returns: The value of the element within the registry key
func CheckKey(key RegistryKey, elem string) string {
	val, _, err := key.GetStringValue(elem)
	if err == nil {
		return val
	} else {
		return "-1"
	}
}
