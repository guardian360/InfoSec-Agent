package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/scan"
)

// DataBase represents the application's database. It is used for invoking functions that interact with the database.
//
// The DataBase struct encapsulates the database operations and provides methods for performing various database operations such as retrieving severity levels for each check performed. It does not contain any fields as it serves as a receiver for methods that interact with the database.
type DataBase struct {
}

// NewDataBase constructs and returns a new instance of the DataBase struct.
//
// This function is used to create a new DataBase object, which represents the database and is used for calling functions that interact with the database.
//
// Parameters: None.
//
// Returns: *DataBase: A pointer to a newly created DataBase object.
func NewDataBase() *DataBase {
	return &DataBase{}
}

// GetAllSeverities retrieves the severity levels for each check performed, based on the provided checks and result IDs.
//
// This function iterates over the provided checks and result IDs, and retrieves the severity level for each check. The severity level is determined by the check's result ID. The function returns a list of severity levels, where each severity level corresponds to a check in the order they were provided.
//
// Parameters:
//   - checks []checks.Check: A list of checks from which to retrieve severity levels.
//   - resultIDs []int: A list of result IDs, each corresponding to a severity level for a check.
//
// Returns:
//   - []scan.Severity: A list of severity levels for each check, in the order the checks were provided.
//   - error: An error object that describes the error, if any occurred. nil if no error occurred.
func (d *DataBase) GetAllSeverities(checks []checks.Check, resultIDs []int) ([]scan.Severity, error) {
	return scan.GetAllSeverities(checks, resultIDs)
}
