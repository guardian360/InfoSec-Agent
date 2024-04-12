package firefox

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/pierrec/lz4"
)

// SearchEngineFirefox is a function that retrieves the default search engine used in the Firefox browser.
//
// Parameters: None
//
// Returns:
//   - checks.Check: A Check object that encapsulates the result of the search engine check. The Check object includes a string that represents the default search engine in the Firefox browser. If an error occurs during the check, the Check object will encapsulate this error.
//
// This function first determines the directory in which the Firefox profile is stored. It then opens and reads the 'search.json.mozlz4' file, which contains information about the default search engine. The function decompresses the file, extracts the default search engine information, and returns this information as a Check object. If an error occurs at any point during this process, it is encapsulated in the Check object and returned.
func SearchEngineFirefox() checks.Check {
	// Determine the directory in which the Firefox profile is stored
	var ffdirectory []string
	var err error
	ffdirectory, err = utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "No firefox directory found", err)
	}
	filePath := ffdirectory[0] + "/search.json.mozlz4"

	// Create a temporary file to copy the compressed json to
	tempSearch := filepath.Join(os.TempDir(), "tempSearch.json.mozlz4")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			log.Println("Error deleting temporary file")
		}
	}(tempSearch)

	// Copy the compressed json to a temporary location
	copyError := utils.CopyFile(filePath, tempSearch, nil, nil)
	if copyError != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to make a copy of the file", copyError)
	}

	fileInfo, err := os.Stat(tempSearch)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to retrieve information about the file", err)
	}
	fileSize := fileInfo.Size()

	// Holds the information from the copied file
	file, err := os.Open(filepath.Clean(tempSearch))
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to open the file", err)
	}
	defer func(file *os.File) {
		tmpFile := mocking.Wrap(file)
		err = utils.CloseFile(tmpFile)
		if err != nil {
			log.Println("Error closing file")
		}
	}(file)

	// Holds the custom magic number for the mozilla lz4 compression
	magicNumber := make([]byte, 8)

	// Retrieves the magicNumber from the file
	_, err = file.Read(magicNumber)
	if err != nil {
		return checks.NewCheckErrorf(checks.SearchFirefoxID, "Unable to read the file", err)
	}

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
	return checks.NewCheckResult(checks.SearchFirefoxID, 0, results(data))
}

// results is a utility function used within the SearchEngineFirefox function.
// It processes the output string from the decompressed 'search.json.mozlz4' file to identify the default search engine.
//
// Parameters:
//   - output (string): Represents the decompressed output string from the 'search.json.mozlz4' file.
//
// Returns:
//   - string: A string that represents the default search engine in the Firefox browser. If the defaultEngineId is empty, the function returns "Google". If the defaultEngineId matches known search engines (ddg, bing, ebay, wikipedia, amazon), the function returns the name of the matched search engine. If the defaultEngineId does not match any known search engines, the function returns "Other Search Engine".
//
// This function first checks if the defaultEngineId in the output string is empty, which indicates that the default search engine is Google. If the defaultEngineId is not empty, the function checks if it matches the ids of other known search engines. If a match is found, the function returns the name of the matched search engine. If no match is found, the function returns "Other Search Engine".
func results(data []byte) string {
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
