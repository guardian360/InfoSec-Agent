package programs

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
)

// func OutdatedSoftware() checks.Check {
// 	var softwareList []software
// 	var err error

// 	softwareListWinget, err := retrieveSoftwareList(retrieveWingetInstalledPrograms, "error listing installed programs in Program Files")
// 	if err != nil {
// 		return checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in Program Files", err)
// 	}
// 	softwareList32, err = retrieveInstalled32BitPrograms(softwareList32)
// 	if err != nil {
// 		return checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 32 bit programs", err)
// 	}
// 	softwareList64, err = retrieveInstalled64BitPrograms(softwareList64, "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*")
// 	if err != nil {
// 		return checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 64bit programs", err)
// 	}
// 	softwareListPackages, err = retrieveInstalledPackages(softwareListPackages, "Get-Package | Select-Object Name, TagId, Version | Sort-Object Name | Format-List")
// 	if err != nil {
// 		return checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in packages", err)
// 	}
// 	softwareListAppPackages, err = retrieveInstalledAppPackages(softwareListAppPackages)
// 	if err != nil {
// 		return checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in app packages", err)
// 	}

// 	softwareList = appendSoftwareLists(softwareList, softwareListWinget, softwareList32, softwareList64, softwareListPackages, softwareListAppPackages)

// 	everythingBut := generateEverythingButLists(softwareListWinget, softwareList32, softwareList64, softwareListPackages, softwareListAppPackages)

// 	uniqueSoftware := filterUniqueSoftware(softwareList)

// 	resultArray := formatSoftwareList(uniqueSoftware)

// 	return checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, resultArray...)
// }

// func retrieveSoftwareList(retrieveFunc func([]software) ([]software, error), errMsg string, args ...string) ([]software, error) {
// 	var list []software
// 	var err error
// 	if len(args) > 0 {
// 		list, err = retrieveFunc(list, args[0])
// 	} else {
// 		list, err = retrieveFunc(list)
// 	}
// 	if err != nil {
// 		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, errMsg, err)
// 	}
// 	return list, nil
// }

// func appendSoftwareLists(lists ...[]software) []software {
// 	var combined []software
// 	for _, list := range lists {
// 		combined = append(combined, list...)
// 	}
// 	return combined
// }

// func generateEverythingButLists(softwareListWinget, softwareList32, softwareList64, softwareListPackages, softwareListAppPackages []software) map[string][]software {
// 	everythingBut := make(map[string][]software)
// 	everythingBut["winget"] = append(append(append(append([]software{}, softwareList32...), softwareList64...), softwareListPackages...), softwareListAppPackages...)
// 	everythingBut["32"] = append(append(append(append([]software{}, softwareListWinget...), softwareList64...), softwareListPackages...), softwareListAppPackages...)
// 	everythingBut["64"] = append(append(append(append([]software{}, softwareListWinget...), softwareList32...), softwareListPackages...), softwareListAppPackages...)
// 	everythingBut["packages"] = append(append(append(append([]software{}, softwareListWinget...), softwareList32...), softwareList64...), softwareListAppPackages...)
// 	everythingBut["appPackages"] = append(append(append(append([]software{}, softwareListWinget...), softwareList32...), softwareList64...), softwareListPackages...)

// 	return everythingBut
// }

// func filterUniqueSoftware(softwareList []software) map[string]software {
// 	uniqueSoftware := make(map[string]software)

// 	for _, sw := range softwareList {
// 		if sw.name == "" || sw.version == "" || strings.Contains(strings.ToLower(sw.name), "microsoft defender") {
// 			continue
// 		}

// 		normalized := normalize(sw.name)

// 		if existing, exists := uniqueSoftware[normalized]; exists {
// 			if compareVersions(sw.version, existing.version) > 0 {
// 				uniqueSoftware[normalized] = sw
// 			}
// 		} else {
// 			uniqueSoftware[normalized] = sw
// 		}
// 	}

// 	return uniqueSoftware
// }

