package mocking

import (
	"golang.org/x/sys/windows"
)

// TODO: Update documentation
// WindowsVersion is an interface that defines a contract for retrieving Windows version information.
// It exposes a single method, RtlGetVersion, which is expected to return an instance of windows.OsVersionInfoEx.
//
// Implementations of this interface can provide either real or mocked Windows version information,
// making it useful in a variety of scenarios such as testing and simulation.
//
// The RtlGetVersion method is designed to return custom Windows version information,
// allowing for flexibility and control over the reported Windows version.
type WindowsVersion interface {
	RtlGetVersion() *windows.OsVersionInfoEx
}

// TODO: Update documentation
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

// TODO: Update documentation
// RtlGetVersion returns a custom Windows version information when called on a MockWindowsVersion instance.
// This method is primarily used in testing scenarios where specific Windows version information is required.
//
// The returned OsVersionInfoEx object contains the following fields:
//   - MajorVersion: Simulates the major version of the Windows OS.
//   - MinorVersion: Simulates the minor version of the Windows OS.
//   - BuildNumber: Simulates the build number of the Windows OS.
//
// Returns:
//   - *windows.OsVersionInfoEx: A pointer to an OsVersionInfoEx object that contains the custom Windows version information.
func (m *MockWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return &windows.OsVersionInfoEx{MajorVersion: m.MajorVersion, MinorVersion: m.MinorVersion, BuildNumber: m.BuildNumber}
}

// TODO: Update documentation
// RealWindowsVersion is a struct that provides the actual version information of the Windows operating system.
//
// It implements the WindowsVersion interface by providing a RtlGetVersion method that returns a genuine windows.OsVersionInfoEx object.
//
// This struct is typically used in production scenarios where accurate Windows version information is required.
type RealWindowsVersion struct {
}

// TODO: Update documentation
// RtlGetVersion retrieves the actual Windows version information when called on a RealWindowsVersion instance.
//
// This method is typically used in production scenarios where accurate Windows version information is required.
//
// The returned OsVersionInfoEx object contains the actual version information of the Windows operating system, including:
//   - MajorVersion: The major version of the Windows OS.
//   - MinorVersion: The minor version of the Windows OS.
//   - BuildNumber: The build number of the Windows OS.
//
// Returns:
//   - *windows.OsVersionInfoEx: A pointer to an OsVersionInfoEx object that contains the actual Windows version information.
func (r *RealWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return windows.RtlGetVersion()
}
