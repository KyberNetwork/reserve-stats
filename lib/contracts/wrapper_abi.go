// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// WrapperABI is the input ABI used to generate the binding from.
const WrapperABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"bytes14\"},{\"name\":\"byteInd\",\"type\":\"uint256\"}],\"name\":\"getInt8FromByte\",\"outputs\":[{\"name\":\"\",\"type\":\"int8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getBalances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ratesContract\",\"type\":\"address\"},{\"name\":\"tokenList\",\"type\":\"address[]\"}],\"name\":\"getTokenIndicies\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\"},{\"name\":\"srcs\",\"type\":\"address[]\"},{\"name\":\"dests\",\"type\":\"address[]\"}],\"name\":\"getReserveRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"bytes14\"},{\"name\":\"byteInd\",\"type\":\"uint256\"}],\"name\":\"getByteFromBytes14\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"network\",\"type\":\"address\"},{\"name\":\"srcs\",\"type\":\"address[]\"},{\"name\":\"dests\",\"type\":\"address[]\"},{\"name\":\"qty\",\"type\":\"uint256[]\"}],\"name\":\"getExpectedRates\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ratesContract\",\"type\":\"address\"},{\"name\":\"tokenList\",\"type\":\"address[]\"}],\"name\":\"getTokenRates\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"int8[]\"},{\"name\":\"\",\"type\":\"int8[]\"},{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Wrapper is an auto generated Go binding around an Ethereum contract.
type Wrapper struct {
	WrapperCaller     // Read-only binding to the contract
	WrapperTransactor // Write-only binding to the contract
	WrapperFilterer   // Log filterer for contract events
}

// WrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type WrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrapperSession struct {
	Contract     *Wrapper          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrapperCallerSession struct {
	Contract *WrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// WrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrapperTransactorSession struct {
	Contract     *WrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// WrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type WrapperRaw struct {
	Contract *Wrapper // Generic contract binding to access the raw methods on
}

// WrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrapperCallerRaw struct {
	Contract *WrapperCaller // Generic read-only contract binding to access the raw methods on
}

// WrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrapperTransactorRaw struct {
	Contract *WrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWrapper creates a new instance of Wrapper, bound to a specific deployed contract.
func NewWrapper(address common.Address, backend bind.ContractBackend) (*Wrapper, error) {
	contract, err := bindWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Wrapper{WrapperCaller: WrapperCaller{contract: contract}, WrapperTransactor: WrapperTransactor{contract: contract}, WrapperFilterer: WrapperFilterer{contract: contract}}, nil
}

// NewWrapperCaller creates a new read-only instance of Wrapper, bound to a specific deployed contract.
func NewWrapperCaller(address common.Address, caller bind.ContractCaller) (*WrapperCaller, error) {
	contract, err := bindWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrapperCaller{contract: contract}, nil
}

// NewWrapperTransactor creates a new write-only instance of Wrapper, bound to a specific deployed contract.
func NewWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*WrapperTransactor, error) {
	contract, err := bindWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrapperTransactor{contract: contract}, nil
}

// NewWrapperFilterer creates a new log filterer instance of Wrapper, bound to a specific deployed contract.
func NewWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*WrapperFilterer, error) {
	contract, err := bindWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrapperFilterer{contract: contract}, nil
}

// bindWrapper binds a generic wrapper to an already deployed contract.
func bindWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wrapper *WrapperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Wrapper.Contract.WrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wrapper *WrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wrapper.Contract.WrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wrapper *WrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wrapper.Contract.WrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wrapper *WrapperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Wrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wrapper *WrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wrapper *WrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wrapper.Contract.contract.Transact(opts, method, params...)
}

// GetBalances is a free data retrieval call binding the contract method 0x6a385ae9.
//
// Solidity: function getBalances(reserve address, tokens address[]) constant returns(uint256[])
func (_Wrapper *WrapperCaller) GetBalances(opts *bind.CallOpts, reserve common.Address, tokens []common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _Wrapper.contract.Call(opts, out, "getBalances", reserve, tokens)
	return *ret0, err
}

