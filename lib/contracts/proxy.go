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

// ProxyABI is the input ABI used to generate the binding from.
const ProxyABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"destAddress\",\"type\":\"address\"},{\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"name\":\"walletId\",\"type\":\"address\"},{\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHint\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapTokenToEther\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxGasPrice\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kyberNetworkContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserCapInWei\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapTokenToToken\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapEtherToToken\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcQty\",\"type\":\"uint256\"}],\"name\":\"getExpectedRate\",\"outputs\":[{\"name\":\"expectedRate\",\"type\":\"uint256\"},{\"name\":\"slippageRate\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getUserCapInTokenWei\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_kyberNetworkContract\",\"type\":\"address\"}],\"name\":\"setKyberNetworkContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"field\",\"type\":\"bytes32\"}],\"name\":\"info\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"destAddress\",\"type\":\"address\"},{\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"name\":\"walletId\",\"type\":\"address\"}],\"name\":\"trade\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"actualSrcAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"actualDestAmount\",\"type\":\"uint256\"}],\"name\":\"ExecuteTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newNetworkContract\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"oldNetworkContract\",\"type\":\"address\"}],\"name\":\"KyberNetworkSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"}]"

// Proxy is an auto generated Go binding around an Ethereum contract.
type Proxy struct {
	ProxyCaller     // Read-only binding to the contract
	ProxyTransactor // Write-only binding to the contract
	ProxyFilterer   // Log filterer for contract events
}

// ProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProxySession struct {
	Contract     *Proxy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProxyCallerSession struct {
	Contract *ProxyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProxyTransactorSession struct {
	Contract     *ProxyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProxyRaw struct {
	Contract *Proxy // Generic contract binding to access the raw methods on
}

// ProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProxyCallerRaw struct {
	Contract *ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProxyTransactorRaw struct {
	Contract *ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxy creates a new instance of Proxy, bound to a specific deployed contract.
func NewProxy(address common.Address, backend bind.ContractBackend) (*Proxy, error) {
	contract, err := bindProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proxy{ProxyCaller: ProxyCaller{contract: contract}, ProxyTransactor: ProxyTransactor{contract: contract}, ProxyFilterer: ProxyFilterer{contract: contract}}, nil
}

// NewProxyCaller creates a new read-only instance of Proxy, bound to a specific deployed contract.
func NewProxyCaller(address common.Address, caller bind.ContractCaller) (*ProxyCaller, error) {
	contract, err := bindProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyCaller{contract: contract}, nil
}

// NewProxyTransactor creates a new write-only instance of Proxy, bound to a specific deployed contract.
func NewProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyTransactor, error) {
	contract, err := bindProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyTransactor{contract: contract}, nil
}

// NewProxyFilterer creates a new log filterer instance of Proxy, bound to a specific deployed contract.
func NewProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyFilterer, error) {
	contract, err := bindProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyFilterer{contract: contract}, nil
}

// bindProxy binds a generic wrapper to an already deployed contract.
func bindProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Proxy *ProxyCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Proxy *ProxySession) Admin() (common.Address, error) {
	return _Proxy.Contract.Admin(&_Proxy.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Proxy *ProxyCallerSession) Admin() (common.Address, error) {
	return _Proxy.Contract.Admin(&_Proxy.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Proxy *ProxyCaller) Enabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Proxy *ProxySession) Enabled() (bool, error) {
	return _Proxy.Contract.Enabled(&_Proxy.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Proxy *ProxyCallerSession) Enabled() (bool, error) {
	return _Proxy.Contract.Enabled(&_Proxy.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Proxy *ProxyCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Proxy *ProxySession) GetAlerters() ([]common.Address, error) {
	return _Proxy.Contract.GetAlerters(&_Proxy.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Proxy *ProxyCallerSession) GetAlerters() ([]common.Address, error) {
	return _Proxy.Contract.GetAlerters(&_Proxy.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_Proxy *ProxyCaller) GetBalance(opts *bind.CallOpts, token common.Address, user common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "getBalance", token, user)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_Proxy *ProxySession) GetBalance(token common.Address, user common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetBalance(&_Proxy.CallOpts, token, user)
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_Proxy *ProxyCallerSession) GetBalance(token common.Address, user common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetBalance(&_Proxy.CallOpts, token, user)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_Proxy *ProxyCaller) GetExpectedRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	ret := new(struct {
		ExpectedRate *big.Int
		SlippageRate *big.Int
	})
	out := ret
	err := _Proxy.contract.Call(opts, out, "getExpectedRate", src, dest, srcQty)
	return *ret, err
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_Proxy *ProxySession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	return _Proxy.Contract.GetExpectedRate(&_Proxy.CallOpts, src, dest, srcQty)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_Proxy *ProxyCallerSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	return _Proxy.Contract.GetExpectedRate(&_Proxy.CallOpts, src, dest, srcQty)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Proxy *ProxyCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Proxy *ProxySession) GetOperators() ([]common.Address, error) {
	return _Proxy.Contract.GetOperators(&_Proxy.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Proxy *ProxyCallerSession) GetOperators() ([]common.Address, error) {
	return _Proxy.Contract.GetOperators(&_Proxy.CallOpts)
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_Proxy *ProxyCaller) GetUserCapInTokenWei(opts *bind.CallOpts, user common.Address, token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "getUserCapInTokenWei", user, token)
	return *ret0, err
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_Proxy *ProxySession) GetUserCapInTokenWei(user common.Address, token common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetUserCapInTokenWei(&_Proxy.CallOpts, user, token)
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_Proxy *ProxyCallerSession) GetUserCapInTokenWei(user common.Address, token common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetUserCapInTokenWei(&_Proxy.CallOpts, user, token)
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_Proxy *ProxyCaller) GetUserCapInWei(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "getUserCapInWei", user)
	return *ret0, err
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_Proxy *ProxySession) GetUserCapInWei(user common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetUserCapInWei(&_Proxy.CallOpts, user)
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_Proxy *ProxyCallerSession) GetUserCapInWei(user common.Address) (*big.Int, error) {
	return _Proxy.Contract.GetUserCapInWei(&_Proxy.CallOpts, user)
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_Proxy *ProxyCaller) Info(opts *bind.CallOpts, field [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "info", field)
	return *ret0, err
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_Proxy *ProxySession) Info(field [32]byte) (*big.Int, error) {
	return _Proxy.Contract.Info(&_Proxy.CallOpts, field)
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_Proxy *ProxyCallerSession) Info(field [32]byte) (*big.Int, error) {
	return _Proxy.Contract.Info(&_Proxy.CallOpts, field)
}

// KyberNetworkContract is a free data retrieval call binding the contract method 0x4f61ff8b.
//
// Solidity: function kyberNetworkContract() constant returns(address)
func (_Proxy *ProxyCaller) KyberNetworkContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "kyberNetworkContract")
	return *ret0, err
}

// KyberNetworkContract is a free data retrieval call binding the contract method 0x4f61ff8b.
//
// Solidity: function kyberNetworkContract() constant returns(address)
func (_Proxy *ProxySession) KyberNetworkContract() (common.Address, error) {
	return _Proxy.Contract.KyberNetworkContract(&_Proxy.CallOpts)
}

