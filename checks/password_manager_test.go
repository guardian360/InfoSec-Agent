package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInstalledProgramsPathsValid tests if the path of installed programs are valid
//
// Parameters: directory string represents the directory to check
//
// Returns: _
func TestInstalledProgramsPathsValid(t *testing.T) {
	testCases := []string{"C:\\Program Files", "C:\\Program Files (x86)"}
	for _, programPath := range testCases {
		programs, err := listInstalledPrograms(programPath)
		assert.Nil(t, err)
		assert.NotNil(t, programs)
	}
}

// TestInstalledProgramsValid tests if the installed programs are valid
//
// Parameters: programPath (string) represents the path in which applications are installed
//
// Returns: _
func TestInstalledProgramsValid(t *testing.T) {
	testCases := []string{"C:\\Program Files", "C:\\Program Files (x86)"}
	for _, programPath := range testCases {
		programs, err := listInstalledPrograms(programPath)
		//Windows defender exist on all windows machines which run windows 10 or 11
		assert.Contains(t, programs, "Windows Defender")
		assert.Nil(t, err)
	}
}

// TestInstalledProgramsInvalid tests if the installed programs are invalid
//
// Parameters: programPath (string) represents the path in which applications are installed
//
// Returns: _
func TestInstalledProgramsInvalid(t *testing.T) {
	testCases := []string{"C:\\nonexistendpathfortest", ""}
	for _, programPath := range testCases {
		programs, err := listInstalledPrograms(programPath)
		assert.Empty(t, programs, "No programs should be found")
		assert.ErrorContainsf(t, err, "The system cannot find the file specified", "Error should contain 'unable to find file' message")
	}
}

// TestNoPasswordManager tests that the 'no password manager' returned when there is none installed
//
// Parameters: _
//
// Returns: _
func TestNoPasswordManager(t *testing.T) {
	check := PasswordManager()
	assert.Contains(t, check.Result, "No password manager found")
}
