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
	db, err := sql.Open("sqlite", "../../reporting-page/database.db")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	// Test case where the severity level is 0
	severity, err := scan.GetSeverity(db, 1, 1)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	require.Equal(t, 4, severity)

	// Test case where the severity level is 1
	severity, err = scan.GetSeverity(db, 3, 0)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	require.Equal(t, 1, severity)
}
