package mathutil

// MinInt64 returns the minimum number between two given params.
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MinUint64 returns the minimum number between two given params.
func MinUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
