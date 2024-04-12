package mocking_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestMain(m *testing.M) {
	// Setup
	logger.SetupTests()

	// Run the tests
	code := m.Run()

	// Teardown

	// Exit with the code returned from the tests
	os.Exit(code)
}

// TestOpenRegistryKeyValidInput tests the OpenRegistryKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyValidInput(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer mocking.CloseRegistryKey(key)
	require.NoError(t, err)
	require.NotNil(t, key)
}

// TestOpenRegistryKeyInvalidKey tests the OpenRegistryKey function with an invalid key
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidKey(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.NewRegistryKeyWrapper(registry.Key(0x0)),
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer mocking.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, mocking.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestOpenRegistryKeyInvalidPath tests the OpenRegistryKey function with an invalid path
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestOpenRegistryKeyInvalidPath(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\nonexistent")
	defer mocking.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, mocking.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestCloseRegistryKeyValidInput tests the CloseRegistryKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyValidInput(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	require.NoError(t, err)
	require.NotNil(t, key)
	mocking.CloseRegistryKey(key)
}

// TestCloseRegistryKeyInvalidKey tests the CloseRegistryKey function with an invalid key
// Because a potential error is only logged, we test if the function panics
//
// Parameters: _ *testing.T - The testing framework
//
// Returns: _
func TestCloseRegistryKeyInvalidKey(_ *testing.T) {
	key := mocking.NewRegistryKeyWrapper(registry.Key(0x0))
	mocking.CloseRegistryKey(key)
}

// TestFindEntriesValidInput tests the FindEntries function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestFindEntriesValidInput(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.CurrentUser,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	defer mocking.CloseRegistryKey(key)
	require.NoError(t, err)
	require.NotNil(t, key)
	entries, err := key.ReadValueNames(0)
	require.NoError(t, err)
	result := mocking.FindEntries(entries, key)
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
	elements := mocking.FindEntries(entries, mocking.NewRegistryKeyWrapper(key))
	require.Empty(t, elements)
}

// TestCheckKeyValidInput tests the CheckKey function with valid input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyValidInput(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer mocking.CloseRegistryKey(key)
	require.NoError(t, err)

	// This test might fail if the element does not exist on the machine where the test is run
	val := mocking.CheckKey(key, "ProductName")
	require.NotEqual(t, "-1", val)
}

// TestCheckKeyInvalidKey tests the CheckKey function with an invalid key
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidKey(t *testing.T) {
	key := registry.Key(0x0)
	val := mocking.CheckKey(mocking.NewRegistryKeyWrapper(key), "ProductName")
	require.Equal(t, "-1", val)
}

// TestCheckKeyInvalidElement tests the CheckKey function with an invalid element
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCheckKeyInvalidElement(t *testing.T) {
	key, err := mocking.OpenRegistryKey(mocking.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer mocking.CloseRegistryKey(key)
	require.NoError(t, err)
	val := mocking.CheckKey(key, "Nonexistent")
	require.Equal(t, "-1", val)
}
