package chrome

import (
	"InfoSec-Agent/checks"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ExtensionsChromium(browser string) checks.Check {
	var browserPath string
	var returnBrowserName string
	if browser == "Chrome" {
		returnBrowserName = "ExtensionsChrome"
		browserPath = "Google/Chrome"
	}
	if browser == "Edge" {
		returnBrowserName = "ExtensionsEdge"
		browserPath = "Microsoft/Edge"
	}
	var extensionIds []string
	var extensionNames []string
	user, err := os.UserHomeDir()
	if err != nil {
		checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	extensionsDir := filepath.Join(user, "AppData", "Local", browserPath, "User Data", "Default", "Extensions")
	files, err := ioutil.ReadDir(extensionsDir)
	if err != nil {
		checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	//Adds the extnesion ID to the extensionIds array
	fmt.Println("Installed", browser, "extensions:")
	for _, f := range files {
		if f.IsDir() {
			extensionIds = append(extensionIds, f.Name())
		}
	}

	//extensionName := ""
	for _, id := range extensionIds {
		extensionName1, err := getExtensionNameChromium(id, "https://chromewebstore.google.com/detail/%s", browser)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Count(extensionName1, "/") > 4 {
			parts := strings.Split(extensionName1, "/")
			extensionNames = append(extensionNames, parts[len(parts)-2])
			fmt.Println(parts[len(parts)-2])
		}
		if browser == "Edge" {
			extensionName2, err := getExtensionNameChromium(id, "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s", browser)
			if err != nil {
				log.Fatal(err)
			}
			extensionNames = append(extensionNames, extensionName2)
			fmt.Println(extensionName2)
		}
	}
	if adblockerInstalled(extensionNames) {
		return checks.NewCheckResult(returnBrowserName, "Adblocker installed")
	} else {
		return checks.NewCheckErrorf(returnBrowserName, "No adblocker installed", errors.New("No adblocker installed"))
	}
}

func getExtensionNameChromium(extensionID string, url string, browser string) (string, error) {
	client := &http.Client{}
	urlToVisit := fmt.Sprintf(url, extensionID)
	req, err := http.NewRequest("GET", urlToVisit, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if browser == "Chrome" && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}
	if browser == "Chrome" {
		return resp.Request.URL.String(), nil
	}
	if browser == "Edge" {
		if strings.Contains(resp.Request.URL.String(), "chromewebstore.google.com") {
			return resp.Request.URL.String(), nil
		} else {
			var data Response
			json.NewDecoder(resp.Body).Decode(&data)
			return data.Name, nil
		}
	} else {
		return "", errors.New("Unknown browser")
	}
}

type Response struct {
	Name string `json:"name"`
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
