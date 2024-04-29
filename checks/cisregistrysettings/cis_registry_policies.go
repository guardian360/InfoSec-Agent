package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/mocking"

// CheckPoliciesHKU is a function that checks various registry settings related to different policies
// to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_USERS for this function.
//
// Returns:
//
//   - []bool: A slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
func CheckPoliciesHKU(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)

	for _, check := range policyChecksHKU {
		results = append(results, check(registryKey)...)
	}

	return results
}

// policyChecksHKU is a collection of registry checks related to different policies.
// Each function in the collection represents a different policy check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var policyChecksHKU = []func(mocking.RegistryKey) []bool{
	policiesAttachments,
	policiesExplorerHKU,
	policiesCloudContentHKU,
	policiesControlPanelDesktop,
	policiesPushNotifications,
	policiesInstallerHKU,
}

// policiesAttachments is a helper function that checks the registry to determine if the system is configured with the correct settings for attachments.
//
// CIS Benchmark Audit list indices: 19.7.4.1-2
func policiesAttachments(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Attachments`

	settings := []string{"SaveZoneInformation", "ScanWithAntiVirus"}

	expectedValues := []interface{}{uint64(2), uint64(3)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesExplorer is a helper function that checks the registry to determine if the system is configured with the correct settings for the Explorer policies.
//
// CIS Benchmark Audit list index: 19.7.28.1
func policiesExplorerHKU(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`

	settings := []string{"NoInplaceSharing"}

	expectedValues := []interface{}{uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesCloudContentHKU is a helper function that checks the registry to determine if the system is configured with the correct settings for cloud content.
//
// CIS Benchmark Audit list indices: 19.7.8.1-2, 19.7.8.5
func policiesCloudContentHKU(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CloudContent`

	settings := []string{"ConfigureWindowsSpotlight", "DisableThirdPartySuggestions", "DisableSpotlightCollectionOnDesktop"}

	expectedValues := []interface{}{uint64(2), uint64(1), uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesControlPanelDesktop is a helper function that checks the registry to determine if the system is configured with the correct settings for the Control Panel Desktop.
//
// CIS Benchmark Audit list indices: 19.1.3.1-3
func policiesControlPanelDesktop(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Control Panel\Desktop`

	settings := []string{"ScreenSaveActive", "ScreenSaverIsSecure", "ScreenSaveTimeOut"}

	expectedValues := []interface{}{uint64(1), uint64(1), []uint64{0, 900}}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPushNotifications is a helper function that checks the registry to determine if the system is configured with the correct settings for push notifications.
//
// CIS Benchmark Audit list index: 19.5.1.1
func policiesPushNotifications(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`

	settings := []string{"NoToastApplicationNotificationOnLockScreen"}

	expectedValues := []interface{}{uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesInstallerHKU is a helper function that checks the registry to determine if the system is configured with the correct settings for the Installer.
//
// CIS Benchmark Audit list index: 19.7.43.1
func policiesInstallerHKU(registryKey mocking.RegistryKey) []bool {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Installer`

	settings := []string{"AlwaysInstallElevated"}

	expectedValues := []interface{}{uint64(0)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
