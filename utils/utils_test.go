package utils

import (
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
	err := CopyFile(src, dst)
	require.NoError(t, err)
}

// TestCopyFileFailNonexistentSource tests the CopyFile function with a nonexistent source file
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCopyFileFailNonexistentSource(T *testing.T) {
	src := "nonexistent.txt"
	dst := "test_copy.txt"
	err := CopyFile(src, dst)
	require.Error(T, err)
}

// TestCopyFileFailNonexistentDestination tests the CopyFile function with a nonexistent destination folder
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCopyFileFailNonexistentDestination(T *testing.T) {
	src := "utils.go"
	dst := "nonexistent/test_copy.txt"
	err := CopyFile(src, dst)
	require.Error(T, err)
	_, err = os.Stat("nonexistent")
	require.Error(T, err)
}

// TestPhishingDomainsReturnsResults ensures the GetPhishingDomains function returns results
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestPhishingDomainsReturnsResults(T *testing.T) {
	domains := GetPhishingDomains()
	require.NotEmpty(T, domains)
}

// TestFirefoxFolder tests the FirefoxFolder function
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
//func TestFirefoxFolder(t *testing.T) {
// 	// This test should not be a part of the test suite, as it is dependent on the user's system
//	// (unless the test suite will be run on a virtual machine)
// 	// It will fail if the user does not have Firefox installed.
//	// It does work properly if you do have it installed.
//	folders, err := FirefoxFolder()
//	require.NoError(T, err)
//	require.NotEmpty(T, folders)
//}

// TestCurrentUserNameReturnsResults ensures the CurrentUserName function returns a result
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestCurrentUsernameReturnsResult(T *testing.T) {
	username, err := CurrentUsername()
	require.NoError(T, err)
	require.NotEmpty(T, username)
}

// TestRemoveDuplicateStrRemovesDuplicates ensures the function works as intended
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestRemoveDuplicateStrRemovesDuplicates(T *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := RemoveDuplicateStr(input)
	require.Equal(T, expected, result)
}

// TestRemoveDuplicateStrEmptyInput ensures the function works as intended with an empty input
//
// Parameters: t *testing.T - The testing framework
//
// Returns: _
func TestRemoveDuplicateStrEmptyInput(T *testing.T) {
	var input []string
	var expected []string
	result := RemoveDuplicateStr(input)
	require.Equal(T, expected, result)
}
