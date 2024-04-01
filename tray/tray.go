// Package tray implements the basic functionality of the system tray application
//
// Exported function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, ChangeLanguage,
// RefreshMenu
package tray

import (
	"log"

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

var scanCounter int
var scanTicker *time.Ticker
var language = 1 // Default language is British English
var menuItems []MenuItem
var rpPage = false
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

	mReportingPage := systray.AddMenuItem(localization.Localize(language, "ReportingPageTitle"), localization.Localize(language, "ReportingPageTooltip"))
	menuItems = append(menuItems, MenuItem{MenuTitle: "ReportingPageTitle", menuTooltip: "ReportingPageTooltip", sysMenuItem: mReportingPage})

	systray.AddSeparator()
	mChangeScanInterval := systray.AddMenuItem(localization.Localize(language, "ScanIntervalTitle"), localization.Localize(language, "ScanIntervalTooltip"))
	menuItems = append(menuItems, MenuItem{MenuTitle: "ScanIntervalTitle", menuTooltip: "ScanIntervalTooltip", sysMenuItem: mChangeScanInterval})

	mScanNow := systray.AddMenuItem(localization.Localize(language, "ScanNowTitle"), localization.Localize(language, "ScanNowTooltip"))
	menuItems = append(menuItems, MenuItem{MenuTitle: "ScanNowTitle", menuTooltip: "ScanNowTooltip", sysMenuItem: mScanNow})

	systray.AddSeparator()
	mChangeLanguage := systray.AddMenuItem(localization.Localize(language, "ChangeLanguageTitle"), localization.Localize(language, "ChangeLanguageTooltip"))
	menuItems = append(menuItems, MenuItem{MenuTitle: "ChangeLanguageTitle", menuTooltip: "ChangeLanguageTooltip", sysMenuItem: mChangeLanguage})

	systray.AddSeparator()
	mQuit = systray.AddMenuItem(localization.Localize(language, "QuitTitle"), localization.Localize(language, "QuitTooltip"))
	menuItems = append(menuItems, MenuItem{MenuTitle: "QuitTitle", menuTooltip: "QuitTooltip", sysMenuItem: mQuit})

	// Set up a channel to receive OS signals, used for termination
	// Can be used to notify the application about system termination signals,
	// allowing it to perform possible cleanup tasks before exiting.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	scanCounter = 0
	// Set a ticker to run a scan at a set interval (default = 1 week)
	scanTicker = time.NewTicker(7 * 24 * time.Hour)

	// Iterate over each menu option/signal
	for {
		select {
		case <-mReportingPage.ClickedCh:
			err := openReportingPage("")
			if err != nil {
				log.Println(err)
			}
		case <-mChangeScanInterval.ClickedCh:
			ChangeScanInterval()
		case <-mScanNow.ClickedCh:
			ScanNow()
		case <-mChangeLanguage.ClickedCh:
			ChangeLanguage()
			RefreshMenu(menuItems)
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		// Executes each time the scanTicker has elapsed the set amount of time
		case <-scanTicker.C:
			scanCounter++
			fmt.Println("Scan:", scanCounter)
			ScanNow()
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

// openReportingPage opens the reporting page using a Wails application
//
// Parameters: _
//
// Returns: _
func openReportingPage(path string) error {
	if rpPage {
		return fmt.Errorf("reporting-page is already running")
	}

	// Get the current working directory
	//TODO: In a release version, there (should be) no need to build the application, just run it
	//Consideration: Wails can also send (termination) signals to the back-end, might be worth investigating
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current directory: %s", err)
	}

	// Change directory to reporting-page folder
	err = os.Chdir(path + "reporting-page")
	if err != nil {
		return fmt.Errorf("Error changing directory: %s", err)
	}

	// Restore the original working directory
	defer func() {
		err := os.Chdir(originalDir)
		if err != nil {
			log.Println("Error changing directory:", err)
		}
		rpPage = false
	}()

	// Build reporting-page executable
	buildCmd := exec.Command("wails", "build", "-windowsconsole")

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Error building reporting-page: %s", err)
	}

	// Set up the reporting-page executable
	runCmd := exec.Command("build/bin/reporting-page")
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	// Set up a listener for the quit function from the system tray
	go func() {
		<-mQuit.ClickedCh
		if err := runCmd.Process.Kill(); err != nil {
			log.Println("Error interrupting reporting-page process:", err)
		}
		rpPage = false
		systray.Quit()
	}()

	rpPage = true
	// Run the reporting page executable
	if err := runCmd.Run(); err != nil {
		rpPage = false
		return fmt.Errorf("Error running reporting-page: %s", err)
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
		res, err = zenity.Entry("Enter the scan interval (in hours):", zenity.Title("Change Scan Interval"), zenity.DefaultItems("24"))
		if err != nil {
			log.Println("Error creating dialog:", err)
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
func ScanNow() ([]checks.Check, error) {
	// scanCounter is not concretely used at the moment
	// might be useful in the future
	scanCounter++
	fmt.Println("Scanning now. Scan:", scanCounter)

	// Display a progress dialog while the scan is running
	dialog, err := zenity.Progress(
		zenity.Title("Security/Privacy Scan"))
	if err != nil {
		log.Println("Error creating dialog:", err)
		return nil, err
	}
	// Defer closing the dialog until the scan completes
	defer func(dialog zenity.ProgressDialog) {
		err := dialog.Close()
		if err != nil {
			log.Println("Error closing dialog:", err)
		}
	}(dialog)

	result, err := scan.Scan(dialog)
	if err != nil {
		log.Println("Error calling scan:", err)
		return result, err
	}

	err = dialog.Complete()
	if err != nil {
		log.Println("Error completing dialog:", err)
		return result, err
	}

	// Fill database with data
	// scan.FillDataBase(result)

	// zeros := make([]int, len(result))
	// zeros[3] = 2
	// sev, _ := scan.GetAllSeverities(result, zeros)
	// fmt.Println("severities", sev)

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
			log.Println("Error creating dialog:", err)
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
// Parameters: items ([]MenuItem)
//
// Returns: _
func RefreshMenu(items []MenuItem) {
	for _, item := range items {
		item.sysMenuItem.SetTitle(localization.Localize(language, item.MenuTitle))
		item.sysMenuItem.SetTooltip(localization.Localize(language, item.menuTooltip))
	}
}
