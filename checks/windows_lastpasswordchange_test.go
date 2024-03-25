package checks

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

type MockUtils struct {
	Username string
	Err      error
}

func (m *MockUtils) CurrentUsername() (string, error) {
	return m.Username, m.Err
}

func TestLastPasswordChange(t *testing.T) {
	tests := []struct {
		name          string
		executorClass *utils.MockCommandExecutor
		want          Check
	}{
		{
			name:          "Password not changed recently",
			executorClass: &utils.MockCommandExecutor{Output: "Gebruikersnaam                           test\nVolledige naam                           test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt                         Nooit\n\nWachtwoord voor het laatst ingesteld     1-1-2022 17:48:16\nWachtwoord verloopt                      Nooit\nWachtwoord mag worden gewijzigd          1-1-2022 17:48:16\nWachtwoord vereist                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\nWerkstations toegestaan                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding                 Nooit\n\nToegestane aanmeldingstijden             Alle\n\nLidmaatschap lokale groep                *Administrators\n                                         *Apparaatbeheerders\n                                         *docker-users\n                                         *Gebruikers\n                                         *Prestatielogboekgebru\nLidmaatschap globale groep               *Geen\nDe opdracht is voltooid.", Err: nil},
			want:          NewCheckResult("LastPasswordChange", "Password last changed on 1-1-2022 , your password was changed more than half a year ago so you should change it again"),
		},
		{
			name:          "Parsing data error",
			executorClass: &utils.MockCommandExecutor{Output: "Gebruikersnaam                           test\nVolledige naam                           test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt                         Nooit\n\nWachtwoord voor het laatst ingesteld     1-0.5-2022 17:48:16\nWachtwoord verloopt                      Nooit\nWachtwoord mag worden gewijzigd          1-0.5-2022 17:48:16\nWachtwoord vereist                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\nWerkstations toegestaan                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding                 Nooit\n\nToegestane aanmeldingstijden             Alle\n\nLidmaatschap lokale groep                *Administrators\n                                         *Apparaatbeheerders\n                                         *docker-users\n                                         *Gebruikers\n                                         *Prestatielogboekgebru\nLidmaatschap globale groep               *Geen\nDe opdracht is voltooid.", Err: nil},
			want:          NewCheckError("LastPasswordChange", fmt.Errorf("error parsing date")),
		},
		{
			name:          "Password changed recently",
			executorClass: &utils.MockCommandExecutor{Output: "Gebruikersnaam                           test\nVolledige naam                           test\nOpmerking\nOpmerking van gebruiker\nLandcode                                 000 (Systeemstandaard)\nAccount actief                           Ja\nAccount verloopt                         Nooit\n\nWachtwoord voor het laatst ingesteld     1-1-2024 17:48:16\nWachtwoord verloopt                      Nooit\nWachtwoord mag worden gewijzigd          1-1-2024 17:48:16\nWachtwoord vereist                       Ja\nGebruiker mag wachtwoord wijzigen        Ja\n\nWerkstations toegestaan                  Alle\nAanmeldingsscript\nGebruikersprofiel\nBasismap\nMeest recente aanmelding                 Nooit\n\nToegestane aanmeldingstijden             Alle\n\nLidmaatschap lokale groep                *Administrators\n                                         *Apparaatbeheerders\n                                         *docker-users\n                                         *Gebruikers\n                                         *Prestatielogboekgebru\nLidmaatschap globale groep               *Geen\nDe opdracht is voltooid.", Err: nil},
			want:          NewCheckResult("LastPasswordChange", "You changed your password recently on 1-1-2024"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastPasswordChange(tt.executorClass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastPasswordChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
