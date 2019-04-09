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

// ConversionRatesABI is the input ABI used to generate the binding from.
const ConversionRatesABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\"}],\"name\":\"setReserveAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"disableTokenTrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validRateDurationInBlocks\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokens\",\"type\":\"address[]\"},{\"name\":\"baseBuy\",\"type\":\"uint256[]\"},{\"name\":\"baseSell\",\"type\":\"uint256[]\"},{\"name\":\"buy\",\"type\":\"bytes14[]\"},{\"name\":\"sell\",\"type\":\"bytes14[]\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"setBaseRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"enableTokenTrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getListedTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numTokensInCurrentCompactData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"command\",\"type\":\"uint256\"},{\"name\":\"param\",\"type\":\"uint256\"}],\"name\":\"getStepFunctionData\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buy\",\"type\":\"bytes14[]\"},{\"name\":\"sell\",\"type\":\"bytes14[]\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"setCompactData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"setValidRateDurationInBlocks\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenBasicData\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getRateUpdateBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"xBuy\",\"type\":\"int256[]\"},{\"name\":\"yBuy\",\"type\":\"int256[]\"},{\"name\":\"xSell\",\"type\":\"int256[]\"},{\"name\":\"ySell\",\"type\":\"int256[]\"}],\"name\":\"setQtyStepFunction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"reserveContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenImbalanceData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"currentBlockNumber\",\"type\":\"uint256\"},{\"name\":\"buy\",\"type\":\"bool\"},{\"name\":\"qty\",\"type\":\"uint256\"}],\"name\":\"getRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"xBuy\",\"type\":\"int256[]\"},{\"name\":\"yBuy\",\"type\":\"int256[]\"},{\"name\":\"xSell\",\"type\":\"int256[]\"},{\"name\":\"ySell\",\"type\":\"int256[]\"}],\"name\":\"setImbalanceStepFunction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"minimalRecordResolution\",\"type\":\"uint256\"},{\"name\":\"maxPerBlockImbalance\",\"type\":\"uint256\"},{\"name\":\"maxTotalImbalance\",\"type\":\"uint256\"}],\"name\":\"setTokenControlInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"buyAmount\",\"type\":\"int256\"},{\"name\":\"rateUpdateBlock\",\"type\":\"uint256\"},{\"name\":\"currentBlock\",\"type\":\"uint256\"}],\"name\":\"recordImbalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"buy\",\"type\":\"bool\"}],\"name\":\"getBasicRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getCompactData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes1\"},{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenControlInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"}]"

// ConversionRates is an auto generated Go binding around an Ethereum contract.
type ConversionRates struct {
	ConversionRatesCaller     // Read-only binding to the contract
	ConversionRatesTransactor // Write-only binding to the contract
	ConversionRatesFilterer   // Log filterer for contract events
}

// ConversionRatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConversionRatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionRatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConversionRatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionRatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConversionRatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConversionRatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConversionRatesSession struct {
	Contract     *ConversionRates  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConversionRatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConversionRatesCallerSession struct {
	Contract *ConversionRatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ConversionRatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConversionRatesTransactorSession struct {
	Contract     *ConversionRatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ConversionRatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConversionRatesRaw struct {
	Contract *ConversionRates // Generic contract binding to access the raw methods on
}

// ConversionRatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConversionRatesCallerRaw struct {
	Contract *ConversionRatesCaller // Generic read-only contract binding to access the raw methods on
}

// ConversionRatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConversionRatesTransactorRaw struct {
	Contract *ConversionRatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConversionRates creates a new instance of ConversionRates, bound to a specific deployed contract.
func NewConversionRates(address common.Address, backend bind.ContractBackend) (*ConversionRates, error) {
	contract, err := bindConversionRates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConversionRates{ConversionRatesCaller: ConversionRatesCaller{contract: contract}, ConversionRatesTransactor: ConversionRatesTransactor{contract: contract}, ConversionRatesFilterer: ConversionRatesFilterer{contract: contract}}, nil
}

// NewConversionRatesCaller creates a new read-only instance of ConversionRates, bound to a specific deployed contract.
func NewConversionRatesCaller(address common.Address, caller bind.ContractCaller) (*ConversionRatesCaller, error) {
	contract, err := bindConversionRates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConversionRatesCaller{contract: contract}, nil
}

// NewConversionRatesTransactor creates a new write-only instance of ConversionRates, bound to a specific deployed contract.
func NewConversionRatesTransactor(address common.Address, transactor bind.ContractTransactor) (*ConversionRatesTransactor, error) {
	contract, err := bindConversionRates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConversionRatesTransactor{contract: contract}, nil
}

// NewConversionRatesFilterer creates a new log filterer instance of ConversionRates, bound to a specific deployed contract.
func NewConversionRatesFilterer(address common.Address, filterer bind.ContractFilterer) (*ConversionRatesFilterer, error) {
	contract, err := bindConversionRates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConversionRatesFilterer{contract: contract}, nil
}

// bindConversionRates binds a generic wrapper to an already deployed contract.
func bindConversionRates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConversionRatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConversionRates *ConversionRatesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConversionRates.Contract.ConversionRatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConversionRates *ConversionRatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConversionRates.Contract.ConversionRatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConversionRates *ConversionRatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConversionRates.Contract.ConversionRatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConversionRates *ConversionRatesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConversionRates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConversionRates *ConversionRatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConversionRates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConversionRates *ConversionRatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConversionRates.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_ConversionRates *ConversionRatesCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_ConversionRates *ConversionRatesSession) Admin() (common.Address, error) {
	return _ConversionRates.Contract.Admin(&_ConversionRates.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_ConversionRates *ConversionRatesCallerSession) Admin() (common.Address, error) {
	return _ConversionRates.Contract.Admin(&_ConversionRates.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_ConversionRates *ConversionRatesCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_ConversionRates *ConversionRatesSession) GetAlerters() ([]common.Address, error) {
	return _ConversionRates.Contract.GetAlerters(&_ConversionRates.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_ConversionRates *ConversionRatesCallerSession) GetAlerters() ([]common.Address, error) {
	return _ConversionRates.Contract.GetAlerters(&_ConversionRates.CallOpts)
}

// GetBasicRate is a free data retrieval call binding the contract method 0xcf8fee11.
//
// Solidity: function getBasicRate(address token, bool buy) constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) GetBasicRate(opts *bind.CallOpts, token common.Address, buy bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getBasicRate", token, buy)
	return *ret0, err
}

// GetBasicRate is a free data retrieval call binding the contract method 0xcf8fee11.
//
// Solidity: function getBasicRate(address token, bool buy) constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) GetBasicRate(token common.Address, buy bool) (*big.Int, error) {
	return _ConversionRates.Contract.GetBasicRate(&_ConversionRates.CallOpts, token, buy)
}

// GetBasicRate is a free data retrieval call binding the contract method 0xcf8fee11.
//
// Solidity: function getBasicRate(address token, bool buy) constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) GetBasicRate(token common.Address, buy bool) (*big.Int, error) {
	return _ConversionRates.Contract.GetBasicRate(&_ConversionRates.CallOpts, token, buy)
}

// GetCompactData is a free data retrieval call binding the contract method 0xe4a2ac62.
//
// Solidity: function getCompactData(address token) constant returns(uint256, uint256, bytes1, bytes1)
func (_ConversionRates *ConversionRatesCaller) GetCompactData(opts *bind.CallOpts, token common.Address) (*big.Int, *big.Int, [1]byte, [1]byte, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new([1]byte)
		ret3 = new([1]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _ConversionRates.contract.Call(opts, out, "getCompactData", token)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetCompactData is a free data retrieval call binding the contract method 0xe4a2ac62.
//
// Solidity: function getCompactData(address token) constant returns(uint256, uint256, bytes1, bytes1)
func (_ConversionRates *ConversionRatesSession) GetCompactData(token common.Address) (*big.Int, *big.Int, [1]byte, [1]byte, error) {
	return _ConversionRates.Contract.GetCompactData(&_ConversionRates.CallOpts, token)
}

// GetCompactData is a free data retrieval call binding the contract method 0xe4a2ac62.
//
// Solidity: function getCompactData(address token) constant returns(uint256, uint256, bytes1, bytes1)
func (_ConversionRates *ConversionRatesCallerSession) GetCompactData(token common.Address) (*big.Int, *big.Int, [1]byte, [1]byte, error) {
	return _ConversionRates.Contract.GetCompactData(&_ConversionRates.CallOpts, token)
}

// GetListedTokens is a free data retrieval call binding the contract method 0x2ba996a5.
//
// Solidity: function getListedTokens() constant returns(address[])
func (_ConversionRates *ConversionRatesCaller) GetListedTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getListedTokens")
	return *ret0, err
}

// GetListedTokens is a free data retrieval call binding the contract method 0x2ba996a5.
//
// Solidity: function getListedTokens() constant returns(address[])
func (_ConversionRates *ConversionRatesSession) GetListedTokens() ([]common.Address, error) {
	return _ConversionRates.Contract.GetListedTokens(&_ConversionRates.CallOpts)
}

// GetListedTokens is a free data retrieval call binding the contract method 0x2ba996a5.
//
// Solidity: function getListedTokens() constant returns(address[])
func (_ConversionRates *ConversionRatesCallerSession) GetListedTokens() ([]common.Address, error) {
	return _ConversionRates.Contract.GetListedTokens(&_ConversionRates.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_ConversionRates *ConversionRatesCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_ConversionRates *ConversionRatesSession) GetOperators() ([]common.Address, error) {
	return _ConversionRates.Contract.GetOperators(&_ConversionRates.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_ConversionRates *ConversionRatesCallerSession) GetOperators() ([]common.Address, error) {
	return _ConversionRates.Contract.GetOperators(&_ConversionRates.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0xb8e9c22e.
//
// Solidity: function getRate(address token, uint256 currentBlockNumber, bool buy, uint256 qty) constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) GetRate(opts *bind.CallOpts, token common.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getRate", token, currentBlockNumber, buy, qty)
	return *ret0, err
}

// GetRate is a free data retrieval call binding the contract method 0xb8e9c22e.
//
// Solidity: function getRate(address token, uint256 currentBlockNumber, bool buy, uint256 qty) constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) GetRate(token common.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.GetRate(&_ConversionRates.CallOpts, token, currentBlockNumber, buy, qty)
}

// GetRate is a free data retrieval call binding the contract method 0xb8e9c22e.
//
// Solidity: function getRate(address token, uint256 currentBlockNumber, bool buy, uint256 qty) constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) GetRate(token common.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.GetRate(&_ConversionRates.CallOpts, token, currentBlockNumber, buy, qty)
}

// GetRateUpdateBlock is a free data retrieval call binding the contract method 0x8036d757.
//
// Solidity: function getRateUpdateBlock(address token) constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) GetRateUpdateBlock(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getRateUpdateBlock", token)
	return *ret0, err
}

// GetRateUpdateBlock is a free data retrieval call binding the contract method 0x8036d757.
//
// Solidity: function getRateUpdateBlock(address token) constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) GetRateUpdateBlock(token common.Address) (*big.Int, error) {
	return _ConversionRates.Contract.GetRateUpdateBlock(&_ConversionRates.CallOpts, token)
}

