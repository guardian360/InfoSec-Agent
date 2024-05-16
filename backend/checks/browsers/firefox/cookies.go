// Package firefox is responsible for running checks on Firefox.
//
// Exported function(s): CookieFirefox, ExtensionFirefox, HistoryFirefox, PasswordFirefox
package firefox

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// CookiesFirefox inspects the cookies stored in the Firefox browser.
//
// Parameters: None
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of cookies stored in the Firefox browser. Each cookie is represented as a string that includes the name, host, and creation time of the cookie. If an error occurs during the check, the result will contain a description of the error.
//
// This function works by locating the Firefox profile directory and copying the cookies.sqlite database to a temporary location. It then opens this database and queries it for the name, host, and creation time of each cookie. The results are returned as a list of strings, each string representing a cookie. If any error occurs during this process, such as an error copying the file or querying the database, this error is returned as the result of the check.
func CookiesFirefox(profileFinder browsers.FirefoxProfileFinder) checks.Check {
	var output []string
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := profileFinder.FirefoxFolder()
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
	copyError := browsers.CopyFile(ffdirectory[0]+"\\cookies.sqlite", tempCookieDbff, nil, nil)
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
	rows, err := db.Query("SELECT name, host FROM moz_cookies")

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

	var possibleTrackingCookie = false
	// Iterate over each found cookie
	for rows.Next() {
		var name, host string
		// Scan the row into variables
		if err = rows.Scan(&name, &host); err != nil {
			return checks.NewCheckError(checks.CookiesFirefoxID, err)
		}
		// Check if the cookie is a (possible) tracking cookie
		// Check is based on the fact that Google Analytics tracking cookies usually contain the substrings "utm" or "ga"
		if strings.Contains(name, "utm") || strings.Contains(name, "ga") {
			possibleTrackingCookie = true
			// Append the cookie to the result list
			output = append(output, name, host)
		}
	}
	if possibleTrackingCookie {
		return checks.NewCheckResult(checks.CookiesFirefoxID, 1, output...)
	}
	return checks.NewCheckResult(checks.CookiesFirefoxID, 0)
}
