package checks

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"reflect"
	"testing"
)

func TestNetworkSharing(t *testing.T) {
	tests := []struct {
		name     string
		executor utils.CommandExecutor
		want     Check
	}{
		// TODO: Add test cases.
		{
			name:     "Get-NetAdapterBinding command error",
			executor: &utils.MockCommandExecutor{Output: "", Err: errors.New("error executing command Get-NetAdapterBinding")},
			want:     NewCheckErrorf("NetworkSharing", "error executing command Get-NetAdapterBinding", errors.New("error executing command Get-NetAdapterBinding")),
		},
		{
			name:     "Network sharing is enabled",
			executor: &utils.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nTrue\r\nTrue\r\n\r\n\r\n", Err: nil},
			want:     NewCheckResult("NetworkSharing", "Network sharing is enabled"),
		},
		{
			name:     "Network sharing is partially enabled",
			executor: &utils.MockCommandExecutor{Output: "\r\n\r\n\r\nTrue\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     NewCheckResult("NetworkSharing", "Network sharing is partially enabled"),
		},
		{
			name:     "Network sharing is disabled",
			executor: &utils.MockCommandExecutor{Output: "\r\n\r\n\r\nFalse\r\n\r\n\r\n", Err: nil},
			want:     NewCheckResult("NetworkSharing", "Network sharing is disabled"),
		},
		{
			name:     "Network sharing status is unknown",
			executor: &utils.MockCommandExecutor{Output: "\r\n\r\n\r\nHelloWorld\r\n\r\n\r\n", Err: nil},
			want:     NewCheckResult("NetworkSharing", "Network sharing status is unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NetworkSharing(tt.executor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworkSharing() = %v, want %v", got, tt.want)
			}
		})
	}
}
