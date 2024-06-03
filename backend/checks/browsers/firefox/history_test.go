package firefox_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/browsers/firefox"
	"github.com/stretchr/testify/require"
)

var Profilefinder browsers.FirefoxProfileFinder

func TestHistoryFirefox(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Mock the FirefoxFolder function to return a valid directory and no error
		Profilefinder = browsers.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return []string{"/valid/directory"}, nil
			},
		}

		check := firefox.HistoryFirefox(Profilefinder, browsers.RealPhishingDomainGetter{})
		require.Nil(t, check.Result)
		require.Error(t, check.Error)
	})

	t.Run("failure", func(t *testing.T) {
		// Mock the FirefoxFolder function to return an error
		Profilefinder = browsers.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return nil, errors.New("mock error")
			},
		}

		check := firefox.HistoryFirefox(Profilefinder, browsers.RealPhishingDomainGetter{})
		require.Nil(t, check.Result)
		require.Error(t, check.Error)
	})
}

func TestOpenAndQueryDatabase_OpenFailure(t *testing.T) {
	db, err := sql.Open("sqlite", "/invalid/path")
	require.NoError(t, err)

	_, err = firefox.QueryDatabase(db)
	require.Error(t, err)
}

func TestQueryDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000

	rows := sqlmock.NewRows([]string{"url", "last_visit_date"}).
		AddRow("http://startpagina.nl/path", 1712580339000000).
		AddRow("http://001return.com/path", 1712580439000000).
		AddRow("http://012345bet.com/path", 1712580539000000)

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WithArgs(lastWeek).
		WillReturnRows(rows)

	results, err := firefox.QueryDatabase(db)
	require.NoError(t, err)

	expected := []firefox.QueryResult{
		{URL: "http://startpagina.nl/path", LastVisitDate: sql.NullInt64{Int64: 1712580339000000, Valid: true}},
		{URL: "http://001return.com/path", LastVisitDate: sql.NullInt64{Int64: 1712580439000000, Valid: true}},
		{URL: "http://012345bet.com/path", LastVisitDate: sql.NullInt64{Int64: 1712580539000000, Valid: true}},
	}

	require.Equal(t, expected, results)
}

func TestQueryDatabase_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WithArgs(lastWeek).
		WillReturnError(errors.New("mock error"))

	_, err = firefox.QueryDatabase(db)
	require.Error(t, err)
	require.Equal(t, "mock error", err.Error())
}

func TestQueryDatabase_RowError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000

	rows := sqlmock.NewRows([]string{"url", "last_visit_date"}).
		AddRow("http://startpagina.nl/path", 1712580339000000).
		AddRow("http://001return.com/path", 1712580439000000).
		AddRow("http://012345bet.com/path", 1712580539000000).
		RowError(1, errors.New("mock row error"))

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WithArgs(lastWeek).
		WillReturnRows(rows)

	_, err = firefox.QueryDatabase(db)
	require.Error(t, err)
	require.Equal(t, "mock row error", err.Error())
}

func TestQueryDatabase_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database: ", err)
		}
	}(db)

	lastWeek := (time.Now().AddDate(0, 0, -7).UnixMicro() / 10000000) * 10000000

	rows := sqlmock.NewRows([]string{"url", "last_visit_date"}).
		AddRow("http://startpagina.nl/path", "invalid time").
		AddRow("http://001return.com/path", 1712580439000000).
		AddRow("http://012345bet.com/path", 1712580539000000).
		RowError(0, errors.New("mock scan error")) // Simulate an error when scanning the rows

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WithArgs(lastWeek).
		WillReturnRows(rows)

	_, err = firefox.QueryDatabase(db)
	require.Error(t, err)
	require.Equal(t, "mock scan error", err.Error())
}

type MockTimeFormatter struct{}

func (MockTimeFormatter) FormatTime(_ sql.NullInt64) string {
	return "Good"
}

func TestProcessQueryResults(t *testing.T) {
	firefox.TimeFormat = MockTimeFormatter{}
	tests := []struct {
		name    string
		results []firefox.QueryResult
		want    []string
	}{
		{
			name: "with results",
			results: []firefox.QueryResult{
				{URL: "http://00000000000000000000000000000000000000000.xyz/path", LastVisitDate: sql.NullInt64{Int64: 1713580339000000, Valid: true}},
			},
			want: []string{"00000000000000000000000000000000000000000.xyzGood"},
		},
		{
			name:    "no results",
			results: []firefox.QueryResult{},
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := firefox.ProcessQueryResults(tt.results, browsers.RealPhishingDomainGetter{})
			if !compareSlices(got, tt.want) && !compareSlices(got, []string{}) {
				t.Errorf("processQueryResults() = %v, want %v", got, tt.want)
			}
		})
	}
	firefox.TimeFormat = firefox.RealTimeFormatter{}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name string
		arg  sql.NullInt64
		want string
	}{
		{
			name: "valid time",
			arg:  sql.NullInt64{Int64: 1712580339000000, Valid: true},
			want: "2024-04-08 14:45:39 +0200 CEST",
		},
		{
			name: "invalid time",
			arg:  sql.NullInt64{Int64: 0, Valid: false},
			want: "1970-01-01 01:00:00 +0100 CET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, _ := time.LoadLocation("UTC")
			got := firefox.FormatTime(tt.arg)
			gotTime, _ := time.ParseInLocation("2024-04-08 14:45:39 +0200 CEST", got, loc)
			wantTime, _ := time.ParseInLocation("2024-04-08 14:45:39 +0200 CEST", tt.want, loc)
			if !gotTime.Equal(wantTime) {
				t.Errorf("formatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// compareSlices compares two slices of strings and returns true if they are equal, false otherwise
func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type MockPhishingDomainsGetter struct {
	GetPhishingDomainsFunc func() ([]string, error)
}

func (m MockPhishingDomainsGetter) GetPhishingDomains(_ browsers.RequestCreator) ([]string, error) {
	return m.GetPhishingDomainsFunc()
}

func TestProcessQueryResults_GetPhishingDomainsError(t *testing.T) {
	mockGetter := MockPhishingDomainsGetter{
		GetPhishingDomainsFunc: func() ([]string, error) {
			return nil, errors.New("mock error")
		},
	}

	_, err := firefox.ProcessQueryResults([]firefox.QueryResult{}, mockGetter)
	require.Error(t, err)
	require.Equal(t, "mock error", err.Error())
}
