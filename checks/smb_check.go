package checks

import (
	"fmt"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// SmbCheck checks whether SMB1 (Server Message Block) and SMB2 are enabled
//
// Parameters: _
//
// Returns: If SMB1 and SMB2 are enabled or not
func SmbCheck(smb1executor mocking.CommandExecutor, smb2executor mocking.CommandExecutor) Check {
	var resultID int
	smb1, resultID, err := SmbEnabled("SMB1", smb1executor, resultID)

	if err != nil {
		return NewCheckError(SmbID, err)
	}
	smb2, resultID, err := SmbEnabled("SMB2", smb2executor, resultID)

	if err != nil {
		return NewCheckError(SmbID, err)
	}

	return NewCheckResult(SmbID, resultID, smb1, smb2)
}

// SmbEnabled checks whether the specified SMB protocol is enabled
//
// Parameters: smb (string) represents the SMB protocol to check
//
// Returns: If the specified SMB protocol is enabled or not
func SmbEnabled(smb string, executor mocking.CommandExecutor, resultID int) (string, int, error) {
	// Get the status of the specified SMB protocol
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)
	output, err := executor.Execute("powershell", command)

	if err != nil {
		return "", 0, err
	}

	outputString := strings.Split(string(output), "\r\n")
	line := strings.TrimSpace(outputString[3])
	if line == "True" {
		if smb == "SMB1" {
			resultID++
		}
		if smb == "SMB2" {
			resultID += 2
		}
		return smb + ": enabled", resultID, nil
	}

	return smb + ": not enabled", resultID, nil
}
