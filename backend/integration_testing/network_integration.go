package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/network"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

func TestIntegrationNetBIOSEnabled(t *testing.T) {
	result := network.NetBIOSEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationNetBIOSDisabled(t *testing.T) {
	result := network.NetBIOSEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationWPADEnabled(t *testing.T) {
	result := network.WPADEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationWPADDisabled(t *testing.T) {
	result := network.WPADEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationOpenPortsPorts(t *testing.T) {
	result := network.OpenPorts(&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{})
	// Check if function did not return an error
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationSmbCheckGoodSetup(t *testing.T) {
	result := network.SmbCheck(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 2, result.ResultID)
}

func TestIntegrationSmbCheckBadSetup(t *testing.T) {
	result := network.SmbCheck(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 3, result.ResultID)
}
