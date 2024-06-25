package programs

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

// InstalledSoftware is a function that checks for outdated software on the system.
// It uses a CommandExecutor to execute system commands for retrieving the list of installed software.
// The function collects all software lists, filters and deduplicates the software list, and formats the result array.
// It returns a Check object that represents the result of the check for outdated software.
//
// Parameters:
// executor: A CommandExecutor that is used to execute system commands for retrieving the list of installed software.
// registryKey: A RegistryKey that is used to access the system's registry.
//
// Returns:
// checks.Check: A Check object that represents the result of the check for outdated software.
func InstalledSoftware(executor mocking.CommandExecutor, registryKey mocking.RegistryKey) checks.Check {
	// Collect all software lists
	softwareList, err := collectAllSoftwareLists(executor, registryKey)
	if softwareList == nil {
		return err
	}

	// Filter and deduplicate the software list
	uniqueSoftware := filterAndDeduplicateSoftware(softwareList)
	// Format the result array
	resultArray := formatResultArray(uniqueSoftware)

	// Return the check result
	return checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, resultArray...)
}

// TODO: Update documentation
// collectAllSoftwareLists is a function that collects all software lists from different sources.
// It uses a CommandExecutor to execute system commands for retrieving the list of installed software.
// The function retrieves installed programs from winget, installed 32 bit programs, and installed 64 bit programs.
// It then appends all software lists and formats the result array.
//
// Parameters:
// executor: A CommandExecutor that is used to execute system commands for retrieving the list of installed software.
// registryKey: A RegistryKey that is used to access the system's registry.
//
// Returns:
//   - []software: A slice of software objects that represents the list of all installed software.
//   - checks.Check: A Check object that represents the result of the check for installed software. If an error occurs during the retrieval of the software list, it returns a CheckError object.
func collectAllSoftwareLists(executor mocking.CommandExecutor, registryKey mocking.RegistryKey) ([]software, checks.Check) {
	var (
		softwareList       []software
		softwareListWinget []software
		softwareList32     []software
		softwareList64     []software
		err                error
	)

	// Retrieve installed programs from winget
	if softwareListWinget, err = retrieveWingetInstalledPrograms(softwareList, executor); err != nil {
		logger.Log.Debug("Error retrieving winget installed programs")
	}
	// Retrieve installed 32 bit programs
	if softwareList32, err = retrieveInstalled32BitPrograms(softwareList, registryKey); err != nil {
		logger.Log.Debug("Error retrieving 32 bit installed programs")
	}
	// Retrieve installed 64 bit programs
	if softwareList64, err = retrieveInstalled64BitPrograms(softwareList, registryKey); err != nil {
		logger.Log.Debug("Error retrieving 64 bit installed programs")
	}

	// Append all software lists
	softwareList = append(softwareList, softwareListWinget...)
	softwareList = append(softwareList, softwareList32...)
	softwareList = append(softwareList, softwareList64...)

	// Format the result array
	resultArray := make([]string, 0)
	for _, v := range softwareList {
		resultArray = append(resultArray, fmt.Sprintf("%s | %s", v.name, v.version))
	}

	return softwareList, checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, resultArray...)
}

// TODO: Update documentation
// filterAndDeduplicateSoftware is a function that filters and deduplicates a list of software.
// It uses the software name as the key in a map to ensure uniqueness. If duplicate software is found, the existing software is kept in the map.
//
// Parameters:
// softwareList: A slice of software objects that represents the list of all installed software.
//
// Returns:
// map[string]software: A map where the key is the normalized software name and the value is the software object. This ensures that each software in the map is unique.
func filterAndDeduplicateSoftware(softwareList []software) map[string]software {
	uniqueSoftware := make(map[string]software)
	for _, sw := range softwareList {
		if sw.name == "" || sw.version == "" || strings.Contains(strings.ToLower(sw.name), "microsoft defender") {
			continue
		}
		normalized := normalize(sw.name)
		if existing, exists := uniqueSoftware[normalized]; exists {
			if compareVersions(sw.version, existing.version) > 0 {
				uniqueSoftware[normalized] = sw
			}
		} else {
			uniqueSoftware[normalized] = sw
		}
	}
	return uniqueSoftware
}

// TODO: Update documentation
// formatResultArray is a function that formats the result array of unique software.
// It iterates over the uniqueSoftware map and appends each software's name and version to the result array.
//
// Parameters:
// uniqueSoftware: A map where the key is the normalized software name and the value is the software object. This ensures that each software in the map is unique.
//
// Returns:
// []string: A slice of strings where each string represents a unique software in the format "name | version".
func formatResultArray(uniqueSoftware map[string]software) []string {
	resultArray := make([]string, 0, len(uniqueSoftware))
	for _, v := range uniqueSoftware {
		resultArray = append(resultArray, fmt.Sprintf("%s | %s", v.name, v.version))
	}
	return resultArray
}

