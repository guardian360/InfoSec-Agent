package apiconnection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
)

// ParseResult is a struct that encapsulates the results of parsing a scan.
// This parsing is done to convert the results of security and privacy checks into a format that
// can be sent to the Guardian360 Lighthouse API.
type ParseResult struct {
	Metadata Metadata    `json:"metadata"`
	Results  []IssueData `json:"results"`
}

// Metadata is a struct that contains metadata about the scan.
// This metadata includes the workstation ID, the user who initiated the scan, and the date of the scan.
type Metadata struct {
	WorkStationID int    `json:"workStationID"`
	User          string `json:"user"`
	Date          string `json:"date"`
}

// IssueData is a struct that encapsulates the data for a checked issue.
// This data includes the issue ID, whether the issue is considered a problem, and any additional data related to the
// issue.
type IssueData struct {
	IssueID        int      `json:"issueID"`
	Detected       bool     `json:"detected"`
	AdditionalData []string `json:"additionalData"`
}

// String returns a string representation of the ParseResult struct.
//
// Parameters: None.
//
// Returns:
//   - string: A string representation of the ParseResult struct.
func (p ParseResult) String() string {
	return fmt.Sprintf("Metadata: %v, Results: %v", p.Metadata, p.Results)
}

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

// SendResultsToAPI sends the results of a scan to the Guardian360 Lighthouse API.
//
// Parameters: result ParseResult: The result of parsing the scan.
//
// Returns: None.
func SendResultsToAPI(result ParseResult) {
	url := "https://localhost"
	jsonData, err := json.Marshal(result)
	if err != nil {
		logger.Log.ErrorWithErr("Error marshalling JSON:", err)
		return
	}

	buffer := bytes.NewBuffer(jsonData)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buffer)
	if err != nil {
		logger.Log.ErrorWithErr("Error creating request:", err)
		return
	}

	settings := usersettings.LoadUserSettings()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.IntegrationKey)
	req.Header.Set("Content-Length", strconv.Itoa(buffer.Len()))

	client := &http.Client{
		Timeout: 60 * time.Second, // Increase timeout for large payloads
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.ErrorWithErr("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	logger.Log.Debug("Response Status:" + resp.Status)
}
