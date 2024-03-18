package firefox

import (
	"encoding/json"
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"github.com/andrewarchi/browser/firefox"
	"os"
)

func PasswordFirefox() checks.Check {
	ffdirectory, _ := utils.FirefoxFolder() //returns the path to the firefox profile directory
	var output []string
	content, err := os.Open(ffdirectory[0] + "\\logins.json") //reads the extensions.json file or returns an error
	if err != nil {
		return checks.NewCheckError("PasswordFirefox", err)
	}
	defer content.Close()
	var extensions firefox.Extensions   //Creates a struct for the json file
	decoder := json.NewDecoder(content) //
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError("PasswordFirefox", err)
	}

	// Can add more data to the output for extension data.
	for _, addon := range extensions.Addons { // Name of addon, type of addon, creator, active or not
		output = append(output, addon.DefaultLocale.Name+addon.Type+addon.DefaultLocale.Creator+fmt.Sprintf("%t", addon.Active))
	}
	return checks.NewCheckResult("PasswordFirefox", output...)
}
