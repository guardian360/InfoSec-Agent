package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/stretchr/testify/require"
)

func TestIntegrationExtensionsFirefoxWithAdBlocker(t *testing.T) {
	result, adblock := firefox.ExtensionFirefox(browsers.RealProfileFinder{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
	require.Equal(t, 0, adblock.ResultID)
}

func TestIntegrationExtensionsFirefoxWithoutAdBlocker(t *testing.T) {
	result, adblock := firefox.ExtensionFirefox(browsers.RealProfileFinder{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
	require.Equal(t, 1, adblock.ResultID)
}

func TestIntegrationExtensionsFirefoxNotInstalled(t *testing.T) {
	result, adblock := firefox.ExtensionFirefox(browsers.RealProfileFinder{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, adblock)
	require.Equal(t, -1, result.ResultID)
	require.Equal(t, -1, adblock.ResultID)
}

func TestIntegrationHistoryFirefoxWithoutPhishing(t *testing.T) {
	result := firefox.HistoryFirefox(browsers.RealProfileFinder{}, browsers.RealPhishingDomainGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationHistoryFirefoxWithPhishing(t *testing.T) {
	result := firefox.HistoryFirefox(browsers.RealProfileFinder{}, browsers.RealPhishingDomainGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationHistoryFirefoxNotInstalled(t *testing.T) {
	result := firefox.HistoryFirefox(browsers.RealProfileFinder{}, browsers.RealPhishingDomainGetter{})
	require.Equal(t, -1, result.ResultID)
	require.NotEmpty(t, result)
}

func TestIntegrationSearchEngineFirefoxWithSearchEngine(t *testing.T) {
	result := firefox.SearchEngineFirefox(browsers.RealProfileFinder{}, false, nil, nil)
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationSearchEngineFirefoxNotInstalled(t *testing.T) {
	result := firefox.SearchEngineFirefox(browsers.RealProfileFinder{}, false, nil, nil)
	require.Equal(t, -1, result.ResultID)
	require.NotEmpty(t, result)
}