// TODO: Update documentation
// normalize function to clean and standardize software names
func normalize(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Remove all non-alphanumeric characters (except spaces)
	var cleaned []rune
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			cleaned = append(cleaned, r)
		}
	}
	return strings.TrimSpace(string(cleaned))
}

// TODO: Update documentation
// compareVersions is a function that compares two version strings.
// It splits the version strings by the dot character and compares each corresponding part as an integer.
// If a part in v1 is greater than the corresponding part in v2, it returns 1.
// If a part in v1 is less than the corresponding part in v2, it returns -1.
// If all parts are equal, it returns 0.
//
// Parameters:
// v1: A string that represents the first version to be compared.
// v2: A string that represents the second version to be compared.
//
// Returns:
// int: An integer that indicates the result of the comparison. If v1 is greater than v2, it returns 1. If v1 is less than v2, it returns -1. If v1 is equal to v2, it returns 0.
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := range maxLen {
		var num1, num2 int
		if i < len(parts1) {
			_, err := fmt.Sscanf(parts1[i], "%d", &num1)
			if err != nil {
				logger.Log.ErrorWithErr("Error parsing version number", err)
			}
		}
		if i < len(parts2) {
			_, err := fmt.Sscanf(parts2[i], "%d", &num2)
			if err != nil {
				logger.Log.ErrorWithErr("Error parsing version number", err)
			}
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}
	return 0
}

// TODO: Update documentation
// retrieveWingetInstalledPrograms is a function that retrieves all installed packages found with the winget package manager.
// It uses a CommandExecutor to execute the winget list command and processes the output to extract the software details.
// The function appends each software to the softwareList and returns the updated softwareList.
//
// Parameters:
// softwareList: A slice of software objects that represents the list of all installed software.
// executor: A CommandExecutor that is used to execute the winget list command.
//
// Returns:
//   - []software: A slice of software objects that represents the updated list of all installed software.
//   - error: An error object that represents any error that occurred during the execution of the winget list command or the processing of the output. If no error occurred, it returns nil.
func retrieveWingetInstalledPrograms(softwareList []software, executor mocking.CommandExecutor) ([]software, error) {
	// Execute the winget list command
	out, err := executor.Execute("powershell", "winget list| Out-String -Stream | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Text.Encoding]::Default.GetBytes($_)) }")
	if err != nil {
		return softwareList, err
	}
	// Process the output
	lines := strings.Split(string(out), "\r\n")
	indexN := -2
	for i, line := range lines {
		if strings.Contains(line, "N") {
			indexN = i
			break
		}
	}
	if indexN < 0 {
		return softwareList, errors.New("error parsing winget output")
	}
	lines[0] = lines[indexN][strings.Index(lines[indexN], "Name")+1:] // Remove the first part of the header
	idIndex := strings.Index(lines[0], "Id")
	versionIndex := strings.Index(lines[0], "Version")
	availableIndex := strings.Index(lines[0], "Available")
	sourcesIndex := strings.Index(lines[0], "Source")
	for _, line := range lines[indexN+2:] { // Skip the header lines
		if len(line) != 0 { // Don't handle the last empty line, and maybe other empty lines
			// Extract the software details
			name := substr(line, 0, idIndex)
			name = strings.TrimSpace(name)
			id := substr(line, idIndex, versionIndex-idIndex)
			id = strings.TrimSpace(id)
			version := substr(line, versionIndex, availableIndex-versionIndex)
			version = strings.TrimSpace(version)
			available := substr(line, availableIndex, sourcesIndex-availableIndex)
			available = strings.TrimSpace(available)
			source := substr(line, sourcesIndex, len(line)-sourcesIndex)
			source = strings.TrimSpace(source)
			// Append the software to the list
			softwareList = append(softwareList, software{
				name:         name,
				identifier:   id,
				version:      version,
				newVersion:   available,
				vendor:       "",
				lastUpdated:  "",
				sourceWinget: source,
				whereFrom:    "winget",
			})
		}
	}
	return softwareList, nil
}

