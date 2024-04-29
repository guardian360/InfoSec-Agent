// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"encoding/json"
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks/cisregistrysettings"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browserutils"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"golang.org/x/sys/windows/registry"

	"github.com/ncruces/zenity"
)

// securityChecks is a slice of functions that return checks.Check objects.
// Each function in the slice represents a different security or privacy check that the application can perform.
// When the Scan function is called, it iterates over this slice and executes each check in turn.
// The result of each check is then appended to the checkResults slice, which is returned by the Scan function.
var securityChecks = []func() checks.Check{
	func() checks.Check {
		return checks.PasswordManager(checks.RealProgramLister{})
	},
	func() checks.Check {
		return checks.WindowsDefender(mocking.LocalMachine, mocking.LocalMachine)
	},
	func() checks.Check {
		return checks.LastPasswordChange(&mocking.RealCommandExecutor{})
	},
	func() checks.Check {
		return checks.LoginMethod(mocking.LocalMachine)
	},
	func() checks.Check {
		return checks.Permission(checks.LocationID, "location", mocking.CurrentUser)
	},
	func() checks.Check {
		return checks.Permission(checks.MicrophoneID, "microphone", mocking.CurrentUser)
	},
	func() checks.Check {
		return checks.Permission(checks.WebcamID, "webcam", mocking.CurrentUser)
	},
	func() checks.Check {
		return checks.Permission(checks.AppointmentsID, "appointments", mocking.CurrentUser)
	},
	func() checks.Check {
		return checks.Permission(checks.ContactsID, "contacts", mocking.CurrentUser)
	},
	func() checks.Check {
		return checks.Bluetooth(mocking.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
	},
	func() checks.Check {
		return checks.OpenPorts(&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{})
	},
	func() checks.Check { return checks.WindowsOutdated(&mocking.RealCommandExecutor{}) },
	func() checks.Check {
		return checks.SecureBoot(mocking.LocalMachine)
	},
	func() checks.Check {
		return checks.SmbCheck(&mocking.RealCommandExecutor{})
	},
	func() checks.Check {
		return checks.Startup(mocking.CurrentUser, mocking.LocalMachine, mocking.LocalMachine)
	},
	func() checks.Check {
		return checks.GuestAccount(&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{},
			&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{})
	},
	func() checks.Check { return checks.UACCheck(&mocking.RealCommandExecutor{}) },
	func() checks.Check {
		return checks.RemoteDesktopCheck(mocking.LocalMachine)
	},
	func() checks.Check { return checks.ExternalDevices(&mocking.RealCommandExecutor{}) },
	func() checks.Check { return checks.NetworkProfileTypes(mocking.LocalMachine) },
	func() checks.Check { return chromium.HistoryChromium("Chrome") },
	func() checks.Check { return chromium.ExtensionsChromium("Chrome") },
	func() checks.Check { return chromium.SearchEngineChromium("Chrome") },
	func() checks.Check { c, _ := firefox.ExtensionFirefox(); return c },
	func() checks.Check { _, c := firefox.ExtensionFirefox(); return c },
	func() checks.Check { return firefox.HistoryFirefox(browserutils.RealProfileFinder{}) },
	firefox.SearchEngineFirefox,
	func() checks.Check {
		return cisregistrysettings.CISRegistrySettings(mocking.LocalMachine, mocking.UserProfiles)
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
	// Define all security/privacy checks that Scan() should execute
	totalChecks := len(securityChecks)

	var checkResults []checks.Check
	// Run all security/privacy checks
	for i, check := range securityChecks {
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

	return checkResults, nil
}
