package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TODO: Update documentation
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