// TODO: Update documentation
// retrieveInstalled32BitPrograms is a function that retrieves all installed 32-bit programs found using a registry query.
// It uses a CommandExecutor to execute the registry query command and processes the output to extract the software details.
// The function appends each software to the softwareList and returns the updated softwareList.
//
// Parameters:
// softwareList: A slice of software objects that represents the list of all installed software.
// registryKey: A RegistryKey that is used to access the system's registry.
//
// Returns:
//   - []software: A slice of software objects that represents the updated list of all installed software.
//   - error: An error object that represents any error that occurred during the execution of the registry query command or the processing of the output. If no error occurred, it returns nil.
func retrieveInstalled32BitPrograms(softwareList []software, registryKey mocking.RegistryKey) ([]software, error) {
	return retrieveInstalledPrograms(softwareList, registryKey, "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "32-bit")
}

// TODO: Update documentation
// retrieveInstalled64BitPrograms is a function that retrieves all installed 64-bit programs found using a registry query.
// It uses a CommandExecutor to execute the registry query command and processes the output to extract the software details.
// The function appends each software to the softwareList and returns the updated softwareList.
//
// Parameters:
// softwareList: A slice of software objects that represents the list of all installed software.
// registryKey: A RegistryKey that is used to access the system's registry.
//
// Returns:
//   - []software: A slice of software objects that represents the updated list of all installed software.
//   - error: An error object that represents any error that occurred during the execution of the registry query command or the processing of the output. If no error occurred, it returns nil.
func retrieveInstalled64BitPrograms(softwareList []software, registryKey mocking.RegistryKey) ([]software, error) {
	return retrieveInstalledPrograms(softwareList, registryKey, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "64-bit")
}

// retrieveInstalledPrograms is a function that retrieves all installed programs found using a registry query.
// It uses a RegistryKey to access the system's registry and processes the output to extract the software details.
// The function appends each software to the softwareList and returns the updated softwareList.
//
// Parameters:
// softwareList: A slice of software objects that represents the list of all installed software.
// registryKey: A RegistryKey that is used to access the system's registry.
// path: A string that represents the registry key to be queried for installed software.
// bitOrigin: A string that represents the origin of the software (e.g., "32-bit", "64-bit").
//
// Returns:
//   - []software: A slice of software objects that represents the updated list of all installed software.
//   - error: An error object that represents any error that occurred during the execution of the registry query command or the processing of the output. If no error occurred, it returns nil.
func retrieveInstalledPrograms(softwareList []software, registryKey mocking.RegistryKey, path string, bitOrigin string) ([]software, error) {
	// Open the registry key
	key, err := mocking.OpenRegistryKey(registryKey, path)
	if err != nil {
		return softwareList, err
	}
	// Close the key after we have received all relevant information
	defer mocking.CloseRegistryKey(key)

	// Read the names of all subkeys
	psChildNames, subErr := key.ReadSubKeyNames(-1)
	if subErr != nil {
		return softwareList, subErr
	}
	// Iterate over each subkey
	for _, psChildName := range psChildNames {
		// Open the subkey
		childKey, childErr := mocking.OpenRegistryKey(key, psChildName)
		if childErr != nil {
			logger.Log.Error("Error opening device subkey " + psChildName)
			continue
		}

		// Close the subkey after we have received all relevant information
		defer mocking.CloseRegistryKey(childKey)
		// Read the DisplayName value
		name, _, nameErr := childKey.GetStringValue("DisplayName")
		if nameErr != nil {
			logger.Log.Debug("Error reading program name " + psChildName)
			continue
		}
		// Read the DisplayVersion value
		version, _, versionErr := childKey.GetStringValue("DisplayVersion")
		if versionErr != nil {
			logger.Log.Debug("Error reading program version " + psChildName)
			continue
		}
		// Read the Publisher value
		publisher, _, publisherErr := childKey.GetStringValue("Publisher")
		if publisherErr != nil {
			logger.Log.Debug("Error reading program publisher " + psChildName)
			continue
		}
		// Append the software to the list
		softwareList = append(softwareList, software{
			name:         name,
			identifier:   psChildName,
			version:      version,
			newVersion:   "",
			vendor:       publisher,
			lastUpdated:  "",
			sourceWinget: "",
			whereFrom:    bitOrigin,
		})
	}
	return softwareList, nil
}

// TODO: Update documentation
// software is a struct that represents a software installed on the system.
// It contains the following fields:
// - name: A string that represents the name of the software.
// - identifier: A string that represents the identifier of the software.
// - version: A string that represents the version of the software.
// - newVersion: A string that represents the new version of the software if available.
// - vendor: A string that represents the vendor of the software.
// - lastUpdated: A string that represents the last updated date of the software.
// - sourceWinget: A string that represents the source of the software if it was installed using the winget package manager.
// - whereFrom: A string that represents the source of the software (e.g., "64-bit", "winget").
type software struct {
	name         string
	identifier   string
	version      string
	newVersion   string
	vendor       string
	lastUpdated  string
	sourceWinget string // For possibly updating the software
	whereFrom    string // tmp for which function found this software
}

// TODO: Update documentation
// substr is a function that returns a substring from the input string.
// It first converts the input string to a slice of runes to handle multi-byte characters correctly.
// It then slices the rune slice from the start index to the start index plus the length.
// If the start index is greater than the length of the rune slice, it returns an empty string.
// If the start index plus the length is greater than the length of the rune slice, it adjusts the length to the end of the rune slice.
//
// Parameters:
// input: A string that represents the input string.
// start: An integer that represents the start index of the substring.
// length: An integer that represents the length of the substring.
//
// Returns:
// string: A string that represents the substring.
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
