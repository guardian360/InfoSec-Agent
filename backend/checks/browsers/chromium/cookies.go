package chromium

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
)

// CookiesChromium inspects the cookies stored in Chromium based browsers.
// It does so by using the browsers.QueryCookieDatabase function to query the Cookies database in the specific browser User Data directory.
// Parameters:
//   - browser: A string representing the name of the browser. Currently, this function supports "Chrome" and "Edge".
//   - getter: A browsers.DefaultDirGetter object that is used to locate the default directory of the browser.
//   - copyGetter: A browsers.CopyFileGetter object that is used to copy the database file to a temporary location.
//   - queryGetter: A browsers.QueryCookieDatabaseGetter object that is used to query the database file.
//
// Returns:
//   - A checks.Check object representing the result of the check. The result contains a list of cookies stored in the Chromium based browser. Each cookie is represented as a string that includes the name and the host of the cookie. If an error occurs during the check, the result will contain a description of the error.
func CookiesChromium(browser string, getter browsers.DefaultDirGetter, copyGetter browsers.CopyFileGetter, queryGetter browsers.QueryCookieDatabaseGetter) checks.Check {
	browserPath, returnID := GetBrowserPathAndIDCookie(browser)
	userDataDir, err := getter.GetDefaultDir(browserPath)
	if err != nil {
		return checks.NewCheckErrorf(returnID, "Error: ", err)
	}
	cookiesDir := userDataDir + "\\Network\\Cookies"

	return queryGetter.QueryCookieDatabase(returnID, browser, cookiesDir, []string{"name", "host_key"}, "cookies", copyGetter)
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
