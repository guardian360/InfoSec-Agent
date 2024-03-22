package checks

import (
	"os/exec"
	"strings"
)

// UACCheck checks the User Account Control (UAC) level
//
// Parameters: _
//
// Returns: The level that the UAC is enabled at
func UACCheck() Check {
	// The UAC level can be retrieved as a property from the ConsentPromptBehaviorAdmin
	key, err := exec.Command("powershell", "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows"+
		"\\CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin").Output()

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
