package mocking

import (
	"errors"
	"os/user"
	"strings"

	"github.com/stretchr/testify/mock"
)

// UsernameRetriever is an interface that defines a method for retrieving the current username (CurrentUsername).
// This interface is used to abstract the retrieval of the current username, allowing for different implementations.
// For example, a real implementation that retrieves the username from the operating system, and a mock implementation for testing.
type UsernameRetriever interface {
	CurrentUsername() (string, error)
}

// MockUsernameRetriever is a struct that implements the UsernameRetriever interface.
// It uses the testify/mock package to simulate the behavior of the CurrentUsername method,
// allowing for controlled testing scenarios.
type MockUsernameRetriever struct {
	mock.Mock
}

// RealUsernameRetriever is a struct that implements the UsernameRetriever interface.
// It provides a real implementation of the CurrentUsername method, which retrieves the username of the currently logged-in user.
type RealUsernameRetriever struct{}

// CurrentUsername is a method of the RealUsernameRetriever struct that implements the UsernameRetriever interface.
// It provides a real implementation for retrieving the username of the currently logged-in user.
// Returns:
//   - string: The username of the currently logged-in user. If the username cannot be retrieved, an empty string is returned.
//   - error: An error object that wraps any error that occurs during the retrieval of the username. If the username is retrieved successfully, it returns nil.
func (r *RealUsernameRetriever) CurrentUsername() (string, error) {
	return CurrentUsername()
}

// CurrentUsername is a method of the MockUsernameRetriever struct that implements the UsernameRetriever interface.
// It simulates the retrieval of the username of the currently logged-in user for testing purposes.
// This method uses the testify/mock package to control the return values of the method call, allowing for controlled testing scenarios.
// Returns:
//   - string: The simulated username of the currently logged-in user. If the username cannot be retrieved, an empty string is returned.
//   - error: An error object that wraps any error that occurs during the simulated retrieval of the username. If the username is retrieved successfully, it returns nil.
func (m *MockUsernameRetriever) CurrentUsername() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

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
