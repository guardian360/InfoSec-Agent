package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/mocking"

// CheckPoliciesHKLM is a function that checks various registry settings related to different policies
// to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_LOCAL_MACHINE for this function.
//
// Returns:
//
//   - []bool: A slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
func CheckPoliciesHKLM(registryKey mocking.RegistryKey) {
	for _, check := range policyChecksHKLM {
		check(registryKey)
	}
}

// policyChecksHKLM is a collection of registry checks related to different policies.
// Each function in the collection represents a different policy check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var policyChecksHKLM = []func(mocking.RegistryKey){
	policiesCredui,
	policiesExplorerHKLM,
	policiesSystem,
	policiesAdmPwd,
	policiesFacialFeatures,
	policiesDsh,
	policiesInputPersonalization,
	policiesIEFeeds,
	policiesMicrosoftAccount,
	policiesPhishingFilter,
	policiesPowerSettings,
	policiesWindowsDefender,
	policiesDNSClient,
	policiesPrinters,
	policiesRPC,
	policiesTerminalServices,
	policiesAppPrivacy,
	policiesAppx,
	policiesCloudContentHKLM,
	policiesConnect,
	policiesCredentialsDelegation,
	policiesGeneralCredui,
	policiesPreviewBuild,
	policiesSandbox,
	policiesDataCollection,
	policiesDeliveryOptimization,
	policiesDeviceMetadata,
	policiesEventLog,
	policiesWindowsExplorer,
	policiesGameDVR,
	policiesGroupPolicy,
	policiesHomeGroup,
	policiesInstallerHKLM,
	policiesLanman,
	policiesNetworkConnections,
	policiesNetworkProvider,
	policiesOneDrive,
	policiesPersonalization,
	policiesPowerShell,
	policiesGeneralSystem,
	policiesWindowsSearch,
	policiesWindowsUpdate,
	policiesWinRM,
	policiesWindowsFirewall,
	policiesWindowsInkWorkspace,
	policiesWindowsStore,
	policiesEarlyLaunch,
}

