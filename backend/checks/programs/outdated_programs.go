package programs

import (
	"fmt"
	"os/exec"
	"strings"
)

func OutdatedSoftware() {
	m := make(map[string]software)
	var softwareList []software
	var softwareListWinget, softwareList32, softwareList64, softwareListPackages, softwareListAppPackages []software
	var err error
	softwareListWinget, err = retrieveWingetInstalledPrograms(softwareListWinget)
	if err != nil {
		fmt.Println("Error retrieving installed programs:", err)
		return
	}
	softwareList32, err = retrieveInstalled32BitPrograms(softwareList32)
	if err != nil {
		fmt.Println("Error retrieving 32-bit installed programs:", err)
		return
	}
	softwareList64, err = retrieveInstalled64BitPrograms(softwareList64)
	if err != nil {
		fmt.Println("Error retrieving 64-bit installed programs:", err)
		return
	}
	softwareListPackages, err = retrieveInstalledPackages(softwareListPackages)
	if err != nil {
		fmt.Println("Error retrieving installed packages:", err)
		return
	}
	softwareListAppPackages, err = retrieveInstalledAppPackages(softwareListAppPackages)
	if err != nil {
		fmt.Println("Error retrieving installed app packages:", err)
		return
	}
	softwareList = append(softwareList, softwareListWinget...)
	softwareList = append(softwareList, softwareList32...)
	softwareList = append(softwareList, softwareList64...)
	softwareList = append(softwareList, softwareListPackages...)
	softwareList = append(softwareList, softwareListAppPackages...)
	runThis := true
	if runThis {
		var everythingButWinget, everythingBut32, everythingBut64, everythingButPackages, everythingButAppPackages []software
		everythingButWinget = append(everythingButWinget, softwareList32...)
		everythingButWinget = append(everythingButWinget, softwareList64...)
		everythingButWinget = append(everythingButWinget, softwareListPackages...)
		everythingButWinget = append(everythingButWinget, softwareListAppPackages...)
		everythingBut32 = append(everythingBut32, softwareListWinget...)
		everythingBut32 = append(everythingBut32, softwareList64...)
		everythingBut32 = append(everythingBut32, softwareListPackages...)
		everythingBut32 = append(everythingBut32, softwareListAppPackages...)
		everythingBut64 = append(everythingBut64, softwareListWinget...)
		everythingBut64 = append(everythingBut64, softwareList32...)
		everythingBut64 = append(everythingBut64, softwareListPackages...)
		everythingBut64 = append(everythingBut64, softwareListAppPackages...)
		everythingButPackages = append(everythingButPackages, softwareListWinget...)
		everythingButPackages = append(everythingButPackages, softwareList32...)
		everythingButPackages = append(everythingButPackages, softwareList64...)
		everythingButPackages = append(everythingButPackages, softwareListAppPackages...)
		everythingButAppPackages = append(everythingButAppPackages, softwareListWinget...)
		everythingButAppPackages = append(everythingButAppPackages, softwareList32...)
		everythingButAppPackages = append(everythingButAppPackages, softwareList64...)
		everythingButAppPackages = append(everythingButAppPackages, softwareListPackages...)
		fmt.Println("Winget:", len(softwareListWinget))
		fmt.Println("32-bit:", len(softwareList32))
		fmt.Println("64-bit:", len(softwareList64))
		fmt.Println("Packages:", len(softwareListPackages))
		fmt.Println("App Packages:", len(softwareListAppPackages))

		uniqueWinget := 0
		unique32 := 0
		unique64 := 0
		uniquePackages := 0
		uniqueAppPackages := 0
		var uniqueWingetList, unique32List, unique64List, uniquePackagesList, uniqueAppPackagesList []software

		winget := make(map[string]software)
		bit32 := make(map[string]software)
		bit64 := make(map[string]software)
		packages := make(map[string]software)
		appPackages := make(map[string]software)
		for i := range everythingButWinget {
			_, ok := winget[everythingButWinget[i].name]
			if !ok {
				winget[everythingButWinget[i].name] = everythingButWinget[i]
			}
		}
		for i := range softwareListWinget {
			_, ok := winget[softwareListWinget[i].name]
			if !ok {
				uniqueWinget++
				uniqueWingetList = append(uniqueWingetList, softwareListWinget[i])
			}
		}
		for i := range everythingBut32 {
			_, ok := bit32[everythingBut32[i].name]
			if !ok {
				bit32[everythingBut32[i].name] = everythingBut32[i]
			}
		}
		for i := range softwareList32 {
			_, ok := bit32[softwareList32[i].name]
			if !ok {
				unique32++
				unique32List = append(unique32List, softwareList32[i])
			}
		}
		for i := range everythingBut64 {
			_, ok := bit64[everythingBut64[i].name]
			if !ok {
				bit64[everythingBut64[i].name] = everythingBut64[i]
			}
		}
		for i := range softwareList64 {
			_, ok := bit64[softwareList64[i].name]
			if !ok {
				unique64++
				unique64List = append(unique64List, softwareList64[i])
			}
		}
		for i := range everythingButPackages {
			_, ok := packages[everythingButPackages[i].name]
			if !ok {
				packages[everythingButPackages[i].name] = everythingButPackages[i]
			}
		}
		for i := range softwareListPackages {
			_, ok := packages[softwareListPackages[i].name]
			if !ok {
				uniquePackages++
				uniquePackagesList = append(uniquePackagesList, softwareListPackages[i])
			}
		}
		for i := range everythingButAppPackages {
			_, ok := appPackages[everythingButAppPackages[i].name]
			if !ok {
				appPackages[everythingButAppPackages[i].name] = everythingButAppPackages[i]
			}
		}
		for i := range softwareListAppPackages {
			_, ok := appPackages[softwareListAppPackages[i].name]
			if !ok {
				uniqueAppPackages++
				uniqueAppPackagesList = append(uniqueAppPackagesList, softwareListAppPackages[i])
			}
		}
		fmt.Println("Unique Winget:", uniqueWinget)
		fmt.Println("Unique 32-bit:", unique32)
		fmt.Println("Unique 64-bit:", unique64)
		fmt.Println("Unique Packages:", uniquePackages)
		fmt.Println("Unique App Packages:", uniqueAppPackages)

		var wingetOnlyMap, bit32OnlyMap, bit64OnlyMap, packagesOnlyMap, appPackagesOnlyMap map[string]software
		wingetOnlyMap = make(map[string]software)
		bit32OnlyMap = make(map[string]software)
		bit64OnlyMap = make(map[string]software)
		packagesOnlyMap = make(map[string]software)
		appPackagesOnlyMap = make(map[string]software)

		for i := range softwareListWinget {
			_, ok := wingetOnlyMap[softwareListWinget[i].name]
			if !ok {
				wingetOnlyMap[softwareListWinget[i].name] = softwareListWinget[i]
			}
		}
		for i := range softwareList32 {
			_, ok := bit32OnlyMap[softwareList32[i].name]
			if !ok {
				bit32OnlyMap[softwareList32[i].name] = softwareList32[i]
			}
		}
		for i := range softwareList64 {
			_, ok := bit64OnlyMap[softwareList64[i].name]
			if !ok {
				bit64OnlyMap[softwareList64[i].name] = softwareList64[i]
			}
		}
		for i := range softwareListPackages {
			_, ok := packagesOnlyMap[softwareListPackages[i].name]
			if !ok {
				packagesOnlyMap[softwareListPackages[i].name] = softwareListPackages[i]
			}
		}
		for i := range softwareListAppPackages {
			_, ok := appPackagesOnlyMap[softwareListAppPackages[i].name]
			if !ok {
				appPackagesOnlyMap[softwareListAppPackages[i].name] = softwareListAppPackages[i]
			}
		}
		fmt.Println("Winget Only:", len(wingetOnlyMap))
		fmt.Println("32-bit Only:", len(bit32OnlyMap))
		fmt.Println("64-bit Only:", len(bit64OnlyMap))
		fmt.Println("Packages Only:", len(packagesOnlyMap))
		fmt.Println("App Packages Only:", len(appPackagesOnlyMap))

		var overlapWinget32, overlapWinget64, overlapWingetPackages, overlapWingetAppPackages, overlap3264, overlap32Packages, overlap32AppPackages, overlap64Packages, overlap64AppPackages, overlapPackagesAppPackages map[string]software
		overlapWinget32 = make(map[string]software)
		overlapWinget64 = make(map[string]software)
		overlapWingetPackages = make(map[string]software)
		overlapWingetAppPackages = make(map[string]software)
		overlap3264 = make(map[string]software)
		overlap32Packages = make(map[string]software)
		overlap32AppPackages = make(map[string]software)
		overlap64Packages = make(map[string]software)
		overlap64AppPackages = make(map[string]software)
		overlapPackagesAppPackages = make(map[string]software)

		for k, v := range wingetOnlyMap {
			_, ok := bit32OnlyMap[k]
			if ok {
				overlapWinget32[k] = v
			}
			_, ok = bit64OnlyMap[k]
			if ok {
				overlapWinget64[k] = v
			}
			_, ok = packagesOnlyMap[k]
			if ok {
				overlapWingetPackages[k] = v
			}
			_, ok = appPackagesOnlyMap[k]
			if ok {
				overlapWingetAppPackages[k] = v
			}
		}
		for k, v := range bit32OnlyMap {
			_, ok := bit64OnlyMap[k]
			if ok {
				overlap3264[k] = v
			}
			_, ok = packagesOnlyMap[k]
			if ok {
				overlap32Packages[k] = v
			}
			_, ok = appPackagesOnlyMap[k]
			if ok {
				overlap32AppPackages[k] = v
			}
		}
		for k, v := range bit64OnlyMap {
			_, ok := packagesOnlyMap[k]
			if ok {
				overlap64Packages[k] = v
			}
			_, ok = appPackagesOnlyMap[k]
			if ok {
				overlap64AppPackages[k] = v
			}
		}
		for k, v := range packagesOnlyMap {
			_, ok := appPackagesOnlyMap[k]
			if ok {
				overlapPackagesAppPackages[k] = v
			}
		}

		for i := range uniqueWingetList {
			fmt.Printf("Name: %s | Version: %s | Id: %s \n", uniqueWingetList[i].name, uniqueWingetList[i].version, uniqueWingetList[i].identifier)
		}
		counter := 0
		for i := range softwareList {
			if softwareList[i].version == "" {
				counter++
			}
		}
		fmt.Println("Empty string counter:", counter)
		fmt.Println("Done")
		return
	}
	counter := 0
	var duplicates []software
	for i := range softwareList {
		val, ok := m[softwareList[i].name]
		if ok {
			//fmt.Println("Found duplicate:", softwareList[i].name, val.version, softwareList[i].version)
			if softwareList[i].version == "" && val.version == "" {
				counter++
			}
			if softwareList[i].version != "" && val.version != "" {
				if softwareList[i].version != val.version {
					//fmt.Printf("Found duplicate: %s | %s from %s | %s from %s \n", softwareList[i].name, val.version, val.whereFrom, softwareList[i].version, softwareList[i].whereFrom)
					duplicates = append(duplicates, val)
					duplicates = append(duplicates, softwareList[i])
				}
				if softwareList[i].version > val.version {
					m[softwareList[i].name] = softwareList[i]
				}
			}
		} else {
			m[softwareList[i].name] = softwareList[i]
		}
	}
	fmt.Println("Installed Programs list:", len(softwareList))
	fmt.Println("Installed Programs map: ", len(m))
	fmt.Println("Duplicates of type empty string: ", counter)
	for _, v := range duplicates {
		fmt.Printf("Name: %s | Version: %s | Id: %s | Vendor: %s | whereFrom: %s \n", v.name, v.version, v.identifier, v.vendor, v.whereFrom)
	}
}

