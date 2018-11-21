package appname

import (
	"encoding/json"
	"os"
	"sync"

	ethereum "github.com/ethereum/go-ethereum/common"
)

//MapAddrAppName return the app name according to address
type MapAddrAppName struct {
	mutex   *sync.RWMutex
	mapping map[ethereum.Address]string
}

//Option sets the initialization behavior MappAddrAppName inSTANCE
type Option func(m *MapAddrAppName)

//NewMapAddrAppName return a new instance of MapAddrAppName
func NewMapAddrAppName(options ...Option) *MapAddrAppName {
	var m = &MapAddrAppName{
		mapping: make(map[ethereum.Address]string),
		mutex:   &sync.RWMutex{},
	}
	for _, opt := range options {
		opt(m)
	}
	return m
}

// WithDataFile returns an option to read data in from a JSON file
func WithDataFile(path string) Option {
	return func(apn *MapAddrAppName) {
		apn.mutex.Lock()
		defer apn.mutex.Unlock()
		var (
			tmp = make(map[string]string)
		)
		fd, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		if err = json.NewDecoder(fd).Decode(&tmp); err != nil {
			panic(err)
		}
		for addr, appName := range tmp {
			apn.mapping[ethereum.HexToAddress(addr)] = appName
		}
	}
}

//UpdateMapAddrAppName update an address map to a name
func (apn *MapAddrAppName) UpdateMapAddrAppName(addr ethereum.Address, name string) error {
	apn.mutex.Lock()
	defer apn.mutex.Unlock()
	apn.mapping[addr] = name
	return nil
}

//GetAddrToAppName return a map of address to app name
func (apn *MapAddrAppName) GetAddrToAppName() map[ethereum.Address]string {
	apn.mutex.RLock()
	defer apn.mutex.RUnlock()
	var (
		result = make(map[ethereum.Address]string)
	)
	for k, v := range apn.mapping {
		result[k] = v
	}
	return result
}
