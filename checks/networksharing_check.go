package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

// Networksharing checks if network sharing is enabled or disabled
func Networksharing() {
	output, err := exec.Command("powershell", "Get-NetAdapterBinding | Where-Object {$_.ComponentID -eq 'ms_server'} | Select-Object Enabled").Output()
	if err != nil {
		fmt.Println(err)
	}
	// Loops keeps count of the number of times "True" appears in the output for each network adapter
	outputString := strings.Split(string(output), "\r\n")
	counter := 0
	total := len(outputString) - 6
	for _, line := range outputString[3 : len(outputString)-3] {
		if strings.Contains(line, "True") {
			counter++
		}
	}
	if counter == total {
		fmt.Println("Network sharing is enabled")
	} else if counter > 0 && counter < total {
		fmt.Println("Network sharing is partially enabled")
	} else {
		fmt.Println("Network sharing is disabled")
	}
}
