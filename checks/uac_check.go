// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"os/exec"
	"strings"
	"syscall"
)

// UACCheck checks the User Account Control (UAC) level
//
// Parameters: _
//
// Returns: The level that the UAC is enabled at
func UACCheck() Check {
	// The UAC level can be retrieved as a property from the ConsentPromptBehaviorAdmin
	cmd := exec.Command("powershell", "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows"+
		"\\CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	key, err := cmd.Output()
	if err != nil {
		return newCheckErrorf("UAC", "error retrieving UAC", err)
	}

	// Based on the value of the key, return the appropriate result
	switch strings.TrimSpace(string(key)) {
	case "0":
		return newCheckResult("UAC", "UAC is disabled.")
	case "2":
		return newCheckResult("UAC", "UAC is turned on for apps making changes to your computer and "+
			"for changing your settings.")
	case "5":
		return newCheckResult("UAC", "UAC is turned on for apps making changes to your computer.")
	default:
		return newCheckResult("UAC", "Unknown UAC level")
	}
}
