package scan

import (
	"database/sql"
	"fmt"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"

	_ "github.com/mattn/go-sqlite3"
)

func FillDataBase(scanResults []checks.Check) {
	fmt.Println("Opening database")
	// Open the database file. If it doesn't exist, it will be created.
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	fmt.Println("Connected to database")

	_, err = db.Exec("ALTER TABLE issues RENAME COLUMN Issue ID TO Issue_ID")
	if err != nil {
		fmt.Println("Error renaming column:", err)
		return
	}

	// for i, s := range scanResults {
	// 	_, err := addIssue(db, s, i, i, 1)
	// 	if err != nil {
	// 		fmt.Println("Error adding issue: ", err)
	// 	}
	// }
	fmt.Println("Closing database")
	defer db.Close()
}

func addIssue(db *sql.DB, check checks.Check, issueId int, resultId int, severity int) (int64, error) {
	result, err := db.Exec("INSERT INTO issues (Issue ID, Result ID, Severity, JSON Key) VALUES (?, ?, ?, ?)", issueId, resultId, severity, check.Id)
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
