package firefox

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

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
func HistoryFirefox() checks.Check {
	var output []string
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(30, "No firefox directory found", err)
	}

	// Copy the database, so we don't have problems with locked files
	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")
	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			log.Println("error removing file: ", err)
		}
	}(tempHistoryDbff)

	// Copy the database to a temporary location
	copyError := utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
	if copyError != nil {
		return checks.NewCheckError(30, copyError)
	}

	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
		return checks.NewCheckError(30, err)
	}
	defer closeDatabase(db)

	rows, err := queryDatabase(db)
	if err != nil {
		return checks.NewCheckError(30, err)
	}
	defer closeRows(rows)

	output, err = processQueryResults(rows)
	if err != nil {
		return checks.NewCheckError(30, err)
	}
	if output == nil {
		return checks.NewCheckResult(30, 0, "No phising domains found in the last week")
	}
	return checks.NewCheckResult(30, 1, output...)
}

// closeDatabase is a helper function used by the HistoryFirefox function.
// It closes the database connection passed to it.
//
// Parameters: db (*sql.DB) - the database connection to close
//
// Returns: _
func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("error closing database: ", err)
	}
}

// queryDatabase is a helper function used by the HistoryFirefox function.
// It queries the Firefox history database for the user's history in the last week.
//
// Parameters: db (*sql.DB) - the database connection to query
//
// Returns: The rows of the query and an error if applicable
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

// closeRows is a helper function used by the HistoryFirefox function.
// It closes the rows of the query passed to it.
//
// Parameters: rows (*sql.Rows) - the rows to close
//
// Returns: _
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("error closing rows: ", err)
	}
}

// processQueryResults is a helper function used by the HistoryFirefox function.
// It processes the results of the query and returns the phishing domains that the user has visited in the last week.
//
// Parameters: rows (*sql.Rows) - the rows of the query
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
// and an error if applicable
func processQueryResults(rows *sql.Rows) ([]string, error) {
	var output []string
	phishingDomainList := utils.GetPhishingDomains()
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

// formatTime is a helper function used by the processQueryResults function.
// It formats the last visit date of a website to a human-readable string.
//
// Parameters: lastVisitDate (sql.NullInt64) - the last visit date of a website
//
// Returns: The last visit date of a website as a human-readable string
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
