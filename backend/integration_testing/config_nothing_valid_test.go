package integration_testing

import "testing"

var testsInvalid = []func(t *testing.T){
	TestIntegrationExtensionsChromiumWithoutAdBlocker,
	TestIntegrationHistoryChromiumWithPhishing,
	TestIntegrationCISRegistrySettingsIncorrect,
	TestIntegrationBluetoothDevices,
	TestIntegrationExternalDevicesDevices,
	TestIntegrationExtensionsFirefoxWithoutAdBlocker,
	TestIntegrationHistoryFirefoxWithPhishing,
	TestIntegrationOpenPortsPorts,
	TestIntegrationSmbCheckBadSetup,
	TestIntegrationPasswordManagerNotPresent,
	TestIntegrationAdvertisementActive,
	TestIntegrationAutomatedLoginActive,
	TestIntegrationDefenderAllNotActive,
	TestIntegrationGuestAccountActive,
	TestIntegrationLastPasswordChangeInvalid,
	TestIntegrationLoginMethodPINOnly,
	TestIntegrationOutdatedWin11NotUpToDate,
	TestIntegrationPermissionWithApps,
	TestIntegrationRemoteDesktopEnabled,
	TestIntegrationSecureBootDisabled,
	TestIntegrationStartupWithApps,
	TestIntegrationUACDisabled,
}

func TestAllInvalid(t *testing.T) {
	for _, test := range testsInvalid {
		test(t)
	}
}
