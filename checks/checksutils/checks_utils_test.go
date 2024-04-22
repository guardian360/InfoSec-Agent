package checksutils_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/checksutils"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
	username, err := checksutils.CurrentUsername()
	require.NoError(t, err)
	require.NotEmpty(t, username)
}
