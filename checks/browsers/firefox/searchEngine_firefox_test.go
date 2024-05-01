package firefox_test

import (
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browserutils"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/stretchr/testify/require"
)

func TestSearchEngineFirefox_WithInvalidDirectory(t *testing.T) {
	// Mock the FirefoxFolder function to return an invalid directory
	Profilefinder = browserutils.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"/invalid/directory"}, nil
		},
	}

	check := firefox.SearchEngineFirefox(Profilefinder)
	require.Nil(t, check.Result)
	require.Error(t, check.Error)
}

func TestResults_WithGoogleEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":""`)
	result := firefox.Results(data)
	require.Equal(t, "Google", result)
}

func TestResults_WithKnownEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":"ddg@search.mozilla.org"`)
	result := firefox.Results(data)
	require.Equal(t, "ddg@search.mozilla.org", result)
}

func TestResults_WithUnknownEngine(t *testing.T) {
	data := []byte(`"defaultEngineId":"unknown@search.mozilla.org"`)
	result := firefox.Results(data)
	require.Equal(t, "Other Search Engine", result)
}

func TestOpenAndStatFile_WithValidFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, _ := os.CreateTemp("", "testfile")
	defer os.Remove(tempFile.Name())

	// Call the function with the temporary file
	file, size, err := firefox.OpenAndStatFile(tempFile.Name())

	// Assert that no error occurred, the returned file is not nil, and the size is correct
	require.NoError(t, err)
	require.NotNil(t, file)
	require.Equal(t, int64(0), size)
}

func TestOpenAndStatFile_WithNonExistentFile(t *testing.T) {
	// Call the function with a non-existent file
	file, size, err := firefox.OpenAndStatFile("/non/existent/file")

	// Assert that an error occurred, the returned file is nil, and the size is 0
	require.Error(t, err)
	require.Nil(t, file)
	require.Equal(t, int64(0), size)
}

func TestYourFunction(t *testing.T) {
	// Mock the OpenAndStatFile function
	// Save the original function and defer its restoration
	originalOpenAndStatFile := firefox.OpenAndStatFile
	defer func() { firefox.OpenAndStatFile = originalOpenAndStatFile }()

	// Mock the OpenAndStatFile function
	firefox.OpenAndStatFile = func(_ string) (mocking.File, int64, error) {
		// Return whatever you need for your test
		file := &mocking.FileMock{
			FileName: "testfile",
			IsOpen:   true,
			Buffer:   []byte("1234567864567744545454443343454"),
			Bytes:    26,
			FileInfo: &mocking.FileInfoMock{},
			Err:      nil,
		}
		return file, 26, nil
	}
	expected := checks.NewCheckErrorf(checks.SearchFirefoxID, "Uncompressed size is 0", nil)
	// Call the function under test
	result := firefox.SearchEngineFirefox(browserutils.RealProfileFinder{}) // Replace with the actual function call

	// Assert the result
	require.Equal(t, expected, result)
}
