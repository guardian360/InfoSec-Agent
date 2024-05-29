// Package tray implements the basic functionality of the system tray application
//
// Exported function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, ChangeLanguage,
// RefreshMenu
package tray

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/usersettings"
	"github.com/go-toast/toast"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/icon"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/scan"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var ScanCounter int

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
// This function sets up the system tray with various menu items such as 'Reporting Page', 'Change Scan Interval', 'Scan Now', 'Change Language', and 'Quit'.
// It then enters a loop where it listens for various events such as clicks on the menu items, system termination signals, and elapse of the scan interval. Depending on the event, it performs actions such as opening the reporting page, changing the scan interval, initiating an immediate scan, changing the application language, refreshing the menu, or quitting the application.
//
// Parameters: None.
//
// Returns: None. The function runs indefinitely, orchestrating the behavior of the system tray application.
func OnReady() {
	// Icon data can be found in the "icon" package
	systray.SetIcon(icon.Data)
	systray.SetTooltip("InfoSec Agent")

	settings := usersettings.LoadUserSettings()
	Language = settings.Language
	scanInterval := settings.ScanInterval

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
	ticker := time.NewTicker(30 * time.Minute)

	// Iterate over each menu option/signal
	for {
		select {
		case <-mReportingPage.ClickedCh:
			err := OpenReportingPage("")
			if err != nil {
				logger.Log.ErrorWithErr("Error opening reporting page:", err)
			}
		case <-mChangeScanInterval.ClickedCh:
			ChangeScanInterval()
		case <-mScanNow.ClickedCh:
			result, err := ScanNow(true)
			if err != nil {
				logger.Log.ErrorWithErr("Error scanning:", err)
			} else {
				// Notify the user that a scan has been completed
				err = Popup(result, "./reporting-page/database.db")
				if err != nil {
					logger.Log.ErrorWithErr("Error notifying user:", err)
				}
			}
		case <-mChangeLanguage.ClickedCh:
			ChangeLanguage()
			RefreshMenu()
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		case <-ticker.C:
			periodicScan(scanInterval)
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

	logger.Log.Debug("opening reporting page")

	// TODO: use build tags to differentiate between development and release versions
	// Get the current working directory
	// Consideration: Wails can also send (termination) signals to the back-end, might be worth investigating
	originalDir, err := os.Getwd()
	logger.Log.Debug("current directory: " + originalDir)
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
			logger.Log.ErrorWithErr("error changing directory:", err)
		}
		ReportingPageOpen = false
	}()

	//TODO: In a release version, there (should be) no need to build the application, just run it
	const build = false
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
			logger.Log.ErrorWithErr("error interrupting reporting-page process:", err)
		}
		ReportingPageOpen = false
		systray.Quit()
	}()

	// Run the reporting page executable
	ReportingPageOpen = true
	if err = runCmd.Run(); err != nil {
		ReportingPageOpen = false
		return fmt.Errorf("error running reporting-page: %w", err)
	}

	logger.Log.Debug("reporting page opened")
	return nil
}

