package testutil

import (
	"os"
	"strconv"
	"testing"
)

const externalTestFlag = "EXTERNAL_TEST"

// SkipExternal will skip the test if mode is not external
func SkipExternal(t *testing.T) {
	external, err := strconv.ParseBool(os.Getenv(externalTestFlag))
	if err != nil && !external {
		t.Skip("disable as this test require external resource")
	}
}
