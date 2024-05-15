package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestIntegrationFirefoxFolderExists(t *testing.T) {
	result, err := browsers.RealProfileFinder{}.FirefoxFolder()
	require.NotEmpty(t, result)
	require.NoError(t, err)
}

func TestIntegrationFirefoxFolderNotExists(t *testing.T) {
	result, err := browsers.RealProfileFinder{}.FirefoxFolder()
	require.Empty(t, result)
	require.Error(t, err)
}

func TestIntegrationGetPreferencesDirExists(t *testing.T) {
	tests := []struct {
		name    string
		browser string
	}{
		{
			name:    "Chrome preferences directory exists",
			browser: "Chrome",
		},
		{
			name:    "Edge preferences directory exists",
			browser: "Edge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := browsers.RealPreferencesDirGetter{}.GetPreferencesDir(tt.browser)
			require.NotEmpty(t, result)
			require.NoError(t, err)
		})

	}
}
