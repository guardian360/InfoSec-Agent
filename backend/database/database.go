// Package database provides functions for interacting with the database.
//
// Exported function(s): FillDatabase, GetSeverity, GetJSONKey, GetData
package database

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

// JSONData represents the structure of the JSON data
type JSONData map[string]Issue

// Issue represents the structure of an issue in the JSON data
type Issue struct {
	Type        string            `json:"Type"`
	Information string            `json:"Information"`
	Results     map[string]Result `json:"Results"`
}

// Result represents the structure of a result in the JSON data
type Result struct {
	Severity    int      `json:"Severity"`
	Name        string   `json:"Name"`
	Solution    []string `json:"Solution"`
	Screenshots []string `json:"Screenshots"`
}

// Data represents a simplified structure with issueID and severity
type Data struct {
	IssueID  int
	Severity int
}

// GetSeverity is a function that retrieves the severity level of a specific issue from the JSON data.
// It takes three parameters:
// - data: a JSONData object that represents the structure of the JSON data.
// - issueID: an integer that represents the ID of the issue.
// - resultID: an integer that represents the ID of the result.
//
// The function first converts the issueID and resultID to strings. It then retrieves the issue from the JSON data
// using the issueID. If the issue does not exist, it returns 0 and an error indicating that the issue was not found.
//
// If the issue exists, the function then retrieves the result from the issue using the resultID. If the result does
// not exist, it returns 0 and an error indicating that the result was not found.
//
// If both the issue and the result exist, the function returns the severity level of the result and nil for the error.
func GetSeverity(data JSONData, issueID int, resultID int) (int, error) {
	// Convert issueID and resultID to string
	issueKey := strconv.Itoa(issueID)
	resultKey := strconv.Itoa(resultID)

	// Retrieve the issue from the JSON data
	issue, exists := data[issueKey]
	if !exists {
		return 0, errors.New("issue not found")
	}

	// Retrieve the result from the issue
	result, exists := issue.Results[resultKey]
	if !exists {
		return 0, errors.New("result not found")
	}

	// Return the severity level
	return result.Severity, nil
}

// GetData is a function that computes the severity for all found checks and puts them into a list of Data.
// It takes two parameters:
// - jsonFilePath: a string that represents the path to the JSON file.
// - checkResults: a slice of checks.Check that represents the check results.
//
// The function first opens the JSON file and decodes the JSON data into a JSONData object.
// If there's an error opening the file or decoding the JSON data, it returns nil and the error.
//
// The function then initializes a slice of Data and iterates through all check results to compute the severities.
// For each check result, it retrieves the issueID and resultID, and calls the GetSeverity function to compute the severity.
// If there's an error retrieving the severity, it logs the error and skips to the next check result.
//
// If the severity is successfully computed, it creates a new Data object with the issueID and severity, and appends it to the dataList.
//
// After iterating through all check results, the function returns the dataList and nil for the error.
func GetData(jsonFilePath string, checkResults []checks.Check) ([]Data, error) {
	// Open the JSON file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// interface
	// Decode the JSON data
	var jsonData JSONData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		return nil, err
	}

	// Initialize the list of Data
	var dataList []Data

	// Iterate through all check results to compute the severities
	for _, checkResult := range checkResults {
		issueID := checkResult.IssueID
		resultID := checkResult.ResultID

		// Compute the severity
		severity, sevErr := GetSeverity(jsonData, issueID, resultID)
		if sevErr != nil {
			logger.Log.ErrorWithErr("Error getting severity:", sevErr)
			continue // Skip if there's an error retrieving severity
		}
		// Add to the dataList
		dataList = append(dataList, Data{IssueID: issueID, Severity: severity})
	}

	return dataList, nil
}
