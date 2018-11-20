package appname

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	ethereum "github.com/ethereum/go-ethereum/common"
)

//MapAddrAppName return the app name according to address
type MapAddrAppName struct {
	mutex   *sync.Mutex
	mapping map[ethereum.Address]string
}

//NewMapAddrAppName return a new instance of MapAddrAppName
func NewMapAddrAppName() *MapAddrAppName {
	var (
		temp = make(map[ethereum.Address]string)
	)
	return &MapAddrAppName{
		mapping: temp,
		mutex:   &sync.Mutex{},
	}
}

//LoadFromFile add mapping addr-appname from a file
func (apn *MapAddrAppName) LoadFromFile(path string) error {
	apn.mutex.Lock()
	defer apn.mutex.Unlock()
	var (
		tmp = make(map[string]string)
	)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	for addr, appName := range tmp {
		apn.mapping[ethereum.HexToAddress(addr)] = appName
	}
	return nil
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
	apn.mutex.Lock()
	defer apn.mutex.Unlock()
	var (
		result = make(map[ethereum.Address]string)
	)
	for k, v := range apn.mapping {
		result[k] = v
	}
	return result
}
