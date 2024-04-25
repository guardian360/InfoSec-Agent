package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

//TODO REPLACE return []bool{} with []bool with correct number of false elements

func CISRegistrySettings(registryKey mocking.RegistryKey) Check {

	return Check{}
}

// AutoConnectHotspot is a helper function that checks the registry to determine if the system is configured to automatically connect to hotspots.
//
// CIS Benchmark Audit list index: 18.5.23.2.1
func AutoConnectHotspot(registryKey mocking.RegistryKey) bool {
	key, err := mocking.OpenRegistryKey(registryKey,
		`SOFTWARE\Microsoft\WcmSvc\wifinetworkmanager\config`)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening AutoConnectHotspot registry key", err)
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "AutoConnectAllowedOEM", 0)
}

// CheckCurrentVersionRegistry is a helper function that checks the registry to determine if the system is configured with the correct settings for the current version.
//
// CIS Benchmark Audit list indices: 2.3.4.1, 2.3.7.8, 2.3.7.9, 18.2.1, 18.4.1, 18.4.10
func CheckCurrentVersionRegistry(registryKey mocking.RegistryKey) []bool {
	result := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"AllocateDASD",
		"PasswordExpiryWarning",
		"ScRemoveOption",
		"AutoAdminLogon",
		"ScreenSaverGracePeriod",
	}

	expectedValues := []interface{}{uint64(2), []uint64{5, 14}, []uint64{1, 2, 3}, uint64(0), []uint64{0, 5}}

	result = append(result, checkMultipleIntegerValues(key, settings, expectedValues)...)

	key, err = openRegistryKeyWithErrHandling(registryKey,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\GPExtensions\{D76B9641-3288-4f75-942D-087DE603E3EA}`)
	if err != nil {
		result = append(result, false)
		return result
	}
	defer mocking.CloseRegistryKey(key)

	result = append(result,
		checkStringValue(key, "DllName", "C:\\Program Files\\LAPS\\CSE\\AdmPwd.dll"))
	return result
}

// EnumerateAdminAccount is a helper function that checks the registry to determine if the system is configured to enumerate administrator accounts.
//
// CIS Benchmark Audit list index: 18.9.16.2
func EnumerateAdminAccount(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Microsoft\Windows\Currentversion\Policies\Credui`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "EnumerateAdministrators", 0)
}

// CheckExplorerPolicies is a helper function that checks the registry to determine if the system is configured with the correct settings for Explorer policies.
//
// CIS Benchmark Audit list indices: 18.8.22.1.6, 18.9.8.2, 18.9.8.3, 18.9.31.4
func CheckExplorerPolicies(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"NoWebServices",
		"NoAutorun",
		"NoDriveTypeAutoRun",
		"PreXPSP2ShellProtocolBehavior",
	}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(255), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckSystemPolicies is a helper function that checks the registry to determine if the system is configured with the correct settings for system policies.
//
// CIS Benchmark Audit list indices: 2.3.1.2, 2.3.7.1-3, 2.3.11.4, 2.3.17.1-8, 18.3.1, 18.8.3.1, 18.8.4.1, 18.9.6.1, 18.9.91.1
func CheckSystemPolicies(registryKey mocking.RegistryKey) []bool {
	result := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"NoConnectedUser",
		"DisableCAD",
		"DontDisplayLastUserName",
		"InactivityTimeoutSecs",
		"FilterAdministratorToken",
		"ConsentPromptBehaviorAdmin",
		"ConsentPromptBehaviorUser",
		"EnableInstallerDetection",
		"EnableSecureUIAPaths",
		"EnableLUA",
		"PromptOnSecureDesktop",
		"EnableVirtualization",
		"LocalAccountTokenFilterPolicy",
		"MSAOptional",
		"DisableAutomaticRestartSignOn",
	}
	expectedValues := []interface{}{uint64(3), uint64(0), uint64(1), []uint64{0, 900}, uint64(1), uint64(2), uint64(0),
		uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(0), uint64(1), uint64(1)}

	result = append(result, checkMultipleIntegerValues(key, settings, expectedValues)...)

	subKeys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\Kerberos\Parameters`,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\Audit`,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\CredSSP\Parameters`,
	}

	subKeysSettings := []string{
		"SupportedEncryptionTypes",
		"ProcessCreationIncludeCmdLine_Enabled",
		"AllowEncryptionOracle",
	}

	subKeysExpected := []interface{}{uint64(2147483640), uint64(1), uint64(0)}

	for i, subKey := range subKeys {
		func() {
			key, err = openRegistryKeyWithErrHandling(registryKey, subKey)
			if err != nil {
				result = append(result, false)
				return
			}
			defer mocking.CloseRegistryKey(key)
			result = append(result, checkIntegerValue(key, subKeysSettings[i], subKeysExpected[i]))
		}()
	}
	return result
}

// CheckAdminPassword is a helper function that checks the registry to determine if the system is configured with the correct settings for the administrator password.
//
// CIS Benchmark Audit list indices: 18.2.2-6
func CheckAdminPassword(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft Services\AdmPwd`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"PwdExpirationProtectionEnabled",
		"AdmPwdEnabled",
		"PasswordComplexity",
		"PasswordLength",
		"PasswordAgeDays",
	}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(4), []uint64{15, ^uint64(0)}, []uint64{0, 30}}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckEnhancedAntiSpoofing is a helper function that checks the registry to determine if the system is configured with the correct settings for enhanced anti-spoofing.
