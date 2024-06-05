package network

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// WPADEnabled checks if the WPAD service is enabled.
// The WPAD service is used to automatically configure proxy settings for a network.
// If the service is running, it is possible that an attacker could use it to redirect traffic.
//
// Parameters:
//   - executor mocking.CommandExecutor: An object that implements the CommandExecutor interface.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the WPAD check, or an error if one occurred.
func WPADEnabled(executor mocking.CommandExecutor) checks.Check {
	output, err := executor.Execute("cmd", "/c", "sc query winhttpautoproxysvc")
	if err != nil {
		logger.Log.ErrorWithErr("Error executing WPAD command: ", err)
		return checks.NewCheckError(checks.WPADID, err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "STATE") {
			if strings.Contains(line, "RUNNING") {
				return checks.NewCheckResult(checks.WPADID, 1)
			}
		}
	}

	return checks.NewCheckResult(checks.WPADID, 0)
}
