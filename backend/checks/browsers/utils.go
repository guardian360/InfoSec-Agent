// Package browsers provides utility functions for handling browser-related operations.
// These functions are used in the context of performing security checks on a system.
package browsers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
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

// TODO: Update documentation
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
		logger.Log.ErrorWithErr("Error closing file", err)
		return err
	}
	return nil
}

// TODO: Update documentation
// FirefoxProfileFinder is an interface that wraps the FirefoxFolder method
type FirefoxProfileFinder interface {
	FirefoxFolder() ([]string, error)
}

// TODO: Update documentation
// RealProfileFinder is a struct that implements the FirefoxProfileFinder interface
type RealProfileFinder struct{}

// TODO: Update documentation
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
		logger.Log.ErrorWithErr("Error getting current user", err)
		return nil, err
	}
	// Specify the path to the firefox profile directory
	profilesDir := currentUser.HomeDir + "\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles"

	dir, err := os.Open(filepath.Clean(profilesDir))
	if err != nil {
		logger.Log.ErrorWithErr("Error getting profiles directory", err)
		return nil, err
	}
	defer func(dir *os.File) {
		err = dir.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing directory", err)
		}
	}(dir)

	// Read the contents of the directory
	files, err := dir.Readdir(0)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading contents", err)
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
		content, err = os.ReadFile(filepath.Clean(profilesDir + "\\" + folder + "\\addons.json"))
		if err != nil {
			continue
		}
		if content != nil {
			profileList = append(profileList, profilesDir+"\\"+folder)
		}
	}
	return profileList, nil
}

// TODO: Update documentation
// MockProfileFinder is a struct that implements the FirefoxProfileFinder interface for testing
type MockProfileFinder struct {
	MockFirefoxFolder func() ([]string, error)
}

