package testutil

import (
	"fmt"
)

// TestMode define a list of test mode for go test
//go:generate stringer -type=TestMode -linecomment
type TestMode int

const (
	//Internal only test all the internal modules without requiring external resources
	Internal TestMode = iota //internal
	//ExternalTest test the modules that requires external resources
	External //external
)

var testModes = map[string]TestMode{
	"internal": Internal,
	"external": External,
}

func getTestModeFromString(key string) (TestMode, error) {
	mode, avail := testModes[key]
	if !avail {
		return 0, fmt.Errorf("mode %s is not available, make sure to set %s ENV", key, testModeFlag)
	}
	return mode, nil
}
