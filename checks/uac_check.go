package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// UACCheck is a function that checks the User Account Control (UAC) level on the system.
//
// Parameters:
//   - uacExecutor commandmock.CommandExecutor: An executor to run the command for checking the UAC level.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates the level at which the UAC is enabled.
//
// The function works by executing a PowerShell command to get the 'ConsentPromptBehaviorAdmin' property from the system registry. This property represents the UAC level. The function then parses the output of the command to determine the UAC level. Based on the value of the key, the function returns a Check instance containing a string that describes the UAC level.
func UACCheck(uacExecutor commandmock.CommandExecutor) Check {
	// The UAC level can be retrieved as a property from the ConsentPromptBehaviorAdmin
	command := "powershell"
	key, err := uacExecutor.Execute(command, "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\"+
		"CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin")

	if err != nil {
		return NewCheckErrorf("UAC", "error retrieving UAC", err)
	}

	// Based on the value of the key, return the appropriate result
	switch strings.TrimSpace(string(key)) {
	case "0":
		return NewCheckResult("UAC", "UAC is disabled.")
	case "2":
		return NewCheckResult("UAC", "UAC is turned on for apps making changes to your computer and "+
			"for changing your settings.")
	case "5":
		return NewCheckResult("UAC", "UAC is turned on for apps making changes to your computer.")
	default:
		return NewCheckResult("UAC", "Unknown UAC level")
	}
}
