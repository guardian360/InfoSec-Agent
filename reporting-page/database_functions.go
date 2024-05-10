package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
)

// Database represents the application's database. It is used for invoking functions that interact with the database.
//
// The Database struct encapsulates the database operations and provides methods for performing various database operations such as retrieving severity levels for each check performed. It does not contain any fields as it serves as a receiver for methods that interact with the database.
type Database struct {
}

// NewDatabase constructs and returns a new instance of the Database struct.
//
// This function is used to create a new Database object, which represents the database and is used for calling functions that interact with the database.
//
// Parameters: None.
//
// Returns: *Database: A pointer to a newly created Database object.
func NewDatabase() *Database {
	return &Database{}
}

// GetData gets all severities and JSON keys found by the check
//
// This function iterates over the provided checks and result IDs, and retrieves the severity level for each check. The severity level is determined by the check's result ID. The function returns a list of severity levels, where each severity level corresponds to a check in the order they were provided.
//
// Parameters:
//   - checks []checks.Check: A list of checks from which to retrieve severity levels.
//   - resultIDs []int: A list of result IDs, each corresponding to a severity level for a check.
//
// Returns: list of severity levels and JSON keys for each issue checked
func (d *Database) GetData(checks []checks.Check) ([]database.Data, error) {
	return database.GetData(checks)
}
