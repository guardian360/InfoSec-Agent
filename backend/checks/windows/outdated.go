package windows

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"golang.org/x/net/html"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

var WinVersion int

// Outdated is a function that checks if the currently installed Windows version is outdated.
//
// Parameters:
//   - mockOS mocking.WindowsVersion: A mock object for retrieving the Windows version information.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the Windows version is up-to-date or if updates are available.
//
// The function works by retrieving the Windows version information using the provided mock object. It then compares the build number of the installed Windows version with the build numbers of the latest Windows 10 and Windows 11 versions. If the installed version's build number matches the latest build number for its major version (10 or 11), the function returns a message indicating that the Windows version is up-to-date. If the build number does not match, the function returns a message indicating that updates are available. If the major version is neither 10 nor 11, the function returns a message suggesting to update to Windows 10 or Windows 11.
func Outdated(mockExecutor mocking.CommandExecutor) checks.Check {
	versionData, err := mockExecutor.Execute("cmd", "/c", "ver")
	if err != nil {
		logger.Log.ErrorWithErr("Error executing command", err)
		return checks.NewCheckError(checks.WindowsOutdatedID, err)
	}

	versionString := string(versionData)
	// Using regular expression to extract numbers after the second dot
	re := regexp.MustCompile(`(\d+)\.\d+\.(\d+)\.(\d+)`)
	match := re.FindStringSubmatch(versionString)
	if len(match) < 4 {
		logger.Log.Error("Error parsing Windows version string")
		return checks.NewCheckError(checks.WindowsOutdatedID, errors.New("error parsing Windows version string"))
	}

	majorVersion, err := strconv.Atoi(match[1])
	if err != nil {
		logger.Log.ErrorWithErr("Error converting major version to integer", err)
		return checks.NewCheckError(checks.WindowsOutdatedID, err)
	}

	buildNumber := match[3]
	minorVersion, err := strconv.Atoi(match[2])
	WinVersion = findWindowsVersion(majorVersion, minorVersion)
	if err != nil {
		logger.Log.ErrorWithErr("Error converting minor version to integer", err)
		return checks.NewCheckError(checks.WindowsOutdatedID, err)
	}

	winVer := strconv.Itoa(minorVersion) + "." + buildNumber

	const win10Url = "https://learn.microsoft.com/en-us/windows/release-health/release-information"
	const win11Url = "https://learn.microsoft.com/en-us/windows/release-health/windows11-release-information"
	win10HTML := GetURLBody(win10Url)
	if win10HTML == nil {
		logger.Log.Error(
			"Error fetching Windows 10 HTML content, this function requires an internet connection")
		return checks.NewCheckError(checks.WindowsOutdatedID, errors.New(
			"error fetching Windows 10 HTML content,this function requires an internet connection"))
	}
	latestWin10Build := FindWindowsBuild(win10HTML)

	win11HTML := GetURLBody(win11Url)
	if win11HTML == nil {
		logger.Log.Error(
			"Error fetching Windows 11 HTML content, this function requires an internet connection")
		return checks.NewCheckError(checks.WindowsOutdatedID, errors.New(
			"error fetching Windows 11 HTML content, this function requires an internet connection"))
	}
	latestWin11Build := FindWindowsBuild(win11HTML)

	// Depending on the major Windows version (10 or 11), act accordingly
	switch {
	case findWindowsVersion(majorVersion, minorVersion) == 11:
		version := []string{"11", "WinVersion: " + winVer}
		if winVer == latestWin11Build {
			return checks.NewCheckResult(checks.WindowsOutdatedID, 0, "11")
		} else {
			return checks.NewCheckResult(checks.WindowsOutdatedID, 1, version...)
		}
	case findWindowsVersion(majorVersion, minorVersion) == 10:
		version := []string{"10", "WinVersion: " + winVer}
		if winVer == latestWin10Build {
			return checks.NewCheckResult(checks.WindowsOutdatedID, 0, "10")
		} else {
			return checks.NewCheckResult(checks.WindowsOutdatedID, 1, version...)
		}
	default:
		return checks.NewCheckResult(checks.WindowsOutdatedID, 2)
	}
}

// GetURLBody fetches and parses the HTML content of a given URL.
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
func GetURLBody(urlStr string) *html.Node {
	// Make HTTP GET request
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		logger.Log.ErrorWithErr("Error creating HTTP request", err)
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting response", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing response body", err)
		}
	}(resp.Body)

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Log.ErrorWithErr("Error parsing HTML", err)
	}
	return doc
}

// FindWindowsBuild searches for the latest Windows build in the HTML content of a given URL.
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
func FindWindowsBuild(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		// Iterate over tbody children
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			// Check if the child is a table row
			if child.Type == html.ElementNode && child.Data == "tr" {
				tdData := checkTableTDs(child)
				if tdData != "" {
					return tdData
				}
			}
		}
	}

	// If data element not found yet, continue searching recursively
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindWindowsBuild(c); result != "" {
			return result
		}
	}

	return ""
}

// checkTableTDs is a helper function that iterates over the children of a given HTML node,
// looking for 'td' elements. It counts the number of 'td' elements and when it finds the fifth one,
// it extracts and returns the data as a string.
//
// Parameters:
//   - child *html.Node: The HTML node to search for 'td' elements.
//
// Returns: The data from the fifth 'td' element in the provided HTML node.
// If no such 'td' element is found, the function returns an empty string.
func checkTableTDs(child *html.Node) string {
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
	return ""
}

// findWindowsVersion is a function that determines the Windows version based on the major and minor version numbers.
// It uses the versioning scheme of Windows 10 and Windows 11, where Windows 11 is identified by a minor version number of 22000 or higher.
//
// Parameters:
//   - majorVersion int: The major version number of the Windows OS. For Windows 10 and 11, this should be 10.
//   - minorVersion int: The minor version number of the Windows OS. For Windows 11, this should be 22000 or higher.
//
// Returns:
//   - int: The identified Windows version. This will be 11 for Windows 11, 10 for Windows 10, and 0 for any other version.
func findWindowsVersion(majorVersion int, minorVersion int) int {
	// If the minor version is 22000 or higher and the major version is 10, this is Windows 11
	if minorVersion >= 22000 && majorVersion == 10 {
		return 11
	}
	// If the major version is 10, and it's not Windows 11, it must be Windows 10
	if majorVersion == 10 {
		return 10
	}
	// If it's neither Windows 10 nor 11, return 0
	return 0
}
