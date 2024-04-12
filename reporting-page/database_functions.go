package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/scan"
)

// Database represents the database, used for calling function which connect to the database
type DataBase struct {
}

// NewDataBase creates a new DataBase struct
//
// Parameters: _
//
// Returns: a pointer to a new DataBase struct
func NewDataBase() *DataBase {
	return &DataBase{}
}

// TODO: Fix this one
// GetDataBaseData gets all severities and JSON keys found by the check
//
// Parameters: checks ([]checks.Check) - list of checks to get severities from
//
// resultIDs ([]int) - list of result ids corresponding to a severity level
//
// Returns: list of severity levels and JSON keys for each issue checked
func (d *DataBase) GetDataBaseData(checks []checks.Check) ([]scan.DataBaseData, error) {
	return scan.GetDataBaseData(checks)
}