// this function takes the list and combines the elements that are the same
func combineSoftwareList(softwareList []software) ([]software, error) {
	return softwareList, nil
}

// This function retrieves all installed packages found with the winget package manager
func retrieveWingetInstalledPrograms(softwareList []software) ([]software, error) {
	out, err := exec.Command("winget", "list").Output()
	if err != nil {
		fmt.Printf("%s \n", err)
		return softwareList, err
	} else {
		lines := strings.Split(string(out), "\r\n")

		lines[0] = lines[0][strings.Index(lines[0], "Name"):] // Remove the first part of the header
		idIndex := strings.Index(lines[0], "Id")
		versionIndex := strings.Index(lines[0], "Version")
		availableIndex := strings.Index(lines[0], "Available")
		sourcesIndex := strings.Index(lines[0], "Source")
		for _, line := range lines[2:] { // Skip the header lines
			//fmt.Println(line)
			if len(line) != 0 { //Don't handle the last empty line, and maybe other empty lines
				name := line[:idIndex]
				name = strings.TrimSpace(name)
				//fmt.Println(name)
				id := line[idIndex : versionIndex-idIndex]
				id = strings.TrimSpace(id)
				//fmt.Println(id)
				version := line[versionIndex : availableIndex-versionIndex]
				version = strings.TrimSpace(version)
				//fmt.Println(version)
				available := line[availableIndex : sourcesIndex-availableIndex]
				available = strings.TrimSpace(available)
				//fmt.Println(available)
				source := line[sourcesIndex : len(line)-sourcesIndex]
				source = strings.TrimSpace(source)
				//fmt.Println(source)
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
	}
	return softwareList, nil
}

// this function returns all the installed 32-bit programs found using registry query
// command run: Get-ItemProperty "HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*" | Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List
func retrieveInstalled32BitPrograms(softwareList []software) ([]software, error) {
	output, err := exec.Command("powershell", "Get-ItemProperty ", "\"HKLM:\\SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*\"", "| Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List").Output()
	if err != nil {
		fmt.Println("Error retrieving 32-bit installed programs:", err)
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
				whereFrom:    "32-bit",
			})
			name, identifier, version, vendor = "", "", "", ""
		}
	}
	return softwareList, nil
}

