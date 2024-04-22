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
// Returns: The phishing domains that the user has visited in the last week and when they visited it
func HistoryFirefox(profileFinder browserutils.FirefoxProfileFinder) checks.Check {
	var output []string
	ffdirectory, err := browserutils.RealProfileFinder{}.FirefoxFolder()
	if err != nil {
		logger.Log.ErrorWithErr("No firefox directory found: ", err)
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
		logger.Log.ErrorWithErr("Unable to make a copy of the file: ", copyError)
		return checks.NewCheckError(checks.HistoryFirefoxID, copyError)
	}

	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil { // Returns an error if a file is not found automatically check not needed?
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	defer CloseDatabase(db)

	rows, err := QueryDatabase(db)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}

	output, err = ProcessQueryResults(rows)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	if output == nil {
		return checks.NewCheckResult(checks.HistoryFirefoxID, 0, "No phishing domains found in the last week")
	}
	return checks.NewCheckResult(checks.HistoryFirefoxID, 1, output...)
}

// CloseDatabase is a utility function that is utilized within the HistoryFirefox function.
// It is responsible for terminating the established database connection.
//
// Parameters:
//   - db (*sql.DB): Represents the active database connection that needs to be closed.
//
// Returns: _
func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing database: ", err)
	}
}

// QueryDatabase is a helper function that is utilized within the HistoryFirefox function.
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
func QueryDatabase(db *sql.DB) ([]QueryResult, error) {
	// Truncate the time to the time of the last week without the milliseconds
	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000
	rows, err := db.Query(
		"SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC",
		lastWeek)
	if err != nil {
		logger.Log.ErrorWithErr("Error querying database: ", err)
		return nil, err
	}

	defer CloseRows(rows)

	var results []QueryResult
	for rows.Next() {
		var result QueryResult
		if err = rows.Scan(&result.URL, &result.LastVisitDate); err != nil {
			logger.Log.ErrorWithErr("Error scanning row: ", err)
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		logger.Log.ErrorWithErr("Error iterating over rows: ", err)
		return nil, err
	}

	return results, nil
}

// CloseRows is a utility function used within the HistoryFirefox function.
// It is responsible for closing the result set of a database query.
//
// Parameters:
//   - rows (*sql.Rows): Represents the result set of a database query that needs to be closed.
//
// This function does not return any value. However, if an error occurs during the closure of the result set,
// it will be logged for debugging purposes.
func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing rows: ", err)
	}
}

// QueryResult is a struct used by the HistoryFirefox function to store the results of the query.
type QueryResult struct {
	URL           string
	LastVisitDate sql.NullInt64
}

// ProcessQueryResults is a function used within the HistoryFirefox function.
// It processes the results of a database query and identifies any visited phishing domains from the past week.
//
// Parameters:
//   - results []QueryResult: Represents the result set of a database query, which includes the URLs visited by the user in the past week and their respective visit dates.
//
// Returns:
//   - []string: A slice of strings, where each string represents a phishing domain that the user has visited in the past week, along with the time of the visit. If no visits to phishing domains are found, the slice will be empty.
//   - error: An error object that encapsulates any error encountered during the processing of the query results. If no error occurred, this will be nil.
//
// This function iterates over each row in the provided result set. For each row, it extracts the URL and last visit date, and checks the URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the visit, and this string is added to the results. If an error occurs during this process, it is returned as well.
func ProcessQueryResults(results []QueryResult) ([]string, error) {
	var output []string
	phishingDomainList := browserutils.GetPhishingDomains()
	re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)

	for _, result := range results {
		timeString := FormatTime(result.LastVisitDate)

		matches := re.FindStringSubmatch(result.URL)
		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				output = append(output, domain+timeString)
			}
		}
	}

	return output, nil
}

// TimeFormatter is an interface that wraps the FormatTime function.
type TimeFormatter interface {
	FormatTime(lastVisitDate sql.NullInt64) string
}

// RealTimeFormatter is a type that implements TimeFormatter using the real FormatTime function.
type RealTimeFormatter struct{}

// FormatTime  is a helper function used by the processQueryResults function.
// It formats the last visit date of a website to a human-readable string.
//
// Parameters:
//   - lastVisitDate (sql.NullInt64): Represents the Unix timestamp of the last visit date of a website. This value can be null.
//
// Returns: The last visit date of a website as a human-readable string
func (RealTimeFormatter) FormatTime(lastVisitDate sql.NullInt64) string {
	if lastVisitDate.Valid && lastVisitDate.Int64 > 0 {
		return time.UnixMicro(lastVisitDate.Int64).String()
	}
	return time.UnixMicro(0).String()
}

// TimeFormat timeFormatter is the TimeFormatter currently in use. It can be reassigned to a mock in tests.
var TimeFormat TimeFormatter = RealTimeFormatter{}

// FormatTime is a function that calls the FormatTime method of the current TimeFormatter.
func FormatTime(lastVisitDate sql.NullInt64) string {
	return TimeFormat.FormatTime(lastVisitDate)
}
