// Package windowsmock provides a mock implementation of the WindowsVersion interface
package windowsmock

import (
	"golang.org/x/sys/windows"
)

// TODO: fix this comment once copilot decides to cooperate

// WindowsVersion is an interface for reading the Windows version
// You can use RtlGetVersion to give custom Windows version information back
type WindowsVersion interface {
	RtlGetVersion() *windows.OsVersionInfoEx
}

// MockWindowsVersion is a struct that simulates the version information of a Windows operating system.
//
// It is designed to be used in testing scenarios where control over the reported version of Windows is required.
//
// Fields:
//   - MajorVersion (uint32): Represents the major version of the Windows OS. For example, for Windows 10 or 11, this would be 10 or 11 respectively.
//   - MinorVersion (uint32): Represents the minor version of the Windows OS. This is typically used for minor updates or revisions.
//   - BuildNumber (uint32): Represents the build number of the Windows OS. This is typically incremented with each build release by Microsoft.
//
// The MockWindowsVersion struct implements the WindowsVersion interface by providing a RtlGetVersion method that returns a custom OsVersionInfoEx object.
type MockWindowsVersion struct {
	MajorVersion uint32
	MinorVersion uint32
	BuildNumber  uint32
}

// TODO: fix this comment once copilot decides to cooperate

// RtlGetVersion with the MockWindowsVersion pointer returns the custom Windows version information
func (m *MockWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return &windows.OsVersionInfoEx{MajorVersion: m.MajorVersion, MinorVersion: m.MinorVersion, BuildNumber: m.BuildNumber}
}

// TODO: fix this comment once copilot decides to cooperate

// RealWindowsVersion is the real implementation of the WindowsVersion interface
type RealWindowsVersion struct {
}

// TODO: fix this comment once copilot decides to cooperate

// RtlGetVersion with the RealWindowsVersion pointer returns the real Windows version information
func (r *RealWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return windows.RtlGetVersion()
}
