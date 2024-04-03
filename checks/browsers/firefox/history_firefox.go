package firefox

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	// Necessary to use the sqlite driver
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	_ "modernc.org/sqlite"
)

// HistoryFirefox checks the user's history in the Firefox browser for phishing domains.
//
// Parameters:
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
func HistoryFirefox() checks.Check {
	var output []string
	db, err := setupDatabase()
	if err != nil {
		return checks.NewCheckError("HistoryFirefox", err)
	}
	defer closeDatabase(db)

	rows, err := queryDatabase(db)
	if err != nil {
		return checks.NewCheckError("HistoryFirefox", err)
	}
	defer closeRows(rows)

	output, err = processQueryResults(rows)
	if err != nil {
		return checks.NewCheckError("HistoryFirefox", err)
	}

	return checks.NewCheckResult("HistoryFirefox", output...)
}

// setupDatabase is a helper function used by the HistoryFirefox function.
// It sets up a temporary copy of the Firefox history database for querying.
//
// Parameters: _
//
// Returns: A pointer to the database and an error if applicable
func setupDatabase() (*sql.DB, error) {
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return nil, err
	}

	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			log.Println("error removing file: ", err)
		}
	}(tempHistoryDbff)

	err = utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
		return nil, err
	}

	return db, nil
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

		timeString, err := formatTime(lastVisitDate)
		if err != nil {
			return nil, err
		}

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
// Returns: The last visit date of a website as a human-readable string and an error if applicable
func formatTime(lastVisitDate sql.NullInt64) (string, error) {
	var lastVisitDateInt64 int64
	if lastVisitDate.Valid {
		lastVisitDateInt64 = lastVisitDate.Int64
	} else {
		lastVisitDateInt64 = -1 // Default value
	}

	if lastVisitDateInt64 > 0 {
		timeofCreation := time.UnixMicro(lastVisitDateInt64)
		return timeofCreation.String(), nil
	} else {
		var timeNow = time.Now()
		return timeNow.String(), nil
	}
}

//func HistoryFirefox() checks.Check {
//	var output []string
//	ffdirectory, err := utils.FirefoxFolder()
//	if err != nil {
//		return checks.NewCheckErrorf("HistoryFirefox", "No firefox directory found", err)
//	}
//
//	// Copy the database, so we don't have problems with locked files
//	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")
//	// Clean up the temporary file when the function returns
//	defer func(name string) {
//		err = os.Remove(name)
//		if err != nil {
//			log.Println("error removing file: ", err)
//		}
//	}(tempHistoryDbff)
//
//	// Copy the database to a temporary location
//	copyError := utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
//	if copyError != nil {
//		return checks.NewCheckError("HistoryFirefox", copyError)
//	}
//
//	db, err := sql.Open("sqlite", tempHistoryDbff)
//	if err != nil {
//		return checks.NewCheckError("HistoryFirefox", err)
//	}
//	defer func(db *sql.DB) {
//		err = db.Close()
//		if err != nil {
//			log.Println("error closing database: ", err)
//		}
//	}(db)
//
//	lastWeek := time.Now().AddDate(0, 0, -7).UnixMicro()
//
//	// Get the phishing domains from up-to-date GitHub list
//	phishingDomainList := utils.GetPhishingDomains()
//
//	// Query the urls and when the sites were visited from the history database
//	rows, err := db.Query(
//		"SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC",
//		lastWeek)
//	// TODO: check if this is error handling is correct
//	if rows.Err() != nil {
//		return checks.NewCheckError("HistoryFirefox", rows.Err())
//	}
//	if err != nil {
//		return checks.NewCheckError("HistoryFirefox", err)
//	}
//	defer func(rows *sql.Rows) {
//		err = rows.Close()
//		if err != nil {
//			log.Println("error closing rows: ", err)
//		}
//	}(rows)
//
//	// Iterate over each found url
//	for rows.Next() {
//		var url string
//		var lastVisitDate sql.NullInt64
//		// Scan the row into variables
//		if err = rows.Scan(&url, &lastVisitDate); err != nil {
//			return checks.NewCheckError("HistoryFirefox", err)
//		}
//		var timeString string
//		// Check if the lastVisitDate is nil
//		var lastVisitDateInt64 int64
//		if lastVisitDate.Valid {
//			lastVisitDateInt64 = lastVisitDate.Int64
//		} else {
//			lastVisitDateInt64 = -1 // Default value
//		}
//		if lastVisitDateInt64 > 0 {
//			timeofCreation := time.UnixMicro(lastVisitDateInt64)
//			timeString = timeofCreation.String()
//		} else {
//			var time = time.Now()
//			timeString = time.String()
//		}
//
//		// The following regex is used to extract the domain from the url,
//		// to use for mapping against the phishing domains
//		re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)
//		matches := re.FindStringSubmatch(url)
//
//		for _, scamDomain := range phishingDomainList {
//			if len(matches) > 1 && matches[1] == scamDomain {
//				domain := matches[1]
//				output = append(output, domain+timeString)
//			}
//		}
//	}
//	return checks.NewCheckResult("HistoryFirefox", output...)
//}
