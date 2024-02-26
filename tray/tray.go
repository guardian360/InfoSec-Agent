// Package tray implements the basic functionality of the system tray application
//
// Function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, GetScanCounter, GetScanTicker
package tray

import (
	"InfoSec-Agent/localization"
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
var language = 1 // Default language is British English

type MenuItem struct {
	menuTitle   string
	menuTooltip string
	sysMenuItem *systray.MenuItem
}

// OnReady handles all actions that should be handled during the application run-time
//
// Parameters: _
//
// Returns: _
func OnReady() {
	var menuItems []MenuItem
	// Icon data can be found in the "icon" package
	systray.SetIcon(icon.Data)
	systray.SetTooltip("InfoSec Agent")

	// Generate the menu for the system tray application

	// Example menu item //
	mGoogleBrowser := systray.AddMenuItem(localization.Localize(language, "GoogleTitle"), localization.Localize(language, "GoogleTooltip"))
	menuItems = append(menuItems, MenuItem{menuTitle: "GoogleTitle", menuTooltip: "GoogleTooltip", sysMenuItem: mGoogleBrowser})

	systray.AddSeparator()
	mChangeScanInterval := systray.AddMenuItem(localization.Localize(language, "ScanIntervalTitle"), localization.Localize(language, "ScanIntervalTooltip"))
	menuItems = append(menuItems, MenuItem{menuTitle: "ScanIntervalTitle", menuTooltip: "ScanIntervalTooltip", sysMenuItem: mChangeScanInterval})

	mScanNow := systray.AddMenuItem(localization.Localize(language, "ScanNowTitle"), localization.Localize(language, "ScanNowTooltip"))
	menuItems = append(menuItems, MenuItem{menuTitle: "ScanNowTitle", menuTooltip: "ScanNowTooltip", sysMenuItem: mScanNow})

	systray.AddSeparator()
	mChangeLang := systray.AddMenuItem(localization.Localize(language, "ChangeLangTitle"), localization.Localize(language, "ChangeLangTooltip"))
	menuItems = append(menuItems, MenuItem{menuTitle: "ChangeLangTitle", menuTooltip: "ChangeLangTooltip", sysMenuItem: mChangeLang})

	systray.AddSeparator()
	mQuit := systray.AddMenuItem(localization.Localize(language, "QuitTitle"), localization.Localize(language, "QuitTooltip"))
	menuItems = append(menuItems, MenuItem{menuTitle: "QuitTitle", menuTooltip: "QuitTooltip", sysMenuItem: mQuit})

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
		case <-mChangeLang.ClickedCh:
			changeLang()
			refreshMenu(menuItems)
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

// changeLang provides the user with a dialog window to change the language of the application
//
// Parameters: _
//
// Returns: _
func changeLang() {
	res, err := zenity.List("Choose a language", []string{"German", "British English", "American English",
		"Spanish", "French", "Dutch", "Portuguese"}, zenity.Title("Change Language"),
		zenity.DefaultItems("British English"))
	if err != nil {
		fmt.Println("Error creating dialog:", err)
		return
	}

	// Assign each language to an index for the localization package
	switch res {
	case "German":
		language = 0
	case "British English":
		language = 1
	case "American English":
		language = 2
	case "Spanish":
		language = 3
	case "French":
		language = 4
	case "Dutch":
		language = 5
	case "Portuguese":
		language = 6
	default:
		language = 1
	}
}

// refreshMenu updates the menu items with the current language
//
// Parameters: items ([]MenuItem)
//
// Returns: _
func refreshMenu(items []MenuItem) {
	for _, item := range items {
		item.sysMenuItem.SetTitle(localization.Localize(language, item.menuTitle))
		item.sysMenuItem.SetTooltip(localization.Localize(language, item.menuTooltip))
	}
}
