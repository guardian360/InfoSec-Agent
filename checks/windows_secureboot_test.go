package checks_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/RegistryKey"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"reflect"
	"testing"
)

func TestSecureBoot(t *testing.T) {
	tests := []struct {
		name string
		key  RegistryKey.RegistryKey
		want checks.Check
	}{
		// TODO: Add test cases.
		{
			name: "SecureBootError",
			key:  &RegistryKey.MockRegistryKey{StringValue: "", BinaryValue: nil, IntegerValue: 5, Err: errors.New("error")},
			want: checks.NewCheckError("SecureBoot", errors.New("error")),
		},
		{
			name: "SecureBootEnabled",
			key:  &RegistryKey.MockRegistryKey{StringValue: "UEFISecureBootEnabled", BinaryValue: nil, IntegerValue: 1, Err: nil},
			want: checks.NewCheckResult("SecureBoot", "Secure boot is enabled"),
		},
		{
			name: "SecureBootDisabled",
			key:  &RegistryKey.MockRegistryKey{StringValue: "UEFISecureBootEnabled", BinaryValue: nil, IntegerValue: 0, Err: nil},
			want: checks.NewCheckResult("SecureBoot", "Secure boot is disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.SecureBoot(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureBoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
