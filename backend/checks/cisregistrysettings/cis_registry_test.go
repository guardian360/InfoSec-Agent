package cisregistrysettings_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/cisregistrysettings"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.SetupTests()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCheckWin11(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "DNS Client set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient",
						IntegerValues: map[string]uint64{"DoHPolicy": uint64(2)}},
				}},
			want: true,
		},
		{
			name: "DNS Client not set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient",
						IntegerValues: map[string]uint64{"DoHPolicy": uint64(1)}},
				}},
			want: false,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckWin11(tt.key)
			if cisregistrysettings.RegistrySettingsMap["SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient\\DoHPolicy"] != tt.want {
				t.Errorf("CheckWin11() = %v, want %v", cisregistrysettings.RegistrySettingsMap["SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient"], tt.want)
			}
		})
	}
}

func TestCheckWin10(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "Nothing set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{}},
			want: false,
		},
		{
			name: "Everything set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Print",
						IntegerValues: map[string]uint64{"RpcAuthnLevelPrivacyEnabled": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows NT\\DNSClient",
						IntegerValues: map[string]uint64{"EnableNetbios": uint64(2)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows NT\\Printers",
						IntegerValues: map[string]uint64{"RedirectionGuardPolicy": uint64(1), "CopyFilesPolicy": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows NT\\Printers\\RPC",
						IntegerValues: map[string]uint64{"RpcUseNamedPipeProtocol": uint64(0), "RpcAuthentication": uint64(0), "RpcProtocols": uint64(5), "ForceKerberosForRpc": uint64(0), "RpcTcpPort": uint64(0)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\System",
						IntegerValues: map[string]uint64{"AllowCustomSSPsAPs": uint64(0)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\AppInstaller",
						IntegerValues: map[string]uint64{"EnableAppInstaller": uint64(0), "EnableExperimentalFeatures": uint64(0), "EnableHashOverride": uint64(0), "EnableMSAppInstallerProtocol": uint64(0)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Internet Explorer\\Main",
						IntegerValues: map[string]uint64{"NotifyDisableIEOptions": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows Defender\\Windows Defender Exploit Guard\\ASR\\Rules",
						IntegerValues: map[string]uint64{"56a863a9-875e-4185-98a7-b882c64b5ce5": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System",
						IntegerValues: map[string]uint64{"EnableMPR": uint64(0)},
					},
				},
			},

			want: true,
		}}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckWin10(tt.key)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckWin10() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckPoliciesHKU(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "Nothing set up correctly",
			key:  &mocking.MockRegistryKey{},
			want: false,
		},
		{
			name: "Everything set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\Attachments",
						IntegerValues: map[string]uint64{"SaveZoneInformation": uint64(2), "ScanWithAntiVirus": uint64(3)},
					},
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\Explorer",
						IntegerValues: map[string]uint64{"NoInplaceSharing": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\CloudContent",
						IntegerValues: map[string]uint64{"ConfigureWindowsSpotlight": uint64(2), "DisableThirdPartySuggestions": uint64(1), "DisableSpotlightCollectionOnDesktop": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\Control Panel\\Desktop",
						IntegerValues: map[string]uint64{"ScreenSaveActive": uint64(1), "ScreenSaverIsSecure": uint64(1), "ScreenSaveTimeOut": uint64(900)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\CurrentVersion\\PushNotifications",
						IntegerValues: map[string]uint64{"NoToastApplicationNotificationOnLockScreen": uint64(1)},
					},
					{
						KeyName:       "SOFTWARE\\Policies\\Microsoft\\Windows\\Installer",
						IntegerValues: map[string]uint64{"AlwaysInstallElevated": uint64(0)},
					},
				},
			},
			want: true,
		}}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckPoliciesHKU(tt.key)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckPoliciesHKU() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckServices(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "Nothing set up correctly",
			key:  &mocking.MockRegistryKey{},
			want: false,
		},
		{
			name: "Everything set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\Eventlog\\Security",
						IntegerValues: map[string]uint64{"WarningLevel": uint64(5)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\LDAP",
						IntegerValues: map[string]uint64{"LDAPClientIntegrity": uint64(1)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\NetBT\\Parameters",
						IntegerValues: map[string]uint64{"NodeType": uint64(2), "nonamereleaseondemand": uint64(1)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\Netlogon\\Parameters",
						IntegerValues: map[string]uint64{"RequireSignOrSeal": uint64(1), "SealSecureChannel": uint64(1), "SignSecureChannel": uint64(1), "DisablePasswordChange": uint64(0), "MaximumPasswordAge": uint64(15), "RequireStrongKey": uint64(1)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters",
						IntegerValues: map[string]uint64{"DisableIPSourceRouting": uint64(2), "EnableICMPRedirect": uint64(0)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\Tcpip6\\Parameters",
						IntegerValues: map[string]uint64{"DisableIPSourceRouting": uint64(2)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\XboxNetApiSvc",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\XboxGipSvc",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\XblGameSave",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\XblAuthManager",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\WMSvc",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\WMPNetworkSvc",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\W3SVC",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\upnphost",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\sshd",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\SSDPSRV",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\simptcp",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\SharedAccess",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\sacsvr",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\RpcLocator",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\RemoteAccess",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\mrxsmb10",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\LxssManager",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\irmon",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\icssvc",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\IISADMIN",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\FTPSVC",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Services\\Browser",
						IntegerValues: map[string]uint64{"Start": uint64(4)},
					},
				}},
			want: true,
		}}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckServices(tt.key)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckServices() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckOtherRegistrySettings(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "Nothing set up correctly",
			key:  &mocking.MockRegistryKey{},
			want: false,
		},
		{
			name: "Everything set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SOFTWARE\\Microsoft\\WcmSvc\\wifinetworkmanager\\config",
						IntegerValues: map[string]uint64{"AutoConnectAllowedOEM": uint64(0)},
					},
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
						IntegerValues: map[string]uint64{"AllocateDASD": uint64(2), "PasswordExpiryWarning": uint64(5), "ScRemoveOption": uint64(1), "AutoAdminLogon": uint64(0), "ScreenSaverGracePeriod": uint64(0)},
					},
					{
						KeyName:      "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon\\GPExtensions\\{D76B9641-3288-4f75-942D-087DE603E3EA}",
						StringValues: map[string]string{"DllName": "C:\\Program Files\\LAPS\\CSE\\AdmPwd.dll"},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Lsa",
						IntegerValues: map[string]uint64{"LimitBlankPasswordUse": uint64(1), "SCENoApplyLegacyAuditPolicy": uint64(1), "CrashOnAuditFail": uint64(0), "RestrictAnonymousSAM": uint64(1), "RestrictAnonymous": uint64(1), "DisableDomainCreds": uint64(1), "EveryoneIncludesAnonymous": uint64(0), "restrictremotesam": uint64(1), "ForceGuest": uint64(0), "UseMachineId": uint64(1), "NoLMHash": uint64(1), "LMCompatibilityLevel": uint64(5)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Lsa\\MSV1_0",
						IntegerValues: map[string]uint64{"AllowNullSessionFallback": uint64(0), "NTLMMinClientSec": uint64(537395200), "NTLMMinServerSec": uint64(537395200)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Lsa\\pku2u",
						IntegerValues: map[string]uint64{"AllowOnlineID": uint64(0)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SAM",
						IntegerValues: map[string]uint64{"RelaxMinimumPasswordLengthLimits": uint64(1)},
					},
					{
						KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecurePipeServers\\Winreg\\AllowedExactPaths",
						StringValues: map[string]string{"Machine": "" +
							"System\\CurrentControlSet\\Control\\ProductOptionsSystem\\CurrentControlSet\\Control\\" +
							"Server ApplicationsSoftware\\Microsoft\\Windows NT\\CurrentVersion"},
					},
					{
						KeyName: "SYSTEM\\CurrentControlSet\\Control\\SecurePipeServers\\Winreg\\AllowedPaths",
						StringValues: map[string]string{"Machine": "" +
							"System\\CurrentControlSet\\Control\\Print\\PrintersSystem\\CurrentControlSet\\Services\\EventlogSoftware" +
							"\\Microsoft\\OLAP ServerSoftware\\Microsoft\\Windows NT\\CurrentVersion\\PrintSoftware\\Microsoft\\Windows NT" +
							"\\CurrentVersion\\WindowsSystem\\CurrentControlSet\\Control\\ContentIndexSystem\\CurrentControlSet\\Control" +
							"\\Terminal ServerSystem\\CurrentControlSet\\Control\\Terminal Server\\UserConfigSystem\\CurrentControlSet" +
							"\\Control\\Terminal Server\\DefaultUserConfigurationSoftware\\Microsoft\\Windows NT\\CurrentVersion" +
							"\\PerflibSystem\\CurrentControlSet\\Services\\Sysmonlog"},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecurityProviders\\WDigest",
						IntegerValues: map[string]uint64{"UseLogonCredential": uint64(0)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Session Manager",
						IntegerValues: map[string]uint64{"ProtectionMode": uint64(1), "SafeDllSearchMode": uint64(1)},
					},
					{
						KeyName:       "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Kernel",
						IntegerValues: map[string]uint64{"ObCaseInsensitive": uint64(1), "DisableExceptionChainValidation": uint64(0)},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckOtherRegistrySettings(tt.key)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckOtherRegistrySettings() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckPoliciesHKLM(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want bool
	}{
		{
			name: "Nothing set up correctly",
			key:  &mocking.MockRegistryKey{},
			want: false,
		},
		{
			name: "Some setting is set up correctly",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
						IntegerValues: map[string]uint64{"EnumerateAdministrators": uint64(0)},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckPoliciesHKLM(tt.key)
			if cisregistrysettings.RegistrySettingsMap[`SOFTWARE\Microsoft\Windows\Currentversion\Policies\Credui\EnumerateAdministrators`] != tt.want {
				t.Errorf("CheckPoliciesHKLM() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCISRegistrySettings(t *testing.T) {
	tests := []struct {
		name   string
		lmKey  mocking.RegistryKey
		usrKey mocking.RegistryKey
		want   checks.Check
	}{
		{
			name:   "Nothing set up correctly",
			lmKey:  &mocking.MockRegistryKey{},
			usrKey: &mocking.MockRegistryKey{},
			want:   checks.NewCheckResult(checks.CISRegistrySettingsID, 0),
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			got := cisregistrysettings.CISRegistrySettings(tt.lmKey, tt.usrKey)
			if got.ResultID != tt.want.ResultID {
				t.Errorf("CISRegistrySettings() = %v, want %v", got, tt)
			}
		})
	}
}

func TestCheckIntegerStringRegistrySettings(t *testing.T) {
	tests := []struct {
		name             string
		key              mocking.RegistryKey
		path             string
		integerSettings  []string
		expectedIntegers []interface{}
		stringSettings   []string
		expectedStrings  []string
		want             bool
	}{
		{
			name:             "Check nothing",
			key:              &mocking.MockRegistryKey{},
			path:             "",
			integerSettings:  []string{},
			expectedIntegers: []interface{}{},
			stringSettings:   []string{},
			expectedStrings:  []string{},
			want:             false,
		},
		{
			name: "Check some int and string setting",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
						IntegerValues: map[string]uint64{"EnumerateAdministrators": uint64(0)},
						StringValues:  map[string]string{"EnumerateAdministrators": "0"},
					},
				},
			},
			path:             "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
			integerSettings:  []string{"EnumerateAdministrators"},
			expectedIntegers: []interface{}{uint64(0)},
			stringSettings:   []string{"EnumerateAdministrators"},
			expectedStrings:  []string{"0"},
			want:             true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckIntegerStringRegistrySettings(tt.key, tt.path, tt.integerSettings, tt.expectedIntegers, tt.stringSettings, tt.expectedStrings)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckIntegerStringRegistrySettings() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckStringRegistrySettings(t *testing.T) {
	tests := []struct {
		name            string
		key             mocking.RegistryKey
		path            string
		stringSettings  []string
		expectedStrings []string
		want            bool
	}{
		{
			name:            "Check nothing",
			key:             &mocking.MockRegistryKey{},
			path:            "",
			stringSettings:  []string{},
			expectedStrings: []string{},
			want:            false,
		},
		{
			name: "Check some string setting",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:      "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
						StringValues: map[string]string{"EnumerateAdministrators": "0"},
					},
				},
			},
			path:            "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
			stringSettings:  []string{"EnumerateAdministrators"},
			expectedStrings: []string{"0"},
			want:            true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckStringRegistrySettings(tt.key, tt.path, tt.stringSettings, tt.expectedStrings)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckStringRegistrySettings() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestCheckIntegerRegistrySettings(t *testing.T) {
	tests := []struct {
		name             string
		key              mocking.RegistryKey
		path             string
		integerSettings  []string
		expectedIntegers []interface{}
		want             bool
	}{
		{
			name:             "Check nothing",
			key:              &mocking.MockRegistryKey{},
			path:             "",
			integerSettings:  []string{},
			expectedIntegers: []interface{}{},
			want:             false,
		},
		{
			name: "Check some integer setting",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName:       "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
						IntegerValues: map[string]uint64{"EnumerateAdministrators": uint64(0)},
					},
				},
			},
			path:             "SOFTWARE\\Microsoft\\Windows\\Currentversion\\Policies\\Credui",
			integerSettings:  []string{"EnumerateAdministrators"},
			expectedIntegers: []interface{}{uint64(0)},
			want:             true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			cisregistrysettings.CheckIntegerRegistrySettings(tt.key, tt.path, tt.integerSettings, tt.expectedIntegers)
			if allValuesTrue(cisregistrysettings.RegistrySettingsMap) != tt.want {
				t.Errorf("CheckIntegerRegistrySettings() = %v, want %v", allValuesTrue(cisregistrysettings.RegistrySettingsMap), tt.want)
			}
		})
	}
}

func TestOpenRegistryKeyWithErrHandling(t *testing.T) {
	tests := []struct {
		name    string
		key     mocking.RegistryKey
		path    string
		want    bool
		wantErr bool
	}{
		{
			name: "Open key successfully",
			key:  &mocking.MockRegistryKey{},
			path: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System",
			want: true,
		},
	}
	for _, tt := range tests {
		cisregistrysettings.RegistrySettingsMap = map[string]bool{}
		t.Run(tt.name, func(t *testing.T) {
			got, _ := cisregistrysettings.OpenRegistryKeyWithErrHandling(tt.key, tt.path)
			if got != nil && tt.want == false {
				t.Errorf("OpenRegistryKeyWithErrHandling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckStringValue(t *testing.T) {
	tests := []struct {
		name           string
		key            mocking.RegistryKey
		stringSetting  string
		expectedString string
		want           bool
	}{
		{
			name:           "Check nothing",
			key:            &mocking.MockRegistryKey{},
			stringSetting:  "",
			expectedString: "",
			want:           false,
		},
		{
			name: "Check some string setting",
			key: &mocking.MockRegistryKey{
				StringValues: map[string]string{"EnumerateAdministrators": "0"},
			},
			stringSetting:  "EnumerateAdministrators",
			expectedString: "0",
			want:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cisregistrysettings.CheckStringValue(tt.key, tt.stringSetting, tt.expectedString)
			if got != tt.want {
				t.Errorf("CheckStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckIntegerValue(t *testing.T) {
	tests := []struct {
		name            string
		key             mocking.RegistryKey
		stringSetting   string
		expectedInteger interface{}
		want            bool
	}{
		{
			name:            "Check nothing",
			key:             &mocking.MockRegistryKey{},
			stringSetting:   "",
			expectedInteger: nil,
			want:            false,
		},
		{
			name: "Check some integer setting",
			key: &mocking.MockRegistryKey{
				IntegerValues: map[string]uint64{"EnumerateAdministrators": uint64(0)},
			},
			stringSetting:   "EnumerateAdministrators",
			expectedInteger: uint64(0),
			want:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cisregistrysettings.CheckIntegerValue(tt.key, tt.stringSetting, tt.expectedInteger)
			if got != tt.want {
				t.Errorf("CheckStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if all values in a map are true
func allValuesTrue(m map[string]bool) bool {
	if (len(m)) == 0 {
		return false
	}
	for _, value := range m {
		if !value {
			return false
		}
	}
	return true
}
