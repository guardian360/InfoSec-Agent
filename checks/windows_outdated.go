package checks

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// The "newest build" is done manually below and should be done automatic.
// So preferably it would read the newest build without us changing the values.
const (
	newest_win10_build uint32 = 19045 // Version 22H2 (2022 update)
	newest_win11_build uint32 = 22631 // Version 23H2 (2023 update)
)

func Outdated() {
	versionData := windows.RtlGetVersion()
	fmt.Printf("You are running Windows version: %d.%d.%d\n", versionData.MajorVersion, versionData.MinorVersion, versionData.BuildNumber)
	// Depending on the major version, act accordingly
	switch versionData.MajorVersion {
	case 11:
		if versionData.BuildNumber == newest_win11_build {
			fmt.Println("You are currently up to date.")
		} else {
			fmt.Println("There are updates available for Windows 11.")
		}
	case 10:
		if versionData.BuildNumber == newest_win10_build {
			fmt.Println("You are currently up to date.")
		} else {
			fmt.Println("There are updates available for Windows 10.")
		}
	default:
		fmt.Println("You are using a Windows version which does not have support anymore. Consider updating to Windows 10 or Windows 11.")
	}
}
