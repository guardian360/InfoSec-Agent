package scan_test

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"os"
	"testing"

	"github.com/ncruces/zenity"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
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
	go localization.Init("../")

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
	_, err = scan.Scan(dialog, 1)
	require.NoError(t, err)

	// Execute the scan without a dialog
	_, err = scan.Scan(nil, 1)
	require.NoError(t, err)
}

// TestGetSeverity tests the GetSeverity function to ensure it returns the correct severity level for a given issue ID and result ID pair.
//
// This test function creates a new SQLite database connection and calls the GetSeverity function with known issue IDs and result IDs.
// It then asserts that the returned severity level matches the expected value for each issue ID and result ID pair.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetSeverity(t *testing.T) {
	// Arrange database connection
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	// Test for valid issue ID and result ID
	severity, err := database.GetSeverity(db, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 4, severity)

	// Test for invalid issue ID and result ID
	_, err = database.GetSeverity(db, 0, 0)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows.Error(), err.Error())
}

// TestGetJSONKey tests the GetJSONKey function to ensure it returns the correct JSON key for a given issue ID and result ID pair.
//
// This test function creates a new SQLite database connection and calls the GetJSONKey function with known issue IDs and result IDs.
// It then asserts that the returned JSON key matches the expected value for each issue ID and result ID pair.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetJSONKey(t *testing.T) {
	// Arrange database connection
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	// Test for valid issue ID and result ID
	jsonKey, err := database.GetJSONKey(db, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 11, jsonKey)

	// Test for invalid issue ID and result ID
	_, err = database.GetSeverity(db, 0, 0)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows.Error(), err.Error())
}

// TestGetDataBaseData tests the GetDataBaseData function to ensure it returns the correct database data for a given list of checks.
//
// This test function creates a list of check results and calls the GetDataBaseData function.
// It then asserts that the returned database data matches the expected data for the given checks.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetDataBaseData(t *testing.T) {
	scanResult := []checks.Check{
		{
			IssueID:  1,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	expectedData := []database.Data{
		{
			CheckID:  1,
			Severity: 4,
			JSONKey:  11,
		},
	}
	emptyScanResult := []checks.Check{}
	emptyExpectedData := []database.Data{}
	invalidScanResult := []checks.Check{
		{
			IssueID:  0,
			ResultID: 0,
			Result:   []string{"Issue 0"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	invalidExpectedData := []database.Data{
		{
			CheckID:  0,
			Severity: 0,
			JSONKey:  0,
		},
	}
	wrongPathExpectedData := []database.Data{
		{
			CheckID:  1,
			Severity: 0,
			JSONKey:  0,
		},
	}
	testCases := []struct {
		scanResult   []checks.Check
		expectedData []database.Data
	}{
		{scanResult, expectedData},
		{emptyScanResult, emptyExpectedData},
		{invalidScanResult, invalidExpectedData},
	}

	for _, tc := range testCases {
		data, err := database.GetData(tc.scanResult, "../../reporting-page/database.db")
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		require.Equal(t, tc.expectedData, data)
		require.Equal(t, tc.expectedData, data)
	}

	// Test for invalid database path
	result, _ := database.GetData(scanResult, "")
	require.Equal(t, wrongPathExpectedData, result)
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
