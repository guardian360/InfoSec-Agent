package checks_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// TestExternalDevices is a test function that validates the behavior of the ExternalDevices function.
// It executes a series of test cases, each with different inputs, to ensure that the function behaves as expected in various scenarios.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If a test case fails, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the ExternalDevices function correctly identifies external devices connected to the system and handles errors as expected.
func TestExternalDevices(t *testing.T) {
	tests := []struct {
		name          string
		executorClass *commandmock.MockCommandExecutor
		want          checks.Check
	}{
		{
			name:          "No external devices connected",
			executorClass: &commandmock.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          checks.NewCheckResult("externaldevices", "", ""),
		},
		{
			name: "External devices connected",
			executorClass: &commandmock.MockCommandExecutor{
				Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want: checks.NewCheckResult("externaldevices", "HD WebCam", "", "HD WebCam", ""),
		},
		{
			name:          "Error checking device",
			executorClass: &commandmock.MockCommandExecutor{Output: "", Err: errors.New("error checking device")},
			want: checks.NewCheckErrorf("externaldevices", "error checking device Mouse",
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

// TestCheckDeviceClass is a test function that validates the behavior of the CheckDeviceClass function.
// It executes a series of test cases, each with different inputs, to ensure that the function behaves as expected in various scenarios.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If a test case fails, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the CheckDeviceClass function correctly identifies devices of a specific class connected to the system and handles errors as expected.
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

// TestCommandOutput is a test function that validates the output of the system command executed in the ExternalDevices function in the 'externaldevices.go' file.
//
// Parameters:
//   - t (*testing.T): An instance of the testing framework provided by the "testing" package. This is used to report test failures and log output.
//
// Returns: None. If the output of the command does not match the expected output, the function calls methods on the *testing.T parameter to report the failure.
//
// This function is part of the test suite for the "checks" package. It is used to verify that the system command executed in the ExternalDevices function produces the expected output. This helps ensure that the function is correctly identifying external devices connected to the system.
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
