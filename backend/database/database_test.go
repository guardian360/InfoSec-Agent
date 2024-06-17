package database_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
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

func TestGetData(t *testing.T) {
	// Create a slice of checks.Check with IssueID and ResultID that exist in your JSON data
	scanResults := []checks.Check{
		{
			IssueID:  29,
			ResultID: 1, // severity 2
		},
		{
			IssueID:  5,
			ResultID: 1, // severity 1
		},
	}

	// Provide a valid JSON file path
	validJSONFilePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"

	// Call PointCalculation method
	got, _ := database.GetData(validJSONFilePath, scanResults)

	// Create a slice of Data with the expected values
	expected := []database.Data{
		{IssueID: 29, Severity: 2},
		{IssueID: 5, Severity: 1},
	}

	// Assert that the points are calculated correctly
	require.Equal(t, expected, got)

}

func TestGetData_UnmarshalError(t *testing.T) {
	scanResults := []checks.Check{}
	invalidJSONFilePath := "../../invalid.json"

	_, err := database.GetData(invalidJSONFilePath, scanResults)

	require.Error(t, err)
}

func TestGetData_ResultError(t *testing.T) {
	// Create a slice of checks.Check with an error
	scanResults := []checks.Check{
		{
			IssueID:  29,
			ResultID: -1, // severity 2
			Error:    errors.New("mock error"),
		},
	}

	// Provide a valid JSON file path
	validJSONFilePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"

	// Call PointCalculation method
	got, _ := database.GetData(validJSONFilePath, scanResults)

	var dataList []database.Data
	// Assert that the points remain the same as the initial GameState
	require.Equal(t, dataList, got)
}

func TestGetData_IssueIDNotFound(t *testing.T) {
	// Create a slice of checks.Check with an IssueID that does not exist in your JSON data
	scanResults := []checks.Check{
		{
			IssueID:  9999, // This IssueID does not exist in the JSON data
			ResultID: 1,
		},
		{
			IssueID:  5,
			ResultID: 22,
		},
	}

	// Provide a valid JSON file path
	validJSONFilePath := "../../reporting-page/frontend/src/databases/database.en-GB.json"

	// Call PointCalculation method
	got, _ := database.GetData(validJSONFilePath, scanResults)
	var dataList []database.Data
	// Assert that the points remain the same as the initial GameState
	require.Equal(t, dataList, got)
}

func TestGetData_SeverityNotFound(t *testing.T) {
	// Create a slice of checks.Check with an IssueID and ResultID that do not exist in your JSON data
	scanResults := []checks.Check{
		{
			IssueID:  9999, // This IssueID does not exist in the JSON data
			ResultID: 9999, // This ResultID does not exist in the JSON data
		},
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example.*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some JSON data to the file
	text := `{"9999": {"9999": {}}}`
	if _, fileErr := tmpfile.WriteString(text); fileErr != nil {
		tmpfile.Close()
		t.Fatal(fileErr)
	}
	if closeErr := tmpfile.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	// Call PointCalculation method
	got, _ := database.GetData(tmpfile.Name(), scanResults)

	var dataList []database.Data
	// Assert that the points remain the same as the initial GameState
	require.Equal(t, dataList, got)
}
