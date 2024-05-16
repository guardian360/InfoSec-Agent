package chromium_test

import (
	"errors"
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	logger.SetupTests()
	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

type MockDefaultDirGetter struct {
	GetDefaultDirFunc func(browserPath string) (string, error)
}

func (m MockDefaultDirGetter) GetDefaultDir(browserPath string) (string, error) {
	return m.GetDefaultDirFunc(browserPath)
}

type MockExtensionIDGetter struct {
	GetExtensionIDsFunc func(extensionsDir string) ([]string, error)
}

func (m MockExtensionIDGetter) GetExtensionIDs(extensionsDir string) ([]string, error) {
	return m.GetExtensionIDsFunc(extensionsDir)
}

func TestExtensionsChromium(t *testing.T) {
	tests := []struct {
		name                      string
		browser                   string
		mockGetDefaultDirFunc     func(browserPath string) (string, error)
		mockGetExtensionIDsFunc   func(extensionsDir string) ([]string, error)
		mockGetExtensionNamesFunc func(extensionID string, url string, browser string) (string, error)
		expectedResult            checks.Check
	}{
		{
			name:    "Test with Chrome and adblocker installed",
			browser: "Chrome",
			mockGetDefaultDirFunc: func(_ string) (string, error) {
				return "/mock/path/to/chrome", nil
			},
			mockGetExtensionIDsFunc: func(_ string) ([]string, error) {
				return []string{"https://chromewebstore.google.com/detail/validExtensionID/fdsfdsfdsfs", "extension2"}, nil
			},
			mockGetExtensionNamesFunc: func(_ string, _ string, _ string) (string, error) {
				return "https://chromewebstore.google.com/detail/adblocker/fdsfdsfdsfs", nil
			},
			expectedResult: checks.NewCheckResult(checks.ExtensionChromiumID, 0),
		},
		{
			name:    "Test with Chrome and no adblocker installed",
			browser: "Chrome",
			mockGetDefaultDirFunc: func(_ string) (string, error) {
				return "/mock/path/to/chrome", nil
			},
			mockGetExtensionIDsFunc: func(_ string) ([]string, error) {
				return []string{"https://chromewebstore.google.com/detail/validExtensionID/fdsfdsfdsfs", "extension2"}, nil
			},
			mockGetExtensionNamesFunc: func(_ string, _ string, _ string) (string, error) {
				return "https://chromewebstore.google.com/detail/privacyGuard/fdsfdsfdsfs", nil
			},
			expectedResult: checks.NewCheckResult(checks.ExtensionChromiumID, 1),
		},
		{
			name:    "Test with error in GetDefaultDir",
			browser: "Chrome",
			mockGetDefaultDirFunc: func(_ string) (string, error) {
				return "", errors.New("mock error")
			},
			mockGetExtensionIDsFunc: func(_ string) ([]string, error) {
				return nil, nil
			},
			mockGetExtensionNamesFunc: func(_ string, _ string, _ string) (string, error) {
				return "", nil
			},
			expectedResult: checks.NewCheckErrorf(checks.ExtensionChromiumID, "Error: ", errors.New("mock error")),
		},
		{
			name:    "Test with error in GetExtensionIDs",
			browser: "Chrome",
			mockGetDefaultDirFunc: func(_ string) (string, error) {
				return "/mock/path/to/chrome", nil
			},
			mockGetExtensionIDsFunc: func(_ string) ([]string, error) {
				return nil, errors.New("mock error")
			},
			mockGetExtensionNamesFunc: func(_ string, _ string, _ string) (string, error) {
				return "", nil
			},
			expectedResult: checks.NewCheckErrorf(checks.ExtensionChromiumID, "Error: ", errors.New("mock error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := MockDefaultDirGetter{
				GetDefaultDirFunc: tt.mockGetDefaultDirFunc,
			}

			mockGetter := MockExtensionIDGetter{
				GetExtensionIDsFunc: tt.mockGetExtensionIDsFunc,
			}

			mockNameGetter := MockExtensionNameGetter{
				GetExtensionNameChromiumFunc: tt.mockGetExtensionNamesFunc,
			}

			result := chromium.ExtensionsChromium(tt.browser, getter, mockGetter, mockNameGetter)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetBrowserPathAndIDExtension(t *testing.T) {
	tests := []struct {
		name     string
		browser  string
		wantPath string
		wantID   int
	}{
		{
			name:     "Test with Chrome",
			browser:  "Chrome",
			wantPath: "Google/Chrome",
			wantID:   checks.ExtensionChromiumID,
		},
		{
			name:     "Test with Edge",
			browser:  "Edge",
			wantPath: "Microsoft/Edge",
			wantID:   checks.ExtensionEdgeID,
		},
		{
			name:     "Test with unknown browser",
			browser:  "Unknown",
			wantPath: "",
			wantID:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotID := chromium.GetBrowserPathAndIDExtension(tt.browser)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}

func TestGetExtensionIDs(t *testing.T) {
	// Create a temporary directory
	tempDir, mKerr := os.MkdirTemp("", "extensions")
	if mKerr != nil {
		t.Fatalf("Failed to create temp dir: %v", mKerr)
	}
	defer os.RemoveAll(tempDir) // clean up

	// Create dummy extension directories
	for i := range []int{0, 1, 2}{
		err := os.Mkdir(filepath.Join(tempDir, fmt.Sprintf("extension%d", i)), 0755)
		if err != nil {
			t.Fatalf("Failed to create dummy extension dir: %v", err)
		}
	}

	getter := chromium.RealExtensionIDGetter{}
	got, err := getter.GetExtensionIDs(tempDir)
	if err != nil {
		t.Fatalf("GetExtensionIDs() error = %v, wantErr nil", err)
	}

	// Check the extension IDs
	want := []string{"extension0", "extension1", "extension2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetExtensionIDs() = %v, want %v", got, want)
	}
}

// Define an interface with the same methods as the os package that you want to mock
type OSInterface interface {
	ReadDir(dirname string) ([]os.DirEntry, error)
	// Add other methods as needed
}

// Create a real implementation that calls the os package
type RealOS struct{}

func (RealOS) ReadDir(dirname string) ([]os.DirEntry, error) {
	return os.ReadDir(dirname)
}

// Create a variable of the OSInterface type and use it in place of the os package
var osFS OSInterface = RealOS{}

// Now in your test, you can replace osFS with a mock implementation
func TestGetExtensionIDs_Error(t *testing.T) {
	// Create a mock implementation
	mockOS := &mocking.FileMock{
		FileName: "unreadable-dir",
		IsOpen:   true,
		Buffer:   []byte("dummyData"),
		Err:      os.ErrPermission, // Simulate a permission error
	}

	// Replace osFS with the mock implementation
	oldOS := osFS
	osFS = mockOS
	defer func() { osFS = oldOS }()

	getter := chromium.RealExtensionIDGetter{}
	_, err := getter.GetExtensionIDs(mockOS.FileName)
	assert.Error(t, err)
}

// Mock the GetExtensionNameChromium function
var MockGetExtensionNameChromium = func(_ string, _ string, _ string) (string, error) {
	// Return a fixed extension name and no error
	return "adblock", nil
}

type MockExtensionNameGetter struct {
	GetExtensionNameChromiumFunc func(_ string, _ string, _ string) (string, error)
}

// Define GetExtensionNameChromium as a method of MockExtensionNameGetter
func (m MockExtensionNameGetter) GetExtensionNameChromium(extensionID string, url string, browser string) (string, error) {
	// Call the function field here
	return m.GetExtensionNameChromiumFunc(extensionID, url, browser)
}

func TestGetExtensionNames(t *testing.T) {
	tests := []struct {
		name         string
		browser      string
		extensionIDs []string
		mockFunc     func(extensionID string, url string, browser string) (string, error)
		want         []string
		wantErr      bool
	}{
		{
			name:         "Chrome browser with valid extension ID",
			browser:      "Chrome",
			extensionIDs: []string{"validExtensionID"},
			mockFunc: func(_ string, _ string, _ string) (string, error) {
				return "https://chromewebstore.google.com/detail/validExtensionID/fdsfdsfdsfs", nil
			},
			want:    []string{"validExtensionID"},
			wantErr: false,
		},
		{
			name:         "Edge browser with valid extension ID",
			browser:      "Edge",
			extensionIDs: []string{"validExtensionID"},
			mockFunc: func(_ string, _ string, _ string) (string, error) {
				return "validExtensionName", nil
			},
			want:    []string{"validExtensionName"},
			wantErr: false,
		},
		{
			name:         "Error getting extension name",
			browser:      "Chrome",
			extensionIDs: []string{"invalidExtensionID"},
			mockFunc: func(_ string, _ string, _ string) (string, error) {
				return "", errors.New("mocked error")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "Error getting extension name Edge",
			browser:      "Edge",
			extensionIDs: []string{"invalidExtensionID"},
			mockFunc: func(_ string, _ string, _ string) (string, error) {
				return "", errors.New("mocked error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := MockExtensionNameGetter{
				GetExtensionNameChromiumFunc: tt.mockFunc,
			}

			got := chromium.GetExtensionNames(getter, tt.extensionIDs, tt.browser)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExtensionNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetExtensionNameChromium(t *testing.T) {
	tests := []struct {
		name        string
		extensionID string
		url         string
		browser     string
		wantErr     bool
	}{
		{
			name:        "uBlock Origin",
			extensionID: "cjpalhdlnbpafiamejdnhcphjbkeiagm",
			url:         "https://chromewebstore.google.com/detail/%s",
			browser:     "Chrome",
			wantErr:     false,
		},
		{
			name:        "Valid Edge extension",
			extensionID: "ndcileolkflehcjpmjnfbnaibdcgglog",
			url:         "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s",
			browser:     "Edge",
			wantErr:     false,
		},
		{
			name:        "Invalid extension ID",
			extensionID: "invalid-extension-id",
			url:         "https://chrome.google.com/webstore/detail/%s",
			browser:     "Chrome",
			wantErr:     true,
		},
		{
			name:        "Unknown browser",
			extensionID: "valid-chrome-extension-id",
			url:         "https://chrome.google.com/webstore/detail/%s",
			browser:     "Unknown",
			wantErr:     true,
		},
	}

	getter := chromium.ChromeExtensionNameGetter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getter.GetExtensionNameChromium(tt.extensionID, tt.url, tt.browser)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExtensionNameChromium() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetExtensionNameChromium_ErrorCreatingRequest(t *testing.T) {
	getter := chromium.ChromeExtensionNameGetter{}
	_, err := getter.GetExtensionNameChromium("invalid-extension-id", ":", "Chrome")
	if err == nil {
		t.Errorf("Expected error due to invalid URL, got nil")
	}
}

func TestGetExtensionNameChromium_ErrorSendingRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://chromewebstore.google.com/detail/invalid-extension-id",
		func(_ *http.Request) (*http.Response, error) {
			return nil, errors.New("mocked error")
		},
	)

	getter := chromium.ChromeExtensionNameGetter{}

	_, err := getter.GetExtensionNameChromium("invalid-extension-id", "https://chromewebstore.google.com/detail/%s", "Chrome")
	if err == nil {
		t.Errorf("Expected error due to mocked error in client.Do, got nil")
	}
}

func TestGetExtensionNameChromium_NonChromeWebStoreURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/valid-extension-id",
		func(_ *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Request = &http.Request{URL: &url.URL{Scheme: "https", Host: "chromewebstore.google.com", Path: "/addons/getproductdetailsbycrxid/valid-extension-id"}}
			return resp, nil
		},
	)

	getter := chromium.ChromeExtensionNameGetter{}

	extensionName, err := getter.GetExtensionNameChromium("valid-extension-id", "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s", "Edge")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if extensionName != "https://chromewebstore.google.com/addons/getproductdetailsbycrxid/valid-extension-id" {
		t.Errorf("Expected extension name to be 'https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/valid-extension-id', got '%s'", extensionName)
	}
}

func TestGetExtensionNameChromium_ErrorDecodingJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/valid-extension-id",
		func(_ *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "invalid json")
			resp.Request = &http.Request{URL: &url.URL{Scheme: "https", Host: "microsoftedge.microsoft.com", Path: "/addons/getproductdetailsbycrxid/valid-extension-id"}}
			return resp, nil
		},
	)

	getter := chromium.ChromeExtensionNameGetter{}

	_, err := getter.GetExtensionNameChromium("valid-extension-id", "https://microsoftedge.microsoft.com/addons/getproductdetailsbycrxid/%s", "Edge")
	if err == nil {
		t.Errorf("Expected error due to invalid JSON, got nil")
	}
}

func TestAdblockerInstalled(t *testing.T) {
	tests := []struct {
		name           string
		extensionNames []string
		want           bool
	}{
		{
			name:           "No extensions",
			extensionNames: []string{},
			want:           false,
		},
		{
			name:           "Non-adblocker extensions",
			extensionNames: []string{"Extension1", "Extension2"},
			want:           false,
		},
		{
			name:           "Adblocker extension",
			extensionNames: []string{"Extension1", "adblock"},
			want:           true,
		},
		{
			name:           "Adblocker extension with different case",
			extensionNames: []string{"Extension1", "AdBlock"},
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chromium.AdblockerInstalled(tt.extensionNames); got != tt.want {
				t.Errorf("adblockerInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}
