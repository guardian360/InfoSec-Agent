package checks_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
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

// TestOpenRegistryKeyValidInput validates the functionality of the OpenRegistryKey function when provided with valid input.
//
// Parameter:
//   - t *testing.T: The testing framework used to run the test.
//
// This function does not return any values. It uses the testing framework to assert that the OpenRegistryKey function behaves as expected when provided with a valid registry key and path. If the OpenRegistryKey function does not behave as expected, this test function will cause the test run to fail.
func TestOpenRegistryKeyValidInput(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer checks.CloseRegistryKey(key)
	require.NoError(t, err)
	require.NotNil(t, key)
}

// TestOpenRegistryKeyInvalidKey is a test function that verifies the behavior of the OpenRegistryKey function when provided with an invalid registry key.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the OpenRegistryKey function behaves as expected when provided with an invalid registry key. Specifically, it checks that the function returns an error and that the returned key is equivalent to the invalid input key. If the OpenRegistryKey function does not behave as expected, this test function will cause the test run to fail.
func TestOpenRegistryKeyInvalidKey(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.NewRegistryKeyWrapper(registry.Key(0x0)),
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	defer checks.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, mocking.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestOpenRegistryKeyInvalidPath is a test function that validates the behavior of the OpenRegistryKey function when provided with an invalid path.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the OpenRegistryKey function behaves as expected when provided with an invalid path. Specifically, it checks that the function returns an error and that the returned key is equivalent to a null key. If the OpenRegistryKey function does not behave as expected, this test function will cause the test run to fail.
func TestOpenRegistryKeyInvalidPath(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\nonexistent")
	defer checks.CloseRegistryKey(key)
	require.Error(t, err)
	require.Equal(t, key, mocking.NewRegistryKeyWrapper(registry.Key(0x0)))
}

// TestCloseRegistryKeyValidInput is a test function that validates the behavior of the CloseRegistryKey function when provided with a valid registry key.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the CloseRegistryKey function behaves as expected when provided with a valid registry key. Specifically, it checks that the function does not return an error and that the key is successfully closed. If the CloseRegistryKey function does not behave as expected, this test function will cause the test run to fail.
func TestCloseRegistryKeyValidInput(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.LocalMachine,
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	require.NoError(t, err)
	require.NotNil(t, key)
	checks.CloseRegistryKey(key)
}

// TestCloseRegistryKeyInvalidKey is a test function that verifies the behavior of the CloseRegistryKey function when provided with an invalid registry key.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the CloseRegistryKey function does not panic when provided with an invalid registry key. As the function only logs potential errors and does not return them, this test function checks for unexpected panics as an indication of error handling.
func TestCloseRegistryKeyInvalidKey(_ *testing.T) {
	key := mocking.NewRegistryKeyWrapper(registry.Key(0x0))
	checks.CloseRegistryKey(key)
}

// TestFindEntriesInvalidInput is a test function that validates the behavior of the FindEntries function when provided with invalid (empty) input.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the FindEntries function behaves as expected when provided with an empty list of entries and an invalid registry key. Specifically, it checks that the function returns an empty list of entries. If the FindEntries function does not behave as expected, this test function will cause the test run to fail.
func TestFindEntriesInvalidInput(t *testing.T) {
	key := registry.Key(0x0)
	var entries []string
	elements := checks.FindEntries(entries, mocking.NewRegistryKeyWrapper(key))
	require.Empty(t, elements)
}

// TestCheckKeyValidInput is a test function that validates the behavior of the CheckKey function when provided with valid input.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the CheckKey function behaves as expected when provided with a valid registry key and a valid element. Specifically, it checks that the function returns the correct value of the specified element within the registry key. If the CheckKey function does not behave as expected, this test function will cause the test run to fail.
func TestCheckKeyValidInput(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer checks.CloseRegistryKey(key)
	require.NoError(t, err)

	// This test might fail if the element does not exist on the machine where the test is run
	val := checks.CheckKey(key, "ProductName")
	require.NotEqual(t, "-1", val)
}

// TestCheckKeyInvalidKey is a test function that validates the behavior of the CheckKey function when provided with an invalid registry key.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the CheckKey function behaves as expected when provided with an invalid registry key. Specifically, it checks that the function returns "-1", indicating that the key does not exist or an error occurred while retrieving its value. If the CheckKey function does not behave as expected, this test function will cause the test run to fail.
func TestCheckKeyInvalidKey(t *testing.T) {
	key := registry.Key(0x0)
	val := checks.CheckKey(mocking.NewRegistryKeyWrapper(key), "ProductName")
	require.Equal(t, "-1", val)
}

// TestCheckKeyInvalidElement is a test function that validates the behavior of the CheckKey function when provided with a non-existent element.
//
// Parameter:
//   - t *testing.T: The testing framework instance used to run the test and report the results.
//
// This function does not return any values. It uses the testing framework to assert that the CheckKey function behaves as expected when provided with a valid registry key and a non-existent element. Specifically, it checks that the function returns "-1", indicating that the element does not exist within the registry key. If the CheckKey function does not behave as expected, this test function will cause the test run to fail.
func TestCheckKeyInvalidElement(t *testing.T) {
	key, err := checks.OpenRegistryKey(mocking.LocalMachine,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	defer checks.CloseRegistryKey(key)
	require.NoError(t, err)
	val := checks.CheckKey(key, "Nonexistent")
	require.Equal(t, "-1", val)
}

// TestCurrentUsernameReturnsResult validates the behavior of the CurrentUsername function by ensuring it returns a valid result.
//
// This test function calls the CurrentUsername function and asserts that it returns a non-empty string and no error.
// It is designed to ensure that the CurrentUsername function correctly retrieves the username of the currently logged-in user.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCurrentUsernameReturnsResult(t *testing.T) {
	username, err := checks.CurrentUsername()
	require.NoError(t, err)
	require.NotEmpty(t, username)
}

func TestCheckAllZero(t *testing.T) {
	entries := []byte{0, 0, 0, 0}
	result := checks.CheckAllZero(entries)
	require.True(t, result)
	entries = []byte{0, 0, 0, 1}
	result = checks.CheckAllZero(entries)
	require.False(t, result)
	entries = []byte{}
	result = checks.CheckAllZero(entries)
	require.True(t, result)
}
