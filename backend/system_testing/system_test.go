package system_test

import (
	"bytes"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/getlantern/systray"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logger.SetupTests()
	// Initialize systray
	go localization.Init("../../")
	time.Sleep(100 * time.Millisecond)

	go systray.Run(tray.OnReady, tray.OnQuit)

	// Wait for the system tray application to initialize
	time.Sleep(100 * time.Millisecond)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTrayFunctionality(t *testing.T) {
	// Change the language to Spanish
	tray.ChangeLanguage("Español")

	// Check that the language was correctly changed to Spanish
	require.Equal(t, 3, tray.Language)

	wg := sync.WaitGroup{}
	wg.Add(2)
	errSlice := make([]error, 2)
	// Run the function without dialog
	go func() {
		defer wg.Done()
		_, err := tray.ScanNow(false)
		errSlice[0] = err
	}()

	// Run the function with dialog
	go func() {
		defer wg.Done()
		_, err := tray.ScanNow(true)
		errSlice[1] = err
	}()

	wg.Wait()
	for _, err := range errSlice {
		require.NoError(t, err)
	}

	// Test ChangeScanInterval function
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
		logger.Log.LogLevelSpecific = -1

		wg = sync.WaitGroup{}
		wg.Add(1)
		// Run the function with mocked user input
		go func() {
			defer wg.Done()
			tray.ChangeScanInterval(tc.input)
		}()

		// Wait for the function to complete
		wg.Wait()

		capturedOutput := buf.String()

		// Assert that the printed message matches the expected message
		require.Contains(t, capturedOutput, tc.expectedMessage)
	}
	// Reset log output to standard output
	logger.Log.SetOutput(os.Stdout)
	logger.Log.LogLevelSpecific = 0

	// Test OpenReportingPage function
	tray.ReportingPageOpen = true
	err := tray.OpenReportingPage()
	require.Error(t, err)

	// Test Popup function
	scanResult := []checks.Check{
		{
			IssueID:  13,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	err = tray.Popup(scanResult, "../../reporting-page/database.db")
	require.NoError(t, err)
}
