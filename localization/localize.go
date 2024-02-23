package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Localizer *i18n.Localizer
var bundle *i18n.Bundle

func Init() { //3
	bundle = i18n.NewBundle(language.BritishEnglish)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("localizations_src/de/tray.json")
	bundle.LoadMessageFile("localizations_src/en-GB/tray.json")
	bundle.LoadMessageFile("localizations_src/en-US/tray.json")
	bundle.LoadMessageFile("localizations_src/es/tray.json")
	bundle.LoadMessageFile("localizations_src/fr/tray.json")
	bundle.LoadMessageFile("localizations_src/nl/tray.json")
	bundle.LoadMessageFile("localizations_src/pt/tray.json")

	Localizer = i18n.NewLocalizer(bundle, language.BritishEnglish.String(), language.German.String(),
		language.AmericanEnglish.String(), language.Spanish.String(),
		language.French.String(), language.Dutch.String(), language.Portuguese.String())
}
