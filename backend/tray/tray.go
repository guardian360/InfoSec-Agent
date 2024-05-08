// Package tray implements the basic functionality of the system tray application
//
// Exported function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, ChangeLanguage,
// RefreshMenu
package tray

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/pkg/errors"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/icon"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"

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

// Language is used to represent the index of the currently selected language.
// The language indices are as follows:
//
// 0: German
//
// 1: British English
//
// 2: American English
//
// 3: Spanish
//
// 4: French
//
// 5: Dutch
//
// 6: Portuguese
//
// Default language is British English
var Language = 1

var MenuItems []MenuItem
var ReportingPageOpen = false
var mQuit *systray.MenuItem

// MenuItem represents a single item in the system tray menu.
//
// This struct encapsulates the title, tooltip text, and the actual system tray menu item object for a single menu item.
// The 'MenuTitle' field is a string that represents the title of the menu item. This is the text that is displayed in the system tray menu.
// The 'menuTooltip' field is a string that represents the tooltip text for the menu item. This is the text that is displayed when the user hovers over the menu item in the system tray menu.
// The 'sysMenuItem' field is a pointer to a systray.MenuItem object. This is the actual menu item object that is added to the system tray menu.
//
// Fields:
//   - MenuTitle string: The title of the menu item. This is the text that is displayed in the system tray menu.
//   - menuTooltip string: The tooltip text for the menu item. This is the text that is displayed when the user hovers over the menu item in the system tray menu.
//   - sysMenuItem *systray.MenuItem: The actual menu item object that is added to the system tray menu.
type MenuItem struct {
	MenuTitle   string
	menuTooltip string
	sysMenuItem *systray.MenuItem
}

// OnReady orchestrates the runtime behavior of the system tray application.
//
// This function sets up the system tray with various menu items such as 'Reporting Page', 'Change Scan Interval', 'Scan Now', 'Change Language', and 'Quit'. It also initializes a ticker for scheduled security scans and a signal listener for system termination signals.
// It then enters a loop where it listens for various events such as clicks on the menu items, system termination signals, and elapse of the scan interval. Depending on the event, it performs actions such as opening the reporting page, changing the scan interval, initiating an immediate scan, changing the application language, refreshing the menu, or quitting the application.
//
// Parameters: None.
//
// Returns: None. The function runs indefinitely, orchestrating the behavior of the system tray application.
func OnReady() {
	// Icon data can be found in the "icon" package
	systray.SetIcon(icon.Data)
	systray.SetTooltip("InfoSec Agent")

	Language = usersettings.LoadUserSettings("backend/usersettings").Language
	scanInterval := usersettings.LoadUserSettings("backend/usersettings").ScanInterval

	// Generate the menu for the system tray application
	mReportingPage := systray.AddMenuItem(localization.Localize(Language, "Tray.ReportingPageTitle"),
		localization.Localize(Language, "Tray.ReportingPageTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ReportingPageTitle",
		menuTooltip: "Tray.ReportingPageTooltip", sysMenuItem: mReportingPage})

	systray.AddSeparator()
	mChangeScanInterval := systray.AddMenuItem(localization.Localize(Language, "Tray.ScanIntervalTitle"),
		localization.Localize(Language, "Tray.ScanIntervalTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ScanIntervalTitle",
		menuTooltip: "Tray.ScanIntervalTooltip", sysMenuItem: mChangeScanInterval})

	mScanNow := systray.AddMenuItem(localization.Localize(Language, "Tray.ScanNowTitle"),
		localization.Localize(Language, "Tray.ScanNowTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ScanNowTitle",
		menuTooltip: "Tray.ScanNowTooltip", sysMenuItem: mScanNow})

	systray.AddSeparator()
	mChangeLanguage := systray.AddMenuItem(localization.Localize(Language, "Tray.ChangeLanguageTitle"),
		localization.Localize(Language, "Tray.ChangeLanguageTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.ChangeLanguageTitle",
		menuTooltip: "Tray.ChangeLanguageTooltip", sysMenuItem: mChangeLanguage})

	systray.AddSeparator()
	mQuit = systray.AddMenuItem(localization.Localize(Language, "Tray.QuitTitle"),
		localization.Localize(Language, "Tray.QuitTooltip"))
	MenuItems = append(MenuItems, MenuItem{MenuTitle: "Tray.QuitTitle",
		menuTooltip: "Tray.QuitTooltip", sysMenuItem: mQuit})

	// Set up a channel to receive OS signals, used for termination
	// Can be used to notify the application about system termination signals,
	// allowing it to perform possible cleanup tasks before exiting.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	ScanCounter = 0
	// Set a ticker to run a scan at a set interval (default = 1 week)
	ScanTicker = time.NewTicker(time.Duration(scanInterval) * time.Hour)

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
			ChangeLanguage("usersettings")
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

// OnQuit manages the cleanup operations that need to be performed when the application is about to terminate.
//
// This function is called when the application is exiting. It is responsible for performing any necessary cleanup operations such as closing open files, terminating active connections, or releasing resources. The specific cleanup operations depend on the resources and services used by the application.
//
// Parameters: None.
//
// Returns: None. The function performs cleanup operations in-place.
func OnQuit() {
	// Perform cleanup tasks here
	// Currently, there are no cleanup tasks to perform
	logger.Log.Info("Quitting the application")
}

// OpenReportingPage launches the reporting page of the application using a Wails application.
//
// This function checks if a reporting page is already open. If it is, it returns an error. If not, it changes the current working directory to the reporting page directory and builds the reporting-page executable using the Wails framework.
// It then runs the executable, opening the reporting page. If the 'Quit' option is selected from the system tray while the reporting page is open, the function kills the reporting-page process and sets the ReportingPageOpen flag to false.
//
// Parameters:
//   - path string: The relative path to the reporting-page directory. This is used to change the current working directory to the reporting-page directory.
//
// Returns:
//   - error: An error object if an error occurred during the process, otherwise nil.
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

	const build = true
	if build {
		err = BuildReportingPage()
		if err != nil {
			return err
		}
	}

	// Set up the reporting-page executable
	runCmd := exec.Command("build/bin/InfoSec-Agent-Reporting-Page")
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

// BuildReportingPage builds the reporting page executable using a Wails application
//
// Parameters: _
//
// Returns: _
func BuildReportingPage() error {
	buildCmd := exec.Command("wails", "build", "-windowsconsole")

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("error building reporting-page: %w", err)
	}
	return nil
}

