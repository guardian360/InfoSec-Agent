package windows

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// CredentialGuardRunning is a function that checks if Credential Guard is currently running on the Windows machine.
// Credential Guard is a security feature that helps to protect credentials on a Windows machine from being stolen by malware.
//
// Parameters:
//   - executor (mocking.CommandExecutor): An object that implements the CommandExecutor interface.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the Credential Guard check or an error if one occurred.
//
// The function works by executing a PowerShell command to list all running processes and checking if the lsaiso.exe process is running.
func CredentialGuardRunning(executor mocking.CommandExecutor) checks.Check {
	// Check for the lsaiso.exe process, which is the process that runs Credential Guard.
	output, err := executor.Execute("powershell", "tasklist /fi \"IMAGENAME eq lsaiso.exe\"")
	if err != nil {
		return checks.NewCheckError(checks.CredentialGuardID, err)
	}
	processes := string(output)
	if strings.Contains(processes, "LsaIso") {
		return checks.NewCheckResult(checks.CredentialGuardID, 0)
	}

	return checks.NewCheckResult(checks.CredentialGuardID, 1)
}
