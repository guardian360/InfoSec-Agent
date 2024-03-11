package checks

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func LoginMethod() Check {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\LogonUI\UserTile`, registry.QUERY_VALUE)
	if err != nil {
		return newCheckErrorf("LoginMethod", "error opening registry key", err)
	}
	defer key.Close()

	keyInfo, err := key.Stat()
	if err != nil {
		return newCheckErrorf("LoginMethod", "error getting key info", err)
	}

	names, err := key.ReadValueNames(int(keyInfo.ValueCount))
	if err != nil {
		return newCheckErrorf("LoginMethod", "error reading value names", err)
	}

	result := newCheckResult("LoginMethod")
	for _, element := range names {
		if checkKey(key, element) == "{D6886603-9D2F-4EB2-B667-1971041FA96B}" {
			result.Result = append(result.Result, "PIN")
		} else if checkKey(key, element) == "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}" {
			result.Result = append(result.Result, "Picture Logon")
		} else if checkKey(key, element) == "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}" {
			result.Result = append(result.Result, "Password")
		} else if checkKey(key, element) == "{BEC09223-B018-416D-A0AC-523971B639F5}" {
			result.Result = append(result.Result, "Fingerprint")
		} else if checkKey(key, element) == "{8AF662BF-65A0-4D0A-A540-A338A999D36F}" {
			result.Result = append(result.Result, "Facial recognition")
		} else if checkKey(key, element) == "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}" {
			result.Result = append(result.Result, "Trust signal")
		}
	}

	return result
}

func checkKey(key registry.Key, el string) string {
	val, _, err := key.GetStringValue(el)
	if err == nil {
		return val
	} else {
		fmt.Printf("Not able to check")
		return "-1"
	}
}
