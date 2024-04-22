package firefox

import (
	"database/sql"
	browser_utils "github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers"
	"os"
	"path/filepath"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// CookieFirefox inspects the cookies stored in the Firefox browser.
//
// Parameters: None
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of cookies stored in the Firefox browser. Each cookie is represented as a string that includes the name, host, and creation time of the cookie. If an error occurs during the check, the result will contain a description of the error.
//
// This function works by locating the Firefox profile directory and copying the cookies.sqlite database to a temporary location. It then opens this database and queries it for the name, host, and creation time of each cookie. The results are returned as a list of strings, each string representing a cookie. If any error occurs during this process, such as an error copying the file or querying the database, this error is returned as the result of the check.
func CookieFirefox() checks.Check {
	var output []string
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := browser_utils.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.CookiesFirefoxID, "No firefox directory found", err)
	}

	// Copy the database, so problems don't arise when the file gets locked
	tempCookieDbff := filepath.Join(os.TempDir(), "tempCookieDbff.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing file: ", err)
		}
	}(tempCookieDbff)

	// Copy the database to a temporary location
	copyError := browser_utils.CopyFile(ffdirectory[0]+"\\cookies.sqlite", tempCookieDbff, nil, nil)
	if copyError != nil {
		return checks.NewCheckErrorf(checks.CookiesFirefoxID, "Unable to make a copy of the file", copyError)
	}

	db, err := sql.Open("sqlite", tempCookieDbff)
	if err != nil {
		return checks.NewCheckError(checks.CookiesFirefoxID, err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	// Query the name, origin and when the cookie was created from the database
	rows, err := db.Query("SELECT name, host, creationTime FROM moz_cookies")
	// TODO: check if this is error handling is correct
	if rows.Err() != nil {
		return checks.NewCheckError(checks.CookiesFirefoxID, rows.Err())
	}
	if err != nil {
		return checks.NewCheckError(checks.CookiesFirefoxID, err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing rows: ", err)
		}
	}(rows)

	// Iterate over each found cookie
	for rows.Next() {
		var name, host string
		var creationTime int64
		// Scan the row into variables
		if err = rows.Scan(&name, &host, &creationTime); err != nil {
			return checks.NewCheckError(checks.CookiesFirefoxID, err)
		}
		// Append the cookie to the result list
		timeofCreation := time.UnixMicro(creationTime)
		timeString := timeofCreation.String()
		output = append(output, name, host, timeString)
	}
	return checks.NewCheckResult(checks.CookiesFirefoxID, 0, output...)
}
