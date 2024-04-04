package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

func TestSecureBoot(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "SecureBootEnabled",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 1}, Err: nil}}},
			want: checks.NewCheckResult(18, 1, "Secure boot is enabled"),
		},
		{
			name: "SecureBootDisabled",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 0}, Err: nil}}},
			want: checks.NewCheckResult(18, 0, "Secure boot is disabled"),
		},
		{
			name: "SecureBootUnknown",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{{
				KeyName:       "SYSTEM\\CurrentControlSet\\Control\\SecureBoot\\State",
				IntegerValues: map[string]uint64{"UEFISecureBootEnabled": 2}, Err: nil}}},
			want: checks.NewCheckResult(18, 2, "Secure boot status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.SecureBoot(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}
