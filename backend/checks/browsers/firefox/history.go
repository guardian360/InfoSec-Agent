package firefox

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryFirefox inspects the user's browsing history in the Firefox browser for any visits to known phishing domains within the last week.
//
// Parameters:   - profileFinder: An object that implements the FirefoxProfileFinder interface. It is used to find the Firefox profile directory.
//   - getter: An object that implements the PhishingDomainGetter interface. It is used to retrieve the list of known phishing domains.
//   - queryGetter: An object that implements the QueryDatabaseGetter interface. It is used to query the Firefox history database.
//   - processGetter: An object that implements the ProcessQueryResultsGetter interface. It is used to process the results of the database query.
//   - copyDBGetter: An object that implements the CopyDBGetter interface. It is used to copy the Firefox history database to a temporary location.
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
func HistoryFirefox(profileFinder browsers.FirefoxProfileFinder, getter browsers.PhishingDomainGetter, queryGetter QueryDatabaseGetter, processGetter ProcessQueryResultsGetter, copyDBGetter CopyDBGetter) checks.Check {
	var output []string
	ffDirectory, err := profileFinder.FirefoxFolder()
	if err != nil {
		logger.Log.ErrorWithErr("No firefox directory found", err)
		return checks.NewCheckErrorf(checks.HistoryFirefoxID, "No firefox directory found", err)
	}

	// Copy the database, so we don't have problems with locked files
	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDBFirefox.sqlite")
	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing file", err)
		}
	}(tempHistoryDbff)

	copyGetter := browsers.RealCopyFileGetter{}
	// Copy the database to a temporary location
	copyErr := copyDBGetter.CopyDatabase(copyGetter, ffDirectory[0], tempHistoryDbff)
	if copyErr != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, copyErr)
	}

	db, _ := sql.Open("sqlite", tempHistoryDbff)
	defer CloseDatabase(db)

	rows, err := queryGetter.QueryDatabase(db)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}

	output, err = processGetter.ProcessQueryResults(rows, getter)
	if err != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, err)
	}
	if output == nil {
		return checks.NewCheckResult(checks.HistoryFirefoxID, 0)
	}
	return checks.NewCheckResult(checks.HistoryFirefoxID, 1, output...)
}

// CopyDBGetter is an interface that wraps the CopyDatabase method.
// It provides a way to abstract the operation of copying the database, allowing for different implementations.
// This can be useful for testing, where a mock implementation can be used.
//
// The CopyDatabase method takes a CopyFileGetter interface, a string representing the Firefox directory,
// and a string representing the temporary history database as parameters.
// It returns an error if any occurs during the copy operation.
type CopyDBGetter interface {
	CopyDatabase(copyGetter browsers.CopyFileGetter, ffDirectory string, tempHistoryDbff string) error
}

// RealCopyDBGetter is a struct that implements the CopyDBGetter interface.
// It provides the real implementation of the CopyDatabase method.
type RealCopyDBGetter struct{}

// CopyDatabase is a method that copies the Firefox history database to a temporary location.
// It takes a CopyFileGetter interface, a string representing the Firefox directory,
// and a string representing the temporary history database as parameters.
// It returns an error if any occurs during the copy operation.
//
// The CopyFileGetter interface is used to copy the file.
// The Firefox directory string represents the directory where the Firefox history database is located.
// The temporary history database string represents the location where the database will be copied to.
//
// If an error occurs during the copy operation, it is logged and returned.
// If the copy operation is successful, nil is returned.
func (r RealCopyDBGetter) CopyDatabase(copyGetter browsers.CopyFileGetter, ffDirectory string, tempHistoryDbff string) error {
	// Copy the database to a temporary location
	copyError := copyGetter.CopyFile(ffDirectory+"\\places.sqlite", tempHistoryDbff, nil, nil)
	if copyError != nil {
		logger.Log.ErrorWithErr("Unable to make a copy of the file", copyError)
		return copyError
	}
	return nil
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
		logger.Log.ErrorWithErr("Error closing database", err)
	}
}

