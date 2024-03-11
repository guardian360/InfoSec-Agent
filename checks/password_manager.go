// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"os"
	"strings"
)

// PasswordManager checks for the presence of known password managers
//
// Parameters: _
//
// Returns: The name of the password manager if found, or "No password manager found" if not found
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
	// List all programs found within the 'Program Files' folder
	programs, err := listInstalledPrograms(programFiles)
	if err != nil {
		return newCheckErrorf("PasswordManager", "error listing installed programs in Program Files", err)
	}

	// Check if any of the listed programs are password managers
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordmanager)) {
				return newCheckResult("PasswordManager", passwordmanager)
			}
		}
	}

	// Check for a password manager within the 'Program Files (x86)' folder
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

// listInstalledPrograms lists the installed programs in a given directory
//
// Parameters: directory (string) representing the directory to check
//
// Returns: A slice of strings containing the names of the installed programs
func listInstalledPrograms(directory string) ([]string, error) {
	var programs []string

	// Read the directory to get a list of files and folders
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// A directory represents an installed program
		if file.IsDir() {
			programs = append(programs, file.Name())
		}
	}

	return programs, nil
}
