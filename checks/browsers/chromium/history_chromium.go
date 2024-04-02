package chromium

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryChromium checks the user's history in a Chromium based browser for phishing domains.
//
// Parameters:
//
//	browser (string) - The name of the browser to check
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
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
	tempHistoryDB := filepath.Join(os.TempDir(), "tempHistoryDB.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("error removing file: ", err)
		}
	}(tempHistoryDB)

	// Copy the database to a temporary location
	copyError := utils.CopyFile(user+"/AppData/Local/"+browserPath+"/User Data/Default/History", tempHistoryDB)
	if copyError != nil {
		return checks.NewCheckError(returnBrowserName, copyError)
	}

	// Open the browser history database
	db, err := sql.Open("sqlite", tempHistoryDB)
	if err != nil {
		return checks.NewCheckError(returnBrowserName, err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("error closing database: ", err)
		}
	}(db)

	// Get the time of one week ago (with Chrome the counting starts from 1601-01-01)
	oneWeekAgo := time.Now().AddDate(369, 0, -7).UnixMicro()
	// Query the history table
	// We limit the results to 50 for the time being
	rows, err := db.Query(
		"SELECT url, title, visit_count, last_visit_time FROM urls "+
			"WHERE last_visit_time > ? ORDER BY last_visit_time DESC", oneWeekAgo)
	if err != nil {
		return checks.NewCheckError(returnBrowserName, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("error closing rows: ", err)
		}
	}(rows)

	// Get the phishing domains from an up-to-date GitHub list
	phishingDomainList := utils.GetPhishingDomains()
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
	}
	return checks.NewCheckResult(returnBrowserName, "No phishing domains found in the last week")
}
