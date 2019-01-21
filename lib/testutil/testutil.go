package testutil

import (
	"os"
	"testing"
)

const testModeFlag = "TEST_MODE"

//GetTestMode return TestMode from environment key
func GetTestMode() (TestMode, error) {
	testMode := os.Getenv(testModeFlag)
	return getTestModeFromString(testMode)
}

func SkipExternal(t *testing.T) {
	testMode, err := GetTestMode()
	if err != nil {
		t.Skip("Can't get test mode. skip this external test")
	}
	if testMode != External {
		t.Skip("disable as this test require external resource")
	}
}