// GetRateUpdateBlock is a free data retrieval call binding the contract method 0x8036d757.
//
// Solidity: function getRateUpdateBlock(address token) constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) GetRateUpdateBlock(token common.Address) (*big.Int, error) {
	return _ConversionRates.Contract.GetRateUpdateBlock(&_ConversionRates.CallOpts, token)
}

// GetStepFunctionData is a free data retrieval call binding the contract method 0x62674e93.
//
// Solidity: function getStepFunctionData(address token, uint256 command, uint256 param) constant returns(int256)
func (_ConversionRates *ConversionRatesCaller) GetStepFunctionData(opts *bind.CallOpts, token common.Address, command *big.Int, param *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "getStepFunctionData", token, command, param)
	return *ret0, err
}

// GetStepFunctionData is a free data retrieval call binding the contract method 0x62674e93.
//
// Solidity: function getStepFunctionData(address token, uint256 command, uint256 param) constant returns(int256)
func (_ConversionRates *ConversionRatesSession) GetStepFunctionData(token common.Address, command *big.Int, param *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.GetStepFunctionData(&_ConversionRates.CallOpts, token, command, param)
}

// GetStepFunctionData is a free data retrieval call binding the contract method 0x62674e93.
//
// Solidity: function getStepFunctionData(address token, uint256 command, uint256 param) constant returns(int256)
func (_ConversionRates *ConversionRatesCallerSession) GetStepFunctionData(token common.Address, command *big.Int, param *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.GetStepFunctionData(&_ConversionRates.CallOpts, token, command, param)
}

