package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

// CheckOtherRegistrySettings is a function that checks various registry settings related to different (unrelated) registry keys
// to ensure they adhere to the CIS Benchmark standards.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_LOCAL_MACHINE for this function.
//
// Returns: None
func CheckOtherRegistrySettings(registryKey mocking.RegistryKey) {
	for _, check := range generalRegistryChecks {
		check(registryKey)
	}
}

// generalRegistryChecks is a collection of registry checks related to different (unrelated) registry keys.
// Each function in the collection represents a different registry setting check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var generalRegistryChecks = []func(mocking.RegistryKey){
	CheckAutoConnectHotspot,
	CheckCurrentVersionRegistry,
	CheckControlLsa,
	CheckControlSAM,
	CheckSecurePipeServers,
	CheckWDigest,
	CheckSessionManager,
}

// CheckAutoConnectHotspot is a helper function that checks the registry to determine if the system is configured to automatically connect to hotspots.
//
// CIS Benchmark Audit list index: 18.5.23.2.1
func CheckAutoConnectHotspot(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\WcmSvc\wifinetworkmanager\config`

	settings := []string{"AutoConnectAllowedOEM"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// CheckCurrentVersionRegistry is a helper function that checks the registry to determine if the system is configured with the correct settings for the current version.
//
// CIS Benchmark Audit list indices: 2.3.4.1, 2.3.7.8, 2.3.7.9, 18.2.1, 18.4.1, 18.4.10
func CheckCurrentVersionRegistry(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`

	settings := []string{
		"AllocateDASD",
		"PasswordExpiryWarning",
		"ScRemoveOption",
		"AutoAdminLogon",
		"ScreenSaverGracePeriod",
	}

	expectedValues := []interface{}{uint64(2), []uint64{5, 14}, []uint64{1, 2, 3}, uint64(0), []uint64{0, 5}}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)

	registryPath =
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\GPExtensions\{D76B9641-3288-4f75-942D-087DE603E3EA}`

	settings = []string{"DllName"}

	expectedStringValues := []string{"C:\\Program Files\\LAPS\\CSE\\AdmPwd.dll"}

	CheckStringRegistrySettings(registryKey, registryPath, settings, expectedStringValues)
}

// CheckControlLsa is a helper function that checks the registry to determine if the system is configured with the correct settings for Control Lsa.
//
// CIS Benchmark Audit list indices: 2.3.1.4, 2.3.2.1-2, 2.3.10.2-5, 2.3.10.10, 2.3.10.12,
// 2.3.11.1-3, 2.3.11.5, 2.3.11.7, 2.3.11.9-10
func CheckControlLsa(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Lsa`

	settings := []string{"LimitBlankPasswordUse", "SCENoApplyLegacyAuditPolicy", "CrashOnAuditFail",
		"RestrictAnonymousSAM", "RestrictAnonymous", "DisableDomainCreds", "EveryoneIncludesAnonymous",
		"restrictremotesam", "ForceGuest", "UseMachineId", "NoLMHash", "LMCompatibilityLevel"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(0), uint64(1), uint64(1), uint64(1), uint64(0),
		uint64(1), uint64(0), uint64(1), uint64(1), uint64(5)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	checkLsaMSV(registryKey)
	checkLsapku2u(registryKey)
}

// checkLsaMSV is a helper function that checks the registry to determine if the system is configured with the correct settings for Lsa MSV.
//
// CIS Benchmark Audit list indices: 2.3.11.2, 2.3.11.9-10
func checkLsaMSV(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Lsa\MSV1_0`

	settings := []string{"AllowNullSessionFallback", "NTLMMinClientSec", "NTLMMinServerSec"}

	expectedValues := []interface{}{uint64(0), uint64(537395200), uint64(537395200)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkLsapku2u is a helper function that checks the registry to determine if the system is configured with the correct settings for Lsa pku2u.
//
// CIS Benchmark Audit list index: 2.3.11.3
func checkLsapku2u(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Lsa\pku2u`

	settings := []string{"AllowOnlineID"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// CheckControlSAM is a helper function that checks the registry to determine if the system is configured with the correct settings for Control SAM.
//
// CIS Benchmark Audit list index: 1.1.6
func CheckControlSAM(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\SAM`

	settings := []string{"RelaxMinimumPasswordLengthLimits"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// CheckSecurePipeServers is a helper function that checks the registry to determine if the system is configured with the correct settings for secure pipe servers.
//
// CIS Benchmark Audit list indices: 2.3.10.7-8
func CheckSecurePipeServers(registryKey mocking.RegistryKey) {
	securePipeServersExactPaths(registryKey)
	securePipeServersPaths(registryKey)
}

// securePipeServersExactPaths is a helper function that checks the registry to determine if the system is configured with the correct settings for secure pipe servers exact paths.
//
// CIS Benchmark Audit list index: 2.3.10.7
func securePipeServersExactPaths(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\SecurePipeServers\Winreg\AllowedExactPaths`

	settings := []string{"Machine"}

	expectedValues := []string{"" +
		"System\\CurrentControlSet\\Control\\ProductOptionsSystem\\CurrentControlSet\\Control\\" +
		"Server ApplicationsSoftware\\Microsoft\\Windows NT\\CurrentVersion"}

	CheckStringRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// securePipeServersPaths is a helper function that checks the registry to determine if the system is configured with the correct settings for secure pipe servers paths.
//
// CIS Benchmark Audit list index: 2.3.10.8
func securePipeServersPaths(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\SecurePipeServers\Winreg\AllowedPaths`

	settings := []string{"Machine"}

	expectedValues := []string{"" +
		"System\\CurrentControlSet\\Control\\Print\\PrintersSystem\\CurrentControlSet\\Services\\EventlogSoftware" +
		"\\Microsoft\\OLAP ServerSoftware\\Microsoft\\Windows NT\\CurrentVersion\\PrintSoftware\\Microsoft\\Windows NT" +
		"\\CurrentVersion\\WindowsSystem\\CurrentControlSet\\Control\\ContentIndexSystem\\CurrentControlSet\\Control" +
		"\\Terminal ServerSystem\\CurrentControlSet\\Control\\Terminal Server\\UserConfigSystem\\CurrentControlSet" +
		"\\Control\\Terminal Server\\DefaultUserConfigurationSoftware\\Microsoft\\Windows NT\\CurrentVersion" +
		"\\PerflibSystem\\CurrentControlSet\\Services\\Sysmonlog"}

	CheckStringRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// CheckWDigest is a helper function that checks the registry to determine if the system is configured with the correct settings for WDigest.
//
// CIS Benchmark Audit list index: 18.3.7
func CheckWDigest(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\SecurityProviders\WDigest`

	settings := []string{"UseLogonCredential"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// CheckSessionManager is a helper function that checks the registry to determine if the system is configured with the correct settings for the session manager.
//
// CIS Benchmark Audit list indices: 2.3.15.1-2, 18.3.4, 18.4.9
func CheckSessionManager(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Session Manager`

	settings := []string{"ProtectionMode", "SafeDllSearchMode"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	sessionManagerKernel(registryKey)
}

// sessionManagerKernel is a helper function that checks the registry to determine if the system is configured with the correct settings for the session manager kernel.
//
// CIS Benchmark Audit list indices: 2.3.15.1, 18.3.4
func sessionManagerKernel(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Session Manager\Kernel`

	settings := []string{"ObCaseInsensitive", "DisableExceptionChainValidation"}

	expectedValues := []interface{}{uint64(1), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
