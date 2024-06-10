package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"testing"
)

var testsNotPresent = []func(t *testing.T){
	i.TestIntegrationFirefoxFolderNotExists,
	// TODO: turn back on when the test is fixed
	// i.TestIntegrationOutdatedWin10UpToDate,
	i.TestIntegrationUACPartialEnabled,
}

func TestNotPresent(t *testing.T) {
	for _, test := range testsNotPresent {
		test(t)
	}
}
