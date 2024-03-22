package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInputPermissionCheck tests the input of the permission check function
//
// Parameters: permission (string) represents the permission to check
//
// Returns: _
func TestInputPermission(t *testing.T) {
	testCases := []string{"/", " ", "test", "cammera"}
	for _, permission := range testCases {
		check := Permission(permission)
		assert.Nil(t, check.Result)
		assert.NotNil(t, check.ErrorMSG)
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
		check := Permission(permission)
		assert.Contains(t, check.Result, "Microsoft.WindowsCamera")
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
		check := Permission(permission)
		assert.NotContains(t, check.Result, "#")
		assert.NotContains(t, check.Result, " ")
		assert.NotContains(t, check.Result, "_")
	}
}
