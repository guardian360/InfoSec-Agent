package windows

import (
	"errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// Startup is a function that checks the Windows registry for startup programs.
//
// Parameters:
//   - key1 mocking.RegistryKey: A registry key object for accessing the first registry key location for startup programs.
//   - key2 mocking.RegistryKey: A registry key object for accessing the second registry key location for startup programs.
//   - key3 mocking.RegistryKey: A registry key object for accessing the third registry key location for startup programs.
//
// Returns:
//   - Check: A struct containing the result of the check. The result includes a list of startup programs if any are found, or a message indicating that no startup programs were found.
//
// The function works by opening three different registry keys where startup programs can be located. It reads the entries within each registry key and concatenates the results. If any startup programs are found, the function returns a Check instance containing a list of the startup programs. If no startup programs are found, the function returns a Check instance with a message indicating that no startup programs were found. If the function encounters an error while opening the registry keys or reading the entries, it returns a Check instance containing an error message.
func Startup(key1 mocking.RegistryKey, key2 mocking.RegistryKey, key3 mocking.RegistryKey) checks.Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := checks.OpenRegistryKey(key1,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := checks.OpenRegistryKey(key2,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := checks.OpenRegistryKey(key3,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)
	if err1 != nil || err2 != nil || err3 != nil {
		return checks.NewCheckError(checks.StartupID, errors.New("error opening registry keys"))
	}

	// Close the keys after we have received all relevant information
	defer checks.CloseRegistryKey(cuKey)
	defer checks.CloseRegistryKey(lmKey)
	defer checks.CloseRegistryKey(lmKey2)

	// Read the entries within the registry key
	cuValueNames, err1 := cuKey.ReadValueNames(0)
	lmValueNames, err2 := lmKey.ReadValueNames(0)
	lm2ValueNames, err3 := lmKey2.ReadValueNames(0)

	if err1 != nil || err2 != nil || err3 != nil {
		return checks.NewCheckError(checks.StartupID, errors.New("error reading value names"))
	}

	output := make([]string, 0)
	output = append(output, FindEntries(cuValueNames, cuKey)...)
	output = append(output, FindEntries(lmValueNames, lmKey)...)
	output = append(output, FindEntries(lm2ValueNames, lmKey2)...)

	if len(output) == 0 {
		return checks.NewCheckResult(checks.StartupID, 0)
	}
	return checks.NewCheckResult(checks.StartupID, 1, output...)
}

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
