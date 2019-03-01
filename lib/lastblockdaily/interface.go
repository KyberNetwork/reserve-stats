package lastblockdaily

// Interface is the common interface of daily last block resolver.
// A block is considered last block of the day if its next block timestamp is on a
// different day as its timestamp.
type Interface interface {
	// Next yields the next last block of the day.
	Next() (blockNumber uint64, err error)
}
