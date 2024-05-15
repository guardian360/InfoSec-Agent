package scan_test

import (
	"database/sql"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
	"github.com/stretchr/testify/require"
)

// TestGetSeverity tests the GetSeverity function to ensure it returns the correct severity level for a given check ID.
//
// This test function creates a new SQLite database connection and calls the GetSeverity function with known check IDs.
// It then asserts that the returned severity level matches the expected value for each check ID.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetSeverity(t *testing.T) {
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
}

// TestGetJSONKey tests the GetJSONKey function to ensure it returns the correct JSON key for a given check ID.
//
// This test function creates a new SQLite database connection and calls the GetJSONKey function with known check IDs.
// It then asserts that the returned JSON key matches the expected value for each check ID.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestGetJSONKey(t *testing.T) {
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
}
