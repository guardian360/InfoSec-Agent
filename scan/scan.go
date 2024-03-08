package scan

import (
	"InfoSec-Agent/checks"
	"encoding/json"
	"fmt"
	"os"
)

func Scan() {
	// Run all checks
	smb := checks.SmbCheck()
	secureBoot := checks.SecureBoot()
	guest := checks.GuestAccount()
	devices := checks.ExternalDevices()
	sharing := checks.NetworkSharing()
	startup := checks.Startup()
	windowsOutdated := checks.WindowsOutdated()
	loginMethod := checks.LoginMethod()
	lastPasswordChange := checks.LastPasswordChange()
	windowsDefender := checks.WindowsDefender()
	uac := checks.UACCheck()
	remoteDesktop := checks.RemoteDesktopCheck()
	bluetooth := checks.Bluetooth()
	location := checks.Permission("location")
	microphone := checks.Permission("microphone")
	webcam := checks.Permission("webcam")
	appointments := checks.Permission("appointments")
	contacts := checks.Permission("contacts")
	passwordManager := checks.PasswordManager()
	ports := checks.OpenPorts()

	// Combine results
	checkResults := []checks.Check{
		smb,
		secureBoot,
		guest,
		devices,
		sharing,
		startup,
		windowsOutdated,
		loginMethod,
		lastPasswordChange,
		windowsDefender,
		uac,
		remoteDesktop,
		bluetooth,
		location,
		microphone,
		webcam,
		appointments,
		contacts,
		passwordManager,
		ports,
	}

	// Serialize check results to JSON
	jsonData, err := json.MarshalIndent(checkResults, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

	// Write JSON data to a file
	file, err := os.Create("checks.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON data to file:", err)
		return
	}
}
