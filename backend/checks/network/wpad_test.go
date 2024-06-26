package network_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWPADEnabled(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.MockRegistryKey
		want checks.Check
		err  bool
	}{
		{
			name: "WPAD enabled",
			key: mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Wpad"},
				}, Err: nil},
			want: checks.NewCheckResult(checks.WPADID, 1),
		},
		{
			name: "WPAD disabled",
			key: mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Wpad",
						IntegerValues: map[string]uint64{"WpadOverride": 1},
					},
				}, Err: nil},
			want: checks.NewCheckResult(checks.WPADID, 0),
		},
		{
			name: "Error opening key",
			key:  mocking.MockRegistryKey{},
			want: checks.NewCheckError(checks.WPADID, errors.New("error")),
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network.WPADEnabled(&tt.key)
			if tt.err {
				require.Equal(t, -1, got.ResultID)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
