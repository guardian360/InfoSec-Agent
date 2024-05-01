package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/mocking"

// CheckWin11 is a function that checks various registry settings specific to the Windows 11 CIS Benchmark Audit list.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_LOCAL_MACHINE for this function.
//
// Returns:
//
//   - []bool: A slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
func CheckWin11(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)

	for _, check := range checksWin11 {
		check(registryKey)
	}

	return results
}

// checksWin11 is a collection of registry checks specific to the Windows 11 CIS Benchmark Audit list.
// Each function in the collection represents a different registry setting check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var checksWin11 = []func(mocking.RegistryKey){
	win11DNSClient,
}

// win11DNSClient is a helper function that checks the registry to determine if the system is configured with the correct settings for the DNS Client.
//
// CIS Benchmark Audit list index: 18.5.4.1
func win11DNSClient(registryKey mocking.RegistryKey) {
	registryPath := DNSClientRegistryPath

	settings := []string{"DoHPolicy"}

	expectedValues := []interface{}{[]uint64{2, 3}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
