// Package scan collects all different privacy/security checks and provides a function that runs them all.
//
// Exported function(s): Scan
package scan

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
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
	totalChecks := 0
	for _, checkSlice := range ChecksList {
		totalChecks += len(checkSlice)
	}

	// Define a channel which serves as a queue for the checks to be executed
	checksChan := make(chan func() checks.Check, totalChecks)
	// Determine the amount of workers to use for concurrent execution of the checks based on the amount of available logical cores
	workerAmount := runtime.NumCPU()
	// Define a WaitGroup and a Mutex to synchronize the concurrent execution of the checks
	// The WaitGroup is used to wait for all checks to complete before returning the results
	// The Mutex is used to synchronize access to the checkResults slice and the progress dialog
	var wg sync.WaitGroup
	var mu sync.Mutex
	counter := 1

	var checkResults []checks.Check

	// Iterate over all checks and add them to the channel
	for _, checkSlice := range ChecksList {
		for _, check := range checkSlice {
			checksChan <- check
		}
	}

	// Start the workers to execute the checks concurrently
	// Each worker can access the channel and take work from it while it is not closed / there are still checks to execute
	for range workerAmount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for check := range checksChan {
				result := check()
				// The following code should only be modified by one worker at a time
				mu.Lock()
				checkResults = append(checkResults, result)
				// Display the currently running check in the progress dialog
				err := dialog.Text(fmt.Sprintf("Running check %d of %d", counter, totalChecks))
				if err != nil {
					logger.Log.ErrorWithErr("Error setting progress text:", err)
					return
				}
				// Update the progress bar within the progress dialog
				progress := float64(counter) / float64(totalChecks) * 100
				counter++
				err = dialog.Value(int(progress))
				if err != nil {
					logger.Log.ErrorWithErr("Error setting progress value:", err)
					return
				}
				mu.Unlock()
			}
		}()
	}

	close(checksChan)

	// Wait until all workers have finished
	wg.Wait()

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
