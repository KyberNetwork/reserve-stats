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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SanityRatesABI is the input ABI used to generate the binding from.
const SanityRatesABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reasonableDiffInBps\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"srcs\",\"type\":\"address[]\"},{\"name\":\"diff\",\"type\":\"uint256[]\"}],\"name\":\"setReasonableDiff\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"}],\"name\":\"getSanityRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"srcs\",\"type\":\"address[]\"},{\"name\":\"rates\",\"type\":\"uint256[]\"}],\"name\":\"setSanityRates\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"}]"

// SanityRates is an auto generated Go binding around an Ethereum contract.
type SanityRates struct {
	SanityRatesCaller     // Read-only binding to the contract
	SanityRatesTransactor // Write-only binding to the contract
	SanityRatesFilterer   // Log filterer for contract events
}

// SanityRatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type SanityRatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SanityRatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SanityRatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SanityRatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SanityRatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SanityRatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SanityRatesSession struct {
	Contract     *SanityRates      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SanityRatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SanityRatesCallerSession struct {
	Contract *SanityRatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SanityRatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SanityRatesTransactorSession struct {
	Contract     *SanityRatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SanityRatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type SanityRatesRaw struct {
	Contract *SanityRates // Generic contract binding to access the raw methods on
}

// SanityRatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SanityRatesCallerRaw struct {
	Contract *SanityRatesCaller // Generic read-only contract binding to access the raw methods on
}

// SanityRatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SanityRatesTransactorRaw struct {
	Contract *SanityRatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSanityRates creates a new instance of SanityRates, bound to a specific deployed contract.
func NewSanityRates(address common.Address, backend bind.ContractBackend) (*SanityRates, error) {
	contract, err := bindSanityRates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SanityRates{SanityRatesCaller: SanityRatesCaller{contract: contract}, SanityRatesTransactor: SanityRatesTransactor{contract: contract}, SanityRatesFilterer: SanityRatesFilterer{contract: contract}}, nil
}

// NewSanityRatesCaller creates a new read-only instance of SanityRates, bound to a specific deployed contract.
func NewSanityRatesCaller(address common.Address, caller bind.ContractCaller) (*SanityRatesCaller, error) {
	contract, err := bindSanityRates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SanityRatesCaller{contract: contract}, nil
}

// NewSanityRatesTransactor creates a new write-only instance of SanityRates, bound to a specific deployed contract.
func NewSanityRatesTransactor(address common.Address, transactor bind.ContractTransactor) (*SanityRatesTransactor, error) {
	contract, err := bindSanityRates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SanityRatesTransactor{contract: contract}, nil
}

// NewSanityRatesFilterer creates a new log filterer instance of SanityRates, bound to a specific deployed contract.
func NewSanityRatesFilterer(address common.Address, filterer bind.ContractFilterer) (*SanityRatesFilterer, error) {
	contract, err := bindSanityRates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SanityRatesFilterer{contract: contract}, nil
}

// bindSanityRates binds a generic wrapper to an already deployed contract.
func bindSanityRates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SanityRatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SanityRates *SanityRatesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SanityRates.Contract.SanityRatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SanityRates *SanityRatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SanityRates.Contract.SanityRatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SanityRates *SanityRatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SanityRates.Contract.SanityRatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SanityRates *SanityRatesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SanityRates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SanityRates *SanityRatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SanityRates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SanityRates *SanityRatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SanityRates.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_SanityRates *SanityRatesCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_SanityRates *SanityRatesSession) Admin() (common.Address, error) {
	return _SanityRates.Contract.Admin(&_SanityRates.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_SanityRates *SanityRatesCallerSession) Admin() (common.Address, error) {
	return _SanityRates.Contract.Admin(&_SanityRates.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_SanityRates *SanityRatesCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_SanityRates *SanityRatesSession) GetAlerters() ([]common.Address, error) {
	return _SanityRates.Contract.GetAlerters(&_SanityRates.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_SanityRates *SanityRatesCallerSession) GetAlerters() ([]common.Address, error) {
	return _SanityRates.Contract.GetAlerters(&_SanityRates.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_SanityRates *SanityRatesCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_SanityRates *SanityRatesSession) GetOperators() ([]common.Address, error) {
	return _SanityRates.Contract.GetOperators(&_SanityRates.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_SanityRates *SanityRatesCallerSession) GetOperators() ([]common.Address, error) {
	return _SanityRates.Contract.GetOperators(&_SanityRates.CallOpts)
}

// GetSanityRate is a free data retrieval call binding the contract method 0xa58092b7.
//
// Solidity: function getSanityRate(src address, dest address) constant returns(uint256)
func (_SanityRates *SanityRatesCaller) GetSanityRate(opts *bind.CallOpts, src common.Address, dest common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "getSanityRate", src, dest)
	return *ret0, err
}

// GetSanityRate is a free data retrieval call binding the contract method 0xa58092b7.
//
// Solidity: function getSanityRate(src address, dest address) constant returns(uint256)
func (_SanityRates *SanityRatesSession) GetSanityRate(src common.Address, dest common.Address) (*big.Int, error) {
	return _SanityRates.Contract.GetSanityRate(&_SanityRates.CallOpts, src, dest)
}

// GetSanityRate is a free data retrieval call binding the contract method 0xa58092b7.
//
// Solidity: function getSanityRate(src address, dest address) constant returns(uint256)
func (_SanityRates *SanityRatesCallerSession) GetSanityRate(src common.Address, dest common.Address) (*big.Int, error) {
	return _SanityRates.Contract.GetSanityRate(&_SanityRates.CallOpts, src, dest)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_SanityRates *SanityRatesCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_SanityRates *SanityRatesSession) PendingAdmin() (common.Address, error) {
	return _SanityRates.Contract.PendingAdmin(&_SanityRates.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_SanityRates *SanityRatesCallerSession) PendingAdmin() (common.Address, error) {
	return _SanityRates.Contract.PendingAdmin(&_SanityRates.CallOpts)
}

// ReasonableDiffInBps is a free data retrieval call binding the contract method 0x5463a2e4.
//
// Solidity: function reasonableDiffInBps( address) constant returns(uint256)
func (_SanityRates *SanityRatesCaller) ReasonableDiffInBps(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "reasonableDiffInBps", arg0)
	return *ret0, err
}

// ReasonableDiffInBps is a free data retrieval call binding the contract method 0x5463a2e4.
//
// Solidity: function reasonableDiffInBps( address) constant returns(uint256)
func (_SanityRates *SanityRatesSession) ReasonableDiffInBps(arg0 common.Address) (*big.Int, error) {
	return _SanityRates.Contract.ReasonableDiffInBps(&_SanityRates.CallOpts, arg0)
}

// ReasonableDiffInBps is a free data retrieval call binding the contract method 0x5463a2e4.
//
// Solidity: function reasonableDiffInBps( address) constant returns(uint256)
func (_SanityRates *SanityRatesCallerSession) ReasonableDiffInBps(arg0 common.Address) (*big.Int, error) {
	return _SanityRates.Contract.ReasonableDiffInBps(&_SanityRates.CallOpts, arg0)
}

// TokenRate is a free data retrieval call binding the contract method 0xc57fbf90.
//
// Solidity: function tokenRate( address) constant returns(uint256)
func (_SanityRates *SanityRatesCaller) TokenRate(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SanityRates.contract.Call(opts, out, "tokenRate", arg0)
	return *ret0, err
}

// TokenRate is a free data retrieval call binding the contract method 0xc57fbf90.
//
// Solidity: function tokenRate( address) constant returns(uint256)
func (_SanityRates *SanityRatesSession) TokenRate(arg0 common.Address) (*big.Int, error) {
	return _SanityRates.Contract.TokenRate(&_SanityRates.CallOpts, arg0)
}

// TokenRate is a free data retrieval call binding the contract method 0xc57fbf90.
//
// Solidity: function tokenRate( address) constant returns(uint256)
func (_SanityRates *SanityRatesCallerSession) TokenRate(arg0 common.Address) (*big.Int, error) {
	return _SanityRates.Contract.TokenRate(&_SanityRates.CallOpts, arg0)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_SanityRates *SanityRatesTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_SanityRates *SanityRatesSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.AddAlerter(&_SanityRates.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_SanityRates *SanityRatesTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.AddAlerter(&_SanityRates.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_SanityRates *SanityRatesTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_SanityRates *SanityRatesSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.AddOperator(&_SanityRates.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_SanityRates *SanityRatesTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.AddOperator(&_SanityRates.TransactOpts, newOperator)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_SanityRates *SanityRatesTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_SanityRates *SanityRatesSession) ClaimAdmin() (*types.Transaction, error) {
	return _SanityRates.Contract.ClaimAdmin(&_SanityRates.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_SanityRates *SanityRatesTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _SanityRates.Contract.ClaimAdmin(&_SanityRates.TransactOpts)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_SanityRates *SanityRatesTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_SanityRates *SanityRatesSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.RemoveAlerter(&_SanityRates.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_SanityRates *SanityRatesTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.RemoveAlerter(&_SanityRates.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_SanityRates *SanityRatesTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_SanityRates *SanityRatesSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.RemoveOperator(&_SanityRates.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_SanityRates *SanityRatesTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.RemoveOperator(&_SanityRates.TransactOpts, operator)
}

// SetReasonableDiff is a paid mutator transaction binding the contract method 0x5c53ec59.
//
// Solidity: function setReasonableDiff(srcs address[], diff uint256[]) returns()
func (_SanityRates *SanityRatesTransactor) SetReasonableDiff(opts *bind.TransactOpts, srcs []common.Address, diff []*big.Int) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "setReasonableDiff", srcs, diff)
}

// SetReasonableDiff is a paid mutator transaction binding the contract method 0x5c53ec59.
//
// Solidity: function setReasonableDiff(srcs address[], diff uint256[]) returns()
func (_SanityRates *SanityRatesSession) SetReasonableDiff(srcs []common.Address, diff []*big.Int) (*types.Transaction, error) {
	return _SanityRates.Contract.SetReasonableDiff(&_SanityRates.TransactOpts, srcs, diff)
}

// SetReasonableDiff is a paid mutator transaction binding the contract method 0x5c53ec59.
//
// Solidity: function setReasonableDiff(srcs address[], diff uint256[]) returns()
func (_SanityRates *SanityRatesTransactorSession) SetReasonableDiff(srcs []common.Address, diff []*big.Int) (*types.Transaction, error) {
	return _SanityRates.Contract.SetReasonableDiff(&_SanityRates.TransactOpts, srcs, diff)
}

// SetSanityRates is a paid mutator transaction binding the contract method 0xf5db370f.
//
// Solidity: function setSanityRates(srcs address[], rates uint256[]) returns()
func (_SanityRates *SanityRatesTransactor) SetSanityRates(opts *bind.TransactOpts, srcs []common.Address, rates []*big.Int) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "setSanityRates", srcs, rates)
}

// SetSanityRates is a paid mutator transaction binding the contract method 0xf5db370f.
//
// Solidity: function setSanityRates(srcs address[], rates uint256[]) returns()
func (_SanityRates *SanityRatesSession) SetSanityRates(srcs []common.Address, rates []*big.Int) (*types.Transaction, error) {
	return _SanityRates.Contract.SetSanityRates(&_SanityRates.TransactOpts, srcs, rates)
}

// SetSanityRates is a paid mutator transaction binding the contract method 0xf5db370f.
//
// Solidity: function setSanityRates(srcs address[], rates uint256[]) returns()
func (_SanityRates *SanityRatesTransactorSession) SetSanityRates(srcs []common.Address, rates []*big.Int) (*types.Transaction, error) {
	return _SanityRates.Contract.SetSanityRates(&_SanityRates.TransactOpts, srcs, rates)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_SanityRates *SanityRatesTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_SanityRates *SanityRatesSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.TransferAdmin(&_SanityRates.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_SanityRates *SanityRatesTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.TransferAdmin(&_SanityRates.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_SanityRates *SanityRatesTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_SanityRates *SanityRatesSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.TransferAdminQuickly(&_SanityRates.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_SanityRates *SanityRatesTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.TransferAdminQuickly(&_SanityRates.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.WithdrawEther(&_SanityRates.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.WithdrawEther(&_SanityRates.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.WithdrawToken(&_SanityRates.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_SanityRates *SanityRatesTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _SanityRates.Contract.WithdrawToken(&_SanityRates.TransactOpts, token, amount, sendTo)
}

// SanityRatesAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the SanityRates contract.
type SanityRatesAdminClaimedIterator struct {
	Event *SanityRatesAdminClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesAdminClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesAdminClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesAdminClaimed represents a AdminClaimed event raised by the SanityRates contract.
type SanityRatesAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: e AdminClaimed(newAdmin address, previousAdmin address)
func (_SanityRates *SanityRatesFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*SanityRatesAdminClaimedIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &SanityRatesAdminClaimedIterator{contract: _SanityRates.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: e AdminClaimed(newAdmin address, previousAdmin address)
func (_SanityRates *SanityRatesFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *SanityRatesAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesAdminClaimed)
				if err := _SanityRates.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SanityRatesAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the SanityRates contract.
type SanityRatesAlerterAddedIterator struct {
	Event *SanityRatesAlerterAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesAlerterAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesAlerterAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesAlerterAdded represents a AlerterAdded event raised by the SanityRates contract.
type SanityRatesAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: e AlerterAdded(newAlerter address, isAdd bool)
func (_SanityRates *SanityRatesFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*SanityRatesAlerterAddedIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &SanityRatesAlerterAddedIterator{contract: _SanityRates.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: e AlerterAdded(newAlerter address, isAdd bool)
func (_SanityRates *SanityRatesFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *SanityRatesAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesAlerterAdded)
				if err := _SanityRates.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SanityRatesEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the SanityRates contract.
type SanityRatesEtherWithdrawIterator struct {
	Event *SanityRatesEtherWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesEtherWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesEtherWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesEtherWithdraw represents a EtherWithdraw event raised by the SanityRates contract.
type SanityRatesEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: e EtherWithdraw(amount uint256, sendTo address)
func (_SanityRates *SanityRatesFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*SanityRatesEtherWithdrawIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &SanityRatesEtherWithdrawIterator{contract: _SanityRates.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: e EtherWithdraw(amount uint256, sendTo address)
func (_SanityRates *SanityRatesFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *SanityRatesEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesEtherWithdraw)
				if err := _SanityRates.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SanityRatesOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the SanityRates contract.
type SanityRatesOperatorAddedIterator struct {
	Event *SanityRatesOperatorAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesOperatorAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesOperatorAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesOperatorAdded represents a OperatorAdded event raised by the SanityRates contract.
type SanityRatesOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: e OperatorAdded(newOperator address, isAdd bool)
func (_SanityRates *SanityRatesFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*SanityRatesOperatorAddedIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &SanityRatesOperatorAddedIterator{contract: _SanityRates.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: e OperatorAdded(newOperator address, isAdd bool)
func (_SanityRates *SanityRatesFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *SanityRatesOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesOperatorAdded)
				if err := _SanityRates.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SanityRatesTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the SanityRates contract.
type SanityRatesTokenWithdrawIterator struct {
	Event *SanityRatesTokenWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesTokenWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesTokenWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesTokenWithdraw represents a TokenWithdraw event raised by the SanityRates contract.
type SanityRatesTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: e TokenWithdraw(token address, amount uint256, sendTo address)
func (_SanityRates *SanityRatesFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*SanityRatesTokenWithdrawIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &SanityRatesTokenWithdrawIterator{contract: _SanityRates.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: e TokenWithdraw(token address, amount uint256, sendTo address)
func (_SanityRates *SanityRatesFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *SanityRatesTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesTokenWithdraw)
				if err := _SanityRates.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SanityRatesTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the SanityRates contract.
type SanityRatesTransferAdminPendingIterator struct {
	Event *SanityRatesTransferAdminPending // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SanityRatesTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SanityRatesTransferAdminPending)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SanityRatesTransferAdminPending)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SanityRatesTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SanityRatesTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SanityRatesTransferAdminPending represents a TransferAdminPending event raised by the SanityRates contract.
type SanityRatesTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: e TransferAdminPending(pendingAdmin address)
func (_SanityRates *SanityRatesFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*SanityRatesTransferAdminPendingIterator, error) {

	logs, sub, err := _SanityRates.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &SanityRatesTransferAdminPendingIterator{contract: _SanityRates.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: e TransferAdminPending(pendingAdmin address)
func (_SanityRates *SanityRatesFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *SanityRatesTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _SanityRates.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SanityRatesTransferAdminPending)
				if err := _SanityRates.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
