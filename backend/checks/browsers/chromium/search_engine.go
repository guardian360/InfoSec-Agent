package chromium

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
)

// OpenFileFunc is a variable that holds the function used to open a file.
var OpenFileFunc = os.Open

// SearchEngineChromium inspects the default search engine setting in Chromium-based browsers.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to inspect. This could be "Chrome", "Edge", etc.
//   - mockBool: A boolean indicating whether to use mock data for testing.
//   - mockFile: A mocking.File object representing the file to use for testing.
//   - getter: A browsers.DefaultDirGetter object that provides the path to the preferences file for the specified browser.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains the name of the default search engine used in the specified browser. If an error occurs during the check, the result will contain a description of the error.
//
// This function works by locating the preferences file in the user's home directory, which contains the browser's settings. It opens and reads this file, then parses it as JSON to access the settings. It specifically looks for the "default_search_provider_data" key in the JSON data, which holds the name of the default search engine. If this key is found, its value is returned as the result of the check. If any error occurs during this process, such as an error reading the file or parsing the JSON, this error is returned as the result of the check.
func SearchEngineChromium(browser string, mockBool bool, mockFile mocking.File, getter browsers.DefaultDirGetter) checks.Check {
	browserPath, returnID := GetBrowserPathAndIDSearch(browser)
	// Holds the return value and sets the default value to Chrome in case you never changed your search engine
	defaultSE := "google.com"

	preferencesDir, err := getter.GetDefaultDir(browserPath)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "error getting preferences directory", err)
	}
	var dev map[string]interface{}
	var file mocking.File

	if !mockBool {
		tmpfile, openErr := os.Open(preferencesDir + "/Preferences")
		if openErr != nil {
			return checks.NewCheckErrorf(returnID, "error opening preferences directory", err)
		}
		file = mocking.Wrap(tmpfile)
		defer func(file mocking.File) {
			openErr = browsers.CloseFile(file)
			if openErr != nil {
				logger.Log.Error("Error closing file")
			}
		}(file)
	} else {
		file = mockFile
	}

	dev, err = ParsePreferencesFile(file)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "error parsing preferences directory", err)
	}

	defaultSE = GetDefaultSearchEngine(dev, defaultSE)

	return checks.NewCheckResult(returnID, 0, defaultSE)
}

// GetBrowserPathAndIDSearch is a function that takes a browser name as input,
// and returns the path to the browser's directory and the ID of the browser.
//
// Parameters:
//   - browser: A string representing the name of the browser. Currently, this function supports "Chrome" and "Edge".
//
// Returns:
//   - A string representing the path to the browser's directory.
//   - An integer representing the ID the check.
//
// This function works by checking the provided browser name against known browser names. If the browser is "Chrome",
// the function returns the path to the Chrome directory and the ID of Chrome. If the browser is "Edge",
// the function returns the path to the Edge directory and the ID of Edge. If the browser is unknown, the function
// returns an empty string and 0.
func GetBrowserPathAndIDSearch(browser string) (string, int) {
	if browser == browsers.Chrome {
		return browsers.ChromePath, checks.SearchChromiumID
	}
	if browser == browsers.Edge {
		return browsers.EdgePath, checks.SearchEdgeID
	}
	return "", 0
}

// ParsePreferencesFile is a function that reads and parses a preferences file.
//
// Parameters:
//   - file mocking.File: The file object representing the preferences file to be parsed.
//
// Returns:
//   - map[string]interface{}: A map representing the parsed JSON data from the preferences file. The keys are the setting names and the values are the setting values.
//   - error: An error object that wraps any error that occurs during the reading and parsing of the preferences file. If the file is read and parsed successfully, it returns nil.
//
// This function works by reading all the bytes from the file, then unmarshalling the bytes into a map. If any error occurs during this process, such as an error reading the file or parsing the JSON, this error is returned.
func ParsePreferencesFile(file mocking.File) (map[string]interface{}, error) {
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var dev map[string]interface{}
	err = json.Unmarshal(byteValue, &dev)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

// GetDefaultSearchEngine is a function that retrieves the default search engine from a parsed preferences file.
//
// Parameters:
//   - dev map[string]interface{}: A map representing the parsed JSON data from the preferences file. The keys are the setting names and the values are the setting values.
//   - defSE string: A string representing the default search engine to return if the "default_search_provider_data" key is not found in the parsed preferences file.
//
// Returns:
//   - string: The name of the default search engine. If the "default_search_provider_data" key is found in the parsed preferences file, the value of this key is returned. If the key is not found, the value of defSE is returned.
//
// This function works by iterating over the keys and values in the parsed preferences file. If it finds the "default_search_provider_data" key, it uses a regular expression to extract the name of the default search engine from the value of this key. If the key is not found, it returns the value of defSE.
func GetDefaultSearchEngine(dev map[string]interface{}, defSE string) string {
	for key, value := range dev {
		if key == "default_search_provider_data" {
			text := fmt.Sprintf("%v", value)
			pattern := `keyword:\s*([^]\s]+)`
			regex := regexp.MustCompile(pattern)
			matches := regex.FindString(text)
			if matches != "" {
				return matches[8:]
			}
		}
	}
	return defSE
}
