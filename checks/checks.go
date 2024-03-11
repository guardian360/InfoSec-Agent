// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import "fmt"

// Check is a struct that holds the result of a  security/privacy check.
//
// Each check returns an Id, a Result, and an Error.
//
// The datatype error can not be (directly) serialised to JSON, so we also include an ErrorMSG field.
//
// A new Check struct can be created with the accompanying functions: newCheckResult, newCheckError, newCheckErrorf
type Check struct {
	Id       string   `json:"id"`
	Result   []string `json:"result,omitempty"`
	Error    error    `json:"-"` // Don't serialize error field to JSON
	ErrorMSG string   `json:"error,omitempty"`
}

// newCheckResult creates a new Check struct with only a result
func newCheckResult(id string, result ...string) Check {
	return Check{Id: id, Result: result}
}

// newCheckError creates a new Check struct with the error and error message
func newCheckError(id string, err error) Check {
	return Check{Id: id, Error: err, ErrorMSG: err.Error()}
}

// newCheckErrorf creates a new Check struct with the error and formatted error message
func newCheckErrorf(id string, message string, err error) Check {
	formatErr := fmt.Errorf(message+": %w", err)
	return Check{Id: id, Error: formatErr, ErrorMSG: formatErr.Error()}
}
