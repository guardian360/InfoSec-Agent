package checks

import "fmt"

type Check struct {
	Id       string   `json:"id"`
	Result   []string `json:"result,omitempty"`
	Error    error    `json:"-"` // Don't serialize error field to JSON
	ErrorMSG string   `json:"error,omitempty"`
}

func newCheckResult(id string, result ...string) Check {
	return Check{Id: id, Result: result}
}

func newCheckError(id string, err error) Check {
	return Check{Id: id, Error: err, ErrorMSG: err.Error()}
}

func newCheckErrorf(id string, message string, err error) Check {
	formatErr := fmt.Errorf(message+": %w", err)
	return Check{Id: id, Error: formatErr, ErrorMSG: formatErr.Error()}
}
