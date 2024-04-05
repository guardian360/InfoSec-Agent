package checks_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TestExternalDevices tests the ExternalDevices function with (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestExternalDevices(t *testing.T) {
	tests := []struct {
		name          string
		executorClass *commandmock.MockCommandExecutor
		want          checks.Check
	}{
		{
			name:          "No external devices connected",
			executorClass: &commandmock.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          checks.NewCheckResult(checks.ExternalDevicesID, 0, "", ""),
		},
		{
			name: "External devices connected",
			executorClass: &commandmock.MockCommandExecutor{
				Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want: checks.NewCheckResult(checks.ExternalDevicesID, 0, "HD WebCam", "", "HD WebCam", ""),
		},
		{
			name:          "Error checking device",
			executorClass: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("error checking device")},
			want: checks.NewCheckErrorf(checks.ExternalDevicesID, "error checking device Mouse",
				errors.New("error checking device")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.ExternalDevices(tt.executorClass)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestCheckDeviceClass tests the CheckDeviceClass with (in)valid inputs
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestCheckDeviceClass(t *testing.T) {
	tests := []struct {
		name          string
		deviceClass   string
		executorClass *commandmock.MockCommandExecutor
		want          []string
		wantErr       error
	}{
		{
			name:          "No devices of the specified class",
			deviceClass:   "Mouse",
			executorClass: &commandmock.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          []string{""},
			wantErr:       nil,
		},
		{
			name:        "Devices of the specified class",
			deviceClass: "Camera",
			executorClass: &commandmock.MockCommandExecutor{
				Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want:    []string{"HD WebCam", ""},
			wantErr: nil,
		},
		{
			name:          "Error checking device",
			deviceClass:   "Camera",
			executorClass: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("error checking device")},
			want:          nil,
			wantErr:       errors.New("error checking device"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := checks.CheckDeviceClass(tt.deviceClass, tt.executorClass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCommandOutput tests that the output of the command run in externaldevices.go is as expected
//
// Parameters: t (testing.T) - the testing framework
//
// Returns: _
func TestCommandOutput(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		arguments string
		expected  string
	}{
		{
			name:      "Get-PnpDevice output",
			command:   "powershell",
			arguments: "Get-PnpDevice | Where-Object -Property Status -eq 'OK' | Select-Object FriendlyName",
			expected:  "FriendlyName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &commandmock.RealCommandExecutor{}
			output, _ := executor.Execute(tt.command, tt.arguments)
			outputList := strings.Split(string(output), "\r\n")
			if res := strings.ReplaceAll(outputList[1], " ", ""); res != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, res)
			}
		})
	}
}
