package scan_test

import (
	"database/sql"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
	"github.com/stretchr/testify/require"
)

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
	logger.SetupTests()

	testCases := []struct {
		issueID          int
		resultID         int
		expectedSeverity int
	}{
		{1, 1, 4},
		{3, 0, 1},
	}
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	for _, tc := range testCases {
		severity, err := scan.GetSeverity(db, tc.issueID, tc.resultID)
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		require.Equal(t, tc.expectedSeverity, severity)
	}

	// Test for invalid issue ID and result ID pair
	_, err = scan.GetSeverity(db, 0, 0)
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
	logger.SetupTests()

	testCases := []struct {
		issueID         int
		resultID        int
		expectedJSONKey int
	}{
		{1, 1, 11},
		{3, 0, 30},
	}

	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	for _, tc := range testCases {
		jsonKey, err := scan.GetJSONKey(db, tc.issueID, tc.resultID)
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		require.Equal(t, tc.expectedJSONKey, jsonKey)
	}

	// Test for invalid issue ID and result ID pair
	_, err = scan.GetSeverity(db, 0, 0)
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
	logger.SetupTests()

	scanResult := []checks.Check{
		{
			IssueID:  1,
			ResultID: 1,
			Result:   []string{"Issue 1"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	expectedData := []scan.DataBaseData{
		{
			CheckID:  1,
			Severity: 4,
			JSONKey:  11,
		},
	}
	emptyScanResult := []checks.Check{}
	emptyExpectedData := []scan.DataBaseData{}
	invalidScanResult := []checks.Check{
		{
			IssueID:  0,
			ResultID: 0,
			Result:   []string{"Issue 0"},
			Error:    nil,
			ErrorMSG: "",
		},
	}
	invalidExpectedData := []scan.DataBaseData{
		{
			CheckID:  0,
			Severity: 0,
			JSONKey:  0,
		},
	}
	wrongPathExpectedData := []scan.DataBaseData{
		{
			CheckID:  1,
			Severity: 0,
			JSONKey:  0,
		},
	}
	testCases := []struct {
		scanResult   []checks.Check
		expectedData []scan.DataBaseData
	}{
		{scanResult, expectedData},
		{emptyScanResult, emptyExpectedData},
	}

	for _, tc := range testCases {
		data, err := scan.GetDataBaseData(tc.scanResult, "../../reporting-page/database.db")
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}
		require.Equal(t, tc.expectedData, data)
		require.Equal(t, tc.expectedData, data)
	}

	// Test for invalid scan result
	result, _ := scan.GetDataBaseData(invalidScanResult, "../../reporting-page/database.db")
	require.Equal(t, invalidExpectedData, result)

	// Test for invalid database path
	result, _ = scan.GetDataBaseData(scanResult, "")
	require.Equal(t, wrongPathExpectedData, result)
}
