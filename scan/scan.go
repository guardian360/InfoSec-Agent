// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"encoding/json"
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"github.com/InfoSec-Agent/InfoSec-Agent/windowsmock"
	"golang.org/x/sys/windows/registry"

	"github.com/ncruces/zenity"
)

// Scan runs all security/privacy checks and serializes the results to JSON.
//
// Parameters: dialog (zenity.ProgressDialog)
// represents the progress dialog window which is displayed while the scan is running
//
// Returns: checks.Check list containing all found issues
func Scan(dialog zenity.ProgressDialog) ([]checks.Check, error) {
	// Define all security/privacy checks that Scan() should execute
	securityChecks := []func() checks.Check{
		func() checks.Check {
			return checks.PasswordManager(checks.RealProgramLister{})
		},
		func() checks.Check {
			return checks.WindowsDefender(registrymock.LocalMachine, registrymock.LocalMachine)
		},
		func() checks.Check {
			return checks.LastPasswordChange(&commandmock.RealCommandExecutor{})
		},
		func() checks.Check {
			return checks.LoginMethod(registrymock.LocalMachine)
		},
		func() checks.Check {
			return checks.Permission("location", registrymock.CurrentUser)
		},
		func() checks.Check {
			return checks.Permission("microphone", registrymock.CurrentUser)
		},
		func() checks.Check {
			return checks.Permission("webcam", registrymock.CurrentUser)
		},
		func() checks.Check {
			return checks.Permission("appointments", registrymock.CurrentUser)
		},
		func() checks.Check {
			return checks.Permission("contacts", registrymock.CurrentUser)
		},
		func() checks.Check {
			return checks.Bluetooth(registrymock.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
		},
		func() checks.Check {
			return checks.OpenPorts(&commandmock.RealCommandExecutor{}, &commandmock.RealCommandExecutor{})
		},
		func() checks.Check { return checks.WindowsOutdated(&windowsmock.RealWindowsVersion{}) },
		func() checks.Check {
			return checks.SecureBoot(registrymock.LocalMachine)
		},
		func() checks.Check {
			return checks.SmbCheck(&commandmock.RealCommandExecutor{}, &commandmock.RealCommandExecutor{})
		},
		func() checks.Check {
			return checks.Startup(registrymock.CurrentUser, registrymock.LocalMachine, registrymock.LocalMachine)
		},
		func() checks.Check {
			return checks.GuestAccount(&commandmock.RealCommandExecutor{}, &commandmock.RealCommandExecutor{},
				&commandmock.RealCommandExecutor{}, &commandmock.RealCommandExecutor{})
		},
		func() checks.Check { return checks.UACCheck(&commandmock.RealCommandExecutor{}) },
		func() checks.Check {
			return checks.RemoteDesktopCheck(registrymock.LocalMachine)
		},
		func() checks.Check { return checks.ExternalDevices(&commandmock.RealCommandExecutor{}) },
		func() checks.Check { return checks.NetworkSharing(&commandmock.RealCommandExecutor{}) },
		func() checks.Check { return chromium.HistoryChromium("Chrome") },
		func() checks.Check { return chromium.ExtensionsChromium("Chrome") },
		func() checks.Check { return chromium.SearchEngineChromium("Chrome") },
		func() checks.Check { c, _ := firefox.ExtensionFirefox(); return c },
		func() checks.Check { _, c := firefox.ExtensionFirefox(); return c },
		firefox.HistoryFirefox,
		firefox.SearchEngineFirefox,
	}
	totalChecks := len(securityChecks)

	var checkResults []checks.Check
	// Run all security/privacy checks
	for i, check := range securityChecks {
		// Display the currently running check in the progress dialog
		err := dialog.Text(fmt.Sprintf("Running check %d of %d", i+1, totalChecks))
		if err != nil {
			logger.Log.Println("Error setting progress text:", err)
			return checkResults, err
		}

		result := check()
		checkResults = append(checkResults, result)

		// Update the progress bar within the progress dialog
		progress := float64(i+1) / float64(totalChecks) * 100
		err = dialog.Value(int(progress))
		if err != nil {
			logger.Log.Println("Error setting progress value:", err)
			return checkResults, err
		}
	}

	// Serialize check results to JSON
	jsonData, err := json.MarshalIndent(checkResults, "", "  ")
	if err != nil {
		logger.Log.Println("Error marshalling JSON:", err)
		return checkResults, err
	}
	logger.Log.Println(string(jsonData))

	return checkResults, nil

	//// Write JSON data to a file
	// file, err := os.Create("checks.json")
	// if err != nil {
	//	fmt.Println("Error creating file:", err)
	//	return
	//}
	// defer file.Close()
	//
	// _, err = file.Write(jsonData)
	// if err != nil {
	//	fmt.Println("Error writing JSON data to file:", err)
	//	return
	//}
}
