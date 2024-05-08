package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.SetupTests()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestAdvertisement(t *testing.T) {
	tests := []struct {
		name string
		key  mocking.RegistryKey
		want checks.Check
		err  bool
	}{
		{
			name: "Advertisement ID enabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\AdvertisingInfo",
						IntegerValues: map[string]uint64{"Enabled": 1}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AdvertisementID, 1),
		},
		{
			name: "Advertisement ID disabled",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\AdvertisingInfo",
						IntegerValues: map[string]uint64{"Enabled": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult(checks.AdvertisementID, 0),
		},
		{
			name: "Error opening registry key",
			key:  &mocking.MockRegistryKey{},
			want: checks.NewCheckError(checks.AdvertisementID, errors.New("error opening registry key: key not found")),
			err:  true,
		},
		{
			name: "Error reading Enabled value",
			key: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{
					{KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\AdvertisingInfo",
						IntegerValues: map[string]uint64{"Enabled2": 0}, Err: nil},
				},
			},
			want: checks.NewCheckError(checks.AdvertisementID, errors.New("")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.Advertisement(tt.key)
			if tt.err {
				require.Error(t, got.Error)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
