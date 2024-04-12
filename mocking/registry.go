package mocking

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"golang.org/x/sys/windows/registry"
)

// OpenRegistryKey is a function that opens a specified registry key and handles any associated errors.
//
// Parameters:
//   - k: A RegistryKey object representing the base registry key from which the specified path will be opened.
//   - path: A string representing the path to the registry key to be opened, relative to the base registry key.
//
// Returns:
//   - A RegistryKey object representing the opened registry key.
//   - An error object that encapsulates any error that occurred while trying to open the registry key. If no error occurred, this will be nil.
//
// Note: This function is designed to handle errors that may occur when opening a registry key, such as the key not existing. If an error occurs, it will be wrapped with additional context and returned, allowing the caller to handle it appropriately.
func OpenRegistryKey(k RegistryKey, path string) (RegistryKey, error) {
	key, err := k.OpenKey(path, registry.READ)

	if err != nil {
		return key, fmt.Errorf("error opening registry key: %w", err)
	}

	return key, nil
}

// CloseRegistryKey is a function that closes a specified registry key and logs any associated errors.
//
// Parameter:
//   - key: A RegistryKey object representing the registry key to be closed.
//
// Returns: None. If an error occurs while closing the registry key, the error is logged and not returned.
//
// Note: This function is designed to handle errors that may occur when closing a registry key. If an error occurs, it is logged with additional context, allowing for easier debugging and error tracking.
func CloseRegistryKey(key RegistryKey) {
	err := key.Close()
	if err != nil {
		logger.Log.ErrorWithErr("Error closing registry key:", err)
	}
}

// TODO: SHOULD BE MOVED TO STARTUP.GO !!

// FindEntries scans a specified registry key for a list of entries and returns the values of those entries.
//
// Parameters:
//   - entries: A slice of strings representing the names of the entries to be checked within the registry key.
//   - key: A RegistryKey object representing the registry key to be scanned.
//
// Returns:
//   - A slice of strings containing the values of the specified entries within the registry key. Only entries that are enabled on startup are included. This is determined by checking the binary values of the entries; entries with a binary value of 0 at indices 4, 5, and 6 are considered enabled.
//
// Note: This function is designed to handle the retrieval of startup-related programs from the registry. It filters out disabled programs to provide a list of only the enabled ones.
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

// CheckKey retrieves the value of a specified element within a given registry key.
//
// Parameters:
//   - key: A RegistryKey object representing the registry key to be checked.
//   - elem: A string representing the name of the element whose value is to be retrieved.
//
// Returns:
//   - A string representing the value of the specified element within the registry key. If the element does not exist or an error occurs while retrieving its value, the function returns "-1".
//
// Note: This function is designed to handle the retrieval of values from the registry. It encapsulates the process of accessing the registry and retrieving a value, providing a simplified interface for this operation.
func CheckKey(key RegistryKey, elem string) string {
	val, _, err := key.GetStringValue(elem)
	if err == nil {
		return val
	}
	return "-1"
}
