package checks

import (
	"InfoSec-Agent/checks"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ExtensionsChrome() check.check {
	var extensionIds []string
	var extensionNames []string
	user, err := os.UserHomeDir()
	if err != nil {
		check.NewCheckErrorf("ExtensionsChrome", "Error: ", err)
	}

	extensionsDir := filepath.Join(user, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Extensions")
	files, err := ioutil.ReadDir(extensionsDir)
	if err != nil {
		check.NewCheckErrorf("ExtensionsChrome", "Error: ", err)
	}

	//Adds the extnesion ID to the extensionIds array
	fmt.Println("Installed Chrome Extensions:")
	for _, f := range files {
		if f.IsDir() {
			extensionIds = append(extensionIds, f.Name())
		}
	}

	for _, id := range extensionIds {
		extensionName, err := getExtensionName(id)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Count(extensionName, "/") > 4 {
			extensionNames = append(extensionNames, strings.Split(extensionName, "/")[4])
		}
	}
	if adblockerInstalled(extensionNames) {
		checks.NewCheckResult("ExtensionsChrome", "Adblocker installed")
	} else {
		checks.NewCheckError("ExtensionsChrome", "No adblocker installed")
	}
}

func getExtensionName(extensionID string) check.Check {
	client := &http.Client{}
	url := fmt.Sprintf("https://chromewebstore.google.com/detail/%s", extensionID)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		check.NewCheckErrorf("ExtensionsChrome", "Error: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		check.NewCheckErrorf("ExtensionsChrome", "Failed to get extension name, status code: %d", resp.StatusCode)
	}

	return resp.Request.URL.String(), nil
}

func adblockerInstalled(extensionNames []string) bool {
	adblockerNames := []string{
		"adblock",
		"adblox",
		"advertentieblokker",
		"ad skip",
		"adkrig",
		"adblokker",
		"advertentieblokkering",
		"ad lock",
		"adlock",
		"privacy badger",
		"ublock",
		"adguard",
		"adaware",
		"adaware adblock",
		"ghostery",
	}
	for _, extensionName := range extensionNames {
		for _, adblockerName := range adblockerNames {
			if strings.Contains(strings.ToLower(extensionName), adblockerName) {
				return true
			}
		}
	}
	return false
}
