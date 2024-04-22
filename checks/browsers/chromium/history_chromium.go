package chromium

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browser_utils"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"

	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryChromium checks the browsing history in a specified Chromium-based browser for visits to phishing domains.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to check. This could be "Chrome", "Edge", etc.
//
// Returns:
//   - A checks.Check object representing the result of the check. If phishing domains are found in the browsing history, the result will contain the domains visited in the last week and the timestamps of the visits. If no phishing domains are found, the result will indicate that no phishing domains were found in the last week.
//
// This function works by accessing the browser's history database, querying for URLs visited in the last week, and checking each URL against a list of known phishing domains. The function uses the utils.GetPhishingDomains helper function to fetch the list of known phishing domains.
func HistoryChromium(browser string) checks.Check {
	var results []string
	var browserPath string
	var returnID int

	if browser == chrome {
		browserPath = chromePath
		returnID = checks.HistoryChromiumID
	}
	if browser == edge {
		browserPath = edgePath
		returnID = checks.HistoryEdgeID
	}
	// Get the current user's home directory, where the history can be found
	user, err := os.UserHomeDir()
	if err != nil {
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}

	// Copy the database, so problems don't arise when the file gets locked
	tempHistoryDB := filepath.Join(os.TempDir(), "tempHistoryDB.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing file: ", err)
		}
	}(tempHistoryDB)

	// Copy the database to a temporary location
	copyError := browser_utils.CopyFile(user+"/AppData/Local/"+browserPath+"/User Data/Default/History", tempHistoryDB, nil, nil)
	if copyError != nil {
		return checks.NewCheckError(returnID, copyError)
	}

	// Open the browser history database
	db, err := sql.Open("sqlite", tempHistoryDB)
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}
	defer closeDatabase(db)

	rows, err := queryDatabase(db)
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}
	defer closeRows(rows)

	results, err = processQueryResults(rows)
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}

	if len(results) > 0 {
		return checks.NewCheckResult(returnID, 0, strings.Join(results, "\n"))
	}
	return checks.NewCheckResult(returnID, 1, "No phishing domains found in the last week")
}

// closeDatabase safely closes the provided database connection.
//
// Parameters:
//   - db: A pointer to an sql.DB object representing the database connection to close.
//
// This function attempts to close the provided database connection and logs an error if the operation fails.
// It does not return any value.
func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing database: ", err)
	}
}

// queryDatabase retrieves the browsing history from the provided database connection.
//
// Parameters:
//   - db: A pointer to an sql.DB object representing the database connection to query.
//
// Returns:
//   - A pointer to an sql.Rows object representing the result set of the query.
//   - An error, which will be nil if the operation was successful.
//
// This function queries the database for URLs visited in the last week, ordered by the last visit time in descending order. It returns the result set of the query and an error, if any occurred.
func queryDatabase(db *sql.DB) (*sql.Rows, error) {
	oneWeekAgo := time.Now().AddDate(369, 0, -7).UnixMicro()
	rows, err := db.Query(
		"SELECT url, title, visit_count, last_visit_time FROM urls "+
			"WHERE last_visit_time > ? ORDER BY last_visit_time DESC", oneWeekAgo)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// closeRows safely closes the provided sql.Rows object.
//
// Parameters:
//   - rows: A pointer to an sql.Rows object that needs to be closed.
//
// This function attempts to close the provided sql.Rows object and logs an error if the operation fails.
// It does not return any value.
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing rows: ", err)
	}
}

// processQueryResults processes the result set of a query to the browsing history database, checking each URL against a list of known phishing domains.
//
// Parameters:
//   - rows: A pointer to an sql.Rows object representing the result set of the query. Each row represents a URL visited in the last week, along with its title, visit count, and last visit time.
//
// Returns:
//   - A slice of strings, each string representing a visit to a known phishing domain. The string includes the domain name and the time of the last visit. If no visits to phishing domains are found, the slice will be empty.
//   - An error, which will be nil if the operation was successful. If an error occurs while scanning the rows or if the rows contain an error, that error will be returned.
//
// This function works by iterating over the rows, extracting the URL from each row, and checking the domain part of the URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the last visit, and this string is added to the results. The function uses the utils.GetPhishingDomains helper function to fetch the list of known phishing domains and a regular expression to extract the domain part of the URL.
func processQueryResults(rows *sql.Rows) ([]string, error) {
	var results []string
	phishingDomainList := browser_utils.GetPhishingDomains()
	re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)

	for rows.Next() {
		var url, title string
		var visitCount, lastVisitTime int
		err := rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		if err != nil {
			return nil, err
		}

		matches := re.FindStringSubmatch(url)
		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				results = append(results, "You visited website: "+domain+" which is a known phishing domain. "+
					"The time of the last visit: "+
					""+time.UnixMicro(int64(lastVisitTime)).AddDate(-369, 0, 0).String())
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
