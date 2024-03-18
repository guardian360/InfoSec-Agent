package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func GetPhisingDomains() []string {
	// Get the phishing domains from up-to-date github list
	client := &http.Client{}
	url := fmt.Sprintf("https://raw.githubusercontent.com/mitchellkrogza/Phishing.Database/master/phishing-domains-ACTIVE.txt")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	//Parse the response of scamdomains and split it into a list of domains
	scamDomainsResponse, err := io.ReadAll(resp.Body)
	return strings.Split(string(scamDomainsResponse), "\n")
}

func FirefoxFolder() ([]string, error) {
	// Get the current user
	currentUser, err := user.Current() //gets the path to the directory of the current user
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	profilesDir := currentUser.HomeDir + "\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles" //Path to the firefox profile directory
	// Open the directory
	dir, err := os.Open(profilesDir)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer dir.Close()

	// Read the contents of the directory
	files, err := dir.Readdir(0)
	if err != nil {
		fmt.Println("Error:", err)
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
	//loop through all the folders to check if they have a logins.json file.
	for _, folder := range folders {
		content, err := os.ReadFile(profilesDir + "\\" + folder + "\\logins.json") // the file is inside the local directory
		if err != nil {
			continue
		}
		if content != nil {
			profileList = append(profileList, profilesDir+"\\"+folder)
		}
	}
	return profileList, nil
}
