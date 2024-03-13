package chrome

import (
	"InfoSec-Agent/checks"
	utils "InfoSec-Agent/utils"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func HistoryChrome() checks.Check {
	// List of the results, this will contain a list of domains which are known to be phishing domains.
	var results []string

	// Get the user's home directory
	user, err := os.UserHomeDir()
	if err != nil {
		return checks.NewCheckErrorf("HistoryChrome", "Error: ", err)
	}

	//Copy the database so we don't have problems with locked files
	tempHistoryDb := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")
	// Copy the database to a temporary location
	copyError := utils.CopyFile(user+"/AppData/Local/Google/Chrome/User Data/Default/History", tempHistoryDb)
	if copyError != nil {
		return checks.NewCheckError("HistoryChrome", copyError)
	}

	// Open the Chrome history database
	db, err := sql.Open("sqlite", tempHistoryDb)
	if err != nil {
		return checks.NewCheckError("HistoryChrome", err)
	}
	defer db.Close()

	// Get the time one week ago (with Chrome the counting starts from 1601-01-01)
	oneWeekAgo := time.Now().AddDate(369, 0, -7).UnixMicro()
	// Query the urls table we limit the results to 50 for the time being
	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls WHERE last_visit_time > ? ORDER BY last_visit_time DESC", oneWeekAgo)
	if err != nil {
		return checks.NewCheckError("HistoryChrome", err)
	}
	defer rows.Close()

	// Get the phishing domains from up-to-date github list
	phishingDomainList := utils.GetPhisingDomains()
	// Iterate over the results and print them
	for rows.Next() {
		var url, title string
		var visitCount, lastVisitTime int
		err = rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		if err != nil {
			return checks.NewCheckError("HistoryChrome", err)
		}
		//We only want to print the url to map against untrustworthy domains so we use the following regex to extract the domain
		re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)
		matches := re.FindStringSubmatch(url)

		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				results = append(results, "You visited website: "+domain+" which is a known phishing domain. The time of the last visit: "+time.UnixMicro(int64(lastVisitTime)).AddDate(-369, 0, 0).String())
			}
		}
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return checks.NewCheckError("HistoryChrome", err)
	}
	defer os.Remove(tempHistoryDb) // Clean up the temporary file when done

	// Print the results
	if len(results) > 0 {
		return checks.NewCheckResult("HistoryChrome", strings.Join(results, "\n"))
	} else {
		return checks.NewCheckResult("HistoryChrome", "No phishing domains found in the last week")
	}
}
