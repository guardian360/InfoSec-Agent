package utils

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
)

// OpenRegistryKey opens registry keys and handles associated errors
//
// Parameters: k (registry.Key) represents the registry key to open,
// path (string) represents the path to the registry key
//
// Returns: The opened registry key
func OpenRegistryKey(k registry.Key, path string) (registry.Key, error) {
	key, err := registry.OpenKey(k, path, registry.READ)

	if err != nil {
		return key, fmt.Errorf("error opening registry key: %w", err)
	}

	return key, nil
}

// CloseRegistryKey closes registry keys and handles associated errors
//
// Parameters: key (registry.Key) represents the registry key to close
//
// Returns: _
func CloseRegistryKey(key registry.Key) {
	err := key.Close()
	if err != nil {
		log.Printf("error closing registry key: %s", err)
	}
}

// FindEntries returns the values of the entries inside the corresponding registry key
//
// Parameters: entries ([]string) represents the entries to check in the registry key,
// key (registry.Key) represents the registry key in which to look
//
// Returns: A slice of the values of the entries inside the registry key
func FindEntries(entries []string, key registry.Key) []string {
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
// Parameters: key (registry.Key) representing the registry key to be checked,
// el (string) representing the element to be checked
//
// Returns: The value of the element within the registry key
func CheckKey(key registry.Key, el string) string {
	val, _, err := key.GetStringValue(el)
	if err == nil {
		return val
	} else {
		return "-1"
	}
}
