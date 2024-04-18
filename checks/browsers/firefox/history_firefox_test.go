package firefox_test

import (
	"errors"
	"fmt"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
	"os"
	"testing"

	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/checks/browsers/firefox"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.SetupTests()
	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

var profileFinder utils.FirefoxProfileFinder

func TestHistoryFirefox(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Mock the FirefoxFolder function to return a valid directory and no error
		profileFinder = utils.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return []string{"/valid/directory"}, nil
			},
		}

		check := firefox.HistoryFirefox(profileFinder)
		require.Nil(t, check.Result)
		require.Error(t, check.Error)
	})

	t.Run("failure", func(t *testing.T) {
		// Mock the FirefoxFolder function to return an error
		profileFinder = utils.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return nil, errors.New("mock error")
			},
		}

		profileFinder = utils.MockProfileFinder{
			MockFirefoxFolder: func() ([]string, error) {
				return nil, errors.New("mock error")
			},
		}

		check := firefox.HistoryFirefox(profileFinder)
		require.Nil(t, check.Result)
		require.Error(t, check.Error)
	})
}

/*
func TestCloseRows_Error(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	rows := sqlmock.NewRows([]string{"column"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	realRows, err := db.Query("SELECT")
	require.NoError(t, err)

	// Simulate an error when closing the rows
	mock.ExpectClose().WillReturnError(errors.New("sql: Rows are already closed"))

	var buf bytes.Buffer
	logger.Log.SetOutput(&buf)

	// This should log an error because the rows are already closed
	firefox.CloseRows(realRows)

	// Wait for the function to complete
	time.Sleep(100 * time.Millisecond)

	capturedOutput := buf.String()

	// Check if the expected error was logged
	require.Contains(t, capturedOutput, "Error closing rows:")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}*/

func TestHistoryFirefox_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	profileFinder = utils.RealProfileFinder{}

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WillReturnRows(sqlmock.NewRows([]string{"url", "last_visit_date"}))

	// now we execute our method
	check := firefox.HistoryFirefox(profileFinder)
	if check.Error != nil {
		t.Errorf("error was not expected while updating stats: %s", check.Error)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHistoryFirefox_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the FirefoxFolder function to return a valid directory
	profileFinder = utils.MockProfileFinder{
		MockFirefoxFolder: func() ([]string, error) {
			return []string{"/valid/directory"}, nil
		},
	}

	mock.ExpectQuery("^SELECT url, last_visit_date FROM moz_places WHERE last_visit_date >= \\? ORDER BY last_visit_date DESC$").
		WillReturnError(fmt.Errorf("some error"))

	// now we execute our method
	check := firefox.HistoryFirefox(profileFinder)
	if check.Error == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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

	defer db.Close()

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

	defer db.Close()

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

	defer db.Close()

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

	defer db.Close()

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
				{URL: "http://startpagina.nl/path", LastVisitDate: sql.NullInt64{Int64: 1712580339000000, Valid: true}},
				{URL: "http://001return.com/path", LastVisitDate: sql.NullInt64{Int64: 1713580339000000, Valid: true}},
				{URL: "http://012345bet.com/path", LastVisitDate: sql.NullInt64{Int64: 1712580379000000, Valid: true}},
			},
			want: []string{"startpagina.nlGood", "001return.comGood", "012345bet.comGood"},
		},
		{
			name:    "no results",
			results: []firefox.QueryResult{},
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := firefox.ProcessQueryResults(tt.results)
			require.Equal(t, tt.want, got)
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
