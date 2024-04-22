package firefox

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/andrewarchi/browser/firefox"
)

// PasswordFirefox is a function that checks the passwords stored in the Firefox browser.
//
// Parameters: None
//
// Returns:
//   - checks.Check: A Check object that encapsulates the results of the password check. The Check object includes a list of strings, where each string represents a saved password in the Firefox browser. If an error occurs during the password check, the Check object will encapsulate this error.
//
// This function first determines the directory in which the Firefox profile is stored. It then opens the 'logins.json' file, which contains a list of all saved Firefox passwords. The function decodes the JSON file into a struct, and then iterates over the struct to extract the saved passwords. These passwords are added to the results, which are returned as a Check object. If an error occurs at any point during this process, it is encapsulated in the Check object and returned.
func PasswordFirefox() checks.Check {
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := utils.RealProfileFinder{}.FirefoxFolder()
	if err != nil {
		logger.Log.ErrorWithErr("No firefox directory found: ", err)
		return checks.NewCheckErrorf(checks.HistoryFirefoxID, "No firefox directory found", err)
	}

	var output []string
	// Open the logins.json file, which contains a list of all saved Firefox passwords
	content, err := os.Open(ffdirectory[0] + "\\logins.json")
	if err != nil {
		return checks.NewCheckError(99, err)
	}
	defer func(content *os.File) {
		err = content.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing file: ", err)
		}
	}(content)

	// Creates a struct for the JSON file
	var extensions firefox.Extensions
	decoder := json.NewDecoder(content)
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError(99, err)
	}

	// TODO: Final functionality currently not implemented yet, should return an analysis on the used passwords
	// Below code is a placeholder
	for _, addon := range extensions.Addons {
		output = append(output,
			addon.DefaultLocale.Name+addon.Type+addon.DefaultLocale.Creator+strconv.FormatBool(addon.Active))
	}
	return checks.NewCheckResult(99, 0, output...)
}
