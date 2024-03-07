package checks

import (
	"fmt"

	"log"

	"golang.org/x/sys/windows/registry"
)

func Startup() {
	// Open the relevant key so we can get the startup data entries out of them
	cuKey := openRegistryKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey := openRegistryKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`)
	lmKey2 := openRegistryKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`)

	defer cuKey.Close()
	defer lmKey.Close()
	defer lmKey2.Close()

	// Read the entries within the registry key
	cuValueNames, err := cuKey.ReadValueNames(0)
	lmValueNames, err := lmKey.ReadValueNames(0)
	lm2ValueNames, err := lmKey2.ReadValueNames(0)

	if err != nil {
		fmt.Println("Error reading value names:", err)
		return
	}

	fmt.Println("The following programs are enabled to run on startup:")
	printEntries(cuValueNames, cuKey)
	printEntries(lmValueNames, lmKey)
	printEntries(lm2ValueNames, lmKey2)
}

// Simple function to open registry keys and handle errors
func openRegistryKey(k registry.Key, path string) registry.Key {
	key, err := registry.OpenKey(k, path, registry.READ)

	if err != nil {
		log.Fatal("Error opening registry key:", err)
	}

	return key
}

// Print the values of the entries inside the corresponding registry key
func printEntries(entries []string, key registry.Key) {

	for _, element := range entries {
		val, _, _ := key.GetBinaryValue(element)

		// We check the binary values to make sure we only print the programs ENABLED to startup
		// For example: in registry it lists all startup programs, but also the ones that are disabled on startup.
		if val[4] == 0 && val[5] == 0 && val[6] == 0 {
			fmt.Println(element)
		}
	}
}
