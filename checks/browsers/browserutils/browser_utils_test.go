package browserutils_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browserutils"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
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
	file := &mocking.FileMock{IsOpen: true, Err: nil}
	err := browserutils.CloseFile(file)
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
	file := &mocking.FileMock{IsOpen: true, Err: nil}
	err := file.Close()
	require.NoError(t, err)
	err = browserutils.CloseFile(file)
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
	var file *mocking.FileMock
	err := browserutils.CloseFile(file)
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
	domains := browserutils.GetPhishingDomains()
	require.NotEmpty(t, domains)
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
	mockSource := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	mockDestination := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	err := browserutils.CopyFile("", "", mockSource, mockDestination)
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
	mockSource := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: os.ErrNotExist}
	mockDestination := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	err := browserutils.CopyFile("", "", mockSource, mockDestination)
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
	mockSource := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	mockDestination := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: os.ErrNotExist}
	err := browserutils.CopyFile("", "", mockSource, mockDestination)
	require.Error(t, err)
}
