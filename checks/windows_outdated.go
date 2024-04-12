package checks

import (
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TODO: The "newest build" is done manually below and should be done automatic.
// So preferably it would read the newest build without us changing the values.
const (
	newestWin10Build uint32 = 19045 // Version 22H2 (2022 update)
	newestWin11Build uint32 = 22631 // Version 23H2 (2023 update)
)

// WindowsOutdated is a function that checks if the currently installed Windows version is outdated.
//
// Parameters:
//   - mockOS mocking.WindowsVersion: A mock object for retrieving the Windows version information.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the Windows version is up-to-date or if updates are available.
//
// The function works by retrieving the Windows version information using the provided mock object. It then compares the build number of the installed Windows version with the build numbers of the latest Windows 10 and Windows 11 versions. If the installed version's build number matches the latest build number for its major version (10 or 11), the function returns a message indicating that the Windows version is up-to-date. If the build number does not match, the function returns a message indicating that updates are available. If the major version is neither 10 nor 11, the function returns a message suggesting to update to Windows 10 or Windows 11.
func WindowsOutdated(mockOS mocking.WindowsVersion) Check {
	versionData := mockOS.RtlGetVersion()
	versionString := fmt.Sprintf("%d.%d.%d", versionData.MajorVersion, versionData.MinorVersion, versionData.BuildNumber)

	// Depending on the major Windows version (10 or 11), act accordingly
	switch versionData.MajorVersion {
	case 11:
		if versionData.BuildNumber == newestWin11Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString, "You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString, "There are updates available for Windows 11.")
		}
	case 10:
		if versionData.BuildNumber == newestWin10Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString, "You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString, "There are updates available for Windows 10.")
		}
	default:
		return NewCheckResult(WindowsOutdatedID, 2, versionString,
			"You are using a Windows version which does not have support anymore. "+
				"Consider updating to Windows 10 or Windows 11.")
	}
}
