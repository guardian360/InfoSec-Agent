package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
)

// Tray serves as an interface between the user interface and the system tray operations.
//
// It provides methods to perform actions such as initiating an immediate scan, changing the language, and altering the scan interval of the system tray application. It does not contain any fields as it is primarily used as a receiver for these methods.
type Tray struct {
	Log *logger.CustomLogger
}

// NewTray is a constructor that returns a pointer to a Tray instance.
//
// The Tray instance serves as a bridge between the front-end and the system tray functions, enabling the invocation of system tray operations.
//
// Parameters: None.
//
// Returns: *Tray: A pointer to a Tray instance.
func NewTray(log *logger.CustomLogger) *Tray {
	return &Tray{
		Log: log,
	}
}

// ScanNow initiates an immediate scan operation via the system tray application.
//
// This method is a bridge between the front-end and the tray package's ScanNow function. It triggers an immediate scan operation, bypassing the regular scan interval. The scan results, represented as a slice of checks, are returned along with any error that might occur during the scan.
//
// Parameters: dialogPresent (bool) - A boolean value that indicates whether a dialog should be displayed during the scan operation.
//
// Returns:
//   - []checks.Check: A slice of checks representing the scan results.
//   - error: An error object that describes the error, if any occurred. nil if no error occurred.
func (t *Tray) ScanNow(dialogPresent bool) ([]checks.Check, error) {
	return tray.ScanNow(dialogPresent, "../"+config.DatabasePath)
}

// GetInstalledPrograms is a method of the Tray struct that retrieves a list of installed programs on the local machine.
// It uses the InstalledSoftware function from the programs package to perform this operation.
// The InstalledSoftware function is called with a RealCommandExecutor and the LocalMachine as arguments.
// The RealCommandExecutor is an implementation of the CommandExecutor interface that executes real commands on the local machine.
// The LocalMachine is a predefined constant that represents the local machine in the context of registry operations.
//
// This method does not take any parameters.
//
// Returns:
//   - checks.Check: A Check object representing the result of the InstalledSoftware function. This object contains information about the installed programs.
func (t *Tray) GetInstalledPrograms() checks.Check {
	return programs.InstalledSoftware(&mocking.RealCommandExecutor{}, mocking.LocalMachine)
}

// ChangeLanguage is responsible for switching the language of the system tray application.
//
// This method invokes the ChangeLanguage function from the tray package, which is responsible for changing the language of the system tray application. The language change is applied immediately upon invocation, and the menu is refreshed to reflect the new language.
//
// Parameters: None. The method uses the receiver (*Tray) to call the function.
//
// Returns: None. The method performs an action and does not return any value.
func (t *Tray) ChangeLanguage() {
	tray.ChangeLanguage()
	tray.RefreshMenu()
}

// ChangeScanInterval triggers a change in the scanning interval of the system tray application.
//
// This method invokes the ChangeScanInterval function from the tray package, which is responsible for altering the frequency at which the system tray application performs its checks. The change in scanning interval is applied immediately upon invocation.
//
// Parameters: None. The method uses the receiver (*Tray) to call the function.
//
// Returns: None. The method performs an action and does not return any value.
func (t *Tray) ChangeScanInterval() {
	tray.ChangeScanInterval()
}

// LogDebug logs a Debug level message to the log file.
// This function is used for logging messages from the front-end JS.
//
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) LogDebug(message string) {
	t.Log.Debug(message)
}

// LogInfo logs an Info level message to the log file.
// This function is used for logging messages from the front-end JS.
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) LogInfo(message string) {
	t.Log.Info(message)
}

// LogError logs an error level message to the log file.
// This function is used for logging messages from the front-end JS.
// Parameters: t (*Tray) - a pointer to the Tray struct
//
// Returns: _
func (t *Tray) LogError(message string) {
	t.Log.Error(message)
}
