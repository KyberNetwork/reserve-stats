# Ethereum Smart Contracts Developer Guide

## Generate Golang Binding for Ethereum Smart Contracts

Dependencies:

	- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)

Copy the minified contract abi JSON to a file in this repository
lib/contracts directory. The file name must have ".abi" extension.

Example: internal_network.abi

Append a line to abi_gen.sh with following format.

```sh
abigen -abi "$OLDPWD"/<filename>.abi -pkg contracts -type <struct_name> -out "$OLDPWD"/<filename>_abi.go
```

Execute abi_gen.sh, the new .go file should be generated.

## Calling Smart Contract with Abitrary Block Number

Following example contains a function that allows calling the
GetReserves method at any block number.

Note: use https://github.com/favadi/go-ethereum/tree/v1.8.16-fork
until https://github.com/ethereum/go-ethereum/pull/17770 is merged and
released.

Example:

```go
func getReserves(blockNumber *big.Int) ([]common.Address, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		return nil, err
	}
	internalNetwork, err := contracts.NewInternalNetwork(
		common.HexToAddress(contracts.InternalNetworkContractAddress),
		client)
	if err != nil {
		return nil, err
	}
	return internalNetwork.GetReserves(&bind.CallOpts{BlockNumber: blockNumber})
}
```
