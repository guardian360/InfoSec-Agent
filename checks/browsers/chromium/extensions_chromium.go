// Package chromium is responsible for running checks on Chromium based browsers.
//
// Exported function(s): ExtensionsChromium, HistoryChromium
package chromium

import (
	"context"
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

// constants to store information of the browsers
const edge = "Edge"
const chrome = "Chrome"
const edgePath = "Microsoft/Edge"
const chromePath = "Google/Chrome"

// ExtensionsChromium checks for the presence of an adblocker extension in a specified Chromium-based browser.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to check.
//     Currently, this function supports "Chrome" and "Edge".
//
// Returns:
//   - A checks.Check object representing the result of the check. If an adblocker is installed,
//     the result will be "Adblocker installed". If no adblocker is found, the result will be
//     "No adblocker installed". If an error occurs during the check, the function will return
//     a checks.CheckError with the error message.
//
// This function works by reading the extensions directory of the specified browser in the user's home directory,
// and checking each installed extension against a list of known adblocker names.
// The function uses the extensionNameChromium helper function to fetch the name of each extension from the
// Chrome Web Store or the Microsoft Edge Addons Store.
func ExtensionsChromium(browser string) checks.Check {
	var browserPath string
	var returnBrowserName string
	// Set the browser path and the return browser name based on the browser to check
	// Currently, supports checking of Google Chrome and Microsoft Edge
	if browser == chrome {
		returnBrowserName = "ExtensionsChrome"
		browserPath = chromePath
	}
	if browser == edge {
		returnBrowserName = "ExtensionsEdge"
		browserPath = edgePath
	}
	var extensionIDs []string
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

	// Construct a list of all extensions ID's
	for _, f := range files {
		if f.IsDir() {
			extensionIDs = append(extensionIDs, f.Name())
		}
	}

	var extensionName1 string
	var extensionName2 string
	for _, id := range extensionIDs {
		// Get the name of the extension from the Chrome Web Store
		extensionName1, err = extensionNameChromium(id,
			"https://chromewebstore.google.com/detail/%s", browser)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Count(extensionName1, "/") > 4 {
			parts := strings.Split(extensionName1, "/")
			extensionNames = append(extensionNames, parts[len(parts)-2])
		}
		if browser == edge {
			// Get the name of the extension from the Microsoft Edge Addons Store
			extensionName2, err = extensionNameChromium(id,
				"https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s", browser)
			if err != nil {
				log.Fatal(err)
			}
			extensionNames = append(extensionNames, extensionName2)
		}
	}
	if adblockerInstalled(extensionNames) {
		return checks.NewCheckResult(returnBrowserName, "Adblocker installed")
	}
	return checks.NewCheckErrorf(returnBrowserName, "No adblocker installed", errors.New("no adblocker installed"))
}

// extensionNameChromium fetches the name of an extension from the Chrome Web Store or the Microsoft Edge Addons Store.
//
// Parameters:
//   - extensionID: The unique identifier of the extension.
//   - url: The URL template of the store to visit, where the extension ID will be inserted.
//   - browser: The name of the browser to check (either "Chrome" or "Edge").
//
// Returns:
//   - The name of the extension.
//   - An error, which will be nil if the operation was successful.
//
// This function sends a HTTP request to the store's URL and parses the response to extract the extension name.
// For Edge, the response is in JSON format and requires decoding.
// If the HTTP request fails or the browser is unknown, the function returns an error.
func extensionNameChromium(extensionID string, url string, browser string) (string, error) {
	client := &http.Client{}
	urlToVisit := fmt.Sprintf(url, extensionID)
	// Generate a new request to visit the extension/addon store
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlToVisit, nil)
	if err != nil {
		log.Println("error creating request: ", err)
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error sending request: ", err)
		return "", err
	}
	// Close the response body after the necessary data is retrieved
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("error closing body: ", err)
		}
	}(resp.Body)

	if browser == chrome && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}
	if browser == chrome {
		return resp.Request.URL.String(), nil
	}
	if browser == edge {
		if strings.Contains(resp.Request.URL.String(), "chromewebstore.google.com") {
			return resp.Request.URL.String(), nil
		}
		// For Edge, the data is stored in a JSON format, so decoding is required
		var data Response
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", err
		}
		return data.Name, nil
	}
	return "", errors.New("unknown browser")
}

// adblockerInstalled determines whether an adblocker extension is installed.
//
// Parameters:
//   - extensionNames: A slice of strings representing the names of the installed extensions.
//
// Returns:
//   - A boolean value indicating whether an adblocker is installed. Returns true if an adblocker is found among the installed extensions, and false otherwise.
//
// This function works by iterating over the provided list of extension names and checking each one against a predefined list of known adblocker names. If a match is found, the function returns true. If no match is found after checking all extension names, the function returns false.
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