// TODO: Update documentation
// FirefoxFolder is a mock function
func (m MockProfileFinder) FirefoxFolder() ([]string, error) {
	return m.MockFirefoxFolder()
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestCreator interface {
	NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
}

// TODO: Update documentation
// RealRequestCreator is a struct that implements the RequestCreator interface.
// It provides the real implementation of the NewRequestWithContext method.
type RealRequestCreator struct {
	Client Doer
}

// TODO: Update documentation
// NewRequestWithContext is a method of RealRequestCreator that creates a new HTTP request with the provided context, method, URL, and body.
// It uses the http.NewRequestWithContext function from the net/http package to create the request.
//
// Parameters:
//   - ctx context.Context: The context to use for the request. This context will be used for timeout and cancellation signals, and for passing request-scoped values.
//   - method string: The HTTP method to use for the request (e.g., "GET", "POST").
//   - url string: The URL to use for the request.
//   - body io.Reader: The body of the request. This can be nil for methods that do not require a body (like GET).
//
// Returns:
//   - *http.Request: The created HTTP request.
//   - error: An error object that wraps any error that occurs during the creation of the request. If the request is created successfully, it returns nil.
func (r RealRequestCreator) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// TODO: Update documentation
// PhishingDomainGetter is an interface that wraps the GetPhishingDomains method.
type PhishingDomainGetter interface {
	// GetPhishingDomains retrieves a list of active phishing domains from a remote database.
	GetPhishingDomains(creator RequestCreator) ([]string, error)
}

// TODO: Update documentation
// RealPhishingDomainGetter is a struct that implements the PhishingDomainGetter interface.
// It provides the real implementation of the GetPhishingDomains method.
type RealPhishingDomainGetter struct {
	Client Doer
}

// TODO: Update documentation
// NewRealPhishingDomainGetter is a constructor function for the RealPhishingDomainGetter struct.
// It takes a Doer interface as an argument, which is used to make HTTP requests.
//
// Parameters:
//   - client Doer: An object that implements the Doer interface. This object is used to make HTTP requests.
//
// Returns:
//   - RealPhishingDomainGetter: A new instance of RealPhishingDomainGetter with the provided client.
func NewRealPhishingDomainGetter(client Doer) RealPhishingDomainGetter {
	return RealPhishingDomainGetter{Client: client}
}

// TODO: Update documentation
// GetPhishingDomains retrieves a list of active phishing domains from a remote database.
//
// This function sends a GET request to the URL of the phishing database hosted on GitHub. It reads the response body,
// which contains a list of active phishing domains, each on a new line. The function then splits this response into a slice
// of strings, where each string represents a single phishing domain.
// Parameters:
//   - creator RequestCreator: An object that implements the RequestCreator interface. It is used to create an HTTP request to fetch the phishing domains.
//
// Returns:
//   - []string: A slice containing the phishing domains. If an error occurs during the retrieval or parsing of the domains, an empty slice is returned.
//   - error: An error object that wraps any error that occurs during the retrieval of the phishing domains. If the domains are retrieved successfully, it returns nil.
func (r RealPhishingDomainGetter) GetPhishingDomains(creator RequestCreator) ([]string, error) {
	if r.Client == nil {
		r.Client = http.DefaultClient
	}
	// Get the phishing domains from up-to-date GitHub list
	url := "https://raw.githubusercontent.com/mitchellkrogza/Phishing.Database/master/phishing-links-ACTIVE-today.txt"
	req, err := creator.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		logger.Log.ErrorWithErr("Error creating HTTP request", err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")

	resp, err := r.Client.Do(req)
	if err != nil {
		logger.Log.ErrorWithErr("Error sending HTTP request", err)
		return nil, err
	}
	// Ensure the response body is closed properly
	defer func() {
		if resp != nil && resp.Body != nil {
			resErr := resp.Body.Close()
			if resErr != nil {
				logger.Log.ErrorWithErr("Error closing response body", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Error("HTTP request failed with status code: " + strconv.Itoa(resp.StatusCode))
		return nil, errors.New("HTTP request failed")
	}

	// Parse the response of potential scam domains and split it into a list of domains
	scamDomainsResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading response body", err)
		return nil, err
	}

	if len(scamDomainsResponse) == 0 {
		logger.Log.Error("Response body is empty, no phishing domains list found online")
		return nil, errors.New("no phishing domains list found online")
	}

	// Split the response into a list of domains and remove the http:// and https:// prefixes
	domains := strings.Split(string(scamDomainsResponse), "\n")
	var result []string
	for _, domain := range domains {
		trimmedDomain := strings.TrimPrefix(domain, "http://")
		trimmedDomain = strings.TrimPrefix(trimmedDomain, "https://")
		if trimmedDomain != "" {
			result = append(result, trimmedDomain)
		}
	}
	return result, nil
}

// TODO: Update documentation
// CopyFileGetter is an interface that defines a method for copying a file from a source path to a destination path.
// This interface is used as a contract that must be fulfilled by any type that wishes to provide functionality
// for copying files.
type CopyFileGetter interface {
	// CopyFile is a method that takes a source file path, a destination file path, a mock source file, and a mock destination file as input.
	// It copies the source file to the destination path. If the mock source file and mock destination file are provided, it uses them instead of the actual files.
	// This method is used when there is a need to create a copy of a file.
	//
	// Parameters:
	//   - src string: The path to the source file that needs to be copied.
	//   - dst string: The path to the destination where the source file should be copied to.
	//   - mockSource mocking.File: A mock file object that represents the source file. If this parameter is not nil, the function uses the mock file for the source.
	//   - mockDest mocking.File: A mock file object that represents the destination file. If this parameter is not nil, the function uses the mock file for the destination.
	//
	// Returns:
	//   - error: An error object that wraps any error that occurs during the file copying process. If the file is copied successfully, it returns nil.
	CopyFile(src, dst string, mockSource mocking.File, mockDest mocking.File) error
}

// TODO: Update documentation
// RealCopyFile is a struct that implements the CopyFileGetter interface.
// It provides the real implementation of the CopyFile method.
type RealCopyFileGetter struct{}

// TODO: Update documentation
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
func (r RealCopyFileGetter) CopyFile(src, dst string, mockSource mocking.File, mockDestination mocking.File) error {
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
			logger.Log.ErrorWithErr("Error closing source file", err)
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
			logger.Log.ErrorWithErr("Error closing destination file", err)
		}
	}(destinationFile)

	_, err = sourceFile.Copy(sourceFile, destinationFile)
	if err != nil {
		logger.Log.ErrorWithErr("Error copying file", err)
		return err
	}
	return nil
}

// TODO: Update documentation
// DefaultDirGetter is an interface that wraps the GetPreferencesDir method.
// It provides a way to get the default directory of a specific browser.
type DefaultDirGetter interface {
	// GetDefaultDir takes a browser name as input and returns the path to the preferences directory of the browser.
	// It returns an error if there is any issue in getting the default directory.
	GetDefaultDir(browser string) (string, error)
}

// TODO: Update documentation
// RealDefaultDirGetter is a struct that implements the DefaultDirGetter interface.
// It provides the real implementation of the GetDefaultDir method.
type RealDefaultDirGetter struct{}

// TODO: Update documentation
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

// TODO: Update documentation
// QueryCookieDatabaseGetter is an interface that defines a method for querying a cookie database.
// This interface is used as a contract that must be fulfilled by any type that wishes to provide functionality
// for querying a cookie database.
//
// The QueryCookieDatabase method takes several parameters:
// - checkID: The ID of the check that is being performed.
// - browser: The name of the browser for which the check is being performed.
// - databasePath: The path to the cookie database file.
// - queryParams: A list of parameters to use in the SQL query for the database.
// - tableName: The name of the table in the database to query.
// - getter: An object that implements the CopyFileGetter interface. It is used to copy the database file to a temporary location.
//
// The method returns a Check object representing the result of the check. If tracking cookies are found, the result contains a list of cookies along with their host stored in the database.
type QueryCookieDatabaseGetter interface {
	QueryCookieDatabase(checkID int, browser string, databasePath string, queryParams []string, tableName string, getter CopyFileGetter) checks.Check
}

// TODO: Update documentation
// RealQueryCookieDatabaseGetter is an empty struct that is used as a receiver for the QueryCookieDatabase method.
// This struct is part of the implementation of the QueryCookieDatabaseGetter interface.
type RealQueryCookieDatabaseGetter struct{}

// TODO: Update documentation
// QueryCookieDatabase is a utility function that queries a cookie database for specific parameters.
// This function is used by the browser-specific (Firefox, Chrome, and Edge) cookie checks to query the cookie database and check for tracking cookies.
//
// Parameters:
//   - checkID int: The ID of the check that is being performed.
//   - browser string: The name of the browser for which the check is being performed.
//   - databasePath string: The path to the cookie database file.
//   - queryParams []string: A list of parameters to use in the SQL query for the database.
//   - tableName string: The name of the table in the database to query.
//   - getter CopyFileGetter: An object that implements the CopyFileGetter interface. It is used to copy the database file to a temporary location.
//
// Returns:
//   - checks.Check: A Check object representing the result of the check. If tracking cookies are found, the result contains a list of cookies along with their host stored in the database.
func (r RealQueryCookieDatabaseGetter) QueryCookieDatabase(checkID int, browser string, databasePath string, queryParams []string, tableName string, getter CopyFileGetter) checks.Check {
	// Copy the database, so problems don't arise when the file gets locked
	tempCookieDB := filepath.Join(os.TempDir(), "tempCookieDb"+browser+".sqlite")

	// Clean up the temporary file when the function returns
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			logger.Log.ErrorWithErr("Error removing temporary "+browser+" cookie database", err)
		}
	}(tempCookieDB)

	// Copy the database to a temporary location
	copyError := getter.CopyFile(databasePath, tempCookieDB, nil, nil)
	if copyError != nil {
		return checks.NewCheckErrorf(checkID, "Unable to make a copy of "+browser+" database", copyError)
	}

	db, err := sql.Open("sqlite", tempCookieDB)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing "+browser+" database", err)
		}
	}(db)

	sqlSelectors := strings.Join(queryParams, ", ")
	// Query the name, origin and when the cookie was created from the database
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM %s", sqlSelectors, tableName))
	if err != nil {
		return checks.NewCheckError(checkID, err)
	}
	if rows.Err() != nil {
		return checks.NewCheckError(checkID, rows.Err())
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing "+browser+" rows", err)
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
