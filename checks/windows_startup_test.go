package checks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
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
		want: checks.NewCheckResult(checks.StartupID, 0, "No startup programs found"),
	}, {
		name: "Startup programs found",
		key1: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
			BinaryValues: map[string][]byte{"MockProgram": {1, 2, 3, 4, 0, 0, 0}}, Err: nil}}},
		key2: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
			BinaryValues: map[string][]byte{"MockProgram2": {0, 0, 0, 0, 1, 0, 0}}, Err: nil}}},
		key3: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
			KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run32",
			BinaryValues: map[string][]byte{"MockProgram3": {0, 0, 0, 0, 0, 1, 0}}, Err: nil}}},
		want: checks.NewCheckResult(checks.StartupID, 1, "MockProgram"),
	}} /*,{
		name: "Error finding startup programs",
		key1:
		key2:
		key3:
		want:
	}}*/

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.Startup(tt.key1, tt.key2, tt.key3)
			require.Equal(t, tt.want, got)
		})
	}
}
