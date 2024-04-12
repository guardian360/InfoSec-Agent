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
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
)

// Response represents the JSON response from the Microsoft Edge Addons Store.
// It contains a single field, Name, which holds the name of the extension.
type Response struct {
	Name string `json:"name"`
}

const edge = "Edge"
const chrome = "Chrome"
const edgePath = "Microsoft/Edge"
const chromePath = "Google/Chrome"

// ExtensionsChromium is a function that determines the presence of an ad-blocker extension in a Chromium-based browser.
//
// Parameters:
//   - browser (string): This parameter represents the name of the Chromium-based browser to be checked. It could be Google Chrome, Microsoft Edge, etc.
//
// Returns:
//   - Check: A Check object is returned. If an ad-blocker is installed in the specified browser, the Check object encapsulates this information. If no adblocker is found, or if an error occurs during the check, the Check object will encapsulate this information as well.
//
// The function works by examining the extensions directory of the specified browser in the user's home directory. It constructs a list of all installed extension IDs, then retrieves the name of each extension from the respective browser's extension store (Chrome Web Store or Microsoft Edge Addons Store). If the name of any extension matches a predefined list of adblocker-related terms, the function concludes that an adblocker is installed.
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
			logger.Log.Fatal(err)
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
				logger.Log.Fatal(err)
			}
			extensionNames = append(extensionNames, extensionName2)
		}
	}
	if adblockerInstalled(extensionNames) {
		return checks.NewCheckResult(returnBrowserName, "Adblocker installed")
	}
	return checks.NewCheckErrorf(returnBrowserName, "No adblocker installed", errors.New("no adblocker installed"))
}

// extensionNameChromium is a function that retrieves the name of a specific extension from either the Chrome Web Store or the Microsoft Edge Addons Store.
//
// Parameters:
//   - extensionID (string): The unique identifier of the extension. This ID is used to locate the extension in the respective store.
//   - url (string): The URL template of the store where the extension is hosted. The extension ID is inserted into this template to form the complete URL.
//
// This function initiates an HTTP request to the provided store URL and processes the response to extract the extension's name. For extensions on the Microsoft Edge Addons Store, the response is in JSON format and thus requires decoding.
// If the HTTP request encounters an error or if the browser is not recognized, the function will return an error.
//
// The function works by first formatting the URL with the extension ID. It then creates a new HTTP request with this URL and sends it. If the HTTP request is successful, the function processes the response. For Chrome, it returns the URL string of the request. For Edge, it decodes the JSON response to extract the extension's name. If the HTTP request fails, the function returns an error. If the browser is neither Chrome nor Edge, the function also returns an error.
//
// Returns:
//   - string: The name of the extension if the HTTP request is successful and the browser is recognized. If the browser is Chrome, the function returns the URL string of the request. If the browser is Edge, the function decodes the JSON response to extract the extension's name.
//   - error: An error object that encapsulates any error that occurred during the HTTP request or the processing of the response. This could be due to the HTTP request failing, the browser not being recognized, or an error occurring while decoding the JSON response for Edge.
func extensionNameChromium(extensionID string, url string, browser string) (string, error) {
	client := &http.Client{}
	urlToVisit := fmt.Sprintf(url, extensionID)
	// Generate a new request to visit the extension/addon store
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlToVisit, nil)
	if err != nil {
		logger.Log.Println("error creating request: ", err)
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Println("error sending request: ", err)
		return "", err
	}
	// Close the response body after the necessary data is retrieved
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Log.Println("error closing body: ", err)
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
