package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TODO: Update documentation
// AllowRemoteRPC checks if the setting to allow remote computers to execute code on your device is enabled.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): A mocker of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether the settings to allow remote computers to execute code on your device is enabled.
func AllowRemoteRPC(registryKey mocking.RegistryKey) checks.Check {
	key, err := mocking.OpenRegistryKey(registryKey, `SYSTEM\CurrentControlSet\Control\Terminal Server`)
	if err != nil {
		return checks.NewCheckError(checks.RemoteRPCID, err)
	}
	defer mocking.CloseRegistryKey(key)

	// Read the value of AllowRemoteRPC, which contains the information if the setting to allow remote computers to execute code on your device is enabled.
	value, _, err := key.GetIntegerValue("AllowRemoteRPC")
	if err != nil {
		logger.Log.ErrorWithErr("Error reading AllowRemoteRPC value", err)
		return checks.NewCheckError(checks.RemoteRPCID, err)
	}

	return checks.NewCheckResult(checks.RemoteRPCID, int(value))
}
