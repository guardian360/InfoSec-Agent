package checks_test

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

// TestLoginMethod tests the LoginMethod function with (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestLoginMethod(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "Login method is PIN",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{
				{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
					StringValues: map[string]string{"S-1-5-21-174387295-310396468-1212757568-1001": "{D6886603-9D2F-4EB2-B667-1971041FA96B}"},
					StatReturn:   &registry.KeyInfo{ValueCount: 1},
					Err:          nil,
				},
			},
			},
			want: checks.NewCheckResult(checks.LoginMethodID, 1, "PIN"),
		},
		{
			name: "Login method is Picture",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{
				{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
					StringValues: map[string]string{"": "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}"},
					StatReturn:   &registry.KeyInfo{ValueCount: 1},
					Err:          nil}}},
			want: checks.NewCheckResult(checks.LoginMethodID, 2, "Picture Logon"),
		},
		{
			name: "Login method is Password",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.LoginMethodID, 4, "Password"),
		},
		{
			name: "Login method is Fingerprint",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{BEC09223-B018-416D-A0AC-523971B639F5}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.LoginMethodID, 8, "Fingerprint"),
		},
		{
			name: "Login method is Facial recognition",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{8AF662BF-65A0-4D0A-A540-A338A999D36F}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.LoginMethodID, 16, "Facial recognition"),
		},
		{
			name: "Login method is Trust signal",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult(checks.LoginMethodID, 32, "Trust signal"),
		},
		{
			name: "Login method is unknown",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "unknown"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckErrorf(checks.LoginMethodID, "error reading value", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.LoginMethod(tt.key)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestRegistryOutput tests that the right registry key is retrieved
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRegistryOutputLoginMethod(t *testing.T) {
	path := "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile"

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	require.NoError(t, err)
	defer func(key registry.Key) {
		err = key.Close()
		require.NoError(t, err)
	}(key)

	valueNames, err := key.ReadValueNames(-1)
	require.NoError(t, err)
	require.NotEmpty(t, valueNames)
}
