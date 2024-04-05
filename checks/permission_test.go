package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"

	"github.com/stretchr/testify/assert"
)

// TestPermission is a function that tests the Permission function's ability to correctly return permissions.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate different sets of permissions. Each test case checks if the Permission function correctly identifies the presence or absence of specific permissions based on the simulated registry keys. The function asserts that the returned permissions match the expected results.
func TestPermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		key        registrymock.RegistryKey
		want       checks.Check
	}{
		{
			name:       "NonPackagedWebcamPermissionExists",
			permission: "webcam",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "NonPackaged",
								SubKeys: []registrymock.MockRegistryKey{
									{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
								},
							},
						},
					},
				},
			},
			want: checks.NewCheckResult("webcam", "microsoft.webcam"),
		},
		{
			name:       "WebcamPermissionExists",
			permission: "webcam",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\webcam",
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
						},
					},
				},
			},
			want: checks.NewCheckResult("webcam", "microsoft.webcam"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checks.Permission(tc.permission, tc.key)
			require.Equal(t, tc.want, result)
		})
	}
}

// TestFormatPermission is a function that tests the Permission function's ability to correctly format the returned permissions.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with a specific scenario where the permission string is in a certain format. It uses a mock implementation of the RegistryKey interface to simulate a specific permission format. The test case checks if the Permission function correctly formats the returned permission string by removing any '#' characters. The function asserts that the returned permission string matches the expected format.
func TestFormatPermission(t *testing.T) {
	key := &registrymock.MockRegistryKey{
		SubKeys: []registrymock.MockRegistryKey{
			{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\location",
				StringValues: map[string]string{"Value": "Allow"},
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "NonPackaged",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "test#test#test.exe"},
						},
					},
				},
			},
		},
	}
	c := checks.Permission("location", key)
	assert.NotContains(t, c.Result, "#")
	assert.Contains(t, c.Result, "test.exe")
}

// TestNonExistingPermission is a function that tests the Permission function's behavior when the requested permission does not exist.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Permission function with a scenario where the requested permission does not exist in the simulated registry keys. It uses a mock implementation of the RegistryKey interface to simulate this scenario. The test case checks if the Permission function correctly returns an error when the requested permission does not exist. The function asserts that the returned Check instance contains the expected error message.
func TestNonExistingPermission(t *testing.T) {
	key := &registrymock.MockRegistryKey{
		SubKeys: []registrymock.MockRegistryKey{
			{KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore\\location",
				StringValues: map[string]string{"Value": "Allow"},
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "NonPackaged",
						StringValues: map[string]string{"Value": "Allow"},
						SubKeys: []registrymock.MockRegistryKey{
							{KeyName: "test.test"},
						},
					},
				},
			},
		},
	}
	c := checks.Permission("hello", key)
	assert.Equal(t, c.Result, []string(nil))
	assert.EqualError(t, c.Error, "error opening registry key: error opening registry key: key not found")
}
