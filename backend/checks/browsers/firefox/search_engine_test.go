package firefox_test

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"io"
	"os"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

func TestSearchEngineFirefox_WithInvalidDirectory(t *testing.T) {
	// Mock the FirefoxFolder function to return an invalid directory
	Profilefinder = browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"/invalid/directory"}, nil
		},
	}

	check := firefox.SearchEngineFirefox(Profilefinder, false, nil, nil)
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
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			logger.Log.Error("error removing file")
		}
	}(tempFile.Name())

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
	Profilefinder = browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"\\valid\\directory"}, nil
		},
	}

	mockFile := &mocking.FileMock{
		FileName: "mockfile",
		IsOpen:   true,
		Buffer:   []byte("This is some test data."),
		Bytes:    22,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}
	// Create a destination mock file
	mockDestinationFile := &mocking.FileMock{
		FileName: "\\valid\\directory\\destination.json.mozlz4",
		IsOpen:   true,
		Buffer:   []byte("This is some test data."),
		Bytes:    22,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}

	// Mock the OpenAndStatFile function
	// Save the original function and defer its restoration
	originalOpenAndStatFile := firefox.OpenAndStatFile
	defer func() { firefox.OpenAndStatFile = originalOpenAndStatFile }()

	// Mock the OpenAndStatFile function
	firefox.OpenAndStatFile = func(_ string) (mocking.File, int64, error) {
		// Return whatever you need for your test
		file := &mocking.FileMock{
			FileName: "\\valid\\directory\\search.json.mozlz4",
			IsOpen:   true,
			Buffer:   []byte{0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00},
			Bytes:    12,
			FileInfo: &mocking.FileInfoMock{},
			Err:      nil,
		}
		return file, 26, nil
	}
	expected := checks.NewCheckErrorf(checks.SearchFirefoxID, "Uncompressed size is 0", nil)
	// Call the function under test
	result := firefox.SearchEngineFirefox(Profilefinder, true, mockFile, mockDestinationFile) // Replace with the actual function call

	// Assert the result
	require.Equal(t, expected, result)
}

func TestYourFunctiontooShortFile(t *testing.T) {
	Profilefinder = browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"\\valid\\directory"}, nil
		},
	}

	mockFile := &mocking.FileMock{
		FileName: "mockfile",
		IsOpen:   true,
		Buffer:   []byte{0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00, 0x07, 0x00, 0x00, 0x00, 0x74},
		Bytes:    22,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}
	// Create a destination mock file
	mockDestinationFile := &mocking.FileMock{
		FileName: "\\valid\\directory\\destination.json.mozlz4",
		IsOpen:   true,
		Buffer:   []byte{0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00, 0x07, 0x00, 0x00, 0x00, 0x74},
		Bytes:    19,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}

	// Mock the OpenAndStatFile function
	// Save the original function and defer its restoration
	originalOpenAndStatFile := firefox.OpenAndStatFile
	defer func() { firefox.OpenAndStatFile = originalOpenAndStatFile }()

	// Mock the OpenAndStatFile function
	firefox.OpenAndStatFile = func(_ string) (mocking.File, int64, error) {
		// Return whatever you need for your test
		file := &mocking.FileMock{
			FileName: "\\valid\\directory\\search.json.mozlz4",
			IsOpen:   true,
			Buffer:   []byte{0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00, 0x07, 0x00, 0x00, 0x00, 0x74},
			Bytes:    26,
			FileInfo: &mocking.FileInfoMock{},
			Err:      nil,
		}
		return file, 26, nil
	}
	expected := checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to skip the first 12 bytes", io.EOF)
	// Call the function under test
	result := firefox.SearchEngineFirefox(Profilefinder, true, mockFile, mockDestinationFile) // Replace with the actual function call

	// Assert the result
	require.Equal(t, expected, result)
}

func TestYourFunctiontooUnCompressFile(t *testing.T) {
	Profilefinder = browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"\\valid\\directory"}, nil
		},
	}
	buff := []byte{0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00,
		0x08, 0x00, 0x00, 0x00, 0x6D, 0x6F, 0x7A, 0x4C,
		0x6D, 0x6F, 0x7A, 0x4C, 0x7A, 0x34, 0x30, 0x00,
		0x80, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
		0x00, 0x00, 0x00, 0x00, 0x54, 0x3a, 0x00, 0xbe}

	mockFile := &mocking.FileMock{
		FileName: "mockfile",
		IsOpen:   true,
		Buffer:   buff,
		Bytes:    28,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}
	// Create a destination mock file
	mockDestinationFile := &mocking.FileMock{
		FileName: "\\valid\\directory\\destination.json.mozlz4",
		IsOpen:   true,
		Buffer:   buff,
		Bytes:    28,
		FileInfo: &mocking.FileInfoMock{},
		Err:      nil,
	}

	// Mock the OpenAndStatFile function
	// Save the original function and defer its restoration
	originalOpenAndStatFile := firefox.OpenAndStatFile
	defer func() { firefox.OpenAndStatFile = originalOpenAndStatFile }()

	// Mock the OpenAndStatFile function
	firefox.OpenAndStatFile = func(_ string) (mocking.File, int64, error) {
		// Return whatever you need for your test
		file := &mocking.FileMock{
			FileName: "\\valid\\directory\\search.json.mozlz4",
			IsOpen:   true,
			Buffer:   buff,
			Bytes:    28,
			FileInfo: &mocking.FileInfoMock{},
			Err:      nil,
		}
		return file, 21, nil
	}
	expected := checks.NewCheckResult(checks.SearchFirefoxID, 0, "Other Search Engine")
	// Call the function under test
	result := firefox.SearchEngineFirefox(Profilefinder, true, mockFile, mockDestinationFile) // Replace with the actual function call

	// Assert the result
	require.Equal(t, expected, result)
}
