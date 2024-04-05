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
	versionString := fmt.Sprintf("%d.%d.%d", versionData.MajorVersion, versionData.MinorVersion, versionData.BuildNumber)

	// Depending on the major Windows version (10 or 11), act accordingly
	switch versionData.MajorVersion {
	case 11:
		if versionData.BuildNumber == newestWin11Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString+"You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString+"There are updates available for Windows 11.")
		}
	case 10:
		if versionData.BuildNumber == newestWin10Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString+"You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString+"There are updates available for Windows 10.")
		}
	default:
		return NewCheckResult(WindowsOutdatedID, 2, versionString+
			"You are using a Windows version which does not have support anymore. "+
			"Consider updating to Windows 10 or Windows 11.")
	}
}
