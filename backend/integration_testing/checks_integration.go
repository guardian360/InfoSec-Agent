package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCurrentUsernameFound(t *testing.T) {
	result, err := mocking.CurrentUsername()
	require.NotEmpty(t, result)
	require.NoError(t, err)
	require.Equal(t, "Test", result)
}
