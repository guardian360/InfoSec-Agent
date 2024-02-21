// Package tray_test is the testing package for the tray.go file, responsible for unit-testing the basic system tray functionality
//
// Function(s): TestChangeScanInterval, TestScanNow, TestOnQuit

package tray_test

import (
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
		w.Close()
		os.Stdout = oldStdout
		capturedOutput, _ := io.ReadAll(r)

		// Assert that the printed message matches the expected message
		if string(capturedOutput) != tc.expectedMessage {
			t.Errorf("Unexpected message: got %q, want %q", string(capturedOutput), tc.expectedMessage)
		}
	}
}

// Test the ScanNow function
func TestScanNow(t *testing.T) {
	// Set up initial scanCounter value
	initialScanCounter := 0

	tickerAdvanced := make(chan struct{})

	// Listen for ticker advancement
	go func() {
		<-tray.GetScanTicker().C
		tickerAdvanced <- struct{}{}
	}()

	// Run the function
	tray.ScanNow()

	// Assert that scanCounter was incremented
	finalScanCounter := tray.GetScanCounter()
	expectedScanCounter := initialScanCounter + 1
	if finalScanCounter != expectedScanCounter {
		t.Errorf("Scan counter mismatch: got %d, want %d", finalScanCounter, expectedScanCounter)
	}
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
