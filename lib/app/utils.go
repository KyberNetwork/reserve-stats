package app

import (
	"encoding/binary"
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

// BytesToUint64 converts the given bytes to a compatible uint64
// value. It is used for reading the uint64 from a stored BoltDB
// record. If returned 0 if the given bytes are incompatible.
func BytesToUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}
