package checks

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func Loginmethod() {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\LogonUI\UserTile`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()
	keyinfo, err := key.Stat()

	names, err := key.ReadValueNames(int(keyinfo.ValueCount))
	if err == nil {
		for _, element := range names {
			if checkKey(key, element) == "{D6886603-9D2F-4EB2-B667-1971041FA96B}" {
				fmt.Println("PIN is enabled")
			} else if checkKey(key, element) == "{2135F72A-90B5-4ED3-A7F1-8BB705AC276A}" {
				fmt.Println("Picture Logon is enabled")
			} else if checkKey(key, element) == "{60B78E88-EAD8-445C-9CFD-0B87F74EA6CD}" {
				fmt.Println("Password is enabled")
			} else if checkKey(key, element) == "{BEC09223-B018-416D-A0AC-523971B639F5}" {
				fmt.Println("Fingerprint is enabled")
			} else if checkKey(key, element) == "{8AF662BF-65A0-4D0A-A540-A338A999D36F}" {
				fmt.Println("Facial recognition")
			} else if checkKey(key, element) == "{27FBDB57-B613-4AF2-9D7E-4FA7A66C21AD}" {
				fmt.Println("Trust signal")
			}
		}
	} else {
		fmt.Print("error here")
	}
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
