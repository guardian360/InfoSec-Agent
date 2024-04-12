package checks_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TestNetworkSharing is a test function that validates the behavior of the NetworkSharing function.
//
// It executes a series of test cases, each with different inputs, to ensure that the function behaves as expected in various scenarios.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If a test case fails, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the NetworkSharing function correctly identifies the status of network sharing on the system and handles errors as expected.
func TestNetworkSharing(t *testing.T) {
	tests := []struct {
		name     string
		executor commandmock.CommandExecutor
		want     checks.Check
	}{
		{
			name: "Get-NetAdapterBinding command error",
			executor: &commandmock.MockCommandExecutor{Output: "",
				Err: errors.New("error executing command Get-NetAdapterBinding")},
			want: checks.NewCheckErrorf(checks.NetworkSharingID,
				"error executing command Get-NetAdapterBinding",
				errors.New("error executing command Get-NetAdapterBinding")),
		},
		{
			name:     "Network sharing is enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nTrue\r\nTrue\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 0, "Network sharing is enabled"),
		},
		{
			name:     "Network sharing is partially enabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 1, "Network sharing is partially enabled"),
		},
		{
			name:     "Network sharing is disabled",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 2, "Network sharing is disabled"),
		},
		{
			name:     "Network sharing status is unknown",
			executor: &commandmock.MockCommandExecutor{Output: "\r\n\r\n\r\nHelloWorld\r\n\r\n\r\n", Err: nil},
			want:     checks.NewCheckResult(checks.NetworkSharingID, 3, "Network sharing status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.NetworkSharing(tt.executor)
			require.Equal(t, tt.want, got)
		})
	}
}
