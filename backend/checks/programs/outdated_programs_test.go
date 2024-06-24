package programs_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.SetupTests()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestOutdatedPrograms(t *testing.T) {
	tests := []struct {
		name     string
		executor mocking.MockCommandExecutor
		key      mocking.RegistryKey
		want     checks.Check
		err      bool
	}{
		{
			name: "Winget program",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"Name": "Device1", "Version": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "Winget error",
			executor: mocking.MockCommandExecutor{
				Output: "test1",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"Name": "Device1", "Version": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{}...),
		},
		{
			name: "DisplayName found program",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Device1", "Version": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "Version found program",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Device1", "DisplayVersion": "1.0.0", "Pdublisher": "Device Manufacturer"},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "64bit program",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Device1", "DisplayVersion": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0", "Device1 | 1.0.0"}...),
		},
		{
			name: "error reading subkeyNames",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"test": "test"},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "error reading subkey",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:       "Git",
								IntegerValues: map[string]uint64{"test": 1},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "No Programs found",
			executor: mocking.MockCommandExecutor{
				Output: "",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"test": "test"},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{}...),
		},
		{
			name: "Empty names and duplicate programs",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Canon Inkjet Print Utility", "DisplayVersion": "1.0.0", "Publisher": "Device Manufacturer"},
							},
							{
								KeyName:      "Gifg",
								StringValues: map[string]string{"DisplayName": "microsoft defender", "DisplayVersion": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | 1.0.0"}...),
		},
		{
			name: "Empty names and duplicate programs error parsing number ",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Canon Inkjet Print Utility", "DisplayVersion": "Ǫ Ǫ 1.0.0.2.4.5.6", "Publisher": "Device Manufacturer"},
							},
							{
								KeyName:      "Gifg",
								StringValues: map[string]string{"DisplayName": "microsoft defender", "DisplayVersion": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
		{
			name: "Empty names and duplicate programs error parsing number ",
			executor: mocking.MockCommandExecutor{
				Output: "\r\n   -\r\n   \\\r\n\r\n\r\n   -\r\n   \\\r\n   |\r\n\r\nName                                  Id                                     Version              Available      Source\r\n-----------------------------------------------------------------------------------------------------------------------\r\nCanon Inkjet Print Utility            34791E63.CanonInkjetPrintUtility_6e5t�Ǫ 3.1.0.0",
			},
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{
						KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
						SubKeys: []mocking.MockRegistryKey{
							{
								KeyName:      "Git",
								StringValues: map[string]string{"DisplayName": "Canon Inkjet Print Utility", "DisplayVersion": "Ǫ 3.1.0.0", "Publisher": "Device Manufacturer"},
							},
							{
								KeyName:      "Gifg",
								StringValues: map[string]string{"DisplayName": "microsoft defender", "DisplayVersion": "1.0.0", "Publisher": "Device Manufacturer"},
							},
						},
					},
				},
			},

			want: checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, []string{"Canon Inkjet Print Utility | Ǫ 3.1.0.0"}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := programs.OutdatedSoftware(&tt.executor, tt.key)
			if tt.err {
				require.Error(t, got.Error)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
