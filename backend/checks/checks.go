// Package checks implements different security/privacy checks.
//
// This package provides a set of functions and types that are used to perform various security and privacy checks.
// The checks themselves are implemented in separate subpackages.
// All checks have an associated issue ID, which is used to identify the check and its results, defined as constants in this package.
//
// The main type is the Check struct, which represents the outcome of a check.
// This package provides several constructor functions for creating instances of the Check struct.
package checks

import "fmt"

// TODO: Update documentation
// Check is a struct that encapsulates the outcome of a security or privacy check.
//
// Fields:
//   - IssueID (int): A unique identifier for the issue. This value is used to distinguish between different checks.
//   - ResultID (int): A unique identifier for the result. This value is used to distinguish between different results of a check.
//   - Result ([]string): The outcome of the check. This could be a list of strings representing various results.
//   - Error (error): An error object that captures any error that occurred during the check. This is not serialized directly to JSON.
//   - ErrorMSG (string): A string representation of the error. This is included because the error datatype cannot be directly serialized to JSON.
//
// The Check struct can be instantiated using the following functions:
//   - NewCheckResult: Creates a new Check instance with only a result.
//   - NewCheckError: Creates a new Check instance with an error and its string representation.
//   - NewCheckErrorf: Creates a new Check instance with a formatted error message and its error object.
//
// This struct is primarily used to standardize the return type across various security and privacy checks in the application.
type Check struct {
	IssueID  int      `json:"issue_id"`
	ResultID int      `json:"result_id"`
	Result   []string `json:"result,omitempty"`
	Error    error    `json:"-"` // Don't serialize error field to JSON
	ErrorMSG string   `json:"error,omitempty"`
}

// TODO: Update documentation
// NewCheckResult is a constructor function that creates and returns a new instance of the Check struct.
// It sets the IssueID, ResultID, and Result fields of the Check struct, leaving the Error and ErrorMSG fields as their zero values.
//
// Parameters:
//   - issID (int): A unique identifier for the issue. This value is assigned to the IssueID field of the Check struct.
//   - resID (int): A unique identifier for the result. This value is assigned to the ResultID field of the Check struct.
//   - result ([]string): The outcome of the check. This could be a list of strings representing various results. This value is assigned to the Result field of the Check struct.
//
// Returns:
//   - Check: A new instance of the Check struct with the IssueID, ResultID, and Result fields set to the provided values, and the Error and ErrorMSG fields set to their zero values.
//
// This function is primarily used when a security or privacy check completes successfully and returns a result without any errors.
func NewCheckResult(issID int, resID int, result ...string) Check {
	return Check{IssueID: issID, ResultID: resID, Result: result}
}

// TODO: Update documentation
// NewCheckError is a constructor function that creates and returns a new instance of the Check struct.
// It sets the ID, Error, and ErrorMSG fields of the Check struct, leaving the Result field as its zero value.
//
// Parameters:
//   - id (int): A unique identifier for the check. This value is assigned to the ID field of the Check struct.
//   - err (error): An error object that captures any error that occurred during the check. This value is assigned to the Error field of the Check struct, and its string representation is assigned to the ErrorMSG field.
//
// Returns:
//   - Check: A new instance of the Check struct with the ID, Error, and ErrorMSG fields set to the provided values, and the Result field set to its zero value.
//
// This function is primarily used when a security or privacy check encounters an error and needs to return a Check instance that encapsulates this error.
func NewCheckError(id int, err error) Check {
	return Check{IssueID: id, ResultID: -1, Error: err, ErrorMSG: err.Error()}
}

// TODO: Update documentation
// NewCheckErrorf is a constructor function that creates and returns a new instance of the Check struct.
// It sets the ID, Error, and ErrorMSG fields of the Check struct, leaving the Result field as its zero value.
//
// Parameters:
//   - id (int): A unique identifier for the check. This value is assigned to the ID field of the Check struct.
//   - message (string): A base error message that provides context about the error. This is used to create a formatted error message.
//   - err (error): An error object that captures any error that occurred during the check. This is used to create a formatted error message, which is assigned to the ErrorMSG field.
//
// Returns:
//   - Check: A new instance of the Check struct with the ID, Error, and ErrorMSG fields set to the provided values, and the Result field set to its zero value.
//
// This function is primarily used when a security or privacy check encounters an error and needs to return a Check instance that encapsulates this error. The formatted error message provides additional context about the error, which can be helpful for debugging and understanding the nature of the error.
func NewCheckErrorf(id int, message string, err error) Check {
	formatErr := fmt.Errorf(message+": %w", err)
	return Check{IssueID: id, ResultID: -1, Error: formatErr, ErrorMSG: formatErr.Error()}
}
