// Package tray implements the basic functionality of the system tray application
//
// Exported function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, ChangeLanguage,
// RefreshMenu
package tray

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"github.com/pkg/errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/icon"
	"github.com/InfoSec-Agent/InfoSec-Agent/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/scan"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"

	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var ScanCounter int
var ScanTicker *time.Ticker
var language = 1 // Default language is British English
var MenuItems []MenuItem
var ReportingPageOpen = false
var mQuit *systray.MenuItem

type MenuItem struct {
	MenuTitle   string
	menuTooltip string
	sysMenuItem *systray.MenuItem
}

// OnReady handles all actions that should be handled during the application run-time
//
// Parameters: _
//
// Returns: _
func OnReady() {
	// Icon data can be found in the "icon" package
	systray.SetIcon(icon.Data)
	systray.SetTooltip("InfoSec Agent")

	// Generate the menu for the system tray application

	mReportingPage := systray.AddMenuItem(localization.Localize(language, "Tray.ReportingPageTitle"),
		localization.Localize(language, "Tray.ReportingPageTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ReportingPageTitle",
		menuTooltip: "Tray.ReportingPageTooltip", sysMenuItem: mReportingPage})

	systray.AddSeparator()
	mChangeScanInterval := systray.AddMenuItem(localization.Localize(language, "Tray.ScanIntervalTitle"),
		localization.Localize(language, "Tray.ScanIntervalTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ScanIntervalTitle",
		menuTooltip: "Tray.ScanIntervalTooltip", sysMenuItem: mChangeScanInterval})

	mScanNow := systray.AddMenuItem(localization.Localize(language, "Tray.ScanNowTitle"),
		localization.Localize(language, "Tray.ScanNowTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ScanNowTitle",
		menuTooltip: "Tray.ScanNowTooltip", sysMenuItem: mScanNow})

	systray.AddSeparator()
	mChangeLanguage := systray.AddMenuItem(localization.Localize(language, "Tray.ChangeLanguageTitle"),
		localization.Localize(language, "Tray.ChangeLanguageTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ChangeLanguageTitle",
		menuTooltip: "Tray.ChangeLanguageTooltip", sysMenuItem: mChangeLanguage})

	systray.AddSeparator()
	mQuit = systray.AddMenuItem(localization.Localize(language, "Tray.QuitTitle"),
		localization.Localize(language, "Tray.QuitTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.QuitTitle",
		menuTooltip: "Tray.QuitTooltip", sysMenuItem: mQuit})

	// Set up a channel to receive OS signals, used for termination
	// Can be used to notify the application about system termination signals,
	// allowing it to perform possible cleanup tasks before exiting.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	ScanCounter = 0
	// Set a ticker to run a scan at a set interval (default = 1 week)
	ScanTicker = time.NewTicker(7 * 24 * time.Hour)

	// Iterate over each menu option/signal
	for {
		select {
		case <-mReportingPage.ClickedCh:
			err := OpenReportingPage("")
			if err != nil {
				logger.Log.Println(err)
			}
		case <-mChangeScanInterval.ClickedCh:
			ChangeScanInterval()
		case <-mScanNow.ClickedCh:
			_, err := ScanNow()
			if err != nil {
				logger.Log.ErrorWithErr("Error scanning:", err)
			}
		case <-mChangeLanguage.ClickedCh:
			ChangeLanguage()
			RefreshMenu()
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		// Executes each time the ScanTicker has elapsed the set amount of time
		case <-ScanTicker.C:
			ScanCounter++
			logger.Log.Println("Scan:", ScanCounter)
			_, err := ScanNow()
			if err != nil {
				logger.Log.ErrorWithErr("Error scanning:", err)
			}
		}
	}
}

// OnQuit handles all actions that should happen when the application exits/terminates
//
// Parameters: _
//
// Returns: _
func OnQuit() {
	// Perform cleanup tasks here
	// Currently, there are no cleanup tasks to perform
	logger.Log.Info("Quitting the application")
}

