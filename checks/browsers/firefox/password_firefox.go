package checks

import (
	utils "InfoSec-Agent/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/andrewarchi/browser/firefox"
)

func password_firefox() ([]string, error) {
	ffdirectory, _ := utils.FirefoxFolder() //returns the path to the firefox profile directory
	var output []string
	content, err := os.Open(ffdirectory[0] + "\\logins.json") //reads the extensions.json file or returns an error
	if err != nil {
		return nil, err
	}
	defer content.Close()
	var extensions firefox.Extensions   //Creates a struct for the json file
	decoder := json.NewDecoder(content) //
	err = decoder.Decode(&extensions)
	if err != nil {
		return nil, err
	} // Can add more data to the output for extension data.
	for _, addon := range extensions.Addons { // Name of addon, type of addon, creator, active or not
		output = append(output, addon.DefaultLocale.Name+addon.Type+addon.DefaultLocale.Creator+fmt.Sprintf("%t", addon.Active))
	}
	return output, nil
}
