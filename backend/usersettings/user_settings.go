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
	Language     int `json:"Language"`
	ScanInterval int `json:"ScanInterval"`
}

// NewUserSettings creates a new UserSettings struct
// TODO: Reformat with new doc standard
// Parameters: _
//
// Returns: a pointer to a new UserSettings struct
func NewUserSettings() *UserSettings {
	return &UserSettings{
		Language:     LoadUserSettings("usersettings").Language,
		ScanInterval: LoadUserSettings("usersettings").ScanInterval,
	}
}

// LoadUserSettings loads the user setting from the user settings file
// TODO: Reformat with new doc standard
// Parameters: _
//
// Returns: language setting (int)
func LoadUserSettings(path string) UserSettings {
	data, err := os.ReadFile(path + "/user_settings.json")
	if err != nil {
		logger.Log.ErrorWithErr("Error reading user settings file:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}

	var settings UserSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		logger.Log.ErrorWithErr("Error unmarshalling user settings JSON:", err)
		return UserSettings{Language: 1, ScanInterval: 24}
	}
	return settings
}

// SaveUserSettings saves the user setting(s) to the user settings file
// TODO: Reformat with new doc standard
// Parameters: lang (int) - the language setting to save
//
// Returns: _
func SaveUserSettings(settings UserSettings, path string) {
	file, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		logger.Log.ErrorWithErr("Error marshalling user settings JSON:", err)
		return
	}
	err = os.WriteFile(path+"/user_settings.json", file, 0600)
	if err != nil {
		logger.Log.ErrorWithErr("Error writing user setting(s) to file:", err)
	}
}
