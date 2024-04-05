package chromium

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/InfoSec-Agent/InfoSec-Agent/utils"

	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"

	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// HistoryChromium checks the user's history in a Chromium based browser for phishing domains.
//
// Parameters:
//
// browser (string) - The name of the browser to check
//
// Returns: The phishing domains that the user has visited in the last week and when they visited it
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
			log.Println("error removing file: ", err)
		}
	}(tempHistoryDB)

	// Copy the database to a temporary location
	copyError := utils.CopyFile(user+"/AppData/Local/"+browserPath+"/User Data/Default/History", tempHistoryDB)
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

func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("error closing database: ", err)
	}
}

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

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("error closing rows: ", err)
	}
}

func processQueryResults(rows *sql.Rows) ([]string, error) {
	var results []string
	phishingDomainList := utils.GetPhishingDomains()
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
