package app

import (
	"fmt"
	"strings"
)

const commonEnvPrefix = "RESERVE_STATS"

// joinEnvVar joins the environment variable with the prefix.
func joinEnvVar(prefix, envVar string) string {
	prefix = strings.TrimRight(prefix, "_")
	envVar = strings.TrimLeft(envVar, "_")
	return fmt.Sprintf("%s_%s", prefix, envVar)
}
