package apiconnection

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
)

// TODO: Update documentation
// ParseResult is a struct that encapsulates the results of parsing a scan.
// This parsing is done to convert the results of security and privacy checks into a format that
// can be sent to the Guardian360 Lighthouse API.
type ParseResult struct {
	Metadata Metadata
	Results  []IssueData
}

// TODO: Update documentation
// Metadata is a struct that contains metadata about the scan.
// This metadata includes the workstation ID, the user who initiated the scan, and the date of the scan.
type Metadata struct {
	WorkStationID int
	User          string
	Date          string
}

// TODO: Update documentation
// IssueData is a struct that encapsulates the data for a checked issue.
// This data includes the issue ID, whether the issue is considered a problem, and any additional data related to the
// issue.
type IssueData struct {
	IssueID        int
	Detected       bool
	AdditionalData []string
}

// TODO: Update documentation
// String returns a string representation of the ParseResult struct.
//
// Parameters: None.
//
// Returns:
//   - string: A string representation of the ParseResult struct.
func (p ParseResult) String() string {
	return fmt.Sprintf("Metadata: %v, Results: %v", p.Metadata, p.Results)
}

// TODO: Update documentation
// ParseScanResults parses the results of a scan into a ParseResult struct.
// This function takes the metadata of the scan and the results of the security and privacy checks and creates a
// ParseResult struct that encapsulates this data.
//
// Parameters:
//   - metaData (Metadata): The metadata of the scan, including the workstation ID, user, and date.
//   - checks ([]checks.Check): A slice of Check objects representing the results of the security and privacy checks.
//
// Returns:
//   - ParseResult: A ParseResult struct that encapsulates the metadata and results of the scan.
func ParseScanResults(metaData Metadata, checks []checks.Check) ParseResult {
	var result []IssueData
	for _, check := range checks {
		result = append(result, ParseCheckResult(check))
	}
	parseResult := ParseResult{Metadata: metaData, Results: result}
	return parseResult
}

// TODO: Update documentation
// ParseCheckResult parses the result of a single security or privacy check into an IssueData struct.
//
// Parameters:
// - check (checks.Check): A Check object representing the result of a security or privacy check.
//
// Returns:
// - IssueData: An IssueData struct that encapsulates the data for the checked issue.
func ParseCheckResult(check checks.Check) IssueData {
	if check.Error != nil {
		return IssueData{IssueID: check.IssueID, Detected: false}
	}
	return IssueData{
		IssueID:        check.IssueID,
		Detected:       IssueMap[IssueResPair{check.IssueID, check.ResultID}],
		AdditionalData: check.Result}
}
