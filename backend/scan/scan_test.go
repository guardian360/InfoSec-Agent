package scan_test

import (
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
	"github.com/ncruces/zenity"
	"github.com/stretchr/testify/require"
)

// TestMain sets up the necessary environment for the system scan package tests and executes them.
//
// This function initializes the logger for the tests and runs the tests.
//
// Parameters:
//   - m *testing.M: The testing framework that manages and runs the tests.
//
// Returns: None. The function calls os.Exit with the exit code returned by m.Run().
func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

// TestScan tests the Scan function to ensure it runs without errors.
//
// This test function calls the Scan function and asserts that it does not return an error.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestScan(t *testing.T) {
	// logger.SetupTests()

	// Display a progress dialog while the scan is running
	dialog, err := zenity.Progress(
		zenity.Title("Security/Privacy Scan"))
	if err != nil {
		logger.Log.ErrorWithErr("Error creating dialog:", err)
	}
	// Defer closing the dialog until the scan completes
	defer func(dialog zenity.ProgressDialog) {
		err = dialog.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing dialog:", err)
		}
	}(dialog)

	// Execute the scan
	_, err = scan.Scan(dialog)
	require.NoError(t, err)
}
