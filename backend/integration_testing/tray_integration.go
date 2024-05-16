package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationScanNowSuccessful(t *testing.T) {
	result, err := tray.ScanNow()
	require.NotEmpty(t, result)
	require.NoError(t, err)
}
