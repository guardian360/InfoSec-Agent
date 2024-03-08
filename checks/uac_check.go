package checks

import (
	"os/exec"
	"strings"
)

func UACCheck() Check {
	// Get the UAC level by performing a command in powershell.
	key, err := exec.Command("powershell", "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin").Output()
	if err != nil {
		return newCheckErrorf("UAC", "error retrieving UAC", err)
	}

	switch strings.TrimSpace(string(key)) {
	case "0":
		return newCheckResult("UAC", "UAC is disabled.")
	case "2":
		return newCheckResult("UAC", "UAC is turned on for apps making changes to your computer and for changing your settings.")
	case "5":
		return newCheckResult("UAC", "UAC is turned on for apps making changes to your computer.")
	default:
		return newCheckResult("UAC", "Unknown UAC level")
	}
}
