package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

// SmbCheck checks whether SMB1 and SMB2 are enabled
func SmbCheck() {
	smbEnabled("SMB1")
	smbEnabled("SMB2")
}

// Check whether specified SMB protocol is enabled
func smbEnabled(smb string) {
	// Format command
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)

	output, err := exec.Command("powershell", command).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	outputString := strings.Split(string(output), "\r\n")
	line := strings.TrimSpace(outputString[3])
	if line == "True" {
		fmt.Println(smb, "is enabled")
	} else {
		fmt.Println(smb, "is not enabled")
	}
}
