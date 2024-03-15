// Package main contains the entry point of the application.
//
// Function(s): main
package main

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/tray"
	"github.com/getlantern/systray"
)

// main is the entry point of the program, starts the system tray application
//
// Parameters: _
//
// Returns: _
func main() {
	localization.Init("")
	systray.Run(tray.OnReady, tray.OnQuit)
}
