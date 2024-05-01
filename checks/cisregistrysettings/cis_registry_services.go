package cisregistrysettings

import "github.com/InfoSec-Agent/InfoSec-Agent/mocking"

// CheckServices is a function that checks various registry settings related to different services
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
func CheckServices(registryKey mocking.RegistryKey) []bool {
	results := make([]bool, 0)
	results = append(results, checkservicesDisabled(registryKey, servicesDisabledRegistryPaths)...)
	for _, check := range serviceChecks {
		results = append(results, check(registryKey)...)
	}

	return results
}

// serviceChecks is a collection of registry checks related to different services.
// Each function in the collection represents a different service check that the application can perform.
// The registry settings get checked against the CIS Benchmark recommendations.
var serviceChecks = []func(mocking.RegistryKey) []bool{
	servicesEventLog,
	servicesLDAP,
	servicesNetBTParameters,
	servicesNetlogonParameters,
	servicesTCPIP,
	servicesTCPIP6,
}

// servicesDisabledRegistryPaths is a collection of paths to services that should be disabled.
//
// CIS Benchmark Audit list indices: 5.3, 5.6-8, 5.10-11, 5.13, 5.24, 5.26, 5.28, 5.30-33, 5.36-37, 5.41-45, 18.3.2
var servicesDisabledRegistryPaths = []string{
	`SYSTEM\CurrentControlSet\Services\XboxNetApiSvc`,
	`SYSTEM\CurrentControlSet\Services\XboxGipSvc`,
	`SYSTEM\CurrentControlSet\Services\XblGameSave`,
	`SYSTEM\CurrentControlSet\Services\XblAuthManager`,
	`SYSTEM\CurrentControlSet\Services\WMSvc`,
	`SYSTEM\CurrentControlSet\Services\WMPNetworkSvc`,
	`SYSTEM\CurrentControlSet\Services\W3SVC`,
	`SYSTEM\CurrentControlSet\Services\upnphost`,
	`SYSTEM\CurrentControlSet\Services\sshd`,
	`SYSTEM\CurrentControlSet\Services\SSDPSRV`,
	`SYSTEM\CurrentControlSet\Services\simptcp`,
	`SYSTEM\CurrentControlSet\Services\SharedAccess`,
	`SYSTEM\CurrentControlSet\Services\sacsvr`,
	`SYSTEM\CurrentControlSet\Services\RpcLocator`,
	`SYSTEM\CurrentControlSet\Services\RemoteAccess`,
	`SYSTEM\CurrentControlSet\Services\mrxsmb10`,
	`SYSTEM\CurrentControlSet\Services\LxssManager`,
	`SYSTEM\CurrentControlSet\Services\irmon`,
	`SYSTEM\CurrentControlSet\Services\icssvc`,
	`SYSTEM\CurrentControlSet\Services\IISADMIN`,
	`SYSTEM\CurrentControlSet\Services\FTPSVC`,
	`SYSTEM\CurrentControlSet\Services\Browser`,
}

// servicesEventLog is a helper function that checks the registry to determine if the system is configured with the correct settings for the event log service.
//
// CIS Benchmark Audit list index: 18.4.13
func servicesEventLog(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\Eventlog\Security`

	settings := []string{"WarningLevel"}

	expectedValues := []interface{}{[]uint64{1, 90}}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// servicesLDAP is a helper function that checks the registry to determine if the system is configured with the correct settings for the LDAP service.
//
// CIS Benchmark Audit list index: 2.3.11.8
func servicesLDAP(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\LDAP`

	settings := []string{"LDAPClientIntegrity"}

	expectedValues := []interface{}{uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// servicesNetBTParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for the NetBT parameters.
//
// CIS Benchmark Audit list indices: 18.3.6, 18.4.7
func servicesNetBTParameters(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\NetBT\Parameters`

	settings := []string{"NodeType", "nonamereleaseondemand"}

	expectedValues := []interface{}{uint64(2), uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// servicesNetlogonParameters is a helper function that checks the registry to determine if the system is configured with the correct settings for the Netlogon parameters.
//
// CIS Benchmark Audit list indices: 2.3.6.1-6
func servicesNetlogonParameters(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\Netlogon\Parameters`

	settings := []string{"RequireSignOrSeal", "SealSecureChannel", "SignSecureChannel",
		"DisablePasswordChange", "MaximumPasswordAge", "RequireStrongKey"}

	expectedValues := []interface{}{uint64(1), uint64(1), uint64(1), uint64(0), []uint64{1, 30}, uint64(1)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// servicesTCPIP is a helper function that checks the registry to determine if the system is configured with the correct settings for the TCPIP service.
//
// CIS Benchmark Audit list indices: 18.4.3, 18.4.5
func servicesTCPIP(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`

	settings := []string{"DisableIPSourceRouting", "EnableICMPRedirect"}

	expectedValues := []interface{}{uint64(2), uint64(0)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// servicesTCPIP6 is a helper function that checks the registry to determine if the system is configured with the correct settings for the TCPIP6 service.
//
// CIS Benchmark Audit list indices: 18.4.2
func servicesTCPIP6(registryKey mocking.RegistryKey) []bool {
	registryPath := `SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`

	settings := []string{"DisableIPSourceRouting"}

	expectedValues := []interface{}{uint64(2)}

	return CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)
}

// checkservicesDisabled is a helper function that checks the registry to determine if the system is configured with the correct settings for the services that should be disabled.
func checkservicesDisabled(registryKey mocking.RegistryKey, registryPaths []string) []bool {
	results := make([]bool, 0)
	for _, registryPath := range registryPaths {
		settings := []string{"Start"}

		expectedValues := []interface{}{uint64(4)}

		results = append(results, CheckIntegerRegistrySettings(registryKey, registryPath, settings, expectedValues)...)
	}
	return results
}
