package checks

import (
	"os/exec"
	"strings"
)

// NetworkSharing checks if network sharing is enabled or disabled
func NetworkSharing() Check {
	output, err := exec.Command("powershell", "Get-NetAdapterBinding | Where-Object {$_.ComponentID -eq 'ms_server'} | Select-Object Enabled").Output()
	if err != nil {
		return newCheckErrorf("NetworkSharing", "error executing command Get-NetAdapterBinding", err)
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
		return newCheckResult("NetworkSharing", "Network sharing is enabled")
	}
	if counter > 0 && counter < total {
		return newCheckResult("NetworkSharing", "Network sharing is partially enabled")
	}
	return newCheckResult("NetworkSharing", "Network sharing is disabled")
}
