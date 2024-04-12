package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TestSecureBoot is a function that tests the behavior of the SecureBoot function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the SecureBoot function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the Secure Boot registry key. Each test case checks if the SecureBoot function correctly identifies the status of Secure Boot (enabled, disabled, or unknown) based on the simulated registry key value. The function asserts that the returned Check instance contains the expected results.
func TestSecureBoot(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
	}{
		{
			name: "SecureBootEnabled",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 1, "Secure boot is enabled"),
		},
		{
			name: "SecureBootDisabled",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 0}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 0, "Secure boot is disabled"),
		},
		{
			name: "SecureBootUnknown",
			key: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 2}, Err: nil}}},
			want: checks.NewCheckResult(checks.SecureBootID, 2, "Secure boot status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.SecureBoot(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
