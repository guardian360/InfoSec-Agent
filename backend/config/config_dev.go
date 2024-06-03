//go:build !prod

// Package config contains the configuration settings for the backend application.
// It defines different constants depending on the build tags.
package config

const (
	LogLevel         = 0
	LogLevelSpecific = -1

	LocalizationPath = "backend/localization/localizations_src/"

	BuildReportingPage    = false
	ReportingPagePath     = "reporting-page/build/bin/InfoSec-Agent-Reporting-Page"
	ReportingPageImageDir = "frontend/src/assets/images/"
)
