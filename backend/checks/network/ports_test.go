package network_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestOpenPorts is a function that tests the OpenPorts function's ability to correctly identify open ports and the processes using them.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the OpenPorts function with different scenarios. It uses mock implementations of the CommandExecutor interface to simulate the output of 'tasklist' and 'netstat' commands. Each test case checks if the OpenPorts function correctly identifies the open ports and the processes using them based on the simulated command outputs. The function asserts that the returned Check instance contains the expected results or error messages.
func TestOpenPorts(t *testing.T) {
	tests := []struct {
		name             string
		executortasklist *mocking.MockCommandExecutor
		executornetstat  *mocking.MockCommandExecutor
		want             checks.Check
	}{
		{
			name: "no open ports",
			executortasklist: &mocking.MockCommandExecutor{
				Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat: &mocking.MockCommandExecutor{
				Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n",
				Err:    nil},
			want: checks.NewCheckResult(checks.PortsID, 0),
		},
		{
			name: "1 open port",
			executortasklist: &mocking.MockCommandExecutor{
				Output: "	\r\n Image Name	PID	Session Name	Session#	Mem Usage\r\n	" +
					"\r\nSystem Idle Process  0  Services  0  8 K\r\n", Err: nil},
			executornetstat: &mocking.MockCommandExecutor{
				Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
					"TCP	123.0.0.1:8080	123:123	ESTABLISHED 0\r\n", Err: nil},
			want: checks.NewCheckResult(checks.PortsID, 0, "process: System Idle Process, port: 8080"),
		},
		{
			name: "multiple open ports",
			executortasklist: &mocking.MockCommandExecutor{
				Output: "	\r\n Image Name	PID	Session Name	Session#	" +
					"Mem Usage\r\n	\r\n" +
					"System Idle Process  0  Services  0  8 K\r\n" +
					"System2  1 Services 0 8 K\r\nSystem3  2 Services 0 8 K\r\n",
				Err: nil},
			executornetstat: &mocking.MockCommandExecutor{
				Output: " \r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
					"TCP 123.0.0.1:8080 123:123 ESTABLISHED 0\r\n" + "TCP 123.0.0.1:8081 123:123 ESTABLISHED 1\r\n" +
					"TCP 123.0.0.1:8082 123:123 ESTABLISHED 2\r\n", Err: nil},
			want: checks.NewCheckResult(checks.PortsID, 0,
				"process: System Idle Process, port: 8080",
				"process: System2, port: 8081",
				"process: System3, port: 8082"),
		},
		{
			name:             "tasklist error",
			executortasklist: &mocking.MockCommandExecutor{Output: "", Err: errors.New("")},
			executornetstat:  &mocking.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(checks.PortsID,
				"error running tasklist", errors.New("")),
		},
		{
			name: "netstat error",
			executortasklist: &mocking.MockCommandExecutor{
				Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat: &mocking.MockCommandExecutor{Output: "", Err: errors.New("")},
			want: checks.NewCheckErrorf(checks.PortsID,
				"error running netstat", errors.New("")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network.OpenPorts(tt.executortasklist, tt.executornetstat)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestPortsExpected is a function that tests if the output of the 'tasklist' and 'netstat' commands match the expected results.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the 'tasklist' and 'netstat' commands by comparing their output with the expected results. It uses the RealCommandExecutor to execute the commands and retrieve their output. The function then checks if the output matches the expected results. If the output does not match the expected results, the function reports a test failure.
func TestPortsExpected(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		arguments string
		expected  string
	}{
		{
			name:     "tasklist returns expected",
			command:  "tasklist",
			expected: "ImageNamePIDSessionNameSession#MemUsage",
		},
		{
			name:      "netstat returns expected",
			command:   "netstat",
			arguments: "-ano",
			expected:  "ActiveConnections",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &mocking.RealCommandExecutor{}
			var res []byte
			if tt.command == "netstat" {
				res, _ = executor.Execute(tt.command, tt.arguments)
			} else {
				res, _ = executor.Execute(tt.command)
			}
			outputList := strings.Split(string(res), "\r\n")
			if strings.ReplaceAll(outputList[1], " ", "") != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, outputList[1])
			}
		})
	}
}
