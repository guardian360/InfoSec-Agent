package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// UACCheck checks the User Account Control (UAC) level
//
// Parameters: _
//
// Returns: The level that the UAC is enabled at
func UACCheck(uacExecutor mocking.CommandExecutor) Check {
	// The UAC level can be retrieved as a property from the ConsentPromptBehaviorAdmin
	command := "powershell"
	key, err := uacExecutor.Execute(command, "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\"+
		"CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin")

	if err != nil {
		return NewCheckErrorf(UacID, "error retrieving UAC", err)
	}

	// Based on the value of the key, return the appropriate result
	switch strings.TrimSpace(string(key)) {
	case "0":
		return NewCheckResult(UacID, 0, "UAC is disabled.")
	case "2":
		return NewCheckResult(UacID, 1, "UAC is turned on for apps making changes to your computer and "+
			"for changing your settings.")
	case "5":
		return NewCheckResult(UacID, 2, "UAC is turned on for apps making changes to your computer.")
	default:
		return NewCheckResult(UacID, 3, "Unknown UAC level")
	}
}
