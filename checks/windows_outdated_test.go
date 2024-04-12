package checks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/windowsmock"
)

// TestWindowsOutdated is a function that tests the behavior of the WindowsOutdated function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the WindowsOutdated function with different scenarios. It uses a mock implementation of the WindowsVersion interface to simulate the behavior of retrieving the Windows version information. Each test case checks if the WindowsOutdated function correctly identifies whether the Windows version is up-to-date, outdated, or unsupported based on the simulated Windows version information. The function asserts that the returned Check instance contains the expected results.
func TestWindowsOutdated(t *testing.T) {
	tests := []struct {
		name   string
		mockOS *windowsmock.MockWindowsVersion
		want   checks.Check
	}{
		{
			name:   "Windows 11 up-to-date",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 11, MinorVersion: 0, BuildNumber: 22631},
			want: checks.NewCheckResult(checks.WindowsOutdatedID, 0, "11.0.22631",
				"You are currently up to date."),
		},
		{
			name:   "Windows 11 outdated",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 11, MinorVersion: 0, BuildNumber: 22630},
			want: checks.NewCheckResult(checks.WindowsOutdatedID, 1, "11.0.22630",
				"There are updates available for Windows 11."),
		},
		{
			name:   "Windows 10 up-to-date",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 10, MinorVersion: 0, BuildNumber: 19045},
			want: checks.NewCheckResult(checks.WindowsOutdatedID, 0, "10.0.19045",
				"You are currently up to date."),
		},
		{
			name:   "Windows 10 outdated",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 10, MinorVersion: 0, BuildNumber: 19044},
			want: checks.NewCheckResult(checks.WindowsOutdatedID, 1, "10.0.19044",
				"There are updates available for Windows 10."),
		},
		{
			name:   "Unsupported Windows version",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 9, MinorVersion: 0, BuildNumber: 0},
			want: checks.NewCheckResult(checks.WindowsOutdatedID, 2, "9.0.0",
				"You are using a Windows version which does not have support anymore. "+
					"Consider updating to Windows 10 or Windows 11."),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.WindowsOutdated(tt.mockOS)
			require.Equal(t, tt.want, got)
		})
	}
}
