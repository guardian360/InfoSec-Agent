package scan

import (
	"os"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/cisregistrysettings"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

var executor = &mocking.RealCommandExecutor{}
var profileFinder = browsers.RealProfileFinder{}
var defaultDirGetter = browsers.RealDefaultDirGetter{}
var copyFileGetter = browsers.RealCopyFileGetter{}
var queryDBGetter = browsers.RealQueryCookieDatabaseGetter{}

const browserChrome = "Chrome"
const browserEdge = "Edge"

// ChecksList is a slice of functions that return checks.Check objects.
// Each function in the slice represents a different security or privacy check that the application can perform.
// When the Scan function is called, it iterates over this slice and executes each check in turn.
// The result of each check is then appended to the checkResults slice, which is returned by the Scan function.
var ChecksList = func() [][]func() checks.Check {
	var checks [][]func() checks.Check
	// Check for the presence of Firefox, Chrome, and Edge profiles. If so, add the corresponding checks
	firefoxDir := GeneratePath("\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles")
	chromeDir := GeneratePath("\\AppData\\Local\\Google\\Chrome\\User Data\\Default")
	edgeDir := GeneratePath("\\AppData\\Local\\Microsoft\\Edge\\User Data\\Default")
	if DirectoryExists(firefoxDir) {
		checks = append(checks, mozillaFirefoxChecks)
	}
	if DirectoryExists(chromeDir) {
		checks = append(checks, googleChromeChecks)
	}
	if DirectoryExists(edgeDir) {
		checks = append(checks, microsoftEdgeChecks)
	}
	checks = append(checks, cisChecks)
	checks = append(checks, devicesChecks)
	checks = append(checks, networkChecks)
	checks = append(checks, programsChecks)
	checks = append(checks, windowsChecks)

	return checks
}()

// googleChromeChecks contains all security/privacy checks that are specific to the Google Chrome browser.
var googleChromeChecks = []func() checks.Check{
	func() checks.Check {
		return chromium.CookiesChromium(browserChrome, defaultDirGetter, copyFileGetter, queryDBGetter)
	},
	func() checks.Check {
		return chromium.ExtensionsChromium("Chrome", defaultDirGetter, chromium.RealExtensionIDGetter{}, chromium.ChromeExtensionNameGetter{})
	},
	func() checks.Check {
		return chromium.HistoryChromium(browserChrome, browsers.RealDefaultDirGetter{}, chromium.RealCopyDBGetter{}, chromium.RealQueryDatabaseGetter{}, chromium.RealProcessQueryResultsGetter{}, browsers.RealPhishingDomainGetter{})
	},
	func() checks.Check {
		return chromium.SearchEngineChromium(browserChrome, false, nil, defaultDirGetter)
	},
}

// microsoftEdgeChecks contains all security/privacy checks that are specific to the Microsoft Edge browser.
var microsoftEdgeChecks = []func() checks.Check{
	func() checks.Check {
		return chromium.CookiesChromium(browserEdge, defaultDirGetter, copyFileGetter, queryDBGetter)
	},
	func() checks.Check {
		return chromium.ExtensionsChromium(browserEdge, defaultDirGetter, chromium.RealExtensionIDGetter{}, chromium.ChromeExtensionNameGetter{})
	},
	func() checks.Check {
		return chromium.HistoryChromium(browserEdge, browsers.RealDefaultDirGetter{}, chromium.RealCopyDBGetter{}, chromium.RealQueryDatabaseGetter{}, chromium.RealProcessQueryResultsGetter{}, browsers.RealPhishingDomainGetter{})
	},
	func() checks.Check {
		return chromium.SearchEngineChromium(browserEdge, false, nil, defaultDirGetter)
	},
}

// mozillaFirefoxChecks contains all security/privacy checks that are specific to the Mozilla Firefox browser.
var mozillaFirefoxChecks = []func() checks.Check{
	func() checks.Check { return firefox.CookiesFirefox(profileFinder, copyFileGetter, queryDBGetter) },
	func() checks.Check { c, _ := firefox.ExtensionFirefox(profileFinder); return c },
	func() checks.Check { _, c := firefox.ExtensionFirefox(profileFinder); return c },
	func() checks.Check {
		return firefox.HistoryFirefox(profileFinder, browsers.RealPhishingDomainGetter{}, firefox.RealQueryDatabaseGetter{}, firefox.RealProcessQueryResultsGetter{}, firefox.RealCopyDBGetter{})
	},
	func() checks.Check { return firefox.SearchEngineFirefox(profileFinder, false, nil, nil) },
}

// cisChecks contains all security/privacy checks that are specific to the CIS benchmark.
var cisChecks = []func() checks.Check{
	func() checks.Check {
		return cisregistrysettings.CISRegistrySettings(mocking.LocalMachine, mocking.UserProfiles)
	},
}

// devicesChecks contains all security/privacy checks that are specific to (external) devices.
var devicesChecks = []func() checks.Check{
	func() checks.Check { return devices.Bluetooth(mocking.LocalMachine) },
	func() checks.Check { return devices.ExternalDevices(executor) },
}

// networkChecks contains all security/privacy checks that are specific to network settings.
var networkChecks = []func() checks.Check{
	func() checks.Check { return network.OpenPorts(executor, executor) },
	func() checks.Check { return network.SmbCheck(executor) },
	func() checks.Check { return network.NetBIOSEnabled(executor) },
	func() checks.Check { return network.WPADEnabled(executor) },
}

// programsChecks contains all security/privacy checks that are specific to installed programs.
var programsChecks = []func() checks.Check{
	func() checks.Check { return programs.PasswordManager(programs.RealProgramLister{}) },
}

// windowsChecks contains all security/privacy checks that are specific to Windows (registry) settings.
var windowsChecks = []func() checks.Check{
	func() checks.Check { return windows.Advertisement(mocking.CurrentUser) },
	func() checks.Check { return windows.AllowRemoteRPC(mocking.LocalMachine) },
	func() checks.Check { return windows.AutomaticLogin(mocking.LocalMachine) },
	func() checks.Check { return windows.Defender(mocking.LocalMachine, mocking.LocalMachine) },
	func() checks.Check { return windows.GuestAccount(executor, executor, executor, executor) },
	func() checks.Check { return windows.LastPasswordChange(executor) },
	func() checks.Check { return windows.LoginMethod(mocking.LocalMachine) },
	func() checks.Check { return windows.Outdated(executor) },
	func() checks.Check {
		return windows.Permission(checks.AppointmentsID, "appointments", mocking.CurrentUser)
	},
	func() checks.Check { return windows.Permission(checks.ContactsID, "contacts", mocking.CurrentUser) },
	func() checks.Check { return windows.Permission(checks.LocationID, "location", mocking.CurrentUser) },
	func() checks.Check { return windows.Permission(checks.MicrophoneID, "microphone", mocking.CurrentUser) },
	func() checks.Check { return windows.Permission(checks.WebcamID, "webcam", mocking.CurrentUser) },
	func() checks.Check { return windows.RemoteDesktopCheck(mocking.LocalMachine) },
	func() checks.Check { return windows.SecureBoot(mocking.LocalMachine) },
	func() checks.Check {
		return windows.Startup(mocking.CurrentUser, mocking.LocalMachine, mocking.LocalMachine)
	},
	func() checks.Check { return windows.UACCheck(executor) },
	func() checks.Check { return windows.FirewallEnabled(executor) },
	func() checks.Check { return windows.PasswordLength(executor) },
	func() checks.Check { return windows.CredentialGuardRunning(executor) },
	func() checks.Check { return windows.ScreenLockEnabled(mocking.CurrentUser) },
}

// DirectoryExists checks if a directory exists at the specified path.
// If the directory exists, the function returns true; otherwise, it returns false.
//
// Parameters:
//   - dirPath (string): The path to the directory to check.
//
// Returns:
//   - bool: A boolean value that indicates whether the directory exists.
func DirectoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GeneratePath generates a path by concatenating the current user's home directory with the specified path.
// If the current user's home directory cannot be determined, the function returns an empty string.
//
// Parameters:
//   - path (string): The path to concatenate with the current user's home directory.
//
// Returns:
//   - string: The generated path.
func GeneratePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDir + path
}
