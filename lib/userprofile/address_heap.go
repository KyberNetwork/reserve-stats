package userprofile

// AddrNode represent a node in heap. It minimize the data that needs to be stored.
type AddrNode struct {
	Timestamp uint64
	Addr      string
}

// AddrHeap is a min-heap of AddrNode
type AddrHeap []AddrNode

// Len return the length of addrHeap
func (a AddrHeap) Len() int {
	return len(a)
}

//Less return bool indicate the comparion between 2 AddrNode by index
func (a AddrHeap) Less(i, j int) bool {
	return a[i].Timestamp < a[j].Timestamp
}

//Swap swap 2 node value
func (a AddrHeap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//Push Add the addrNode to heap
func (a *AddrHeap) Push(x interface{}) {
	*a = append(*a, x.(AddrNode))
}

//Pop return the root and delete it from heap
func (a *AddrHeap) Pop() interface{} {
	old := *a
	n := len(old)
	x := old[n-1]
	*a = old[0 : n-1]
	return x
}
