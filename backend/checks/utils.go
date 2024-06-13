package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

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
func FindEntries(entries []string, key mocking.RegistryKey) []string {
	elements := make([]string, 0)

	for _, element := range entries {
		val, _, _ := key.GetBinaryValue(element)

		// Check the binary values to make sure we only return the programs that are ENABLED on startup
		// This is because the registry lists all programs that are related to the start-up,
		// including those that are disabled.
		// Start up programs that are enabled have a binary signature of non-zero at index 0 and zero at the rest of the indices.
		if val[0] == 0 {
			continue
		}
		if CheckAllZero(val[1:]) {
			elements = append(elements, element)
		}
	}
	return elements
}

// CheckAllZero is a helper function that checks if all elements in a byte slice are zero.
//
// Parameters:
//   - entries: A byte slice representing the elements to be checked.
//
// Returns:
//   - A boolean value indicating whether all elements in the byte slice are zero. If all elements are zero, the function returns true; otherwise, it returns false.
func CheckAllZero(entries []byte) bool {
	for _, entry := range entries {
		if entry != 0 {
			return false
		}
	}
	return true
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
func CheckKey(key mocking.RegistryKey, elem string) string {
	val, _, err := key.GetStringValue(elem)
	if err == nil {
		return val
	}
	return "-1"
}
