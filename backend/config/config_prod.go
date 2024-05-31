//go:build prod

package config

const (
	LogLevel         = 2
	LogLevelSpecific = -1

	LocalizationPath = "localization/"

	BuildReportingPage = false
	ReportingPagePath  = "reporting-page/InfoSec-Agent-Reporting-Page"
)
