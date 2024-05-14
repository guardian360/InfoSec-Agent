package windows

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// AutomaticLogin checks if automatic log-in is enabled on the system.
//
// Parameters:
//   - registryKey (mocking.RegistryKey): A mocker of a Windows registry key. This is used to simulate the behavior of the Windows registry for testing purposes.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates whether automatic log-in is enabled on the system.
func AutomaticLogin(registryKey mocking.RegistryKey) checks.Check {
	key, err := checks.OpenRegistryKey(registryKey, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`)
	if err != nil {
		return checks.NewCheckError(checks.AutoLoginID, err)
	}
	defer checks.CloseRegistryKey(key)

	// Read the value of AutoAdminLogon, which contains the information if automatic log-in is enabled on the system.
	// If the registry key does not exist or its value is 0, then automatic log-in is not enabled.
	value, _, err := key.GetIntegerValue("AutoAdminLogon")
	if err != nil {
		logger.Log.ErrorWithErr("Error reading AutoAdminLogon value", err)
		return checks.NewCheckResult(checks.AutoLoginID, 0)
	}

	if value == 1 {
		return checks.NewCheckResult(checks.AutoLoginID, 1)
	}
	return checks.NewCheckResult(checks.AutoLoginID, 0)
}
