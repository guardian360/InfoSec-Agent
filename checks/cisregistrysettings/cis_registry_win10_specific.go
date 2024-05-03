package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/mocking"

// CheckWin10 is a function that checks various registry settings specific to the Windows 10 CIS Benchmark Audit list.
// It takes a RegistryKey object as an argument, which represents the root key from which the registry settings will be checked.
// The function returns a slice of boolean values, where each boolean represents whether a particular registry setting adheres to the CIS Benchmark standards.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): The root key from which the registry settings will be checked. Should be HKEY_LOCAL_MACHINE for this function.
//
// Returns: None
func CheckWin10(registryKey mocking.RegistryKey) {
	for _, check := range checksWin10 {
		check(registryKey)
	}
}

// checksWin10 is a collection of registry checks specific to the Windows 10 CIS Benchmark Audit list.
// Each function in the collection represents a different registry setting check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var checksWin10 = []func(mocking.RegistryKey){
	win10Print,
	win10DNSClient,
	win10Printers,
	win10System,
	win10AppInstaller,
	win10InternetExplorerMain,
	win10ExploitGuardRules,
	win10PoliciesSystem,
}

// win10Print is a helper function that checks the registry to determine if the system is configured with the correct settings for the Print control.
//
// CIS Benchmark Audit list index: 18.4.2
func win10Print(registryKey mocking.RegistryKey) {
	registryPath := `SYSTEM\CurrentControlSet\Control\Print`

	settings := []string{"RpcAuthnLevelPrivacyEnabled"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// wind10DNSClient is a helper function that checks the registry to determine if the system is configured with the correct settings for the DNS Client.
//
// CIS Benchmark Audit list index: 18.6.4.1
func win10DNSClient(registryKey mocking.RegistryKey) {
	registryPath := DNSClientRegistryPath

	settings := []string{"EnableNetbios"}

	expectedValues := []interface{}{uint64(2)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10Printers is a helper function that checks the registry to determine if the system is configured with the correct settings for the Printers.
//
// CIS Benchmark Audit list index: 18.7.2-7 18.7.9
func win10Printers(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows NT\Printers`

	settings := []string{"RedirectionGuardPolicy", "CopyFilesPolicy"}

	expectedValues := []interface{}{uint64(1), uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
	win10PrintersRPC(registryKey)
}

// win10PrintersRPC is a helper function that checks the registry to determine if the system is configured with the correct settings for the Printers RPC.
//
// CIS Benchmark Audit list index: 18.7.3-7
func win10PrintersRPC(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows NT\Printers\RPC`

	settings := []string{"RpcUseNamedPipeProtocol", "RpcAuthentication",
		"RpcProtocols", "ForceKerberosForRpc", "RpcTcpPort"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(5), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10System is a helper function that checks the registry to determine if the system is configured with the correct settings for the System.
//
// CIS Benchmark Audit list index: 18.9.25.1
func win10System(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\System`

	settings := []string{"AllowCustomSSPsAPs"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10AppInstaller is a helper function that checks the registry to determine if the system is configured with the correct settings for the AppInstaller.
//
// CIS Benchmark Audit list indices: 18.10.17.1-4
func win10AppInstaller(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows\AppInstaller`

	settings := []string{"EnableAppInstaller", "EnableExperimentalFeatures",
		"EnableHashOverride", "EnableMSAppInstallerProtocol"}

	expectedValues := []interface{}{uint64(0), uint64(0), uint64(0), uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10InternetExplorerMain is a helper function that checks the registry to determine if the system is configured with the correct settings for Internet Explorer
//
// CIS Benchmark Audit list index: 18.10.35.1
func win10InternetExplorerMain(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Internet Explorer\Main`

	settings := []string{"NotifyDisableIEOptions"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10ExploitGuardRules is a helper function that checks the registry to determine if the system is configured with the correct settings for the Exploit Guard Rules.
//
// CIS Benchmark Audit list index: 18.10.43.6.1.2
func win10ExploitGuardRules(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Windows Defender Exploit Guard\ASR\Rules`

	settings := []string{"56a863a9-875e-4185-98a7-b882c64b5ce5"}

	expectedValues := []interface{}{uint64(1)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// win10PoliciesSystem is a helper function that checks the registry to determine if the system is configured with the correct policy settings for the System.
//
// CIS Benchmark Audit list index: 18.10.82.1
func win10PoliciesSystem(registryKey mocking.RegistryKey) {
	registryPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`

	settings := []string{"EnableMPR"}

	expectedValues := []interface{}{uint64(0)}

	CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}
