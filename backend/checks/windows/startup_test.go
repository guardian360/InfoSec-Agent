package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"golang.org/x/sys/windows/registry"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestStartup is a function that tests the behavior of the Startup function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the Startup function with different scenarios. It uses a mock implementation of the RegistryKey interface to simulate the behavior of the registry keys where startup programs can be located. Each test case checks if the Startup function correctly identifies the presence of startup programs based on the simulated registry key values. The function asserts that the returned Check instance contains the expected results.
func TestStartup(t *testing.T) {
	tests := []struct {
		name string
		key1 mocking.RegistryKey
		key2 mocking.RegistryKey
		key3 mocking.RegistryKey
		want checks.Check
	}{{
		name: "No startup programs found",
		key1: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"}}},
		key2: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"}}},
		key3: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName: "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run32"}}},
		want: checks.NewCheckResult(checks.StartupID, 0),
	}, {
		name: "Startup programs found",
		key1: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
			BinaryValues: map[string][]byte{"MockProgram": {1, 0, 0, 0, 0, 0, 0}}, Err: nil}}},
		key2: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
			BinaryValues: map[string][]byte{"MockProgram2": {0, 0, 0, 0, 1, 0, 0}}, Err: nil}}},
		key3: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run32",
			BinaryValues: map[string][]byte{"MockProgram3": {0, 0, 0, 0, 0, 1, 0}}, Err: nil}}},
		want: checks.NewCheckResult(checks.StartupID, 1, "MockProgram"),
	},
		{
			name: "Error opening registry keys",
			key1: &mocking.MockRegistryKey{},
			key2: &mocking.MockRegistryKey{},
			key3: &mocking.MockRegistryKey{},
			want: checks.NewCheckError(checks.StartupID, errors.New("error opening registry keys")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.Startup(tt.key1, tt.key2, tt.key3)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestFindEntriesInvalidInput is a test function that validates the behavior of the FindEntries function when provided with invalid (empty) input.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the FindEntries function behaves as expected when provided with an empty list of entries and an invalid registry key. Specifically, it checks that the function returns an empty list of entries. If the FindEntries function does not behave as expected, this test function will cause the test run to fail.
func TestFindEntriesInvalidInput(t *testing.T) {
	key := registry.Key(0x0)
	var entries []string
	elements := windows.FindEntries(entries, mocking.NewRegistryKeyWrapper(key))
	require.Empty(t, elements)
}

func TestCheckAllZero(t *testing.T) {
	entries := []byte{0, 0, 0, 0}
	result := windows.CheckAllZero(entries)
	require.True(t, result)
	entries = []byte{0, 0, 0, 1}
	result = windows.CheckAllZero(entries)
	require.False(t, result)
	entries = []byte{}
	result = windows.CheckAllZero(entries)
	require.True(t, result)
}
