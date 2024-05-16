package integration_testing

import (
	"testing"
)

var testsNotPresent = []func(t *testing.T){
	TestIntegrationFirefoxFolderNotExists,
	TestIntegrationExtensionsChromiumNotInstalled,
	TestIntegrationHistoryChromiumNotInstalled,
	TestIntegrationSearchEngineChromiumNotInstalled,
	TestIntegrationSearchEngineFirefoxNotInstalled,
	TestIntegrationHistoryFirefoxNotInstalled,
	TestIntegrationExtensionsFirefoxNotInstalled,
	TestIntegrationDefenderPeriodicScanActive,
	TestIntegrationLoginMethodPasswordAndPIN,
	TestIntegrationOutdatedWin10UpToDate,
	TestIntegrationUACPartialEnabled,
}

func TestNotPresent(t *testing.T) {
	for _, test := range testsNotPresent {
		test(t)
	}
}
