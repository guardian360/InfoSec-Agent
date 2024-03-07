package checks

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func OpenPorts() {
	re := regexp.MustCompile("  +")

	// Get pids and names of all processes
	output, _ := exec.Command("tasklist").Output()
	tasklist := strings.Split(string(output), "\r\n")

	// Map of pids to process names
	pids := make(map[string]string)
	for _, line := range tasklist[3 : len(tasklist)-1] {
		words := re.Split(line, 2)[:2]
		words[1] = strings.Fields(words[1])[0]
		pids[words[1]] = words[0]
	}

	// Get all open ports
	cmd := exec.Command("netstat", "-ano")
	output, _ = cmd.Output()
	netstat := strings.Split(string(output), "\n")

	for _, line := range netstat[4 : len(netstat)-1] {
		words := strings.Fields(line)

		// Separate ip and port
		ip := strings.Split(words[1], ":")

		// Skip local host
		if ip[0] == "0.0.0.0" {
			continue
		}

		port := ip[len(ip)-1]

		// Get pid from line
		var pid string
		if words[0] == "TCP" {
			pid = words[4]
		} else {
			pid = words[3]
		}

		// Print process name from pid
		name, ok := pids[pid]
		if ok {
			fmt.Println("port:", port, "\tprocess:", name)
		}
	}
}
