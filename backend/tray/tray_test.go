package tray_test

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"

	"github.com/getlantern/systray"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
)

// TestMain sets up the necessary environment for the system tray application tests and executes them.
//
// This function initializes the localization and system tray, waits for the system tray application to be ready, and then runs the tests.
// After the tests are completed, it ensures the system tray is properly cleaned up and the application exits with the appropriate exit code.
//
// Parameters:
//   - m *testing.M: The testing framework that manages and runs the tests.
//
// Returns: None. The function calls os.Exit with the exit code returned by m.Run().
func TestMain(m *testing.M) {
	logger.SetupTests()

	// Initialize systray
	go localization.Init("../../")
	time.Sleep(100 * time.Millisecond)

	go systray.Run(tray.OnReady, tray.OnQuit)

	// Wait for the system tray application to initialize
	time.Sleep(100 * time.Millisecond)

	// Run tests
	exitCode := m.Run()

	// Clean up systray
	systray.Quit()

	os.Exit(exitCode)
}

// TestChangeScanInterval validates the ChangeScanInterval function by testing it with both valid and invalid inputs.
//
// This test function defines a set of test cases with various input values, including valid scan intervals, non-numeric values, negative values, zero, and a large valid interval.
// For each test case, it simulates user input for the ChangeScanInterval function and captures the output message.
// The function asserts that the output message matches the expected result for each input, validating that the ChangeScanInterval function correctly changes the scan interval for valid inputs and handles invalid inputs appropriately.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestChangeScanInterval(t *testing.T) {
	logger.Log.LogLevelSpecific = -1
	// Define test cases with input values and expected results
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		// Valid input
		{"24", "Changing scan interval to 24 day(s)"},
		// Invalid input (non-numeric)
		{"abc", "Invalid scan interval input"},
		// Invalid input (negative)
		{"-1", "Invalid scan interval input"},
		// Invalid input (zero)
		{"0", "Invalid scan interval input"},
		// Valid large input
		{"1000", "Changing scan interval to 1000 day(s)"},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		var buf bytes.Buffer
		logger.Log.SetOutput(&buf)

		// Run the function with mocked user input
		tray.ChangeScanInterval(tc.input)

		capturedOutput := buf.String()

		// Assert that the printed message matches the expected message
		require.Contains(t, capturedOutput, tc.expectedMessage)
	}
	// Reset logger
	logger.Log.SetOutput(os.Stdout)
	logger.Log.LogLevelSpecific = 0
}

// TestScanNow validates the behavior of the ScanNow function.
//
// This test function initiates a scan by calling the ScanNow function and verifies that the ScanCounter increments correctly, indicating that a scan has been performed.
// It also checks that the ScanNow function does not return any errors during its execution.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestScanNow(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	errSlice := make([]error, 2)
	// Run the function without dialog
	go func() {
		defer wg.Done()
		_, err := tray.ScanNow(false, config.DatabasePath)
		errSlice[0] = err
	}()

	wg.Wait()
	wg.Add(1)
	// Run the function with dialog
	go func() {
		defer wg.Done()
		_, err := tray.ScanNow(true, config.DatabasePath)
		errSlice[1] = err
	}()

	wg.Wait()
	for _, err := range errSlice {
		require.NoError(t, err)
	}
}

// TestOnQuit validates the behavior of the OnQuit function by simulating an application quit scenario.
//
// This test function simulates an application quit by sending an os.Interrupt signal and then verifies that the OnQuit function completes its execution within a reasonable time frame.
// This ensures that the OnQuit function responds appropriately to quit signals and completes its cleanup tasks in a timely manner.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestOnQuit(t *testing.T) {
	// Mock OS signals channel
	sigc := make(chan os.Signal, 1)

	// Create a channel to signal completion
	quitCompleted := make(chan struct{})

	// Run OnQuit in a separate goroutine
	go func() {
		tray.OnQuit()
		quitCompleted <- struct{}{}
	}()

	// Simulate quitting the application
	sigc <- os.Interrupt

	// Wait for OnQuit to complete
	select {
	case <-quitCompleted:
		// OnQuit completed
	case <-time.After(1 * time.Second):
		t.Error("OnQuit did not complete within the timeout")
	}
}

// TestTranslation verifies the functionality of the localization package by ensuring that it correctly translates strings based on the set language.
//
// This test function retrieves a localized string for a given message ID using the initial language setting.
// It then changes the language setting and retrieves the localized string for the same message ID.
// The test asserts that the two localized strings are different, validating that the localization package correctly translates strings based on the set language.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestTranslation(t *testing.T) {
	var localizer = localization.Localizers()[0]
	s1 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "Tray.ScanIntervalTitle",
	})
	// Change the language, then check if the translation is different
	localizer = localization.Localizers()[1]
	s2 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "Tray.ScanIntervalTitle",
	})
	require.NotEqual(t, s1, s2)
}

