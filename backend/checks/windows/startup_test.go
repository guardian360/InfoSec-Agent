package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
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
		name  string
		key1  mocking.RegistryKey
		key2  mocking.RegistryKey
		key3  mocking.RegistryKey
		want  checks.Check
		error bool
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
			name:  "Error opening registry keys",
			key1:  &mocking.MockRegistryKey{},
			key2:  &mocking.MockRegistryKey{},
			key3:  &mocking.MockRegistryKey{},
			want:  checks.NewCheckError(checks.StartupID, errors.New("error")),
			error: true,
		},
		{
			name: "Error reading value names",
			key1: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
				StringValues: map[string]string{"test": "test"}, Err: nil}}},
			key2: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run",
				BinaryValues: map[string][]byte{"MockProgram2": {0, 0, 0, 0, 1, 0, 0}}, Err: nil}}},
			key3: &mocking.MockRegistryKey{SubKeys: []mocking.MockRegistryKey{{
				KeyName:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run32",
				BinaryValues: map[string][]byte{"MockProgram3": {0, 0, 0, 0, 0, 1, 0}}, Err: nil}}},
			want:  checks.NewCheckError(checks.StartupID, errors.New("error")),
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.Startup(tt.key1, tt.key2, tt.key3)
			if tt.error {
				require.Equal(t, -1, got.ResultID)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
