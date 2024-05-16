package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCurrentUsernameFound(t *testing.T) {
	result, err := checks.CurrentUsername()
	require.NotEmpty(t, result)
	require.NoError(t, err)
}
