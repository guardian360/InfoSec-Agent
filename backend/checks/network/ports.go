// Package network provides functions related to security/privacy checks of network settings
package network

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// OpenPorts is a function that checks for open ports on the system and identifies the processes that are using them.
//
// Parameters:
//   - tasklistexecutor (mocking.CommandExecutor): An executor to run the 'tasklist' command which retrieves the list of currently running tasks.
//   - netstatexecutor (mocking.CommandExecutor): An executor to run the 'netstat' command which provides network statistics.
//
// Returns:
//   - Check: A struct containing the result of the check. The result is a list of open ports along with the names of the processes that are using them.
//
// The function works by first running the 'tasklist' command to get a list of all running tasks. It then maps each process ID to its corresponding process name. Next, it runs the 'netstat' command to get a list of all open ports. For each open port, it identifies the process ID and maps it back to the process name using the previously created map. The function then returns a list of open ports along with the names of the processes that are using them.
func OpenPorts(tasklistexecutor, netstatexecutor mocking.CommandExecutor) checks.Check {
	// Regular expression to clean up multiple spaces in the output
	re := regexp.MustCompile("  +")

	// Get process ids (pids) and names of all processes
	command := "tasklist"
	output, err := tasklistexecutor.Execute(command)

	if err != nil {
		return checks.NewCheckErrorf(checks.PortsID, "error running tasklist", err)
	}
	tasklist := strings.Split(string(output), "\r\n")

	// Map each pid to its process name
	pids := make(map[string]string)
	for _, line := range tasklist[3 : len(tasklist)-1] {
		words := re.Split(line, 2)[:2]
		// Extract the process name
		words[1] = strings.Fields(words[1])[0]
		// Map the pid to its process name
		pids[words[1]] = words[0]
	}

	// Get all open ports
	command = "netstat"
	output, err = netstatexecutor.Execute(command, "-ano")

	if err != nil {
		return checks.NewCheckErrorf(checks.PortsID, "error running netstat", err)
	}
	netstat := strings.Split(string(output), "\n")

	result := checks.NewCheckResult(checks.PortsID, 0)
	processPorts := make(map[string][]string)
	for _, line := range netstat[4 : len(netstat)-1] {
		words := strings.Fields(line)

		// Separate ip and port
		ip := strings.Split(words[1], ":")

		// Skip local host
		if ip[0] == "0.0.0.0" {
			continue
		}

		// Extract the port from the ip
		port := ip[len(ip)-1]

		// Get pid from line depending on the protocol used (TCP or UDP)
		var pid string
		if words[0] == "TCP" {
			pid = words[4]
		} else {
			pid = words[3]
		}

		// Return the process name from the pid
		name, ok := pids[pid]
		if ok {
			processPorts[name] = append(processPorts[name], port)
		}
	}

	// Construct the output strings
	for name, ports := range processPorts {
		result.Result = append(result.Result, fmt.Sprintf("process: %s, port: %s", name, strings.Join(ports, ", ")))
	}

	return result
}
