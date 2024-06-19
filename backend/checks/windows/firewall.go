package windows

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// FirewallEnabled is a function that checks if the Windows firewall is enabled for all 3 profile types (private, public, and domain).
// It does this by executing a command to show the state of the firewall for all profiles.
//
// Parameters:
//   - executor mocking.CommandExecutor: An object that implements the CommandExecutor interface.
//
// Returns:
//   - Check: A Check object that encapsulates the outcome of the firewall check.
//     If the firewall is enabled for all profiles, the Result field of the Check object is set to 0.
//     If the firewall is disabled for any profile, the Result field is set to 1. If an error occurs during the check,
//     the Error field of the Check object is set to the error encountered.
func FirewallEnabled(executor mocking.CommandExecutor) checks.Check {
	firewallCommand := "netsh advfirewall show allprofiles state"
	output, err := executor.Execute("cmd", "/c", firewallCommand)
	if err != nil {
		logger.Log.ErrorWithErr("Error executing firewall command", err)
		return checks.NewCheckError(checks.FirewallID, err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "State") && !strings.Contains(line, "ON") {
			return checks.NewCheckResult(checks.FirewallID, 1)
		}
	}

	return checks.NewCheckResult(checks.FirewallID, 0)
}
