package checks

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func Passwordmanager() {
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
		fmt.Println("Error listing installed programs in Program Files:", err)
		return
	}
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), strings.ToLower(passwordmanager)) {
				fmt.Printf("Password manager %s found\n", passwordmanager)
				return
			}
		}
	}

	programs, err = listInstalledPrograms(programFilesx86)
	if err != nil {
		fmt.Println("Error listing installed programs in Program Files (x86):", err)
		return
	}
	for _, program := range programs {
		for _, passwordmanager := range passwordManagerNames {
			if strings.Contains(strings.ToLower(program), passwordmanager) {
				fmt.Printf("Password manager %s found\n", strings.ToLower(passwordmanager))
				return
			}
		}
	}

	fmt.Printf("Password manager not found")
}

func listInstalledPrograms(directory string) ([]string, error) {
	var programs []string

	files, err := ioutil.ReadDir(directory)
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
