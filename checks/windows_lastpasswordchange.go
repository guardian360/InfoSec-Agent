package checks

import (
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"
)

func LastPasswordChange() {
	username, err := getCurrentUsername()
	if err != nil {
		println("error retrieving username")
		return
	}

	cmd := exec.Command("net", "user", username)
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	datePattern := `\b(\d{1,2}(-|/)\d{1,2}(-|/)\d{4})\b` //regex expression for the date
	regex := regexp.MustCompile(datePattern)
	match := regex.FindString(lines[8]) //gets the string which matches the regex expression

	var date time.Time
	dateFormats := []string{"2/1/2006", "2-1-2006", "1/2/2006", "1-2-2006", "2/1/06", "2-1-06", "1/2/06", "1-2-06"}

	// Try parsing the date with different formats
	for _, format := range dateFormats {
		date, err = time.Parse(format, match)

		// Stop on successful parse
		if err == nil {
			break
		}
	}

	if err != nil {
		println("error parsing date")
		return
	}

	fmt.Println(date)
	currentTime := time.Now() //get current time
	difference := currentTime.Sub(date)
	// Define the duration for half a year
	halfYear := 365 / 2 * 24 * time.Hour
	if difference > halfYear {
		println("On", match, "the password was changed for the last time, this is more than half a year ago so you should change it again")
	} else {
		println("You changed your password recently on", match)
	}
}

func getCurrentUsername() (string, error) {
	currentUser, err := user.Current()
	if currentUser.Username == "" || err != nil {
		return "", fmt.Errorf("failed to retrieve current username")
	}
	return strings.Split(currentUser.Username, "\\")[1], nil
}