//
// CIS Benchmark Audit list index: 18.9.10.1.1
func CheckEnhancedAntiSpoofing(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Biometrics\FacialFeatures`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "EnhancedAntiSpoofing", 1)
}

// CheckWidgetAllowance is a helper function that checks the registry to determine if the system is configured to allow widgets.
//
// CIS Benchmark Audit list index: 18.9.81.1
func CheckWidgetAllowance(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Dsh`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "AllowNewsAndInterests", 0)
}

// CheckOnlineSpeechRecognitionServices is a helper function that checks the registry to determine if the system is configured to allow online speech recognition services.
//
// CIS Benchmark Audit list index: 18.1.2.2
func CheckOnlineSpeechRecognitionServices(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\InputPersonalization`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "AllowInputPersonalization", 0)
}

// CheckDownloadEnclosures is a helper function that checks the registry to determine if the system is configured to download enclosures.
//
// CIS Benchmark Audit list index: 18.9.66.1
func CheckDownloadEnclosures(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Internet Explorer\Feeds`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "DisableEnclosureDownload", 1)
}

// CheckBlockConsumerUserAuthentication is a helper function that checks the registry to determine if the system is configured to block consumer user authentication.
//
// CIS Benchmark Audit list index: 18.9.46.1
func CheckBlockConsumerUserAuthentication(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\MicrosoftAccount`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)
	return checkIntegerValue(key, "DisableUserAuth", 1)
}

// CheckPhishingFilter is a helper function that checks the registry to determine if the system is configured with the correct settings for the phishing filter.
//
// CIS Benchmark Audit list indices: 18.9.85.2.1-2
func CheckPhishingFilter(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\MicrosoftEdge\PhishingFilter`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"EnabledV9",
		"PreventOverride",
	}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)

}

// CheckPowerSettings is a helper function that checks the registry to determine if the system is configured with the correct power settings.
//
// CIS Benchmark Audit list indices: 18.8.34.6.1-2, 18.8.34.6.5-6
func CheckPowerSettings(registryKey mocking.RegistryKey) []bool {
	result := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Power\PowerSettings\f15576e8-98b7-4186-b944-eafa664402d9`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"DCSettingIndex",
		"ACSettingIndex",
	}

	expectedValues := []interface{}{uint64(0), uint64(0)}
	result = append(result, checkMultipleIntegerValues(key, settings, expectedValues)...)

	key, err = openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Power\PowerSettings\0e796bdb-100d-47d6-a2d5-f7d2daa51f51`)
	if err != nil {
		result = append(result, false, false)
		return result
	}
	defer mocking.CloseRegistryKey(key)

	expectedValues = []interface{}{uint64(1), uint64(1)}
	result = append(result, checkMultipleIntegerValues(key, settings, expectedValues)...)
	return result
}

