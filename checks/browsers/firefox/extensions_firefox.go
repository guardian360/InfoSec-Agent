package firefox

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"InfoSec-Agent/checks"
	utils "InfoSec-Agent/utils"

	"github.com/andrewarchi/browser/firefox"
)

func ExtensionFirefox() (checks.Check, checks.Check) {
	ffdirectory, _ := utils.FirefoxFolder() //returns the path to the firefox profile directory
	addBlocker := false                     //Variable used for checking if addblocker is used
	var output []string
	content, err := os.Open(ffdirectory[0] + "\\extensions.json") //reads the extensions.json file or returns an error
	if err != nil {
		return checks.NewCheckError("ExtensionsFirefox", err), checks.NewCheckError("AdblockerFirefox", err)
	}
	defer content.Close()
	var extensions firefox.Extensions   //Creates a struct for the json file
	decoder := json.NewDecoder(content) //
	err = decoder.Decode(&extensions)
	if err != nil {
		return checks.NewCheckError("ExtensionsFirefox", err), checks.NewCheckError("AdblockerFirefox", err)
	} // Can add more data to the output for extension data.
	
	for _, addon := range extensions.Addons { // Name of addon, type of addon, creator, active or not
		output = append(output, addon.DefaultLocale.Name+addon.Type+addon.DefaultLocale.Creator+fmt.Sprintf("%t", addon.Active))
		if adblockerFirefox(addon.DefaultLocale.Name) {
			addBlocker = true
		}
	}
	adBlockused := fmt.Sprintf("%t", addBlocker)
	return checks.NewCheckResult("ExtensionsFirefox", output...), checks.NewCheckResult("AdblockerFirefox", adBlockused)
}

func adblockerFirefox(extensionName string) bool { //Check for the most used adblockers in firefox
	adblockerNames := []string{
		"adblocker ultimate",
		"adguard adblocker",
		"adblocker for youtube",
		"ublock origin",
		"adblock plus",
		"adblock for firefox",
	}
	for _, adblockerName := range adblockerNames {
		if strings.Contains(strings.ToLower(extensionName), adblockerName) {
			return true
		}
	}
	return false
}
