package common

import (
	"encoding/json"
	"io/ioutil"

	ethereum "github.com/ethereum/go-ethereum/common"
)

//AddrToAppName return the app name according to address
type AddrToAppName map[ethereum.Address]string

func AddrAppNameFromFile(path string) AddrToAppName {
	var (
		result = make(AddrToAppName)
		tmp    = make(map[string]string)
	)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}
	for addr, appName := range tmp {
		result[ethereum.HexToAddress(addr)] = appName
	}
	return result
}
