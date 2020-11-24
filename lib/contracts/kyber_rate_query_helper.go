// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RateHelperABI is the input ABI used to generate the binding from.
const RateHelperABI = "[{\"inputs\":[{\"internalType\":\"contractIKyberReserve\",\"name\":\"reserve\",\"type\":\"address\"},{\"internalType\":\"contractIERC20[]\",\"name\":\"srcs\",\"type\":\"address[]\"},{\"internalType\":\"contractIERC20[]\",\"name\":\"dests\",\"type\":\"address[]\"}],\"name\":\"getReserveRates\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"pricingRates\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"sanityRates\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// RateHelper is an auto generated Go binding around an Ethereum contract.
type RateHelper struct {
	RateHelperCaller     // Read-only binding to the contract
	RateHelperTransactor // Write-only binding to the contract
	RateHelperFilterer   // Log filterer for contract events
}

// RateHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type RateHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RateHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RateHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RateHelperSession struct {
	Contract     *RateHelper       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RateHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RateHelperCallerSession struct {
	Contract *RateHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// RateHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RateHelperTransactorSession struct {
	Contract     *RateHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// RateHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type RateHelperRaw struct {
	Contract *RateHelper // Generic contract binding to access the raw methods on
}

// RateHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RateHelperCallerRaw struct {
	Contract *RateHelperCaller // Generic read-only contract binding to access the raw methods on
}

// RateHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RateHelperTransactorRaw struct {
	Contract *RateHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRateHelper creates a new instance of RateHelper, bound to a specific deployed contract.
func NewRateHelper(address common.Address, backend bind.ContractBackend) (*RateHelper, error) {
	contract, err := bindRateHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RateHelper{RateHelperCaller: RateHelperCaller{contract: contract}, RateHelperTransactor: RateHelperTransactor{contract: contract}, RateHelperFilterer: RateHelperFilterer{contract: contract}}, nil
}

// NewRateHelperCaller creates a new read-only instance of RateHelper, bound to a specific deployed contract.
func NewRateHelperCaller(address common.Address, caller bind.ContractCaller) (*RateHelperCaller, error) {
	contract, err := bindRateHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RateHelperCaller{contract: contract}, nil
}

// NewRateHelperTransactor creates a new write-only instance of RateHelper, bound to a specific deployed contract.
func NewRateHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*RateHelperTransactor, error) {
	contract, err := bindRateHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RateHelperTransactor{contract: contract}, nil
}

// NewRateHelperFilterer creates a new log filterer instance of RateHelper, bound to a specific deployed contract.
func NewRateHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*RateHelperFilterer, error) {
	contract, err := bindRateHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RateHelperFilterer{contract: contract}, nil
}

// bindRateHelper binds a generic wrapper to an already deployed contract.
func bindRateHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RateHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateHelper *RateHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RateHelper.Contract.RateHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateHelper *RateHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateHelper.Contract.RateHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateHelper *RateHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateHelper.Contract.RateHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateHelper *RateHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RateHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateHelper *RateHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateHelper *RateHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateHelper.Contract.contract.Transact(opts, method, params...)
}

// GetReserveRates is a free data retrieval call binding the contract method 0xcd85afd4.
//
// Solidity: function getReserveRates(address reserve, address[] srcs, address[] dests) view returns(uint256[] pricingRates, uint256[] sanityRates)
func (_RateHelper *RateHelperCaller) GetReserveRates(opts *bind.CallOpts, reserve common.Address, srcs []common.Address, dests []common.Address) (struct {
	PricingRates []*big.Int
	SanityRates  []*big.Int
}, error) {
	ret := new(struct {
		PricingRates []*big.Int
		SanityRates  []*big.Int
	})
	out := ret
	err := _RateHelper.contract.Call(opts, out, "getReserveRates", reserve, srcs, dests)
	return *ret, err
}

// GetReserveRates is a free data retrieval call binding the contract method 0xcd85afd4.
//
// Solidity: function getReserveRates(address reserve, address[] srcs, address[] dests) view returns(uint256[] pricingRates, uint256[] sanityRates)
func (_RateHelper *RateHelperSession) GetReserveRates(reserve common.Address, srcs []common.Address, dests []common.Address) (struct {
	PricingRates []*big.Int
	SanityRates  []*big.Int
}, error) {
	return _RateHelper.Contract.GetReserveRates(&_RateHelper.CallOpts, reserve, srcs, dests)
}

// GetReserveRates is a free data retrieval call binding the contract method 0xcd85afd4.
//
// Solidity: function getReserveRates(address reserve, address[] srcs, address[] dests) view returns(uint256[] pricingRates, uint256[] sanityRates)
func (_RateHelper *RateHelperCallerSession) GetReserveRates(reserve common.Address, srcs []common.Address, dests []common.Address) (struct {
	PricingRates []*big.Int
	SanityRates  []*big.Int
}, error) {
	return _RateHelper.Contract.GetReserveRates(&_RateHelper.CallOpts, reserve, srcs, dests)
}
