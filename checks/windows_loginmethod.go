package checks

import (
	"github.com/InfoSec-Agent/InfoSec-Agent/registrymock"
)

// LoginMethod checks which login method(s) the user has enabled
//
// Parameters: _
//
// Returns: List of login methods enabled
func LoginMethod(registryKey registrymock.RegistryKey) Check {
	var resultID int
	// Open the registry key related to log-in methods
	key, err := registrymock.OpenRegistryKey(registryKey,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\LogonUI\UserTile`)
	if err != nil {
		return NewCheckErrorf(LoginMethodID, "error opening registry key", err)
	}
	// Close the key after we have received all relevant information
	defer registrymock.CloseRegistryKey(key)

	// Read the info of the key
	keyInfo, err := key.Stat()
	if err != nil {
		return NewCheckErrorf(LoginMethodID, "error getting key info", err)
	}

	// Read the value names, which correspond to different log-in methods
	names, err := key.ReadValueNames(int(keyInfo.ValueCount))
	if err != nil {
		return NewCheckErrorf(LoginMethodID, "error reading value names", err)
	}

	var resultString []string

	// Each log-in method corresponds to a unique GUID
	// Check whether the GUID is present in the registry key, and if it is, that log-in method is enabled
	for _, element := range names {
		switch {
		case registrymock.CheckKey(key, element) == "{D6886603-9D2F-4EB2-B667-1971041FA96B}":
			resultID |= 1 << 0
			resultString = append(resultString, "PIN")
		case registrymock.CheckKey(key, element) == "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}":
			resultID |= 1 << 1
			resultString = append(resultString, "Picture Logon")
		case registrymock.CheckKey(key, element) == "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}":
			resultID |= 1 << 2
			resultString = append(resultString, "Password")
		case registrymock.CheckKey(key, element) == "{BEC09223-B018-416D-A0AC-523971B639F5}":
			resultID |= 1 << 3
			resultString = append(resultString, "Fingerprint")
		case registrymock.CheckKey(key, element) == "{8AF662BF-65A0-4D0A-A540-A338A999D36F}":
			resultID |= 1 << 4
			resultString = append(resultString, "Facial recognition")
		case registrymock.CheckKey(key, element) == "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}":
			resultID |= 1 << 5
			resultString = append(resultString, "Trust signal")
		default:
			return NewCheckErrorf(LoginMethodID, "error reading value", err)
		}
	}

	return NewCheckResult(LoginMethodID, resultID, resultString...)
}
