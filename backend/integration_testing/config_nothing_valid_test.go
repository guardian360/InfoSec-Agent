package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"testing"
)

var testsInvalid = []func(t *testing.T){
	i.TestIntegrationExtensionsChromiumWithoutAdBlocker,
	i.TestIntegrationHistoryChromiumWithPhishing,
	i.TestIntegrationCISRegistrySettingsIncorrect,
	i.TestIntegrationExternalDevicesDevices,
	i.TestIntegrationExtensionsFirefoxWithoutAdBlocker,
	i.TestIntegrationHistoryFirefoxWithPhishing,
	i.TestIntegrationOpenPortsPorts,
	i.TestIntegrationSmbCheckBadSetup,
	i.TestIntegrationPasswordManagerNotPresent,
	i.TestIntegrationAdvertisementActive,
	i.TestIntegrationAutomatedLoginActive,
	i.TestIntegrationDefenderAllNotActive,
	i.TestIntegrationGuestAccountActive,
	i.TestIntegrationLastPasswordChangeInvalid,
	i.TestIntegrationLoginMethodPINOnly,
	i.TestIntegrationOutdatedWin11NotUpToDate,
	i.TestIntegrationPermissionWithApps,
	i.TestIntegrationRemoteDesktopEnabled,
	i.TestIntegrationSecureBootDisabled,
	i.TestIntegrationStartupWithApps,
	i.TestIntegrationUACDisabled,
	i.TestIntegrationCookiesFirefoxWithCookies,
	i.TestIntegrationCookiesChromiumWithCookies,
	i.TestIntegrationRemoteRPCEnabled,
}

func TestAllInvalid(t *testing.T) {
	for _, test := range testsInvalid {
		test(t)
	}
}
