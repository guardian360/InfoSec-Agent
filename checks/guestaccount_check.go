// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"os/exec"
	"strings"
)

// GuestAccount checks if the Windows guest account is active
//
// Parameters: _
//
// Returns: If the guest account is active or not
func GuestAccount() Check {
	// Get localgroup name using GetWmiObject
	output, err := exec.Command("powershell",
		"Get-WmiObject", "Win32_Group", "|", "Select-Object", "SID,Name").Output()

	if err != nil {
		return NewCheckErrorf("Guest account", "error executing command Get-WmiObject", err)
	}
	outputString := strings.Split(string(output), "\r\n")
	found := false
	guestGroup := ""
	for _, line := range outputString {
		// Check for the guest account SID
		if strings.Contains(line, "S-1-5-32-546") {
			line = line[13 : len(line)-1]
			line = strings.TrimSpace(line)
			found = true
			guestGroup = line
		}
	}
	if !found {
		return NewCheckResult("Guest account", "Guest group not found")
	}

	// Get local group members using net localgroup command
	output, err = exec.Command("net", "localgroup", guestGroup).Output()

	if err != nil {
		return NewCheckErrorf("Guest account", "error executing command net localgroup", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	guestUser := ""
	for i := range outputString {
		// Find the line containing the guest account
		if strings.Contains(outputString[i], "-----") {
			guestUser = outputString[i+1]
		}
	}
	if guestUser == "" {
		return NewCheckResult("Guest account", "Guest account not found")
	}

	// Retrieve current username
	currentUser, err := getCurrentUsername()
	if err != nil {
		return NewCheckErrorf("Guest account", "error retrieving current username", err)
	}

	// Retrieve the word for 'yes' from the currentUser language
	output, err = exec.Command("net", "user", currentUser).Output()

	if err != nil {
		return NewCheckErrorf("Guest account", "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	line := strings.Split(outputString[5], " ")
	yesWord := line[len(line)-1]

	// Get all users using net user command
	output, err = exec.Command("net", "user", guestUser).Output()

	if err != nil {
		return NewCheckErrorf("Guest account", "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")

	// Check if the guest account is active based on the presence of the word 'yes' in the user's language
	if strings.Contains(outputString[5], yesWord) {
		return NewCheckResult("Guest account",
			"Guest account is active")
	}

	return NewCheckResult("Guest account", "Guest account is not active")
}
