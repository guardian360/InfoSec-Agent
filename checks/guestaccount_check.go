package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

func GuestAccount() {
	// Get localgroup name using GetWmiObject
	output, err := exec.Command("powershell", "Get-WmiObject", "Win32_Group", "|", "Select-Object", "SID,Name").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	outputString := strings.Split(string(output), "\r\n")
	found := false
	guestGroup := ""
	for _, line := range outputString {
		if strings.Contains(line, "S-1-5-32-546") {
			line = line[13 : len(line)-1]
			line = strings.TrimSpace(line)
			found = true
			guestGroup = line
		}
	}
	if !found {
		fmt.Println("Guest group not found")
		return
	}

	// Get local group members using net localgroup command
	output, err = exec.Command("net", "localgroup", guestGroup).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	outputString = strings.Split(string(output), "\r\n")
	guestUser := ""
	for i := range outputString {
		if strings.Contains(outputString[i], "-----") {
			guestUser = outputString[i+1]
		}
	}
	if guestUser == "" {
		fmt.Println("Guest account not found")
		return
	}
	// Retrieve current username
	currentUser, err := getCurrentUsername()
	if err != nil {
		fmt.Println("Error retrieving current username:", err)
		return
	}
	// Retrieve the word for yes from the currentUser language
	output, err = exec.Command("net", "user", currentUser).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	outputString = strings.Split(string(output), "\r\n")
	line := strings.Split(outputString[5], " ")
	yesWord := line[len(line)-1]
	// Get all users using net user command
	output, err = exec.Command("net", "user", guestUser).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	outputString = strings.Split(string(output), "\r\n")
	if strings.Contains(outputString[5], yesWord) {
		fmt.Println("Guest account is active")
		return
	}
	fmt.Println("Guest account is not active")
}
