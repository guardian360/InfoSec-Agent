package integration

import (
	"testing"

	apiconnection "github.com/InfoSec-Agent/InfoSec-Agent/backend/api_connection"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"
	"github.com/ncruces/zenity"
	"github.com/stretchr/testify/require"
)

func TestIntegrationScanSuccess(t *testing.T) {
	// Display a progress dialog while the scan is running
	dialog, err := zenity.Progress(
		zenity.Title("Security/Privacy Scan"))
	if err != nil {
		logger.Log.ErrorWithErr("Error creating dialog during test:", err)
		return
	}
	// Defer closing the dialog until the scan completes
	defer func(dialog zenity.ProgressDialog) {
		err = dialog.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing dialog during test:", err)
		}
	}(dialog)

	// Run the scan
	checks, err := scan.Scan(dialog)
	require.NotEmpty(t, checks)
	require.NoError(t, err)
	totalLength := 0
	for _, checkSlice := range scan.ChecksList {
		totalLength += len(checkSlice)
	}
	require.Len(t, totalLength, len(checks))

	// Get database data
	data, err := scan.GetDataBaseData(checks, "../../reporting-page/database.db")
	if err != nil {
		logger.Log.ErrorWithErr("Error getting database data during test:", err)
		return
	}
	require.NotEmpty(t, data)
	require.NoError(t, err)
	require.Equal(t, len(checks), len(data))

	// Parse scan results
	metaData := apiconnection.Metadata{WorkStationID: 0, User: "user", Date: "2022-01-01T00:00:00Z"}
	parseResult := apiconnection.ParseScanResults(metaData, checks)
	require.NotEmpty(t, parseResult)
	require.Equal(t, metaData, parseResult.Metadata)
	require.Equal(t, len(checks), len(parseResult.Results))
}
