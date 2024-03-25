package checks

import (
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"strings"
)

// NetworkSharing checks if network sharing is enabled or disabled
//
// Parameters: _
//
// Returns: If network sharing is enabled or not
func NetworkSharing(executor utils.CommandExecutor) Check {
	// Execute a powershell command to get the network adapter binding status
	command := fmt.Sprintf("Get-NetAdapterBinding | Where-Object {$_.ComponentID -eq 'ms_server'} | Select-Object Enabled")
	output, err := executor.Execute("powershell", command)
	if err != nil {
		return NewCheckErrorf("NetworkSharing",
			"error executing command Get-NetAdapterBinding", err)
	}

	outputString := strings.Split(string(output), "\r\n")
	trueCounter := 0               // Counter to keep track of the number of enabled network adapters
	falseCounter := 0              // Counter to keep track of the number of disabled network adapters
	total := len(outputString) - 6 // Expected number of enabled network adapters for network sharing to be enabled

	for _, line := range outputString[3 : len(outputString)-3] {
		// Check if the line contains "True" indicating network sharing is enabled for the adapter
		if strings.Contains(line, "True") {
			trueCounter++
		} else if strings.Contains(line, "False") {
			falseCounter++
		}
	}

	// Check the status of network sharing based on the number of enabled network adapters
	if trueCounter == total && falseCounter == 0 {
		return NewCheckResult("NetworkSharing", "Network sharing is enabled")
	}
	if trueCounter > 0 && trueCounter < total && falseCounter > 0 {
		return NewCheckResult("NetworkSharing", "Network sharing is partially enabled")
	}
	if trueCounter == 0 && falseCounter == total {
		return NewCheckResult("NetworkSharing", "Network sharing is disabled")
	}
	return NewCheckResult("NetworkSharing", "Network sharing status is unknown")
}
