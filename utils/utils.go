// Package utils contains helper functions that can be used throughout the project
//
// Exported function(s): CopyFile, GetPhishingDomains, FirefoxFolder
package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// CopyFile copies a file from the source to the destination
//
// Parameters: src - the source file
//
// dst - the destination file
//
// Returns: an error if the file cannot be copied, nil if the file is copied successfully
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

	_, err = sourceFile.Copy(destinationFile, sourceFile)
	if err != nil {
		logger.Log.Println("Error copying file:", err)
		return err
	}
	return nil
}

// GetPhishingDomains gets the phishing domains from a remote GitHub list
//
// Parameters: _
//
// Returns: a list of phishing domains
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

// FirefoxFolder gets the path to the Firefox profile folder
//
// Parameters: _
//
// Returns: a list of paths to the Firefox profile folder, and an optional error which should be nil on success
func FirefoxFolder() ([]string, error) {
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

// CurrentUsername retrieves the current Windows username
//
// Parameters: _
//
// Returns: The current Windows username
func CurrentUsername() (string, error) {
	currentUser, err := user.Current()
	if currentUser.Username == "" || err != nil {
		return "", errors.New("failed to retrieve current username")
	}
	return strings.Split(currentUser.Username, "\\")[1], nil
}

// RemoveDuplicateStr removes duplicate strings from a slice
//
// Parameters: strSlice (string slice) - the slice to remove duplicates from
//
// Returns: A slice with the duplicates removed
func RemoveDuplicateStr(strSlice []string) []string {
	// Keep a map of found values, where true means the value has (already) been found
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			// If the value is found for the first time, append it to the list of results
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// CloseFile closes a file and handles associated errors
//
// Parameters: file *filemock.File - the file to close
//
// Returns: _
func CloseFile(file mocking.File) error {
	err := file.Close()
	if err != nil {
		logger.Log.ErrorWithErr("Error closing file: %s", err)
		return err
	}
	return nil
}