// GetTokenBasicData is a free data retrieval call binding the contract method 0x721bba59.
//
// Solidity: function getTokenBasicData(address token) constant returns(bool, bool)
func (_ConversionRates *ConversionRatesCaller) GetTokenBasicData(opts *bind.CallOpts, token common.Address) (bool, bool, error) {
	var (
		ret0 = new(bool)
		ret1 = new(bool)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _ConversionRates.contract.Call(opts, out, "getTokenBasicData", token)
	return *ret0, *ret1, err
}

// GetTokenBasicData is a free data retrieval call binding the contract method 0x721bba59.
//
// Solidity: function getTokenBasicData(address token) constant returns(bool, bool)
func (_ConversionRates *ConversionRatesSession) GetTokenBasicData(token common.Address) (bool, bool, error) {
	return _ConversionRates.Contract.GetTokenBasicData(&_ConversionRates.CallOpts, token)
}

// GetTokenBasicData is a free data retrieval call binding the contract method 0x721bba59.
//
// Solidity: function getTokenBasicData(address token) constant returns(bool, bool)
func (_ConversionRates *ConversionRatesCallerSession) GetTokenBasicData(token common.Address) (bool, bool, error) {
	return _ConversionRates.Contract.GetTokenBasicData(&_ConversionRates.CallOpts, token)
}

// GetTokenControlInfo is a free data retrieval call binding the contract method 0xe7d4fd91.
//
// Solidity: function getTokenControlInfo(address token) constant returns(uint256, uint256, uint256)
func (_ConversionRates *ConversionRatesCaller) GetTokenControlInfo(opts *bind.CallOpts, token common.Address) (*big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _ConversionRates.contract.Call(opts, out, "getTokenControlInfo", token)
	return *ret0, *ret1, *ret2, err
}

// GetTokenControlInfo is a free data retrieval call binding the contract method 0xe7d4fd91.
//
// Solidity: function getTokenControlInfo(address token) constant returns(uint256, uint256, uint256)
func (_ConversionRates *ConversionRatesSession) GetTokenControlInfo(token common.Address) (*big.Int, *big.Int, *big.Int, error) {
	return _ConversionRates.Contract.GetTokenControlInfo(&_ConversionRates.CallOpts, token)
}

// GetTokenControlInfo is a free data retrieval call binding the contract method 0xe7d4fd91.
//
// Solidity: function getTokenControlInfo(address token) constant returns(uint256, uint256, uint256)
func (_ConversionRates *ConversionRatesCallerSession) GetTokenControlInfo(token common.Address) (*big.Int, *big.Int, *big.Int, error) {
	return _ConversionRates.Contract.GetTokenControlInfo(&_ConversionRates.CallOpts, token)
}

// NumTokensInCurrentCompactData is a free data retrieval call binding the contract method 0x5085c9f1.
//
// Solidity: function numTokensInCurrentCompactData() constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) NumTokensInCurrentCompactData(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "numTokensInCurrentCompactData")
	return *ret0, err
}

// NumTokensInCurrentCompactData is a free data retrieval call binding the contract method 0x5085c9f1.
//
// Solidity: function numTokensInCurrentCompactData() constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) NumTokensInCurrentCompactData() (*big.Int, error) {
	return _ConversionRates.Contract.NumTokensInCurrentCompactData(&_ConversionRates.CallOpts)
}

// NumTokensInCurrentCompactData is a free data retrieval call binding the contract method 0x5085c9f1.
//
// Solidity: function numTokensInCurrentCompactData() constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) NumTokensInCurrentCompactData() (*big.Int, error) {
	return _ConversionRates.Contract.NumTokensInCurrentCompactData(&_ConversionRates.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_ConversionRates *ConversionRatesCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_ConversionRates *ConversionRatesSession) PendingAdmin() (common.Address, error) {
	return _ConversionRates.Contract.PendingAdmin(&_ConversionRates.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_ConversionRates *ConversionRatesCallerSession) PendingAdmin() (common.Address, error) {
	return _ConversionRates.Contract.PendingAdmin(&_ConversionRates.CallOpts)
}

// ReserveContract is a free data retrieval call binding the contract method 0xa7f43acd.
//
// Solidity: function reserveContract() constant returns(address)
func (_ConversionRates *ConversionRatesCaller) ReserveContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "reserveContract")
	return *ret0, err
}

// ReserveContract is a free data retrieval call binding the contract method 0xa7f43acd.
//
// Solidity: function reserveContract() constant returns(address)
func (_ConversionRates *ConversionRatesSession) ReserveContract() (common.Address, error) {
	return _ConversionRates.Contract.ReserveContract(&_ConversionRates.CallOpts)
}

// ReserveContract is a free data retrieval call binding the contract method 0xa7f43acd.
//
// Solidity: function reserveContract() constant returns(address)
func (_ConversionRates *ConversionRatesCallerSession) ReserveContract() (common.Address, error) {
	return _ConversionRates.Contract.ReserveContract(&_ConversionRates.CallOpts)
}

// TokenImbalanceData is a free data retrieval call binding the contract method 0xa80c609e.
//
// Solidity: function tokenImbalanceData(address , uint256 ) constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) TokenImbalanceData(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "tokenImbalanceData", arg0, arg1)
	return *ret0, err
}

// TokenImbalanceData is a free data retrieval call binding the contract method 0xa80c609e.
//
// Solidity: function tokenImbalanceData(address , uint256 ) constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) TokenImbalanceData(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.TokenImbalanceData(&_ConversionRates.CallOpts, arg0, arg1)
}

// TokenImbalanceData is a free data retrieval call binding the contract method 0xa80c609e.
//
// Solidity: function tokenImbalanceData(address , uint256 ) constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) TokenImbalanceData(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _ConversionRates.Contract.TokenImbalanceData(&_ConversionRates.CallOpts, arg0, arg1)
}