// GetBalances is a free data retrieval call binding the contract method 0x6a385ae9.
//
// Solidity: function getBalances(reserve address, tokens address[]) constant returns(uint256[])
func (_Wrapper *WrapperSession) GetBalances(reserve common.Address, tokens []common.Address) ([]*big.Int, error) {
	return _Wrapper.Contract.GetBalances(&_Wrapper.CallOpts, reserve, tokens)
}

// GetBalances is a free data retrieval call binding the contract method 0x6a385ae9.
//
// Solidity: function getBalances(reserve address, tokens address[]) constant returns(uint256[])
func (_Wrapper *WrapperCallerSession) GetBalances(reserve common.Address, tokens []common.Address) ([]*big.Int, error) {
	return _Wrapper.Contract.GetBalances(&_Wrapper.CallOpts, reserve, tokens)
}

// GetByteFromBytes14 is a free data retrieval call binding the contract method 0xa609f034.
//
// Solidity: function getByteFromBytes14(x bytes14, byteInd uint256) constant returns(bytes1)
func (_Wrapper *WrapperCaller) GetByteFromBytes14(opts *bind.CallOpts, x [14]byte, byteInd *big.Int) ([1]byte, error) {
	var (
		ret0 = new([1]byte)
	)
	out := ret0
	err := _Wrapper.contract.Call(opts, out, "getByteFromBytes14", x, byteInd)
	return *ret0, err
}

// GetByteFromBytes14 is a free data retrieval call binding the contract method 0xa609f034.
//
// Solidity: function getByteFromBytes14(x bytes14, byteInd uint256) constant returns(bytes1)
func (_Wrapper *WrapperSession) GetByteFromBytes14(x [14]byte, byteInd *big.Int) ([1]byte, error) {
	return _Wrapper.Contract.GetByteFromBytes14(&_Wrapper.CallOpts, x, byteInd)
}

// GetByteFromBytes14 is a free data retrieval call binding the contract method 0xa609f034.
//
// Solidity: function getByteFromBytes14(x bytes14, byteInd uint256) constant returns(bytes1)
func (_Wrapper *WrapperCallerSession) GetByteFromBytes14(x [14]byte, byteInd *big.Int) ([1]byte, error) {
	return _Wrapper.Contract.GetByteFromBytes14(&_Wrapper.CallOpts, x, byteInd)
}

// GetExpectedRates is a free data retrieval call binding the contract method 0xf1838fe4.
//
// Solidity: function getExpectedRates(network address, srcs address[], dests address[], qty uint256[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCaller) GetExpectedRates(opts *bind.CallOpts, network common.Address, srcs []common.Address, dests []common.Address, qty []*big.Int) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Wrapper.contract.Call(opts, out, "getExpectedRates", network, srcs, dests, qty)
	return *ret0, *ret1, err
}

// GetExpectedRates is a free data retrieval call binding the contract method 0xf1838fe4.
//
// Solidity: function getExpectedRates(network address, srcs address[], dests address[], qty uint256[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperSession) GetExpectedRates(network common.Address, srcs []common.Address, dests []common.Address, qty []*big.Int) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetExpectedRates(&_Wrapper.CallOpts, network, srcs, dests, qty)
}

// GetExpectedRates is a free data retrieval call binding the contract method 0xf1838fe4.
//
// Solidity: function getExpectedRates(network address, srcs address[], dests address[], qty uint256[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCallerSession) GetExpectedRates(network common.Address, srcs []common.Address, dests []common.Address, qty []*big.Int) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetExpectedRates(&_Wrapper.CallOpts, network, srcs, dests, qty)
}

// GetInt8FromByte is a free data retrieval call binding the contract method 0x67c33c80.
//
// Solidity: function getInt8FromByte(x bytes14, byteInd uint256) constant returns(int8)
func (_Wrapper *WrapperCaller) GetInt8FromByte(opts *bind.CallOpts, x [14]byte, byteInd *big.Int) (int8, error) {
	var (
		ret0 = new(int8)
	)
	out := ret0
	err := _Wrapper.contract.Call(opts, out, "getInt8FromByte", x, byteInd)
	return *ret0, err
}

