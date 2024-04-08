package usersettings

import (
	"encoding/json"
	"github.com/InfoSec-Agent/InfoSec-Agent/logger"
	"os"
)

type UserSettings struct {
	Language     int `json:"Language"`
	ScanInterval int `json:"ScanInterval"`
}

// NewUserSettings creates a new UserSettings struct
//
// Parameters: _
//
// Returns: a pointer to a new UserSettings struct
func NewUserSettings() *UserSettings {
	return &UserSettings{
		Language:     LoadUserSettings().Language,
		ScanInterval: LoadUserSettings().ScanInterval,
	}
}

// LoadUserSettings loads the user setting from the user settings file
//
// Parameters: _
//
// Returns: language setting (int)
func LoadUserSettings() UserSettings {
	data, err := os.ReadFile("./usersettings/usersettings.json")
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
//
// Parameters: lang (int) - the language setting to save
//
// Returns: _
func SaveUserSettings(settings UserSettings) {
	file, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		logger.Log.ErrorWithErr("Error marshalling user settings JSON:", err)
		return
	}
	err = os.WriteFile("./usersettings/usersettings.json", file, 0644)
	if err != nil {
		logger.Log.ErrorWithErr("Error writing user setting(s) to file:", err)
	}
}
