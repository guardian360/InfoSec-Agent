package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"testing"
)

var testsInvalid = []func(t *testing.T){
	i.TestIntegrationExtensionsChromiumWithoutAdBlocker,
	// TODO: Look into this
	i.TestIntegrationHistoryChromiumWithPhishing,
	i.TestIntegrationCISRegistrySettingsIncorrect,
	i.TestIntegrationExtensionsFirefoxWithoutAdBlocker,
	// TODO: Look into this
	i.TestIntegrationHistoryFirefoxWithPhishing,
	i.TestIntegrationOpenPortsPorts,
	i.TestIntegrationSmbCheckBadSetup,
	i.TestIntegrationPasswordManagerNotPresent,
	i.TestIntegrationAdvertisementActive,
	i.TestIntegrationAutomatedLoginActive,
	i.TestIntegrationGuestAccountActive,
	i.TestIntegrationLoginMethodPINOnly,
	i.TestIntegrationOutdatedWinNotUpToDate,
	i.TestIntegrationPermissionWithApps,
	i.TestIntegrationRemoteDesktopEnabled,
	i.TestIntegrationSecureBootDisabled,
	i.TestIntegrationStartupWithApps,
	i.TestIntegrationUACDisabled,
	i.TestIntegrationCookiesFirefoxWithCookies,
	i.TestIntegrationCookiesChromiumWithCookies,
	i.TestIntegrationRemoteRPCEnabled,
	i.TestIntegrationNetBIOSEnabled,
	i.TestIntegrationWPADEnabled,
	i.TestIntegrationCredentialGuardDisabled,
	i.TestIntegrationFirewallDisabled,
	i.TestIntegrationPasswordComplexityInvalid,
	i.TestIntegrationScreenLockDisabled,
}

func TestAllInvalid(t *testing.T) {
	for _, test := range testsInvalid {
		test(t)
	}
}
