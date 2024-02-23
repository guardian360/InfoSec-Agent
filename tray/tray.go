// Package tray implements the basic functionality of the system tray application
//
// Function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, GetScanCounter, GetScanTicker
package tray

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"InfoSec-Agent/icon"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
	"github.com/skratchdot/open-golang/open"
)

var scanCounter int
var scanTicker *time.Ticker

// OnReady handles all actions that should be handled during the application run-time
//
// Parameters: _
//
// Returns: _
func OnReady() {
	// Icon data can be found in the "icon" package
	systray.SetIcon(icon.Data)

	// Generate the menu for the system tray application

	// Example menu item //

	mGoogleBrowser := systray.AddMenuItem("Google in Browser", "Opens Google in a normal browser")
	systray.AddSeparator()
	///////////////////////
	mChangeScanInterval := systray.AddMenuItem("Change Scan Interval", "Change the interval for scanning")
	mScanNow := systray.AddMenuItem("Scan now", "Scan your device now")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit example tray application")

	// Set up a channel to receive OS signals, used for termination
	// Can be used to notify the application about system termination signals,
	// allowing it to perform possible cleanup tasks before exiting.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	scanCounter = 0
	// Set a ticker to run a scan at a set interval (default = 1 week)
	// Possible future bug: if NewTicker is set to low amount (10 seconds) then calling ScanNow crashes program
	scanTicker = time.NewTicker(7 * 24 * time.Hour)

	// Iterate over each menu option/signal
	for {
		select {
		case <-mGoogleBrowser.ClickedCh:
			err := open.Run("https://www.google.com")
			if err != nil {
				fmt.Println(err)
			}
		case <-mChangeScanInterval.ClickedCh:
			ChangeScanInterval()
		case <-mScanNow.ClickedCh:
			ScanNow()
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		// Executes each time the scanTicker has elapsed the set amount of time
		case <-scanTicker.C:
			scanCounter++
			fmt.Println("Scan:", scanCounter)
		}
	}
}

// OnQuit handles all actions that should happen when the application exits/terminates
//
// Parameters: _
//
// Returns: _
func OnQuit() {
}

// ChangeScanInterval provides the user with a dialog window to set the (new) scan interval
//
// Parameters: optional string testInput, used in tray_test.go
//
// Returns: _
func ChangeScanInterval(testInput ...string) {
	var res string
	// If testInput is provided, use it for testing
	if len(testInput) > 0 {
		res = testInput[0]
	} else {
		// Get user input by creating a dialog window
		var err error
		res, err = zenity.Entry("Enter the scan interval (in hours):", zenity.Title("Change Scan Interval"), zenity.DefaultItems("24"))
		if err != nil {
			fmt.Println("Error creating dialog:", err)
			return
		}
	}

	// Parse the user input
	interval, err := strconv.Atoi(res)
	if err != nil || interval <= 0 {
		fmt.Printf("Invalid input. Using default interval of 24 hours.")
		interval = 24
	}

	// Restart the ticker with the new interval
	scanTicker.Stop()
	scanTicker = time.NewTicker(time.Duration(interval) * time.Hour)
	fmt.Printf("Scan interval changed to %d hours\n", interval)
}

// ScanNow performs one scan iteration (without checking if it is scheduled)
//
// Parameters: _
//
// Returns: _
func ScanNow() {
	scanCounter++
	fmt.Println("Scanning now. Scan:", scanCounter)
	// Manually advance the ticker by sending a signal to its channel
	select {
	case <-scanTicker.C:
	default:
	}
}

// GetScanCounter returns the scanCounter, for use in tray_test.go
//
// Parameters: _
//
// Returns: scanCounter (int)
func GetScanCounter() int {
	return scanCounter
}

// GetScanTicker returns the scanTicker, for use in tray_test.go
//
// Parameters: _
//
// Returns: scanTicker (*time.Ticker)
func GetScanTicker() *time.Ticker {
	return scanTicker
}
