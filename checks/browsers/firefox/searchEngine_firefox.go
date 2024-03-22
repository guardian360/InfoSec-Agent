package firefox

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/pierrec/lz4"
)

// SearchEngineFireFox checks the standard search engine in firefox.
//
// Parameters: _
//
// Returns: The standard search engine for firefox
func SearchEngineFirefox() checks.Check {
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, _ := utils.FirefoxFolder()

	// Path to the search.json.mozlz4 file
	filePath := ffdirectory[0] + "/search.json.mozlz4"

	//Copy the json file so we don't have problems with locked files
	tempSearch := filepath.Join(os.TempDir(), "tempSearch.json.mozlz4")
	// Clean up the temporary file when the function returns
	defer os.Remove(tempSearch)

	// Copy the compressed json to a temporary location
	copyError := utils.CopyFile(filePath, tempSearch)
	if copyError != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to make a copy of the file", copyError)
	}

	// Get file information
	fileInfo, err := os.Stat(tempSearch)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to retrieve information about the file", err)
	}
	// Get the size of the compressed file
	fileSize := fileInfo.Size()

	//Holds the information from the copied file
	file, err := os.Open(tempSearch)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to open the file", err)
	}
	defer utils.CloseFile(file)

	//Holds the custom magig number for the mozzila lz4 compression
	magicNumber := make([]byte, 8)

	//Retrieves the magicNumber from the file
	_, err = file.Read(magicNumber)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	//Holds the size of the file after decompressing it
	uncompressSize := make([]byte, 4)

	// Skip the first 8 bytes to take the bytes 8-11 that hold the size after decompression
	_, err = file.Seek(8, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to skip the first 8 bytes", err)
	}

	//Here we read the 8-11 bytes to find the size of the file
	_, err = file.Read(uncompressSize)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	//Transforms the size of the file after decompression from Little Endian to a normal 32-bit integer
	unCompressedSize := binary.LittleEndian.Uint32(uncompressSize)

	// Seek to skip the first 12 bytes
	_, err = file.Seek(12, io.SeekStart)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to skip the first 12 bytes", err)
	}

	//Byte slice to hold the compressed data without the header (first 12 bytes)
	compressedData := make([]byte, fileSize-12)

	_, err = file.Read(compressedData)
	if err != nil {
		return checks.NewCheckErrorf("SearchEngineFirefox", "Unable to read the file", err)
	}

	//Byte slice to hold all the uncompressed data
	data := make([]byte, unCompressedSize)
	// Uncompresses the file and puts it into the data slice
	lz4.UncompressBlock(compressedData, data)
	// Transforms the data into a readable string
	output := string(data)
	var result string
	// Regex to check if the defaultEngineId is empty which means that the engine is Google
	re := regexp.MustCompile(`"defaultEngineId":""`)
	// Performs the regex on the string and returns either google or goes into the next check
	matches := re.FindString(output)
	if matches != "" {
		result = "Google"
	} else {
		// This pattern looks for the values of the other known search engines and returns them
		pattern := `defaultEngineId":"(?:ddg|bing|ebay|wikipedia|amazon)@search\.mozilla\.org`
		// Compile the regex pattern
		re := regexp.MustCompile(pattern)

		// Find all matches in the text
		matches := re.FindAllString(output, -1)

		// adds the match to the result
		if len(matches) > 0 {
			result = matches[0]
		} else {
			result = "Other search engine"
		}
	}
	return checks.NewCheckResult("SearchEngineFirefox", result)
}
