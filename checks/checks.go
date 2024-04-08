// Package checks implements different security/privacy checks
//
// Exported function(s): PasswordManager, WindowsDefender, LastPasswordChange, LoginMethod, Permission, Bluetooth,
// OpenPorts, WindowsOutdated, SecureBoot, SmbCheck, Startup, GuestAccount, UACCheck, RemoteDesktopCheck,
// ExternalDevices, NetworkSharing
package checks

import "fmt"

// Check is a struct that holds the result of a security/privacy check.
//
// Each check returns an ID, a Result, and an Error.
//
// The datatype error can not be (directly) serialised to JSON, so we also include an ErrorMSG field.
//
// A new Check struct can be created with the accompanying functions: NewCheckResult, NewCheckError, NewCheckErrorf
type Check struct {
	IssueID  int      `json:"issue_id"`
	ResultID int      `json:"result_id"`
	Result   []string `json:"result,omitempty"`
	Error    error    `json:"-"` // Don't serialize error field to JSON
	ErrorMSG string   `json:"error,omitempty"`
}

// NewCheckResult creates a new Check struct with only a result
func NewCheckResult(issID int, resID int, result ...string) Check {
	return Check{IssueID: issID, ResultID: resID, Result: result}
}

// NewCheckError creates a new Check struct with the error and error message
func NewCheckError(id int, err error) Check {
	return Check{IssueID: id, Error: err, ErrorMSG: err.Error()}
}

// NewCheckErrorf creates a new Check struct with the error and formatted error message
func NewCheckErrorf(id int, message string, err error) Check {
	formatErr := fmt.Errorf(message+": %w", err)
	return Check{IssueID: id, Error: formatErr, ErrorMSG: formatErr.Error()}
}
