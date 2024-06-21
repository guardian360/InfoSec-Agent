package windows

import (
	"errors"
	"strconv"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TODO: Update documentation
// PasswordLength is a function that checks if the minimum password length is at least 15 characters.
// It does this by executing a command to show information about the password complexity settings and parsing the output.
//
// Parameters:
//   - executor mocking.CommandExecutor: An object that implements the CommandExecutor interface.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the password length check.
func PasswordLength(executor mocking.CommandExecutor) checks.Check {
	passwordCommand := "net accounts"
	output, err := executor.Execute("cmd", "/c", passwordCommand)
	if err != nil {
		logger.Log.ErrorWithErr("Error executing password complexity command", err)
		return checks.NewCheckError(checks.PasswordComplexityID, err)
	}

	lines := strings.Split(string(output), "\n")
	// Check that the output has the desired format
	if lines == nil || len(lines) < 4 {
		logger.Log.Error("Error parsing password complexity output")
		return checks.NewCheckError(checks.PasswordComplexityID, errors.New("command output does not have expected structure"))
	}
	// Minimum password length is the fourth entry of the output
	passwordLengthStr := strings.Split(lines[3], ":")[1]
	passwordLengthStr = strings.TrimSpace(passwordLengthStr)
	passwordLength, err := strconv.Atoi(passwordLengthStr)
	if err != nil {
		logger.Log.ErrorWithErr("Error parsing password length", err)
		return checks.NewCheckError(checks.PasswordComplexityID, err)
	}
	if passwordLength < 15 {
		return checks.NewCheckResult(checks.PasswordComplexityID, 1)
	}
	return checks.NewCheckResult(checks.PasswordComplexityID, 0)
}
