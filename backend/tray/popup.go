// Package tray implements the basic functionality of the system tray application
//
// Exported function(s): OnReady, OnQuit, ChangeScanInterval, ScanNow, ChangeLanguage,
// RefreshMenu
package tray

import (
	"fmt"
	"os"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/database"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/localization"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/go-toast/toast"
)

const iconPath string = "/InfoSec-Agent/icon/icon128.ico"

// Popup displays a notification to the user when a scan is completed.
//
// This function creates a notification with a title, message, and icon to inform the user that a scan has been completed.
// The notification also includes an action button that lets the user open the reporting page.
//
// Parameters: scanResult []checks.Check: A slice of checks representing the scan results.
//
// Returns: error: An error object if an error occurred during the scan, otherwise nil.
func Popup(scanResult []checks.Check, path string) error {
	logger.Log.Trace("Displaying popup for scan result")

	// Generate notification message based on the severity of the issues found during the scan
	resultMessage := PopupMessage(scanResult, path)

	// Get the path to the popup icon
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.ErrorWithErr("error getting icon path: error getting user config dir", err)
	}

	// Create a notification to inform the user that the scan is complete
	notification := toast.Notification{
		AppID:               "InfoSec Agent",
		Title:               localization.Localize(Language, "Dialogs.Popup.Title"),
		Message:             resultMessage,
		Icon:                appDataPath + iconPath,
		ActivationArguments: "infosecagent:",
		Actions: []toast.Action{
			{Type: "protocol", Label: localization.Localize(Language, "Dialogs.Popup.Button"), Arguments: "infosecagent:"},
		},
	}
	if err = notification.Push(); err != nil {
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
	logger.Log.Trace("Generating popup message")

	dbData, err := database.GetData(scanResult, path)
	if err != nil {
		logger.Log.ErrorWithErr("Error getting database data:", err)
		return localization.Localize(Language, "Dialogs.Popup.Default")
	}

	// Count the number of issues at each severity level
	severityCounters := make(map[int]int)
	for _, issue := range dbData {
		severityCounters[issue.Severity]++
	}

	// Generate the notification message based on the number of issues found at each severity level
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

// StartPopup displays a notification to the user when the tray application starts.
//
// This function creates a notification with a title, message, and icon to inform the user that the tray application has started.
// The notification provides information on how to access the application through the system tray.
//
// Parameters: None.
//
// Returns: None.
func StartPopup() {
	// Get the path to the popup icon
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.ErrorWithErr("error getting icon path: error getting user config dir", err)
	}

	// Create a notification to inform the user that the tray application has started
	notification := toast.Notification{
		AppID:   "InfoSec Agent",
		Title:   localization.Localize(Language, "Dialogs.Popup.StartupTitle"),
		Message: localization.Localize(Language, "Dialogs.Popup.Startup"),
		Icon:    appDataPath + iconPath,
	}
	if err = notification.Push(); err != nil {
		logger.Log.ErrorWithErr("error pushing scan notification", err)
	}
}

// AlreadyRunningPopup displays a notification to the user when the tray application is already running.
//
// This function creates a notification with a title, message, and icon to inform the user that the tray application is already running in the background.
// The notification provides information on how to access the application through the system tray.
//
// Parameters: None.
//
// Returns: None.
func AlreadyRunningPopup() {
	// Get the path to the popup icon
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.ErrorWithErr("error getting icon path: error getting user config dir", err)
	}

	// Create a notification to inform the user that the tray application has started
	notification := toast.Notification{
		AppID:   "InfoSec Agent",
		Title:   localization.Localize(Language, "Dialogs.Popup.AlreadyRunningTitle"),
		Message: localization.Localize(Language, "Dialogs.Popup.AlreadyRunning"),
		Icon:    appDataPath + iconPath,
	}
	if err = notification.Push(); err != nil {
		logger.Log.ErrorWithErr("error pushing scan notification", err)
	}
}
