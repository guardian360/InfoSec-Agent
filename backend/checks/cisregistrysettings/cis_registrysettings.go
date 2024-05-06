// Package cisregistrysettings provides a set of functions to check various registry settings
// to ensure they adhere to the CIS Benchmark standards. Each function takes a RegistryKey object
// as an argument, which represents the root key from which the registry settings will be checked.
// The functions return a slice of boolean values, where each boolean represents whether a particular
// registry setting adheres to the CIS Benchmark standards.
package cisregistrysettings

import (
	"slices"
	"sort"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

const DNSClientRegistryPath = `SOFTWARE\Policies\Microsoft\Windows NT\DNSClient`

// RegistrySettingsMap is a map that stores the results of the registry settings checks.
// The key is the registry path and the value is a boolean indicating whether the registry setting adheres to the CIS Benchmark standards.
var RegistrySettingsMap = map[string]bool{}

// CISRegistrySettings is a function that checks various registry settings to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked.
//
// Returns:
//   - checks.Check: A check object containing the settings that do not adhere to the CIS Benchmark standards.
func CISRegistrySettings(localMachineKey mocking.RegistryKey, usersKey mocking.RegistryKey) checks.Check {
	// Following function(s) need the HKEY_LOCAL_MACHINE registry key
	CheckServices(localMachineKey)
	CheckPoliciesHKLM(localMachineKey)
	CheckOtherRegistrySettings(localMachineKey)
	// Following function(s) need the HKEY_USERS registry key
	CheckPoliciesHKU(usersKey)

	// Following function(s) need the HKEY_LOCAL_MACHINE registry key
	if windows.WinVersion == 10 {
		CheckWin10(localMachineKey)
	}
	if windows.WinVersion == 11 {
		CheckWin11(localMachineKey)
	}
	// Check if all registry settings adhere to the CIS Benchmark standards
	fullyTrue, incorrectSettings := getIncorrectSettings()
	if fullyTrue {
		return checks.NewCheckResult(checks.CISRegistrySettingsID, 1, "All registry settings adhere to the CIS Benchmark standards")
	}
	resultString := "Not all registry settings adhere to the CIS Benchmark standards"
	return checks.NewCheckResult(checks.CISRegistrySettingsID, 0, append([]string{resultString}, incorrectSettings...)...)
}

// CheckIntegerValue is a helper function that checks if the integer value of a registry key matches the expected value.
//
// Parameters:
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - value (string): The name of the value to check.
//   - expected (interface{}): The expected value of the registry key.
//
// Returns:
//   - bool: A boolean value indicating whether the integer value of the registry key matches the expected value.
func CheckIntegerValue(openKey mocking.RegistryKey, value string, expected interface{}) bool {
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
		// Slice of exactly 2 uint64 values, check if registry value is within the range
		if len(v) == 2 {
			return val >= v[0] && val <= v[1]
		}
		for _, i := range v {
			if val == i {
				return true
			}
		}
	default:
		return false
	}
	return false
}

// CheckStringValue is a helper function that checks if the string value of a registry key matches the expected value.
//
// Parameters:
//   - openKey (mocking.RegistryKey): The registry key to check.
//   - value (string): The name of the value to check.
//   - expected (string): The expected value of the registry key.
//
// Returns:
//   - bool: A boolean value indicating whether the string value of the registry key matches the expected value.
func CheckStringValue(openKey mocking.RegistryKey, value string, expected string) bool {
	val, _, err := openKey.GetStringValue(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading registry value of "+value, err)
		return false
	}
	return val == expected
}

// OpenRegistryKeyWithErrHandling is a helper function that opens a registry key and handles any errors that occur.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The registry key to open.
//   - path (string): The path of the registry key to open.
//
// Returns:
//   - mocking.RegistryKey: The opened registry key.
//   - error: An error object that describes the error (if any) that occurred while opening the registry key. If no error occurred, this value is nil.
func OpenRegistryKeyWithErrHandling(registryKey mocking.RegistryKey, path string) (mocking.RegistryKey, error) {
	key, err := checks.OpenRegistryKey(registryKey, path)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening registry key for CIS Audit list", err)
	}
	return key, err
}

// CheckIntegerRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct integer settings.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]interface{}): A slice of interface values representing the expected values of the registry keys.
//
// Returns: None
func CheckIntegerRegistrySettings(registryKey mocking.RegistryKey, registryPath string, settings []string, expectedValues []interface{}) {
	key, err := OpenRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		for _, setting := range settings {
			RegistrySettingsMap[registryPath+"\\"+setting] = false
			return
		}
	}
	defer checks.CloseRegistryKey(key)

	for i, setting := range settings {
		RegistrySettingsMap[registryPath+"\\"+setting] = CheckIntegerValue(key, setting, expectedValues[i])
	}
}

// CheckStringRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct string settings.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - settings ([]string): A slice of strings representing the names of the values to check.
//   - expectedValues ([]string): A slice of strings representing the expected values of the registry keys.
//
// Returns: None
func CheckStringRegistrySettings(registryKey mocking.RegistryKey, registryPath string, settings []string, expectedValues []string) {
	key, err := OpenRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		for _, setting := range settings {
			RegistrySettingsMap[registryPath+"\\"+setting] = false
			return
		}
	}
	defer checks.CloseRegistryKey(key)

	for i, setting := range settings {
		RegistrySettingsMap[registryPath+"\\"+setting] = CheckStringValue(key, setting, expectedValues[i])
	}
}

