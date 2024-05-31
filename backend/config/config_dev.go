//go:build !prod

// Package config contains the configuration settings for the backend application.
// It defines different constants depending on the build tags.
package config

const (
	LocalizationPath = "backend/localization/localizations_src/"
)
