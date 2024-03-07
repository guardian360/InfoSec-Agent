package checks

import (
	"fmt"
	"os/exec"
)

// to do: Improve formatting of output

func Externaldevices() {
	// All the classes you want to check within the Get-PnpDevice command
	classesToCheck := [2]string{"Mouse", "Camera"}
	for _, s := range classesToCheck {
		printDeviceClass(s)
	}
}

// Run the command for a specific class within the Get-PnpDevice and print its results
func printDeviceClass(deviceClass string) {
	fmt.Printf("The following %s devices are detected:", deviceClass)
	cmd := exec.Command("powershell", "-Command", "Get-PnpDevice -Class", deviceClass, " | Where-Object -Property Status -eq 'OK'")
	output, _ := cmd.Output()
	fmt.Println(string(output))
}
