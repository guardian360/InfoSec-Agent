// Package usersettings contains functions for loading and saving user settings
//
// Exported function(s): NewUserSettings, LoadUserSettings, SaveUserSettings
//
// Exported type(s): UserSettings
package usersettings

import (
	"encoding/json"
	"os"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
)

type UserSettings struct {
	Language     int  `json:"Language"`
	ScanInterval int  `json:"ScanInterval"`
	Integration  bool `json:"Integration"`
}

// LoadUserSettings loads the user settings from a JSON file in the Windows AppData folder.
//
// The function uses the os package to get the path to the AppData folder, and reads the user settings from a file named "user_settings.json" in this folder.
// If there is an error while getting the AppData folder path, reading the JSON data from the file, or unmarshalling the JSON data to a UserSettings struct,
//
// Parameters: None
//
// Returns:
//   - settings (UserSettings): The loaded user settings. This is a UserSettings struct.
func LoadUserSettings() UserSettings {
	logger.Log.Debug("Getting user config directory")
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.ErrorWithErr("Error getting user config directory:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}
	dirPath := appDataPath + `\InfoSec-Agent`
	logger.Log.Debug("Creating/reading directory at:" + dirPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		logger.Log.ErrorWithErr("Error creating directory:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}

	filePath := dirPath + `\user_settings.json`

	logger.Log.Debug("Reading user settings from file:" + filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Log.ErrorWithErr("Error reading user settings file:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}

	var settings UserSettings
	logger.Log.Debug("Unmarshalling user settings JSON")
	err = json.Unmarshal(data, &settings)
	if err != nil {
		logger.Log.ErrorWithErr("Error unmarshalling user settings JSON:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}
	return settings
}

// SaveUserSettings saves the user settings to a JSON file in the Windows AppData\Roaming folder.
//
// The function takes a UserSettings struct as input, which contains the user settings to be saved.
// It uses the os package to get the path to the AppData folder, and saves the user settings to a file named "user_settings.json" in this folder.
//
// Parameters:
//   - settings (UserSettings): The user settings to be saved.
//
// Returns: None
func SaveUserSettings(settings UserSettings) {
	logger.Log.Debug("Getting user config directory")
	appDataPath, err := os.UserConfigDir()
	if err != nil {
		logger.Log.ErrorWithErr("Error getting user config directory:", err)
		return
	}

	dirPath := appDataPath + `\InfoSec-Agent`
	logger.Log.Debug("Creating/reading directory at:" + dirPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		logger.Log.ErrorWithErr("Error creating directory:", err)
		return
	}
	filePath := dirPath + `\user_settings.json`

	logger.Log.Debug("Marshalling data to JSON")
	file, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		logger.Log.ErrorWithErr("Error marshalling user settings JSON:", err)
		return
	}
	logger.Log.Debug("Writing user settings to file:" + filePath)
	err = os.WriteFile(filePath, file, 0600)
	if err != nil {
		logger.Log.ErrorWithErr("Error writing user setting(s) to file:", err)
	}
}
