package integration_test

import (
	i "github.com/InfoSec-Agent/InfoSec-Agent/backend/integration_testing"
	"testing"
)

var testsOther = []func(t *testing.T){
	i.TestIntegrationDefenderRealTimeActive,
}

func TestOther(t *testing.T) {
	for _, test := range testsOther {
		test(t)
	}
}