// CheckWindowsDefender is a helper function that checks the registry to determine if Windows Defender is configured with the correct settings.
//
// CIS Benchmark Audit list indices: 18.9.47.15-16
func CheckWindowsDefender(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{
		"PUAProtection",
		"DisableAntiSpyware",
	}

	expectedValues := []interface{}{uint64(1), uint64(0)}
	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWindowsDefenderScan(registryKey)...)
	results = append(results, checkWindowsDefenderRealTime(registryKey)...)
	results = append(results, checkWindowsDefenderASR(registryKey)...)
	results = append(results, checkWindowsDefenderSpyNet(registryKey),
		checkWindowsDefenderNetworkProtection(registryKey),
		checkWindowsDefenderAppBrowserProtection(registryKey))
	return results
}

// checkWindowsDefenderScan is a helper function that checks the registry to determine if Windows Defender is configured with the correct scan settings.
//
// CIS Benchmark Audit list indices: 18.9.47.12.1-2
func checkWindowsDefenderScan(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Scan`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"DisableRemovableDriveScanning", "DisableEmailScanning"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkWindowsDefenderRealTime is a helper function that checks the registry to determine if Windows Defender is configured with the correct real-time protection settings.
//
// CIS Benchmark Audit list indices: 18.9.47.9.1-4
func checkWindowsDefenderRealTime(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"DisableIOAVProtection", "DisableRealtimeMonitoring", "DisableBehaviorMonitoring", "DisableScriptScanning"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkWindowsDefenderASR is a helper function that checks the registry to determine if Windows Defender is configured with the correct ASR settings.
//
// CIS Benchmark Audit list indices: 18.9.47.5.1.1-2
func checkWindowsDefenderASR(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\ASR`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	results = append(results, checkIntegerValue(key, "ExploitGuard_ASR_Rules", 1))

	key, err = openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\ASR\Rules`)
	if err != nil {
		results = append(results, make([]bool, 11)...)
		return results
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"26190899-1602-49e8-8b27-eb1d0a1ce869",
		"3b576869-a4ec-4529-8536-b80a7769e899",
		"5beb7efe-fd9a-4556-801d-275e5ffc04cc",
		"75668c1f-73b5-4cf0-bb93-3ecf5cb7cc84",
		"7674ba52-37eb-4a4f-a9a1-f0f9a1619a2c",
		"92e97fa1-2edf-4476-bdd6-9dd0b4dddc7b",
		"9e6c4e1f-7d60-472f-ba1a-a39ef669e4b2",
		"b2b3f03d-6a65-4f7b-a9c7-1c7ef74a9ba4",
		"be9ba2d9-53ea-4cdc-84e5-9b1eeee46550",
		"d3e037e1-3eb8-44c8-a917-57927947596d",
		"d4f940ab-401b-4efc-aadc-ad5f3c50688a",
	}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(1),
		uint64(1), uint64(1), uint64(1), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	return results
}

// checkWindowsDefenderSpyNet is a helper function that checks the registry to determine if Windows Defender is configured with the correct SpyNet settings.
//
// CIS Benchmark Audit list index: 18.9.47.4.1
func checkWindowsDefenderSpyNet(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Spynet`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "LocalSettingOverrideSpynetReporting", 0)
}

// checkWindowsDefenderAppBrowserProtection is a helper function that checks the registry to determine if Windows Defender is configured with the correct app and browser protection settings.
//
// CIS Benchmark Audit list index: 18.9.105.2.1
func checkWindowsDefenderAppBrowserProtection(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender Security Center\App and Browser protection`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "DisallowExploitProtectionOverride", 1)
}

// checkWindowsDefenderNetworkProtection is a helper function that checks the registry to determine if Windows Defender is configured with the correct network protection settings.
//
// CIS Benchmark Audit list index: 18.9.47.5.3.1
func checkWindowsDefenderNetworkProtection(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\Network Protection`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "EnableNetworkProtection", 1)
}

