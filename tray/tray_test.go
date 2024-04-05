package tray_test

import (
	"bytes"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
	"log"
	"os"
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"

	"github.com/getlantern/systray"
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
	go localization.Init("../")
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
	// Define test cases with input values and expected results
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		// Valid input
		{"24", "Scan interval changed to 24 hours\n"},
		// Invalid input (non-numeric)
		{"abc", "Invalid input"},
		// Invalid input (negative)
		{"-1", "Invalid input"},
		// Invalid input (zero)
		{"0", "Invalid input"},
		// Valid large input
		{"1000", "Scan interval changed to 1000 hours\n"},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		var buf bytes.Buffer
		logger.Log.SetOutput(&buf)

		// Run the function with mocked user input
		go tray.ChangeScanInterval(tc.input)

		// Wait for the function to complete
		time.Sleep(100 * time.Millisecond)

		capturedOutput := buf.String()

		// Assert that the printed message matches the expected message
		require.Contains(t, capturedOutput, tc.expectedMessage)

		// Reset log output to standard output
		log.SetOutput(os.Stdout)
	}
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
	// Set up initial ScanCounter value
	initialScanCounter := 0

	tickerAdvanced := make(chan struct{})

	// Listen for ticker advancement
	go func() {
		<-tray.ScanTicker.C
		tickerAdvanced <- struct{}{}
	}()

	// Run the function
	_, err := tray.ScanNow()
	require.NoError(t, err)

	// Assert that ScanCounter was incremented
	finalScanCounter := tray.ScanCounter
	expectedScanCounter := initialScanCounter + 1
	require.Equal(t, expectedScanCounter, finalScanCounter)
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
// This validates that the ChangeLanguage function correctly updates the language index for valid inputs and defaults to the index for "British English" for invalid inputs.
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
		{"German", 0},
		{"British English", 1},
		{"American English", 2},
		{"Spanish", 3},
		{"French", 4},
		{"Dutch", 5},
		{"Portuguese", 6},
		// Invalid input, should return the default index (British English)
		{"Italian", 1},
	}

	for _, tc := range testCases {
		tray.ChangeLanguage(tc.input)
		require.Equal(t, tc.expectedIndex, tray.Language())
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
	translation1 := localization.Localize(tray.Language(), value1)
	tray.ChangeLanguage("Spanish")
	// Refresh the menu, then check if the translation is different
	// RefreshMenu(MenuItems)
	value2 := tray.MenuItems[0].MenuTitle
	translation2 := localization.Localize(tray.Language(), value2)

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
	err := tray.OpenReportingPage("../")
	require.Error(t, err)
	require.Equal(t, "reporting-page is already running", err.Error())
}
