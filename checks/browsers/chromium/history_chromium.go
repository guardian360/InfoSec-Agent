// Package chromium is responsible for running checks on Chromium based browsers.
//
// Exported function(s): ExtensionsChromium, HistoryChromium
package chromium

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// HistoryChromium checks the user's history in a Chromium based browser for phishing domains.
//
// Parameters:
//
//	browser (string) - The name of the browser to check
//
// Returns: If the user has visited a phishing domain in the last week
func HistoryChromium(browser string) checks.Check {
	// List of the results, this will contain a list of domains which are known to be phishing domains.
	var results []string
	var browserPath string
	var returnBrowserName string
	// Set the browser path and the return browser name based on the browser to check
	// Currently, supports checking of Google Chrome and Microsoft Edge
	if browser == "Chrome" {
		returnBrowserName = "HistoryChrome"
		browserPath = "Google/Chrome"
	}
	if browser == "Edge" {
		returnBrowserName = "HistoryEdge"
		browserPath = "Microsoft/Edge"
	}

	// Get the current user's home directory, where the history can be found
	user, err := os.UserHomeDir()
	if err != nil {
		return checks.NewCheckErrorf(returnBrowserName, "Error: ", err)
	}

	// Copy the database, so problems don't arise when the file gets locked
	tempHistoryDb := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")

	// Clean up the temporary file when the function returns
	defer os.Remove(tempHistoryDb)

	// Copy the database to a temporary location
	copyError := utils.CopyFile(user+"/AppData/Local/"+browserPath+"/User Data/Default/History", tempHistoryDb)
	if copyError != nil {
		return checks.NewCheckError(returnBrowserName, copyError)
	}

	// Open the browser history database
	db, err := sql.Open("sqlite", tempHistoryDb)
	if err != nil {
		return checks.NewCheckError(returnBrowserName, err)
	}
	defer db.Close()

	// Get the time of one week ago (with Chrome the counting starts from 1601-01-01)
	oneWeekAgo := time.Now().AddDate(369, 0, -7).UnixMicro()
	// Query the history table
	// We limit the results to 50 for the time being
	rows, err := db.Query(
		"SELECT url, title, visit_count, last_visit_time FROM urls WHERE last_visit_time > ? ORDER BY last_visit_time DESC", oneWeekAgo)
	if err != nil {
		return checks.NewCheckError(returnBrowserName, err)
	}
	defer rows.Close()

	// Get the phishing domains from an up-to-date GitHub list
	phishingDomainList := utils.GetPhisingDomains()
	// Compare the visited domains with the phishing domains
	for rows.Next() {
		var url, title string
		var visitCount, lastVisitTime int
		err = rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		if err != nil {
			return checks.NewCheckError(returnBrowserName, err)
		}
		// The following regex is used to extract the domain from the url,
		// to use for mapping against the phishing domains
		re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)
		matches := re.FindStringSubmatch(url)

		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				// Return which domain is associated with phishing and when it was visited
				results = append(results, "You visited website: "+domain+" which is a known phishing domain. "+
					"The time of the last visit: "+
					""+time.UnixMicro(int64(lastVisitTime)).AddDate(-369, 0, 0).String())
			}
		}
	}

	// Check for errors from iterating over the database rows
	if err = rows.Err(); err != nil {
		return checks.NewCheckError(returnBrowserName, err)
	}

	// Return the domains found, and when they were visited
	if len(results) > 0 {
		return checks.NewCheckResult(returnBrowserName, strings.Join(results, "\n"))
	} else {
		return checks.NewCheckResult(returnBrowserName, "No phishing domains found in the last week")
	}
}