// this function returns all the installed 64-bit programs found using registry query
// command run: Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*" | Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List
func retrieveInstalled64BitPrograms(softwareList []software) ([]software, error) {
	output, err := exec.Command("powershell", "Get-ItemProperty ", "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*", "| Select-Object DisplayName, PSChildName, DisplayVersion, Publisher | Sort-Object DisplayName | Format-List").Output()
	if err != nil {
		fmt.Println("Error retrieving 64-bit installed programs:", err)
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

// this function returns all the installed packages
// command run: Get-Package | Select-Object Name, TagId, Version | Sort-Object Name | Format-List
func retrieveInstalledPackages(softwareList []software) ([]software, error) {
	output, err := exec.Command("powershell", "Get-Package", "| Select-Object Name, TagId, Version | Sort-Object Name | Format-List").Output()
	if err != nil {
		fmt.Println("Error retrieving installed packages:", err)
		return softwareList, err
	}
	outputString := strings.Split(string(output), "\r\n")
	var name, identifier, version string
	for i, line := range outputString[2:] {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Name") {
			name = strings.Split(line, ":")[1]
			name = strings.TrimSpace(name)
		}
		if strings.Contains(line, "TagId") {
			identifier = strings.Split(line, ":")[1]
			identifier = strings.TrimSpace(identifier)
		}
		if strings.Contains(line, "Version") {
			version = strings.Split(line, ":")[1]
			version = strings.TrimSpace(version)
		}
		if i%4 == 3 {
			softwareList = append(softwareList, software{
				name:         name,
				identifier:   identifier,
				version:      version,
				newVersion:   "",
				vendor:       "",
				lastUpdated:  "",
				sourceWinget: "",
				whereFrom:    "Get-Package",
			})
			name, identifier, version = "", "", ""
		}
	}
	return softwareList, nil
}

// this function returns all the installed app packages
// command run: Get-AppxPackage | Select-Object Name, Version, Publisher | Sort-Object Name | Format-List
func retrieveInstalledAppPackages(softwareList []software) ([]software, error) {
	output, err := exec.Command("powershell", "Get-AppxPackage", "| Select-Object Name, Version, Publisher | Sort-Object Name | Format-List").Output()
	if err != nil {
		fmt.Println("Error retrieving installed app packages:", err)
		return softwareList, err
	}
	outputString := strings.Split(string(output), "\r\n")
	var name, version, vendor string
	for i, line := range outputString[2:] {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Name") {
			name = strings.Split(line, ":")[1]
			name = strings.TrimSpace(name)
		}
		if strings.Contains(line, "Version") {
			version = strings.Split(line, ":")[1]
			version = strings.TrimSpace(version)
		}
		if strings.Contains(line, "Publisher") {
			vendor = strings.Split(line, ":")[1]
			vendor = strings.TrimSpace(vendor)
		}
		if i%4 == 3 {
			softwareList = append(softwareList, software{
				name:         name,
				identifier:   "",
				version:      version,
				newVersion:   "",
				vendor:       vendor,
				lastUpdated:  "",
				sourceWinget: "",
				whereFrom:    "Get-AppxPackage",
			})
			name, version, vendor = "", "", ""
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
