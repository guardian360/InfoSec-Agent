package checks

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// LastPasswordChange is a function that checks the last time the Windows password was changed.
//
// Parameters:
//   - executor commandmock.CommandExecutor: An executor to run the command for retrieving the last password change date.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates the date when the password was last changed.
//
// The function works by executing a 'net user' command to get the user's password last set date. It then parses the output of the command to extract the date. The function compares this date with the current date and if the difference is more than half a year, it returns a warning suggesting the user to change the password. Otherwise, it returns a message indicating that the password was changed recently.
func LastPasswordChange(executor commandmock.CommandExecutor) Check {
	// Get the current Windows username
	username, err := utils.CurrentUsername()
	if err != nil {
		return NewCheckErrorf("LastPasswordChange", "error retrieving username", err)
	}

	output, err := executor.Execute("net", "user", username)
	if err != nil {
		return NewCheckErrorf("LastPasswordChange", "error executing net user", err)
	}

	lines := strings.Split(string(output), "\n")
	// Define the regex pattern for the date
	datePattern := `\b(\d{1,2}(-|/)\d{1,2}(-|/)\d{4})\b`
	regex := regexp.MustCompile(datePattern)
	// Find the date in the output
	match := regex.FindString(lines[8])

	var date time.Time
	// Define different valid date formats
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
		return NewCheckError("LastPasswordChange", errors.New("error parsing date"))
	}

	// Get the current time
	currentTime := time.Now()
	difference := currentTime.Sub(date)
	// Define the duration of half a year
	halfYear := 365 / 2 * 24 * time.Hour
	// If it has been more than half a year since the password was last changed, return a warning
	if difference > halfYear {
		return NewCheckResult("LastPasswordChange",
			fmt.Sprintf("Password last changed on %s , "+
				"your password was changed more than half a year ago so you should change it again", match))
	}
	return NewCheckResult("LastPasswordChange", fmt.Sprintf("You changed your password recently on %s",
		match))
}
