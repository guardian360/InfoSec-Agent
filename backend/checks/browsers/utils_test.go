package browsers_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Setup
	logger.SetupTests()

	// Run the tests
	code := m.Run()

	// Teardown

	// Exit with the code returned from the tests
	os.Exit(code)
}

// TestCloseFileNoError validates the CloseFile function's ability to close a file without errors.
//
// This test function creates a mock file, then calls the CloseFile function with this file as an argument.
// It asserts that no error is returned by the CloseFile function, indicating that the file was successfully closed.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileNoError(t *testing.T) {
	file := &mocking.FileMock{IsOpen: true, Err: nil}
	err := browsers.CloseFile(file)
	require.NoError(t, err)
}

// TestCloseFileWhenFileWasAlreadyClosed verifies the behavior of the CloseFile function when the file has already been closed.
//
// This test function asserts that the CloseFile function returns an error when it is called with a file that has already been closed.
// It is designed to ensure that the CloseFile function handles this edge case correctly, contributing to the robustness of the file handling process.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileWhenFileWasAlreadyClosed(t *testing.T) {
	file := &mocking.FileMock{IsOpen: true, Err: nil}
	err := file.Close()
	require.NoError(t, err)
	err = browsers.CloseFile(file)
	require.Error(t, err)
}

// TestCloseFileWhenFileIsNil verifies the behavior of the CloseFile function when the provided file is nil.
//
// This test function asserts that the CloseFile function returns an error when it is called with a nil file.
// It is designed to ensure that the CloseFile function handles this edge case correctly, contributing to the robustness of the file handling process.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCloseFileWhenFileIsNil(t *testing.T) {
	var file *mocking.FileMock
	err := browsers.CloseFile(file)
	require.Error(t, err)
}

// TestPhishingDomainsReturnsResults validates the behavior of the GetPhishingDomains function by ensuring it returns results.
//
// This test function calls the GetPhishingDomains function and asserts that the returned slice is not empty.
// It is designed to ensure that the GetPhishingDomains function correctly retrieves a list of phishing domains.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestPhishingDomainsReturnsResults(t *testing.T) {
	getter := browsers.RealPhishingDomainGetter{}
	requestCreator := browsers.RealRequestCreator{}
	domains, _ := getter.GetPhishingDomains(requestCreator)

	require.NotEmpty(t, domains)
}

type MockRequestCreator struct {
	NewRequestWithContextFunc func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
}

func (m MockRequestCreator) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return m.NewRequestWithContextFunc(ctx, method, url, body)
}

func TestGetPhishingDomains_NewRequestError(t *testing.T) {
	// Mock the http.NewRequestWithContext function to return an error
	mockRequestCreator := MockRequestCreator{
		NewRequestWithContextFunc: func(_ context.Context, _, _ string, _ io.Reader) (*http.Request, error) {
			return nil, errors.New("mock error")
		},
	}

	getter := browsers.RealPhishingDomainGetter{}
	requestCreator := mockRequestCreator
	domains, err := getter.GetPhishingDomains(requestCreator)

	// Assert there was an error and no domains were returned
	require.Error(t, err)
	require.Nil(t, domains)
}

type MockReadCloser struct {
	io.Reader
}

func (MockReadCloser) Close() error {
	return errors.New("mock error")
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetPhishingDomains_ResponseHandling(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		respStatusCode int
		respBody       string
		expectError    bool
		respError      error
	}{
		{
			name:           "Client Do method error",
			respStatusCode: http.StatusOK,
			respBody:       "",
			expectError:    true,
			respError:      errors.New("mock error"),
		},
		{
			name:           "HTTP request failed",
			respStatusCode: http.StatusNotFound,
			respBody:       "",
			expectError:    true,
		},
		{
			name:           "Error reading response body",
			respStatusCode: http.StatusOK,
			respBody:       "",
			expectError:    true,
		},
		{
			name:           "Response body is empty",
			respStatusCode: http.StatusOK,
			respBody:       "",
			expectError:    true,
		},
		{
			name:           "Successful request",
			respStatusCode: http.StatusOK,
			respBody:       "domain1\ndomain2\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the http.Client's Do function to return a custom response
			mockClient := MockClient{
				DoFunc: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.respStatusCode,
						Body:       MockReadCloser{strings.NewReader(tt.respBody)},
					}, tt.respError // Use the respError here
				},
			}

			getter := browsers.RealPhishingDomainGetter{Client: &mockClient}
			_, err := getter.GetPhishingDomains(browsers.RealRequestCreator{})

			// Assert the expected behavior
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type ErrorReadCloser struct{}

func (ErrorReadCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("mock read error")
}

func (ErrorReadCloser) Close() error {
	return nil
}

type MockDoer struct{}

func (MockDoer) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ErrorReadCloser{},
	}, nil
}

