package windows

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/logger"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// LastPasswordChange is a function that checks the last time the Windows password was changed.
//
// Parameters:
//   - executor mocking.CommandExecutor: An executor to run the command for retrieving the last password change date.
//   - usernameRetriever mocking.UsernameRetriever: An instance of UsernameRetriever used to retrieve the current username.
//
// Returns:
//   - Check: A struct containing the result of the check. The result indicates the date when the password was last changed.
//
// The function works by executing a 'net user' command to get the user's password last set date. It then parses the output of the command to extract the date. The function compares this date with the current date and if the difference is more than half a year, it returns a warning suggesting the user to change the password. Otherwise, it returns a message indicating that the password was changed recently.
func LastPasswordChange(executor mocking.CommandExecutor, usernameRetriever mocking.UsernameRetriever) checks.Check {
	// Get the current Windows username
	username, err := usernameRetriever.CurrentUsername()
	if err != nil {
		return checks.NewCheckErrorf(checks.LastPasswordChangeID, "error retrieving username", err)
	}
	// Get the date format from the system
	dateOutput, dateErr := executor.Execute("powershell", "(Get-Culture).DateTimeFormat.ShortDatePattern")
	if dateErr != nil {
		logger.Log.ErrorWithErr("Error getting date format", dateErr)
	}
	// Trim the output to remove any leading or trailing whitespace
	dateFormat := strings.TrimSpace(string(dateOutput))
	dateFormat = strings.ReplaceAll(dateFormat, "/", "-")

	output, err := executor.Execute("net", "user", username)
	if err != nil {
		return checks.NewCheckErrorf(checks.LastPasswordChangeID, "error executing net user", err)
	}

	lines := strings.Split(string(output), "\n")
	// Define the regex pattern for the date
	datePattern := `\b(\d{1,2}(-|/)\d{1,2}(-|/)\d{4})\b`
	regex := regexp.MustCompile(datePattern)

	// Determine the current user's date format
	var goDateFormat string
	internationalDate := "02-01-2006"
	usDate := "01-02-2006"
	switch dateFormat {
	case "d-M-yyyy":
		goDateFormat = internationalDate
	case "M-d-yyyy":
		goDateFormat = usDate
	default:
		logger.Log.Error("Unknown date format:")
		goDateFormat = internationalDate
	}
	if len(lines) < 9 {
		return checks.NewCheckError(checks.LastPasswordChangeID, errors.New("error parsing output"))
	}
	// Find the date in the output
	match := regex.FindString(lines[8])
	match = strings.ReplaceAll(match, "/", "-")

	// Split the string on both "-" and "/"
	parts := strings.Split(match, "-")

	// Check each part and add a leading zero if the length is 1
	for i, part := range parts {
		if len(part) == 1 {
			parts[i] = "0" + part
		}
	}

	// Join the parts back together
	formattedDate := strings.Join(parts, "-")
	parsedDate, err := time.Parse(goDateFormat, formattedDate)
	if err != nil {
		logger.Log.ErrorWithErr("Error parsing date", err)
		parsedDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// Get the current time
	currentTime := time.Now()
	difference := currentTime.Sub(parsedDate)
	// Define the duration of half a year
	halfYear := 365 / 2 * 24 * time.Hour
	// If it has been more than half a year since the password was last changed, return a warning
	if difference > halfYear {
		return checks.NewCheckResult(checks.LastPasswordChangeID, 0, match)
	}
	return checks.NewCheckResult(checks.LastPasswordChangeID, 1, match)
}
