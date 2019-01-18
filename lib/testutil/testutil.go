package testutil

import (
	"os"
)

const testModeFlag = "TEST_MODE"

//GetTestMode return TestMode from environment key
func GetTestMode() (TestMode, error) {
	testMode := os.Getenv(testModeFlag)
	return getTestModeFromString(testMode)
}
