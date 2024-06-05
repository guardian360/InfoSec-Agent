package windows_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScreenLockEnabled(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
		err  bool
	}{
		{
			name: "Screen lock correctly enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop",
						StringValues: map[string]string{"ScreenSaveActive": "1", "ScreenSaverIsSecure": "1", "ScreenSaveTimeOut": "0"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.ScreenLockID, 0),
		},
		{
			name: "Screen lock not correctly enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop",
						StringValues: map[string]string{"ScreenSaveActive": "0", "ScreenSaverIsSecure": "1", "ScreenSaveTimeOut": "0"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.ScreenLockID, 1),
		},
		{
			name: "Error opening registry key",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: ""}}},
			want: checks.NewCheckErrorf(checks.ScreenLockID, "", nil),
			err:  true,
		},
		{
			name: "Error reading registry value 1",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop"}}},
			want: checks.NewCheckResult(checks.ScreenLockID, 1),
		},
		{
			name: "Error reading registry value 2",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop",
						StringValues: map[string]string{"ScreenSaveActive": "0"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.ScreenLockID, 1),
		},
		{
			name: "Error reading registry value 3",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop",
						StringValues: map[string]string{"ScreenSaveActive": "0", "ScreenSaverIsSecure": "0"}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.ScreenLockID, 1),
		},
		{
			name: "Error parsing interval",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "Control Panel\\Desktop",
						StringValues: map[string]string{"ScreenSaveActive": "0", "ScreenSaverIsSecure": "0", "ScreenSaveTimeOut": "abc"}, Err: nil},
				},
			},
			want: checks.NewCheckErrorf(checks.ScreenLockID, "", nil),
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.ScreenLockEnabled(tt.key)
			if tt.err {
				require.Equal(t, -1, got.ResultID)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
