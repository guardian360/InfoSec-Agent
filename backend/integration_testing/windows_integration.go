package integration

import (
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAdvertisementActive(t *testing.T) {
	result := windows.Advertisement(mocking.CurrentUser)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationAdvertisementNotActive(t *testing.T) {
	result := windows.Advertisement(mocking.CurrentUser)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationAutomatedLoginActive(t *testing.T) {
	result := windows.AutomaticLogin(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationAutomatedLoginNotActive(t *testing.T) {
	result := windows.AutomaticLogin(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationCredentialGuardDisabled(t *testing.T) {
	result := windows.CredentialGuardRunning(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationDefenderAllActive(t *testing.T) {
	result := windows.Defender(mocking.LocalMachine, mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationFirewallEnabled(t *testing.T) {
	result := windows.FirewallEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationFirewallDisabled(t *testing.T) {
	result := windows.FirewallEnabled(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationGuestAccountActive(t *testing.T) {
	result := windows.GuestAccount(
		&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{},
		&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationGuestAccountNotActive(t *testing.T) {
	result := windows.GuestAccount(
		&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{},
		&mocking.RealCommandExecutor{}, &mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 2, result.ResultID)
}

func TestIntegrationLastPasswordChangeValid(t *testing.T) {
	result := windows.LastPasswordChange(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationLoginMethodPasswordOnly(t *testing.T) {
	result := windows.LoginMethod(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 4, result.ResultID)
}

func TestIntegrationLoginMethodPINOnly(t *testing.T) {
	result := windows.LoginMethod(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationOutdatedWinUpToDate(t *testing.T) {
	result := windows.Outdated(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationOutdatedWinNotUpToDate(t *testing.T) {
	result := windows.Outdated(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationPasswordComplexityValid(t *testing.T) {
	result := windows.PasswordLength(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationPasswordComplexityInvalid(t *testing.T) {
	result := windows.PasswordLength(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationPermissionWithApps(t *testing.T) {
	result := windows.Permission(checks.MicrophoneID, "microphone", mocking.CurrentUser)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationPermissionWithoutApps(t *testing.T) {
	result := windows.Permission(checks.MicrophoneID, "microphone", mocking.CurrentUser)
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationRemoteDesktopEnabled(t *testing.T) {
	result := windows.RemoteDesktopCheck(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationRemoteDesktopDisabled(t *testing.T) {
	result := windows.RemoteDesktopCheck(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationRemoteRPCEnabled(t *testing.T) {
	result := windows.AllowRemoteRPC(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationRemoteRPCDisabled(t *testing.T) {
	result := windows.AllowRemoteRPC(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationScreenLockDisabled(t *testing.T) {
	result := windows.ScreenLockEnabled(mocking.CurrentUser)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationSecureBootEnabled(t *testing.T) {
	result := windows.SecureBoot(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationSecureBootDisabled(t *testing.T) {
	result := windows.SecureBoot(mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationStartupWithApps(t *testing.T) {
	result := windows.Startup(mocking.CurrentUser, mocking.LocalMachine, mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationStartupWithoutApps(t *testing.T) {
	result := windows.Startup(mocking.CurrentUser, mocking.LocalMachine, mocking.LocalMachine)
	require.NotEmpty(t, result)
	require.Empty(t, result.Result)
	require.Equal(t, 0, result.ResultID)
}

func TestIntegrationUACFullEnabled(t *testing.T) {
	result := windows.UACCheck(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 1, result.ResultID)
}

func TestIntegrationUACPartialEnabled(t *testing.T) {
	result := windows.UACCheck(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 2, result.ResultID)
}

func TestIntegrationUACDisabled(t *testing.T) {
	result := windows.UACCheck(&mocking.RealCommandExecutor{})
	require.NotEmpty(t, result)
	require.Equal(t, 0, result.ResultID)
}
