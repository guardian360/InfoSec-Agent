package windows

import (
	"strconv"

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
	key, err := mocking.OpenRegistryKey(registryKey, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`)
	if err != nil {
		return checks.NewCheckError(checks.AutoLoginID, err)
	}
	defer mocking.CloseRegistryKey(key)

	// Read the value of AutoAdminLogon, which contains the information if automatic log-in is enabled on the system.
	// If the registry key does not exist or its value is 0, then automatic log-in is not enabled.
	value, _, err := key.GetStringValue("AutoAdminLogon")
	if err != nil {
		logger.Log.ErrorWithErr("Error reading AutoAdminLogon value", err)
		return checks.NewCheckResult(checks.AutoLoginID, 0)
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		logger.Log.ErrorWithErr("Error converting AutoAdminLogon value to int", err)
		return checks.NewCheckResult(checks.AutoLoginID, 0)
	}
	return checks.NewCheckResult(checks.AutoLoginID, intVal)
}
