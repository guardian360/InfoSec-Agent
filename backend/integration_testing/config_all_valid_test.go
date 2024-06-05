package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"os"
	"testing"
)

var testsValid = []func(t *testing.T){
	i.TestIntegrationFirefoxFolderExists,
	i.TestIntegrationGetDefaultDirExists,
	i.TestIntegrationCurrentUsernameFound,
	i.TestIntegrationExtensionsChromiumWithAdBlocker,
	i.TestIntegrationHistoryChromiumWithoutPhishing,
	i.TestIntegrationSearchEngineChromiumWithSearchEngine,
	i.TestIntegrationBluetoothNoDevices,
	i.TestIntegrationExternalDevicesNoDevices,
	i.TestIntegrationExtensionsFirefoxWithAdBlocker,
	i.TestIntegrationHistoryFirefoxWithoutPhishing,
	i.TestIntegrationSearchEngineFirefoxWithSearchEngine,
	i.TestIntegrationSmbCheckGoodSetup,
	i.TestIntegrationPasswordManagerPresent,
	i.TestIntegrationAdvertisementNotActive,
	i.TestIntegrationAutomatedLoginNotActive,
	i.TestIntegrationDefenderAllActive,
	i.TestIntegrationGuestAccountNotActive,
	i.TestIntegrationLastPasswordChangeValid,
	i.TestIntegrationLoginMethodPasswordOnly,
	// TODO: turn back on when the test is fixed
	// i.TestIntegrationOutdatedWin11UpToDate,
	i.TestIntegrationPermissionWithoutApps,
	// TODO: turn back on when the test is fixed
	i.TestIntegrationRemoteDesktopDisabled,
	i.TestIntegrationSecureBootEnabled,
	i.TestIntegrationStartupWithoutApps,
	i.TestIntegrationUACFullEnabled,
	i.TestIntegrationScanNowSuccessful,
	i.TestIntegrationScanSuccess,
	i.TestIntegrationCookiesFirefoxWithoutCookies,
	i.TestIntegrationCookiesChromiumWithoutCookies,
	i.TestIntegrationRemoteRPCDisabled,
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
