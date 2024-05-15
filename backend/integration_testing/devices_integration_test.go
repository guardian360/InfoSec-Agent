package integration_testing

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/devices"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
	"testing"
)

func TestIntegrationBluetoothDevices(t *testing.T) {
	result := devices.Bluetooth(mocking.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
	// Check if function did not return an error
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationBluetoothNoDevices(t *testing.T) {
	result := devices.Bluetooth(mocking.NewRegistryKeyWrapper(registry.LOCAL_MACHINE))
	// Check if function did not return an error
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationExternalDevicesDevices(t *testing.T) {
	result := devices.ExternalDevices(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationExternalDevicesNoDevices(t *testing.T) {
	result := devices.ExternalDevices(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}
