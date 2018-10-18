package boltutils

import "encoding/binary"

// Uint64ToBytes converts the given uint64 to a stream of bytes to
// store in BoltDB.
func Uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
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
