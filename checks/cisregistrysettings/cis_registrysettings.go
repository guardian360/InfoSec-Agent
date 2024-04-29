// Package cisregistrysettings provides a set of functions to check various registry settings
// to ensure they adhere to the CIS Benchmark standards. Each function takes a RegistryKey object
// as an argument, which represents the root key from which the registry settings will be checked.
// The functions return a slice of boolean values, where each boolean represents whether a particular
// registry setting adheres to the CIS Benchmark standards.
package cisregistrysettings

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// CISRegistrySettings is a function that checks various registry settings to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked.
//
// Returns:
//
//   - []bool: A slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
func CISRegistrySettings(localMachineKey mocking.RegistryKey, usersKey mocking.RegistryKey) checks.Check {
	results := make([]bool, 0)
	// Following function(s) need the HKEY_LOCAL_MACHINE registry key
	results = append(results, CheckServices(localMachineKey)...)
	results = append(results, CheckPoliciesHKLM(localMachineKey)...)
	results = append(results, CheckOtherRegistrySettings(localMachineKey)...)
	// Following function(s) need the HKEY_USERS registry key
	results = append(results, CheckPoliciesHKU(usersKey)...)

	// Following function(s) need the HKEY_LOCAL_MACHINE registry key
	if checks.WinVersion == 10 {
		results = append(results, CheckWin10(localMachineKey)...)
	}
	if checks.WinVersion == 11 {
		results = append(results, CheckWin11(localMachineKey)...)
	}
	// check if results are all true
	for _, result := range results {
		if !result {
			return checks.NewCheckResult(checks.CISRegistrySettingsID, 0, "Not all registry settings adhere to the CIS Benchmark standards")
		}
	}
	return checks.NewCheckResult(checks.CISRegistrySettingsID, 1, "All registry settings adhere to the CIS Benchmark standards")
}

// checkIntegerValue is a helper function that checks if the integer value of a registry key matches the expected value.
//
// Parameters:
//
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - value (string): The name of the value to check.
//   - expected (interface{}): The expected value of the registry key.
//
// Returns:
//
//   - bool: A boolean value indicating whether the integer value of the registry key matches the expected value.
func checkIntegerValue(openKey mocking.RegistryKey, value string, expected interface{}) bool {
	val, _, err := openKey.GetIntegerValue(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading registry value of "+value, err)
		return false
	}
	// Determine functionality based on the value type of the expected parameter
	switch v := expected.(type) {
	// Single uint64, check if registry value is equal to expected value
	case uint64:
		return val == v
	// Slice of uint64 values, check if registry value is in the slice
	case []uint64:
		for _, i := range v {
			if val == i {
				return true
			}
		}
	// Slice of exactly 2 uint64 values, check if registry value is within the range
	case [2]uint64:
		return val >= v[0] && val <= v[1]
	default:
		return false
	}
	return false
}

// checkStringValue is a helper function that checks if the string value of a registry key matches the expected value.
//
// Parameters:
//
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - value (string): The name of the value to check.
//   - expected (string): The expected value of the registry key.
//
// Returns:
//
//   - bool: A boolean value indicating whether the string value of the registry key matches the expected value.
func checkStringValue(openKey mocking.RegistryKey, value string, expected string) bool {
	val, _, err := openKey.GetStringValue(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading registry value of "+value, err)
		return false
	}
	return val == expected
}

// openRegistryKeyWithErrHandling is a helper function that opens a registry key and handles any errors that occur.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The registry key to open.
//   - path (string): The path of the registry key to open.
//
// Returns:
//
//   - mocking.RegistryKey: The opened registry key.
//   - error: An error object that describes the error (if any) that occurred while opening the registry key. If no error occurred, this value is nil.
func openRegistryKeyWithErrHandling(registryKey mocking.RegistryKey, path string) (mocking.RegistryKey, error) {
	key, err := mocking.OpenRegistryKey(registryKey, path)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening registry key for CIS Audit list", err)
	}
	return key, err
}

// checkMultipleIntegerValues is a helper function that checks multiple integer values of a registry key against their expected values.
//
// Parameters:
//
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]interface{}): A slice of interface values representing the expected values of the registry keys.
//
// Returns:
//
//   - []bool: A slice of boolean values indicating whether the integer values of the registry keys match the expected values.
func checkMultipleIntegerValues(openKey mocking.RegistryKey, settings []string, expectedValues []interface{}) []bool {
	results := make([]bool, len(settings))
	for i, val := range settings {
		results[i] = checkIntegerValue(openKey, val, expectedValues[i])
	}
	return results
}

// checkMultipleStringValues is a helper function that checks multiple string values of a registry key against their expected values.
//
// Parameters:
//
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]string): A slice of strings representing the expected values of the registry keys.
//
// Returns:
//
//   - []bool: A slice of boolean values indicating whether the string values of the registry keys match the expected values.
func checkMultipleStringValues(openKey mocking.RegistryKey, settings []string, expectedValues []string) []bool {
	results := make([]bool, len(settings))
	for i, val := range settings {
		results[i] = checkStringValue(openKey, val, expectedValues[i])
	}
	return results
}

// checkIntegerRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct integer settings.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]interface{}): A slice of interface values representing the expected values of the registry keys.
//
// Returns:
//
//   - []bool: A slice of boolean values indicating whether the integer settings of the registry keys match the expected values.
func checkIntegerRegistrySettings(registryKey mocking.RegistryKey, registryPath string, settings []string, expectedValues []interface{}) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		return make([]bool, len(settings))
	}
	defer mocking.CloseRegistryKey(key)

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkStringRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct string settings.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]string): A slice of strings representing the expected values of the registry keys.
//
// Returns:
//
//   - []bool: A slice of boolean values indicating whether the string settings of the registry keys match the expected values.
func checkStringRegistrySettings(registryKey mocking.RegistryKey, registryPath string, settings []string, expectedValues []string) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		return make([]bool, len(settings))
	}
	defer mocking.CloseRegistryKey(key)

	return checkMultipleStringValues(key, settings, expectedValues)
}

// checkIntegerStringRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct integer and string settings.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - integerSettings ([]string): A slice of strings representing the names of the integer values to check.
//   - expectedIntegers ([]interface{}): A slice of interface values representing the expected integer values of the registry keys.
//   - stringSettings ([]string): A slice of strings representing the names of the string values to check.
//   - expectedStrings ([]string): A slice of strings representing the expected string values of the registry keys.
//
// Returns:
//
//   - []bool: A slice of boolean values indicating whether the integer and string settings of the registry keys match the expected values.
func checkIntegerStringRegistrySettings(registryKey mocking.RegistryKey, registryPath string,
	integerSettings []string, expectedIntegers []interface{}, stringSettings []string,
	expectedStrings []string) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		return make([]bool, len(integerSettings)+len(stringSettings))
	}
	defer mocking.CloseRegistryKey(key)

	results := make([]bool, 0)
	results = append(results, checkMultipleIntegerValues(key, integerSettings, expectedIntegers)...)
	results = append(results, checkMultipleStringValues(key, stringSettings, expectedStrings)...)
	return results
}
