package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
func TestPhishingDomainsReturnsResults(t *testing.T) {
	domains := GetPhishingDomains()
	require.NotEmpty(t, domains)
}
