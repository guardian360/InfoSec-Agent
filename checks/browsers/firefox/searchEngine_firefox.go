package firefox

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browserutils"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/pierrec/lz4"
)

// SearchEngineFirefox is a function that retrieves the default search engine used in the Firefox browser.
//
// Parameters:
//   - profileFinder: An object that implements the FirefoxProfileFinder interface. It is used to find the Firefox profile directory.
//   - boolMock: A boolean value that determines whether to use the mockSource and mockDest files for testing.
//   - mockSource: A mock file used for testing. It is used as the source file when boolMock is true.
//   - mockDest: A mock file used for testing. It is used as the destination file when boolMock is true.
//
// Returns:
//   - checks.Check: A Check object that encapsulates the result of the search engine check. The Check object includes a string that represents the default search engine in the Firefox browser. If an error occurs during the check, the Check object will encapsulate this error.
//
// This function first determines the directory in which the Firefox profile is stored. It then opens and reads the 'search.json.mozlz4' file, which contains information about the default search engine. The function decompresses the file, extracts the default search engine information, and returns this information as a Check object. If an error occurs at any point during this process, it is encapsulated in the Check object and returned.
func SearchEngineFirefox(profileFinder browserutils.FirefoxProfileFinder, boolMock bool, mockSource mocking.File, mockDest mocking.File) checks.Check {
	// Determine the directory in which the Firefox profile is stored
	var ffDirectory []string
	var err error
	ffDirectory, err = profileFinder.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "No firefox directory found", err)
	}
	filePath := ffDirectory[0] + "/search.json.mozlz4"

	// Create a temporary file to copy the compressed json to
	tempSearch := filepath.Join(os.TempDir(), "tempSearch.json.mozlz4")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			log.Println("Error deleting temporary file")
		}
	}(tempSearch)

	if !boolMock {
		// Copy the compressed json to a temporary location
		copyError := browserutils.CopyFile(filePath, tempSearch, nil, nil)
		if copyError != nil {
			return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to make a copy of the file", copyError)
		}
	} else {
		// Copy the compressed json to a temporary location
		copyError := browserutils.CopyFile(filePath, tempSearch, mockSource, mockDest)
		if copyError != nil {
			return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to make a copy of the file", copyError)
		}
	}

	file, fileSize, err := OpenAndStatFile(tempSearch)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to open the file", err)
	}
	defer func(file mocking.File) {
		err = browserutils.CloseFile(file)
		if err != nil {
			log.Println("Error closing file")
		}
	}(file)

	// Holds the size of the file after decompressing it
	uncompressSize := make([]byte, 4)

	// Skip the first 8 bytes to take the bytes 8-11 that hold the size after decompression
	_, err = file.Seek(8, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to skip the first 8 bytes", err)
	}

	// Retrieves bytes 8-11 to find the size of the file
	_, err = file.Read(uncompressSize)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to read the file", err)
	}

	// Transforms the size of the file after decompression from Little Endian to a normal 32-bit integer
	unCompressedSize := binary.LittleEndian.Uint32(uncompressSize)
	if unCompressedSize == 0 {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Uncompressed size is 0", nil)
	}

	// Skip the first 12 bytes because that is the start of the data
	_, err = file.Seek(12, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to skip the first 12 bytes", err)
	}

	// Byte slice to hold the compressed data without the header (first 12 bytes)
	compressedData := make([]byte, fileSize-12)

	_, err = file.Read(compressedData)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to read the file", err)
	}

	data := make([]byte, unCompressedSize)
	_, err = lz4.UncompressBlock(compressedData, data)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to uncompress", err)
	}
	return checks.NewCheckResult(checks.SearchFirefoxID, 0, Results(data))
}

// Results is a utility function used within the SearchEngineFirefox function.
// It processes the output string from the decompressed 'search.json.mozlz4' file to identify the default search engine.
//
// Parameters:
//   - output (string): Represents the decompressed output string from the 'search.json.mozlz4' file.
//
// Returns:
//   - string: A string that represents the default search engine in the Firefox browser. If the defaultEngineId is empty, the function returns "Google". If the defaultEngineId matches known search engines (ddg, bing, ebay, wikipedia, amazon), the function returns the name of the matched search engine. If the defaultEngineId does not match any known search engines, the function returns "Other Search Engine".
//
// This function first checks if the defaultEngineId in the output string is empty, which indicates that the default search engine is Google. If the defaultEngineId is not empty, the function checks if it matches the ids of other known search engines. If a match is found, the function returns the name of the matched search engine. If no match is found, the function returns "Other Search Engine".
func Results(data []byte) string {
	output := string(data)
	var result string
	var re *regexp.Regexp
	var matches string
	// Regex to check if the defaultEngineId is empty which means that the engine is Google
	re = regexp.MustCompile(`"defaultEngineId":""`)
	// Performs the regex on the string and returns either google or goes into the next check
	matches = re.FindString(output)
	if matches != "" {
		result = "Google"
	} else {
		// This pattern looks for the values of the other known search engines and returns them
		pattern := `defaultEngineId":"(?:ddg|bing|ebay|wikipedia|amazon)@search\.mozilla\.org`
		re = regexp.MustCompile(pattern)
		matches = re.FindString(output)
		if matches == "" {
			return "Other Search Engine"
		}
		result = matches[18:]
	}
	return result
}

// OpenAndStatFile is a function that opens a file and retrieves its size.
//
// Parameters:
//   - tempSearch: A string that represents the path to the file that should be opened.
//
// Returns:
//   - mocking.File: A File object that represents the opened file. This object is wrapped in a mocking layer for testing purposes.
//   - int64: An integer that represents the size of the file in bytes.
//   - error: An error object that encapsulates any errors that occurred during the execution of the function. If no errors occurred, this object is nil.
//
// This function first calls the os.Stat function to retrieve the FileInfo object for the file. If an error occurs during this call, the function returns nil, 0, and the error. If no error occurs, the function retrieves the size of the file from the FileInfo object.
// The function then calls the os.Open function to open the file. If an error occurs during this call, the function returns nil, 0, and the error. If no error occurs, the function wraps the opened file in a mocking layer and returns the wrapped file, the size of the file, and nil for the error.
var OpenAndStatFile = func(tempSearch string) (mocking.File, int64, error) {
	// Retrieve the FileInfo object for the file
	fileInfo, err := os.Stat(tempSearch)
	if err != nil {
		// If an error occurred, return nil, 0, and the error
		return nil, 0, err
	}
	// Retrieve the size of the file from the FileInfo object
	fileSize := fileInfo.Size()

	// Open the file
	file, err := os.Open(filepath.Clean(tempSearch))
	if err != nil {
		// If an error occurred, return nil, 0, and the error
		return nil, 0, err
	}

	// Wrap the opened file in a mocking layer and return the wrapped file, the size of the file, and nil for the error
	return mocking.Wrap(file), fileSize, nil
}