// QueryDatabaseGetter is an interface that wraps the QueryDatabase method.
// It provides a way to abstract the database querying operation, allowing for different implementations.
// This can be useful for testing, where a mock implementation can be used.
//
// The QueryDatabase method takes a pointer to an sql.DB object and returns a slice of QueryResult objects and an error.
// The sql.DB object represents the active database connection that will be queried.
// The slice of QueryResult objects represents the result set of the database query.
// The error represents any error encountered during the execution of the query. If no error occurred, this will be nil.
type QueryDatabaseGetter interface {
	QueryDatabase(db *sql.DB) ([]QueryResult, error)
}

// RealQueryDatabaseGetter is a struct that implements the QueryDatabaseGetter interface.
// It provides the real implementation of the QueryDatabase method.
type RealQueryDatabaseGetter struct{}

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
func (r RealQueryDatabaseGetter) QueryDatabase(db *sql.DB) ([]QueryResult, error) {
	// Truncate the time to the time of the last week without the milliseconds
	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000
	rows, err := db.Query(
		"SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC",
		lastWeek)
	if err != nil {
		logger.Log.ErrorWithErr("Error querying database", err)
		return nil, err
	}

	defer CloseRows(rows)

	var results []QueryResult
	for rows.Next() {
		var result QueryResult
		if err = rows.Scan(&result.URL, &result.LastVisitDate); err != nil {
			logger.Log.ErrorWithErr("Error scanning row", err)
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		logger.Log.ErrorWithErr("Error iterating over rows", err)
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
		logger.Log.ErrorWithErr("Error closing rows", err)
	}
}

// QueryResult is a struct used by the HistoryFirefox function to store the results of the query.
type QueryResult struct {
	URL           string
	LastVisitDate sql.NullInt64
}

// ProcessQueryResultsGetter is an interface that wraps the ProcessQueryResults method.
// It provides a way to abstract the operation of processing query results, allowing for different implementations.
// This can be useful for testing, where a mock implementation can be used.
//
// The ProcessQueryResults method takes a slice of QueryResult objects and a PhishingDomainGetter interface as parameters.
// It returns a slice of strings and an error.
// The slice of QueryResult objects represents the result set of a database query.
// The PhishingDomainGetter interface is used to retrieve the list of known phishing domains.
// The slice of strings represents the phishing domains that the user has visited in the past week, along with the time of the visit.
// The error represents any error encountered during the processing of the query results. If no error occurred, this will be nil.
type ProcessQueryResultsGetter interface {
	ProcessQueryResults(results []QueryResult, getter browsers.PhishingDomainGetter) ([]string, error)
}

// RealProcessQueryResultsGetter is a struct that implements the ProcessQueryResultsGetter interface.
// It provides the real implementation of the ProcessQueryResults method.
type RealProcessQueryResultsGetter struct{}

// ProcessQueryResults is a function used within the HistoryFirefox function.
// It processes the results of a database query and identifies any visited phishing domains from the past week.
//
// Parameters:
//   - results []QueryResult: Represents the result set of a database query, which includes the URLs visited by the user in the past week and their respective visit dates.
//   - getter browsers.PhishingDomainGetter: An object that implements the PhishingDomainGetter interface. It is used to retrieve the list of known phishing domains.
//
// Returns:
//   - []string: A slice of strings, where each string represents a phishing domain that the user has visited in the past week, along with the time of the visit. If no visits to phishing domains are found, the slice will be empty.
//   - error: An error object that encapsulates any error encountered during the processing of the query results. If no error occurred, this will be nil.
//
// This function iterates over each row in the provided result set. For each row, it extracts the URL and last visit date, and checks the URL against a list of known phishing domains. If a match is found, a string is generated that includes the domain name and the time of the visit, and this string is added to the results. If an error occurs during this process, it is returned as well.
func (r RealProcessQueryResultsGetter) ProcessQueryResults(results []QueryResult, getter browsers.PhishingDomainGetter) ([]string, error) {
	var output []string
	creator := browsers.RealRequestCreator{}
	phishingDomainList, err := getter.GetPhishingDomains(creator)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting phishing domains", err)
		return nil, err
	}

	for _, result := range results {
		timeString := FormatTime(result.LastVisitDate)

		for _, scamDomain := range phishingDomainList {
			if strings.Contains(result.URL, scamDomain) {
				output = append(output, result.URL+timeString)
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