// CheckIntegerStringRegistrySettings is a helper function that checks the registry to determine if the system is configured with the correct integer and string settings.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The registry key used to access the system's registry.
//   - registryPath (string): The path to the registry key to check.
//   - integerSettings ([]string): A slice of strings representing the names of the integer values to check.
//   - expectedIntegers ([]interface{}): A slice of interface values representing the expected integer values of the registry keys.
//   - stringSettings ([]string): A slice of strings representing the names of the string values to check.
//   - expectedStrings ([]string): A slice of strings representing the expected string values of the registry keys.
//
// Returns: None
func CheckIntegerStringRegistrySettings(registryKey mocking.RegistryKey, registryPath string,
	integerSettings []string, expectedIntegers []interface{}, stringSettings []string,
	expectedStrings []string) {
	key, err := OpenRegistryKeyWithErrHandling(registryKey, registryPath)
	if err != nil {
		for _, setting := range integerSettings {
			RegistrySettingsMap[registryPath+"\\"+setting] = false
		}
		for _, setting := range stringSettings {
			RegistrySettingsMap[registryPath+"\\"+setting] = false
		}
		return
	}
	defer checks.CloseRegistryKey(key)

	for i, integerSetting := range integerSettings {
		RegistrySettingsMap[registryPath+"\\"+integerSetting] = CheckIntegerValue(key, integerSetting, expectedIntegers[i])
	}
	for i, stringSetting := range stringSettings {
		RegistrySettingsMap[registryPath+"\\"+stringSetting] = CheckStringValue(key, stringSetting, expectedStrings[i])
	}
}

// getIncorrectSettings iterates over the RegistrySettingsMap and returns a slice of incorrect settings.
//
// Parameters: None
//
// Returns:
//   - bool: A boolean value indicating whether all registry settings adhere to the CIS Benchmark standards.
//   - []string: A slice of strings representing the incorrect registry settings.
func getIncorrectSettings() (bool, []string) {
	var incorrectSettings []string
	fullyTrue := true
	for key, value := range RegistrySettingsMap {
		if !value {
			fullyTrue = false
			// If a registry setting does not adhere to the CIS Benchmark standards, store the setting to be returned
			incorrectSettings = append(incorrectSettings, key)
		}
	}
	slices.Sort(incorrectSettings)
	incorrectSettings = trimCommonPrefixes(incorrectSettings, commonPrefixes)
	return fullyTrue, incorrectSettings
}

// commonPrefixes is a slice of common prefixes of registry paths that are checked.
// These prefixes are used to group the registry paths in the output.
var commonPrefixes = []string{
	"SOFTWARE\\Microsoft\\",
	"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon\\",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\Explorer\\",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System\\",
	"SOFTWARE\\Policies\\",
	"SOFTWARE\\Policies\\Microsoft\\",
	"SOFTWARE\\Policies\\Microsoft\\Power\\PowerSettings\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows Defender\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows Defender\\Windows Defender Exploit Guard\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows Defender\\Windows Defender Exploit Guard\\ASR\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows NT\\Printers\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows NT\\Terminal Services\\",
	"SOFTWARE\\Policies\\Microsoft\\WindowsFirewall\\DomainProfile\\",
	"SOFTWARE\\Policies\\Microsoft\\WindowsFirewall\\PrivateProfile\\",
	"SOFTWARE\\Policies\\Microsoft\\WindowsFirewall\\PublicProfile\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\CloudContent\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\DataCollection\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\EventLog\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\Installer\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\Network Connections\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\NetworkProvider\\HardenedPaths\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\Powershell\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\System\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\WcmSvc\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\WinRM\\",
	"SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate\\",
	"SYSTEM\\CurrentControlSet\\",
	"SYSTEM\\CurrentControlSet\\Control\\",
	"SYSTEM\\CurrentControlSet\\Control\\Lsa\\",
	"SYSTEM\\CurrentControlSet\\Control\\SecurePipeServers\\Winreg\\",
	"SYSTEM\\CurrentControlSet\\Control\\Session Manager\\",
	"SYSTEM\\CurrentControlSet\\Services\\",
	"SYSTEM\\CurrentControlSet\\Services\\LanmanServer\\Parameters\\",
	"SYSTEM\\CurrentControlSet\\Services\\NetBT\\Parameters\\",
	"SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters\\",
}

// trimCommonPrefixes trims the common prefixes from the paths and returns the trimmed paths.
//
// Parameters:
//   - paths ([]string): A slice of strings representing the paths to trim.
//   - commonPrefixes ([]string): A slice of strings representing the common prefixes to trim.
//
// Returns:
//   - []string: A slice of strings representing the trimmed paths.
func trimCommonPrefixes(paths []string, commonPrefixes []string) []string {
	// Sort commonPrefixes in descending order of length, so that the longest prefix is trimmed first
	sort.Slice(commonPrefixes, func(i, j int) bool {
		return len(commonPrefixes[i]) > len(commonPrefixes[j])
	})

	var result []string

	for _, prefix := range commonPrefixes {
		var trimmedPaths []string
		for i, path := range paths {
			if strings.HasPrefix(path, prefix) {
				paths[i] = strings.TrimPrefix(path, prefix)
				trimmedPaths = append(trimmedPaths, paths[i])
			}
		}
		if len(trimmedPaths) > 0 {
			// Store the prefix that was trimmed, followed by the trimmed paths
			result = append(result, prefix)
			result = append(result, trimmedPaths...)
		}
	}

	return result
}
