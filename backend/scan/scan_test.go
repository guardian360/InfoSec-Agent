package scan_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"os"
	"testing"

	"github.com/ncruces/zenity"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
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
	go localization.Init("../../")

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
	// Display a progress dialog while the scan is running
	dialog, err := zenity.Progress(
		zenity.Title("Security/Privacy Scan"))
	if err != nil {
		logger.Log.ErrorWithErr("Error creating dialog", err)
	}
	// Defer closing the dialog until the scan completes
	defer func(dialog zenity.ProgressDialog) {
		err = dialog.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing dialog", err)
		}
	}(dialog)

	// Execute the scan
	_, err = scan.Scan(dialog, 1)
	require.NoError(t, err)

	// Execute the scan without a dialog
	_, err = scan.Scan(nil, 1)
	require.NoError(t, err)
}

// TestDirectoryExists tests the DirectoryExists function to ensure it correctly identifies whether a directory exists.
//
// This test function calls the DirectoryExists function with a path to an existing directory and asserts that the function returns true.
// It also tests the function with a path to a non-existing directory and asserts that the function returns false.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestDirectoryExists(t *testing.T) {
	// Test for existing directory
	exists := scan.DirectoryExists("../../reporting-page")
	require.True(t, exists)

	// Test for non-existing directory
	exists = scan.DirectoryExists("non-existing-directory")
	require.False(t, exists)
}

// TestGeneratePath tests the GeneratePath function to ensure it correctly generates the path.
//
// This test function calls the GeneratePath function with a given path and asserts that the returned path matches the expected value.
// It also tests the function with an empty string and asserts that the returned path is the current user's home directory.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGeneratePath(t *testing.T) {
	// Test for valid generated path
	path := scan.GeneratePath("\\test")
	currHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Test failed: error getting user home directory: %v", err)
	}
	require.Equal(t, currHomeDir+"\\test", path)

	// Test that given no path, it returns the path to the current user's home directory
	path = scan.GeneratePath("")
	require.Equal(t, currHomeDir, path)
}

func TestCheckInstalled(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		path string
		want bool
	}{
		{
			name: "Path that exists",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\msedge.exe"}},
			},
			path: "msedge.exe",
			want: true,
		},
		{
			name: "Path that does not exist",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\msedge.exe"}},
			},
			path: "chrome.exe",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scan.CheckInstalled(tt.key, tt.path)
			if got != tt.want {
				t.Errorf("CheckInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}
