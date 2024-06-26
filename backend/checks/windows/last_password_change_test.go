package windows_test

import (
	"errors"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks/windows"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// TestLastPasswordChange is a function that tests the behavior of the LastPasswordChange function with various inputs.
//
// Parameters:
//   - t *testing.T: The testing framework provided by the Go testing package. It provides methods for reporting test failures and logging additional information.
//
// Returns: None
//
// This function tests the LastPasswordChange function with different scenarios. It uses a mock implementation of the CommandExecutor interface to simulate the behavior of the command execution for retrieving the last password change date. Each test case checks if the LastPasswordChange function correctly identifies the date when the password was last changed based on the simulated command output. The function asserts that the returned Check instance contains the expected results.
func TestLastPasswordChange(t *testing.T) {
	usernameRetriever := new(mocking.MockUsernameRetriever)
	usernameRetriever.On("CurrentUsername").Return("", errors.New("current username error"))
	tests := []struct {
		name              string
		executorClass     *mocking.MockCommandExecutor
		usernameRetriever mocking.UsernameRetriever
		want              checks.Check
	}{
		{
			name: "Password not changed recently",
			executorClass: &mocking.MockCommandExecutor{
				Output: "Gebruikersnaam                           test\nVolledige naam                           " +
					"test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 " +
					"000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt" +
					"                         Nooit\n\nWachtwoord voor het laatst ingesteld     1-1-2022 17:48:16\n" +
					"Wachtwoord verloopt                      Nooit\nWachtwoord mag worden gewijzigd          " +
					"1-1-2022 17:48:16\nWachtwoord vereist                       Ja\n" +
					"Gebruiker mag wachtwoord wijzigen" +
					"        Ja\n\nWerkstations toegestaan                  Alle\n" +
					"Aanmeldingsscript\nGebruikersprofiel" +
					"\nBasismap\nMeest recente aanmelding                 Nooit\n\nToegestane aanmeldingstijden" +
					"             Alle\n\nLidmaatschap lokale groep                *Administrators\n" +
					"                                         *Apparaatbeheerders\n" +
					"                                         *docker-users\n" +
					"                                         *Gebruikers\n" +
					"                                         *Prestatielogboekgebru\nLidmaatschap globale groep" +
					"               *Geen\nDe opdracht is voltooid.", Err: nil},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckResult(checks.LastPasswordChangeID, 0, "1-1-2022"),
		},
		{
			name: "Parsing data error",
			executorClass: &mocking.MockCommandExecutor{
				Output: "Gebruikersnaam                           test\nVolledige naam                           " +
					"test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 " +
					"000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt" +
					"                         Nooit\n\nWachtwoord voor het laatst ingesteld     " +
					"1-0.5-2022 17:48:16\nWachtwoord verloopt                      Nooit\n" +
					"Wachtwoord mag worden gewijzigd          1-0.5-2022 17:48:16\nWachtwoord vereist" +
					"                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\n" +
					"Werkstations toegestaan" +
					"                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding" +
					"                 Nooit\n\nToegestane aanmeldingstijden             " +
					"Alle\n\nLidmaatschap lokale groep                *Administrators\n" +
					"                                         *Apparaatbeheerders\n" +
					"                                         *docker-users\n" +
					"                                         *Gebruikers\n" +
					"                                         *Prestatielogboekgebru\nLidmaatschap globale groep" +
					"               *Geen\nDe opdracht is voltooid.", Err: nil},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckResult(checks.LastPasswordChangeID, 0, ""),
		},
		{
			name: "Password changed recently",
			executorClass: &mocking.MockCommandExecutor{
				Output: "Gebruikersnaam                           test\nVolledige naam                           " +
					"test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 " +
					"000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt" +
					"                         Nooit\n\nWachtwoord voor het laatst ingesteld     " +
					"1-1-2024 17:48:16\nWachtwoord verloopt                      Nooit\n" +
					"Wachtwoord mag worden gewijzigd          1-1-2024 17:48:16\nWachtwoord vereist" +
					"                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\n" +
					"Werkstations toegestaan" +
					"                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding" +
					"                 Nooit\n\nToegestane aanmeldingstijden             Alle\n\n" +
					"Lidmaatschap lokale groep                *Administrators\n" +
					"                                         *Apparaatbeheerders\n" +
					"                                         *docker-users\n" +
					"                                         *Gebruikers\n" +
					"                                         *Prestatielogboekgebru\nLidmaatschap globale groep" +
					"               *Geen\nDe opdracht is voltooid.", Err: nil},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckResult(checks.LastPasswordChangeID, 1, "1-1-2024"),
		},
		{
			name: "Error executing net user",
			executorClass: &mocking.MockCommandExecutor{
				Output: ".", Err: errors.New("error executing net user")},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckErrorf(checks.LastPasswordChangeID, "error executing net user", errors.New("error executing net user")),
		},
		{
			name: "CurrentUsername Error",
			executorClass: &mocking.MockCommandExecutor{
				Output: "Gebruikersnaam                           test\nVolledige naam                           " +
					"test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 " +
					"000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt" +
					"                         Nooit\n\nWachtwoord voor het laatst ingesteld     " +
					"1-1-2024 17:48:16\nWachtwoord verloopt                      Nooit\n" +
					"Wachtwoord mag worden gewijzigd          1-1-2024 17:48:16\nWachtwoord vereist" +
					"                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\n" +
					"Werkstations toegestaan" +
					"                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding" +
					"                 Nooit\n\nToegestane aanmeldingstijden             Alle\n\n" +
					"Lidmaatschap lokale groep                *Administrators\n" +
					"                                         *Apparaatbeheerders\n" +
					"                                         *docker-users\n" +
					"                                         *Gebruikers\n" +
					"                                         *Prestatielogboekgebru\nLidmaatschap globale groep" +
					"               *Geen\nDe opdracht is voltooid.", Err: nil},
			usernameRetriever: usernameRetriever,
			want:              checks.NewCheckErrorf(checks.LastPasswordChangeID, "error retrieving username", errors.New("current username error")),
		},
		{
			name: "Other date format",
			executorClass: &mocking.MockCommandExecutor{
				Output: "d-M-yyyy", Err: nil},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckError(checks.LastPasswordChangeID, errors.New("error parsing output")),
		},
		{
			name: "Other date format",
			executorClass: &mocking.MockCommandExecutor{
				Output: "M-d-yyyy", Err: nil},
			usernameRetriever: &mocking.RealUsernameRetriever{},
			want:              checks.NewCheckError(checks.LastPasswordChangeID, errors.New("error parsing output")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windows.LastPasswordChange(tt.executorClass, tt.usernameRetriever)
			require.Equal(t, tt.want, got)
		})
	}
}
