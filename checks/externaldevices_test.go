package checks_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"reflect"
	"testing"
)

func TestExternalDevices(t *testing.T) {
	tests := []struct {
		name          string
		executorClass *utils.MockCommandExecutor
		want          checks.Check
	}{
		{
			name:          "No external devices connected",
			executorClass: &utils.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          checks.NewCheckResult("externaldevices", "", ""),
		},
		{
			name:          "External devices connected",
			executorClass: &utils.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want:          checks.NewCheckResult("externaldevices", "HD WebCam", "", "HD WebCam", ""),
		},
		{
			name:          "Error checking device",
			executorClass: &utils.MockCommandExecutor{Output: "", Err: errors.New("error checking device")},
			want:          checks.NewCheckErrorf("externaldevices", "error checking device Mouse", errors.New("error checking device")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checks.ExternalDevices(tt.executorClass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDeviceClass(t *testing.T) {
	tests := []struct {
		name          string
		deviceClass   string
		executorClass *utils.MockCommandExecutor
		want          []string
		wantErr       error
	}{
		{
			name:          "No devices of the specified class",
			deviceClass:   "Mouse",
			executorClass: &utils.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\n\r\n\r\n\r\n", Err: nil},
			want:          []string{""},
			wantErr:       nil,
		},
		{
			name:          "Devices of the specified class",
			deviceClass:   "Camera",
			executorClass: &utils.MockCommandExecutor{Output: "\r\nFriendlyName\r\n-\r\nHD WebCam\r\n\r\n\r\n\r\n", Err: nil},
			want:          []string{"HD WebCam", ""},
			wantErr:       nil,
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
