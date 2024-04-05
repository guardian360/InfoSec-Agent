package checks

import (
	"slices"

	"github.com/getlantern/errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// Startup is a function that checks the Windows registry for startup programs.
//
// Parameters:
//   - key1 registrymock.RegistryKey: A registry key object for accessing the first registry key location for startup programs.
//   - key2 registrymock.RegistryKey: A registry key object for accessing the second registry key location for startup programs.
//   - key3 registrymock.RegistryKey: A registry key object for accessing the third registry key location for startup programs.
//
// Returns:
//   - Check: A struct containing the result of the check. The result includes a list of startup programs if any are found, or a message indicating that no startup programs were found.
//
// The function works by opening three different registry keys where startup programs can be located. It reads the entries within each registry key and concatenates the results. If any startup programs are found, the function returns a Check instance containing a list of the startup programs. If no startup programs are found, the function returns a Check instance with a message indicating that no startup programs were found. If the function encounters an error while opening the registry keys or reading the entries, it returns a Check instance containing an error message.
func Startup(key1 registrymock.RegistryKey, key2 registrymock.RegistryKey, key3 registrymock.RegistryKey) Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := registrymock.OpenRegistryKey(key1,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := registrymock.OpenRegistryKey(key2,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := registrymock.OpenRegistryKey(key3,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)
	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError("Startup", errors.New("error opening registry keys"))
	}

	// Close the keys after we have received all relevant information
	defer registrymock.CloseRegistryKey(cuKey)
	defer registrymock.CloseRegistryKey(lmKey)
	defer registrymock.CloseRegistryKey(lmKey2)

	// Read the entries within the registry key
	cuValueNames, err1 := cuKey.ReadValueNames(0)
	lmValueNames, err2 := lmKey.ReadValueNames(0)
	lm2ValueNames, err3 := lmKey2.ReadValueNames(0)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError("Startup", errors.New("error reading value names"))
	}

	if len(slices.Concat(cuValueNames, lmValueNames, lm2ValueNames)) == 0 {
		return NewCheckResult("Startup", "No startup programs found")
	}

	output := make([]string, 0)
	output = append(output, registrymock.FindEntries(cuValueNames, cuKey)...)
	output = append(output, registrymock.FindEntries(lmValueNames, lmKey)...)
	output = append(output, registrymock.FindEntries(lm2ValueNames, lmKey2)...)

	return NewCheckResult("Startup", output...)
}
