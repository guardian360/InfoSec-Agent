// Package browsers provides utility functions for handling browser-related operations.
// These functions are used in the context of performing security checks on a system.
//
// Exported function(s): CloseFile, FirefoxFolder, GetPhishingDomains, CopyFile
package browsers

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	// Necessary to use the sqlite driver
	_ "modernc.org/sqlite"
)

// constants to store information of the browsers

const Edge = "Edge"
const Chrome = "Chrome"
const EdgePath = "Microsoft/Edge"
const ChromePath = "Google/Chrome"

var UserHomeDirFunc = os.UserHomeDir

// CloseFile is a utility function that closes a given file and logs any errors that occur during the process.
//
// Parameters:
//   - file *filemock.File: The file that needs to be closed. It is an instance of a File from the filemock package.
//
// Returns:
//   - error: An error object that wraps any error that occurs during file closure. If the file is closed successfully, it returns nil.
func CloseFile(file mocking.File) error {
	err := file.Close()
	if err != nil {
		logger.Log.ErrorWithErr("Error closing file: %s", err)
		return err
	}
	return nil
}

// FirefoxProfileFinder is an interface that wraps the FirefoxFolder method
type FirefoxProfileFinder interface {
	FirefoxFolder() ([]string, error)
}

// RealProfileFinder is a struct that implements the FirefoxProfileFinder interface
type RealProfileFinder struct{}

// FirefoxFolder retrieves the paths to all Firefox profile folders for the currently logged-in user.
//
// This function uses the os/user package to access the current user's information and constructs the path to the Firefox profile directory.
// It then reads the directory and filters out all non-directory files. For each directory, it checks if a 'logins.json' file exists.
// If such a file exists, the directory is considered a Firefox profile folder and its path is added to the returned list.
//
// Returns:
//   - []string: A slice containing the paths to all Firefox profile folders. If no profile folders are found or an error occurs, an empty slice is returned.
//   - error: An error object that wraps any error that occurs during the retrieval of the Firefox profile folders. If the folders are retrieved successfully, it returns nil.
func (r RealProfileFinder) FirefoxFolder() ([]string, error) {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		logger.Log.ErrorWithErr("Error getting current user:", err)
		return nil, err
	}
	// Specify the path to the firefox profile directory
	profilesDir := currentUser.HomeDir + "\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles"

	dir, err := os.Open(filepath.Clean(profilesDir))
	if err != nil {
		logger.Log.ErrorWithErr("Error getting profiles directory:", err)
		return nil, err
	}
	defer func(dir *os.File) {
		err = dir.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing directory: %v", err)
		}
	}(dir)

	// Read the contents of the directory
	files, err := dir.Readdir(0)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading contents:", err)
		return nil, err
	}

	// Iterate through the files and get only directories
	var folders []string
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}
	var profileList []string
	var content []byte
	// Loop through all the folders to check if they have a logins.json file.
	for _, folder := range folders {
		content, err = os.ReadFile(filepath.Clean(profilesDir + "\\" + folder + "\\logins.json"))
		if err != nil {
			continue
		}
		if content != nil {
			profileList = append(profileList, profilesDir+"\\"+folder)
		}
	}
	return profileList, nil
}

// MockProfileFinder is a struct that implements the FirefoxProfileFinder interface for testing
type MockProfileFinder struct {
	MockFirefoxFolder func() ([]string, error)
}

// FirefoxFolder is a mock function
func (m MockProfileFinder) FirefoxFolder() ([]string, error) {
	return m.MockFirefoxFolder()
}

// GetPhishingDomains retrieves a list of active phishing domains from a remote database.
//
// This function sends a GET request to the URL of the phishing database hosted on GitHub. It reads the response body,
// which contains a list of active phishing domains, each on a new line. The function then splits this response into a slice
// of strings, where each string represents a single phishing domain.
//
// Returns:
//   - []string: A slice containing the phishing domains. If an error occurs during the retrieval or parsing of the domains, an empty slice is returned.
func GetPhishingDomains() []string {
	// Get the phishing domains from up-to-date GitHub list
	client := &http.Client{}
	url := "https://raw.githubusercontent.com/mitchellkrogza/Phishing.Database/master/phishing-domains-ACTIVE.txt"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	if err != nil {
		logger.Log.FatalWithErr("Error creating HTTP request:", err)
	}
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Log.Printf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Parse the response of potential scam domains and split it into a list of domains
	scamDomainsResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading response body:", err)
	}
	return strings.Split(string(scamDomainsResponse), "\n")
}