// BuildReportingPage builds the reporting page executable using a Wails application
//
// Parameters:
//   - path string: The relative path to the reporting-page directory. This is used to change the current working directory to the reporting-page directory.
//
// Returns:
//   - error: An error object if an error occurred during the process, otherwise nil.
func BuildReportingPage() error {
	// Build reporting page
	buildCmd := exec.Command("wails", "build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("error building reporting-page: %w", err)
	}
	logger.Log.Debug("reporting page built successfully")
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
// Returns: None.
func ChangeScanInterval(testInput ...string) {
	var res string
	test := len(testInput) > 0
	// If testInput is provided, use it for testing
	if test {
		res = testInput[0]
	} else {
		scanInterval := usersettings.LoadUserSettings().ScanInterval

		// Get user input by creating a dialog window
		var err error
		res, err = zenity.Entry(localization.Localize(Language, "Dialogs.ScanInterval.Content"),
			zenity.Title(localization.Localize(Language, "Dialogs.ScanInterval.Title")),
			zenity.EntryText(strconv.Itoa(scanInterval)),
			zenity.DefaultItems("24"))
		if err != nil {
			logger.Log.ErrorWithErr("Error creating dialog:", err)
			return
		}
	}

	// Parse the user input
	interval, err := strconv.Atoi(res)
	scanInterval := usersettings.LoadUserSettings().ScanInterval
	if err != nil || interval <= 0 {
		if !test {
			err = zenity.Info("Invalid input. Using previous interval of "+strconv.Itoa(scanInterval)+" hours.",
				zenity.Title("Invalid scan interval input"))
			if err != nil {
				logger.Log.ErrorWithErr("Error creating invalid interval confirmation dialog:", err)
			}
		}
		interval = scanInterval
		updateScanInterval(interval, test)
		return
	}
	if !test {
		err = zenity.Info("Scan interval changed to "+strconv.Itoa(interval)+" hours",
			zenity.Title("Scan Interval Changed"))
		if err != nil {
			logger.Log.ErrorWithErr("Error creating interval confirmation dialog:", err)
		}
	}
	updateScanInterval(interval, test)
}

// ScanNow initiates an immediate security scan, bypassing the scheduled intervals.
//
// This function triggers a security scan regardless of the scheduled intervals. It is useful for situations where an immediate scan is required, such as after a significant system change or when manually requested by the user.
// During the scan, a progress dialog is displayed to keep the user informed about the scan progress. Once the scan is complete, the dialog is closed and the results of the scan are returned.
//
// Parameters:
//   - dialogPresent bool: A boolean value indicating whether a progress dialog should be displayed during the scan. If true, a dialog is shown; if false, no dialog is displayed.
//
// Returns:
//   - []checks.Check: A list of checks performed during the scan.
//   - error: An error object if an error occurred during the scan, otherwise nil.
func ScanNow(dialogPresent bool) ([]checks.Check, error) {
	// ScanCounter is not concretely used at the moment
	// might be useful in the future
	ScanCounter++
	logger.Log.Info("Scanning now. Scan:" + strconv.Itoa(ScanCounter))

	var result []checks.Check
	var err error
	var dialog zenity.ProgressDialog
	if dialogPresent {
		dialog, result, err = runScanWithDialog()
		if err != nil {
			logger.Log.ErrorWithErr("Error running scan with dialog:", err)
			return result, err
		}
		// Defer closing the dialog until the scan completes
		defer func(dialog zenity.ProgressDialog) {
			err = dialog.Close()
			if err != nil {
				logger.Log.ErrorWithErr("Error closing dialog:", err)
			}
		}(dialog)
	} else {
		result, err = scan.Scan(nil, Language)
		if err != nil {
			logger.Log.ErrorWithErr("Error calling scan:", err)
			return result, err
		}
	}
	/*// Uncomment for points printing

	//Temporary dummy game state. For future changed to the current saved game state.
	//gsDummy := gamification.GameState{Points: 0, PointsHistory: nil, LighthouseState: 0}
	gsDummy := gamification.GameState{Points: 0, PointsHistory: []int{}, LighthouseState: 0}

	//Calculate points based on the scan results
	gs, err := gamification.PointCalculation(gsDummy, result, "reporting-page/database.db")
	if err != nil {
		logger.Log.ErrorWithErr("Error calculating points:", err)
		return result, err
	}

	fmt.Print(gs)*/

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
//   - testInput ...string: Optional parameter used for testing. If provided, the function uses this as the user's language selection instead of displaying the dialog window.
//
// Returns: None. The function updates the 'language' variable in-place.
func ChangeLanguage(testInput ...string) {
	var res string
	test := testInput != nil
	if test {
		res = testInput[0]
	} else {
		var err error
		res, err = zenity.List(localization.Localize(Language, "Dialogs.Language.Content"), []string{"German", "British English", "American English",
			"Spanish", "French", "Dutch", "Portuguese"}, zenity.Title(localization.Localize(Language, "Dialogs.Language.Title")),
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

	if test {
		return
	}
	err := usersettings.SaveUserSettings(usersettings.UserSettings{
		Language:     Language,
		ScanInterval: usersettings.LoadUserSettings().ScanInterval,
	})
	if err != nil {
		logger.Log.Warning("Language setting not saved to file")
	}
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

// Popup displays a notification to the user when a scan is completed.
//
// This function creates a notification with a title, message, and icon to inform the user that a scan has been completed.
// The notification also includes an action button that lets the user open the reporting page.
//
// Parameters: scanResult []checks.Check: A slice of checks representing the scan results.
//
// Returns: error: An error object if an error occurred during the scan, otherwise nil.
func Popup(scanResult []checks.Check, path string) error {
	// Generate notification message based on the severity of the issues found during the scan
	resultMessage := PopupMessage(scanResult, path)

	// Create a notification to inform the user that the scan is complete
	notification := toast.Notification{
		AppID:   "InfoSec Agent",
		Title:   localization.Localize(Language, "Dialogs.Popup.Title"),
		Message: resultMessage,
		// Icon:    "",
		ActivationArguments: "infosecagent:",
		Actions: []toast.Action{
			{Type: "protocol", Label: localization.Localize(Language, "Dialogs.Popup.Button"), Arguments: "infosecagent:"},
		},
	}
	if err := notification.Push(); err != nil {
		return fmt.Errorf("error pushing scan notification: %w", err)
	}
	return nil
}

// PopupMessage generates a notification message based on the severity of the issues found during the scan.
//
// This function takes a slice of checks representing the scan results and generates a notification message based on the number of issues found at each severity level.
// The message informs the user about the number of issues found during the scan and prompts them to open the reporting page for more information.
//
// Parameters: scanResult []checks.Check: A slice of checks representing the scan results.
//
// Returns: string: A notification message based on the severity of the issues found during the scan.
func PopupMessage(scanResult []checks.Check, path string) string {
	dbData, err := database.GetData(scanResult, path)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting database data:", err)
	}
	severityCounters := make(map[int]int)
	for _, issue := range dbData {
		severityCounters[issue.Severity]++
	}
	if severityCounters[3] > 0 {
		if severityCounters[3] == 1 {
			return localization.Localize(Language, "Dialogs.Popup.OneHigh")
		}
		return fmt.Sprintf(localization.Localize(Language, "Dialogs.Popup.MultipleHigh"), severityCounters[3])
	} else if severityCounters[2] > 0 {
		if severityCounters[2] == 1 {
			return localization.Localize(Language, "Dialogs.Popup.OneMedium")
		}
		return fmt.Sprintf(localization.Localize(Language, "Dialogs.Popup.MultipleMedium"), severityCounters[2])
	}
	return localization.Localize(Language, "Dialogs.Popup.Default")
}

// changeNextScan updates the next scan time based on the current time and the scan interval.
//
// Parameters:
//   - settings usersettings.UserSettings: The user settings object containing the current scan interval and next scan time.
//   - value int: The new scan interval in hours.
//
// Returns: None.
func changeNextScan(settings usersettings.UserSettings, value int) {
	settings.NextScan = time.Now().Add(time.Duration(value) * time.Hour)
	err := usersettings.SaveUserSettings(settings)
	if err != nil {
		logger.Log.Warning("Next scan time not saved to file")
	}
}

// periodicScan checks if a scan is due based on the scan interval and the current time.
// If a scan is due, it performs a scan and notifies the user using a pop-up.
//
// Parameters:
//   - scanInterval int: The scan interval in hours.
//
// Returns: None.
func periodicScan(scanInterval int) {
	settings := usersettings.LoadUserSettings()
	if time.Now().After(settings.NextScan) {
		result, err := ScanNow(false)
		if err != nil {
			logger.Log.ErrorWithErr("Error performing periodic scan:", err)
		} else {
			// Notify the user that a scan has been completed
			err = Popup(result, "./reporting-page/database.db")
			if err != nil {
				logger.Log.ErrorWithErr("Error notifying user:", err)
			}
		}
		// Update the next scan time
		changeNextScan(settings, scanInterval)
	}
}

// runScanWithDialog runs a scan with a progress dialog to keep the user informed about the scan progress.
// It returns the progress dialog, the scan results, and any error that occurred during the scan.
//
// Parameters: None.
//
// Returns:
//   - zenity.ProgressDialog: A progress dialog that displays the scan progress to the user.
//   - []checks.Check: A slice of checks representing the scan results.
//   - error: An error object that describes the error (if any) that occurred during the scan.
func runScanWithDialog() (zenity.ProgressDialog, []checks.Check, error) {
	dialog, err := zenity.Progress(
		zenity.Title(localization.Localize(Language, "Dialogs.Scan.Title")))
	if err != nil {
		logger.Log.ErrorWithErr("Error creating dialog:", err)
	}
	result, err := scan.Scan(dialog, Language)
	if err != nil {
		logger.Log.ErrorWithErr("Error calling scan:", err)
		return dialog, result, err
	}

	err = dialog.Complete()
	if err != nil {
		logger.Log.ErrorWithErr("Error completing dialog:", err)
		return dialog, result, err
	}
	return dialog, result, err
}

// updateScanInterval updates the scan interval in the user settings file.
//
// Parameters:
//   - interval int: The new scan interval in hours.
//
// Returns: None.
func updateScanInterval(interval int, test bool) {
	logger.Log.Printf("INFO: Scan interval changed to " + strconv.Itoa(interval) + " hours")
	if !test {
		err := usersettings.SaveUserSettings(usersettings.UserSettings{
			Language:     usersettings.LoadUserSettings().Language,
			ScanInterval: interval,
			NextScan:     time.Now().Add(time.Duration(interval) * time.Hour),
		})
		if err != nil {
			logger.Log.Warning("Scan interval setting not saved to file")
		}
	}
}