func TestGetPhishingDomains_ReadError(t *testing.T) {
	getter := browsers.RealPhishingDomainGetter{Client: MockDoer{}}
	_, err := getter.GetPhishingDomains(browsers.RealRequestCreator{})
	assert.Error(t, err)
}

// TestCopyFileSuccess validates the behavior of the CopyFile function when provided with a valid source and destination file.
//
// This test function creates a source file and a destination file, then calls the CopyFile function with these files as arguments.
// It asserts that no error is returned by the CopyFile function, indicating that the file was successfully copied from the source to the destination.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileSuccess(t *testing.T) {
	mockSource := &mocking.FileMock{IsOpen: true, Buffer: []byte{96, 96, 97, 97, 98, 98, 99, 99, 100, 100}, Bytes: 10, Err: nil}
	mockDestination := &mocking.FileMock{IsOpen: true, Buffer: []byte{96, 96, 97, 97, 98, 98, 99, 99, 100, 100}, Bytes: 10, Err: nil}
	getter := browsers.RealCopyFileGetter{}
	err := getter.CopyFile("", "", mockSource, mockDestination)
	require.NoError(t, err)
}

// TestCopyFileFailNonexistentSource validates the behavior of the CopyFile function when provided with a nonexistent source file.
//
// This test function calls the CopyFile function with a source file path that does not exist and a valid destination path.
// It asserts that an error is returned by the CopyFile function, indicating that the file could not be copied from the nonexistent source.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileFailNonexistentSource(t *testing.T) {
	mockSource := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: os.ErrNotExist}
	mockDestination := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	getter := browsers.RealCopyFileGetter{}
	err := getter.CopyFile("", "", mockSource, mockDestination)
	require.Error(t, err)
}

// TestCopyFileFailNonexistentDestination validates the behavior of the CopyFile function when provided with a nonexistent destination folder.
//
// This test function calls the CopyFile function with a valid source file and a destination path that does not exist.
// It asserts that an error is returned by the CopyFile function, indicating that the file could not be copied to the nonexistent destination.
//
// Parameters:
//   - t *testing.T: The testing framework used for assertions.
//
// No return values.
func TestCopyFileFailNonexistentDestination(t *testing.T) {
	mockSource := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: nil}
	mockDestination := &mocking.FileMock{IsOpen: true, Bytes: 10, Err: os.ErrNotExist}
	getter := browsers.RealCopyFileGetter{}
	err := getter.CopyFile("", "", mockSource, mockDestination)
	require.Error(t, err)
}

