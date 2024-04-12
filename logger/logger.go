// Package logger provides a logging mechanism for the application.
// The logging mechanism writes to a log.txt file in the root folder of the application.
package logger

import (
	"log"
	"os"
)

var (
	Log *CustomLogger
)

// CustomLogger is a custom logger struct that wraps the standard logger.
// TODO: fix this docstring according to new documentation standard
// It implements 6 error levels: Trace, Debug, Info, Warning, Error, and Fatal.
type CustomLogger struct {
	*log.Logger
}

// Setup initializes a new logger for the runtime of the application
// TODO: fix this docstring according to new documentation standard
// Parameters: _
//
// Returns: _
func Setup() {
	Log = NewCustomLogger(false)
}

// SetupTests initializes a logger for the runtime of the tests
// This logger does not write to a file, but simply writes to standard output
// TODO: fix this docstring according to new documentation standard
// Parameters: _
//
// Returns: _
func SetupTests() {
	Log = NewCustomLogger(true)
}

// NewCustomLogger creates a new CustomLogger struct.
// The parameter test can specify whether the logger should write to a file.
// During testing, we wish for the logger to write to standard output.
// TODO: fix this docstring according to new documentation standard
// Parameters: test bool - a boolean value that specifies whether the logger is used for testing
//
// Returns: a pointer to a new CustomLogger struct
func NewCustomLogger(test bool) *CustomLogger {
	if test {
		return &CustomLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return &CustomLogger{
		Logger: log.New(file, "", log.LstdFlags),
	}
}

// Print writes a message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Print(message string) {
	l.Println(message)
}

// Trace writes a trace level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Trace(message string) {
	l.Println("TRACE: " + message)
}

// Debug writes a debug level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Debug(message string) {
	l.Println("DEBUG: " + message)
}

// Info writes an info level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Info(message string) {
	l.Println("INFO: " + message)
}

// Warning writes a warning level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Warning(message string) {
	l.Println("WARNING: " + message)
}

// Error writes an error level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Error(message string) {
	l.Println("ERROR: " + message)
}

// ErrorWithErr writes an error level message to the log file, including the error variable
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// err error - the error variable to write to the log file
//
// Returns: _
func (l *CustomLogger) ErrorWithErr(message string, err error) {
	l.Println("ERROR: " + message + " " + err.Error())
}

// Fatal writes a fatal level message to the log file
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// Returns: _
func (l *CustomLogger) Fatal(message string) {
	l.Fatalln("FATAL: " + message)
}

// FatalWithErr writes a fatal level message to the log file, including the error variable
// TODO: fix this docstring according to new documentation standard
// Parameters: message string - the message to write to the log file
//
// err error - the error variable to write to the log file
//
// Returns: _
func (l *CustomLogger) FatalWithErr(message string, err error) {
	l.Fatalln("FATAL: " + message + " " + err.Error())
}
