package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationCurrentUsernameFound(t *testing.T) {
	result, err := checks.CurrentUsername()
	require.NotEmpty(t, result)
	require.NoError(t, err)
}
