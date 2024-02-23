package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Localizers []*i18n.Localizer
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

	// Localizes for each language
	Localizers[0] = i18n.NewLocalizer(bundle, language.German.String())
	Localizers[1] = i18n.NewLocalizer(bundle, language.BritishEnglish.String())
	Localizers[2] = i18n.NewLocalizer(bundle, language.AmericanEnglish.String())
	Localizers[3] = i18n.NewLocalizer(bundle, language.Spanish.String())
	Localizers[4] = i18n.NewLocalizer(bundle, language.French.String())
	Localizers[5] = i18n.NewLocalizer(bundle, language.Dutch.String())
	Localizers[6] = i18n.NewLocalizer(bundle, language.Portuguese.String())
}
