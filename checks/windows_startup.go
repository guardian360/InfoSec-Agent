package checks

import (
	"slices"

	"github.com/getlantern/errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// Startup checks the registry for startup programs
//
// Parameters: _
//
// Returns: A list of start-up programs
func Startup(key1 mocking.RegistryKey, key2 mocking.RegistryKey, key3 mocking.RegistryKey) Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := mocking.OpenRegistryKey(key1,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := mocking.OpenRegistryKey(key2,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := mocking.OpenRegistryKey(key3,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)
	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError(StartupID, errors.New("error opening registry keys"))
	}

	// Close the keys after we have received all relevant information
	defer mocking.CloseRegistryKey(cuKey)
	defer mocking.CloseRegistryKey(lmKey)
	defer mocking.CloseRegistryKey(lmKey2)

	// Read the entries within the registry key
	cuValueNames, err1 := cuKey.ReadValueNames(0)
	lmValueNames, err2 := lmKey.ReadValueNames(0)
	lm2ValueNames, err3 := lmKey2.ReadValueNames(0)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError(StartupID, errors.New("error reading value names"))
	}

	if len(slices.Concat(cuValueNames, lmValueNames, lm2ValueNames)) == 0 {
		return NewCheckResult(StartupID, 0, "No startup programs found")
	}

	output := make([]string, 0)
	output = append(output, mocking.FindEntries(cuValueNames, cuKey)...)
	output = append(output, mocking.FindEntries(lmValueNames, lmKey)...)
	output = append(output, mocking.FindEntries(lm2ValueNames, lmKey2)...)

	return NewCheckResult(StartupID, 1, output...)
}
