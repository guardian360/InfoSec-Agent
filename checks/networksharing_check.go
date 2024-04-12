package checks

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
)

// NetworkSharing is a function that checks the status of network sharing on the system.
//
// Parameters:
//   - executor (mocking.CommandExecutor): An instance of CommandExecutor used to execute the PowerShell command to retrieve the network adapter binding status.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the network sharing check. The Result field of the Check instance will contain one of the following messages:
//   - "Network sharing is enabled" if all network adapters have network sharing enabled.
//   - "Network sharing is partially enabled" if some network adapters have network sharing enabled and others do not.
//   - "Network sharing is disabled" if no network adapters have network sharing enabled.
//   - "Network sharing status is unknown" if the function is unable to determine the status of network sharing.
//
// This function is primarily used to identify potential security risks associated with network sharing on the system.
func NetworkSharing(executor mocking.CommandExecutor) Check {
	// Execute a powershell command to get the network adapter binding status
	command := "Get-NetAdapterBinding | Where-Object {$_.ComponentID -eq 'ms_server'} | Select-Object Enabled"
	output, err := executor.Execute("powershell", command)
	if err != nil {
		return NewCheckErrorf(NetworkSharingID,
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
		return NewCheckResult(NetworkSharingID, 0, "Network sharing is enabled")
	}
	if trueCounter > 0 && trueCounter < total && falseCounter > 0 {
		return NewCheckResult(NetworkSharingID, 1, "Network sharing is partially enabled")
	}
	if trueCounter == 0 && falseCounter == total {
		return NewCheckResult(NetworkSharingID, 2, "Network sharing is disabled")
	}
	return NewCheckResult(NetworkSharingID, 3, "Network sharing status is unknown")
}
