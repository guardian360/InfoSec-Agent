package utils_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/filemock"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCopyFileSuccess tests the CopyFile function with a valid source and destination file
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
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

// TestCopyFileFailNonexistentSource tests the CopyFile function with a nonexistent source file
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCopyFileFailNonexistentSource(t *testing.T) {
	src := "nonexistent.txt"
	dst := "test_copy.txt"
	err := utils.CopyFile(src, dst)
	require.Error(t, err)
}

// TestCopyFileFailNonexistentDestination tests the CopyFile function with a nonexistent destination folder
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCopyFileFailNonexistentDestination(t *testing.T) {
	src := "utils.go"
	dst := "nonexistent/test_copy.txt"
	err := utils.CopyFile(src, dst)
	require.Error(t, err)
	_, err = os.Stat("nonexistent")
	require.Error(t, err)
}

// TestPhishingDomainsReturnsResults ensures the GetPhishingDomains function returns results
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestPhishingDomainsReturnsResults(t *testing.T) {
	domains := utils.GetPhishingDomains()
	require.NotEmpty(t, domains)
}

// TestFirefoxFolder tests the FirefoxFolder function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
// func TestFirefoxFolder(t *testing.T) {
// 	// This test should not be a part of the test suite, as it is dependent on the user's system
//	// (unless the test suite will be run on a virtual machine)
// 	// It will fail if the user does not have Firefox installed.
//	// It does work properly if you do have it installed.
//	folders, err := FirefoxFolder()
//	require.NoError(t, err)
//	require.NotEmpty(t, folders)
//}

// TestCurrentUserNameReturnsResults ensures the CurrentUserName function returns a result
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCurrentUsernameReturnsResult(t *testing.T) {
	username, err := utils.CurrentUsername()
	require.NoError(t, err)
	require.NotEmpty(t, username)
}

// TestRemoveDuplicateStrRemovesDuplicates ensures the function works as intended
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestRemoveDuplicateStrRemovesDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := utils.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}

// TestRemoveDuplicateStrEmptyInput ensures the function works as intended with an empty input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestRemoveDuplicateStrEmptyInput(t *testing.T) {
	var input []string
	var expected []string
	result := utils.RemoveDuplicateStr(input)
	require.Equal(t, expected, result)
}

// TestCloseFileNoError ensures the CloseFile function works as intended
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseFileNoError(t *testing.T) {
	file := &filemock.FileMock{Open: true, Err: nil}
	err := utils.CloseFile(file)
	require.NoError(t, err)
}

// TestCloseFileWhenFileWasAlreadyClosed ensures the CloseFile function works as intended when the file was already closed
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseFileWhenFileWasAlreadyClosed(t *testing.T) {
	file := &filemock.FileMock{Open: true, Err: nil}
	err := file.Close()
	require.NoError(t, err)
	err = utils.CloseFile(file)
	require.Error(t, err)
}

// TestCloseFileWhenFileIsNil ensures the CloseFile function works as intended when the file is nil
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCloseFileWhenFileIsNil(t *testing.T) {
	var file *filemock.FileMock
	err := utils.CloseFile(file)
	require.Error(t, err)
}
