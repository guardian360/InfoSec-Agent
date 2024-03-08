package scan

import (
	"InfoSec-Agent/checks"
	"encoding/json"
	"fmt"
	"os"
)

func Scan() {
	// Run all checks
	smb := checks.SmbCheck()
	secureBoot := checks.SecureBoot()

	// Combine results
	checkResults := []checks.Check{
		smb,
		secureBoot,
	}

	// Serialize check results to JSON
	jsonData, err := json.Marshal(checkResults)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

	// Write JSON data to a file
	file, err := os.Create("checks.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON data to file:", err)
		return
	}
}
