package common

import (
	"encoding/json"
	"io/ioutil"

	ethereum "github.com/ethereum/go-ethereum/common"
)

//AddrToAppName return the app name according to address
type AddrToAppName map[ethereum.Address]string

//AddrAppNameFromFile return a map WalletAddr to AppName from a file
func AddrAppNameFromFile(path string) (AddrToAppName, error) {
	var (
		result = make(AddrToAppName)
		tmp    = make(map[string]string)
	)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}
	for addr, appName := range tmp {
		result[ethereum.HexToAddress(addr)] = appName
	}
	return result, nil
}
