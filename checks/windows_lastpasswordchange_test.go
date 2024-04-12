package checks_test

import (
	"errors"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/checks"
	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
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
	tests := []struct {
		name          string
		executorClass *mocking.MockCommandExecutor
		want          checks.Check
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
			want: checks.NewCheckResult(checks.LastPasswordChangeID, 0, "Password last changed on 1-1-2022 , "+
				"your password was changed more than half a year ago so you should change it again"),
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
			want: checks.NewCheckError(checks.LastPasswordChangeID, errors.New("error parsing date")),
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
			want: checks.NewCheckResult(checks.LastPasswordChangeID, 1,
				"You changed your password recently on 1-1-2024"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checks.LastPasswordChange(tt.executorClass)
			require.Equal(t, tt.want, got)
		})
	}
}
