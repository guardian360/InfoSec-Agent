// Package integration contains integration tests for the entire project.
// The integration tests are run on Virtual Machines with a custom configuration and are not meant
// to succeed on any machine.
package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/stretchr/testify/require"
)

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
