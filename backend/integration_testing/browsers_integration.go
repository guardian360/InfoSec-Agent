// Package integration contains integration tests for the entire project.
// The integration tests are run on Virtual Machines with a custom configuration to ensure that the project works as expected in different environments.
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

func TestIntegrationGetDefaultDirExists(t *testing.T) {
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
			result, err := browsers.RealDefaultDirGetter{}.GetDefaultDir(tt.browser)
			require.NotEmpty(t, result)
			require.NoError(t, err)
		})
	}
}
