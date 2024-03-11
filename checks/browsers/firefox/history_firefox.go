package checks

import (
	utils "InfoSec-Agent/utils"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func history_firefox() ([]string, error) {
	var output []string
	ffdirectory, _ := utils.FirefoxFolder()

	//Copy the database so we don't have problems with locked files
	tempHistoryDbff := filepath.Join(os.TempDir(), "tempHistoryDb.sqlite")

	copyError := utils.CopyFile(ffdirectory[0]+"\\places.sqlite", tempHistoryDbff)
	if copyError != nil {
		return nil, copyError
	}
	//OpenDatabase
	db, err := sql.Open("sqlite", tempHistoryDbff)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	last30Days := time.Now().AddDate(0, 0, -30).UnixMicro()

	// Execute a query
	rows, err := db.Query("SELECT url, title, visit_count, last_visit_date FROM moz_places WHERE last_visit_date >= ? ORDER BY last_visit_date DESC", last30Days)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var url string
		var title sql.NullString
		var visitCount int
		var lastVisitDate sql.NullInt64
		// Scan the row into variables
		if err := rows.Scan(&url, &title, &visitCount, &lastVisitDate); err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Check if the title is NULL
		var titleStr string
		if title.Valid {
			titleStr = title.String
		} else {
			titleStr = "<NULL>"
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

		fmt.Println(url, timeString)
		output = append(output, url, titleStr, string(visitCount), timeString)
	}
	os.Remove(tempHistoryDbff)
	return output, nil
}
