package chromium_test

import (
	"database/sql"
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/chromium"
	"github.com/stretchr/testify/assert"
)

type MockCopyDBGetter struct {
	CopyDatabaseFunc func(source string, browser string, getter browsers.CopyFileGetter) (string, error)
}

func (m MockCopyDBGetter) CopyDatabase(source string, browser string, getter browsers.CopyFileGetter) (string, error) {
	return m.CopyDatabaseFunc(source, browser, getter)
}

type MockQueryDatabaseGetter struct {
	QueryDatabaseFunc func(db *sql.DB) (*sql.Rows, error)
}

func (m MockQueryDatabaseGetter) QueryDatabase(db *sql.DB) (*sql.Rows, error) {
	return m.QueryDatabaseFunc(db)
}

type MockProcessQueryResultsGetter struct {
	ProcessQueryResultsFunc func(rows *sql.Rows, getter browsers.PhishingDomainGetter) ([]string, error)
}

func (m MockProcessQueryResultsGetter) ProcessQueryResults(rows *sql.Rows, getter browsers.PhishingDomainGetter) ([]string, error) {
	return m.ProcessQueryResultsFunc(rows, getter)
}

type MockPhishingDomainsGetter struct {
	GetPhishingDomainsFunc func() ([]string, error)
}

func (m MockPhishingDomainsGetter) GetPhishingDomains(_ browsers.RequestCreator) ([]string, error) {
	return m.GetPhishingDomainsFunc()
}

func TestHistoryChromium_Success(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "/path/to/browser", nil
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "/path/to/temp", nil
		},
	}

	mockQueryDBGetter := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			db, mock, err := sqlmock.New()
			if err != nil {
				return nil, err
			}
			defer func() {
				_ = db.Close()
			}()

			rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"})
			rows.AddRow("http://example.com", "Example Site", 1, 1000000000)

			mock.ExpectQuery("^SELECT (.+) FROM websites$").WillReturnRows(rows)

			realRows, err := db.Query("SELECT * FROM websites")
			if err != nil {
				return nil, err
			}

			return realRows, nil
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockQueryDBGetter, mockGetterQR, mockPhishingGetter)

	assert.Equal(t, checks.NewCheckResult(checks.HistoryChromiumID, 0, "phishing.com"), result)
}

func TestHistoryChromium_Success_NoPhishing(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "/path/to/browser", nil
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "/path/to/temp", nil
		},
	}

	mockQueryDBGetter := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			db, mock, _ := sqlmock.New()
			defer func() {
			}()

			rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"})
			rows.AddRow("http://example.com", "Example Site", 1, 1000000000)

			mock.ExpectQuery("^SELECT (.+) FROM websites$").WillReturnRows(rows)

			realRows, err := db.Query("SELECT * FROM websites")
			if err != nil {
				return nil, err
			}

			return realRows, nil
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return []string{}, nil
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockQueryDBGetter, mockGetterQR, mockPhishingGetter)

	assert.Equal(t, checks.NewCheckResult(checks.HistoryChromiumID, 1), result)
}

func TestHistoryChromium_Error_GetDefaultDir(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "", errors.New("mock error")
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "/path/to/temp", nil
		},
	}

	mockQueryDBGetter := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			return &sql.Rows{}, nil
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return []string{}, nil
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockQueryDBGetter, mockGetterQR, mockPhishingGetter)

	assert.Equal(t, checks.NewCheckErrorf(checks.HistoryChromiumID, "error getting preferences directory", errors.New("mock error")), result)
}

func TestHistoryChromium_Error_CopyDatabase(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "/path/to/browser", nil
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "", errors.New("mock error")
		},
	}

	mockQueryDBGetter := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			return &sql.Rows{}, nil
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return []string{}, nil
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockQueryDBGetter, mockGetterQR, mockPhishingGetter)

	assert.Error(t, result.Error)
}

func TestHistoryChromium_Error_QueryDb(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "/path/to/browser", nil
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "/invalid/path", nil
		},
	}

	mockGetterQDB := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			return nil, errors.New("mock error")
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return []string{}, nil
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockGetterQDB, mockGetterQR, mockPhishingGetter)

	require.Error(t, result.Error)
	// require.Contains(t, result.Error.Error(), "unable to open database file")
}

func TestHistoryChromium_Error_ProcessQueryResults(t *testing.T) {
	mockGetter := MockDefaultDirGetter{
		GetDefaultDirFunc: func(_ string) (string, error) {
			return "/path/to/browser", nil
		},
	}

	mockCopyGetter := MockCopyDBGetter{
		CopyDatabaseFunc: func(_ string, _ string, _ browsers.CopyFileGetter) (string, error) {
			return "/path/to/temp", nil
		},
	}

	mockQueryDBGetter := MockQueryDatabaseGetter{
		QueryDatabaseFunc: func(_ *sql.DB) (*sql.Rows, error) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"})
			mock.ExpectQuery("^SELECT (.+) FROM websites$").WillReturnRows(rows)

			realRows, err := db.Query("SELECT * FROM websites")
			if err != nil {
				t.Fatalf("an error '%s' was not expected while querying", err)
			}

			return realRows, nil
		},
	}

	mockGetterQR := MockProcessQueryResultsGetter{
		ProcessQueryResultsFunc: func(_ *sql.Rows, _ browsers.PhishingDomainGetter) ([]string, error) {
			return nil, errors.New("mock error")
		},
	}

	mockPhishingGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	result := chromium.HistoryChromium("Chrome", mockGetter, mockCopyGetter, mockQueryDBGetter, mockGetterQR, mockPhishingGetter)

	require.Error(t, result.Error)
}