// ValidRateDurationInBlocks is a free data retrieval call binding the contract method 0x16265694.
//
// Solidity: function validRateDurationInBlocks() constant returns(uint256)
func (_ConversionRates *ConversionRatesCaller) ValidRateDurationInBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConversionRates.contract.Call(opts, out, "validRateDurationInBlocks")
	return *ret0, err
}

// ValidRateDurationInBlocks is a free data retrieval call binding the contract method 0x16265694.
//
// Solidity: function validRateDurationInBlocks() constant returns(uint256)
func (_ConversionRates *ConversionRatesSession) ValidRateDurationInBlocks() (*big.Int, error) {
	return _ConversionRates.Contract.ValidRateDurationInBlocks(&_ConversionRates.CallOpts)
}

// ValidRateDurationInBlocks is a free data retrieval call binding the contract method 0x16265694.
//
// Solidity: function validRateDurationInBlocks() constant returns(uint256)
func (_ConversionRates *ConversionRatesCallerSession) ValidRateDurationInBlocks() (*big.Int, error) {
	return _ConversionRates.Contract.ValidRateDurationInBlocks(&_ConversionRates.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_ConversionRates *ConversionRatesTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_ConversionRates *ConversionRatesSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddAlerter(&_ConversionRates.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_ConversionRates *ConversionRatesTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddAlerter(&_ConversionRates.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_ConversionRates *ConversionRatesTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_ConversionRates *ConversionRatesSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddOperator(&_ConversionRates.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_ConversionRates *ConversionRatesTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddOperator(&_ConversionRates.TransactOpts, newOperator)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address token) returns()
func (_ConversionRates *ConversionRatesTransactor) AddToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "addToken", token)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address token) returns()
func (_ConversionRates *ConversionRatesSession) AddToken(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddToken(&_ConversionRates.TransactOpts, token)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address token) returns()
func (_ConversionRates *ConversionRatesTransactorSession) AddToken(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.AddToken(&_ConversionRates.TransactOpts, token)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_ConversionRates *ConversionRatesTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_ConversionRates *ConversionRatesSession) ClaimAdmin() (*types.Transaction, error) {
	return _ConversionRates.Contract.ClaimAdmin(&_ConversionRates.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_ConversionRates *ConversionRatesTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _ConversionRates.Contract.ClaimAdmin(&_ConversionRates.TransactOpts)
}

// DisableTokenTrade is a paid mutator transaction binding the contract method 0x158859f7.
//
// Solidity: function disableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesTransactor) DisableTokenTrade(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "disableTokenTrade", token)
}

// DisableTokenTrade is a paid mutator transaction binding the contract method 0x158859f7.
//
// Solidity: function disableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesSession) DisableTokenTrade(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.DisableTokenTrade(&_ConversionRates.TransactOpts, token)
}

// DisableTokenTrade is a paid mutator transaction binding the contract method 0x158859f7.
//
// Solidity: function disableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesTransactorSession) DisableTokenTrade(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.DisableTokenTrade(&_ConversionRates.TransactOpts, token)
}

// EnableTokenTrade is a paid mutator transaction binding the contract method 0x1d6a8bda.
//
// Solidity: function enableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesTransactor) EnableTokenTrade(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "enableTokenTrade", token)
}

// EnableTokenTrade is a paid mutator transaction binding the contract method 0x1d6a8bda.
//
// Solidity: function enableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesSession) EnableTokenTrade(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.EnableTokenTrade(&_ConversionRates.TransactOpts, token)
}

// EnableTokenTrade is a paid mutator transaction binding the contract method 0x1d6a8bda.
//
// Solidity: function enableTokenTrade(address token) returns()
func (_ConversionRates *ConversionRatesTransactorSession) EnableTokenTrade(token common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.EnableTokenTrade(&_ConversionRates.TransactOpts, token)
}

// RecordImbalance is a paid mutator transaction binding the contract method 0xc6fd2103.
//
// Solidity: function recordImbalance(address token, int256 buyAmount, uint256 rateUpdateBlock, uint256 currentBlock) returns()
func (_ConversionRates *ConversionRatesTransactor) RecordImbalance(opts *bind.TransactOpts, token common.Address, buyAmount *big.Int, rateUpdateBlock *big.Int, currentBlock *big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "recordImbalance", token, buyAmount, rateUpdateBlock, currentBlock)
}

// RecordImbalance is a paid mutator transaction binding the contract method 0xc6fd2103.
//
// Solidity: function recordImbalance(address token, int256 buyAmount, uint256 rateUpdateBlock, uint256 currentBlock) returns()
func (_ConversionRates *ConversionRatesSession) RecordImbalance(token common.Address, buyAmount *big.Int, rateUpdateBlock *big.Int, currentBlock *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.RecordImbalance(&_ConversionRates.TransactOpts, token, buyAmount, rateUpdateBlock, currentBlock)
}

