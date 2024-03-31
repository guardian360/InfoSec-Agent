package checks

import (
	"fmt"
	"slices"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// Startup checks the registry for startup programs
//
// Parameters: _
//
// Returns: A list of start-up programs
func Startup() Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := registrymock.OpenRegistryKey(registrymock.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := registrymock.OpenRegistryKey(registrymock.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := registrymock.OpenRegistryKey(registrymock.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError("Startup", fmt.Errorf("error opening registry keys"))
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
		return NewCheckError("Startup", fmt.Errorf("error reading value names"))
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
