package scan

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite3 driver
	_ "modernc.org/sqlite"
)

// Severity is a struct that represents the severity level of a check.
// It contains two fields: CheckID and Level.
// CheckID is a string that represents the ID of the check.
// Level is an integer that represents the severity level of the check.
// The severity level is represented as an integer where a higher value indicates a higher severity.
type Severity struct {
	CheckID string `json:"checkid"`
	Level   int    `json:"level"`
}

// SeverityLevel is a type that represents the severity level of an issue.
// It is intended to replace the current integer-based severity level in future versions.
// The possible values are "Safe", "Low", "Medium", and "High".
type SeverityLevel string

const (
	Safe   SeverityLevel = "Safe"
	Low    SeverityLevel = "Low"
	Medium SeverityLevel = "Medium"
	High   SeverityLevel = "High"
)

// SeverityLevels is a struct that encapsulates the severity levels of a security issue.
// It comprises two fields: Value and TSName.
// Value is of type SeverityLevel and represents the severity level of an issue. It can be "Safe", "Low", "Medium", or "High".
// TSName is a string that represents the name of the TypeScript equivalent of the SeverityLevel.
type SeverityLevels struct {
	Value  SeverityLevel
	TSName string
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
	logger.Log.Println("Opening database")
	var err error
	var db *sql.DB
	// Open the database file. If it doesn't exist, it will be created.
	db, err = sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Log.Println("Error opening database:", err)
		return
	}
	logger.Log.Println("Connected to database")

	// Drop the existing table if it exists
	_, err = db.Exec("DROP TABLE IF EXISTS issues")
	if err != nil {
		logger.Log.Println("Error dropping table:", err)
		return
	}

	// Create a new table
	_, err = db.Exec(`CREATE TABLE issues (
                        [JSON Key] TEXT PRIMARY KEY,
                        [Issue ID] INTEGER,
                        [Result ID] INTEGER,
                        severity INTEGER
                    )`)
	if err != nil {
		logger.Log.Println("Error creating table:", err)
		return
	}

	// Clear issues of rows if they are still there
	_, err = db.Exec("DELETE FROM issues")
	if err != nil {
		logger.Log.Println("Error deleting from table:", err)
	}

	var val int64
	// Add dummy values to table
	for i, s := range scanResults {
		val, err = addIssue(db, s, i, 0, 0)
		if err != nil {
			logger.Log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 1, 1)
		if err != nil {
			logger.Log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 2, 2)
		if err != nil {
			logger.Log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 3, 3)
		if err != nil {
			logger.Log.Println("Error adding issue: ", err, val)
		}
	}

	// Close the database
	logger.Log.Println("Closing database")
	defer db.Close()
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
		"INSERT INTO issues ([JSON Key], [Issue ID], [Result ID], Severity) VALUES (?, ?, ?, ?)",
		check.ID+"_"+strconv.Itoa(resultID), issueID, resultID, severity)
	if err != nil {
		return 0, fmt.Errorf("addIssue: %w", err)
	}
	logger.Log.Println("Inserted issue")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addIssue: %w", err)
	}
	return id, nil
}

// GetAllSeverities retrieves the severity levels for all checks performed during a scan.
//
// Parameters:
//   - checks ([]checks.Check): A slice of Check objects obtained from a scan. Each Check object represents a security check that has been performed.
//   - resultIDs ([]int): A slice of integers representing the unique identifiers of the results corresponding to each check.
//
// This function performs the following operations:
//  1. Opens a connection to the SQLite database located at "./database.db". If the database does not exist, it is created.
//  2. Iterates over the checks slice and for each check, retrieves its severity level from the 'issues' table in the database using the GetSeverity function.
//  3. Appends the retrieved severity level to the severities slice.
//  4. Closes the connection to the database.
//
// Returns:
//   - []Severity: A slice of Severity objects representing the severity levels of all checks. Each Severity object contains the ID of the check and its severity level.
//   - error: An error object that describes the error (if any) that occurred while retrieving the severity levels. If no error occurred, this value is nil.
func GetAllSeverities(checks []checks.Check, resultIDs []int) ([]Severity, error) {
	logger.Log.Println("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Log.Println("Error opening database:", err)
		return nil, err
	}
	logger.Log.Println("Connected to database")

	var val int
	severities := make([]Severity, len(checks))
	for i, s := range checks {
		val, err = GetSeverity(db, i, resultIDs[i])
		if err != nil {
			logger.Log.Println("Error getting severity value")
		}
		severities[i] = Severity{s.ID, val}
	}

	logger.Log.Println("Closing database")
	defer db.Close()
	return severities, nil
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
		if err == sql.ErrNoRows {
			logger.Log.Println("No rows found.")
		} else {
			logger.Log.Println("Error scanning row:", err)
		}
		return 0, err
	}
	return result, nil
}
