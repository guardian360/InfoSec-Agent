package firefox

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/browserutils"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryFirefox inspects the user's browsing history in the Firefox browser for any visits to known phishing domains within the last week.
//
// Parameters: None
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of phishing domains that the user has visited in the last week. Each domain is represented as a string that includes the domain name and the time of the visit. If no visits to phishing domains are found, the result will be a string stating "No phishing domains found in the last week".
//
// This function works by locating the Firefox profile directory and copying the places.sqlite database to a temporary location. It then opens this database and queries it for the URLs visited in the last week. It processes the results of the query by checking each URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the visit, and this string is added to the results. If any error occurs during this process, such as an error copying the file or querying the database, this error is returned as the result of the check.
func HistoryFirefox() checks.Check {
	var output []string
	ffdirectory, err := browserutils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.HistoryFirefoxID, "No firefox directory found", err)
	}

	// Copy the database, so we don't have problems with locked files
	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")
	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing file: ", err)
		}
	}(tempHistoryDbff)

	// Copy the database to a temporary location
	copyError := browserutils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff, nil, nil)
	if copyError != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, copyError)
	}

	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	defer closeDatabase(db)

	rows, err := queryDatabase(db)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	defer closeRows(rows)

	output, err = processQueryResults(rows)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	if output == nil {
		return checks.NewCheckResult(checks.HistoryFirefoxID, 0, "No phising domains found in the last week")
	}
	return checks.NewCheckResult(checks.HistoryFirefoxID, 1, output...)
}

// closeDatabase is a utility function that is utilized within the HistoryFirefox function.
// It is responsible for terminating the established database connection.
//
// Parameters:
//   - db (*sql.DB): Represents the active database connection that needs to be closed.
//
// This function does not return any value. However, if an error occurs during the closure of the database connection,
// it will be logged for debugging purposes.
func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing database: ", err)
	}
}

// queryDatabase is a helper function that is utilized within the HistoryFirefox function.
// It performs a query on the Firefox history database to retrieve the user's browsing history from the past week.
//
// Parameters:
//   - db (*sql.DB): Represents the active database connection that will be queried.
//
// Returns:
//   - *sql.Rows: The result set of the query, which includes the URLs visited by the user in the past week and their respective visit dates.
//   - error: An error object that encapsulates any error encountered during the execution of the query. If no error occurred, this will be nil.
//
// This function constructs a SQL query to select the URL and last visit date from the 'moz_places' table for entries where the last visit date is within the past week. The query is executed against the provided database connection, and the resulting rows are returned. If an error occurs during the execution of the query, it is returned as well.
func queryDatabase(db *sql.DB) (*sql.Rows, error) {
	lastWeek := time.Now().AddDate(0, 0, -7).UnixMicro()
	rows, err := db.Query(
		"SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC",
		lastWeek)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// closeRows is a utility function used within the HistoryFirefox function.
// It is responsible for closing the result set of a database query.
//
// Parameters:
//   - rows (*sql.Rows): Represents the result set of a database query that needs to be closed.
//
// This function does not return any value. However, if an error occurs during the closure of the result set,
// it will be logged for debugging purposes.
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing rows: ", err)
	}
}

// processQueryResults is a function used within the HistoryFirefox function.
// It processes the results of a database query and identifies any visited phishing domains from the past week.
//
// Parameters:
//   - rows (*sql.Rows): Represents the result set of a database query, which includes the URLs visited by the user in the past week and their respective visit dates.
//
// Returns:
//   - []string: A slice of strings, where each string represents a phishing domain that the user has visited in the past week, along with the time of the visit. If no visits to phishing domains are found, the slice will be empty.
//   - error: An error object that encapsulates any error encountered during the processing of the query results. If no error occurred, this will be nil.
//
// This function iterates over each row in the provided result set. For each row, it extracts the URL and last visit date, and checks the URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the visit, and this string is added to the results. If an error occurs during this process, it is returned as well.
func processQueryResults(rows *sql.Rows) ([]string, error) {
	var output []string
	phishingDomainList := browserutils.GetPhishingDomains()
	re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)

	for rows.Next() {
		var url string
		var lastVisitDate sql.NullInt64
		if err := rows.Scan(&url, &lastVisitDate); err != nil {
			return nil, err
		}

		timeString := formatTime(lastVisitDate)

		matches := re.FindStringSubmatch(url)
		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				output = append(output, domain+timeString)
			}
		}
	}

	return output, nil
}

// formatTime is a utility function used within the processQueryResults function.
// It converts the last visit date of a website from a Unix timestamp to a human-readable string.
//
// Parameters:
//   - lastVisitDate (sql.NullInt64): Represents the Unix timestamp of the last visit date of a website. This value can be null.
//
// Returns:
//   - string: A human-readable string representing the last visit date of a website. If the provided timestamp is valid, the string will represent this timestamp. If the timestamp is not valid (i.e., null), the string will represent the current time.
//
// This function first checks if the provided timestamp is valid. If it is, the function converts the timestamp to a time.Time object and then formats this object to a string. If the timestamp is not valid, the function gets the current time, formats it to a string, and returns this string.
func formatTime(lastVisitDate sql.NullInt64) string {
	var lastVisitDateInt64 int64
	if lastVisitDate.Valid {
		lastVisitDateInt64 = lastVisitDate.Int64
	} else {
		lastVisitDateInt64 = -1 // Default value
	}

	if lastVisitDateInt64 > 0 {
		timeofCreation := time.UnixMicro(lastVisitDateInt64)
		return timeofCreation.String()
	}
	var timeNow = time.Now()
	return timeNow.String()
}
