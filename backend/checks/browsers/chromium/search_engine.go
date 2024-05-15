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

var OpenFileFunc = os.Open

// SearchEngineChromium inspects the default search engine setting in Chromium-based browsers.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to inspect. This could be "Chrome", "Edge", etc.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains the name of the default search engine used in the specified browser. If an error occurs during the check, the result will contain a description of the error.
//
// This function works by locating the preferences file in the user's home directory, which contains the browser's settings. It opens and reads this file, then parses it as JSON to access the settings. It specifically looks for the "default_search_provider_data" key in the JSON data, which holds the name of the default search engine. If this key is found, its value is returned as the result of the check. If any error occurs during this process, such as an error reading the file or parsing the JSON, this error is returned as the result of the check.
func SearchEngineChromium(browser string, mockBool bool, mockFile mocking.File, getter browsers.PreferencesDirGetter) checks.Check {
	var returnID int
	if browser == browsers.Chrome {
		returnID = checks.SearchChromiumID
	}
	if browser == browsers.Edge {
		returnID = checks.SearchEdgeID
	}
	defaultSE := "google.com"

	preferencesDir, err := getter.GetPreferencesDir(browser)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}
	var dev map[string]interface{}
	var file mocking.File

	if !mockBool {
		tmpfile, openErr := os.Open(preferencesDir + "/Preferences")
		if openErr != nil {
			return checks.NewCheckErrorf(returnID, "Error: ", err)
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
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}

	defaultSE = GetDefaultSearchEngine(dev, defaultSE)

	return checks.NewCheckResult(returnID, 0, defaultSE)
}

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