// func formatSoftwareList(uniqueSoftware map[string]software) []string {
// 	resultArray := make([]string, 0)
// 	for _, v := range uniqueSoftware {
// 		resultArray = append(resultArray, fmt.Sprintf("%s | %s", v.name, v.version))
// 	}
// 	return resultArray
// }

func OutdatedSoftware() checks.Check {
	softwareList, err := collectAllSoftwareLists()
	if softwareList == nil {
		return err
	}

	uniqueSoftware := filterAndDeduplicateSoftware(softwareList)
	resultArray := formatResultArray(uniqueSoftware)

	return checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, resultArray...)
}

func collectAllSoftwareLists() ([]software, checks.Check) {
	var (
		softwareList       []software
		softwareListWinget []software
		softwareList32     []software
		softwareList64     []software
		err                error
	)

	if softwareListWinget, err = retrieveWingetInstalledPrograms(softwareList); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in Program Files", err)
	}
	if softwareList32, err = retrieveInstalled32BitPrograms(softwareList); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 32 bit programs", err)
	}
	if softwareList64, err = retrieveInstalled64BitPrograms(softwareList, "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*"); err != nil {
		return nil, checks.NewCheckErrorf(checks.OutdatedSoftwareID, "error listing installed programs in 64bit programs", err)
	}

	softwareList = append(softwareList, softwareListWinget...)
	softwareList = append(softwareList, softwareList32...)
	softwareList = append(softwareList, softwareList64...)

	resultArray := make([]string, 0)
	for _, v := range softwareList {
		resultArray = append(resultArray, fmt.Sprintf("%s | %s", v.name, v.version))
	}

	return softwareList, checks.NewCheckResult(checks.OutdatedSoftwareID, checks.OutdatedSoftwareID, resultArray...)
}

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
			fmt.Sscanf(parts1[i], "%d", &num1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &num2)
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}
	return 0
}

// This function retrieves all installed packages found with the winget package manager
func retrieveWingetInstalledPrograms(softwareList []software) ([]software, error) {
	out, err := exec.Command("winget", "list").Output()
	if err != nil {
		// fmt.Printf("%s \n", err)
		return softwareList, err
	}
	lines := strings.Split(string(out), "\r\n")
	lines[0] = lines[0][strings.Index(lines[0], "Name")+1:] // Remove the first part of the header
	idIndex := strings.Index(lines[0], "Id")
	versionIndex := strings.Index(lines[0], "Version")
	availableIndex := strings.Index(lines[0], "Available")
	sourcesIndex := strings.Index(lines[0], "Source")
	for _, line := range lines[2:] { // Skip the header lines
		// fmt.Println(line)
		if len(line) != 0 { // Don't handle the last empty line, and maybe other empty lines
			name := substr(line, 0, idIndex)
			name = strings.TrimSpace(name)
			// fmt.Println(name)
			id := substr(line, idIndex, versionIndex-idIndex)
			id = strings.TrimSpace(id)
			// fmt.Println(id)
			version := substr(line, versionIndex, availableIndex-versionIndex)
			version = strings.TrimSpace(version)
			// fmt.Println(version)
			available := substr(line, availableIndex, sourcesIndex-availableIndex)
			available = strings.TrimSpace(available)
			// fmt.Println(available)
			source := substr(line, sourcesIndex, len(line)-sourcesIndex)
			source = strings.TrimSpace(source)
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

// this function returns all the installed 32-bit programs found using registry query
// command run: Get-ItemProperty "HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*" | Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List
func retrieveInstalled32BitPrograms(softwareList []software) ([]software, error) {
	return retrieveInstalled64BitPrograms(softwareList, "\"HKLM:\\SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*\"")
}

// this function returns all the installed 64-bit programs found using registry query
// command run: Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*" | Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List
func retrieveInstalled64BitPrograms(softwareList []software, bits string) ([]software, error) {
	output, err := exec.Command("powershell", "Get-ItemProperty ", bits, "| Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List").Output()
	if err != nil {
		// fmt.Println("Error retrieving 64-bit installed programs:", err)
		return softwareList, err
	}
	outputString := strings.Split(string(output), "\r\n")
	var name, identifier, version, vendor string
	for i, line := range outputString[2:] {
		line = strings.TrimSpace(line)
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
