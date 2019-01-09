package mathutil

// MintInt64 returns the minimum number between two given params.
func MintInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MintUint64 returns the minimum number between two given params.
func MintUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