// RecordImbalance is a paid mutator transaction binding the contract method 0xc6fd2103.
//
// Solidity: function recordImbalance(address token, int256 buyAmount, uint256 rateUpdateBlock, uint256 currentBlock) returns()
func (_ConversionRates *ConversionRatesTransactorSession) RecordImbalance(token common.Address, buyAmount *big.Int, rateUpdateBlock *big.Int, currentBlock *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.RecordImbalance(&_ConversionRates.TransactOpts, token, buyAmount, rateUpdateBlock, currentBlock)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_ConversionRates *ConversionRatesTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_ConversionRates *ConversionRatesSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.RemoveAlerter(&_ConversionRates.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_ConversionRates *ConversionRatesTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.RemoveAlerter(&_ConversionRates.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_ConversionRates *ConversionRatesTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_ConversionRates *ConversionRatesSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.RemoveOperator(&_ConversionRates.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_ConversionRates *ConversionRatesTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.RemoveOperator(&_ConversionRates.TransactOpts, operator)
}

// SetBaseRate is a paid mutator transaction binding the contract method 0x1a4813d7.
//
// Solidity: function setBaseRate(address[] tokens, uint256[] baseBuy, uint256[] baseSell, bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesTransactor) SetBaseRate(opts *bind.TransactOpts, tokens []common.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setBaseRate", tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

// SetBaseRate is a paid mutator transaction binding the contract method 0x1a4813d7.
//
// Solidity: function setBaseRate(address[] tokens, uint256[] baseBuy, uint256[] baseSell, bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesSession) SetBaseRate(tokens []common.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetBaseRate(&_ConversionRates.TransactOpts, tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

// SetBaseRate is a paid mutator transaction binding the contract method 0x1a4813d7.
//
// Solidity: function setBaseRate(address[] tokens, uint256[] baseBuy, uint256[] baseSell, bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetBaseRate(tokens []common.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetBaseRate(&_ConversionRates.TransactOpts, tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

// SetCompactData is a paid mutator transaction binding the contract method 0x64887334.
//
// Solidity: function setCompactData(bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesTransactor) SetCompactData(opts *bind.TransactOpts, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setCompactData", buy, sell, blockNumber, indices)
}

// SetCompactData is a paid mutator transaction binding the contract method 0x64887334.
//
// Solidity: function setCompactData(bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesSession) SetCompactData(buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetCompactData(&_ConversionRates.TransactOpts, buy, sell, blockNumber, indices)
}

// SetCompactData is a paid mutator transaction binding the contract method 0x64887334.
//
// Solidity: function setCompactData(bytes14[] buy, bytes14[] sell, uint256 blockNumber, uint256[] indices) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetCompactData(buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetCompactData(&_ConversionRates.TransactOpts, buy, sell, blockNumber, indices)
}

// SetImbalanceStepFunction is a paid mutator transaction binding the contract method 0xbc9cbcc8.
//
// Solidity: function setImbalanceStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesTransactor) SetImbalanceStepFunction(opts *bind.TransactOpts, token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setImbalanceStepFunction", token, xBuy, yBuy, xSell, ySell)
}

// SetImbalanceStepFunction is a paid mutator transaction binding the contract method 0xbc9cbcc8.
//
// Solidity: function setImbalanceStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesSession) SetImbalanceStepFunction(token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetImbalanceStepFunction(&_ConversionRates.TransactOpts, token, xBuy, yBuy, xSell, ySell)
}

// SetImbalanceStepFunction is a paid mutator transaction binding the contract method 0xbc9cbcc8.
//
// Solidity: function setImbalanceStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetImbalanceStepFunction(token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetImbalanceStepFunction(&_ConversionRates.TransactOpts, token, xBuy, yBuy, xSell, ySell)
}

// SetQtyStepFunction is a paid mutator transaction binding the contract method 0x80d8b380.
//
// Solidity: function setQtyStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesTransactor) SetQtyStepFunction(opts *bind.TransactOpts, token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setQtyStepFunction", token, xBuy, yBuy, xSell, ySell)
}

// SetQtyStepFunction is a paid mutator transaction binding the contract method 0x80d8b380.
//
// Solidity: function setQtyStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesSession) SetQtyStepFunction(token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetQtyStepFunction(&_ConversionRates.TransactOpts, token, xBuy, yBuy, xSell, ySell)
}

// SetQtyStepFunction is a paid mutator transaction binding the contract method 0x80d8b380.
//
// Solidity: function setQtyStepFunction(address token, int256[] xBuy, int256[] yBuy, int256[] xSell, int256[] ySell) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetQtyStepFunction(token common.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetQtyStepFunction(&_ConversionRates.TransactOpts, token, xBuy, yBuy, xSell, ySell)
}

