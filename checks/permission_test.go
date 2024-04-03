package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"

	"github.com/stretchr/testify/assert"
)

// TestPermission tests if the correct permission is returned
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
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

// TestFormatPermission tests if the correct format is returned
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
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

// TestNonExistingPermission tests if the correct error is returned when the permission does not exist
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
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
