package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/stretchr/testify/require"
)

func TestIntegrationScanNowSuccessful(t *testing.T) {
	result, err := tray.ScanNow(false, "reporting-page/frontend/src/databases/database.en-GB.json")
	require.NotEmpty(t, result)
	require.NoError(t, err)
}
