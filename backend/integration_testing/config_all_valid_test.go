package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"os"
	"testing"
)

var testsValid = []func(t *testing.T){
	i.TestIntegrationFirefoxFolderExists,
	i.TestIntegrationGetPreferencesDirExists,
	i.TestIntegrationCurrentUsernameFound,
	i.TestIntegrationExtensionsChromiumWithAdBlocker,
	i.TestIntegrationHistoryChromiumWithoutPhishing,
	i.TestIntegrationSearchEngineChromiumWithSearchEngine,
	i.TestIntegrationBluetoothNoDevices,
	i.TestIntegrationExternalDevicesNoDevices,
	i.TestIntegrationExtensionsFirefoxWithAdBlocker,
	i.TestIntegrationHistoryFirefoxWithoutPhishing,
	i.TestIntegrationSearchEngineFirefoxWithSearchEngine,
	i.TestIntegrationOpenPortsNoPorts,
	i.TestIntegrationSmbCheckGoodSetup,
	i.TestIntegrationPasswordManagerPresent,
	i.TestIntegrationAdvertisementNotActive,
	i.TestIntegrationAutomatedLoginNotActive,
	i.TestIntegrationDefenderAllActive,
	i.TestIntegrationGuestAccountNotActive,
	i.TestIntegrationLastPasswordChangeValid,
	i.TestIntegrationLoginMethodPasswordOnly,
	i.TestIntegrationOutdatedWin11UpToDate,
	i.TestIntegrationPermissionWithoutApps,
	i.TestIntegrationRemoteDesktopDisabled,
	i.TestIntegrationSecureBootEnabled,
	i.TestIntegrationStartupWithoutApps,
	i.TestIntegrationUACFullEnabled,
	i.TestIntegrationScanNowSuccessful,
	i.TestIntegrationScanSuccess,
}

func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestAllValid(t *testing.T) {
	for _, test := range testsValid {
		test(t)
	}
}
