// Package firefox is responsible for running checks on Firefox.
//
// Exported function(s): CookieFirefox, ExtensionFirefox, HistoryFirefox, PasswordFirefox
package firefox

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	// Necessary to use the sqlite driver
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	_ "modernc.org/sqlite"
)

// CookiesFirefox inspects the cookies stored in the Firefox browser.
// It does so by using the browsers.QueryCookieDatabase function to query the cookies.sqlite database in the Firefox profile directory.
// Parameters:
//   - profileFinder: A browsers.FirefoxProfileFinder object that is used to locate the Firefox profile directory.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of cookies stored in the Firefox browser. Each cookie is represented as a string that includes the name and the host of the cookie. If an error occurs during the check, the result will contain a description of the error.
func CookiesFirefox(profileFinder browsers.FirefoxProfileFinder) checks.Check {
	// Determine the directory in which the Firefox profile is stored
	ffdirectory, err := profileFinder.FirefoxFolder()
	if err != nil {
		return checks.NewCheckErrorf(checks.CookiesFirefoxID, "No firefox directory found", err)
	}

	return browsers.QueryCookieDatabase(
		checks.CookiesFirefoxID, "Firefox", ffdirectory[0]+"\\cookies.sqlite",
		[]string{"name", "host"}, "moz_cookies")
}