// CheckDNSClient is a helper function that checks the registry to determine if the system is configured with the correct settings for the DNS client.
//
// CIS Benchmark Audit list indices: 18.5.4.1-2
func CheckDNSClient(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows NT\DNSClient`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"DoHPolicy", "EnableMulticast"}

	expectedValues := []interface{}{[]uint64{2, 3}, uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckPrinters is a helper function that checks the registry to determine if the system is configured with the correct settings for printers.
//
// CIS Benchmark Audit list indices: 18.3.5, 18.6.1-3, 18.8.22.1.2
func CheckPrinters(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows NT\Printers`)
	if err != nil {
		return []bool{}
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"RegisterSpoolerRemoteRpcEndPoint", "DisableWebPnPDownload"}

	expectedValues := []interface{}{uint64(2), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)

	key, err = openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows NT\Printers\PointAndPrint`)
	if err != nil {
		results = append(results, make([]bool, 3)...)
		return results
	}
	defer mocking.CloseRegistryKey(key)

	settings = []string{"RestrictDriverInstallationToAdministrators", "NoWarningNoElevationOnInstall", "UpdatePromptSettings"}

	expectedValues = []interface{}{uint64(1), uint64(0), uint64(0)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	return results
}

// CheckRPC is a helper function that checks the registry to determine if the system is configured with the correct settings for RPC.
//
// CIS Benchmark Audit list indices: 18.8.37.1-2
func CheckRPC(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows NT\Rpc`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableAuthEpResolution", "RestrictRemoteClients"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckTerminalServices is a helper function that checks the registry to determine if the system is configured with the correct settings for terminal services.
//
// CIS Benchmark Audit list indices: 18.8.36.1-2, 18.9.65.2.2, 18.9.65.3.3.3, 18.9.65.3.9.1-5, 18.9.65.3.11.1
func CheckTerminalServices(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows NT\Terminal Services`)
	if err != nil {
		return make([]bool, 10)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"fAllowUnsolicited", "fAllowToGetHelp", "DisablePasswordSaving", "fDisableCdm",
		"fPromptForPassword", "fEncryptRPCTraffic", "SecurityLayer", "UserAuthentication", "MinEncryptionLevel",
		"DeleteTempDirsOnExit"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1), uint64(1), uint64(1), uint64(1), uint64(2),
		uint64(1), uint64(3), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckAppPrivacy is a helper function that checks the registry to determine if the system is configured with the correct settings for app privacy.
//
// CIS Benchmark Audit list index: 18.9.5.1
func CheckAppPrivacy(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\AppPrivacy`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "LetAppsActivateWithVoiceAboveLock", 2)
}

// CheckAppx is a helper function that checks the registry to determine if the system is configured with the correct settings for Appx.
//
// CIS Benchmark Audit list index: 18.9.4.2
func CheckAppx(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Appx`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "BlockNonAdminUserInstall", 1)
}

// CheckCloudContent is a helper function that checks the registry to determine if the system is configured with the correct settings for cloud content.
//
// CIS Benchmark Audit list indices: 18.9.14.1, 18.9.14.3
func CheckCloudContent(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\CloudContent`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"DisableConsumerAccountStateContent", "DisableWindowsConsumerFeatures"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckConnect is a helper function that checks the registry to determine if the system is configured with the correct settings for Connect.
//
// CIS Benchmark Audit list index: 18.9.15.1
func CheckConnect(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Connect`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "RequirePinForPairing", []uint64{1, 2})
}

// CheckCredentialsDelegation is a helper function that checks the registry to determine if the system is configured with the correct settings for credentials delegation.
//
// CIS Benchmark Audit list index: 18.8.4.2
func CheckCredentialsDelegation(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\CredentialsDelegation`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "AllowProtectedCreds", 1)
}

// CheckCredui is a helper function that checks the registry to determine if the system is configured with the correct settings for Credui.
//
// CIS Benchmark Audit list index: 18.9.16.1
func CheckCredui(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Credui`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "DisablePasswordReveal", 1)
}

// CheckDataCollection is a helper function that checks the registry to determine if the system is configured with the correct settings for data collection.
//
// CIS Benchmark Audit list indices: 18.9.17.1, 18.9.17.3-7
func CheckDataCollection(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\DataCollection`)
	if err != nil {
		return make([]bool, 6)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AllowTelemetry", "DisableOneSettingsDownloads", "DoNotShowFeedbackNotifications",
		"EnableOneSettingsAuditing", "LimitDiagnosticLogCollection", "LimitDumpCollection"}

	expectedValues := []interface{}{[]uint64{0, 1}, uint64(1), uint64(1), uint64(1), uint64(1), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)

}

