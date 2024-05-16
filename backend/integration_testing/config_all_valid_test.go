package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"os"
	"testing"
)

var testsValid = []func(t *testing.T){
	TestIntegrationFirefoxFolderExists,
	TestIntegrationGetPreferencesDirExists,
	TestIntegrationCurrentUsernameFound,
	TestIntegrationExtensionsChromiumWithAdBlocker,
	TestIntegrationHistoryChromiumWithoutPhishing,
	TestIntegrationSearchEngineChromiumWithSearchEngine,
	TestIntegrationBluetoothNoDevices,
	TestIntegrationExternalDevicesNoDevices,
	TestIntegrationExtensionsFirefoxWithAdBlocker,
	TestIntegrationHistoryFirefoxWithoutPhishing,
	TestIntegrationSearchEngineFirefoxWithSearchEngine,
	TestIntegrationOpenPortsNoPorts,
	TestIntegrationSmbCheckGoodSetup,
	TestIntegrationPasswordManagerPresent,
	TestIntegrationAdvertisementNotActive,
	TestIntegrationAutomatedLoginNotActive,
	TestIntegrationDefenderAllActive,
	TestIntegrationGuestAccountNotActive,
	TestIntegrationLastPasswordChangeValid,
	TestIntegrationLoginMethodPasswordOnly,
	TestIntegrationOutdatedWin11UpToDate,
	TestIntegrationPermissionWithoutApps,
	TestIntegrationRemoteDesktopDisabled,
	TestIntegrationSecureBootEnabled,
	TestIntegrationStartupWithoutApps,
	TestIntegrationUACFullEnabled,
	TestIntegrationScanNowSuccessful,
	TestIntegrationScanSuccess,
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
