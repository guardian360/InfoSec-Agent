package localization

import go_i18n "github.com/nicksnyder/go-i18n/v2/i18n"

func init() {
	// Register the localizations
	go_i18n.MustLoadTranslationFile("localization/en-US.all.json")
	go_i18n.MustLoadTranslationFile("localization/nl-NL.all.json")
}
