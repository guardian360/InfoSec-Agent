package checks

import (
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// TODO: The "newest build" is done manually below and should be done automatic.
// So preferably it would read the newest build without us changing the values.
const (
	newestWin10Build uint32 = 19045 // Version 22H2 (2022 update)
	newestWin11Build uint32 = 22631 // Version 23H2 (2023 update)
)

// WindowsOutdated is a function that checks if the currently installed Windows version is outdated.
//
// Parameters:
//   - mockOS mocking.WindowsVersion: A mock object for retrieving the Windows version information.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the Windows version is up-to-date or if updates are available.
//
// The function works by retrieving the Windows version information using the provided mock object. It then compares the build number of the installed Windows version with the build numbers of the latest Windows 10 and Windows 11 versions. If the installed version's build number matches the latest build number for its major version (10 or 11), the function returns a message indicating that the Windows version is up-to-date. If the build number does not match, the function returns a message indicating that updates are available. If the major version is neither 10 nor 11, the function returns a message suggesting to update to Windows 10 or Windows 11.
func WindowsOutdated(mockOS mocking.WindowsVersion) Check {
	versionData := mockOS.RtlGetVersion()
	versionString := fmt.Sprintf("%d.%d.%d", versionData.MajorVersion, versionData.MinorVersion, versionData.BuildNumber)

	win10HTML := getUrlBody("https://learn.microsoft.com/en-us/windows/release-health/release-information")
	latestWin10Build := findWindowsBuild(win10HTML)

	win11HTML := getUrlBody("https://learn.microsoft.com/en-us/windows/release-health/windows11-release-information")
	latestWin11Build := findWindowsBuild(win11HTML)

	// Depending on the major Windows version (10 or 11), act accordingly
	switch versionData.MajorVersion {
	case 11:
		if versionData.BuildNumber == latestWin11Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString, "You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString, "There are updates available for Windows 11.")
		}
	case 10:
		if versionData.BuildNumber == latestWin10Build {
			return NewCheckResult(WindowsOutdatedID, 0, versionString, "You are currently up to date.")
		} else {
			return NewCheckResult(WindowsOutdatedID, 1, versionString, "There are updates available for Windows 10.")
		}
	default:
		return NewCheckResult(WindowsOutdatedID, 2, versionString,
			"You are using a Windows version which does not have support anymore. "+
				"Consider updating to Windows 10 or Windows 11.")
	}
}

// getUrlBody fetches and parses the HTML content of a given URL.
//
// This function makes an HTTP GET request to the provided URL and parses the HTML content of the response.
// It logs any errors that occur during the HTTP request or the HTML parsing.
// The function returns the root node of the parsed HTML document.
//
// Parameters:
//
//   - url string - The URL to fetch and parse the HTML content from.
//
// Returns: The root node of the parsed HTML document.
func getUrlBody(url string) *html.Node {
	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.ErrorWithErr("Error fetching URL: ", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing response body: ", err)
		}
	}(resp.Body)

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Log.ErrorWithErr("Error parsing HTML: ", err)
	}
	return doc
}

// findWindowsBuild searches for the latest Windows build in the HTML content of a given URL.
//
// This function iterates over the children of the provided HTML node. If the node is a table body (tbody),
// the function iterates over its children.
// If a child is a table row (tr), the function counts the number of table data (td) elements in the row.
// When it finds the fifth td element, it extracts and returns the data as a string.
// If the function does not find a tbody or a tr with five td elements,
// it continues the search recursively on the node's children.
//
// The function is designed to work for the specific layout of the HTML content at the provided URL.
// Should this layout change, the function may need to be updated to reflect the new structure.
//
// Parameters:
//
//   - n *html.Node - The HTML node to search for the data element.
//
// Returns: The data from the fifth td element in the first tr of the tbody of the provided HTML node.
// If no such data element is found, the function returns an empty string.
func findWindowsBuild(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		// Iterate over tbody children
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			// Check if the child is a table row
			if child.Type == html.ElementNode && child.Data == "tr" {
				// Count the number of table data elements in this row
				tdCount := 0
				for td := child.FirstChild; td != nil; td = td.NextSibling {
					if td.Type == html.ElementNode && td.Data == "td" {
						// Increment td count
						tdCount++
						if tdCount == 5 { // Fifth data element
							// Extract and return the data
							return strings.TrimSpace(td.FirstChild.Data)
						}
					}
				}
			}
		}
	}

	// If data element not found yet, continue searching recursively
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findWindowsBuild(c); result != "" {
			return result
		}
	}

	return ""
}
