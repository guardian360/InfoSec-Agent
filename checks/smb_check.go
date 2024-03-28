package checks

import (
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
	"strings"
)

// SmbCheck checks whether SMB1 (Server Message Block) and SMB2 are enabled
//
// Parameters: _
//
// Returns: If SMB1 and SMB2 are enabled or not
func SmbCheck(smb1executor commandmock.CommandExecutor, smb2executor commandmock.CommandExecutor) Check {
	smb1, err := SmbEnabled("SMB1", smb1executor)

	if err != nil {
		return NewCheckError("smb", err)
	}
	smb2, err := SmbEnabled("SMB2", smb2executor)

	if err != nil {
		return NewCheckError("smb", err)
	}

	return NewCheckResult("smb", smb1, smb2)
}

// SmbEnabled checks whether the specified SMB protocol is enabled
//
// Parameters: smb (string) represents the SMB protocol to check
//
// Returns: If the specified SMB protocol is enabled or not
func SmbEnabled(smb string, executor commandmock.CommandExecutor) (string, error) {
	// Get the status of the specified SMB protocol
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)
	output, err := executor.Execute("powershell", command)

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
