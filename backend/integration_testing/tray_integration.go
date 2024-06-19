package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/stretchr/testify/require"
)

func TestIntegrationScanNowSuccessful(t *testing.T) {
	result, err := tray.ScanNow(false, config.DatabasePath)
	require.NotEmpty(t, result)
	require.NoError(t, err)
}
