package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

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
