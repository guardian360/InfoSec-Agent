package apiconnection_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	apiconnection "github.com/InfoSec-Agent/InfoSec-Agent/backend/api_connection"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCheckResult(t *testing.T) {
	tests := []struct {
		name  string
		check checks.Check
		want  apiconnection.IssueData
	}{
		{
			name: "Check error",
			check: checks.Check{
				IssueID: 1,
				Error:   errors.New("error"),
			},
			want: apiconnection.IssueData{IssueID: 1, Detected: false},
		},
		{
			name: "Check no error",
			check: checks.Check{
				IssueID:  3,
				ResultID: 1,
				Error:    nil,
			},
			want: apiconnection.IssueData{IssueID: 3, Detected: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiconnection.ParseCheckResult(tt.check)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseScanResults(t *testing.T) {
	tests := []struct {
		name     string
		metaData apiconnection.Metadata
		checks   []checks.Check
		want     apiconnection.ParseResult
	}{
		{
			name: "Empty checks",
			metaData: apiconnection.Metadata{
				WorkStationID: 1,
				User:          "test",
				Date:          "2021-01-01",
			},
			checks: []checks.Check{},
			want: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: nil,
			},
		},
		{
			name: "Non-empty checks",
			metaData: apiconnection.Metadata{
				WorkStationID: 1,
				User:          "test",
				Date:          "2021-01-01",
			},
			checks: []checks.Check{
				{
					IssueID:  3,
					ResultID: 1,
					Error:    nil,
				}},
			want: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: []apiconnection.IssueData{
					{
						IssueID:  3,
						Detected: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiconnection.ParseScanResults(tt.metaData, tt.checks)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name string
		p    apiconnection.ParseResult
		want string
	}{
		{
			name: "Empty results",
			p: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: nil,
			},
			want: "Metadata: {1 test 2021-01-01}, Results: []",
		},
		{
			name: "Non-empty results",
			p: apiconnection.ParseResult{
				Metadata: apiconnection.Metadata{
					WorkStationID: 1,
					User:          "test",
					Date:          "2021-01-01",
				},
				Results: []apiconnection.IssueData{
					{
						IssueID:  3,
						Detected: true,
					},
				},
			},
			want: "Metadata: {1 test 2021-01-01}, Results: [{3 true []}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.String()
			require.Equal(t, tt.want, got)
		})
	}
}

// Create a ParseResult instance
type ParseResult struct {
	Status string `json:"status"`
}

func TestSendResultsToAPI(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer", r.Header.Get("Authorization"))

		var result ParseResult

		// Validate the ParseResult fields
		assert.Equal(t, "success", result.Status)

		// Send a response
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Override the URL for the test
	url := testServer.URL

	// Create a ParseResult instance
	result := ParseResult{
		Status: "success",
	}

	// Convert result to JSON
	jsonData, err := json.Marshal(result)
	require.NoError(t, err)

	// Act
	buffer := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest(http.MethodPost, url, buffer)
	require.NoError(t, err)

	settings := usersettings.DefaultUserSettings
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.IntegrationKey)
	req.Header.Set("Content-Length", strconv.Itoa(buffer.Len()))

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
