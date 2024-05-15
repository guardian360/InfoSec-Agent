package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/cisregistrysettings"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationCISRegistrySettingsCorrect(t *testing.T) {
	result := cisregistrysettings.CISRegistrySettings(mocking.LocalMachine, mocking.UserProfiles)
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationCISRegistrySettingsIncorrect(t *testing.T) {
	result := cisregistrysettings.CISRegistrySettings(mocking.LocalMachine, mocking.UserProfiles)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}