// SetReserveAddress is a paid mutator transaction binding the contract method 0x14673d31.
//
// Solidity: function setReserveAddress(address reserve) returns()
func (_ConversionRates *ConversionRatesTransactor) SetReserveAddress(opts *bind.TransactOpts, reserve common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setReserveAddress", reserve)
}

// SetReserveAddress is a paid mutator transaction binding the contract method 0x14673d31.
//
// Solidity: function setReserveAddress(address reserve) returns()
func (_ConversionRates *ConversionRatesSession) SetReserveAddress(reserve common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetReserveAddress(&_ConversionRates.TransactOpts, reserve)
}

// SetReserveAddress is a paid mutator transaction binding the contract method 0x14673d31.
//
// Solidity: function setReserveAddress(address reserve) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetReserveAddress(reserve common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetReserveAddress(&_ConversionRates.TransactOpts, reserve)
}

// SetTokenControlInfo is a paid mutator transaction binding the contract method 0xbfee3569.
//
// Solidity: function setTokenControlInfo(address token, uint256 minimalRecordResolution, uint256 maxPerBlockImbalance, uint256 maxTotalImbalance) returns()
func (_ConversionRates *ConversionRatesTransactor) SetTokenControlInfo(opts *bind.TransactOpts, token common.Address, minimalRecordResolution *big.Int, maxPerBlockImbalance *big.Int, maxTotalImbalance *big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setTokenControlInfo", token, minimalRecordResolution, maxPerBlockImbalance, maxTotalImbalance)
}

// SetTokenControlInfo is a paid mutator transaction binding the contract method 0xbfee3569.
//
// Solidity: function setTokenControlInfo(address token, uint256 minimalRecordResolution, uint256 maxPerBlockImbalance, uint256 maxTotalImbalance) returns()
func (_ConversionRates *ConversionRatesSession) SetTokenControlInfo(token common.Address, minimalRecordResolution *big.Int, maxPerBlockImbalance *big.Int, maxTotalImbalance *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetTokenControlInfo(&_ConversionRates.TransactOpts, token, minimalRecordResolution, maxPerBlockImbalance, maxTotalImbalance)
}

// SetTokenControlInfo is a paid mutator transaction binding the contract method 0xbfee3569.
//
// Solidity: function setTokenControlInfo(address token, uint256 minimalRecordResolution, uint256 maxPerBlockImbalance, uint256 maxTotalImbalance) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetTokenControlInfo(token common.Address, minimalRecordResolution *big.Int, maxPerBlockImbalance *big.Int, maxTotalImbalance *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetTokenControlInfo(&_ConversionRates.TransactOpts, token, minimalRecordResolution, maxPerBlockImbalance, maxTotalImbalance)
}

// SetValidRateDurationInBlocks is a paid mutator transaction binding the contract method 0x6c6295b8.
//
// Solidity: function setValidRateDurationInBlocks(uint256 duration) returns()
func (_ConversionRates *ConversionRatesTransactor) SetValidRateDurationInBlocks(opts *bind.TransactOpts, duration *big.Int) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "setValidRateDurationInBlocks", duration)
}

// SetValidRateDurationInBlocks is a paid mutator transaction binding the contract method 0x6c6295b8.
//
// Solidity: function setValidRateDurationInBlocks(uint256 duration) returns()
func (_ConversionRates *ConversionRatesSession) SetValidRateDurationInBlocks(duration *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetValidRateDurationInBlocks(&_ConversionRates.TransactOpts, duration)
}

