package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

// SmbCheck checks whether SMB1 (Server Message Block) and SMB2 are enabled
//
// Parameters: _
//
// Returns: If SMB1 and SMB2 are enabled or not
func SmbCheck() Check {
	smb1, err := smbEnabled("SMB1")

	if err != nil {
		return NewCheckError("smb", err)
	}
	smb2, err := smbEnabled("SMB2")

	if err != nil {
		return NewCheckError("smb", err)
	}

	return NewCheckResult("smb", smb1, smb2)
}

// smbEnabled checks whether the specified SMB protocol is enabled
//
// Parameters: smb (string) represents the SMB protocol to check
//
// Returns: If the specified SMB protocol is enabled or not
func smbEnabled(smb string) (string, error) {
	// Get the status of the specified SMB protocol
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
