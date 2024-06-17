package programs

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

// OutdatedSoftware function checks for outdated software on the system
func OutdatedSoftware() checks.Check {
	// Collect all software lists
	softwareList, err := collectAllSoftwareLists()
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

// collectAllSoftwareLists function collects all software lists from different sources
func collectAllSoftwareLists() ([]software, checks.Check) {
	var (
		softwareList       []software
		softwareListWinget []software
		softwareList32     []software
		softwareList64     []software
		err                error
	)

	// Retrieve installed programs from winget
	if softwareListWinget, err = retrieveWingetInstalledPrograms(softwareList); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in Program Files", err)
	}
	// Retrieve installed 32 bit programs
	if softwareList32, err = retrieveInstalled32BitPrograms(softwareList); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 32 bit programs", err)
	}
	// Retrieve installed 64 bit programs
	if softwareList64, err = retrieveInstalled64BitPrograms(softwareList, "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*"); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 64bit programs", err)
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

// filterAndDeduplicateSoftware function filters and deduplicates the software list
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

// formatResultArray function formats the result array
func formatResultArray(uniqueSoftware map[string]software) []string {
	resultArray := make([]string, 0, len(uniqueSoftware))
	for _, v := range uniqueSoftware {
		resultArray = append(resultArray, fmt.Sprintf("%s | %s", v.name, v.version))
	}
	return resultArray
}

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

// compareVersions function to compare two version strings
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
				logger.Log.ErrorWithErr("error parsing version number", err)
			}
		}
		if i < len(parts2) {
			_, err := fmt.Sscanf(parts2[i], "%d", &num2)
			if err != nil {
				logger.Log.ErrorWithErr("error parsing version number", err)
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

// retrieveWingetInstalledPrograms function retrieves all installed packages found with the winget package manager
func retrieveWingetInstalledPrograms(softwareList []software) ([]software, error) {
	// Execute the winget list command
	// winget list | Out-String -Stream | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Text.Encoding]::Default.GetBytes($_)) }
	out, err := exec.Command("powershell", "winget list| Out-String -Stream | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Text.Encoding]::Default.GetBytes($_)) }").Output()
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
	lines[0] = lines[indexN][strings.Index(lines[indexN], "Name"):] // Remove the first part of the header
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

// retrieveInstalled32BitPrograms function returns all the installed 32-bit programs found using registry query
func retrieveInstalled32BitPrograms(softwareList []software) ([]software, error) {
	return retrieveInstalled64BitPrograms(softwareList, "\"HKLM:\\SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*\"")
}

// retrieveInstalled64BitPrograms function returns all the installed 64-bit programs found using registry query
func retrieveInstalled64BitPrograms(softwareList []software, bits string) ([]software, error) {
	// Execute the powershell command to get the installed programs
	output, err := exec.Command("powershell", "Get-ItemProperty ", bits, "| Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List").Output()
	if err != nil {
		return softwareList, err
	}
	// Process the output
	outputString := strings.Split(string(output), "\r\n")
	var name, identifier, version, vendor string
	for i, line := range outputString[2:] {
		line = strings.TrimSpace(line)
		// Extract the software details
		if strings.Contains(line, "DisplayName") {
			name = strings.Split(line, ":")[1]
			name = strings.TrimSpace(name)
		}
		if strings.Contains(line, "PSChildName") {
			identifier = strings.Split(line, ":")[1]
			identifier = strings.TrimSpace(identifier)
		}
		if strings.Contains(line, "DisplayVersion") {
			version = strings.Split(line, ":")[1]
			version = strings.TrimSpace(version)
		}
		if strings.Contains(line, "Publisher") {
			vendor = strings.Split(line, ":")[1]
			vendor = strings.TrimSpace(vendor)
		}
		// Append the software to the list
		if i%5 == 4 {
			softwareList = append(softwareList, software{
				name:         name,
				identifier:   identifier,
				version:      version,
				newVersion:   "",
				vendor:       vendor,
				lastUpdated:  "",
				sourceWinget: "",
				whereFrom:    "64-bit",
			})
			name, identifier, version, vendor = "", "", "", ""
		}
	}
	return softwareList, nil
}

// software struct represents a software
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

// substr function returns a substring from the input string
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