// ChangeScanInterval prompts the user to set a new scan interval through a dialog window.
//
// This function displays a dialog window asking the user to input the desired scan interval in hours. If the user input is valid, the function updates the scan interval accordingly. If the input is invalid or less than or equal to zero, the function defaults to a 24-hour interval.
//
// For testing purposes, an optional string parameter 'testInput' can be provided. If 'testInput' is provided, the function uses this as the user's input instead of displaying the dialog window.
//
// Parameters:
//   - testInput ...string: Optional parameter used for testing. If provided, the function uses this as the user's input instead of displaying the dialog window.
//
// Returns: None. The function updates the 'ScanTicker' variable in-place.
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
	if ScanTicker != nil {
		ScanTicker.Stop()
	}
	ScanTicker = time.NewTicker(time.Duration(interval) * time.Hour)
	logger.Log.Printf("Scan interval changed to %d hours\n", interval)
	usersettings.SaveUserSettings(usersettings.UserSettings{
		Language:     usersettings.LoadUserSettings("backend/usersettings").Language,
		ScanInterval: interval,
	}, "backend/usersettings")
}

// ScanNow initiates an immediate security scan, bypassing the scheduled intervals.
//
// This function triggers a security scan regardless of the scheduled intervals. It is useful for situations where an immediate scan is required, such as after a significant system change or when manually requested by the user.
// During the scan, a progress dialog is displayed to keep the user informed about the scan progress. Once the scan is complete, the dialog is closed and the results of the scan are returned.
//
// Parameters: None.
//
// Returns:
//   - []checks.Check: A list of checks performed during the scan.
//   - error: An error object if an error occurred during the scan, otherwise nil.
func ScanNow() ([]checks.Check, error) {
	// ScanCounter is not concretely used at the moment
	// might be useful in the future
	ScanCounter++
	logger.Log.Info("Scanning now. Scan:" + strconv.Itoa(ScanCounter))

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

// ChangeLanguage allows the user to select a new language for the application via a dialog window.
//
// This function presents a dialog window with a list of available languages. The user can select a language from this list, and the application's language setting is updated accordingly.
// The function maps each language to an index, which is used internally for localization. If the function is called with a test input, it uses the test input instead of displaying the dialog window.
//
// The language indices are as follows:
// 0: German
// 1: British English
// 2: American English
// 3: Spanish
// 4: French
// 5: Dutch
// 6: Portuguese
//
// Parameters:
//
//   - path string: The relative path to the user settings file. This is used to save the updated language setting.
//   - testInput ...string: Optional parameter used for testing. If provided, the function uses this as the user's language selection instead of displaying the dialog window.
//
// Returns: None. The function updates the 'language' variable in-place.
func ChangeLanguage(path string, testInput ...string) {
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
		Language = 0
	case "British English":
		Language = 1
	case "American English":
		Language = 2
	case "Spanish":
		Language = 3
	case "French":
		Language = 4
	case "Dutch":
		Language = 5
	case "Portuguese":
		Language = 6
	default:
		Language = 1
	}
	usersettings.SaveUserSettings(usersettings.UserSettings{
		Language:     Language,
		ScanInterval: usersettings.LoadUserSettings(path).ScanInterval,
	}, path)
}

// RefreshMenu updates the system tray menu items to reflect the current language setting.
//
// This function iterates over each menu item in the system tray and updates its title and tooltip text to match the current language setting.
// The language setting is determined by the 'language' variable, which stores the index of the currently active language.
// The function uses the 'Localize' function from the 'localization' package to translate the title and tooltip text of each menu item.
//
// Parameters: None.
//
// Returns: None. The function updates the system tray menu items in-place.
func RefreshMenu() {
	for _, item := range MenuItems {
		item.sysMenuItem.SetTitle(localization.Localize(Language, item.MenuTitle))
		item.sysMenuItem.SetTooltip(localization.Localize(Language, item.menuTooltip))
	}
}
