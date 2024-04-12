package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// GuestAccount checks if the Windows guest account is active
//
// Parameters: _
//
// Returns: If the guest account is active or not
func GuestAccount(
	executorLocalGroup mocking.CommandExecutor,
	executorLocalGroupMembers mocking.CommandExecutor,
	executorYesWord mocking.CommandExecutor,
	executorNetUser mocking.CommandExecutor,
) Check {
	// Get localgroup name using GetWmiObject
	// output, err := GuestAccountLocalGroup(executorLocalGroup)

	command := "Get-WmiObject Win32_Group | Select-Object SID,Name"
	output, err := executorLocalGroup.Execute("powershell", command)
	if err != nil {
		return NewCheckErrorf(GuestAccountID, "error executing command Get-WmiObject", err)
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
		return NewCheckResult(GuestAccountID, 0, "Guest localgroup not found")
	}

	// Get local group members using net localgroup command
	output, err = executorLocalGroupMembers.Execute("net", "localgroup", guestGroup)
	if err != nil {
		return NewCheckErrorf(GuestAccountID, "error executing command net localgroup", err)
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
		return NewCheckResult(GuestAccountID, 0, "Guest account not found")
	}

	// Retrieve current username
	currentUser, err := utils.CurrentUsername()
	if err != nil {
		return NewCheckErrorf(GuestAccountID, "error retrieving current username", err)
	}

	// Retrieve the word for 'yes' from the currentUser language
	output, err = executorYesWord.Execute("net", "user", currentUser)
	if err != nil {
		return NewCheckErrorf(GuestAccountID, "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	line := strings.Split(outputString[5], " ")
	yesWord := line[len(line)-1]

	// Get all users using net user command
	output, err = executorNetUser.Execute("net", "user", guestUser)
	if err != nil {
		return NewCheckErrorf(GuestAccountID, "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	// Check if the guest account is active based on the presence of the word 'yes' in the user's language
	if strings.Contains(outputString[5], yesWord) {
		return NewCheckResult(GuestAccountID, 1,
			"Guest account is active")
	}
	return NewCheckResult(GuestAccountID, 2, "Guest account is not active")
}
