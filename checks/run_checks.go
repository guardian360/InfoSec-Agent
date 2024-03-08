package checks

type Check struct {
	Id     string
	Result []string
	Error  error
}

func newCheckResult(id string, result []string) Check {
	return Check{Id: id, Result: result, Error: nil}
}

func newCheckError(id string, err error) Check {
	return Check{Id: id, Result: []string{}, Error: err}
}