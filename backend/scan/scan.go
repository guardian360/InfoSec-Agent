// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"encoding/json"
	"fmt"
	"time"

	apiconnection "github.com/InfoSec-Agent/InfoSec-Agent/backend/api_connection"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/cisregistrysettings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"golang.org/x/sys/windows/registry"

	"github.com/ncruces/zenity"
)

var executor = &mocking.RealCommandExecutor{}

// SecurityChecks is a slice of functions that return checks.Check objects.
// Each function in the slice represents a different security or privacy check that the application can perform.
// When the Scan function is called, it iterates over this slice and executes each check in turn.
// The result of each check is then appended to the checkResults slice, which is returned by the Scan function.
var SecurityChecks = []func() checks.Check{
	func() checks.Check {
		return programs.PasswordManager(programs.RealProgramLister{})
	},
	func() checks.Check {
		return windows.Defender(mocking.LocalMachine, mocking.LocalMachine)
	},
	func() checks.Check {
		return windows.LastPasswordChange(executor)
	},
	func() checks.Check {
		return windows.LoginMethod(mocking.LocalMachine)
	},
	func() checks.Check {
		return windows.Permission(checks.LocationID, "location", mocking.CurrentUser)
	},
	func() checks.Check {
		return windows.Permission(checks.MicrophoneID, "microphone", mocking.CurrentUser)
	},
	func() checks.Check {
		return windows.Permission(checks.WebcamID, "webcam", mocking.CurrentUser)
	},
	func() checks.Check {
		return windows.Permission(checks.AppointmentsID, "appointments", mocking.CurrentUser)
	},
	func() checks.Check {
		return windows.Permission(checks.ContactsID, "contacts", mocking.CurrentUser)
	},
	func() checks.Check {
		return devices.Bluetooth(mocking.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
	},
	func() checks.Check {
		return network.OpenPorts(executor, executor)
	},
	func() checks.Check { return windows.Outdated(executor) },
	func() checks.Check {
		return windows.SecureBoot(mocking.LocalMachine)
	},
	func() checks.Check {
		return network.SmbCheck(executor)
	},
	func() checks.Check {
		return windows.Startup(mocking.CurrentUser, mocking.LocalMachine, mocking.LocalMachine)
	},
	func() checks.Check {
		return windows.GuestAccount(executor, executor,
			executor, executor)
	},
	func() checks.Check { return windows.UACCheck(executor) },
	func() checks.Check {
		return windows.RemoteDesktopCheck(mocking.LocalMachine)
	},
	func() checks.Check { return devices.ExternalDevices(executor) },
	func() checks.Check { return windows.Advertisement(mocking.LocalMachine) },
	func() checks.Check { return chromium.HistoryChromium("Chrome") },
	func() checks.Check { return chromium.ExtensionsChromium("Chrome") },
	func() checks.Check {
		return chromium.SearchEngineChromium("Chrome", false, nil, browsers.RealPreferencesDirGetter{})
	},
	func() checks.Check { c, _ := firefox.ExtensionFirefox(browsers.RealProfileFinder{}); return c },
	func() checks.Check { _, c := firefox.ExtensionFirefox(browsers.RealProfileFinder{}); return c },
	func() checks.Check { return firefox.HistoryFirefox(browsers.RealProfileFinder{}) },
	func() checks.Check {
		return firefox.SearchEngineFirefox(browsers.RealProfileFinder{}, false, nil, nil)
	},
	func() checks.Check {
		return cisregistrysettings.CISRegistrySettings(mocking.LocalMachine, mocking.UserProfiles)
	},
	func() checks.Check {
		return windows.AutomaticLogin(mocking.LocalMachine)
	},
}

// Scan executes all security/privacy checks, serializes the results to JSON, and returns a list of all found issues.
//
// Parameters:
//   - dialog (zenity.ProgressDialog): A dialog window that displays the progress of the scan as it runs.
//
// This function performs the following operations:
//  1. Iterates over each check, displaying the currently running check in the progress dialog and executing the check.
//  2. Appends the result of each check to the checkResults slice.
//  3. Serializes the checkResults slice to JSON and logs the JSON data.
//
// Returns:
//   - []checks.Check: A slice of Check objects representing all found issues.
//   - error: An error object that describes the error (if any) that occurred while running the checks or serializing the results to JSON. If no error occurred, this value is nil.
func Scan(dialog zenity.ProgressDialog) ([]checks.Check, error) {
	date := time.Now().Format(time.RFC3339)
	// TODO: Replace with actual workstation ID and user
	workStationID := 0
	user := "user"
	// Define all security/privacy checks that Scan() should execute
	totalChecks := len(SecurityChecks)

	var checkResults []checks.Check
	// Run all security/privacy checks
	for i, check := range SecurityChecks {
		// Display the currently running check in the progress dialog
		err := dialog.Text(fmt.Sprintf("Running check %d of %d", i+1, totalChecks))
		if err != nil {
			logger.Log.ErrorWithErr("Error setting progress text:", err)
			return checkResults, err
		}

		result := check()
		checkResults = append(checkResults, result)

		// Update the progress bar within the progress dialog
		progress := float64(i+1) / float64(totalChecks) * 100
		err = dialog.Value(int(progress))
		if err != nil {
			logger.Log.ErrorWithErr("Error setting progress value:", err)
			return checkResults, err
		}
	}

	// Serialize check results to JSON
	jsonData, err := json.MarshalIndent(checkResults, "", "  ")
	if err != nil {
		logger.Log.ErrorWithErr("Error marshalling JSON:", err)
		return checkResults, err
	}
	logger.Log.Info(string(jsonData))

	// TODO: Set usersettings.Integration to true depending on whether user has connected with the API
	if usersettings.LoadUserSettings().Integration {
		apiconnection.ParseScanResults(apiconnection.Metadata{WorkStationID: workStationID, User: user, Date: date}, checkResults)
	}
	return checkResults, nil
}
