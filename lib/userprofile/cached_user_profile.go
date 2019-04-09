package userprofile

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// CachedClient is the wrapper of User Profile Client with caching ability.
type CachedClient struct {
	*Client
	mu           *sync.RWMutex
	cached       map[ethereum.Address]UserProfile
	MaxCacheSize int64
	cacheSize    int64
	minHeap      *AddrHeap
}

// NewCachedClient creates a new User Profile cached client instance.
func NewCachedClient(client *Client, maxCache int64) *CachedClient {
	var h = &AddrHeap{}
	heap.Init(h)
	client.sugar.Debugw("Creating cache client ...", "max cache", maxCache)
	return &CachedClient{
		Client:       client,
		mu:           &sync.RWMutex{},
		cached:       make(map[ethereum.Address]UserProfile),
		MaxCacheSize: maxCache,
		cacheSize:    0,
		minHeap:      h,
	}
}

// LookUpCache will lookup the Userprofile from cache
func (cc *CachedClient) LookUpCache(addr ethereum.Address) (UserProfile, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	p, ok := cc.cached[addr]
	return p, ok
}

// LookUpUserProfile will look for the UserProfile of input addr in cache first
// If this fail then it will query from endpoint
func (cc *CachedClient) LookUpUserProfile(addr ethereum.Address) (UserProfile, error) {
	logger := cc.sugar.With(
		"func", "lib/core/CachedClient.Token",
		"address", addr.Hex(),
	)
	p, ok := cc.LookUpCache(addr)
	if ok {
		logger.Debugw("cache hit")
		return p, nil
	}

	logger.Debugf("cache missed. Lookup from API endpoint and caching... Current cache size : %d", cc.cacheSize)
	p, err := cc.Client.LookUpUserProfile(addr)
	if err != nil {
		return p, err
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	//if MaxCacheSize reached, delete the oldest member
	if cc.cacheSize >= cc.MaxCacheSize {
		oldest := heap.Pop(cc.minHeap)
		addrNode, ok := oldest.(AddrNode)
		if !ok {
			return p, fmt.Errorf("cannot assert address node %v from heap", oldest)
		}
		logger.Debugf("removing %s from address heap", addrNode.Addr)
		delete(cc.cached, ethereum.HexToAddress(addrNode.Addr))
		cc.cacheSize--
	}
	//write the result return from API query to heap and cached
	cc.cached[addr] = p
	heap.Push(cc.minHeap, AddrNode{
		Timestamp: timeutil.UnixMilliSecond(),
		Addr:      addr.Hex(),
	})
	cc.cacheSize++
	return p, nil
}
