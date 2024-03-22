package checks

import (
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"golang.org/x/sys/windows/registry"
)

// Startup checks the registry for startup programs
//
// Parameters: _
//
// Returns: A list of start-up programs
func Startup() Check {
	// Start-up programs can be found in different locations within the registry
	// Both the current user and local machine registry keys are checked
	cuKey, err1 := utils.OpenRegistryKey(registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey, err2 := utils.OpenRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2, err3 := utils.OpenRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError("Startup", fmt.Errorf("error opening registry keys"))
	}

	// Close the keys after we have received all relevant information
	defer utils.CloseRegistryKey(cuKey)
	defer utils.CloseRegistryKey(lmKey)
	defer utils.CloseRegistryKey(lmKey2)

	// Read the entries within the registry key
	cuValueNames, err1 := cuKey.ReadValueNames(0)
	lmValueNames, err2 := lmKey.ReadValueNames(0)
	lm2ValueNames, err3 := lmKey2.ReadValueNames(0)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewCheckError("Startup", fmt.Errorf("error reading value names"))
	}

	output := make([]string, 0)
	output = append(output, utils.FindEntries(cuValueNames, cuKey)...)
	output = append(output, utils.FindEntries(lmValueNames, lmKey)...)
	output = append(output, utils.FindEntries(lm2ValueNames, lmKey2)...)

	return NewCheckResult("Startup", output...)
}
