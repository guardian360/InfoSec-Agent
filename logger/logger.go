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

type CustomLogger struct {
	*log.Logger
}

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

func NewCustomLogger() *CustomLogger {
	file, _ := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	return &CustomLogger{
		Logger: log.New(file, "", log.LstdFlags),
	}
}
func (l *CustomLogger) Print(message string) {
	l.Println(message)
}

func (l *CustomLogger) Trace(message string) {
	l.Println("TRACE: " + message)
}

func (l *CustomLogger) Debug(message string) {
	l.Println("DEBUG: " + message)
}

func (l *CustomLogger) Info(message string) {
	l.Println("INFO: " + message)
}

func (l *CustomLogger) Warning(message string) {
	l.Println("WARNING: " + message)
}

func (l *CustomLogger) Error(message string) {
	l.Println("ERROR: " + message)
}

func (l *CustomLogger) Fatal(message string) {
	l.Fatalln("FATAL: " + message)
}
