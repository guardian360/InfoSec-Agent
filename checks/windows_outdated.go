package checks

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/windowsmock"
)

// TODO: The "newest build" is done manually below and should be done automatic.
// So preferably it would read the newest build without us changing the values.
const (
	newestWin10Build uint32 = 19045 // Version 22H2 (2022 update)
	newestWin11Build uint32 = 22631 // Version 23H2 (2023 update)
)

// WindowsOutdated checks if the current installed Windows version is outdated
//
// Parameters: _
//
// Returns: If the Windows version is up-to-date or if there are updated available
func WindowsOutdated(mockOS windowsmock.WindowsVersion) Check {
	versionData := mockOS.RtlGetVersion()
	// Prepare the result
	result := NewCheckResult("Windows Version", fmt.Sprintf("%d.%d.%d",
		versionData.MajorVersion, versionData.MinorVersion, versionData.BuildNumber))

	// Depending on the major Windows version (10 or 11), act accordingly
	switch versionData.MajorVersion {
	case 11:
		if versionData.BuildNumber == newestWin11Build {
			result.Result = append(result.Result, "You are currently up to date.")
		} else {
			result.Result = append(result.Result, "There are updates available for Windows 11.")

		}
	case 10:
		if versionData.BuildNumber == newestWin10Build {
			result.Result = append(result.Result, "You are currently up to date.")
		} else {
			result.Result = append(result.Result, "There are updates available for Windows 10.")
		}
	default:
		result.Result = append(result.Result,
			"You are using a Windows version which does not have support anymore. "+
				"Consider updating to Windows 10 or Windows 11.")
	}

	return result
}
