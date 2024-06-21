// Package localization is responsible for localizing strings for different languages.
// The supported languages are German, English (UK), English (US), Spanish, French, Dutch, and Portuguese.
//
// Exported function(s): Init, Localize, Localizers
package localization

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/config"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizers [7]*i18n.Localizer
var bundle *i18n.Bundle

// Init is a function that sets up the localization bundle and localizers for different languages.
//
// It loads localization files from a specified path and creates localizers for each supported language.
// The supported languages are German, English (UK), English (US), Spanish, French, Dutch, and Portuguese.
//
// Parameter:
//   - path: A string representing the path to the directory containing the localization files. Each file should be named as "active.<language_code>.json".
//
// Returns: None. This function initializes global variables within the package.
//
// Note: This function should be called before using the Localize function to ensure that the localizers are properly set up.
func Init(path string) {
	logger.Log.Debug("Initializing localization files")

	path += config.LocalizationPath
	logger.Log.Debug("Localization path: " + path)

	bundle = i18n.NewBundle(language.BritishEnglish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	messageFiles := []string{
		"de/active.de.json",
		"en-GB/active.en-GB.json",
		"en-US/active.en-US.json",
		"es/active.es.json",
		"fr/active.fr.json",
		"nl/active.nl.json",
		"pt/active.pt.json",
	}
	for _, file := range messageFiles {
		if _, err := bundle.LoadMessageFile(path + file); err != nil {
			logger.Log.ErrorWithErr("Error loading localization file "+file, err)
		}
	}
	logger.Log.Debug("Localization files initialized")

	// Localizes for each language
	localizers[0] = i18n.NewLocalizer(bundle, language.German.String())
	localizers[1] = i18n.NewLocalizer(bundle, language.BritishEnglish.String())
	localizers[2] = i18n.NewLocalizer(bundle, language.AmericanEnglish.String())
	localizers[3] = i18n.NewLocalizer(bundle, language.Spanish.String())
	localizers[4] = i18n.NewLocalizer(bundle, language.French.String())
	localizers[5] = i18n.NewLocalizer(bundle, language.Dutch.String())
	localizers[6] = i18n.NewLocalizer(bundle, language.Portuguese.String())
}

// Localize is a function that retrieves and returns a localized string based on the provided language and message ID.
//
// Parameters:
//   - language: An integer representing the index of the language in the localizers array. The language should correspond to one of the supported languages (0: German, 1: English (UK), 2: English (US), 3: Spanish, 4: French, 5: Dutch, 6: Portuguese).
//   - messageID: A string representing the ID of the message to be localized. This ID should correspond to a key in the localization files.
//
// Returns:
//   - A string containing the localized message corresponding to the provided message ID and language. If the message ID does not exist in the localization files for the specified language, the function will return the message ID as is.
//
// Note: The Init function should be called before using this function to ensure that the localizers are properly set up.
func Localize(language int, messageID string) string {
	return localizers[language].MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}

// Localizers is a function that provides access to the array of localizer objects used for string localization.
//
// Parameters: None.
//
// Returns:
//   - An array of pointers to i18n.Localizer objects. Each localizer corresponds to a supported language (0: German, 1: English (UK), 2: English (US), 3: Spanish, 4: French, 5: Dutch, 6: Portuguese).
//
// Note: The Init function should be called before using this function to ensure that the localizers are properly set up.
func Localizers() [7]*i18n.Localizer {
	return localizers
}
