package checks

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/mocking"
	"github.com/InfoSec-Agent/InfoSec-Agent/utils"
)

// LastPasswordChange checks when the Windows password was last changed
//
// Parameters: _
//
// Returns: When the password was last changed
func LastPasswordChange(executor mocking.CommandExecutor) Check {
	// Get the current Windows username
	username, err := utils.CurrentUsername()
	if err != nil {
		return NewCheckErrorf(LastPasswordChangeID, "error retrieving username", err)
	}

	output, err := executor.Execute("net", "user", username)
	if err != nil {
		return NewCheckErrorf(LastPasswordChangeID, "error executing net user", err)
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
		return NewCheckError(LastPasswordChangeID, errors.New("error parsing date"))
	}

	// Get the current time
	currentTime := time.Now()
	difference := currentTime.Sub(date)
	// Define the duration of half a year
	halfYear := 365 / 2 * 24 * time.Hour
	// If it has been more than half a year since the password was last changed, return a warning
	if difference > halfYear {
		return NewCheckResult(LastPasswordChangeID, 0,
			fmt.Sprintf("Password last changed on %s , "+
				"your password was changed more than half a year ago so you should change it again", match))
	}
	return NewCheckResult(LastPasswordChangeID, 1, fmt.Sprintf("You changed your password recently on %s",
		match))
}
