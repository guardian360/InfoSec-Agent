// Package chromium is responsible for running checks on Chromium based browsers.
//
// Exported function(s): ExtensionsChromium, HistoryChromium
package chromium

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
)

// Response is a struct that represents the JSON response from the Microsoft Edge Addons Store
//
// Fields: Name (string) - The name of the extension
type Response struct {
	Name string `json:"name"`
}

// ExtensionsChromium checks if an adblocker is installed in a Chromium based browser.
//
// Parameters:
//
//	browser (string) - The name of the browser to check
//
// Returns: If the user has an ad blocker installed
func ExtensionsChromium(browser string) checks.Check {
	var browserPath string
	var returnBrowserName string
	// Set the browser path and the return browser name based on the browser to check
	// Currently, supports checking of Google Chrome and Microsoft Edge
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
	// Get the current user's home directory, where the extensions are stored
	user, err := os.UserHomeDir()
	if err != nil {
		checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	extensionsDir := filepath.Join(user, "AppData", "Local", browserPath, "User Data", "Default", "Extensions")
	files, err := os.ReadDir(extensionsDir)
	if err != nil {
		checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	// Construct a list of all extensions Id's
	for _, f := range files {
		if f.IsDir() {
			extensionIds = append(extensionIds, f.Name())
		}
	}

	// extensionName := ""
	for _, id := range extensionIds {
		// Get the name of the extension from the Chrome Web Store
		extensionName1, err := getExtensionNameChromium(id,
			"https://chromewebstore.google.com/detail/%s", browser)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Count(extensionName1, "/") > 4 {
			parts := strings.Split(extensionName1, "/")
			extensionNames = append(extensionNames, parts[len(parts)-2])
		}
		if browser == "Edge" {
			// Get the name of the extension from the Microsoft Edge Addons Store
			extensionName2, err := getExtensionNameChromium(id,
				"https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s", browser)
			if err != nil {
				log.Fatal(err)
			}
			extensionNames = append(extensionNames, extensionName2)
		}
	}
	if adblockerInstalled(extensionNames) {
		return checks.NewCheckResult(returnBrowserName, "Adblocker installed")
	} else {
		return checks.NewCheckErrorf(returnBrowserName, "No adblocker installed", errors.New("no adblocker installed"))
	}
}

// getExtensionNameChromium gets the name of an extension from the Chrome Web Store or the Microsoft Edge Addons Store
//
// Parameters:
//
//	extensionID (string) - The ID of the extension
//
//	url (string) - The URL of the store to visit
//
//	browser (string) - The name of the browser to check
//
// Returns: The name of the extension and an optional error (string,error)
// Error should be nil on success
func getExtensionNameChromium(extensionID string, url string, browser string) (string, error) {
	client := &http.Client{}
	urlToVisit := fmt.Sprintf(url, extensionID)
	// Generate a new request to visit the extension/addon store
	req, err := http.NewRequest("GET", urlToVisit, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// Close the response body after the necessary data is retrieved
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing body: ", err)
		}
	}(resp.Body)

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
			// For Edge, the data is stored in a JSON format, so decoding is required
			var data Response
			err := json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				return "", err
			}
			return data.Name, nil
		}
	} else {
		return "", errors.New("unknown browser")
	}
}

// adblockerInstalled checks if an adblocker is installed
//
// Parameters: extensionNames ([]string) - A list of the names of the installed extensions
//
// Returns: If an adblocker is installed (bool)
func adblockerInstalled(extensionNames []string) bool {
	// List of adblocker (related) terms to check for in the name of the extension/addon
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
	// If any of these terms appear in the name, it is assumed the extension is an adblocker
	for _, extensionName := range extensionNames {
		for _, adblockerName := range adblockerNames {
			if strings.Contains(strings.ToLower(extensionName), adblockerName) {
				return true
			}
		}
	}
	return false
}
