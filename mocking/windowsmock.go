package mocking

import (
	"golang.org/x/sys/windows"
)

// WindowsVersion is an interface for reading the Windows version
// You can use RtlGetVersion to give custom Windows version information back
type WindowsVersion interface {
	RtlGetVersion() *windows.OsVersionInfoEx
}

// MockWindowsVersion is a custom implementation of the version information
//
// Fields:
//
// MajorVersion (uint32) - the major version of the Windows OS (10 or 11)
//
// # MinorVersion (uint32) - the minor version of the Windows OS
//
// BuildNumber (uint32) - the build number of the Windows OS
type MockWindowsVersion struct {
	MajorVersion uint32
	MinorVersion uint32
	BuildNumber  uint32
}

// RtlGetVersion with the MockWindowsVersion pointer returns the custom Windows version information
func (m *MockWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return &windows.OsVersionInfoEx{MajorVersion: m.MajorVersion, MinorVersion: m.MinorVersion, BuildNumber: m.BuildNumber}
}

// RealWindowsVersion is the real implementation of the WindowsVersion interface
type RealWindowsVersion struct {
}

// RtlGetVersion with the RealWindowsVersion pointer returns the real Windows version information
func (r *RealWindowsVersion) RtlGetVersion() *windows.OsVersionInfoEx {
	return windows.RtlGetVersion()
}
