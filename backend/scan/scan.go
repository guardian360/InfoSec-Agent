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
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"

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
func Scan(dialog zenity.ProgressDialog, language int) ([]checks.Check, error) {
	dialogPresent := false
	if dialog != nil {
		dialogPresent = true
	}

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
	startWorkers(workerAmount, &wg, checksChan, &mu, &checkResults, dialogPresent, &counter, totalChecks, dialog, language)

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

// startWorkers creates and starts the specified number of workers to execute the checks concurrently.
// The workers receive checks from the checksChan channel, execute them, and store the results in the checkResults slice.
//
// Parameters:
//   - workerAmount (int): The number of workers to create and start.
//   - wg (*sync.WaitGroup): A WaitGroup object that allows the workers to signal when they have completed their work.
//   - checksChan (chan func() checks.Check): A channel that provides the workers with checks to execute.
//   - mu (*sync.Mutex): A Mutex object that synchronizes access to the checkResults slice and the progress dialog.
//   - checkResults (*[]checks.Check): A pointer to a slice of Check objects that stores the results of the executed checks.
//   - dialogPresent (bool): A boolean value that indicates whether a progress dialog is present.
//   - counter (*int): A pointer to an integer that represents the current check number.
//   - totalChecks (int): An integer value that represents the total number of checks to be executed.
//   - dialog (zenity.ProgressDialog): A dialog window that displays the progress of the scan as it runs.
//
// Returns: None.
func startWorkers(workerAmount int, wg *sync.WaitGroup, checksChan chan func() checks.Check, mu *sync.Mutex,
	checkResults *[]checks.Check, dialogPresent bool, counter *int, totalChecks int, dialog zenity.ProgressDialog, language int) {
	for range workerAmount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for check := range checksChan {
				result := check()
				mu.Lock()
				*checkResults = append(*checkResults, result)
				if dialogPresent {
					err := dialog.Text(fmt.Sprintf(localization.Localize(language, "Dialogs.Scan.Content"), *counter, totalChecks))
					if err != nil {
						logger.Log.ErrorWithErr("Error setting progress text:", err)
						mu.Unlock()
						return
					}
					progress := float64(*counter) / float64(totalChecks) * 100
					*counter++
					err = dialog.Value(int(progress))
					if err != nil {
						logger.Log.ErrorWithErr("Error setting progress value:", err)
						mu.Unlock()
						return
					}
				}
				mu.Unlock()
			}
		}()
	}
}
