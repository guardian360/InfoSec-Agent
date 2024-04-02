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

// future enumerator maybe for severity
var Severities = []scan.SeverityLevels{}

// func (d *DataBase) CreateIssues() []scan.SeverityLevels {
// 	return []scan.SeverityLevels{{scan.Safe, "0"},
// 		{scan.Low, "1"},
// 		{scan.Medium, "2"},
// 		{scan.High, "3"}}
// }

// GetDataBaseData gets all severities and JSON keys found by the check
//
// Parameters: checks ([]checks.Check) - list of checks to get severities from
//
// resultIDs ([]int) - list of result ids corresponding to a severity level
//
// Returns: list of severity levels and JSON keys for each issue checked
func (d *DataBase) GetDataBaseData(checks []checks.Check, resultIDs []int) ([]scan.DataBaseData, error) {
	return scan.GetDataBaseData(checks, resultIDs)
}