// TestChangeLang validates the behavior of the ChangeLanguage function with both valid and invalid inputs.
//
// This test function iterates over a set of test cases that include valid language inputs and an invalid language input.
// For each test case, it calls the ChangeLanguage function with the test input and asserts that the language index returned by the Language function matches the expected index.
// This validates that the ChangeLanguage function correctly updates the language index for valid inputs and defaults to the index for "English (UK)" for invalid inputs.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestChangeLang(t *testing.T) {
	testCases := []struct {
		input         string
		expectedIndex int
	}{
		// Valid input
		{"Deutsch", 0},
		{"English (UK)", 1},
		{"English (US)", 2},
		{"Español", 3},
		{"Français", 4},
		{"Nederlands", 5},
		{"Português", 6},
		// Invalid input, should return the default index ( English (UK))
		{"Italian", 1},
	}

	for _, tc := range testCases {
		tray.ChangeLanguage(tc.input)
		require.Equal(t, tc.expectedIndex, tray.Language)
	}
}

// TestRefreshMenu validates the behavior of the RefreshMenu function by ensuring it correctly updates the menu items with the current language.
//
// This test function initially captures the translation of a menu item, then changes the application language and refreshes the menu.
// It asserts that the translation of the same menu item is different after the language change and menu refresh, validating that the RefreshMenu function correctly updates the menu items to reflect the current language.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestRefreshMenu(t *testing.T) {
	value1 := tray.MenuItems[0].MenuTitle
	translation1 := localization.Localize(tray.Language, value1)
	tray.ChangeLanguage("Español")
	// Refresh the menu, then check if the translation is different
	// RefreshMenu(MenuItems)
	value2 := tray.MenuItems[0].MenuTitle
	translation2 := localization.Localize(tray.Language, value2)

	require.NotEqual(t, translation1, translation2)
}

// TestOpenReportingPageWhenAlreadyOpen verifies the behavior of the OpenReportingPage function when a reporting page is already open.
//
// This test function sets the ReportingPageOpen flag to true, simulating a scenario where a reporting page is already running.
// It then calls the OpenReportingPage function and asserts that an error is returned, indicating that a new reporting page cannot be opened while one is already running.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestOpenReportingPageWhenAlreadyOpen(t *testing.T) {
	tray.ReportingPageOpen = true
	err := tray.OpenReportingPage()
	require.Error(t, err)
	require.Equal(t, "reporting-page is already running", err.Error())
}

// TestPopup verifies the behavior of the Popup function.
//
// This test function sets up a scan result and calls the Popup function to display a popup notification with the scan result.
// It then asserts that the Popup function does not return any errors during its execution.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestPopup(t *testing.T) {
	tray.ChangeLanguage("English (UK)")
	// Define check result
	scanResult := []checks.Check{
		{
			IssueID:  13,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}

	// Run the function
	er := tray.Popup(scanResult, "../../reporting-page/database.db")
	require.NoError(t, er)
}

// TestPopupMessage varifies the behavior of the PopupMessage function by entering scan results and verifying that it returns a correct message.
//
// This test function sets up a scan result and calls the PopupMessage function to generate a message based on the scan result.
// It then asserts that the generated message matches the expected message based on the scan result, validating that the PopupMessage function correctly formats messages based on scan results.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestPopupMessage(t *testing.T) {
	// Define test cases with input values and expected results
	scanResult1 := []checks.Check{
		{
			IssueID:  13,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
		{
			IssueID:  13,
			ResultID: 0,
			Result:   []string{"Issue 2"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	scanResult2 := []checks.Check{
		{
			IssueID:  12,
			ResultID: 0,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
		{
			IssueID:  13,
			ResultID: 1,
			Result:   []string{"Issue 2"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	scanResult3 := []checks.Check{
		{
			IssueID:  3,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	scanResult4 := []checks.Check{
		{
			IssueID:  13,
			ResultID: 0,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
		{
			IssueID:  3,
			ResultID: 1,
			Result:   []string{"Issue 2"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	scanResult5 := []checks.Check{
		{
			IssueID:  1,
			ResultID: 0,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	testCases := []struct {
		input           []checks.Check
		expectedMessage string
	}{
		// Scanresult with 1 high risk issue and medium risk issues
		{scanResult1, "The privacy and security scan has been completed. You have 1 high risk issue. Open the reporting page to see more information."},
		// Scanresult with multiple high risk issues
		{scanResult2, "The privacy and security scan has been completed. You have 2 high risk issues. Open the reporting page to see more information."},
		// Scan result with 1 medium risk issue
		{scanResult3, "The privacy and security scan has been completed. You have 1 medium risk issue. Open the reporting page to see more information."},
		// Scanresult with no high risk issues an multiple medium risk issues
		{scanResult4, "The privacy and security scan has been completed. You have 2 medium risk issues. Open the reporting page to see more information."},
		// Scanresult with no high risk issues an no medium risk issues
		{scanResult5, "The privacy and security scan has been completed. Open the reporting page to view the results."},
		// Empty scanresult
		{[]checks.Check{}, "The privacy and security scan has been completed. Open the reporting page to view the results."},
	}
	// Iterate over test cases
	for _, tc := range testCases {
		// Run the function with mocked scan
		result := tray.PopupMessage(tc.input, "../../"+config.DatabasePath)

		// Assert that the message matches the expected message
		require.Equal(t, tc.expectedMessage, result)
	}
}
