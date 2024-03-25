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
			key: &registrymock.MockRegistryKey{
				StringValue:          "{D6886603-9D2F-4EB2-B667-1971041FA96B}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "PIN"),
		},
		{
			name: "Login method is Picture",
			key: &registrymock.MockRegistryKey{
				StringValue:          "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "Picture Logon"),
		},
		{
			name: "Login method is Password",
			key: &registrymock.MockRegistryKey{
				StringValue:          "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "Password"),
		},
		{
			name: "Login method is Fingerprint",
			key: &registrymock.MockRegistryKey{
				StringValue:          "{BEC09223-B018-416D-A0AC-523971B639F5}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "Fingerprint"),
		},
		{
			name: "Login method is Facial recognition",
			key: &registrymock.MockRegistryKey{
				StringValue:          "{8AF662BF-65A0-4D0A-A540-A338A999D36F}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "Facial recognition"),
		},
		{
			name: "Login method is Trust signal",
			key: &registrymock.MockRegistryKey{
				StringValue:          "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
			want: checks.NewCheckResult("LoginMethod", "Trust signal"),
		},
		{
			name: "Login method is unknown",
			key: &registrymock.MockRegistryKey{
				StringValue:          "unknown",
				StatReturn:           &registry.KeyInfo{ValueCount: 1},
				ReadValueNamesReturn: []string{""}, Err: nil},
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