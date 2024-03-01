// Package tray_test is the testing package for the tray.go file, responsible for unit-testing the basic system tray functionality
//
// Function(s): TestChangeScanInterval, TestScanNow, TestOnQuit

package tray_test

import (
	"InfoSec-Agent/localization"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"

	"InfoSec-Agent/tray"

	"github.com/getlantern/systray"
)

// Setup for the tests
func TestMain(m *testing.M) {
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

// Test the ChangeScanInterval function
func TestChangeScanInterval(t *testing.T) {
	// Define test cases with input values and expected results
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		// Valid input
		{"24", "Scan interval changed to 24 hours\n"},
		// Invalid input (non-numeric)
		{"abc", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Invalid input (negative)
		{"-1", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Invalid input (zero)
		{"0", "Invalid input. Using default interval of 24 hours.Scan interval changed to 24 hours\n"},
		// Valid large input
		{"1000", "Scan interval changed to 1000 hours\n"},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Capture standard output to check the printed message
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the function with mocked user input
		go tray.ChangeScanInterval(tc.input)

		// Wait for the function to complete
		time.Sleep(100 * time.Millisecond)

		// Restore standard output
		err := w.Close()
		if err != nil {
			t.Errorf("Error closing pipe: %v", err)
		}
		os.Stdout = oldStdout
		capturedOutput, _ := io.ReadAll(r)

		// Assert that the printed message matches the expected message
		require.Equal(t, string(capturedOutput), tc.expectedMessage)
	}
}

// Test the ScanNow function
func TestScanNow(t *testing.T) {
	// Set up initial scanCounter value
	initialScanCounter := 0

	tickerAdvanced := make(chan struct{})

	// Listen for ticker advancement
	go func() {
		<-tray.ScanTicker().C
		tickerAdvanced <- struct{}{}
	}()

	// Run the function
	tray.ScanNow()

	// Assert that scanCounter was incremented
	finalScanCounter := tray.ScanCounter()
	expectedScanCounter := initialScanCounter + 1
	require.Equal(t, finalScanCounter, expectedScanCounter)
}

// Test the OnQuit function
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

// TestTranslation tests the localization package, ensuring that strings are translated correctly
func TestTranslation(t *testing.T) {
	var localizer = localization.Localizers()[0]
	s1 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "ScanIntervalTitle",
	})
	// Change the language, then check if the translation is different
	localizer = localization.Localizers()[1]
	s2 := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "ScanIntervalTitle",
	})
	require.NotEqual(t, s1, s2)
}

// TestChangeLang tests the tray.ChangeLang function on valid and invalid inputs
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

func TestRefreshMenu(t *testing.T) {
	value1 := tray.MenuItems()[0].MenuTitle
	translation1 := localization.Localize(tray.Language(), value1)
	tray.ChangeLanguage("Spanish")
	// Refresh the menu, then check if the translation is different
	tray.RefreshMenu(tray.MenuItems())
	value2 := tray.MenuItems()[0].MenuTitle
	translation2 := localization.Localize(tray.Language(), value2)

	require.NotEqual(t, translation1, translation2)
}
