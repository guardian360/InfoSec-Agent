package windows

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// NetBIOSEnabled is a function that checks if NetBIOS over TCP/IP is enabled.
// It does this by executing a command to show the state of NetBIOS over TCP/IP.
//
// Parameters:
//   - executor mocking.CommandExecutor: An object that implements the CommandExecutor interface.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the NetBIOS check, or an error if one occurred.
func NetBIOSEnabled(executor mocking.CommandExecutor) checks.Check {
	output, err := executor.Execute("cmd", "/c", "ipconfig /all")
	if err != nil {
		logger.Log.ErrorWithErr("Error executing firewall command: ", err)
		return checks.NewCheckError(checks.NetBIOSID, err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "NetBIOS over Tcpip") {
			if strings.Contains(line, "Enabled") {
				return checks.NewCheckResult(checks.NetBIOSID, 1)
			}
		}
	}

	return checks.NewCheckResult(checks.NetBIOSID, 0)
}
