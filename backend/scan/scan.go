// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"encoding/json"
	"fmt"
	"time"

	apiconnection "github.com/InfoSec-Agent/InfoSec-Agent/backend/api_connection"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/ncruces/zenity"
)

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
	totalChecks := len(ChecksList)

	var checkResults []checks.Check
	// Run all security/privacy checks
	for i, check := range ChecksList {
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
