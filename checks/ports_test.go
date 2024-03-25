package checks_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"reflect"
	"testing"
)

// TestOpenPorts tests the OpenPorts function with different (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestOpenPorts(t *testing.T) {
	tests := []struct {
		name             string
		executortasklist *utils.MockCommandExecutor
		executornetstat  *utils.MockCommandExecutor
		want             checks.Check
	}{
		{
			name:             "no open ports",
			executortasklist: &utils.MockCommandExecutor{Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat:  &utils.MockCommandExecutor{Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n", Err: nil},
			want:             checks.NewCheckResult("OpenPorts"),
		},
		{
			name:             "1 open port",
			executortasklist: &utils.MockCommandExecutor{Output: "	\r\n Image Name	PID	Session Name	Session#	Mem Usage\r\n	\r\nSystem Idle Process  0  Services  0  8 K\r\n", Err: nil},
			executornetstat: &utils.MockCommandExecutor{Output: "	\r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
				"TCP	123.0.0.1:8080	123:123	ESTABLISHED 0\r\n", Err: nil},
			want: checks.NewCheckResult("OpenPorts", "port: 8080, process: System Idle Process"),
		},
		{
			name: "multiple open ports",
			executortasklist: &utils.MockCommandExecutor{Output: "	\r\n Image Name	PID	Session Name	Session#	" +
				"Mem Usage\r\n	\r\n" +
				"System Idle Process  0  Services  0  8 K\r\nSystem2  1 Services 0 8 K\r\nSystem3  2 Services 0 8 K\r\n",
				Err: nil},
			executornetstat: &utils.MockCommandExecutor{Output: " \r\n Active Connections \r\n	\r\n	Proto	Local Address	Foreign Address	State\r\n" +
				"TCP 123.0.0.1:8080 123:123 ESTABLISHED 0\r\n" + "TCP 123.0.0.1:8081 123:123 ESTABLISHED 1\r\n" +
				"TCP 123.0.0.1:8082 123:123 ESTABLISHED 2\r\n", Err: nil},
			want: checks.NewCheckResult("OpenPorts", "port: 8080, process: System Idle Process", "port: 8081, process: System2", "port: 8082, process: System3"),
		},
		{
			name:             "tasklist error",
			executortasklist: &utils.MockCommandExecutor{Output: "", Err: errors.New("")},
			executornetstat:  &utils.MockCommandExecutor{Output: "", Err: nil},
			want:             checks.NewCheckErrorf("OpenPorts", "error running tasklist", errors.New("")),
		},
		{
			name:             "netstat error",
			executortasklist: &utils.MockCommandExecutor{Output: "	\r\n Image Name 	PID 	Session Name 	Session# 	Mem Usage\r\n	\r\n", Err: nil},
			executornetstat:  &utils.MockCommandExecutor{Output: "", Err: errors.New("")},
			want:             checks.NewCheckErrorf("OpenPorts", "error running netstat", errors.New("")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.OpenPorts(tt.executortasklist, tt.executornetstat)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}