// OpenReportingPage opens the reporting page using a Wails application
//
// Parameters: _
//
// Returns: _
func OpenReportingPage(path string) error {
	if ReportingPageOpen {
		return errors.New("reporting-page is already running")
	}

	// Get the current working directory
	//TODO: In a release version, there (should be) no need to build the application, just run it
	//Consideration: Wails can also send (termination) signals to the back-end, might be worth investigating
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Change directory to reporting-page folder
	err = os.Chdir(path + "reporting-page")
	if err != nil {
		return fmt.Errorf("error changing directory: %w", err)
	}

	// Restore the original working directory
	defer func() {
		err = os.Chdir(originalDir)
		if err != nil {
			logger.Log.ErrorWithErr("Error changing directory:", err)
		}
		ReportingPageOpen = false
	}()

	// Build reporting-page executable
	buildCmd := exec.Command("wails", "build", "-windowsconsole")

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err = buildCmd.Run(); err != nil {
		return fmt.Errorf("error building reporting-page: %w", err)
	}

	// Set up the reporting-page executable
	runCmd := exec.Command("build/bin/reporting-page")
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	// Set up a listener for the quit function from the system tray
	go func() {
		<-mQuit.ClickedCh
		if err = runCmd.Process.Kill(); err != nil {
			logger.Log.ErrorWithErr("Error interrupting reporting-page process:", err)
		}
		ReportingPageOpen = false
		systray.Quit()
	}()

	ReportingPageOpen = true
	// Run the reporting page executable
	if err = runCmd.Run(); err != nil {
		ReportingPageOpen = false
		return fmt.Errorf("error running reporting-page: %w", err)
	}
	return nil
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
		res, err = zenity.Entry("Enter the scan interval (in hours):", zenity.Title("Change Scan Interval"),
			zenity.DefaultItems("24"))
		if err != nil {
			logger.Log.ErrorWithErr("Error creating dialog:", err)
			return
		}
	}

	// Parse the user input
	interval, err := strconv.Atoi(res)
	if err != nil || interval <= 0 {
		logger.Log.Printf("Invalid input. Using default interval of 24 hours.")
		interval = 24
	}

	// Restart the ticker with the new interval
	ScanTicker.Stop()
	ScanTicker = time.NewTicker(time.Duration(interval) * time.Hour)
	logger.Log.Printf("Scan interval changed to %d hours\n", interval)
}

// ScanNow performs one scan iteration (without checking if it is scheduled)
//
// Parameters: _
//
// Returns: list of checks
func ScanNow() ([]checks.Check, error) {
	// ScanCounter is not concretely used at the moment
	// might be useful in the future
	ScanCounter++
	logger.Log.Info("Scanning now. Scan:" + string(rune(ScanCounter)))

	// Display a progress dialog while the scan is running
	dialog, err := zenity.Progress(
		zenity.Title("Security/Privacy Scan"))
	if err != nil {
		logger.Log.ErrorWithErr("Error creating dialog:", err)
		return nil, err
	}
	// Defer closing the dialog until the scan completes
	defer func(dialog zenity.ProgressDialog) {
		err = dialog.Close()
		if err != nil {
			logger.Log.ErrorWithErr("Error closing dialog:", err)
		}
	}(dialog)

	result, err := scan.Scan(dialog)
	if err != nil {
		logger.Log.ErrorWithErr("Error calling scan:", err)
		return result, err
	}

	err = dialog.Complete()
	if err != nil {
		logger.Log.ErrorWithErr("Error completing dialog:", err)
		return result, err
	}

	return result, nil
}

// ChangeLanguage provides the user with a dialog window to change the language of the application
//
// Parameters: _
//
// Returns: _
func ChangeLanguage(testInput ...string) {
	var res string
	if len(testInput) > 0 {
		res = testInput[0]
	} else {
		var err error
		res, err = zenity.List("Choose a language", []string{"German", "British English", "American English",
			"Spanish", "French", "Dutch", "Portuguese"}, zenity.Title("Change Language"),
			zenity.DefaultItems("British English"))
		if err != nil {
			logger.Log.ErrorWithErr("Error creating dialog:", err)
			return
		}
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

// RefreshMenu updates the menu items with the current language
//
// Parameters: _
//
// Returns: _
func RefreshMenu() {
	for _, item := range MenuItems {
		item.sysMenuItem.SetTitle(localization.Localize(language, item.MenuTitle))
		item.sysMenuItem.SetTooltip(localization.Localize(language, item.menuTooltip))
	}
}

// Language returns the current language index
//
// Parameters: _
//
// Returns: language index
func Language() int {
	return language
}
