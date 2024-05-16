package integration_testing

import (
	"testing"
)

var testsOther = []func(t *testing.T){
	TestIntegrationDefenderRealTimeActive,
}

func TestOther(t *testing.T) {
	for _, test := range testsOther {
		test(t)
	}
}
