package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

func Uac_check() {
	// Get the UAC level by performing a command in powershell.
	key, err := exec.Command("powershell", "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System').ConsentPromptBehaviorAdmin").Output()
	if err != nil {
		fmt.Println("Error no UAC:", err)
		return
	}

	switch strings.TrimSpace(string(key)) {
	case "0":
		fmt.Println("UAC is disabled.")
	case "2":
		fmt.Println("UAC is turned on for apps making changes to your computer and for changing your settings.")
	case "5":
		fmt.Println("UAC is turned on for apps making changes to your computer.")
	default:
		fmt.Println("Error reading UAC")
	}
}
