// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"InfoSec-Agent/checks"
	"InfoSec-Agent/checks/browsers/chrome"
	"InfoSec-Agent/checks/browsers/firefox"
	"encoding/json"
	"fmt"
)

// Scan runs all security/privacy checks and serializes the results to JSON.
//
// Parameters: _
//
// Returns: checks.json file containing the results of all security/privacy checks
func Scan() {
	// Run all checks
	passwordManager := checks.PasswordManager()
	windowsDefender := checks.WindowsDefender()
	lastPasswordChange := checks.LastPasswordChange()
	loginMethod := checks.LoginMethod()
	location := checks.Permission("location")
	microphone := checks.Permission("microphone")
	webcam := checks.Permission("webcam")
	appointments := checks.Permission("appointments")
	contacts := checks.Permission("contacts")
	bluetooth := checks.Bluetooth()
	ports := checks.OpenPorts()
	windowsOutdated := checks.WindowsOutdated()
	secureBoot := checks.SecureBoot()
	smb := checks.SmbCheck()
	startup := checks.Startup()
	guest := checks.GuestAccount()
	uac := checks.UACCheck()
	remoteDesktop := checks.RemoteDesktopCheck()
	devices := checks.ExternalDevices()
	sharing := checks.NetworkSharing()
	//cookieFF := firefox.CookieFirefox()
	extensionFF, adblockFF := firefox.ExtensionFirefox()
	historyFF := firefox.HistoryFirefox()
	historyChrome := chrome.HistoryChrome()
	extensionChrome := chrome.ExtensionsChrome()

	// Combine results
	checkResults := []checks.Check{
		passwordManager,
		windowsDefender,
		lastPasswordChange,
		loginMethod,
		location,
		microphone,
		webcam,
		appointments,
		contacts,
		bluetooth,
		ports,
		windowsOutdated,
		secureBoot,
		smb,
		startup,
		guest,
		uac,
		remoteDesktop,
		devices,
		sharing,
		//cookieFF,
		extensionFF,
		adblockFF,
		historyFF,
		extensionChrome,
		historyChrome,
	}

	// Serialize check results to JSON
	jsonData, err := json.MarshalIndent(checkResults, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

	//// Write JSON data to a file
	//file, err := os.Create("checks.json")
	//if err != nil {
	//	fmt.Println("Error creating file:", err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = file.Write(jsonData)
	//if err != nil {
	//	fmt.Println("Error writing JSON data to file:", err)
	//	return
	//}
}
