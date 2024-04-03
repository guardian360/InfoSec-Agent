package checks_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/windowsmock"
)

// TestWindowsOutdated tests the WindowsOutdated function with different (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestWindowsOutdated(t *testing.T) {
	tests := []struct {
		name   string
		mockOS *windowsmock.MockWindowsVersion
		want   checks.Check
	}{
		{
			name:   "Windows 11 up-to-date",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 11, MinorVersion: 0, BuildNumber: 22631},
			want: checks.NewCheckResult("Windows Version", "11.0.22631",
				"You are currently up to date."),
		},
		{
			name:   "Windows 11 outdated",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 11, MinorVersion: 0, BuildNumber: 22630},
			want: checks.NewCheckResult("Windows Version", "11.0.22630",
				"There are updates available for Windows 11."),
		},
		{
			name:   "Windows 10 up-to-date",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 10, MinorVersion: 0, BuildNumber: 19045},
			want: checks.NewCheckResult("Windows Version", "10.0.19045",
				"You are currently up to date."),
		},
		{
			name:   "Windows 10 outdated",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 10, MinorVersion: 0, BuildNumber: 19044},
			want: checks.NewCheckResult("Windows Version", "10.0.19044",
				"There are updates available for Windows 10."),
		},
		{
			name:   "Unsupported Windows version",
			mockOS: &windowsmock.MockWindowsVersion{MajorVersion: 9, MinorVersion: 0, BuildNumber: 0},
			want: checks.NewCheckResult("Windows Version", "9.0.0",
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
