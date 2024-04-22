package firefox_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSearchEngineFirefox_WithInvalidDirectory(t *testing.T) {
	// Mock the FirefoxFolder function to return an invalid directory
	Profilefinder = utils.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"/invalid/directory"}, nil
		},
	}

	check := firefox.SearchEngineFirefox(Profilefinder)
	assert.Nil(t, check.Result)
	assert.NotNil(t, check.Error)
}

func TestResults_WithGoogleEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":""`)
	result := firefox.Results(data)
	assert.Equal(t, "Google", result)
}

func TestResults_WithKnownEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":"ddg@search.mozilla.org"`)
	result := firefox.Results(data)
	assert.Equal(t, "ddg@search.mozilla.org", result)
}

func TestResults_WithUnknownEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":"unknown@search.mozilla.org"`)
	result := firefox.Results(data)
	assert.Equal(t, "Other Search Engine", result)
}

func TestOpenAndStatFile_WithValidFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
	}
	defer os.Remove(tempFile.Name())

	// Call the function with the temporary file
	file, size, err := firefox.OpenAndStatFile(tempFile.Name())

	// Assert that no error occurred, the returned file is not nil, and the size is correct
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, int64(0), size)
}

func TestOpenAndStatFile_WithNonExistentFile(t *testing.T) {
	// Call the function with a non-existent file
	file, size, err := firefox.OpenAndStatFile("/non/existent/file")

	// Assert that an error occurred, the returned file is nil, and the size is 0
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Equal(t, int64(0), size)
}

func TestOpenAndStatFile_WithDirectory(t *testing.T) {
	// Call the function with a directory
	file, size, err := firefox.OpenAndStatFile("/")

	// Assert that an error occurred, the returned file is nil, and the size is 0
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Equal(t, int64(0), size)
}