// policiesCredui is a helper function that checks the registry to determine if the system is configured to enumerate administrator accounts.
//
// CIS Benchmark Audit list index: 18.9.16.2
func policiesCredui(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\Currentversion\Policies\Credui`

	settings := []string{"EnumerateAdministrators"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesExplorerHKLM is a helper function that checks the registry to determine if the system is configured with the correct settings for Explorer policies.
//
// CIS Benchmark Audit list indices: 18.8.22.1.6, 18.9.8.2, 18.9.8.3, 18.9.31.4
func policiesExplorerHKLM(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`

	settings := []string{
		"NoWebServices",
		"NoAutorun",
		"NoDriveTypeAutoRun",
		"PreXPSP2ShellProtocolBehavior",
	}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(255), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesSystem is a helper function that checks the registry to determine if the system is configured with the correct settings for system policies.
//
// CIS Benchmark Audit list indices: 2.3.1.2, 2.3.7.1-3, 2.3.11.4, 2.3.17.1-8, 18.3.1, 18.8.3.1, 18.8.4.1, 18.9.6.1, 18.9.91.1
func policiesSystem(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`

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

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

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
			registryPath = subKey
			CheckIntegerRegistrySettings(
				registryKey, registryPath, []string{subKeysSettings[i]}, []interface{}{subKeysExpected[i]})
		}()
	}
}

// policiesAdmPwd is a helper function that checks the registry to determine if the system is configured with the correct settings for the administrator password.
//
// CIS Benchmark Audit list indices: 18.2.2-6
func policiesAdmPwd(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft Services\AdmPwd`

	settings := []string{
		"PwdExpirationProtectionEnabled",
		"AdmPwdEnabled",
		"PasswordComplexity",
		"PasswordLength",
		"PasswordAgeDays",
	}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(4), []uint64{15, ^uint64(0)}, []uint64{0, 30}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesFacialFeatures is a helper function that checks the registry to determine if the system is configured with the correct settings for enhanced anti-spoofing.
//
// CIS Benchmark Audit list index: 18.9.10.1.1
func policiesFacialFeatures(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Biometrics\FacialFeatures`

	settings := []string{"EnhancedAntiSpoofing"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesDsh is a helper function that checks the registry to determine if the system is configured to allow widgets.
//
// CIS Benchmark Audit list index: 18.9.81.1
func policiesDsh(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Dsh`

	settings := []string{"AllowNewsAndInterests"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesInputPersonalization is a helper function that checks the registry to determine if the system is configured to allow online speech recognition services.
//
// CIS Benchmark Audit list index: 18.1.2.2
func policiesInputPersonalization(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\InputPersonalization`

	settings := []string{"AllowInputPersonalization"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesIEFeeds is a helper function that checks the registry to determine if the system is configured to download enclosures.
//
// CIS Benchmark Audit list index: 18.9.66.1
func policiesIEFeeds(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Internet Explorer\Feeds`

	settings := []string{"DisableEnclosureDownload"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesMicrosoftAccount is a helper function that checks the registry to determine if the system is configured to block consumer user authentication.
//
// CIS Benchmark Audit list index: 18.9.46.1
func policiesMicrosoftAccount(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\MicrosoftAccount`

	settings := []string{"DisableUserAuth"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPhishingFilter is a helper function that checks the registry to determine if the system is configured with the correct settings for the phishing filter.
//
// CIS Benchmark Audit list indices: 18.9.85.2.1-2
func policiesPhishingFilter(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\MicrosoftEdge\PhishingFilter`

	settings := []string{"EnabledV9", "PreventOverride"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPowerSettings is a helper function that checks the registry to determine if the system is configured with the correct power settings.
//
// CIS Benchmark Audit list indices: 18.8.34.6.1-2, 18.8.34.6.5-6
func policiesPowerSettings(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Power\PowerSettings\f15576e8-98b7-4186-b944-eafa664402d9`

	settings := []string{
		"DCSettingIndex",
		"ACSettingIndex",
	}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

	registryPath = `SOFTWARE\Policies\Microsoft\Power\PowerSettings\0e796bdb-100d-47d6-a2d5-f7d2daa51f51`

	expectedValues = []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsDefender is a helper function that checks the registry to determine if Windows Defender is configured with the correct settings.
//
// CIS Benchmark Audit list indices: 18.9.47.15-16
func policiesWindowsDefender(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender`

	settings := []string{
		"PUAProtection",
		"DisableAntiSpyware",
	}

	expectedValues := []interface{}{uint64(1), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkWindowsDefenderScan(registryKey)
	checkWindowsDefenderRealTime(registryKey)
	checkWindowsDefenderASR(registryKey)
	checkWindowsDefenderSpyNet(registryKey)
	checkWindowsDefenderNetworkProtection(registryKey)
	checkWindowsDefenderAppBrowserProtection(registryKey)

}

// checkWindowsDefenderScan is a helper function that checks the registry to determine if Windows Defender is configured with the correct scan settings.
//
// CIS Benchmark Audit list indices: 18.9.47.12.1-2
func checkWindowsDefenderScan(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Scan`

	settings := []string{"DisableRemovableDriveScanning", "DisableEmailScanning"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsDefenderRealTime is a helper function that checks the registry to determine if Windows Defender is configured with the correct real-time protection settings.
//
// CIS Benchmark Audit list indices: 18.9.47.9.1-4
func checkWindowsDefenderRealTime(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection`

	settings := []string{"DisableIOAVProtection", "DisableRealtimeMonitoring", "DisableBehaviorMonitoring", "DisableScriptScanning"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsDefenderASR is a helper function that checks the registry to determine if Windows Defender is configured with the correct ASR settings.
//
// CIS Benchmark Audit list indices: 18.9.47.5.1.1-2
func checkWindowsDefenderASR(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\ASR`

	settings := []string{"ExploitGuard_ASR_Rules"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

	registryPath = `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\ASR\Rules`

	settings = []string{"26190899-1602-49e8-8b27-eb1d0a1ce869",
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

	expectedValues = []interface{}{uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(1),
		uint64(1), uint64(1), uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsDefenderSpyNet is a helper function that checks the registry to determine if Windows Defender is configured with the correct SpyNet settings.
//
// CIS Benchmark Audit list index: 18.9.47.4.1
func checkWindowsDefenderSpyNet(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Spynet`

	settings := []string{"LocalSettingOverrideSpynetReporting"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsDefenderAppBrowserProtection is a helper function that checks the registry to determine if Windows Defender is configured with the correct app and browser protection settings.
//
// CIS Benchmark Audit list index: 18.9.105.2.1
func checkWindowsDefenderAppBrowserProtection(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender Security Center\App and Browser protection`

	settings := []string{"DisallowExploitProtectionOverride"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsDefenderNetworkProtection is a helper function that checks the registry to determine if Windows Defender is configured with the correct network protection settings.
//
// CIS Benchmark Audit list index: 18.9.47.5.3.1
func checkWindowsDefenderNetworkProtection(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\Network Protection`

	settings := []string{"EnableNetworkProtection"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesDNSClient is a helper function that checks the registry to determine if the system is configured with the correct settings for the DNS client.
//
// CIS Benchmark Audit list indices: 18.5.4.2
func policiesDNSClient(registryKey mocking.RegistryKey) {
	registryPath := DNSClientRegistryPath

	settings := []string{"EnableMulticast"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPrinters is a helper function that checks the registry to determine if the system is configured with the correct settings for printers.
//
// CIS Benchmark Audit list indices: 18.3.5, 18.6.1-3, 18.8.22.1.2
func policiesPrinters(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows NT\Printers`

	settings := []string{"RegisterSpoolerRemoteRpcEndPoint", "DisableWebPnPDownload"}

	expectedValues := []interface{}{uint64(2), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

	registryPath = `SOFTWARE\Policies\Microsoft\Windows NT\Printers\PointAndPrint`

	settings = []string{"RestrictDriverInstallationToAdministrators", "NoWarningNoElevationOnInstall", "UpdatePromptSettings"}

	expectedValues = []interface{}{uint64(1), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesRPC is a helper function that checks the registry to determine if the system is configured with the correct settings for RPC.
//
// CIS Benchmark Audit list indices: 18.8.37.1-2
func policiesRPC(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows NT\Rpc`

	settings := []string{"EnableAuthEpResolution", "RestrictRemoteClients"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesTerminalServices is a helper function that checks the registry to determine if the system is configured with the correct settings for terminal services.
//
// CIS Benchmark Audit list indices: 18.8.36.1-2, 18.9.65.2.2, 18.9.65.3.3.3, 18.9.65.3.9.1-5, 18.9.65.3.11.1
func policiesTerminalServices(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows NT\Terminal Services`

	settings := []string{"fAllowUnsolicited", "fAllowToGetHelp", "DisablePasswordSaving", "fDisableCdm",
		"fPromptForPassword", "fEncryptRPCTraffic", "SecurityLayer", "UserAuthentication", "MinEncryptionLevel",
		"DeleteTempDirsOnExit"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1), uint64(1), uint64(1), uint64(1), uint64(2),
		uint64(1), uint64(3), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesAppPrivacy is a helper function that checks the registry to determine if the system is configured with the correct settings for app privacy.
//
// CIS Benchmark Audit list index: 18.9.5.1
func policiesAppPrivacy(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\AppPrivacy`

	settings := []string{"LetAppsActivateWithVoiceAboveLock"}

	expectedValues := []interface{}{uint64(2)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesAppx is a helper function that checks the registry to determine if the system is configured with the correct settings for Appx.
//
// CIS Benchmark Audit list index: 18.9.4.2
func policiesAppx(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Appx`

	settings := []string{"BlockNonAdminUserInstall"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesCloudContentHKLM is a helper function that checks the registry to determine if the system is configured with the correct settings for cloud content.
//
// CIS Benchmark Audit list indices: 18.9.14.1, 18.9.14.3
func policiesCloudContentHKLM(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CloudContent`

	settings := []string{"DisableConsumerAccountStateContent", "DisableWindowsConsumerFeatures"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesConnect is a helper function that checks the registry to determine if the system is configured with the correct settings for Connect.
//
// CIS Benchmark Audit list index: 18.9.15.1
func policiesConnect(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Connect`

	settings := []string{"RequirePinForPairing"}

	expectedValues := []interface{}{[]uint64{1, 2}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesCredentialsDelegation is a helper function that checks the registry to determine if the system is configured with the correct settings for credentials delegation.
//
// CIS Benchmark Audit list index: 18.8.4.2
func policiesCredentialsDelegation(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\CredentialsDelegation`

	settings := []string{"AllowProtectedCreds"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesGeneralCredui is a helper function that checks the registry to determine if the system is configured with the correct settings for Credui.
//
// CIS Benchmark Audit list index: 18.9.16.1
func policiesGeneralCredui(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Credui`

	settings := []string{"DisablePasswordReveal"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesDataCollection is a helper function that checks the registry to determine if the system is configured with the correct settings for data collection.
//
// CIS Benchmark Audit list indices: 18.9.17.1, 18.9.17.3-7
func policiesDataCollection(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\DataCollection`

	settings := []string{"AllowTelemetry", "DisableOneSettingsDownloads", "DoNotShowFeedbackNotifications",
		"EnableOneSettingsAuditing", "LimitDiagnosticLogCollection", "LimitDumpCollection"}

	expectedValues := []interface{}{[]uint64{0, 1}, uint64(1), uint64(1), uint64(1), uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesDeliveryOptimization is a helper function that checks the registry to determine if the system is configured with the correct settings for delivery optimization.
//
// CIS Benchmark Audit list index: 18.9.18.1
func policiesDeliveryOptimization(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\DeliveryOptimization`

	settings := []string{"DODownloadMode"}

	expectedValues := []interface{}{[]uint64{0, 1, 2, 99, 100}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesDeviceMetadata is a helper function that checks the registry to determine if the system is configured with the correct settings for device metadata.
//
// CIS Benchmark Audit list index: 18.8.7.2
func policiesDeviceMetadata(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Device Metadata`

	settings := []string{"PreventDeviceMetadataFromNetwork"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesEventLog is a helper function that checks the registry to determine if the system is configured with the correct settings for the event log.
//
// CIS Benchmark Audit list indices: 18.9.27.1.1-2
func policiesEventLog(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\EventLog\Application`

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkEventLogSecurity(registryKey)
	checkEventLogSetup(registryKey)
	checkEventLogSystem(registryKey)
}

// checkEventLogSecurity is a helper function that checks the registry to determine if the system is configured with the correct settings for the security event log.
//
// CIS Benchmark Audit list indices: 18.9.27.2.1-2
func checkEventLogSecurity(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\EventLog\Security`

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{196608, ^uint64(0)}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkEventLogSetup is a helper function that checks the registry to determine if the system is configured with the correct settings for the setup event log.
//
// CIS Benchmark Audit list indices: 18.9.27.3.1-2
func checkEventLogSetup(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Eventlog\Setup`

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkEventLogSystem is a helper function that checks the registry to determine if the system is configured with the correct settings for the system event log.
//
// CIS Benchmark Audit list indices: 18.9.27.4.1-2
func checkEventLogSystem(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\EventLog\System`

	settings := []string{"Retention", "MaxSize"}

	expectedValues := []interface{}{uint64(0), []uint64{32768, ^uint64(0)}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsExplorer is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Explorer.
//
// CIS Benchmark Audit list indices: 18.9.8.1, 18.9.31.2-3,
func policiesWindowsExplorer(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Explorer`

	settings := []string{"NoAutoplayfornonVolume", "NoDataExecutionPrevention", "NoHeapTerminationOnCorruption"}

	expectedValues := []interface{}{uint64(1), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesGameDVR is a helper function that checks the registry to determine if the system is configured with the correct settings for GameDVR.
//
// CIS Benchmark Audit list index: 18.9.87.1
func policiesGameDVR(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\GameDVR`

	settings := []string{"AllowGameDVR"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesGroupPolicy is a helper function that checks the registry to determine if the system is configured with the correct settings for Group Policy.
//
// CIS Benchmark Audit list indices: 18.8.21.2-3
func policiesGroupPolicy(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Group Policy\{35378EAC-683F-11D2-A89A-00C04FBBCFA2}`

	settings := []string{"NoBackgroundPolicy", "NoGPOListChanges"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkWcmSvcGroupPolicy(registryKey)
}

// checkWcmSvcGroupPolicy is a helper function that checks the registry to determine if the system is configured with the correct settings for WcmSvc group policy.
//
// CIS Benchmark Audit list indices: 18.5.21.1-2
func checkWcmSvcGroupPolicy(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\WcmSvc\GroupPolicy`

	settings := []string{"fMinimizeConnections", "fBlockNonDomain"}

	expectedValues := []interface{}{uint64(3), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesHomeGroup is a helper function that checks the registry to determine if the system is configured with the correct settings for HomeGroup.
//
// CIS Benchmark Audit list index: 18.9.36.1
func policiesHomeGroup(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\HomeGroup`

	settings := []string{"DisableHomeGroup"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesInstallerHKLM is a helper function that checks the registry to determine if the system is configured with the correct settings for the installer.
//
// CIS Benchmark Audit list indices: 18.9.90.1-2
func policiesInstallerHKLM(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Installer`

	settings := []string{"EnableUserControl", "AlwaysInstallElevated"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesLanman is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman.
//
// CIS Benchmark Audit list index: 18.5.8.1
func policiesLanman(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\LanmanWorkstation`

	settings := []string{"AllowInsecureGuestAuth"}

	expectedValues := []interface{}{uint64(0)}
	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkLanmanParameters(registryKey)
	checkLanmanServerParameters(registryKey)
}

// checkLanmanParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman parameters.
//
// CIS Benchmark Audit list indices: 2.3.8.1-3
func checkLanmanParameters(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Services\LanmanWorkstation\Parameters`

	settings := []string{"RequireSecuritySignature", "EnableSecuritySignature", "EnablePlainTextPassword"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkLanmanServerParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for Lanman server parameters.
//
// CIS Benchmark Audit list indices: 2.3.9.1-5, 2.3.10.6, 2.3.10.9, 2.3.10.11, 18.3.3
func checkLanmanServerParameters(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Services\LanmanServer\Parameters`

	settings := []string{"AutoDisconnect", "RequireSecuritySignature", "EnableSecuritySignature", "enableforcedlogoff",
		"SMBServerNameHardeningLevel", "NullSessionPipes", "RestrictNullSessAccess", "NullSessionShares", "SMB1"}

	expectedValues := []interface{}{[]uint64{1, 15}, uint64(1), uint64(1), uint64(1), []uint64{1, 2}, nil, uint64(1),
		nil, uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesNetworkConnections is a helper function that checks the registry to determine if the system is configured with the correct settings for network connections.
//
// CIS Benchmark Audit list indices: 18.5.11.2-4
func policiesNetworkConnections(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Network Connections`

	settings := []string{"NC_AllowNetBridge_NLA", "NC_ShowSharedAccessUI", "NC_StdDomainUserSetLocation"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesNetworkProvider is a helper function that checks the registry to determine if the system is configured with the correct settings for the network provider.
//
// CIS Benchmark Audit list index: 18.5.14.1
// TODO: NEEDS CHECKING, IF THIS WORKS AS INTENDED, COULD NOT TEST DUE TO NON-EXISTENT REGISTRY KEY
func policiesNetworkProvider(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\NetworkProvider\HardenedPaths`

	settings := []string{"\\\\*\\NETLOGON", "\\\\*\\SYSVOL"}

	expectedValues := []string{
		"[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1.*[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1",
		"[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1.*[Rr]equire([Mm]utual[Aa]uthentication|[Ii]ntegrity)=1",
	}

	CheckStringRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesOneDrive is a helper function that checks the registry to determine if the system is configured with the correct settings for OneDrive.
//
// CIS Benchmark Audit list index: 18.9.58.1
func policiesOneDrive(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\OneDrive`

	settings := []string{"DisableFileSyncNGSC"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPersonalization is a helper function that checks the registry to determine if the system is configured with the correct settings for personalization.
//
// CIS Benchmark Audit list indices: 18.1.1.1-2
func policiesPersonalization(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Personalization`

	settings := []string{"NoLockScreenCamera", "NoLockScreenSlideshow"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPowerShell is a helper function that checks the registry to determine if the system is configured with the correct settings for PowerShell.
//
// CIS Benchmark Audit list indices: 18.9.100.1-2
func policiesPowerShell(registryKey mocking.RegistryKey) {
	checkPowershellScriptblocklogging(registryKey)
	checkPowershellTranscription(registryKey)
}

// checkPowershellScriptblocklogging is a helper function that checks the registry to determine if the system is configured with the correct settings for PowerShell script block logging.
//
// CIS Benchmark Audit list index: 18.9.100.1
func checkPowershellScriptblocklogging(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Powershell\Scriptblocklogging`

	settings := []string{"EnableScriptBlockLogging"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkPowershellTranscription is a helper function that checks the registry to determine if the system is configured with the correct settings for PowerShell transcription.
//
// CIS Benchmark Audit list index: 18.9.100.1
func checkPowershellTranscription(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Powershell\Transcription`

	settings := []string{"EnableTranscripting"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesPreviewBuild is a helper function that checks the registry to determine if the system is configured with the correct settings for preview builds.
//
// CIS Benchmark Audit list index: 18.9.17.8
func policiesPreviewBuild(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Previewbuilds`

	settings := []string{"AllowBuildPreview"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesSandbox is a helper function that checks the registry to determine if the system is configured with the correct settings for the sandbox.
//
// CIS Benchmark Audit list indices: 18.9.104.1-2
func policiesSandbox(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Sandbox`

	settings := []string{"AllowClipboardRedirection", "AllowNetworking"}

	expectedValues := []interface{}{uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesGeneralSystem is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows System.
//
// CIS Benchmark Audit list indices: 18.8.21.4, 18.8.28.1-7, 18.9.16.3, 18.9.85.1.1
func policiesGeneralSystem(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\System`

	settings := []string{"EnableCdp", "BlockUserFromShowingAccountDetailsOnSignin", "DontDisplayNetworkSelectionUI",
		"DontEnumerateConnectedUsers", "EnumerateLocalUsers", "DisableLockScreenAppNotifications",
		"BlockDomainPicturePassword", "AllowDomainPINLogon", "NoLocalPasswordResetQuestions", " EnableSmartScreen",
		"ShellSmartScreenLevel"}

	expectedValues := []interface{}{uint64(0), uint64(1), uint64(1), uint64(1), uint64(0), uint64(1), uint64(1),
		uint64(0), uint64(1), uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsSearch is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Search.
//
// CIS Benchmark Audit list indices: 18.9.67.3-6
func policiesWindowsSearch(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Windows Search`

	settings := []string{"AllowCortana", "AllowCortanaAboveLock", "AllowIndexingEncryptedStoresOrItems", "AllowSearchToUseLocation"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsUpdate is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Update.
//
// CIS Benchmark Audit list indices: 18.9.108.2.2-3, 18.9.108.4.1-3
func policiesWindowsUpdate(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`

	settings := []string{"SetDisablePauseUXAccess", "ManagePreviewBuildsPolicyValue", "DeferFeatureUpdates",
		"DeferFeatureUpdatesPeriodInDays", "DeferQualityUpdates", "DeferQualityUpdatesPeriodInDays"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(1), []uint64{180, ^uint64(0)}, uint64(1), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkWindowsUpdateAu(registryKey)
}

// checkWindowsUpdateAu is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Update AU.
//
// CIS Benchmark Audit list indices: 18.9.108.1.1, 18.9.108.2.1-2
func checkWindowsUpdateAu(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\Windowsupdate\Au`

	settings := []string{"NoAutoRebootWithLoggedOnUsers", "NoAutoUpdate", "ScheduledInstallDay"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWinRM is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM.
func policiesWinRM(registryKey mocking.RegistryKey) {

	checkWinRMClient(registryKey)
	checkWinRMService(registryKey)

}

// checkWinRMClient is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM client.
//
// CIS Benchmark Audit list indices: 18.9.102.1.1-3
func checkWinRMClient(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\WinRM\Client`

	settings := []string{"AllowBasic", "AllowUnencryptedTraffic", "AllowDigest"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWinRMService is a helper function that checks the registry to determine if the system is configured with the correct settings for WinRM service.
//
// CIS Benchmark Audit list indices: 18.9.102.2.1, 18.9.102.2.3-4
func checkWinRMService(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\WinRM\Service`

	settings := []string{"AllowBasic", "AllowUnencryptedTraffic", "DisableRunAs"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsFirewall is a helper function that checks the registry to determine if the system is configured with the correct settings for the Windows firewall.
//
// CIS Benchmark Audit list indices: 9.1.1-4, 9.1.5-8, 9.2.1-4, 9.2.5-8, 9.3.1-6, 9.3.7-10
func policiesWindowsFirewall(registryKey mocking.RegistryKey) {
	checkWindowsFirewallPrivateProfile(registryKey)
	checkWindowsFirewallPublicProfile(registryKey)
	checkWindowsFirewallDomainProfile(registryKey)
}

// checkWindowsFirewallDomainProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the domain profile.
//
// CIS Benchmark Audit list indices: 9.1.1-4
func checkWindowsFirewallDomainProfile(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\DomainProfile`

	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1)}

	checkWindowsFirewallDomainProfileLogging(registryKey)
	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsFirewallDomainProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the domain profile logging.
//
// CIS Benchmark Audit list indices: 9.1.5-8
func checkWindowsFirewallDomainProfileLogging(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\DomainProfile\Logging`

	stringSetting := []string{"LogFilePath"}
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := []string{`%SYSTEMROOT%\System32\logfiles\firewall\domainfw.log`}
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	CheckIntegerStringRegistrySettings(registryKey, registryPath, settings, expectedValues,
		stringSetting, expectedString)
}

// checkWindowsFirewallPublicProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile.
//
// CIS Benchmark Audit list indices: 9.2.1-4
func checkWindowsFirewallPrivateProfile(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\PrivateProfile`

	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1)}

	checkWindowsFirewallPrivateProfileLogging(registryKey)
	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkWindowsFirewallPrivateProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the private profile logging.
//
// CIS Benchmark Audit list indices: 9.2.5-8
func checkWindowsFirewallPrivateProfileLogging(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\PrivateProfile\Logging`

	stringSetting := []string{"LogFilePath"}
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := []string{`%SYSTEMROOT%\System32\logfiles\firewall\privatefw.log`}
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	CheckIntegerStringRegistrySettings(registryKey, registryPath, settings, expectedValues, stringSetting, expectedString)
}

// checkWindowsFirewallPublicProfile is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile.
//
// CIS Benchmark Audit list indices: 9.3.1-6
func checkWindowsFirewallPublicProfile(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\PublicProfile`
	settings := []string{"EnableFirewall", "DefaultInboundAction", "DefaultOutboundAction", "DisableNotifications",
		"AllowLocalPolicyMerge", "AllowLocalIPsecPolicyMerge"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1), uint64(0), uint64(0)}

	checkWindowsFirewallPublicProfileLogging(registryKey)
	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

}

// checkWindowsFirewallPublicProfileLogging is a helper function that checks the registry to determine if the system is configured with the correct settings for the public profile logging.
//
// CIS Benchmark Audit list indices: 9.3.7-10
func checkWindowsFirewallPublicProfileLogging(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsFirewall\PublicProfile\Logging`
	stringSetting := []string{"LogFilePath"}
	settings := []string{"LogFileSize", "LogDroppedPackets", "LogSuccessfulConnections"}

	expectedString := []string{`%SYSTEMROOT%\System32\logfiles\firewall\publicfw.log`}
	expectedValues := []interface{}{[]uint64{16384, ^uint64(0)}, uint64(1), uint64(1)}

	CheckIntegerStringRegistrySettings(registryKey, registryPath, settings, expectedValues,
		stringSetting, expectedString)
}

// policiesWindowsInkWorkspace is a helper function that checks the registry to determine if the system is configured with the correct settings for Windows Ink Workspace.
//
// CIS Benchmark Audit list index: 18.9.89.2
func policiesWindowsInkWorkspace(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsInkWorkspace`

	settings := []string{"AllowWindowsInkWorkspace"}

	expectedValues := []interface{}{[]uint64{0, 1}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesWindowsStore is a helper function that checks the registry to determine if the system is configured with the correct settings for the Windows Store.
//
// CIS Benchmark Audit list indices: 18.9.75.2-4
func policiesWindowsStore(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\WindowsStore`

	settings := []string{"RequirePrivateStoreOnly", "AutoDownload", "DisableOSUpgrade"}

	expectedValues := []interface{}{uint64(1), uint64(4), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// policiesEarlyLaunch is a helper function that checks the registry to determine if the system is configured with the correct settings for early launch.
//
// CIS Benchmark Audit list index: 18.8.14.1
func policiesEarlyLaunch(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Policies\EarlyLaunch`

	settings := []string{"DriverLoadPolicy"}

	expectedValues := []interface{}{uint64(3)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
