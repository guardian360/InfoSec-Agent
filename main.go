// Package main contains the entry point of the application.
//
// Exported function(s): _
package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
	"github.com/getlantern/systray"
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
	localization.Init("")
	systray.Run(tray.OnReady, tray.OnQuit)
}
