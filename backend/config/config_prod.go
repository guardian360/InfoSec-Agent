//go:build prod

package config

const (
	LogLevel         = 2
	LogLevelSpecific = -1

	LocalizationPath = "localization/"
	DatabasePath     = "localization/en-GB"

	BuildReportingPage    = false
	ReportingPagePath     = "reporting-page/InfoSec-Agent-Reporting-Page"
	ReportingPageImageDir = "images/"
)
