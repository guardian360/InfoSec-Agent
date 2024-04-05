package scan

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	// Necessary to use the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type Severity struct {
	CheckID Int `json:"checkid"`
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
	log.Println("Opening database")
	var err error
	var db *sql.DB
	// Open the database file. If it doesn't exist, it will be created.
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		return
	}
	log.Println("Connected to database")

	// Drop the existing table if it exists
	_, err = db.Exec("DROP TABLE IF EXISTS issues")
	if err != nil {
		log.Println("Error dropping table:", err)
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
		log.Println("Error creating table:", err)
		return
	}

	// Clear issues of rows if they are still there
	_, err = db.Exec("DELETE FROM issues")
	if err != nil {
		log.Println("Error deleting from table:", err)
	}

	var val int64
	// Add dummy values to table
	for i, s := range scanResults {
		val, err = addIssue(db, s, i, 0, 0)
		if err != nil {
			log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 1, 1)
		if err != nil {
			log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 2, 2)
		if err != nil {
			log.Println("Error adding issue: ", err, val)
		}
		val, err = addIssue(db, s, i, 3, 3)
		if err != nil {
			log.Println("Error adding issue: ", err, val)
		}
	}

	// Close the database
	log.Println("Closing database")
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
		"INSERT INTO issues ([JSON Key], [Issue ID], [Result ID], Severity) VALUES (?, ?, ?, ?)",
		strconv.Itoa(check.Issue_ID) ,resultID, issueID, resultID, severity)
	if err != nil {
		return 0, fmt.Errorf("addIssue: %w", err)
	}
	log.Println("Inserted issue")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addIssue: %w", err)
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
	log.Println("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, err
	}
	log.Println("Connected to database")

	var val int
	severities := make([]Severity, len(checks))
	for i, s := range checks {
		val, err = GetSeverity(db, i, resultIDs[i])
		if err != nil {
			log.Println("Error getting severity value")
		}
		severities[i] = Severity{s.ID, val}
	}

	log.Println("Closing database")
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
			log.Println("No rows found.")
		} else {
			log.Println("Error scanning row:", err)
		}
		return 0, err
	}
	return result, nil
}
