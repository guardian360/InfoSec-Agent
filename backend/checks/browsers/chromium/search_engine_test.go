package chromium_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

func TestSearchEngineChromium(t *testing.T) {
	tests := []struct {
		name     string
		browser  string
		mockBool bool
		mockFile *mocking.FileMock
		want     checks.Check
	}{
		{
			name:     "Test with Chrome",
			browser:  "Chrome",
			mockBool: true,
			mockFile: &mocking.FileMock{
				FileName: "\\valid\\directory\\Preferences",
				IsOpen:   true,
				Buffer:   []byte("{\"default_search_provider_data\":{\"template_url_data\":{\"keyword\":\"google.com\"}}}"),
				Bytes:    3,
				FileInfo: &mocking.FileInfoMock{},
			},
			want: checks.NewCheckResult(checks.SearchChromiumID, 0, "google.com"),
		},
		{
			name:     "Test with Edge",
			browser:  "Edge",
			mockBool: true,
			mockFile: &mocking.FileMock{
				FileName: "\\valid\\directory\\Preferences",
				IsOpen:   true,
				Buffer:   []byte("{\"default_search_provider_data\":{\"template_url_data\":{\"keyword\":\"google.com\"}}}"),
				Bytes:    3,
				FileInfo: &mocking.FileInfoMock{},
			},
			want: checks.NewCheckResult(checks.SearchEdgeID, 0, "google.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := browsers.RealDefaultDirGetter{}
			got := chromium.SearchEngineChromium(tt.browser, tt.mockBool, tt.mockFile, getter)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetBrowserPathAndIDSearch(t *testing.T) {
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
			wantID:   checks.SearchChromiumID,
		},
		{
			name:     "Test with Edge",
			browser:  "Edge",
			wantPath: "Microsoft/Edge",
			wantID:   checks.SearchEdgeID,
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
			gotPath, gotID := chromium.GetBrowserPathAndIDSearch(tt.browser)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}

// Create a mock implementation of the interface
type MockPreferencesDirGetter struct{}

func (m MockPreferencesDirGetter) GetDefaultDir(_ string) (string, error) {
	return "", errors.New("mock error")
}

func TestSearchEngineChromium_GetPreferencesDirError(t *testing.T) {
	// Use the mock implementation in your test
	getter := MockPreferencesDirGetter{}

	mockFile := &mocking.FileMock{
		FileName: "\\valid\\directory\\Preferences",
		IsOpen:   true,
		Buffer:   []byte("{\"default_search_provider_data\":{\"template_url_data\":{\"keyword\":\"google.com\"}}}"),
		Bytes:    3,
		FileInfo: &mocking.FileInfoMock{},
	}
	want := checks.NewCheckErrorf(checks.SearchChromiumID, "Error: ", errors.New("mock error"))

	// Call SearchEngineChromium with the mock implementation
	result := chromium.SearchEngineChromium("Chrome", true, mockFile, getter)

	// Assert that the result is a CheckError
	require.Equal(t, want, result)
}

func TestSearchEngineChromium_ParsePreferencesFileError(t *testing.T) {
	// Use the mock implementation in your test
	getter := browsers.RealDefaultDirGetter{}

	mockFile := &mocking.FileMock{
		FileName: "\\valid\\directory\\Preferences",
		IsOpen:   true,
		Buffer:   []byte("invalid json"), // This will cause ParsePreferencesFile to return an error
		Bytes:    11,
		FileInfo: &mocking.FileInfoMock{},
	}
	want := checks.NewCheckErrorf(checks.SearchChromiumID, "Error: ", errors.New("invalid character 'i' looking for beginning of value"))

	// Call SearchEngineChromium with the mock implementation
	result := chromium.SearchEngineChromium("Chrome", true, mockFile, getter)

	// Assert that the result is a CheckError
	require.Equal(t, want.ErrorMSG, result.ErrorMSG)
}

func TestParsePreferencesFile(t *testing.T) {
	// Create a temporary file with some JSON content
	file := &mocking.FileMock{
		FileName: "\\valid\\directory\\search.json.mozlz4",
		IsOpen:   true,
		Buffer:   []byte("{\"key1\":\"value1\",\"key2\":\"value2\"}"),
		Bytes:    3,
		FileInfo: &mocking.FileInfoMock{},
	}

	// Call ParsePreferencesFile
	result, err := chromium.ParsePreferencesFile(file)
	if err != nil {
		t.Fatalf("ParsePreferencesFile() error = %v", err)
	}

	// Check the result
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParsePreferencesFile() = %v, want %v", result, expected)
	}
}

func TestParsePreferencesFileErrorRead(t *testing.T) {
	errRead := errors.New("read error")
	// Create a temporary file with some JSON content
	file := &mocking.FileMock{
		FileName: "\\valid\\directory\\search.json.mozlz4",
		IsOpen:   true,
		Buffer:   []byte{},
		Bytes:    2,
		FileInfo: &mocking.FileInfoMock{},
		Err:      errRead,
	}

	// Call ParsePreferencesFile
	result, err := chromium.ParsePreferencesFile(file)

	// Check the result
	require.Nil(t, result)
	require.Error(t, err)
}

func TestParsePreferencesFileErrorParse(t *testing.T) {
	// Create a temporary file with some JSON content
	file := &mocking.FileMock{
		FileName: "\\valid\\directory\\search.json.mozlz4",
		IsOpen:   true,
		Buffer:   []byte{2, 3},
		Bytes:    2,
		FileInfo: &mocking.FileInfoMock{},
	}

	// Call ParsePreferencesFile
	result, err := chromium.ParsePreferencesFile(file)

	// Check the result
	require.Nil(t, result)
	require.Error(t, err)
}

func TestGetDefaultSearchEngine(t *testing.T) {
	// Test cases
	tests := []struct {
		name  string
		dev   map[string]interface{}
		defSE string
		want  string
	}{
		{
			name: "Google",
			dev: map[string]interface{}{
				"default_search_provider_data": map[string]interface{}{
					"template_url_data": map[string]interface{}{
						"keyword": "google.com ",
					},
				},
			},
			defSE: "google.com",
			want:  "google.com",
		},
		{
			name: "Bing",
			dev: map[string]interface{}{
				"default_search_provider_data": map[string]interface{}{
					"template_url_data": map[string]interface{}{
						"keyword": "bing.com ",
					},
				},
			},
			defSE: "google.com",
			want:  "bing.com",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chromium.GetDefaultSearchEngine(tt.dev, tt.defSE)
			if got != tt.want {
				t.Errorf("GetDefaultSearchEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultSearchEngineError(t *testing.T) {
	// Test cases
	tests := []struct {
		name  string
		dev   map[string]interface{}
		defSE string
		want  string
	}{
		{
			name: "Emptykeyword",
			dev: map[string]interface{}{
				"default_search_provider_data": map[string]interface{}{
					"template_url_data": map[string]interface{}{
						"keyword": "",
					},
				},
			},
			defSE: "google.com",
			want:  "google.com",
		},
		{
			name:  "NoDefaultSearchProviderData",
			dev:   map[string]interface{}{},
			defSE: "google.com",
			want:  "google.com",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chromium.GetDefaultSearchEngine(tt.dev, tt.defSE)
			if got != tt.want {
				t.Errorf("GetDefaultSearchEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}
