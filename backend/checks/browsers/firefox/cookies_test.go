package firefox_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.SetupTests()
	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

// MockProfileFinder is a mock implementation of the FirefoxProfileFinder interface
type MockProfileFinder struct {
	FirefoxFolderFunc func() ([]string, error)
}

func (m *MockProfileFinder) FirefoxFolder() ([]string, error) {
	return m.FirefoxFolderFunc()
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

func TestCookiesFirefox(t *testing.T) {
	// Mock the FirefoxProfileFinder to return a specific directory and no error
	// Mock the FirefoxProfileFinder to return a specific directory and no error
	profilefinder := &browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"\\valid\\directory"}, nil
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
	result := firefox.CookiesFirefox(profilefinder, mockCopyGetter, mockQueryGetter)

	// Assert there was no error and the result is as expected
	assert.NoError(t, result.Error)
}

func TestCookiesFirefox_Error(t *testing.T) {
	// Mock the FirefoxProfileFinder to return a specific directory and no error
	Profilefinder = browsers.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"\\invalid\\directory"}, errors.New("No firefox directory found")
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
	result := firefox.CookiesFirefox(Profilefinder, mockCopyGetter, mockQueryGetter)

	// Assert there was an error
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "No firefox directory found")
}
