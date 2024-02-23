package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Localizers [7]*i18n.Localizer
var bundle *i18n.Bundle

func Init(path string) { //3
	bundle = i18n.NewBundle(language.BritishEnglish)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(path + "localization/localizations_src/de/active.de.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/en-GB/active.en-GB.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/en-US/active.en-US.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/es/active.es.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/fr/active.fr.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/nl/active.nl.json")
	bundle.MustLoadMessageFile(path + "localization/localizations_src/pt/active.pt.json")

	// Localizes for each language
	Localizers[0] = i18n.NewLocalizer(bundle, language.German.String())
	Localizers[1] = i18n.NewLocalizer(bundle, language.BritishEnglish.String())
	Localizers[2] = i18n.NewLocalizer(bundle, language.AmericanEnglish.String())
	Localizers[3] = i18n.NewLocalizer(bundle, language.Spanish.String())
	Localizers[4] = i18n.NewLocalizer(bundle, language.French.String())
	Localizers[5] = i18n.NewLocalizer(bundle, language.Dutch.String())
	Localizers[6] = i18n.NewLocalizer(bundle, language.Portuguese.String())
}
