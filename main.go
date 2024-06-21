// The InfoSec Agent is a security and privacy tool for Windows 10 and 11.
// This module is a lightweight, user-friendly system tray application that identifies security and privacy weaknesses.
// It does this by scanning the device and programs for vulnerabilities and misconfigurations.
//
// This package links together with the reporting page module, which is a wails application that displays the results of the scan and provides insights and solutions.
//
// This project is a collaborative effort involving nine students from Utrecht University in The Netherlands, in partnership with the Dutch IT company Guardian360.
// It serves as the Software Project for the Bachelor's Programme in Computing Sciences at the UU.
// This project is also supported by funding from the SIDN Fund (Stichting Internet Domeinregistratie Nederland), the Dutch domain name registrar.
package main

//go:generate go-winres make --in scripts/winres/winres.json --product-version=git-tag --file-version=git-tag

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/tray"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/getlantern/systray"
	"github.com/rodolfoag/gow32"
)

// main is the entry point of the application. It initializes the localization settings and starts the system tray application.
//
// This function first calls the Init function from the localization package to set up the localization settings for the application. The empty string argument means that the default language setting will be used.
// After setting up localization, it calls the Run function from the systray package to start the system tray application. The OnReady and OnQuit functions from the tray package are passed as arguments to the Run function. OnReady is called to set up the system tray when the application starts, and OnQuit is called to perform cleanup operations when the application is about to terminate.
//
// Parameters: None.
//
// Returns: None. This function does not return a value as it is the entry point of the application.
func main() {
	// Create a mutex to ensure only one instance of the application is running
	// If the mutex already exists, it means another instance of the application is running, so we exit
	// This also ensures the program is not running when uninstalling the application
	_, mutexErr := gow32.CreateMutex("InfoSec-Agent")
	if mutexErr != nil {
		// Initialize localization settings for startup popups
		logger.SetupTests()
		localization.Init("")
		settings := usersettings.LoadUserSettings()
		tray.Language = settings.Language
		tray.AlreadyRunningPopup()
		return
	}

	// Set up the logger, passing the log-level you desire (it logs everything equal and lower to the log-level):
	// 0 - Trace
	// 1 - Debug
	// 2 - Info
	// 3 - Warning
	// 4 - Error
	// 5 - Fatal
	// The second argument is the specific log-level you want to log, giving this a value will only log that level.
	// If you want to log all levels up to the specified level, pass -1.
	logger.Setup("log.txt", config.LogLevel, config.LogLevelSpecific)
	logger.Log.Info("Starting InfoSec Agent")

	// Initialize localization settings
	localization.Init("")
	settings := usersettings.LoadUserSettings()
	tray.Language = settings.Language

	// Start Tray
	tray.StartPopup()
	systray.Run(tray.OnReady, tray.OnQuit)
}
