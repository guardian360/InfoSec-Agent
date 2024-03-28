package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestPermissionRegistry(t *testing.T) {
//	subKey := registrymock.MockRegistryKey{KeyName: "Microsoft.Gaming", StringValues: map[string]string{"Value": "Allow"}, Err: nil}
//	tests := []struct {
//		name string
//		key  registrymock.RegistryKey
//		want []checks.Check
//	}{
//		{
//			name: "webcam",
//			key:  &registrymock.MockRegistryKey{KeyName: "webcam", SubKeys: []registrymock.MockRegistryKey{subKey}, Err: nil},
//			want: []checks.Check{checks.NewCheckResult("webcam", "Allow")},
//		},
//	}
//	//tests[0].key.(*registrymock.MockRegistryKey).SubKeys = []registrymock.MockRegistryKey{KeyName: "Microsoft.Gaming", StringValues: map[string]string{"Value": "Allow"}, Err: nil}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := checks.Permission("webcam", tt.key); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Permission() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestPermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		key        registrymock.RegistryKey
		want       checks.Check
	}{
		{
			name:       "PermissionExistsWithApps",
			permission: "webcam",
			key: &registrymock.MockRegistryKey{
				KeyName: "Software\\Microsoft\\Windows\\CurrentVersion\\CapabilityAccessManager\\ConsentStore",
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "microsoft.webcam", StringValues: map[string]string{"Value": "Allow"}},
				},
			},
			want: checks.NewCheckResult("webcam", "microsoft.webcam"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checks.Permission(tc.permission, tc.key)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Test %s failed. Expected %#v, got %#v", tc.name, tc.want, result)
			}
		})
	}
}

// TestInputPermissionCheck tests the input of the permission check function
//
// Parameters: permission (string) represents the permission to check
//
// Returns: _
func TestInputPermission(t *testing.T) {
	testCases := []string{"/", " ", "test", "camera"}
	for _, permission := range testCases {
		c := checks.Permission(permission, &registrymock.MockRegistryKey{})
		assert.Nil(t, c.Result)
		assert.NotNil(t, c.ErrorMSG)
	}
}

// TestValidInputPermission tests valid permissions
//
// Parameters: permission (string) represents the permission to check
//
// Returns: _
func TestValidPermissions(t *testing.T) {
	testCases := []string{"webcam", "microphone"}
	for _, permission := range testCases {
		c := checks.Permission(permission, &registrymock.MockRegistryKey{})
		assert.Contains(t, c.Result, "Microsoft.WindowsCamera")
	}
}

// TestFormatPermission tests if the correct format is returned
//
// Parameters: permission (string) represents the permission to check
//
// Returns: _
func TestFormatPermission(t *testing.T) {
	testCases := []string{"location"}
	for _, permission := range testCases {
		c := checks.Permission(permission, &registrymock.MockRegistryKey{})
		assert.NotContains(t, c.Result, "#")
		assert.NotContains(t, c.Result, " ")
		assert.NotContains(t, c.Result, "_")
	}
}
