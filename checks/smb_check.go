package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

// SmbCheck checks whether SMB1 and SMB2 are enabled
func SmbCheck() Check {
	smb1, err := smbEnabled("SMB1")

	if err != nil {
		return newCheckError("smb", err)
	}
	smb2, err := smbEnabled("SMB2")

	if err != nil {
		return newCheckError("smb", err)
	}

	return newCheckResult("smb", []string{smb1, smb2})
}

// Check whether specified SMB protocol is enabled
func smbEnabled(smb string) (string, error) {
	// Format command
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)

	output, err := exec.Command("powershell", command).Output()
	if err != nil {
		return "", err
	}

	outputString := strings.Split(string(output), "\r\n")
	line := strings.TrimSpace(outputString[3])
	if line == "True" {
		return smb + ": enabled", nil
	}

	return smb + ": not enabled", nil
}
