package scan

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"

	_ "github.com/mattn/go-sqlite3"
)

type Severity struct {
	CheckId string `json:"checkid"`
	Level   int    `json:"level"`
}

// future enumerator replacing type of level int with SeverityLevel
type SeverityLevel string

const (
	Safe   SeverityLevel = "Safe"
	Low    SeverityLevel = "Low"
	Medium SeverityLevel = "Medium"
	High   SeverityLevel = "High"
)

type SeverityLevels struct {
	Value  SeverityLevel
	TSName string
}

// FillDataBase will remove the current issues table and create a new one filled with dummy values
//
// Parameters: scanResults ([]checks.Check) - the list of checks from a scan
//
// Returns: _
func FillDataBase(scanResults []checks.Check) {
	fmt.Println("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	fmt.Println("Connected to database")

	// Drop the existing table if it exists
	_, err = db.Exec("DROP TABLE IF EXISTS issues")
	if err != nil {
		fmt.Println("Error dropping table:", err)
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
		fmt.Println("Error creating table:", err)
		return
	}

	// Clear issues of rows if they are still there
	_, err = db.Exec("DELETE FROM issues")
	if err != nil {
		fmt.Println("Error deleting from table:", err)
	}

	// Add dummy values to table
	for i, s := range scanResults {
		_, err := addIssue(db, s, i, 0, 0)
		if err != nil {
			fmt.Println("Error adding issue: ", err)
		}
		_, err = addIssue(db, s, i, 1, 1)
		if err != nil {
			fmt.Println("Error adding issue: ", err)
		}
		_, err = addIssue(db, s, i, 2, 2)
		if err != nil {
			fmt.Println("Error adding issue: ", err)
		}
		_, err = addIssue(db, s, i, 3, 3)
		if err != nil {
			fmt.Println("Error adding issue: ", err)
		}
	}

	// Close the database
	fmt.Println("Closing database")
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
// issueId (int) - id of the issue
//
// resultId (int) - id of the result of the issue
//
// severity (int) - severity of the result
//
// Returns: returns index of the added row in the table
func addIssue(db *sql.DB, check checks.Check, issueId int, resultId int, severity int) (int64, error) {
	result, err := db.Exec("INSERT INTO issues ([JSON Key], [Issue ID], [Result ID], Severity) VALUES (?, ?, ?, ?)", check.Id+"_"+strconv.Itoa(resultId), issueId, resultId, severity)
	if err != nil {
		return 0, fmt.Errorf("addIssue: %v", err)
	}
	fmt.Println("Inserted issue")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addIssue: %v", err)
	}
	return id, nil
}

// GetAllSeverities get the severities for all checks passed
//
// Parameters: checks ([]checks.Check) - the list of checks from a scan
//
// resultIDs ([]int) - the list of results corresponding to each check
//
// Returns: list of all severities
func GetAllSeverities(checks []checks.Check, resultIDs []int) ([]Severity, error) {
	fmt.Println("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}
	fmt.Println("Connected to database")

	severities := make([]Severity, len(checks))
	for i, s := range checks {
		val, err := GetSeverity(db, i, resultIDs[i])
		if err != nil {
			fmt.Println("Error getting severity value")
		}
		severities[i] = Severity{s.Id, val}
	}

	fmt.Println("Closing database")
	defer db.Close()
	return severities, nil
}

// GetSeverity gets the single severity of an issue
//
// Parameters:
//
// db (*sql.DB) - database connection where table resides
//
// issueId (int) - id of the issue
//
// resultId (int) - id of the result of the issue
//
// Returns: severity of the issue
func GetSeverity(db *sql.DB, issueID int, resultId int) (int, error) {
	// Prepare the SQL query
	query := "SELECT Severity FROM issues WHERE [Issue ID] = ? AND [Result ID] = ?"

	// Query the database
	row := db.QueryRow(query, issueID, resultId)

	var result int
	// Scan the value from the row into the integer variable
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found.")
		} else {
			fmt.Println("Error scanning row:", err)
		}
		return 0, err
	}
	return result, nil
}
