package checks

import (
	"InfoSec-Agent/utils"
	"database/sql"
	"os"
	"path/filepath"
	"time"
	"InfoSec-Agent/checks"

	_ "modernc.org/sqlite"
)

func CookieFirefox() checks.Check {
	var output []string
	ffdirectory, _ := utils.FirefoxFolder()

	//Copy the database so we don't have problems with locked files
	tempCookieDbff := filepath.Join(os.TempDir(), "tempCookieDbff.sqlite")

	copyError := utils.CopyFile(ffdirectory[0]+"\\cookies.sqlite", tempCookieDbff)
	if copyError != nil {
		return checks.newCheckErrorf("CookieFirefox", "Unable to make a copy of the file", copyError)
	}

	//OpenDatabase
	db, err := sql.Open("sqlite", tempCookieDbff)
	if err != nil {
		return newCheckError("CookieFirefox", err)
	}
	defer db.Close()
	// Execute a query
	rows, err := db.Query("SELECT name, host, creationTime FROM moz_cookies")
	if err != nil {
		return newCheckError("CookieFirefox", err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var name, host string
		var creationTime int64
		// Scan the row into variables
		if err := rows.Scan(&name, &host, &creationTime); err != nil {
			return  newCheckError("CookieFirefox", err)
		}
		// Print the row
		timeofCreation := time.UnixMicro(creationTime)
		timeString := timeofCreation.String()
		output = append(output, name, host, timeString)
	}
	os.Remove(tempCookieDbff)
	return newCheckResult("CookieFirefox", output...)
}
