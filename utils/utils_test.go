package utils_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/filemock"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

// TestCopyFileSuccess validates the behavior of the CopyFile function when provided with a valid source and destination file.
//
// This test function creates a source file and a destination file, then calls the CopyFile function with these files as arguments.
// It asserts that no error is returned by the CopyFile function, indicating that the file was successfully copied from the source to the destination.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileSuccess(t *testing.T) {
	src := "utils.go"
	dst := "utils_copy.go"
	defer func() {
		err := os.Remove("utils_copy.go")
		require.NoError(t, err)
	}()
	err := utils.CopyFile(src, dst)
	require.NoError(t, err)
}

// TestCopyFileFailNonexistentSource validates the behavior of the CopyFile function when provided with a nonexistent source file.
//
// This test function calls the CopyFile function with a source file path that does not exist and a valid destination path.
// It asserts that an error is returned by the CopyFile function, indicating that the file could not be copied from the nonexistent source.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileFailNonexistentSource(t *testing.T) {
	src := "nonexistent.txt"
	dst := "test_copy.txt"
	err := utils.CopyFile(src, dst)
	require.Error(t, err)
}

// TestCopyFileFailNonexistentDestination validates the behavior of the CopyFile function when provided with a nonexistent destination folder.
//
// This test function calls the CopyFile function with a valid source file and a destination path that does not exist.
// It asserts that an error is returned by the CopyFile function, indicating that the file could not be copied to the nonexistent destination.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileFailNonexistentDestination(t *testing.T) {
	src := "utils.go"
	dst := "nonexistent/test_copy.txt"
	err := utils.CopyFile(src, dst)
	require.Error(t, err)
	_, err = os.Stat("nonexistent")
	require.Error(t, err)
}

// TestPhishingDomainsReturnsResults validates the behavior of the GetPhishingDomains function by ensuring it returns results.
//
// This test function calls the GetPhishingDomains function and asserts that the returned slice is not empty.
// It is designed to ensure that the GetPhishingDomains function correctly retrieves a list of phishing domains.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestPhishingDomainsReturnsResults(t *testing.T) {
	domains := utils.GetPhishingDomains()
	require.NotEmpty(t, domains)
}

// TestFirefoxFolder verifies the functionality of the FirefoxFolder function.
//
// This test function is designed to ensure that the FirefoxFolder function correctly retrieves the Firefox folder paths. It is not included in the test suite due to its dependency on the user's system and the presence of Firefox installation.
//
// Parameters: t *testing.T - The testing tool used for managing and reporting test cases.
//
// Returns: None. The function does not return any value.
// func TestFirefoxFolder(t *testing.T) {
// 	// This test should not be a part of the test suite, as it is dependent on the user's system
//	// (unless the test suite will be run on a virtual machine)
// 	// It will fail if the user does not have Firefox installed.
//	// It does work properly if you do have it installed.
//	folders, err := FirefoxFolder()
//	require.NoError(t, err)
//	require.NotEmpty(t, folders)
//}

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
	username, err := utils.CurrentUsername()
	require.NoError(t, err)
	require.NotEmpty(t, username)
}

// TestRemoveDuplicateStrRemovesDuplicates validates the functionality of the RemoveDuplicateStr function by ensuring it correctly removes duplicate string values from a given slice.
//
// This test function creates a slice with duplicate string values and passes it to the RemoveDuplicateStr function.
// It asserts that the returned slice contains only the unique string values from the input slice, in the order of their first occurrence.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestRemoveDuplicateStrRemovesDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := utils.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}

// TestRemoveDuplicateStrEmptyInput validates the behavior of the RemoveDuplicateStr function when provided with an empty input.
//
// This test function creates an empty string slice and passes it to the RemoveDuplicateStr function.
// It asserts that the returned slice is also empty, confirming that the function handles empty input correctly.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestRemoveDuplicateStrEmptyInput(t *testing.T) {
	var input []string
	var expected []string
	result := utils.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}

// TestCloseFileNoError validates the CloseFile function's ability to close a file without errors.
//
// This test function creates a mock file, then calls the CloseFile function with this file as an argument.
// It asserts that no error is returned by the CloseFile function, indicating that the file was successfully closed.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileNoError(t *testing.T) {
	file := &filemock.FileMock{Open: true, Err: nil}
	err := utils.CloseFile(file)
	require.NoError(t, err)
}

// TestCloseFileWhenFileWasAlreadyClosed verifies the behavior of the CloseFile function when the file has already been closed.
//
// This test function asserts that the CloseFile function returns an error when it is called with a file that has already been closed.
// It is designed to ensure that the CloseFile function handles this edge case correctly, contributing to the robustness of the file handling process.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileWhenFileWasAlreadyClosed(t *testing.T) {
	file := &filemock.FileMock{Open: true, Err: nil}
	err := file.Close()
	require.NoError(t, err)
	err = utils.CloseFile(file)
	require.Error(t, err)
}

// TestCloseFileWhenFileIsNil verifies the behavior of the CloseFile function when the provided file is nil.
//
// This test function asserts that the CloseFile function returns an error when it is called with a nil file.
// It is designed to ensure that the CloseFile function handles this edge case correctly, contributing to the robustness of the file handling process.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileWhenFileIsNil(t *testing.T) {
	var file *filemock.FileMock
	err := utils.CloseFile(file)
	require.Error(t, err)
}