// CheckDeliveryOptimization is a helper function that checks the registry to determine if the system is configured with the correct settings for delivery optimization.
//
// CIS Benchmark Audit list index: 18.9.18.1
func CheckDeliveryOptimization(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\DeliveryOptimization`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "DODownloadMode", []uint64{0, 1, 2, 99, 100})
}

// CheckDeviceMetaData is a helper function that checks the registry to determine if the system is configured with the correct settings for device metadata.
//
// CIS Benchmark Audit list index: 18.8.7.2
func CheckDeviceMetaData(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Device Metadata`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "PreventDeviceMetadataFromNetwork", 1)
}

// CheckEventLog is a helper function that checks the registry to determine if the system is configured with the correct settings for the event log.
//
// CIS Benchmark Audit list indices: 18.9.27.1.1-2
func CheckEventLog(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\EventLog\Application`)
	if err != nil {
		results = append(results, make([]bool, 2)...)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkEventLogSecurity(registryKey)...)
	results = append(results, checkEventLogSetup(registryKey)...)
	results = append(results, checkEventLogSystem(registryKey)...)
	return results
}

// checkEventLogSecurity is a helper function that checks the registry to determine if the system is configured with the correct settings for the security event log.
//
// CIS Benchmark Audit list indices: 18.9.27.2.1-2
func checkEventLogSecurity(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\EventLog\Security`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{196608, ^uint64(0)}}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkEventLogSetup is a helper function that checks the registry to determine if the system is configured with the correct settings for the setup event log.
//
// CIS Benchmark Audit list indices: 18.9.27.3.1-2
func checkEventLogSetup(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Eventlog\Setup`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkEventLogSystem is a helper function that checks the registry to determine if the system is configured with the correct settings for the system event log.
//
// CIS Benchmark Audit list indices: 18.9.27.4.1-2
func checkEventLogSystem(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\EventLog\System`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckWindowsExplorer is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Explorer.
//
// CIS Benchmark Audit list indices: 18.9.8.1, 18.9.31.2-3,
func CheckWindowsExplorer(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Explorer`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"NoAutoplayfornonVolume", "NoDataExecutionPrevention", "NoHeapTerminationOnCorruption"}

	expectedValues := []interface{}{uint64(1), uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)

}

// CheckGameDVR is a helper function that checks the registry to determine if the system is configured with the correct settings for GameDVR.
//
// CIS Benchmark Audit list index: 18.9.87.1
func CheckGameDVR(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\GameDVR`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "AllowGameDVR", 0)
}

