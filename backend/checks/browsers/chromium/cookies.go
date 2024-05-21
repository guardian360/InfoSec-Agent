package chromium

import (
	"database/sql"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"os"
	"path/filepath"
	"strings"
)

func CookiesChromium(browser string, getter browsers.DefaultDirGetter) checks.Check {
	browserPath, returnID := GetBrowserPathAndIDCookie(browser)
	userDataDir, err := getter.GetDefaultDir(browserPath)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}
	// TODO: Check if it is in the same position for Win10
	cookiesDir := userDataDir + "\\Network"

	// Copy the database, so problems don't arise when the file gets locked
	tempCookieDbchr := filepath.Join(os.TempDir(), "tempCookieDbchr.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing temporary "+browser+" cookie database: ", err)
		}
	}(tempCookieDbchr)

	// Copy the database to a temporary location
	copyError := browsers.CopyFile(cookiesDir+"\\Cookies", tempCookieDbchr, nil, nil)
	if copyError != nil {
		return checks.NewCheckErrorf(returnID, "Unable to make a copy of "+browser+" database: ", copyError)
	}

	db, err := sql.Open("sqlite", tempCookieDbchr)
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	// Query the name, origin and when the cookie was created from the database
	rows, err := db.Query("SELECT host_key, name FROM cookies")

	if rows.Err() != nil {
		return checks.NewCheckError(returnID, rows.Err())
	}
	if err != nil {
		return checks.NewCheckError(returnID, err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing rows: ", err)
		}
	}(rows)

	var possibleTrackingCookie = false
	var output []string
	// Iterate over each found cookie
	for rows.Next() {
		var name, host string
		// Scan the row into variables
		if err = rows.Scan(&host, &name); err != nil {
			return checks.NewCheckError(returnID, err)
		}
		// Check if the cookie is a (possible) tracking cookie
		// Check is based on the fact that Google Analytics tracking cookies usually contain the substrings "utm" or "ga"
		if strings.Contains(name, "_utm") || strings.Contains(name, "_ga") {
			possibleTrackingCookie = true
			// Append the cookie to the result list
			output = append(output, name, host)
		}
	}
	if possibleTrackingCookie {
		return checks.NewCheckResult(returnID, 1, output...)
	}
	return checks.NewCheckResult(returnID, 0)
}

// GetBrowserPathAndIDCookie is a function that takes a browser name as input,
// and returns the path to the browser's directory and the ID of the browser for the cookie check.
//
// Parameters:
//   - browser: A string representing the name of the browser. Currently, this function supports "Chrome" and "Edge".
//
// Returns:
//   - A string representing the path to the browser's directory.
//   - An integer representing the ID of the check.
//
// If the browser is unknown/unsupported, the function returns an empty string and 0.
func GetBrowserPathAndIDCookie(browser string) (string, int) {
	if browser == browsers.Chrome {
		return browsers.ChromePath, checks.CookiesChromiumID
	}
	if browser == browsers.Edge {
		return browsers.EdgePath, checks.CookiesEdgeID
	}
	return "", 0
}
