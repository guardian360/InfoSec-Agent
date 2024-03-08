package checks

import (
	"os"
	"strings"
)

func PasswordManager() Check {
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
	programs, err := listInstalledPrograms(programFiles)
	if err != nil {
		return newCheckErrorf("PasswordManager", "error listing installed programs in Program Files", err)
	}

	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordmanager)) {
				return newCheckResult("PasswordManager", passwordmanager)
			}
		}
	}

	programs, err = listInstalledPrograms(programFilesx86)
	if err != nil {
		return newCheckErrorf("PasswordManager", "error listing installed programs in Program Files (x86)", err)
	}
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), passwordmanager) {
				return newCheckResult("PasswordManager", passwordmanager)
			}
		}
	}

	return newCheckResult("PasswordManager", "No password manager found")
}

func listInstalledPrograms(directory string) ([]string, error) {
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
