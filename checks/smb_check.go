package checks

import (
	"fmt"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// SmbCheck checks whether SMB1 (Server Message Block) and SMB2 are enabled
//
// Parameters: _
//
// Returns: If SMB1 and SMB2 are enabled or not
func SmbCheck(smb1executor commandmock.CommandExecutor, smb2executor commandmock.CommandExecutor) Check {
	var resultInt int
	smb1, err := SmbEnabled("SMB1", smb1executor, resultInt)

	if err != nil {
		return NewCheckError(12, err)
	}
	smb2, err := SmbEnabled("SMB2", smb2executor, resultInt)

	if err != nil {
		return NewCheckError(12, err)
	}

	return NewCheckResult(12, resultInt, smb1, smb2)
}

// SmbEnabled checks whether the specified SMB protocol is enabled
//
// Parameters: smb (string) represents the SMB protocol to check
//
// Returns: If the specified SMB protocol is enabled or not
func SmbEnabled(smb string, executor commandmock.CommandExecutor, resultID int) (string, error) {
	// Get the status of the specified SMB protocol
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)
	output, err := executor.Execute("powershell", command)

	if err != nil {
		return "", err
	}

	outputString := strings.Split(string(output), "\r\n")
	line := strings.TrimSpace(outputString[3])
	if line == "True" {
		if smb == "SMB1" {
			resultID += 1
		}
		if smb == "SMB2" {
			resultID += 2
		}
		return smb + ": enabled", nil
	}

	return smb + ": not enabled", nil
}
