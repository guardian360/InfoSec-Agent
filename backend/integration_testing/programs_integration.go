package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/programs"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPasswordManagerPresent(t *testing.T) {
	result := programs.PasswordManager(mocking.RealProgramLister{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationPasswordManagerNotPresent(t *testing.T) {
	result := programs.PasswordManager(mocking.RealProgramLister{})
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}
