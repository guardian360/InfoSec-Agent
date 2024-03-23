package utils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
	"testing"
)

// TestOpenRegistryKeyValidInput tests the OpenRegistryKey function with valid input
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyValidInput(T *testing.T) {
	key, err := OpenRegistryKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer CloseRegistryKey(key)
	require.NoError(T, err)
	require.NotNil(T, key)
}

// TestOpenRegistryKeyInvalidKey tests the OpenRegistryKey function with an invalid key
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidKey(T *testing.T) {
	key, err := OpenRegistryKey(registry.Key(0x0), "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer CloseRegistryKey(key)
	require.Error(T, err)
	require.Equal(T, key, registry.Key(0x0))
}

// TestOpenRegistryKeyInvalidPath tests the OpenRegistryKey function with an invalid path
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidPath(T *testing.T) {
	key, err := OpenRegistryKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\nonexistent")
	defer CloseRegistryKey(key)
	require.Error(T, err)
	require.Equal(T, key, registry.Key(0x0))
}

// TestCloseRegistryKeyValidInput tests the CloseRegistryKey function with valid input
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyValidInput(T *testing.T) {
	key, err := OpenRegistryKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	require.NoError(T, err)
	require.NotNil(T, key)
	CloseRegistryKey(key)
}

// TestCloseRegistryKeyInvalidKey tests the CloseRegistryKey function with an invalid key
// Because a potential error is only logged, we test if the function panics
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyInvalidKey(T *testing.T) {
	key := registry.Key(0x0)
	CloseRegistryKey(key)
}

// TestFindEntriesValidInput tests the FindEntries function with valid input
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestFindEntriesValidInput(T *testing.T) {
	key, err := OpenRegistryKey(registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	defer CloseRegistryKey(key)
	require.NoError(T, err)
	require.NotNil(T, key)
	entries, err := key.ReadValueNames(0)
	require.NoError(T, err)
	result := FindEntries(entries, key)
	require.NotEmpty(T, result)
}

// TestFindEntriesInvalidInput tests the FindEntries function with invalid (empty) input
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestFindEntriesInvalidInput(T *testing.T) {
	key := registry.Key(0x0)
	var entries []string
	elements := FindEntries(entries, key)
	require.Empty(T, elements)
}

// TestCheckKeyValidInput tests the CheckKey function with valid input
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyValidInput(T *testing.T) {
	key, err := OpenRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer CloseRegistryKey(key)
	require.NoError(T, err)

	// This test might fail if the element does not exist on the machine where the test is run
	val := CheckKey(key, "ProductName")
	require.NotEqual(T, "-1", val)
}

// TestCheckKeyInvalidKey tests the CheckKey function with an invalid key
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidKey(T *testing.T) {
	key := registry.Key(0x0)
	val := CheckKey(key, "ProductName")
	require.Equal(T, "-1", val)
}

// TestCheckKeyInvalidElement tests the CheckKey function with an invalid element
//
// Parameters: T *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidElement(T *testing.T) {
	key, err := OpenRegistryKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer CloseRegistryKey(key)
	require.NoError(T, err)
	val := CheckKey(key, "Nonexistent")
	require.Equal(T, "-1", val)
}
