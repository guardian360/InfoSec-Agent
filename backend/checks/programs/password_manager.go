// Package programs provides functions related to security/privacy checks of installed programs
package programs

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
)

// PasswordManager is a function that checks for the presence of known password managers on the system.
//
// Parameters:
//   - pl (ProgramLister): An instance of ProgramLister used to list installed programs.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the password manager check.
//
// This function uses the ListInstalledPrograms method of the provided ProgramLister to list installed programs in the 'Program Files' and 'Program Files (x86)' directories.
// It then checks if any of the listed programs match the names of known password managers.
// If a match is found, it returns a Check instance with the name of the password manager.
// If no match is found, it returns a Check instance with the message "No password manager found".
func PasswordManager(pl mocking.ProgramLister) checks.Check {
	// List of known password manager registry keys
	passwordManagerNames := []string{
		`LastPass`,
		`1Password`,
		`Dashlane`,
		`enpass`,
		`Bitwarden`,
		`Keeper`,
		`RoboForm`,
		`NordPass`,
		`Sticky Password`,
		`KeePass`,
	}

	programFiles := "C:\\Program Files"
	programFilesx86 := "C:\\Program Files (x86)"
	// List all programs found within the 'Program Files' folder
	programs, err := pl.ListInstalledPrograms(programFiles)
	if err != nil {
		return checks.NewCheckErrorf(checks.PasswordManagerID,
			"error listing installed programs in Program Files", err)
	}

	// Check if any of the listed programs are password managers
	for _, program := range programs {
		for _, passwordManager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordManager)) {
				return checks.NewCheckResult(checks.PasswordManagerID, 0, passwordManager)
			}
		}
	}

	// Check for a password manager within the 'Program Files (x86)' folder
	programs, err = pl.ListInstalledPrograms(programFilesx86)
	if err != nil {
		return checks.NewCheckErrorf(checks.PasswordManagerID,
			"error listing installed programs in Program Files (x86)", err)
	}
	for _, program := range programs {
		for _, passwordManager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordManager)) {
				return checks.NewCheckResult(checks.PasswordManagerID, 0, passwordManager)
			}
		}
	}

	return checks.NewCheckResult(checks.PasswordManagerID, 1)
}