// CheckGroupPolicy is a helper function that checks the registry to determine if the system is configured with the correct settings for Group Policy.
//
// CIS Benchmark Audit list indices: 18.8.21.2-3
func CheckGroupPolicy(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Group Policy\{35378EAC-683F-11D2-A89A-00C04FBBCFA2}`)
	if err != nil {
		return make([]bool, 4)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"NoBackgroundPolicy", "NoGPOListChanges"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWcmSvcGroupPolicy(registryKey)...)
	return results
}

// checkWcmSvcGroupPolicy is a helper function that checks the registry to determine if the system is configured with the correct settings for WcmSvc group policy.
//
// CIS Benchmark Audit list indices: 18.5.21.1-2
func checkWcmSvcGroupPolicy(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\WcmSvc\GroupPolicy`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"fMinimizeConnections", "fBlockNonDomain"}

	expectedValues := []interface{}{uint64(3), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckHomeGroup is a helper function that checks the registry to determine if the system is configured with the correct settings for HomeGroup.
//
// CIS Benchmark Audit list index: 18.9.36.1
func CheckHomeGroup(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\HomeGroup`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "DisableHomeGroup", 1)
}

// CheckInstaller is a helper function that checks the registry to determine if the system is configured with the correct settings for the installer.
//
// CIS Benchmark Audit list indices: 18.9.90.1-2
func CheckInstaller(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Installer`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableUserControl", "AlwaysInstallElevated"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckLanman is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman.
//
// CIS Benchmark Audit list index: 18.5.8.1
func CheckLanman(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\LanmanWorkstation`)
	if err != nil {
		return make([]bool, 13)
	}
	defer mocking.CloseRegistryKey(key)

	results = append(results, checkIntegerValue(key, "AllowInsecureGuestAuth", 0))
	results = append(results, checkLanmanParameters(registryKey)...)
	results = append(results, checkLanmanServerParameters(registryKey)...)
	return results
}

// checkLanmanParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman parameters.
//
// CIS Benchmark Audit list indices: 2.3.8.1-3
func checkLanmanParameters(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SYSTEM\CurrentControlSet\Services\LanmanWorkstation\Parameters`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"RequireSecuritySignature", "EnableSecuritySignature", "EnablePlainTextPassword"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)

	return results
}

// checkLanmanServerParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman server parameters.
//
// CIS Benchmark Audit list indices: 2.3.9.1-5, 2.3.10.6, 2.3.10.9, 2.3.10.11, 18.3.3
func checkLanmanServerParameters(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SYSTEM\CurrentControlSet\Services\LanmanServer\Parameters`)
	if err != nil {
		return make([]bool, 9)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AutoDisconnect", "RequireSecuritySignature", "EnableSecuritySignature", "enableforcedlogoff",
		"SMBServerNameHardeningLevel", "NullSessionPipes", "RestrictNullSessAccess", "NullSessionShares", "SMB1"}

	expectedValues := []interface{}{[]uint64{1, 15}, uint64(1), uint64(1), uint64(1), []uint64{1, 2}, nil, uint64(1),
		nil, uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckNetworkConnections is a helper function that checks the registry to determine if the system is configured with the correct settings for network connections.
//
// CIS Benchmark Audit list indices: 18.5.11.2-4
func CheckNetworkConnections(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Network Connections`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"NC_AllowNetBridge_NLA", "NC_ShowSharedAccessUI", "NC_StdDomainUserSetLocation"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckNetworkProvider is a helper function that checks the registry to determine if the system is configured with the correct settings for the network provider.
//
// CIS Benchmark Audit list index: 18.5.14.1
// TODO: NEEDS CHECKING, IF THIS WORKS AS INTENDED, COULD NOT TEST DUE TO NON-EXISTENT REGISTRY KEY
func CheckNetworkProvider(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\NetworkProvider\HardenedPaths`)
	if err != nil {
		return make([]bool, 1)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"\\\\*\\NETLOGON", "\\\\*\\SYSVOL"}

	expectedValues := []string{
		"[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1.*[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1",
		"[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1.*[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1",
	}

	return checkMultipleStringValues(key, settings, expectedValues)
}

// CheckOneDrive is a helper function that checks the registry to determine if the system is configured with the correct settings for OneDrive.
//
// CIS Benchmark Audit list index: 18.9.58.1
func CheckOneDrive(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\OneDrive`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "DisableFileSyncNGSC", 1)
}

// CheckPersonalization is a helper function that checks the registry to determine if the system is configured with the correct settings for personalization.
//
// CIS Benchmark Audit list indices: 18.1.1.1-2
func CheckPersonalization(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Personalization`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"NoLockScreenCamera", "NoLockScreenSlideshow"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckPowerShell is a helper function that checks the registry to determine if the system is configured with the correct settings for PowerShell.
//
// CIS Benchmark Audit list indices: 18.9.100.1-2
func CheckPowerShell(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Powershell\Scriptblocklogging`)
	if err != nil {
		results = append(results, false)
	} else {
		results = append(results, checkIntegerValue(key, "EnableScriptBlockLogging", 1))
	}
	defer mocking.CloseRegistryKey(key)

	key, err = openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\PowerShell\Transcription`)
	if err != nil {
		results = append(results, false)
	} else {
		results = append(results, checkIntegerValue(key, "EnableTranscripting", 0))

	}
	defer mocking.CloseRegistryKey(key)

	return results
}

// CheckPreviewBuild is a helper function that checks the registry to determine if the system is configured with the correct settings for preview builds.
//
// CIS Benchmark Audit list index: 18.9.17.8
func CheckPreviewBuild(registryKey mocking.RegistryKey) bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Previewbuilds`)
	if err != nil {
		return false
	}
	defer mocking.CloseRegistryKey(key)

	return checkIntegerValue(key, "AllowBuildPreview", 0)
}

