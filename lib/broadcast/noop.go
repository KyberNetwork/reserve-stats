package broadcast

// Noop is an implementation of broadcast client that returns no data,
// to be used when broadcast service is not available.
type Noop struct{}

// NewNoop creates a new Noop instance.
func NewNoop() *Noop {
	return &Noop{}
}

func (*Noop) GetTxInfo(tx string) (string, string, string, error) {
	return "", "", "", nil
}
