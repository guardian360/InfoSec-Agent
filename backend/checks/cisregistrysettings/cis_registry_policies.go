package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

// TODO: Update documentation
// CheckPoliciesHKU is a function that checks various registry settings related to different policies
// to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_USERS for this function.
//
// Returns: None
func CheckPoliciesHKU(registryKey mocking.RegistryKey) {
	for _, check := range policyChecksHKU {
		check(registryKey)
	}
}

// TODO: Update documentation
// policyChecksHKU is a collection of registry checks related to different policies.
// Each function in the collection represents a different policy check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var policyChecksHKU = []func(mocking.RegistryKey){
	policiesAttachments,
	policiesExplorerHKU,
	policiesCloudContentHKU,
	policiesControlPanelDesktop,
	policiesPushNotifications,
	policiesInstallerHKU,
}

// TODO: Update documentation
// policiesAttachments is a helper function that checks the registry to determine if the system is configured with the correct settings for attachments.
//
// CIS Benchmark Audit list indices: 19.7.4.1-2
func policiesAttachments(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Attachments`

	settings := []string{"SaveZoneInformation", "ScanWithAntiVirus"}

	expectedValues := []interface{}{uint64(2), uint64(3)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// TODO: Update documentation
// policiesExplorer is a helper function that checks the registry to determine if the system is configured with the correct settings for the Explorer policies.
//
// CIS Benchmark Audit list index: 19.7.28.1
func policiesExplorerHKU(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`

	settings := []string{"NoInplaceSharing"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// TODO: Update documentation
// policiesCloudContentHKU is a helper function that checks the registry to determine if the system is configured with the correct settings for cloud content.
//
// CIS Benchmark Audit list indices: 19.7.8.1-2, 19.7.8.5
func policiesCloudContentHKU(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CloudContent`

	settings := []string{"ConfigureWindowsSpotlight", "DisableThirdPartySuggestions", "DisableSpotlightCollectionOnDesktop"}

	expectedValues := []interface{}{uint64(2), uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// TODO: Update documentation
// policiesControlPanelDesktop is a helper function that checks the registry to determine if the system is configured with the correct settings for the Control Panel Desktop.
//
// CIS Benchmark Audit list indices: 19.1.3.1-3
func policiesControlPanelDesktop(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Control Panel\Desktop`

	settings := []string{"ScreenSaveActive", "ScreenSaverIsSecure", "ScreenSaveTimeOut"}

	expectedValues := []interface{}{uint64(1), uint64(1), []uint64{0, 900}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// TODO: Update documentation
// policiesPushNotifications is a helper function that checks the registry to determine if the system is configured with the correct settings for push notifications.
//
// CIS Benchmark Audit list index: 19.5.1.1
func policiesPushNotifications(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`

	settings := []string{"NoToastApplicationNotificationOnLockScreen"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// TODO: Update documentation
// policiesInstallerHKU is a helper function that checks the registry to determine if the system is configured with the correct settings for the Installer.
//
// CIS Benchmark Audit list index: 19.7.43.1
func policiesInstallerHKU(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Installer`

	settings := []string{"AlwaysInstallElevated"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
