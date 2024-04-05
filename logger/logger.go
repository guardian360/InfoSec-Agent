// Package logger provides a logging mechanism for the application.
// The logging mechanism writes to a log.txt file in the root folder of the application.
package logger

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

// Setup initializes a logger for the runtime of the application
// The logger writes to a log.txt file in the root folder
//
// Parameters: _
//
// Returns: _
func Setup() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Log = log.New(file, "", log.LstdFlags)
}

// SetupTests initializes a logger for the runtime of the tests
// This logger does not write to a file, but simply writes to standard output
//
// Parameters: _
//
// Returns: _
func SetupTests() {
	Log = log.New(os.Stdout, "", log.LstdFlags)
}
