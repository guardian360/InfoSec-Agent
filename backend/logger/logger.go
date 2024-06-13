// Package logger provides a logging mechanism for the application.
// The logging mechanism writes to a log.txt file in the root folder of the application.
package logger

import (
	"log"
	"os"
)

// Constants representing the log levels for the logger.
// These constants are used to set the LogLevel and LogLevelSpecific fields of the CustomLogger struct.
// TraceLevel represents the lowest level of logs, used for detailed information about program execution.
// DebugLevel is used for information that may be helpful in diagnosing problems.
// InfoLevel is for general operational information about the program.
// WarningLevel is for information about events that may indicate a problem.
// ErrorLevel is for reporting errors that occur during program execution.
// FatalLevel is for reporting severe errors that may prevent program execution.
const (
	TraceLevel = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

var (
	Log *CustomLogger
)

// CustomLogger is a custom logger struct that wraps the standard logger.
// It implements 6 error levels: Trace, Debug, Info, Warning, Error, and Fatal.
// LogLevel determines the maximum level of messages that will be logged.
// LogLevelSpecific determines the specific level of messages that will be logged.
// If LogLevelSpecific is set to a log-level value, then only that log-level will be written to the log file.
// If LogLevelSpecific is set to -1, then all log-levels up to the specified LogLevel will be written to the log file.
type CustomLogger struct {
	*log.Logger
	LogLevel         int
	LogLevelSpecific int
}

// Setup initializes a new logger for the runtime of the application
//
// This function is used to set up the logger with a specific log level and a specific log level filter.
//
// Parameters:
//
// fileName string - The name of the log file to write to.
//
// logLevel int - The log level to log up to. This should be a value between 0 (TraceLevel) and 5 (FatalLevel).
//
// logLevelSpecific int - The specific log level to log. This should be a value between
// 0 (TraceLevel) and 5 (FatalLevel), or -1 to log all levels up to the specified log level.
//
// Returns: None
func Setup(fileName string, logLevel int, logLevelSpecific int) {
	Log = NewCustomLogger(false, fileName, logLevel, logLevelSpecific)
}

// SetupTests initializes a logger for the runtime of the tests.
//
// This function is used to set up the logger with a specific log level and a specific log level filter.
// The log level is set to 0 (TraceLevel) and the log level specific is set to -1, meaning all log levels will be logged.
// This logger does not write to a file, but simply writes to standard output.
//
// Parameters: None
//
// Returns: None
func SetupTests() {
	Log = NewCustomLogger(true, "", 0, 0)
}

// NewCustomLogger creates a new CustomLogger struct.
//
// This function is used to create a new CustomLogger with a specific log level and a specific log level filter.
// The test parameter determines whether the logger should write to a file or to standard output.
// During testing, we wish for the logger to write to standard output.
//
// Parameters:
//
// test bool - A boolean value that specifies whether the logger is used for testing.
//
// fileName string - The name of the log file to write to.
//
// logLevel int - The log level to log up to. This should be a value between 0 (TraceLevel) and 5 (FatalLevel).
//
// logLevelSpecific int - The specific log level to log. This should be a value between 0 (TraceLevel) and 5 (FatalLevel), or -1 to log all levels up to the specified log level.
//
// Returns: a pointer to a new CustomLogger struct
func NewCustomLogger(test bool, fileName string, logLevel int, logLevelSpecific int) *CustomLogger {
	if test {
		return &CustomLogger{
			Logger:           log.New(os.Stdout, "", log.LstdFlags),
			LogLevel:         logLevel,
			LogLevelSpecific: logLevelSpecific,
		}
	}
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("error setting up logger: error getting user config dir", err)
	}
	// Create the InfoSec-Agent directory in the AppData folder if it does not exist
	dirPath := appDataPath + `\InfoSec-Agent\`
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal("error setting up logger: error creating InfoSec-Agent dir", err)
	}

	// Create the log file in the InfoSec-Agent directory or truncate it if it already exists
	logPath := dirPath + fileName
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("error setting up logger: error opening log file", err)
	}
	return &CustomLogger{
		Logger:           log.New(file, "", log.LstdFlags),
		LogLevel:         logLevel,
		LogLevelSpecific: logLevelSpecific,
	}
}

// Print writes a general message to the log file.
//
// This method is used to write a message to the log file. The message is not associated with any specific log level.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Print(message string) {
	l.Println(message)
}

// Trace writes a trace level message to the log file.
//
// The message will only be logged if the LogLevelSpecific of the logger is set to TraceLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to TraceLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Trace(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= TraceLevel || l.LogLevelSpecific == TraceLevel {
		l.Println("TRACE: " + message)
	}
}

// Debug writes a debug level message to the log file.
//
// The message will only be logged if the LogLevelSpecific of the logger is set to DebugLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to DebugLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Debug(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= DebugLevel || l.LogLevelSpecific == DebugLevel {
		l.Println("DEBUG: " + message)
	}
}

// Info writes an info level message to the log file.
//
// The message will only be logged if the LogLevelSpecific of the logger is set to InfoLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to Info.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Info(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= InfoLevel || l.LogLevelSpecific == InfoLevel {
		l.Println("INFO: " + message)
	}
}

// Warning writes a warning level message to the log file.
//
// The message will only be logged if the LogLevelSpecific of the logger is set to WarningLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to WarningLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Warning(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= WarningLevel || l.LogLevelSpecific == WarningLevel {
		l.Println("WARNING: " + message)
	}
}

// Error writes an error level message to the log file.
//
// The message will only be logged if the LogLevelSpecific of the logger is set to ErrorLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to ErrorLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Error(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= ErrorLevel || l.LogLevelSpecific == ErrorLevel {
		l.Println("ERROR: " + message)
	}
}

// ErrorWithErr writes an error level message to the log file, including the error variable.
//
// The message and the error variable will only be logged if the LogLevelSpecific of the logger is set to ErrorLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to ErrorLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// err error - The error variable to write to the log file.
//
// Returns: None
func (l *CustomLogger) ErrorWithErr(message string, err error) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= ErrorLevel || l.LogLevelSpecific == ErrorLevel {
		l.Println("ERROR: " + message + " " + err.Error())
	}
}

// Fatal writes a fatal level message to the log file, including the error variable.
//
// The message and the error variable will only be logged if the LogLevelSpecific of the logger is set to FatalLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to FatalLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// Returns: None
func (l *CustomLogger) Fatal(message string) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= FatalLevel || l.LogLevelSpecific == FatalLevel {
		l.Fatalln("FATAL: " + message)
	}
}

// FatalWithErr writes a fatal level message to the log file, including the error variable.
//
// The message and the error variable will only be logged if the LogLevelSpecific of the logger is set to FatalLevel,
// or if the LogLevelSpecific is set to -1 and the LogLevel is less than or equal to FatalLevel.
//
// Parameters:
//
// message string - The message to write to the log file.
//
// err error - The error variable to write to the log file.
//
// Returns: None
func (l *CustomLogger) FatalWithErr(message string, err error) {
	if l.LogLevelSpecific == -1 && l.LogLevel <= FatalLevel || l.LogLevelSpecific == FatalLevel {
		l.Fatalln("FATAL: " + message + " " + err.Error())
	}
}
