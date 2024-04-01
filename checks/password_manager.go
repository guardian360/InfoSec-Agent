package checks

import (
	"os"
	"strings"
)

type ProgramLister interface {
	ListInstalledPrograms(directory string) ([]string, error)
}

type RealProgramLister struct{}

func (rpl RealProgramLister) ListInstalledPrograms(directory string) ([]string, error) {
	var programs []string
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			programs = append(programs, file.Name())
		}
	}
	return programs, nil
}

// PasswordManager checks for the presence of known password managers
//
// Parameters: _
//
// Returns: The name of the password manager if found, or "No password manager found" if not found
func PasswordManager(pl ProgramLister) Check {
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
		return NewCheckErrorf("PasswordManager",
			"error listing installed programs in Program Files", err)
	}

	// Check if any of the listed programs are password managers
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordmanager)) {
				return NewCheckResult("PasswordManager", passwordmanager)
			}
		}
	}

	// Check for a password manager within the 'Program Files (x86)' folder
	programs, err = pl.ListInstalledPrograms(programFilesx86)
	if err != nil {
		return NewCheckErrorf("PasswordManager",
			"error listing installed programs in Program Files (x86)", err)
	}
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), passwordmanager) {
				return NewCheckResult("PasswordManager", passwordmanager)
			}
		}
	}

	return NewCheckResult("PasswordManager", "No password manager found")
}
