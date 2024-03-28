package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"golang.org/x/sys/windows/registry"
	"reflect"
	"testing"
)

// TestRemoteDesktopCheck tests the RemoteDesktopCheck function on (in)valid input
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRemoteDesktopCheck(t *testing.T) {
	tests := []struct {
		name string
		key  registrymock.RegistryKey
		want checks.Check
	}{
		{
			name: "Remote Desktop enabled",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 0}, Err: nil},
				},
			},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is enabled"),
		},
		{
			name: "Remote Desktop disabled",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 1}, Err: nil},
				},
			},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is disabled"),
		},
		{
			name: "Unknown status",
			key: &registrymock.MockRegistryKey{
				SubKeys: []registrymock.MockRegistryKey{
					{KeyName: "System\\CurrentControlSet\\Control\\Terminal Server",
						IntegerValues: map[string]uint64{"fDenyTSConnections": 3}, Err: nil},
				},
			},
			want: checks.NewCheckResult("RemoteDesktop", "Remote Desktop is disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.RemoteDesktopCheck(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteDesktopCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRegistryOutput ensures the registry output has the expected format
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestRegistryOutput(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  registry.KeyInfo
		want2 []uint64
	}{
		{
			name: "Terminal Server key",
			path: `System\CurrentControlSet\Control\Terminal Server`,
			// We do not assign a value to lastWritetime, since it can be overwritten by the system
			want: registry.KeyInfo{
				SubKeyCount:     11,
				MaxSubKeyLen:    24,
				ValueCount:      14,
				MaxValueNameLen: 21,
				MaxValueLen:     64},
			want2: []uint64{0, 1},
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
			val, _, err := key.GetIntegerValue("fDenyTSConnections")
			if !reflect.DeepEqual(val, tt.want2[0]) &&
				!reflect.DeepEqual(val, tt.want2[1]) {
				t.Errorf("Key integer value= %v, want %v or %v", val, tt.want2[0], tt.want2[1])
			}
		})
	}
}
