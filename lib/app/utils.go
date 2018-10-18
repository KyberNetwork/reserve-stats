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

// JoinEnvVar joins the common env prefix with given env var name.
func JoinEnvVar(envVar string) string {
	return joinEnvVar(commonEnvPrefix, envVar)
}