func TestGetDefaultDir(t *testing.T) {
	tests := []struct {
		name           string
		browserPath    string
		userHomeDir    string
		userHomeDirErr error
		expected       string
		expectErr      bool
	}{
		{
			name:           "Test with Chrome",
			browserPath:    browsers.ChromePath,
			userHomeDir:    "/mock/home/dir",
			userHomeDirErr: nil,
			expected:       "\\mock\\home\\dir\\AppData\\Local\\Google\\Chrome\\User Data\\Default",
			expectErr:      false,
		},
		{
			name:           "Test with Edge",
			browserPath:    browsers.EdgePath,
			userHomeDir:    "/mock/home/dir",
			userHomeDirErr: nil,
			expected:       "\\mock\\home\\dir\\AppData\\Local\\Microsoft\\Edge\\User Data\\Default",
			expectErr:      false,
		},
		{
			name:           "Test with error",
			browserPath:    browsers.ChromePath,
			userHomeDir:    "",
			userHomeDirErr: errors.New("mock error"),
			expected:       "",
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock UserHomeDirFunc
			browsers.UserHomeDirFunc = func() (string, error) {
				return tt.userHomeDir, tt.userHomeDirErr
			}

			getter := browsers.RealDefaultDirGetter{}
			got, err := getter.GetDefaultDir(tt.browserPath)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func createTestDatabase(t *testing.T, testData string) string {
	// Create a temporary file for the SQLite database
	tempFile, err := os.CreateTemp("", "testdb-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// Open the database connection
	db, err := sql.Open("sqlite", tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		t.Fatalf("Failed to ping database: %v", pingErr)
	}

	// Execute the test data SQL statements
	_, err = db.Exec(testData)
	if err != nil {
		t.Fatalf("Failed to execute test data SQL: %v", err)
	}

	return tempFile.Name()
}

func TestQueryCookieDatabase(t *testing.T) {
	// Setup test database
	testData := `
CREATE TABLE cookies (name TEXT, host_key TEXT);
INSERT INTO cookies (name, host_key) VALUES ('_ga', 'example.com');
INSERT INTO cookies (name, host_key) VALUES ('_utm', 'example.com');
INSERT INTO cookies (name, host_key) VALUES ('sessionid', 'example.com');
`
	dbPath := createTestDatabase(t, testData)
	dbPathNothing := createTestDatabase(t, "CREATE TABLE cookies (name TEXT, host_key TEXT);")
	defer os.Remove(dbPath)

	tests := []struct {
		name         string
		checkID      int
		browser      string
		databasePath string
		queryParams  []string
		tableName    string
		expected     checks.Check
	}{
		{
			name:         "Tracking cookies found",
			checkID:      1,
			browser:      "chrome",
			databasePath: dbPath,
			queryParams:  []string{"name", "host_key"},
			tableName:    "cookies",
			expected:     checks.NewCheckResult(1, 1, "_ga", "example.com", "_utm", "example.com"),
		},
		{
			name:         "No tracking cookies found",
			checkID:      2,
			browser:      "chrome",
			databasePath: dbPathNothing,
			queryParams:  []string{"name", "host_key"},
			tableName:    "cookies",
			expected:     checks.NewCheckResult(2, 0),
		},
	}
	getter := browsers.RealQueryCookieDatabaseGetter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getter.QueryCookieDatabase(tt.checkID, tt.browser, tt.databasePath, tt.queryParams, tt.tableName, browsers.RealCopyFileGetter{})
			assert.Equal(t, tt.expected, result)
		})
	}
}

type MockCopyFileGetter struct {
	CopyFileFunc func(src, dst string, mockSource mocking.File, mockDestination mocking.File) error
}

func (m MockCopyFileGetter) CopyFile(src, dst string, mockSource mocking.File, mockDest mocking.File) error {
	return m.CopyFileFunc(src, dst, mockSource, mockDest)
}

func TestQueryCookieDatabase_CopyError(t *testing.T) {
	mockGetter := MockCopyFileGetter{
		CopyFileFunc: func(_, _ string, _ mocking.File, _ mocking.File) error {
			return errors.New("mock error")
		},
	}
	getter := browsers.RealQueryCookieDatabaseGetter{}
	// Call the QueryCookieDatabase function
	result := getter.QueryCookieDatabase(1, "chrome", "/path/to/database", []string{"name", "host_key"}, "cookies", mockGetter)

	// Assert there was an error
	assert.Error(t, result.Error)
}

func TestQueryCookieDatabaseRowsError(t *testing.T) {
	mockGetter := MockCopyFileGetter{
		CopyFileFunc: func(_, _ string, _ mocking.File, _ mocking.File) error {
			return nil
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	checkID := 1
	browser := "test_browser"
	databasePath := "test_db_path"
	queryParams := []string{"name", "host"}
	tableName := "cookies"

	// Simulate error during query execution
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", strings.Join(queryParams, ", "), tableName)).
		WillReturnError(errors.New("query execution error"))

	getter := browsers.RealQueryCookieDatabaseGetter{}
	check := getter.QueryCookieDatabase(checkID, browser, databasePath, queryParams, tableName, mockGetter)
	require.Error(t, check.Error)

	// Simulate error after retrieving rows
	rows := sqlmock.NewRows([]string{"name", "host"}).AddRow("cookie_name", "cookie_host")
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", strings.Join(queryParams, ", "), tableName)).
		WillReturnRows(rows)
	//

	check = getter.QueryCookieDatabase(checkID, browser, databasePath, queryParams, tableName, mockGetter)
	require.Error(t, check.Error)
}