func TestGetBrowserPathAndIDHistory(t *testing.T) {
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
			wantID:   checks.HistoryChromiumID,
		},
		{
			name:     "Test with Edge",
			browser:  "Edge",
			wantPath: "Microsoft/Edge",
			wantID:   checks.HistoryEdgeID,
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
			gotPath, gotID := chromium.GetBrowserPathAndIDHistory(tt.browser)
			require.Equal(t, tt.wantPath, gotPath)
			require.Equal(t, tt.wantID, gotID)
		})
	}
}

func TestCopyDatabase_Success(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "source")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating a temporary file", err)
	}
	defer os.Remove(tempFile.Name())

	// Call the function
	getter := chromium.RealCopyDBGetter{}
	_, err = getter.CopyDatabase(tempFile.Name(), "", browsers.RealCopyFileGetter{})

	// Assert there was no error
	require.NoError(t, err)
}

func TestCopyDatabase_Error(t *testing.T) {
	// Define an invalid file path
	invalidFilePath := "/invalid/file/path"

	// Call the function
	getter := chromium.RealCopyDBGetter{}
	_, err := getter.CopyDatabase(invalidFilePath, "", browsers.RealCopyFileGetter{})

	// Assert there was an error
	require.Error(t, err)
}

func TestQueryDatabase(t *testing.T) {
	// Mock the sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Calculate the time one week ago
	oneWeekAgo := (time.Now().AddDate(369, 0, -7).UnixMicro() / 1000000000) * 1000000000

	// Mock the rows
	rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"})
	mock.ExpectQuery("^SELECT (.+) FROM urls WHERE last_visit_time > \\? ORDER BY last_visit_time DESC$").
		WithArgs(oneWeekAgo).
		WillReturnRows(rows)

	// Call the function
	getter := chromium.RealQueryDatabaseGetter{}
	rows2, err := getter.QueryDatabase(db)
	if rows2.Err() != nil {
		t.Fatalf("an error '%s' was not expected while querying", rows2.Err())
	}
	if err != nil {
		t.Fatalf("an error '%s' was not expected while querying", err)
	}

	// Assert that all expectations were met
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
	}
}

func TestQueryDatabase_Error(t *testing.T) {
	// Mock the sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define that an error should be returned when db.Query is called
	mock.ExpectQuery("^SELECT (.+) FROM urls WHERE last_visit_time > \\? ORDER BY last_visit_time DESC$").
		WillReturnError(errors.New("mock error"))

	// Call the function
	getter := chromium.RealQueryDatabaseGetter{}
	rows, err := getter.QueryDatabase(db)
	if err != nil {
		require.Equal(t, "mock error", err.Error())
		return
	}
	if rows.Err() != nil {
		logger.Log.ErrorWithErr("Error closing rows", rows.Err())
	}
}

func TestProcessQueryResults(t *testing.T) {
	// Mock the sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the rows
	rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"}).
		AddRow("http://phishing.com", "Phishing Site", 1, 1000000000).
		AddRow("http://safe.com", "Safe Site", 1, 1000000000)
	mock.ExpectQuery("^SELECT (.+) FROM websites$").WillReturnRows(rows)

	// Execute the query to get *sql.Rows
	realRows, err := db.Query("SELECT * FROM websites")
	if realRows.Err() != nil {
		logger.Log.ErrorWithErr("Error closing rows", realRows.Err())
	}
	if err != nil {
		t.Fatalf("an error '%s' was not expected while querying", err)
	}

	// Mock the GetPhishingDomains function to return a list containing "phishing.com"
	mockGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"http://phishing.com"}, nil
		},
	}
	getter := chromium.RealProcessQueryResultsGetter{}
	results, err := getter.ProcessQueryResults(realRows, mockGetter)

	// Assert there was no error
	require.NoError(t, err)

	// Assert the results contain the phishing domain

	require.True(t, strings.Contains(results[0], "phishing.com"))
}

func TestProcessQueryResults_Error(t *testing.T) {
	// Mock the sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the rows
	rows := sqlmock.NewRows([]string{"url", "title", "visit_count", "last_visit_time"}).
		AddRow("http://phishing.com", "Phishing Site", 1, 1000000000).
		RowError(0, errors.New("mock error"))
	mock.ExpectQuery("^SELECT (.+) FROM websites$").WillReturnRows(rows)

	// Execute the query to get *sql.Rows
	realRows, err := db.Query("SELECT * FROM websites")
	if realRows.Err() != nil {
		logger.Log.ErrorWithErr("Error closing rows", realRows.Err())
	}
	if err != nil {
		t.Fatalf("an error '%s' was not expected while querying", err)
	}

	// Mock the GetPhishingDomains function to return a list containing "phishing.com"
	mockGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return []string{"phishing.com"}, nil
		},
	}

	getter := chromium.RealProcessQueryResultsGetter{}
	_, err = getter.ProcessQueryResults(realRows, mockGetter)

	// Assert there was an error
	require.Error(t, err)
}

func TestProcessQueryResults_GetPhishingDomainsError(t *testing.T) {
	mockGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return nil, errors.New("mock error")
		},
	}

	mockRows := &sql.Rows{}

	getter := chromium.RealProcessQueryResultsGetter{}
	_, err := getter.ProcessQueryResults(mockRows, mockGetter)

	require.Error(t, err)
	require.Equal(t, "mock error", err.Error())
}
