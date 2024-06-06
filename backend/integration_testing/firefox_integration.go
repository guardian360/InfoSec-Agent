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

func TestIntegrationHistoryFirefoxWithoutPhishing(t *testing.T) {
	result := firefox.HistoryFirefox(browsers.RealProfileFinder{}, browsers.RealPhishingDomainGetter{}, firefox.RealQueryDatabaseGetter{}, firefox.RealProcessQueryResultsGetter{}, firefox.RealCopyDBGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationHistoryFirefoxWithPhishing(t *testing.T) {
	result := firefox.HistoryFirefox(browsers.RealProfileFinder{}, browsers.RealPhishingDomainGetter{}, firefox.RealQueryDatabaseGetter{}, firefox.RealProcessQueryResultsGetter{}, firefox.RealCopyDBGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationSearchEngineFirefoxWithSearchEngine(t *testing.T) {
	result := firefox.SearchEngineFirefox(browsers.RealProfileFinder{}, false, nil, nil)
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationCookiesFirefoxWithCookies(t *testing.T) {
	result := firefox.CookiesFirefox(browsers.RealProfileFinder{}, browsers.RealCopyFileGetter{}, browsers.RealQueryCookieDatabaseGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationCookiesFirefoxWithoutCookies(t *testing.T) {
	result := firefox.CookiesFirefox(browsers.RealProfileFinder{}, browsers.RealCopyFileGetter{}, browsers.RealQueryCookieDatabaseGetter{})
	require.NotEqual(t, -1, result.ResultID)
	require.Empty(t, result)
	require.Equal(t, 0, result.ResultID)
}
