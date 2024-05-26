package chromium

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryChromium checks the browsing history in a specified Chromium-based browser for visits to phishing domains.
//
// Parameters:
//   - browser: A string representing the name of the Chromium-based browser to check. This could be "Chrome", "Edge", etc.
//   - getter: A function that returns the default directory for the browser's history database.
//   - getterCopyDb: A function that copies the browser's history database to a temporary location.
//   - getterQDb: A function that queries the browser's history database for URLs visited in the last week.
//   - getterQR: A function that processes the query results to check for visits to phishing domains.
//   - phishingGetter: A function that returns the list of known phishing domains.
//
// Returns:
//   - A checks.Check object representing the result of the check. If phishing domains are found in the browsing history, the result will contain the domains visited in the last week and the timestamps of the visits. If no phishing domains are found, the result will indicate that no phishing domains were found in the last week.
//
// This function works by accessing the browser's history database, querying for URLs visited in the last week, and checking each URL against a list of known phishing domains. The function uses the utils.GetPhishingDomains helper function to fetch the list of known phishing domains.
func HistoryChromium(browser string, getter browsers.DefaultDirGetter, getterCopyDB CopyDBGetter, getterQDb QueryDatabaseGetter, getterQR ProcessQueryResultsGetter, phishingGetter browsers.PhishingDomainGetter) checks.Check {
	var results []string

	browserPath, returnID := GetBrowserPathAndIDHistory(browser)
	extensionsDir, err := getter.GetDefaultDir(browserPath)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}

	tempHistoryDB, err := getterCopyDB.CopyDatabase(extensionsDir + "/History")
	if err != nil {
		logger.Log.ErrorWithErr("Error copying database: ", err)
		return checks.NewCheckError(returnID, err)
	}
	defer os.Remove(tempHistoryDB)

	// Open the browser history database
	db, _ := sql.Open("sqlite", tempHistoryDB)

	defer CloseDatabase(db)

	rows, rowErr := getterQDb.QueryDatabase(db)
	if rowErr != nil {
		return checks.NewCheckError(returnID, rowErr)
	}
	defer CloseRows(rows)

	results, err = getterQR.ProcessQueryResults(rows, phishingGetter)
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}

	if len(results) > 0 {
		return checks.NewCheckResult(returnID, 0, strings.Join(results, "\n"))
	}
	return checks.NewCheckResult(returnID, 1)
}

// GetBrowserPathAndIDHistory is a function that takes a browser name as input,
// and returns the path to the browser's directory and the ID of the browser.
//
// Parameters:
//   - browser: A string representing the name of the browser. Currently, this function supports "Chrome" and "Edge".
//
// Returns:
//   - A string representing the path to the browser's directory.
//   - An integer representing the ID the check.
//
// This function works by checking the provided browser name against known browser names. If the browser is "Chrome",
// the function returns the path to the Chrome directory and the ID of Chrome. If the browser is "Edge",
// the function returns the path to the Edge directory and the ID of Edge. If the browser is unknown, the function
// returns an empty string and 0.
func GetBrowserPathAndIDHistory(browser string) (string, int) {
	if browser == browsers.Chrome {
		return browsers.ChromePath, checks.HistoryChromiumID
	}
	if browser == browsers.Edge {
		return browsers.EdgePath, checks.HistoryEdgeID
	}
	return "", 0
}

type CopyDBGetter interface {
	CopyDatabase(src string) (string, error)
}

type RealCopyDBGetter struct{}

// CopyDatabase creates a temporary copy of the database file at the given source path.
//
// Parameters:
//   - src: A string representing the path to the source database file.
//
// Returns:
//   - A string representing the path to the temporary copy of the database file.
//   - An error, which will be nil if the operation was successful. If an error occurs while copying the file, that error will be returned.
//
// This function works by joining the path to the system's temporary directory with the name of the temporary database file ("tempHistoryDB.sqlite"), and then calling the CopyFile function to copy the source file to the temporary location. If the CopyFile function returns an error, the function returns an empty string and the error. Otherwise, it returns the path to the temporary file and nil.
func (r RealCopyDBGetter) CopyDatabase(src string) (string, error) {
	tempDB := filepath.Join(os.TempDir(), "tempHistoryDB.sqlite")
	err := browsers.CopyFile(src, tempDB, nil, nil)
	if err != nil {
		return "", err
	}
	return tempDB, nil
}

// CloseDatabase safely closes the provided database connection.
//
// Parameters:
//   - db: A pointer to a sql.DB object representing the database connection to close.
//
// This function attempts to close the provided database connection and logs an error if the operation fails.
// It does not return any value.
func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing database: ", err)
	}
}

type QueryDatabaseGetter interface {
	QueryDatabase(db *sql.DB) (*sql.Rows, error)
}

type RealQueryDatabaseGetter struct{}

// QueryDatabase retrieves the browsing history from the provided database connection.
//
// Parameters:
//   - db: A pointer to a sql.DB object representing the database connection to query.
//
// Returns:
//   - A pointer to a sql.Rows object representing the result set of the query.
//   - An error, which will be nil if the operation was successful.
//
// This function queries the database for URLs visited in the last week, ordered by the last visit time in descending order. It returns the result set of the query and an error, if any occurred.
func (r RealQueryDatabaseGetter) QueryDatabase(db *sql.DB) (*sql.Rows, error) {
	oneWeekAgo := (time.Now().AddDate(369, 0, -7).UnixMicro() / 1000000000) * 1000000000
	rows, err := db.Query(
		"SELECT url, title, visit_count, last_visit_time FROM urls "+
			"WHERE last_visit_time > ? ORDER BY last_visit_time DESC", oneWeekAgo)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// CloseRows safely closes the provided sql.Rows object.
//
// Parameters:
//   - rows: A pointer to a sql.Rows object that needs to be closed.
//
// This function attempts to close the provided sql.Rows object and logs an error if the operation fails.
// It does not return any value.
func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing rows: ", err)
	}
}

type ProcessQueryResultsGetter interface {
	ProcessQueryResults(rows *sql.Rows, getter browsers.PhishingDomainGetter) ([]string, error)
}

type RealProcessQueryResultsGetter struct{}

// ProcessQueryResults processes the result set of a query to the browsing history database, checking each URL against a list of known phishing domains.
//
// Parameters:
//   - rows: A pointer to a sql.Rows object representing the result set of the query. Each row represents a URL visited in the last week, along with its title, visit count, and last visit time.
//   - getter: A function that returns the list of known phishing domains.
//
// Returns:
//   - A slice of strings, each string representing a visit to a known phishing domain. The string includes the domain name and the time of the last visit. If no visits to phishing domains are found, the slice will be empty.
//   - An error, which will be nil if the operation was successful. If an error occurs while scanning the rows or if the rows contain an error, that error will be returned.
//
// This function works by iterating over the rows, extracting the URL from each row, and checking the domain part of the URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the last visit, and this string is added to the results. The function uses the utils.GetPhishingDomains helper function to fetch the list of known phishing domains and a regular expression to extract the domain part of the URL.
func (r RealProcessQueryResultsGetter) ProcessQueryResults(rows *sql.Rows, getter browsers.PhishingDomainGetter) ([]string, error) {
	var results []string
	phishingDomainList, err := getter.GetPhishingDomains()
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?:https?://)?(?:[^@\n]+@)?(?:www\.)?([^:/\n?]+\.[^:/\n?]+)`)

	for rows.Next() {
		var url, title string
		var visitCount, lastVisitTime int
		scanErr := rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		if scanErr != nil {
			return nil, scanErr
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
