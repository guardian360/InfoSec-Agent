package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"golang.org/x/sys/windows/registry"
	"reflect"
	"testing"
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
			want: checks.NewCheckResult("LoginMethod", "PIN"),
		},
		{
			name: "Login method is Picture",
			key: &registrymock.MockRegistryKey{SubKeys: []registrymock.MockRegistryKey{
				{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
					StringValues: map[string]string{"": "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}"},
					StatReturn:   &registry.KeyInfo{ValueCount: 1},
					Err:          nil}}},
			want: checks.NewCheckResult("LoginMethod", "Picture Logon"),
		},
		{
			name: "Login method is Password",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"S-1-5-21-174387295-310396468-1212757568-500": "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult("LoginMethod", "Password"),
		},
		{
			name: "Login method is Fingerprint",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{BEC09223-B018-416D-A0AC-523971B639F5}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult("LoginMethod", "Fingerprint"),
		},
		{
			name: "Login method is Facial recognition",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{8AF662BF-65A0-4D0A-A540-A338A999D36F}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult("LoginMethod", "Facial recognition"),
		},
		{
			name: "Login method is Trust signal",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult("LoginMethod", "Trust signal"),
		},
		{
			name: "Login method is unknown",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI\\UserTile",
						StringValues: map[string]string{"": "unknown"},
						StatReturn:   &registry.KeyInfo{ValueCount: 1}, Err: nil}}},
			want: checks.NewCheckResult("LoginMethod"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.LoginMethod(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRegistryOutput tests that the right registry key is retrieved
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRegistryOutputLoginMethod(t *testing.T) {
	tests := []struct {
		name string
		path string
		want registry.KeyInfo
	}{
		{
			name: "UserTitle key",
			path: `SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\LogonUI\UserTile`,
			// We do not assign a value to lastWritetime, since it can be overwritten by the system
			want: registry.KeyInfo{
				SubKeyCount:     0,
				MaxSubKeyLen:    0,
				ValueCount:      2,
				MaxValueNameLen: 46,
				MaxValueLen:     78},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := registry.OpenKey(registry.LOCAL_MACHINE, tt.path, registry.QUERY_VALUE)
			if err != nil {
				t.Fail()
			}
			defer func(key registry.Key) {
				err := key.Close()
				if err != nil {
					t.Fail()
				}
			}(key)
			res, err := key.Stat()
			if !reflect.DeepEqual(res.SubKeyCount, tt.want.SubKeyCount) ||
				!reflect.DeepEqual(res.MaxSubKeyLen, tt.want.MaxSubKeyLen) ||
				!reflect.DeepEqual(res.ValueCount, tt.want.ValueCount) ||
				!reflect.DeepEqual(res.MaxValueNameLen, tt.want.MaxValueNameLen) ||
				!reflect.DeepEqual(res.MaxValueLen, tt.want.MaxValueLen) {
				t.Errorf("Registry key info = %v, want %v", res, tt.want)
			}
		})
	}
}