// GetInt8FromByte is a free data retrieval call binding the contract method 0x67c33c80.
//
// Solidity: function getInt8FromByte(x bytes14, byteInd uint256) constant returns(int8)
func (_Wrapper *WrapperSession) GetInt8FromByte(x [14]byte, byteInd *big.Int) (int8, error) {
	return _Wrapper.Contract.GetInt8FromByte(&_Wrapper.CallOpts, x, byteInd)
}

// GetInt8FromByte is a free data retrieval call binding the contract method 0x67c33c80.
//
// Solidity: function getInt8FromByte(x bytes14, byteInd uint256) constant returns(int8)
func (_Wrapper *WrapperCallerSession) GetInt8FromByte(x [14]byte, byteInd *big.Int) (int8, error) {
	return _Wrapper.Contract.GetInt8FromByte(&_Wrapper.CallOpts, x, byteInd)
}

// GetReserveRate is a free data retrieval call binding the contract method 0x91eb1c69.
//
// Solidity: function getReserveRate(reserve address, srcs address[], dests address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCaller) GetReserveRate(opts *bind.CallOpts, reserve common.Address, srcs []common.Address, dests []common.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Wrapper.contract.Call(opts, out, "getReserveRate", reserve, srcs, dests)
	return *ret0, *ret1, err
}

// GetReserveRate is a free data retrieval call binding the contract method 0x91eb1c69.
//
// Solidity: function getReserveRate(reserve address, srcs address[], dests address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperSession) GetReserveRate(reserve common.Address, srcs []common.Address, dests []common.Address) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetReserveRate(&_Wrapper.CallOpts, reserve, srcs, dests)
}

// GetReserveRate is a free data retrieval call binding the contract method 0x91eb1c69.
//
// Solidity: function getReserveRate(reserve address, srcs address[], dests address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCallerSession) GetReserveRate(reserve common.Address, srcs []common.Address, dests []common.Address) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetReserveRate(&_Wrapper.CallOpts, reserve, srcs, dests)
}

// GetTokenIndicies is a free data retrieval call binding the contract method 0x7c80feff.
//
// Solidity: function getTokenIndicies(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCaller) GetTokenIndicies(opts *bind.CallOpts, ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Wrapper.contract.Call(opts, out, "getTokenIndicies", ratesContract, tokenList)
	return *ret0, *ret1, err
}

// GetTokenIndicies is a free data retrieval call binding the contract method 0x7c80feff.
//
// Solidity: function getTokenIndicies(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperSession) GetTokenIndicies(ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetTokenIndicies(&_Wrapper.CallOpts, ratesContract, tokenList)
}

// GetTokenIndicies is a free data retrieval call binding the contract method 0x7c80feff.
//
// Solidity: function getTokenIndicies(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[])
func (_Wrapper *WrapperCallerSession) GetTokenIndicies(ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, error) {
	return _Wrapper.Contract.GetTokenIndicies(&_Wrapper.CallOpts, ratesContract, tokenList)
}

// GetTokenRates is a free data retrieval call binding the contract method 0xf37f8345.
//
// Solidity: function getTokenRates(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[], int8[], int8[], uint256[])
func (_Wrapper *WrapperCaller) GetTokenRates(opts *bind.CallOpts, ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, []int8, []int8, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
		ret2 = new([]int8)
		ret3 = new([]int8)
		ret4 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Wrapper.contract.Call(opts, out, "getTokenRates", ratesContract, tokenList)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// GetTokenRates is a free data retrieval call binding the contract method 0xf37f8345.
//
// Solidity: function getTokenRates(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[], int8[], int8[], uint256[])
func (_Wrapper *WrapperSession) GetTokenRates(ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, []int8, []int8, []*big.Int, error) {
	return _Wrapper.Contract.GetTokenRates(&_Wrapper.CallOpts, ratesContract, tokenList)
}

// GetTokenRates is a free data retrieval call binding the contract method 0xf37f8345.
//
// Solidity: function getTokenRates(ratesContract address, tokenList address[]) constant returns(uint256[], uint256[], int8[], int8[], uint256[])
func (_Wrapper *WrapperCallerSession) GetTokenRates(ratesContract common.Address, tokenList []common.Address) ([]*big.Int, []*big.Int, []int8, []int8, []*big.Int, error) {
	return _Wrapper.Contract.GetTokenRates(&_Wrapper.CallOpts, ratesContract, tokenList)
}
