package checks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TestNetworkProfileTypes is a test function that validates the behavior of the NetworkProfileTypes function.
//
// It executes a series of test cases, each with different inputs, to ensure that the function behaves as expected in various scenarios.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If a test case fails, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the NetworkProfileTypes function correctly identifies the types of network profiles on the system and handles errors as expected.
func TestNetworkProfileTypes(t *testing.T) {
	tests := []struct {
		name        string
		registryKey mocking.RegistryKey
		want        checks.Check
	}{
		{
			name:        "No network profiles found",
			registryKey: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles"}}},
			want:        checks.NewCheckResult(checks.NetworkProfileTypeID, 0, "No network profiles found"),
		},
		{
			name: "All network profiles are public",
			registryKey: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles",
					SubKeys: []mocking.MockRegistryKey{
						{KeyName: "Profile2", StringValues: map[string]string{"ProfileName": "Profile2"}, BinaryValues: map[string][]byte{"Category": {0}}},
						{KeyName: "Profile1", StringValues: map[string]string{"ProfileName": "Profile1"}, BinaryValues: map[string][]byte{"Category": {0}}},
					}}}},
			want: checks.NewCheckResult(checks.NetworkProfileTypeID, 1, []string{"Network Profile1 is Public", "Network Profile2 is Public"}...),
		},
		{
			name: "All network profiles are private",
			registryKey: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles",
					SubKeys: []mocking.MockRegistryKey{
						{KeyName: "Profile1", StringValues: map[string]string{"ProfileName": "Profile1"}, BinaryValues: map[string][]byte{"Category": {1}}},
						{KeyName: "Profile2", StringValues: map[string]string{"ProfileName": "Profile2"}, BinaryValues: map[string][]byte{"Category": {1}}},
					}}}},
			want: checks.NewCheckResult(checks.NetworkProfileTypeID, 1, []string{"Network Profile1 is Private", "Network Profile2 is Private"}...),
		},
		{
			name: "All network profiles are domain",
			registryKey: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles",
					SubKeys: []mocking.MockRegistryKey{
						{KeyName: "Profile1", StringValues: map[string]string{"ProfileName": "Profile1"}, BinaryValues: map[string][]byte{"Category": {2}}},
						{KeyName: "Profile2", StringValues: map[string]string{"ProfileName": "Profile2"}, BinaryValues: map[string][]byte{"Category": {2}}},
					}}}},
			want: checks.NewCheckResult(checks.NetworkProfileTypeID, 1, []string{"Network Profile1 is Domain", "Network Profile2 is Domain"}...),
		},
		{
			name: "Mixed network profile types",
			registryKey: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles",
					SubKeys: []mocking.MockRegistryKey{
						{KeyName: "Profile1", StringValues: map[string]string{"ProfileName": "Profile1"}, BinaryValues: map[string][]byte{"Category": {0}}},
						{KeyName: "Profile2", StringValues: map[string]string{"ProfileName": "Profile2"}, BinaryValues: map[string][]byte{"Category": {1}}},
						{KeyName: "Profile3", StringValues: map[string]string{"ProfileName": "Profile3"}, BinaryValues: map[string][]byte{"Category": {2}}},
					}}}},
			want: checks.NewCheckResult(checks.NetworkProfileTypeID, 1, []string{"Network Profile1 is Public", "Network Profile2 is Private", "Network Profile3 is Domain"}...),
		},
		{
			name: "Unknown network profile type",
			registryKey: &mocking.MockRegistryKey{
				SubKeys: []mocking.MockRegistryKey{{KeyName: "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles",
					SubKeys: []mocking.MockRegistryKey{
						{KeyName: "Profile1", StringValues: map[string]string{"ProfileName": "Profile1"}, BinaryValues: map[string][]byte{"Category": {3}}},
					}}}},
			want: checks.NewCheckResult(checks.NetworkProfileTypeID, 1, []string{"Network Profile1 is Unknown"}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.NetworkProfileTypes(tt.registryKey)
			require.Equal(t, tt.want, got)
		})
	}
}
