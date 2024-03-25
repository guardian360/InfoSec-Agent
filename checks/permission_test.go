package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionRegistry(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want []checks.Check
	}{
		{
			name: "webcam",
			key:  &registrymock.MockRegistryKey{Name: "Microsoft.Gaming", BinaryValue: nil, IntegerValue: 1, Err: nil},
			want: []checks.Check{checks.NewCheckResult("webcam", "Allow")},
		},
	}
	tests[0].key.(*registrymock.MockRegistryKey).SubKeys = []registrymock.MockRegistryKey{registrymock.MockRegistryKey{Name: "Value", StringValue: "Allow", IntegerValue: 1, Err: nil}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.Permission("webcam", tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Permission() = %v, want %v", got, tt.want)
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
	testCases := []string{"/", " ", "test", "cammera"}
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
