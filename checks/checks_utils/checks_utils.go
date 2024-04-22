package checks_utils

import (
	"errors"
	"os/user"
	"strings"
)

// CurrentUsername retrieves the username of the currently logged-in user in a Windows environment.
//
// This function uses the os/user package to access the current user's information.
// It then parses the Username field to extract the actual username, discarding the domain if present.
//
// Returns:
//   - string: The username of the currently logged-in user. If the username cannot be retrieved, an empty string is returned.
//   - error: An error object that wraps any error that occurs during the retrieval of the username. If the username is retrieved successfully, it returns nil.
func CurrentUsername() (string, error) {
	currentUser, err := user.Current()
	if currentUser.Username == "" || err != nil {
		return "", errors.New("failed to retrieve current username")
	}
	return strings.Split(currentUser.Username, "\\")[1], nil
}
