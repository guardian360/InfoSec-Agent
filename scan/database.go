package scan

import (
	"database/sql"
	"errors"
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

// FillDataBase clears the existing issues table and populates it with the results from a scan.
//
// Parameters:
//   - scanResults ([]checks.Check): A slice of Check objects obtained from a scan. Each Check object represents a security check that has been performed.
//
// This function performs the following operations:
//  1. Opens a connection to the SQLite database located at "./database.db". If the database does not exist, it is created.
//  2. Drops the existing "issues" table if it exists.
//  3. Creates a new "issues" table with columns for JSON Key, Issue ID, Result ID, and Severity.
//  4. Clears any existing rows in the "issues" table.
//  5. Iterates over the scanResults slice and adds each Check object to the "issues" table as a new row.
//  6. Closes the connection to the database.
//
// Note: This function logs any errors that occur during its execution and does not return any values.
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
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database:", err)
		}
	}(db)
}

// addIssue inserts a new issue into the 'issues' table in the database.
//
// Parameters:
//   - db (*sql.DB): The database connection where the 'issues' table resides.
//   - check (checks.Check): The issue to be added to the table. Each issue represents a security check that has been performed.
//   - issueID (int): The unique identifier of the issue. This is used as a reference in the 'issues' table.
//   - resultID (int): The unique identifier of the result of the issue. This is used as a reference in the 'issues' table.
//   - severity (int): The severity level of the issue. This is represented as an integer where a higher value indicates a higher severity.
//
// Returns:
//   - int64: The index of the newly added row in the 'issues' table.
//   - error: An error object that describes the error (if any) that occurred while adding the issue to the table. If no error occurred, this value is nil.
//
// This function performs the following operations:
//  1. Prepares an SQL INSERT statement to add a new row to the 'issues' table.
//  2. Executes the SQL statement, passing in the parameters of the function.
//  3. If an error occurs while executing the SQL statement, it logs the error and returns it along with a zero value for the index.
//  4. If the SQL statement executes successfully, it retrieves the index of the newly added row and returns it along with a nil error.
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

// GetSeverity retrieves the severity level of a specific issue from the 'issues' table in the database.
//
// Parameters:
//   - db (*sql.DB): The database connection where the 'issues' table resides.
//   - issueID (int): The unique identifier of the issue. This is used as a reference in the 'issues' table.
//   - resultID (int): The unique identifier of the result of the issue. This is used as a reference in the 'issues' table.
//
// This function performs the following operations:
//  1. Prepares an SQL SELECT statement to retrieve the severity level of the issue from the 'issues' table.
//  2. Executes the SQL statement, passing in the issueID and resultID as parameters.
//  3. If an error occurs while executing the SQL statement, it logs the error and returns it along with a zero value for the severity level.
//  4. If the SQL statement executes successfully, it retrieves the severity level from the result and returns it along with a nil error.
//
// Returns:
//   - int: The severity level of the issue. This is represented as an integer where a higher value indicates a higher severity.
//   - error: An error object that describes the error (if any) that occurred while retrieving the severity level. If no error occurred, this value is nil.
func GetSeverity(db *sql.DB, issueID int, resultID int) (int, error) {
	// Prepare the SQL query
	query := "SELECT Severity FROM issues WHERE [Issue ID] = ? AND [Result ID] = ?"

	// Query the database
	row := db.QueryRow(query, issueID, resultID)

	var result int
	// Scan the value from the row into the integer variable
	err := row.Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
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
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing database:", err)
		}
	}(db)
	return dbData, nil
}