// KyberNetworkContract is a free data retrieval call binding the contract method 0x4f61ff8b.
//
// Solidity: function kyberNetworkContract() constant returns(address)
func (_Proxy *ProxyCallerSession) KyberNetworkContract() (common.Address, error) {
	return _Proxy.Contract.KyberNetworkContract(&_Proxy.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_Proxy *ProxyCaller) MaxGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "maxGasPrice")
	return *ret0, err
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_Proxy *ProxySession) MaxGasPrice() (*big.Int, error) {
	return _Proxy.Contract.MaxGasPrice(&_Proxy.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_Proxy *ProxyCallerSession) MaxGasPrice() (*big.Int, error) {
	return _Proxy.Contract.MaxGasPrice(&_Proxy.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Proxy *ProxyCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Proxy.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Proxy *ProxySession) PendingAdmin() (common.Address, error) {
	return _Proxy.Contract.PendingAdmin(&_Proxy.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Proxy *ProxyCallerSession) PendingAdmin() (common.Address, error) {
	return _Proxy.Contract.PendingAdmin(&_Proxy.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Proxy *ProxyTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Proxy *ProxySession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.AddAlerter(&_Proxy.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Proxy *ProxyTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.AddAlerter(&_Proxy.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Proxy *ProxyTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Proxy *ProxySession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.AddOperator(&_Proxy.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Proxy *ProxyTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.AddOperator(&_Proxy.TransactOpts, newOperator)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Proxy *ProxyTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Proxy *ProxySession) ClaimAdmin() (*types.Transaction, error) {
	return _Proxy.Contract.ClaimAdmin(&_Proxy.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Proxy *ProxyTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _Proxy.Contract.ClaimAdmin(&_Proxy.TransactOpts)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Proxy *ProxyTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Proxy *ProxySession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.RemoveAlerter(&_Proxy.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Proxy *ProxyTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.RemoveAlerter(&_Proxy.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Proxy *ProxyTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Proxy *ProxySession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.RemoveOperator(&_Proxy.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Proxy *ProxyTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.RemoveOperator(&_Proxy.TransactOpts, operator)
}

// SetKyberNetworkContract is a paid mutator transaction binding the contract method 0xabd188a8.
//
// Solidity: function setKyberNetworkContract(address _kyberNetworkContract) returns()
func (_Proxy *ProxyTransactor) SetKyberNetworkContract(opts *bind.TransactOpts, _kyberNetworkContract common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "setKyberNetworkContract", _kyberNetworkContract)
}

// SetKyberNetworkContract is a paid mutator transaction binding the contract method 0xabd188a8.
//
// Solidity: function setKyberNetworkContract(address _kyberNetworkContract) returns()
func (_Proxy *ProxySession) SetKyberNetworkContract(_kyberNetworkContract common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.SetKyberNetworkContract(&_Proxy.TransactOpts, _kyberNetworkContract)
}

// SetKyberNetworkContract is a paid mutator transaction binding the contract method 0xabd188a8.
//
// Solidity: function setKyberNetworkContract(address _kyberNetworkContract) returns()
func (_Proxy *ProxyTransactorSession) SetKyberNetworkContract(_kyberNetworkContract common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.SetKyberNetworkContract(&_Proxy.TransactOpts, _kyberNetworkContract)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactor) SwapEtherToToken(opts *bind.TransactOpts, token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "swapEtherToToken", token, minConversionRate)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxySession) SwapEtherToToken(token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapEtherToToken(&_Proxy.TransactOpts, token, minConversionRate)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactorSession) SwapEtherToToken(token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapEtherToToken(&_Proxy.TransactOpts, token, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactor) SwapTokenToEther(opts *bind.TransactOpts, token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "swapTokenToEther", token, srcAmount, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxySession) SwapTokenToEther(token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapTokenToEther(&_Proxy.TransactOpts, token, srcAmount, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactorSession) SwapTokenToEther(token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapTokenToEther(&_Proxy.TransactOpts, token, srcAmount, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactor) SwapTokenToToken(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "swapTokenToToken", src, srcAmount, dest, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxySession) SwapTokenToToken(src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapTokenToToken(&_Proxy.TransactOpts, src, srcAmount, dest, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_Proxy *ProxyTransactorSession) SwapTokenToToken(src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _Proxy.Contract.SwapTokenToToken(&_Proxy.TransactOpts, src, srcAmount, dest, minConversionRate)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId) returns(uint256)
func (_Proxy *ProxyTransactor) Trade(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "trade", src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId) returns(uint256)
func (_Proxy *ProxySession) Trade(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.Trade(&_Proxy.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId) returns(uint256)
func (_Proxy *ProxyTransactorSession) Trade(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.Trade(&_Proxy.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_Proxy *ProxyTransactor) TradeWithHint(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "tradeWithHint", src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_Proxy *ProxySession) TradeWithHint(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _Proxy.Contract.TradeWithHint(&_Proxy.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_Proxy *ProxyTransactorSession) TradeWithHint(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _Proxy.Contract.TradeWithHint(&_Proxy.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Proxy *ProxyTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Proxy *ProxySession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.TransferAdmin(&_Proxy.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Proxy *ProxyTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.TransferAdmin(&_Proxy.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Proxy *ProxyTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Proxy *ProxySession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.TransferAdminQuickly(&_Proxy.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Proxy *ProxyTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.TransferAdminQuickly(&_Proxy.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Proxy *ProxyTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Proxy *ProxySession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.WithdrawEther(&_Proxy.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Proxy *ProxyTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.WithdrawEther(&_Proxy.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Proxy *ProxyTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Proxy *ProxySession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.WithdrawToken(&_Proxy.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Proxy *ProxyTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Proxy.Contract.WithdrawToken(&_Proxy.TransactOpts, token, amount, sendTo)
}

// ProxyAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the Proxy contract.
type ProxyAdminClaimedIterator struct {
	Event *ProxyAdminClaimed // Event containing the contract specifics and raw log

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
func (it *ProxyAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyAdminClaimed)
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
		it.Event = new(ProxyAdminClaimed)
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
func (it *ProxyAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyAdminClaimed represents a AdminClaimed event raised by the Proxy contract.
type ProxyAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_Proxy *ProxyFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*ProxyAdminClaimedIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &ProxyAdminClaimedIterator{contract: _Proxy.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_Proxy *ProxyFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *ProxyAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyAdminClaimed)
				if err := _Proxy.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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

// ProxyAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the Proxy contract.
type ProxyAlerterAddedIterator struct {
	Event *ProxyAlerterAdded // Event containing the contract specifics and raw log

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
func (it *ProxyAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyAlerterAdded)
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
		it.Event = new(ProxyAlerterAdded)
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
func (it *ProxyAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyAlerterAdded represents a AlerterAdded event raised by the Proxy contract.
type ProxyAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_Proxy *ProxyFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*ProxyAlerterAddedIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &ProxyAlerterAddedIterator{contract: _Proxy.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_Proxy *ProxyFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *ProxyAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyAlerterAdded)
				if err := _Proxy.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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

// ProxyEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the Proxy contract.
type ProxyEtherWithdrawIterator struct {
	Event *ProxyEtherWithdraw // Event containing the contract specifics and raw log

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
func (it *ProxyEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyEtherWithdraw)
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
		it.Event = new(ProxyEtherWithdraw)
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
func (it *ProxyEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyEtherWithdraw represents a EtherWithdraw event raised by the Proxy contract.
type ProxyEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_Proxy *ProxyFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*ProxyEtherWithdrawIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &ProxyEtherWithdrawIterator{contract: _Proxy.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_Proxy *ProxyFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *ProxyEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyEtherWithdraw)
				if err := _Proxy.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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

// ProxyExecuteTradeIterator is returned from FilterExecuteTrade and is used to iterate over the raw logs and unpacked data for ExecuteTrade events raised by the Proxy contract.
type ProxyExecuteTradeIterator struct {
	Event *ProxyExecuteTrade // Event containing the contract specifics and raw log

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
func (it *ProxyExecuteTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyExecuteTrade)
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
		it.Event = new(ProxyExecuteTrade)
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
func (it *ProxyExecuteTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyExecuteTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyExecuteTrade represents a ExecuteTrade event raised by the Proxy contract.
type ProxyExecuteTrade struct {
	Trader           common.Address
	Src              common.Address
	Dest             common.Address
	ActualSrcAmount  *big.Int
	ActualDestAmount *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterExecuteTrade is a free log retrieval operation binding the contract event 0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39.
//
// Solidity: event ExecuteTrade(address indexed trader, address src, address dest, uint256 actualSrcAmount, uint256 actualDestAmount)
func (_Proxy *ProxyFilterer) FilterExecuteTrade(opts *bind.FilterOpts, trader []common.Address) (*ProxyExecuteTradeIterator, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "ExecuteTrade", traderRule)
	if err != nil {
		return nil, err
	}
	return &ProxyExecuteTradeIterator{contract: _Proxy.contract, event: "ExecuteTrade", logs: logs, sub: sub}, nil
}

// WatchExecuteTrade is a free log subscription operation binding the contract event 0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39.
//
// Solidity: event ExecuteTrade(address indexed trader, address src, address dest, uint256 actualSrcAmount, uint256 actualDestAmount)
func (_Proxy *ProxyFilterer) WatchExecuteTrade(opts *bind.WatchOpts, sink chan<- *ProxyExecuteTrade, trader []common.Address) (event.Subscription, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "ExecuteTrade", traderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyExecuteTrade)
				if err := _Proxy.contract.UnpackLog(event, "ExecuteTrade", log); err != nil {
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

// ProxyKyberNetworkSetIterator is returned from FilterKyberNetworkSet and is used to iterate over the raw logs and unpacked data for KyberNetworkSet events raised by the Proxy contract.
type ProxyKyberNetworkSetIterator struct {
	Event *ProxyKyberNetworkSet // Event containing the contract specifics and raw log

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
func (it *ProxyKyberNetworkSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyKyberNetworkSet)
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
		it.Event = new(ProxyKyberNetworkSet)
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
func (it *ProxyKyberNetworkSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyKyberNetworkSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyKyberNetworkSet represents a KyberNetworkSet event raised by the Proxy contract.
type ProxyKyberNetworkSet struct {
	NewNetworkContract common.Address
	OldNetworkContract common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkSet is a free log retrieval operation binding the contract event 0x8936e1f096bf0a8c9df862b3d1d5b82774cad78116200175f00b5b7ba3010b02.
//
// Solidity: event KyberNetworkSet(address newNetworkContract, address oldNetworkContract)
func (_Proxy *ProxyFilterer) FilterKyberNetworkSet(opts *bind.FilterOpts) (*ProxyKyberNetworkSetIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "KyberNetworkSet")
	if err != nil {
		return nil, err
	}
	return &ProxyKyberNetworkSetIterator{contract: _Proxy.contract, event: "KyberNetworkSet", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkSet is a free log subscription operation binding the contract event 0x8936e1f096bf0a8c9df862b3d1d5b82774cad78116200175f00b5b7ba3010b02.
//
// Solidity: event KyberNetworkSet(address newNetworkContract, address oldNetworkContract)
func (_Proxy *ProxyFilterer) WatchKyberNetworkSet(opts *bind.WatchOpts, sink chan<- *ProxyKyberNetworkSet) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "KyberNetworkSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyKyberNetworkSet)
				if err := _Proxy.contract.UnpackLog(event, "KyberNetworkSet", log); err != nil {
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

// ProxyOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the Proxy contract.
type ProxyOperatorAddedIterator struct {
	Event *ProxyOperatorAdded // Event containing the contract specifics and raw log

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
func (it *ProxyOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyOperatorAdded)
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
		it.Event = new(ProxyOperatorAdded)
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
func (it *ProxyOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyOperatorAdded represents a OperatorAdded event raised by the Proxy contract.
type ProxyOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_Proxy *ProxyFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*ProxyOperatorAddedIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &ProxyOperatorAddedIterator{contract: _Proxy.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_Proxy *ProxyFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *ProxyOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyOperatorAdded)
				if err := _Proxy.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ProxyTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the Proxy contract.
type ProxyTokenWithdrawIterator struct {
	Event *ProxyTokenWithdraw // Event containing the contract specifics and raw log

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
func (it *ProxyTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyTokenWithdraw)
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
		it.Event = new(ProxyTokenWithdraw)
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
func (it *ProxyTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyTokenWithdraw represents a TokenWithdraw event raised by the Proxy contract.
type ProxyTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_Proxy *ProxyFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*ProxyTokenWithdrawIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &ProxyTokenWithdrawIterator{contract: _Proxy.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_Proxy *ProxyFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *ProxyTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyTokenWithdraw)
				if err := _Proxy.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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

// ProxyTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the Proxy contract.
type ProxyTransferAdminPendingIterator struct {
	Event *ProxyTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *ProxyTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyTransferAdminPending)
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
		it.Event = new(ProxyTransferAdminPending)
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
func (it *ProxyTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyTransferAdminPending represents a TransferAdminPending event raised by the Proxy contract.
type ProxyTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_Proxy *ProxyFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*ProxyTransferAdminPendingIterator, error) {

	logs, sub, err := _Proxy.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &ProxyTransferAdminPendingIterator{contract: _Proxy.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_Proxy *ProxyFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *ProxyTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _Proxy.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyTransferAdminPending)
				if err := _Proxy.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
