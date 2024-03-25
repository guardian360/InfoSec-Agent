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
func TestUACCheck(t *testing.T) {
	tests := []struct {
		name        string
		executorUAC *utils.MockCommandExecutor
		want        checks.Check
	}{
		{
			name:        "UAC disabled",
			executorUAC: &utils.MockCommandExecutor{Output: "0", Err: nil},
			want:        checks.NewCheckResult("UAC", "UAC is disabled."),
		},
		{
			name:        "UAC enabled for apps and settings",
			executorUAC: &utils.MockCommandExecutor{Output: "2", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to your computer "+
				"and for changing your settings."),
		},
		{
			name:        "UAC enabled for apps but not for settings",
			executorUAC: &utils.MockCommandExecutor{Output: "5", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to "+
				"your computer."),
		},
		{
			name:        "unknown UAC level",
			executorUAC: &utils.MockCommandExecutor{Output: "3", Err: nil},
			want:        checks.NewCheckResult("UAC", "Unknown UAC level"),
		},
		{
			name:        "UAC error",
			executorUAC: &utils.MockCommandExecutor{Output: "", Err: errors.New("error retrieving UAC")},
			want: checks.NewCheckErrorf("UAC", "error retrieving UAC",
				errors.New("error retrieving UAC")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.UACCheck(tt.executorUAC)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}
