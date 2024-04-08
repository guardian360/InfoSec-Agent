package scan

import (
	"database/sql"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite3 driver
	_ "modernc.org/sqlite"
)

// DataBaseData is a struct which is used to format extracted information from the database
//
// CheckId is used as the identifier to connect the severity level and JSON key to
type DataBaseData struct {
	CheckID  int `json:"id"`
	Severity int `json:"severity"`
	JSONKey  int `json:"jsonkey"`
}

// FillDataBase will remove the current issues table and create a new one filled with dummy values
//
// Parameters: scanResults ([]checks.Check) - the list of checks from a scan
//
// Returns: _
func FillDataBase(scanResults []checks.Check) {
	logger.Log.Info("Opening database")
	var err error
	var db *sql.DB
	// Open the database file. If it doesn't exist, it will be created.
	db, err = sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Log.ErrorWithErr("Error opening database:", err)
		return
	}
	logger.Log.Info("Connected to database")

	// Clear rows of issues if they are still there
	_, err = db.Exec("DELETE FROM issues")
	if err != nil {
		logger.Log.ErrorWithErr("Error deleting from table:", err)
	}

	var val int64
	// Add dummy values to table
	// addIssue's second argument should become s.id and the specific results and severities should be used
	for _, s := range scanResults {
		val, err = addIssue(db, s, s.IssueID, 0, 0)
		if err != nil {
			logger.Log.Printf("Error adding issue: %s %s", err, strconv.FormatInt(val, 10))
		}
		val, err = addIssue(db, s, s.IssueID, 1, 1)
		if err != nil {
			logger.Log.Printf("Error adding issue: %s %s", err, strconv.FormatInt(val, 10))
		}
		val, err = addIssue(db, s, s.IssueID, 2, 2)
		if err != nil {
			logger.Log.Printf("Error adding issue: %s %s", err, strconv.FormatInt(val, 10))
		}
		val, err = addIssue(db, s, s.IssueID, 3, 3)
		if err != nil {
			logger.Log.Printf("Error adding issue: %s %s", err, strconv.FormatInt(val, 10))
		}
	}

	// Close the database
	logger.Log.Info("Closing database")
	defer db.Close()
}

// addIssue will add a single entry in the issues table
//
// Parameters:
//
// db (*sql.DB) - database connection where table resides
//
// check (checks.Check) - the issue to be added
//
// issueID (int) - id of the issue
//
// resultID (int) - id of the result of the issue
//
// severity (int) - severity of the result
//
// Returns: returns index of the added row in the table
func addIssue(db *sql.DB, check checks.Check, issueID int, resultID int, severity int) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO issues ([Issue ID], [Result ID], Severity, [JSON Key]) VALUES (?, ?, ?, ?)",
		resultID, issueID, resultID, severity, strconv.Itoa(check.IssueID))
	if err != nil {
		return 0, err
	}
	logger.Log.Info("Inserted issue with issueID: " + strconv.Itoa(issueID))
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetSeverity gets the single severity of an issue
//
// Parameters:
//
// db (*sql.DB) - database connection where table resides
//
// issueId (int) - id of the issue
//
// resultID (int) - id of the result of the issue
//
// Returns: severity of the issue
func GetSeverity(db *sql.DB, issueID int, resultID int) (int, error) {
	// Prepare the SQL query
	query := "SELECT Severity FROM issues WHERE [Issue ID] = ? AND [Result ID] = ?"

	// Query the database
	row := db.QueryRow(query, issueID, resultID)

	var result int
	// Scan the value from the row into the integer variable
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Warning("No rows found.")
		} else {
			logger.Log.ErrorWithErr("Error scanning row:", err)
		}
		return 0, err
	}
	return result, nil
}

// GetJSONKey gets the single JSON key of an issue
//
// Parameters:
//
// db (*sql.DB) - database connection where table resides
//
// issueId (int) - id of the issue
//
// resultId (int) - id of the result of the issue
//
// Returns: JSON key of the issue
func GetJSONKey(db *sql.DB, issueID int, resultID int) (int, error) {
	// Prepare the SQL query
	query := "SELECT [JSON Key] FROM issues WHERE [Issue ID] = ? AND [Result ID] = ?"

	// Query the database
	row := db.QueryRow(query, issueID, resultID)

	var result int
	// Scan the value from the row into the integer variable
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found.")
		} else {
			logger.Log.ErrorWithErr("Error scanning row:", err)
		}
		return result, err
	}
	return result, nil
}

// GetDataBaseData gets the severities and JSON keys for all checks passed
//
// Parameters: checks ([]checks.Check) - the list of checks from a scan
//
// resultIDs ([]int) - the list of results corresponding to each check
//
// Returns: list of all severities and JSON keys
func GetDataBaseData(checks []checks.Check) ([]DataBaseData, error) {
	logger.Log.Info("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Log.ErrorWithErr("Error opening database:", err)
		return nil, err
	}
	logger.Log.Info("Connected to database")

	dbData := make([]DataBaseData, len(checks))
	for i, s := range checks {
		sev, err2 := GetSeverity(db, s.IssueID, s.ResultID)
		if err2 != nil {
			logger.Log.Printf("Error getting severity value for IssueID:%v and ResultID:%v", s.IssueID, s.ResultID)
		}
		jsn, err3 := GetJSONKey(db, s.IssueID, s.ResultID)
		if err3 != nil {
			logger.Log.Printf("Error getting severity value for IssueID:%v and ResultID:%v", s.IssueID, s.ResultID)
		}
		dbData[i] = DataBaseData{s.IssueID, sev, jsn}
	}

	logger.Log.Info("Closing database")
	defer db.Close()
	return dbData, nil
}
