package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAutomaticLogin(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
		err  bool
	}{
		{
			name: "Automatic log-in enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
						StringValues: map[string]string{"AutoAdminLogon": "1"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AutoLoginID, 1),
		},
		{
			name: "Automatic log-in disabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
						StringValues: map[string]string{"AutoAdminLogon": "0"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AutoLoginID, 0),
		},
		{
			name: "Error opening registry key",
			key:  &mocking.MockRegistryKey{},
			want: checks.NewCheckError(checks.AutoLoginID, errors.New("error opening registry key: key not found")),
			err:  true,
		},
		{
			name: "Error reading AutoAdminLogon value",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
						IntegerValues: map[string]uint64{"AutoAdminLogon2": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AutoLoginID, 0),
		},
		{
			name: "Error converting string",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
						StringValues: map[string]string{"AutoAdminLogon": "test"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AutoLoginID, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.AutomaticLogin(tt.key)
			if tt.err {
				require.Error(t, got.Error)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