// CopyFile is a utility function that copies a file from a source path to a destination path.
//
// Parameters:
//   - src string: The path to the source file that needs to be copied.
//   - dst string: The path to the destination where the source file should be copied to.
//   - mockSource mocking.File: A mock file object that represents the source file. If this parameter is not nil, the function uses the mock file for the source.
//   - mockDestination mocking.File: A mock file object that represents the destination file. If this parameter is not nil, the function uses the mock file for the destination.
//
// Returns:
//   - error: An error object that wraps any error that occurs during the file copying process. If the file is copied successfully, it returns nil.
func CopyFile(src, dst string, mockSource mocking.File, mockDestination mocking.File) error {
	var sourceFile mocking.File
	var err error
	if mockSource != nil {
		sourceFile, err = mockSource, nil
	} else {
		var tmp *os.File
		tmp, err = os.Open(filepath.Clean(src))
		sourceFile = mocking.Wrap(tmp)
	}
	if err != nil {
		return err
	}
	defer func(sourceFile mocking.File) {
		err = sourceFile.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing source file:", err)
		}
	}(sourceFile)
	var destinationFile mocking.File
	if mockDestination != nil {
		destinationFile, err = mockDestination, nil
	} else {
		var tmp *os.File
		tmp, err = os.Create(filepath.Clean(dst))
		destinationFile = mocking.Wrap(tmp)
	}
	if err != nil {
		return err
	}
	defer func(destinationFile mocking.File) {
		err = destinationFile.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing destination file:", err)
		}
	}(destinationFile)

	_, err = sourceFile.Copy(sourceFile, destinationFile)
	if err != nil {
		logger.Log.Println("Error copying file:", err)
		return err
	}
	return nil
}

// DefaultDirGetter is an interface that wraps the GetPreferencesDir method.
// It provides a way to get the default directory of a specific browser.
type DefaultDirGetter interface {
	// GetDefaultDir takes a browser name as input and returns the path to the preferences directory of the browser.
	// It returns an error if there is any issue in getting the default directory.
	GetDefaultDir(browser string) (string, error)
}

// RealDefaultDirGetter is a struct that implements the DefaultDirGetter interface.
// It provides the real implementation of the GetDefaultDir method.
type RealDefaultDirGetter struct{}

// GetDefaultDir is a method of RealDefaultDirGetter that gets the default directory of a specific browser.
// It takes a browser name as input and returns the path to the default directory of the browser.
// It returns an error if there is any issue in getting the default directory.
func (r RealDefaultDirGetter) GetDefaultDir(browserPath string) (string, error) {
	userDir, err := UserHomeDirFunc()
	if err != nil {
		return "", err
	}
	return filepath.Join(userDir, "AppData", "Local", browserPath, "User Data", "Default"), nil
}

// QueryCookieDatabase is a utility function that queries a cookie database for specific parameters.
// This function is used by the browser-specific (Firefox, Chrome, and Edge) cookie checks to query the cookie database and check for tracking cookies.
//
// Parameters:
//   - checkID int: The ID of the check that is being performed.
//   - browser string: The name of the browser for which the check is being performed.
//   - databasePath string: The path to the cookie database file.
//   - queryParams []string: A list of parameters to use in the SQL query for the database.
//   - tableName string: The name of the table in the database to query.
//
// Returns:
//   - checks.Check: A Check object representing the result of the check. If tracking cookies are found, the result contains a list of cookies along with their host stored in the database.
func QueryCookieDatabase(checkID int, browser string, databasePath string, queryParams []string, tableName string) checks.Check {
	// Copy the database, so problems don't arise when the file gets locked
	tempCookieDB := filepath.Join(os.TempDir(), "tempCookieDbchr.sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing temporary "+browser+" cookie database: ", err)
		}
	}(tempCookieDB)

	// Copy the database to a temporary location
	copyError := CopyFile(databasePath, tempCookieDB, nil, nil)
	if copyError != nil {
		return checks.NewCheckErrorf(checkID, "Unable to make a copy of "+browser+" database: ", copyError)
	}

	db, err := sql.Open("sqlite", tempCookieDB)
	if err != nil {
		return checks.NewCheckError(checkID, err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing "+browser+" database: ", err)
		}
	}(db)

	sqlSelectors := strings.Join(queryParams, ", ")
	// Query the name, origin and when the cookie was created from the database
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM %s", sqlSelectors, tableName))

	if rows.Err() != nil {
		return checks.NewCheckError(checkID, rows.Err())
	}
	if err != nil {
		return checks.NewCheckError(checkID, err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing "+browser+" rows: ", err)
		}
	}(rows)

	var possibleTrackingCookie = false
	var output []string
	// Iterate over each found cookie
	for rows.Next() {
		var name, host string
		// Scan the row into variables
		if err = rows.Scan(&name, &host); err != nil {
			return checks.NewCheckError(checkID, err)
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
		return checks.NewCheckResult(checkID, 1, output...)
	}
	return checks.NewCheckResult(checkID, 0)
}
