package chromium_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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

type MockCopyFileGetter struct {
	CopyFileFunc func(string, string, mocking.File, mocking.File) error
}

func (m *MockCopyFileGetter) CopyFile(src, dst string, mocksrc mocking.File, mockdst mocking.File) error {
	return m.CopyFileFunc(src, dst, mocksrc, mockdst)
}

type MockQueryCookieDatabaseGetter struct {
	QueryCookieDatabaseFunc func(int, string, string, []string, string, browsers.CopyFileGetter) checks.Check
}

func (m *MockQueryCookieDatabaseGetter) QueryCookieDatabase(id int, browser, dbPath string, columns []string, table string, getter browsers.CopyFileGetter) checks.Check {
	return m.QueryCookieDatabaseFunc(id, browser, dbPath, columns, table, getter)
}

func TestCookiesChromium(t *testing.T) {
	mockDefaultDirGetter := &MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "valid\\directory", nil
		},
	}

	mockCopyGetter := &MockCopyFileGetter{
		CopyFileFunc: func(_ string, _ string, _ mocking.File, _ mocking.File) error {
			return nil
		},
	}

	mockQueryGetter := &MockQueryCookieDatabaseGetter{
		QueryCookieDatabaseFunc: func(_ int, _ string, _ string, _ []string, _ string, _ browsers.CopyFileGetter) checks.Check {
			return checks.Check{
				Result: []string{"cookie1", "cookie2"},
			}
		},
	}

	// Call the CookiesFirefox function
	result := chromium.CookiesChromium("Chrome", mockDefaultDirGetter, mockCopyGetter, mockQueryGetter)

	// Assert there was no error and the result is as expected
	assert.NoError(t, result.Error)
}

func TestCookiesChromium_Error(t *testing.T) {
	mockDefaultDirGetter := &MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "valid\\directory", errors.New("No chromium directory found")
		},
	}

	mockCopyGetter := &MockCopyFileGetter{
		CopyFileFunc: func(_ string, _ string, _ mocking.File, _ mocking.File) error {
			return nil
		},
	}

	mockQueryGetter := &MockQueryCookieDatabaseGetter{
		QueryCookieDatabaseFunc: func(_ int, _ string, _ string, _ []string, _ string, _ browsers.CopyFileGetter) checks.Check {
			return checks.Check{
				Result: []string{"cookie1", "cookie2"},
			}
		},
	}

	// Call the CookiesFirefox function
	result := chromium.CookiesChromium("Chrome", mockDefaultDirGetter, mockCopyGetter, mockQueryGetter)

	// Assert there was an error
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "No chromium directory found")
}

func TestGetBrowserPathAndIDCookies(t *testing.T) {
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
			wantID:   checks.CookiesChromiumID,
		},
		{
			name:     "Test with Edge",
			browser:  "Edge",
			wantPath: "Microsoft/Edge",
			wantID:   checks.CookiesEdgeID,
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
			gotPath, gotID := chromium.GetBrowserPathAndIDCookie(tt.browser)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}