// SetValidRateDurationInBlocks is a paid mutator transaction binding the contract method 0x6c6295b8.
//
// Solidity: function setValidRateDurationInBlocks(uint256 duration) returns()
func (_ConversionRates *ConversionRatesTransactorSession) SetValidRateDurationInBlocks(duration *big.Int) (*types.Transaction, error) {
	return _ConversionRates.Contract.SetValidRateDurationInBlocks(&_ConversionRates.TransactOpts, duration)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_ConversionRates *ConversionRatesTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_ConversionRates *ConversionRatesSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.TransferAdmin(&_ConversionRates.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_ConversionRates *ConversionRatesTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.TransferAdmin(&_ConversionRates.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_ConversionRates *ConversionRatesTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_ConversionRates *ConversionRatesSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.TransferAdminQuickly(&_ConversionRates.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_ConversionRates *ConversionRatesTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.TransferAdminQuickly(&_ConversionRates.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.WithdrawEther(&_ConversionRates.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.WithdrawEther(&_ConversionRates.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.WithdrawToken(&_ConversionRates.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_ConversionRates *ConversionRatesTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _ConversionRates.Contract.WithdrawToken(&_ConversionRates.TransactOpts, token, amount, sendTo)
}

// ConversionRatesAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the ConversionRates contract.
type ConversionRatesAdminClaimedIterator struct {
	Event *ConversionRatesAdminClaimed // Event containing the contract specifics and raw log

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
func (it *ConversionRatesAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesAdminClaimed)
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
		it.Event = new(ConversionRatesAdminClaimed)
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
func (it *ConversionRatesAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesAdminClaimed represents a AdminClaimed event raised by the ConversionRates contract.
type ConversionRatesAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_ConversionRates *ConversionRatesFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*ConversionRatesAdminClaimedIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesAdminClaimedIterator{contract: _ConversionRates.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_ConversionRates *ConversionRatesFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *ConversionRatesAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesAdminClaimed)
				if err := _ConversionRates.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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

// ConversionRatesAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the ConversionRates contract.
type ConversionRatesAlerterAddedIterator struct {
	Event *ConversionRatesAlerterAdded // Event containing the contract specifics and raw log

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
func (it *ConversionRatesAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesAlerterAdded)
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
		it.Event = new(ConversionRatesAlerterAdded)
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
func (it *ConversionRatesAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesAlerterAdded represents a AlerterAdded event raised by the ConversionRates contract.
type ConversionRatesAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_ConversionRates *ConversionRatesFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*ConversionRatesAlerterAddedIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesAlerterAddedIterator{contract: _ConversionRates.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_ConversionRates *ConversionRatesFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *ConversionRatesAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesAlerterAdded)
				if err := _ConversionRates.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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

// ConversionRatesEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the ConversionRates contract.
type ConversionRatesEtherWithdrawIterator struct {
	Event *ConversionRatesEtherWithdraw // Event containing the contract specifics and raw log

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
func (it *ConversionRatesEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesEtherWithdraw)
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
		it.Event = new(ConversionRatesEtherWithdraw)
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
func (it *ConversionRatesEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesEtherWithdraw represents a EtherWithdraw event raised by the ConversionRates contract.
type ConversionRatesEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_ConversionRates *ConversionRatesFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*ConversionRatesEtherWithdrawIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesEtherWithdrawIterator{contract: _ConversionRates.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_ConversionRates *ConversionRatesFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *ConversionRatesEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesEtherWithdraw)
				if err := _ConversionRates.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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

// ConversionRatesOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the ConversionRates contract.
type ConversionRatesOperatorAddedIterator struct {
	Event *ConversionRatesOperatorAdded // Event containing the contract specifics and raw log

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
func (it *ConversionRatesOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesOperatorAdded)
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
		it.Event = new(ConversionRatesOperatorAdded)
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
func (it *ConversionRatesOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesOperatorAdded represents a OperatorAdded event raised by the ConversionRates contract.
type ConversionRatesOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_ConversionRates *ConversionRatesFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*ConversionRatesOperatorAddedIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesOperatorAddedIterator{contract: _ConversionRates.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_ConversionRates *ConversionRatesFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *ConversionRatesOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesOperatorAdded)
				if err := _ConversionRates.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ConversionRatesTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the ConversionRates contract.
type ConversionRatesTokenWithdrawIterator struct {
	Event *ConversionRatesTokenWithdraw // Event containing the contract specifics and raw log

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
func (it *ConversionRatesTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesTokenWithdraw)
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
		it.Event = new(ConversionRatesTokenWithdraw)
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
func (it *ConversionRatesTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesTokenWithdraw represents a TokenWithdraw event raised by the ConversionRates contract.
type ConversionRatesTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_ConversionRates *ConversionRatesFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*ConversionRatesTokenWithdrawIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesTokenWithdrawIterator{contract: _ConversionRates.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_ConversionRates *ConversionRatesFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *ConversionRatesTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesTokenWithdraw)
				if err := _ConversionRates.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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

// ConversionRatesTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the ConversionRates contract.
type ConversionRatesTransferAdminPendingIterator struct {
	Event *ConversionRatesTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *ConversionRatesTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConversionRatesTransferAdminPending)
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
		it.Event = new(ConversionRatesTransferAdminPending)
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
func (it *ConversionRatesTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConversionRatesTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConversionRatesTransferAdminPending represents a TransferAdminPending event raised by the ConversionRates contract.
type ConversionRatesTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_ConversionRates *ConversionRatesFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*ConversionRatesTransferAdminPendingIterator, error) {

	logs, sub, err := _ConversionRates.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &ConversionRatesTransferAdminPendingIterator{contract: _ConversionRates.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_ConversionRates *ConversionRatesFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *ConversionRatesTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _ConversionRates.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConversionRatesTransferAdminPending)
				if err := _ConversionRates.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
