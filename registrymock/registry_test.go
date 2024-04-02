package registrymock_test

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

// TestOpenRegistryKeyValidInput tests the OpenRegistryKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyValidInput(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.LocalMachine, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer registrymock.CloseRegistryKey(key)
	require.NoError(t, err)
	require.NotNil(t, key)
}

// TestOpenRegistryKeyInvalidKey tests the OpenRegistryKey function with an invalid key
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidKey(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.NewRegistryKeyWrapper(registry.Key(0x0)), "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer registrymock.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, registrymock.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestOpenRegistryKeyInvalidPath tests the OpenRegistryKey function with an invalid path
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidPath(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\nonexistent")
	defer registrymock.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, registrymock.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestCloseRegistryKeyValidInput tests the CloseRegistryKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyValidInput(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.LocalMachine, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	require.NoError(t, err)
	require.NotNil(t, key)
	registrymock.CloseRegistryKey(key)
}

// TestCloseRegistryKeyInvalidKey tests the CloseRegistryKey function with an invalid key
// Because a potential error is only logged, we test if the function panics
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyInvalidKey(t *testing.T) {
	key := registrymock.NewRegistryKeyWrapper(registry.Key(0x0))
	registrymock.CloseRegistryKey(key)
}

// TestFindEntriesValidInput tests the FindEntries function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestFindEntriesValidInput(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.CurrentUser,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	defer registrymock.CloseRegistryKey(key)
	require.NoError(t, err)
	require.NotNil(t, key)
	entries, err := key.ReadValueNames(0)
	require.NoError(t, err)
	result := registrymock.FindEntries(entries, key)
	require.NotEmpty(t, result)
}

// TestFindEntriesInvalidInput tests the FindEntries function with invalid (empty) input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestFindEntriesInvalidInput(t *testing.T) {
	key := registry.Key(0x0)
	var entries []string
	elements := registrymock.FindEntries(entries, registrymock.NewRegistryKeyWrapper(key))
	require.Empty(t, elements)
}

// TestCheckKeyValidInput tests the CheckKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyValidInput(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer registrymock.CloseRegistryKey(key)
	require.NoError(t, err)

	// This test might fail if the element does not exist on the machine where the test is run
	val := registrymock.CheckKey(key, "ProductName")
	require.NotEqual(t, "-1", val)
}

// TestCheckKeyInvalidKey tests the CheckKey function with an invalid key
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidKey(t *testing.T) {
	key := registry.Key(0x0)
	val := registrymock.CheckKey(registrymock.NewRegistryKeyWrapper(key), "ProductName")
	require.Equal(t, "-1", val)
}

// TestCheckKeyInvalidElement tests the CheckKey function with an invalid element
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidElement(t *testing.T) {
	key, err := registrymock.OpenRegistryKey(registrymock.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer registrymock.CloseRegistryKey(key)
	require.NoError(t, err)
	val := registrymock.CheckKey(key, "Nonexistent")
	require.Equal(t, "-1", val)
}
