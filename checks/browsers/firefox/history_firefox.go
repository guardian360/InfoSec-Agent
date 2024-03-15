package firefox

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"os"
	"path/filepath"
	"regexp"
	"time"

	_ "modernc.org/sqlite"
)

func HistoryFirefox() checks.Check {
	var output []string
	ffdirectory, err := utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf("HistoryFirefox", "No firefox directory found", err)
	}

	//Copy the database so we don't have problems with locked files
	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")

	copyError := utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
	if copyError != nil {
		return checks.NewCheckError("HistoryFirefox", copyError)
	}
	//OpenDatabase
	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
		return checks.NewCheckError("HistoryFirefox", err)
	}
	defer db.Close()

	last30Days := time.Now().AddDate(0, 0, -30).UnixMicro()

	// Get the phishing domains from up-to-date github list
	phishingDomainList := utils.GetPhisingDomains()

	// Execute a query
	rows, err := db.Query("SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC", last30Days)
	if err != nil {
		return checks.NewCheckError("HistoryFirefox", err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var url string
		var lastVisitDate sql.NullInt64
		// Scan the row into variables
		if err := rows.Scan(&url, &lastVisitDate); err != nil {
			return checks.NewCheckError("HistoryFirefox", err)
		}
		var timeString = ""
		// Check if the lastVisitDate is NULL
		var lastVisitDateInt64 int64
		if lastVisitDate.Valid {
			lastVisitDateInt64 = lastVisitDate.Int64
		} else {
			lastVisitDateInt64 = -1 // Or any other default value you prefer
		}
		if lastVisitDateInt64 > 0 {
			timeofCreation := time.UnixMicro(lastVisitDateInt64)
			timeString = timeofCreation.String()
		} else {
			var time = time.Now()
			timeString = time.String()
		}

		//We only want to print the url to map against untrustworthy domains so we use the following regex to extract the domain
		re := regexp.MustCompile(`(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+\.[^:\/\n?]+)`)
		matches := re.FindStringSubmatch(url)

		for _, scamDomain := range phishingDomainList {
			if len(matches) > 1 && matches[1] == scamDomain {
				domain := matches[1]
				output = append(output, domain+timeString)
			}
		}
	}
	os.Remove(tempHistoryDbff)
	return checks.NewCheckResult("HistoryFirefox", output...)
}
