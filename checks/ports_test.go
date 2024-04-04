package checks_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TestOpenPorts tests the OpenPorts function with different (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestOpenPorts(t *testing.T) {
	tests := []struct {
		name             string
		executortasklist *commandmock.MockCommandExecutor
		executornetstat  *commandmock.MockCommandExecutor
		want             checks.Check
	}{
		{
			name: "no open ports",
			executortasklist: &commandmock.MockCommandExecutor{
				Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat: &commandmock.MockCommandExecutor{
				Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n",
				Err:    nil},
			want: checks.NewCheckResult(10, 0),
		},
		{
			name: "1 open port",
			executortasklist: &commandmock.MockCommandExecutor{
				Output: "	\r\n Image Name	PID	Session Name	Session#	Mem Usage\r\n	" +
					"\r\nSystem Idle Process  0  Services  0  8 K\r\n", Err: nil},
			executornetstat: &commandmock.MockCommandExecutor{
				Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
					"TCP	123.0.0.1:8080	123:123	ESTABLISHED 0\r\n", Err: nil},
			want: checks.NewCheckResult(10, 0, "port: 8080, process: System Idle Process"),
		},
		{
			name: "multiple open ports",
			executortasklist: &commandmock.MockCommandExecutor{
				Output: "	\r\n Image Name	PID	Session Name	Session#	" +
					"Mem Usage\r\n	\r\n" +
					"System Idle Process  0  Services  0  8 K\r\n" +
					"System2  1 Services 0 8 K\r\nSystem3  2 Services 0 8 K\r\n",
				Err: nil},
			executornetstat: &commandmock.MockCommandExecutor{
				Output: " \r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
					"TCP 123.0.0.1:8080 123:123 ESTABLISHED 0\r\n" + "TCP 123.0.0.1:8081 123:123 ESTABLISHED 1\r\n" +
					"TCP 123.0.0.1:8082 123:123 ESTABLISHED 2\r\n", Err: nil},
			want: checks.NewCheckResult(10, 0,
				"port: 8080, process: System Idle Process", "port: 8081, process: System2",
				"port: 8082, process: System3"),
		},
		{
			name:             "tasklist error",
			executortasklist: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("")},
			executornetstat:  &commandmock.MockCommandExecutor{Output: "", Err: nil},
			want: checks.NewCheckErrorf(10,
				"error running tasklist", errors.New("")),
		},
		{
			name: "netstat error",
			executortasklist: &commandmock.MockCommandExecutor{
				Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("")},
			want: checks.NewCheckErrorf(10,
				"error running netstat", errors.New("")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.OpenPorts(tt.executortasklist, tt.executornetstat)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestPortsExpected
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
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
			executor := &commandmock.RealCommandExecutor{}
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
