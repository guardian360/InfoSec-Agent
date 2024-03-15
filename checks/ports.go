// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// OpenPorts checks for open ports and the processes that are using them
//
// Parameters: _
//
// Returns: A list of open ports and the processes that are using them
func OpenPorts() Check {
	// Regular expression to clean up multiple spaces in the output
	re := regexp.MustCompile("  +")

	// Get process ids (pids) and names of all processes
	cmd := exec.Command("tasklist")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()

	if err != nil {
		return newCheckErrorf("OpenPorts", "error running tasklist", err)
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
	cmd = exec.Command("netstat", "-ano")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err = cmd.Output()

	if err != nil {
		return newCheckErrorf("OpenPorts", "error running netstat", err)
	}
	netstat := strings.Split(string(output), "\n")

	result := newCheckResult("OpenPorts")
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
			result.Result = append(result.Result, fmt.Sprintf("port: %s, process: %s", port, name))
		}
	}

	return result
}
