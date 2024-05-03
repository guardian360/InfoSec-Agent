// Package programs provides functions related to security/privacy checks of installed programs
package programs

import (
	"os"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
)

// ProgramLister is an interface that defines a method for listing installed programs.
//
// The ListInstalledPrograms method takes a directory path as input and returns a slice of strings representing the names of installed programs, or an error if the operation fails.
//
// This interface is used in the PasswordManager function to abstract the operation of listing installed programs, allowing for different implementations that can be swapped out as needed. This is particularly useful for testing, as a mock implementation can be used to simulate different scenarios.
type ProgramLister interface {
	ListInstalledPrograms(directory string) ([]string, error)
}

// RealProgramLister is a struct that implements the ProgramLister interface.
//
// It provides a real-world implementation of the ListInstalledPrograms method, which lists all installed programs in a given directory by reading the directory's contents and returning the names of all subdirectories, which represent installed programs.
//
// This struct is used in the PasswordManager function to list installed programs when checking for the presence of known password managers.
type RealProgramLister struct{}

// ListInstalledPrograms is a method of the RealProgramLister struct that lists all installed programs in a given directory.
//
// Parameters:
//   - directory (string): The path of the directory to list the installed programs from.
//
// Returns:
//   - []string: A slice of strings representing the names of installed programs.
//   - error: An error object that describes the error, if any occurred.
//
// This method reads the contents of the specified directory and returns the names of all subdirectories, which represent installed programs. If an error occurs during the operation, it returns the error.
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

// PasswordManager is a function that checks for the presence of known password managers on the system.
//
// Parameters:
//   - pl (ProgramLister): An instance of ProgramLister used to list installed programs.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the password manager check. The Result field of the Check instance will contain one of the following messages:
//   - The name of the password manager if found.
//   - "No password manager found" if no known password managers are found.
//
// This function uses the ListInstalledPrograms method of the provided ProgramLister to list installed programs in the 'Program Files' and 'Program Files (x86)' directories. It then checks if any of the listed programs match the names of known password managers. If a match is found, it returns a Check instance with the name of the password manager. If no match is found, it returns a Check instance with the message "No password manager found".
func PasswordManager(pl ProgramLister) checks.Check {
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
			if strings.Contains(strings.ToLower(program), passwordManager) {
				return checks.NewCheckResult(checks.PasswordManagerID, 0, passwordManager)
			}
		}
	}

	return checks.NewCheckResult(checks.PasswordManagerID, 1, "No password manager found")
}
