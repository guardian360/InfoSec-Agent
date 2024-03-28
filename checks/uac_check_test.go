package checks_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
)

// TestUACCheck tests the OpenPorts function with different (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestUACCheck(t *testing.T) {
	tests := []struct {
		name        string
		executorUAC *commandmock.MockCommandExecutor
		want        checks.Check
	}{
		{
			name:        "UAC disabled",
			executorUAC: &commandmock.MockCommandExecutor{Output: "0", Err: nil},
			want:        checks.NewCheckResult("UAC", "UAC is disabled."),
		},
		{
			name:        "UAC enabled for apps and settings",
			executorUAC: &commandmock.MockCommandExecutor{Output: "2", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to your computer "+
				"and for changing your settings."),
		},
		{
			name:        "UAC enabled for apps but not for settings",
			executorUAC: &commandmock.MockCommandExecutor{Output: "5", Err: nil},
			want: checks.NewCheckResult("UAC", "UAC is turned on for apps making changes to "+
				"your computer."),
		},
		{
			name:        "unknown UAC level",
			executorUAC: &commandmock.MockCommandExecutor{Output: "3", Err: nil},
			want:        checks.NewCheckResult("UAC", "Unknown UAC level"),
		},
		{
			name:        "UAC error",
			executorUAC: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("error retrieving UAC")},
			want: checks.NewCheckErrorf("UAC", "error retrieving UAC",
				errors.New("error retrieving UAC")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.UACCheck(tt.executorUAC)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UACCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
