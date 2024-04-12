package firefox

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/utils"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryFirefox checks the user's history in the Firefox browser for phishing domains.
//
// Parameters:
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
func HistoryFirefox(profileFinder utils.FirefoxProfileFinder) checks.Check {
	var output []string
	ffdirectory, err := profileFinder.FirefoxFolder()
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
	copyError := utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
	if copyError != nil {
		return checks.NewCheckError(checks.HistoryFirefoxID, copyError)
	}

	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
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

// CloseDatabase is a helper function used by the HistoryFirefox function.
// It closes the database connection passed to it.
//
// Parameters: db (*sql.DB) - the database connection to close
//
// Returns: _
func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Log.ErrorWithErr("Error closing database: ", err)
	}
}

// QueryDatabase is a helper function used by the HistoryFirefox function.
// It queries the Firefox history database for the user's history in the last week.
//
// Parameters: db (*sql.DB) - the database connection to query
//
// Returns: The QueryResults in a slice and an error if applicable
func QueryDatabase(db *sql.DB) ([]QueryResult, error) {
	// Truncate the time to the time of the last week without the milliseconds
	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000
	rows, err := db.Query(
		"SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC",
		lastWeek)
	if err != nil {
		return nil, err
	}

	defer CloseRows(rows)

	var results []QueryResult
	for rows.Next() {
		var result QueryResult
		if err := rows.Scan(&result.URL, &result.LastVisitDate); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// CloseRows is a helper function used by the HistoryFirefox function.
// It closes the rows of the query passed to it.
//
// Parameters: rows (*sql.Rows) - the rows to close
//
// Returns: _
func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logger.Log.ErrorWithErr("Error closing rows: ", err)
	}
}

// QueryResult is a struct used by the HistoryFirefox function to store the results of the query.
type QueryResult struct {
	URL           string
	LastVisitDate sql.NullInt64
}

// ProcessQueryResults is a helper function used by the HistoryFirefox function.
// It processes the results of the query and returns the phishing domains that the user has visited in the last week.
//
// Parameters: QueryResult ([]QueryResult) - A slice of QueryResults
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
// and an error if applicable
func ProcessQueryResults(results []QueryResult) ([]string, error) {
	var output []string
	phishingDomainList := utils.GetPhishingDomains()
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
// Parameters: lastVisitDate (sql.NullInt64) - the last visit date of a website
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
