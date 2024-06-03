package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"testing"
)

var testsNotPresent = []func(t *testing.T){
	i.TestIntegrationFirefoxFolderNotExists,
	i.TestIntegrationExtensionsChromiumNotInstalled,
	i.TestIntegrationHistoryChromiumNotInstalled,
	i.TestIntegrationSearchEngineChromiumNotInstalled,
	i.TestIntegrationSearchEngineFirefoxNotInstalled,
	i.TestIntegrationHistoryFirefoxNotInstalled,
	i.TestIntegrationExtensionsFirefoxNotInstalled,
	i.TestIntegrationDefenderPeriodicScanActive,
	i.TestIntegrationLoginMethodPasswordAndPIN,
	//i.TestIntegrationOutdatedWin10UpToDate,
	i.TestIntegrationUACPartialEnabled,
}

func TestNotPresent(t *testing.T) {
	for _, test := range testsNotPresent {
		test(t)
	}
}