// CheckSandbox is a helper function that checks the registry to determine if the system is configured with the correct settings for the sandbox.
//
// CIS Benchmark Audit list indices: 18.9.104.1-2
func CheckSandbox(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Sandbox`)
	if err != nil {
		return make([]bool, 2)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AllowClipboardRedirection", "AllowNetworking"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)

}

// CheckWindowsSystem is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows System.
//
// CIS Benchmark Audit list indices: 18.8.21.4, 18.8.28.1-7, 18.9.16.3, 18.9.85.1.1
func CheckWindowsSystem(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\System`)
	if err != nil {
		return make([]bool, 11)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableCdp", "BlockUserFromShowingAccountDetailsOnSignin", "DontDisplayNetworkSelectionUI",
		"DontEnumerateConnectedUsers", "EnumerateLocalUsers", "DisableLockScreenAppNotifications",
		"BlockDomainPicturePassword", "AllowDomainPINLogon", "NoLocalPasswordResetQuestions", " EnableSmartScreen",
		"ShellSmartScreenLevel"}

	expectedValues := []interface{}{uint64(0), uint64(1), uint64(1), uint64(1), uint64(0), uint64(1), uint64(1),
		uint64(0), uint64(1), uint64(1), uint64(1)}
	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckWindowsSearch is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Search.
//
// CIS Benchmark Audit list indices: 18.9.67.3-6
func CheckWindowsSearch(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Windows Search`)
	if err != nil {
		return make([]bool, 4)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AllowCortana", "AllowCortanaAboveLock", "AllowIndexingEncryptedStoresOrItems", "AllowSearchToUseLocation"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckWindowsUpdate is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Update.
//
// CIS Benchmark Audit list indices: 18.9.108.2.2-3, 18.9.108.4.1-3
func CheckWindowsUpdate(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`)
	if err != nil {
		return make([]bool, 9)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"SetDisablePauseUXAccess", "ManagePreviewBuildsPolicyValue", "DeferFeatureUpdates",
		"DeferFeatureUpdatesPeriodInDays", "DeferQualityUpdates", "DeferQualityUpdatesPeriodInDays"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(1), []uint64{180, ^uint64(0)}, uint64(1), uint64(0)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWindowsUpdateAu(registryKey)...)
	return results
}

// checkWindowsUpdateAu is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Update AU.
//
// CIS Benchmark Audit list indices: 18.9.108.1.1, 18.9.108.2.1-2
func checkWindowsUpdateAu(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\Windowsupdate\Au`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"NoAutoRebootWithLoggedOnUsers", "NoAutoUpdate", "ScheduledInstallDay"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// CheckWinRM is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM.
func CheckWinRM(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)

	results = append(results, checkWinRMClient(registryKey)...)
	results = append(results, checkWinRMService(registryKey)...)

	return results
}

// checkWinRMClient is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM client.
//
// CIS Benchmark Audit list indices: 18.9.102.1.1-3
func checkWinRMClient(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\WinRM\Client`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AllowBasic", "AllowUnencryptedTraffic", "AllowDigest"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

// checkWinRMService is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM service.
//
// CIS Benchmark Audit list indices: 18.9.102.2.1, 18.9.102.2.3-4
func checkWinRMService(registryKey mocking.RegistryKey) []bool {
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\Windows\WinRM\Service`)
	if err != nil {
		return make([]bool, 3)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"AllowBasic", "AllowUnencryptedTraffic", "DisableRunAs"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1)}

	return checkMultipleIntegerValues(key, settings, expectedValues)
}

func CheckWindowsFirewall(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)

	results = append(results, checkWindowsFirewallPrivateProfile(registryKey)...)
	results = append(results, checkWindowsFirewallPublicProfile(registryKey)...)
	results = append(results, checkWindowsFirewallDomainProfile(registryKey)...)
	return results
}

// checkWindowsFirewallDomainProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the domain profile.
//
// CIS Benchmark Audit list indices: 9.1.1-4
func checkWindowsFirewallDomainProfile(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\DomainProfile`)
	if err != nil {
		results = append(results, make([]bool, 4)...)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWindowsFirewallDomainProfileLogging(registryKey)...)

	return results
}

// checkWindowsFirewallDomainProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the domain profile logging.
//
// CIS Benchmark Audit list indices: 9.1.5-8
func checkWindowsFirewallDomainProfileLogging(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\DomainProfile\Logging`)
	if err != nil {
		return make([]bool, 4)
	}
	defer mocking.CloseRegistryKey(key)

	stringSetting := "LogFilePath"
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := `%SYSTEMROOT%\System32\logfiles\firewall\domainfw.log`
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkStringValue(key, stringSetting, expectedString))
	return results
}

// checkWindowsFirewallPublicProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile.
//
// CIS Benchmark Audit list indices: 9.2.1-4
func checkWindowsFirewallPrivateProfile(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\PrivateProfile`)
	if err != nil {
		results = append(results, make([]bool, 4)...)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWindowsFirewallPrivateProfileLogging(registryKey)...)

	return results
}

// checkWindowsFirewallPrivateProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the private profile logging.
//
// CIS Benchmark Audit list indices: 9.2.5-8
func checkWindowsFirewallPrivateProfileLogging(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\PrivateProfile\Logging`)
	if err != nil {
		return make([]bool, 4)
	}
	defer mocking.CloseRegistryKey(key)

	stringSetting := "LogFilePath"
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := `%SYSTEMROOT%\System32\logfiles\firewall\privatefw.log`
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkStringValue(key, stringSetting, expectedString))
	return results
}

// checkWindowsFirewallPublicProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile.
//
// CIS Benchmark Audit list indices: 9.3.1-6
func checkWindowsFirewallPublicProfile(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\PublicProfile`)
	if err != nil {
		results = append(results, make([]bool, 6)...)
	}
	defer mocking.CloseRegistryKey(key)

	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications",
		"AllowLocalPolicyMerge", "AllowLocalIPsecPolicyMerge"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1), uint64(0), uint64(0)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkWindowsFirewallPublicProfileLogging(registryKey)...)

	return results

}

// checkWindowsFirewallPublicProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile logging.
//
// CIS Benchmark Audit list indices: 9.3.7-10
func checkWindowsFirewallPublicProfileLogging(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	key, err := openRegistryKeyWithErrHandling(registryKey, `SOFTWARE\Policies\Microsoft\WindowsFirewall\PublicProfile\Logging`)
	if err != nil {
		return make([]bool, 4)
	}
	defer mocking.CloseRegistryKey(key)

	stringSetting := "LogFilePath"
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := `%SYSTEMROOT%\System32\logfiles\firewall\publicfw.log`
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	results = append(results, checkMultipleIntegerValues(key, settings, expectedValues)...)
	results = append(results, checkStringValue(key, stringSetting, expectedString))
	return results
}

func checkIntegerValue(openKey mocking.RegistryKey, value string, expected interface{}) bool {
	val, _, err := openKey.GetIntegerValue(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading registry value of "+value, err)
		return false
	}
	// Determine functionality based on the value type of the expected parameter
	switch v := expected.(type) {
	// Single uint64, check if registry value is equal to expected value
	case uint64:
		return val == v
	// Slice of uint64 values, check if registry value is in the slice
	case []uint64:
		for _, i := range v {
			if val == i {
				return true
			}
		}
	// Slice of exactly 2 uint64 values, check if registry value is within the range
	case [2]uint64:
		return val >= v[0] && val <= v[1]
	default:
		return false
	}
	return false
}

func checkStringValue(openKey mocking.RegistryKey, value string, expected string) bool {
	val, _, err := openKey.GetStringValue(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading registry value of "+value, err)
		return false
	}
	return val == expected
}

func openRegistryKeyWithErrHandling(registryKey mocking.RegistryKey, path string) (mocking.RegistryKey, error) {
	key, err := mocking.OpenRegistryKey(registryKey, path)
	if err != nil {
		logger.Log.ErrorWithErr("Error opening registry key for CIS Audit list", err)
	}
	return key, err
}

func checkMultipleIntegerValues(openKey mocking.RegistryKey, settings []string, expectedValues []interface{}) []bool {
	results := make([]bool, len(settings))
	for i, val := range settings {
		results[i] = checkIntegerValue(openKey, val, expectedValues[i])
	}
	return results
}

func checkMultipleStringValues(openKey mocking.RegistryKey, settings []string, expectedValues []string) []bool {
	results := make([]bool, len(settings))
	for i, val := range settings {
		results[i] = checkStringValue(openKey, val, expectedValues[i])
	}
	return results
}
