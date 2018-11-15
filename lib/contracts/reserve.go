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

// ReserveABI is the input ABI used to generate the binding from.
const ReserveABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"enableTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sanityRatesContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"approve\",\"type\":\"bool\"}],\"name\":\"approveWithdrawAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"disableTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"srcToken\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"destToken\",\"type\":\"address\"},{\"name\":\"destAddress\",\"type\":\"address\"},{\"name\":\"conversionRate\",\"type\":\"uint256\"},{\"name\":\"validate\",\"type\":\"bool\"}],\"name\":\"trade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcQty\",\"type\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getConversionRate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"dstQty\",\"type\":\"uint256\"},{\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"getSrcQty\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_kyberNetwork\",\"type\":\"address\"},{\"name\":\"_conversionRates\",\"type\":\"address\"},{\"name\":\"_sanityRates\",\"type\":\"address\"}],\"name\":\"setContracts\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kyberNetwork\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"conversionRatesContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tradeEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"approvedWithdrawAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcQty\",\"type\":\"uint256\"},{\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"getDestQty\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_kyberNetwork\",\"type\":\"address\"},{\"name\":\"_ratesContract\",\"type\":\"address\"},{\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"destToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"destAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"destAddress\",\"type\":\"address\"}],\"name\":\"TradeExecute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"TradeEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approve\",\"type\":\"bool\"}],\"name\":\"WithdrawAddressApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"WithdrawFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"network\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"rate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sanity\",\"type\":\"address\"}],\"name\":\"SetContractAddresses\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"}]"

// Reserve is an auto generated Go binding around an Ethereum contract.
type Reserve struct {
	ReserveCaller     // Read-only binding to the contract
	ReserveTransactor // Write-only binding to the contract
	ReserveFilterer   // Log filterer for contract events
}

// ReserveCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReserveCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReserveTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReserveFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReserveSession struct {
	Contract     *Reserve          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReserveCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReserveCallerSession struct {
	Contract *ReserveCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ReserveTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReserveTransactorSession struct {
	Contract     *ReserveTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ReserveRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReserveRaw struct {
	Contract *Reserve // Generic contract binding to access the raw methods on
}

// ReserveCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReserveCallerRaw struct {
	Contract *ReserveCaller // Generic read-only contract binding to access the raw methods on
}

// ReserveTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReserveTransactorRaw struct {
	Contract *ReserveTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReserve creates a new instance of Reserve, bound to a specific deployed contract.
func NewReserve(address common.Address, backend bind.ContractBackend) (*Reserve, error) {
	contract, err := bindReserve(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reserve{ReserveCaller: ReserveCaller{contract: contract}, ReserveTransactor: ReserveTransactor{contract: contract}, ReserveFilterer: ReserveFilterer{contract: contract}}, nil
}

// NewReserveCaller creates a new read-only instance of Reserve, bound to a specific deployed contract.
func NewReserveCaller(address common.Address, caller bind.ContractCaller) (*ReserveCaller, error) {
	contract, err := bindReserve(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReserveCaller{contract: contract}, nil
}

// NewReserveTransactor creates a new write-only instance of Reserve, bound to a specific deployed contract.
func NewReserveTransactor(address common.Address, transactor bind.ContractTransactor) (*ReserveTransactor, error) {
	contract, err := bindReserve(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReserveTransactor{contract: contract}, nil
}

// NewReserveFilterer creates a new log filterer instance of Reserve, bound to a specific deployed contract.
func NewReserveFilterer(address common.Address, filterer bind.ContractFilterer) (*ReserveFilterer, error) {
	contract, err := bindReserve(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReserveFilterer{contract: contract}, nil
}

// bindReserve binds a generic wrapper to an already deployed contract.
func bindReserve(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReserveABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reserve *ReserveRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Reserve.Contract.ReserveCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reserve *ReserveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.Contract.ReserveTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reserve *ReserveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reserve.Contract.ReserveTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reserve *ReserveCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Reserve.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reserve *ReserveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reserve *ReserveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reserve.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Reserve *ReserveCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Reserve *ReserveSession) Admin() (common.Address, error) {
	return _Reserve.Contract.Admin(&_Reserve.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_Reserve *ReserveCallerSession) Admin() (common.Address, error) {
	return _Reserve.Contract.Admin(&_Reserve.CallOpts)
}

// ApprovedWithdrawAddresses is a free data retrieval call binding the contract method 0xd7b7024d.
//
// Solidity: function approvedWithdrawAddresses( bytes32) constant returns(bool)
func (_Reserve *ReserveCaller) ApprovedWithdrawAddresses(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "approvedWithdrawAddresses", arg0)
	return *ret0, err
}

// ApprovedWithdrawAddresses is a free data retrieval call binding the contract method 0xd7b7024d.
//
// Solidity: function approvedWithdrawAddresses( bytes32) constant returns(bool)
func (_Reserve *ReserveSession) ApprovedWithdrawAddresses(arg0 [32]byte) (bool, error) {
	return _Reserve.Contract.ApprovedWithdrawAddresses(&_Reserve.CallOpts, arg0)
}

// ApprovedWithdrawAddresses is a free data retrieval call binding the contract method 0xd7b7024d.
//
// Solidity: function approvedWithdrawAddresses( bytes32) constant returns(bool)
func (_Reserve *ReserveCallerSession) ApprovedWithdrawAddresses(arg0 [32]byte) (bool, error) {
	return _Reserve.Contract.ApprovedWithdrawAddresses(&_Reserve.CallOpts, arg0)
}

// ConversionRatesContract is a free data retrieval call binding the contract method 0xd5847d33.
//
// Solidity: function conversionRatesContract() constant returns(address)
func (_Reserve *ReserveCaller) ConversionRatesContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "conversionRatesContract")
	return *ret0, err
}

// ConversionRatesContract is a free data retrieval call binding the contract method 0xd5847d33.
//
// Solidity: function conversionRatesContract() constant returns(address)
func (_Reserve *ReserveSession) ConversionRatesContract() (common.Address, error) {
	return _Reserve.Contract.ConversionRatesContract(&_Reserve.CallOpts)
}

// ConversionRatesContract is a free data retrieval call binding the contract method 0xd5847d33.
//
// Solidity: function conversionRatesContract() constant returns(address)
func (_Reserve *ReserveCallerSession) ConversionRatesContract() (common.Address, error) {
	return _Reserve.Contract.ConversionRatesContract(&_Reserve.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Reserve *ReserveCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Reserve *ReserveSession) GetAlerters() ([]common.Address, error) {
	return _Reserve.Contract.GetAlerters(&_Reserve.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_Reserve *ReserveCallerSession) GetAlerters() ([]common.Address, error) {
	return _Reserve.Contract.GetAlerters(&_Reserve.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(token address) constant returns(uint256)
func (_Reserve *ReserveCaller) GetBalance(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getBalance", token)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(token address) constant returns(uint256)
func (_Reserve *ReserveSession) GetBalance(token common.Address) (*big.Int, error) {
	return _Reserve.Contract.GetBalance(&_Reserve.CallOpts, token)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(token address) constant returns(uint256)
func (_Reserve *ReserveCallerSession) GetBalance(token common.Address) (*big.Int, error) {
	return _Reserve.Contract.GetBalance(&_Reserve.CallOpts, token)
}

// GetConversionRate is a free data retrieval call binding the contract method 0x7cd44272.
//
// Solidity: function getConversionRate(src address, dest address, srcQty uint256, blockNumber uint256) constant returns(uint256)
func (_Reserve *ReserveCaller) GetConversionRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int, blockNumber *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getConversionRate", src, dest, srcQty, blockNumber)
	return *ret0, err
}

// GetConversionRate is a free data retrieval call binding the contract method 0x7cd44272.
//
// Solidity: function getConversionRate(src address, dest address, srcQty uint256, blockNumber uint256) constant returns(uint256)
func (_Reserve *ReserveSession) GetConversionRate(src common.Address, dest common.Address, srcQty *big.Int, blockNumber *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetConversionRate(&_Reserve.CallOpts, src, dest, srcQty, blockNumber)
}

// GetConversionRate is a free data retrieval call binding the contract method 0x7cd44272.
//
// Solidity: function getConversionRate(src address, dest address, srcQty uint256, blockNumber uint256) constant returns(uint256)
func (_Reserve *ReserveCallerSession) GetConversionRate(src common.Address, dest common.Address, srcQty *big.Int, blockNumber *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetConversionRate(&_Reserve.CallOpts, src, dest, srcQty, blockNumber)
}

// GetDestQty is a free data retrieval call binding the contract method 0xfa64dffa.
//
// Solidity: function getDestQty(src address, dest address, srcQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveCaller) GetDestQty(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int, rate *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getDestQty", src, dest, srcQty, rate)
	return *ret0, err
}

// GetDestQty is a free data retrieval call binding the contract method 0xfa64dffa.
//
// Solidity: function getDestQty(src address, dest address, srcQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveSession) GetDestQty(src common.Address, dest common.Address, srcQty *big.Int, rate *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetDestQty(&_Reserve.CallOpts, src, dest, srcQty, rate)
}

// GetDestQty is a free data retrieval call binding the contract method 0xfa64dffa.
//
// Solidity: function getDestQty(src address, dest address, srcQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveCallerSession) GetDestQty(src common.Address, dest common.Address, srcQty *big.Int, rate *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetDestQty(&_Reserve.CallOpts, src, dest, srcQty, rate)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Reserve *ReserveCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Reserve *ReserveSession) GetOperators() ([]common.Address, error) {
	return _Reserve.Contract.GetOperators(&_Reserve.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_Reserve *ReserveCallerSession) GetOperators() ([]common.Address, error) {
	return _Reserve.Contract.GetOperators(&_Reserve.CallOpts)
}

// GetSrcQty is a free data retrieval call binding the contract method 0xa7fca953.
//
// Solidity: function getSrcQty(src address, dest address, dstQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveCaller) GetSrcQty(opts *bind.CallOpts, src common.Address, dest common.Address, dstQty *big.Int, rate *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getSrcQty", src, dest, dstQty, rate)
	return *ret0, err
}

// GetSrcQty is a free data retrieval call binding the contract method 0xa7fca953.
//
// Solidity: function getSrcQty(src address, dest address, dstQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveSession) GetSrcQty(src common.Address, dest common.Address, dstQty *big.Int, rate *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetSrcQty(&_Reserve.CallOpts, src, dest, dstQty, rate)
}

// GetSrcQty is a free data retrieval call binding the contract method 0xa7fca953.
//
// Solidity: function getSrcQty(src address, dest address, dstQty uint256, rate uint256) constant returns(uint256)
func (_Reserve *ReserveCallerSession) GetSrcQty(src common.Address, dest common.Address, dstQty *big.Int, rate *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetSrcQty(&_Reserve.CallOpts, src, dest, dstQty, rate)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() constant returns(address)
func (_Reserve *ReserveCaller) KyberNetwork(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "kyberNetwork")
	return *ret0, err
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() constant returns(address)
func (_Reserve *ReserveSession) KyberNetwork() (common.Address, error) {
	return _Reserve.Contract.KyberNetwork(&_Reserve.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() constant returns(address)
func (_Reserve *ReserveCallerSession) KyberNetwork() (common.Address, error) {
	return _Reserve.Contract.KyberNetwork(&_Reserve.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Reserve *ReserveCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Reserve *ReserveSession) PendingAdmin() (common.Address, error) {
	return _Reserve.Contract.PendingAdmin(&_Reserve.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_Reserve *ReserveCallerSession) PendingAdmin() (common.Address, error) {
	return _Reserve.Contract.PendingAdmin(&_Reserve.CallOpts)
}

// SanityRatesContract is a free data retrieval call binding the contract method 0x47e6924f.
//
// Solidity: function sanityRatesContract() constant returns(address)
func (_Reserve *ReserveCaller) SanityRatesContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "sanityRatesContract")
	return *ret0, err
}

// SanityRatesContract is a free data retrieval call binding the contract method 0x47e6924f.
//
// Solidity: function sanityRatesContract() constant returns(address)
func (_Reserve *ReserveSession) SanityRatesContract() (common.Address, error) {
	return _Reserve.Contract.SanityRatesContract(&_Reserve.CallOpts)
}

// SanityRatesContract is a free data retrieval call binding the contract method 0x47e6924f.
//
// Solidity: function sanityRatesContract() constant returns(address)
func (_Reserve *ReserveCallerSession) SanityRatesContract() (common.Address, error) {
	return _Reserve.Contract.SanityRatesContract(&_Reserve.CallOpts)
}

// TradeEnabled is a free data retrieval call binding the contract method 0xd621e813.
//
// Solidity: function tradeEnabled() constant returns(bool)
func (_Reserve *ReserveCaller) TradeEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "tradeEnabled")
	return *ret0, err
}

// TradeEnabled is a free data retrieval call binding the contract method 0xd621e813.
//
// Solidity: function tradeEnabled() constant returns(bool)
func (_Reserve *ReserveSession) TradeEnabled() (bool, error) {
	return _Reserve.Contract.TradeEnabled(&_Reserve.CallOpts)
}

// TradeEnabled is a free data retrieval call binding the contract method 0xd621e813.
//
// Solidity: function tradeEnabled() constant returns(bool)
func (_Reserve *ReserveCallerSession) TradeEnabled() (bool, error) {
	return _Reserve.Contract.TradeEnabled(&_Reserve.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_Reserve *ReserveTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_Reserve *ReserveSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddAlerter(&_Reserve.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(newAlerter address) returns()
func (_Reserve *ReserveTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddAlerter(&_Reserve.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_Reserve *ReserveTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_Reserve *ReserveSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddOperator(&_Reserve.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(newOperator address) returns()
func (_Reserve *ReserveTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddOperator(&_Reserve.TransactOpts, newOperator)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(token address, addr address, approve bool) returns()
func (_Reserve *ReserveTransactor) ApproveWithdrawAddress(opts *bind.TransactOpts, token common.Address, addr common.Address, approve bool) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "approveWithdrawAddress", token, addr, approve)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(token address, addr address, approve bool) returns()
func (_Reserve *ReserveSession) ApproveWithdrawAddress(token common.Address, addr common.Address, approve bool) (*types.Transaction, error) {
	return _Reserve.Contract.ApproveWithdrawAddress(&_Reserve.TransactOpts, token, addr, approve)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(token address, addr address, approve bool) returns()
func (_Reserve *ReserveTransactorSession) ApproveWithdrawAddress(token common.Address, addr common.Address, approve bool) (*types.Transaction, error) {
	return _Reserve.Contract.ApproveWithdrawAddress(&_Reserve.TransactOpts, token, addr, approve)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Reserve *ReserveTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Reserve *ReserveSession) ClaimAdmin() (*types.Transaction, error) {
	return _Reserve.Contract.ClaimAdmin(&_Reserve.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_Reserve *ReserveTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _Reserve.Contract.ClaimAdmin(&_Reserve.TransactOpts)
}

// DisableTrade is a paid mutator transaction binding the contract method 0x6940030f.
//
// Solidity: function disableTrade() returns(bool)
func (_Reserve *ReserveTransactor) DisableTrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "disableTrade")
}

// DisableTrade is a paid mutator transaction binding the contract method 0x6940030f.
//
// Solidity: function disableTrade() returns(bool)
func (_Reserve *ReserveSession) DisableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.DisableTrade(&_Reserve.TransactOpts)
}

// DisableTrade is a paid mutator transaction binding the contract method 0x6940030f.
//
// Solidity: function disableTrade() returns(bool)
func (_Reserve *ReserveTransactorSession) DisableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.DisableTrade(&_Reserve.TransactOpts)
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns(bool)
func (_Reserve *ReserveTransactor) EnableTrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "enableTrade")
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns(bool)
func (_Reserve *ReserveSession) EnableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.EnableTrade(&_Reserve.TransactOpts)
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns(bool)
func (_Reserve *ReserveTransactorSession) EnableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.EnableTrade(&_Reserve.TransactOpts)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_Reserve *ReserveTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_Reserve *ReserveSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveAlerter(&_Reserve.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(alerter address) returns()
func (_Reserve *ReserveTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveAlerter(&_Reserve.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_Reserve *ReserveTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_Reserve *ReserveSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveOperator(&_Reserve.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(operator address) returns()
func (_Reserve *ReserveTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveOperator(&_Reserve.TransactOpts, operator)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(_kyberNetwork address, _conversionRates address, _sanityRates address) returns()
func (_Reserve *ReserveTransactor) SetContracts(opts *bind.TransactOpts, _kyberNetwork common.Address, _conversionRates common.Address, _sanityRates common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setContracts", _kyberNetwork, _conversionRates, _sanityRates)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(_kyberNetwork address, _conversionRates address, _sanityRates address) returns()
func (_Reserve *ReserveSession) SetContracts(_kyberNetwork common.Address, _conversionRates common.Address, _sanityRates common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetContracts(&_Reserve.TransactOpts, _kyberNetwork, _conversionRates, _sanityRates)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(_kyberNetwork address, _conversionRates address, _sanityRates address) returns()
func (_Reserve *ReserveTransactorSession) SetContracts(_kyberNetwork common.Address, _conversionRates common.Address, _sanityRates common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetContracts(&_Reserve.TransactOpts, _kyberNetwork, _conversionRates, _sanityRates)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(srcToken address, srcAmount uint256, destToken address, destAddress address, conversionRate uint256, validate bool) returns(bool)
func (_Reserve *ReserveTransactor) Trade(opts *bind.TransactOpts, srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, validate bool) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "trade", srcToken, srcAmount, destToken, destAddress, conversionRate, validate)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(srcToken address, srcAmount uint256, destToken address, destAddress address, conversionRate uint256, validate bool) returns(bool)
func (_Reserve *ReserveSession) Trade(srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, validate bool) (*types.Transaction, error) {
	return _Reserve.Contract.Trade(&_Reserve.TransactOpts, srcToken, srcAmount, destToken, destAddress, conversionRate, validate)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(srcToken address, srcAmount uint256, destToken address, destAddress address, conversionRate uint256, validate bool) returns(bool)
func (_Reserve *ReserveTransactorSession) Trade(srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, validate bool) (*types.Transaction, error) {
	return _Reserve.Contract.Trade(&_Reserve.TransactOpts, srcToken, srcAmount, destToken, destAddress, conversionRate, validate)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_Reserve *ReserveTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_Reserve *ReserveSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdmin(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(newAdmin address) returns()
func (_Reserve *ReserveTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdmin(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_Reserve *ReserveTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_Reserve *ReserveSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdminQuickly(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(newAdmin address) returns()
func (_Reserve *ReserveTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdminQuickly(&_Reserve.TransactOpts, newAdmin)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(token address, amount uint256, destination address) returns(bool)
func (_Reserve *ReserveTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdraw", token, amount, destination)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(token address, amount uint256, destination address) returns(bool)
func (_Reserve *ReserveSession) Withdraw(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.Withdraw(&_Reserve.TransactOpts, token, amount, destination)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(token address, amount uint256, destination address) returns(bool)
func (_Reserve *ReserveTransactorSession) Withdraw(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.Withdraw(&_Reserve.TransactOpts, token, amount, destination)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_Reserve *ReserveTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_Reserve *ReserveSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawEther(&_Reserve.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(amount uint256, sendTo address) returns()
func (_Reserve *ReserveTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawEther(&_Reserve.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_Reserve *ReserveTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_Reserve *ReserveSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawToken(&_Reserve.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(token address, amount uint256, sendTo address) returns()
func (_Reserve *ReserveTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawToken(&_Reserve.TransactOpts, token, amount, sendTo)
}

// ReserveAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the Reserve contract.
type ReserveAdminClaimedIterator struct {
	Event *ReserveAdminClaimed // Event containing the contract specifics and raw log

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
func (it *ReserveAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveAdminClaimed)
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
		it.Event = new(ReserveAdminClaimed)
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
func (it *ReserveAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveAdminClaimed represents a AdminClaimed event raised by the Reserve contract.
type ReserveAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: e AdminClaimed(newAdmin address, previousAdmin address)
func (_Reserve *ReserveFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*ReserveAdminClaimedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &ReserveAdminClaimedIterator{contract: _Reserve.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: e AdminClaimed(newAdmin address, previousAdmin address)
func (_Reserve *ReserveFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *ReserveAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveAdminClaimed)
				if err := _Reserve.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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

// ReserveAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the Reserve contract.
type ReserveAlerterAddedIterator struct {
	Event *ReserveAlerterAdded // Event containing the contract specifics and raw log

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
func (it *ReserveAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveAlerterAdded)
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
		it.Event = new(ReserveAlerterAdded)
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
func (it *ReserveAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveAlerterAdded represents a AlerterAdded event raised by the Reserve contract.
type ReserveAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: e AlerterAdded(newAlerter address, isAdd bool)
func (_Reserve *ReserveFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*ReserveAlerterAddedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &ReserveAlerterAddedIterator{contract: _Reserve.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: e AlerterAdded(newAlerter address, isAdd bool)
func (_Reserve *ReserveFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *ReserveAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveAlerterAdded)
				if err := _Reserve.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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

// ReserveDepositTokenIterator is returned from FilterDepositToken and is used to iterate over the raw logs and unpacked data for DepositToken events raised by the Reserve contract.
type ReserveDepositTokenIterator struct {
	Event *ReserveDepositToken // Event containing the contract specifics and raw log

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
func (it *ReserveDepositTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveDepositToken)
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
		it.Event = new(ReserveDepositToken)
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
func (it *ReserveDepositTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveDepositTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveDepositToken represents a DepositToken event raised by the Reserve contract.
type ReserveDepositToken struct {
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDepositToken is a free log retrieval operation binding the contract event 0x2d0c0a8842b9944ece1495eb61121621b5e36bd6af3bba0318c695f525aef79f.
//
// Solidity: e DepositToken(token address, amount uint256)
func (_Reserve *ReserveFilterer) FilterDepositToken(opts *bind.FilterOpts) (*ReserveDepositTokenIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "DepositToken")
	if err != nil {
		return nil, err
	}
	return &ReserveDepositTokenIterator{contract: _Reserve.contract, event: "DepositToken", logs: logs, sub: sub}, nil
}

// WatchDepositToken is a free log subscription operation binding the contract event 0x2d0c0a8842b9944ece1495eb61121621b5e36bd6af3bba0318c695f525aef79f.
//
// Solidity: e DepositToken(token address, amount uint256)
func (_Reserve *ReserveFilterer) WatchDepositToken(opts *bind.WatchOpts, sink chan<- *ReserveDepositToken) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "DepositToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveDepositToken)
				if err := _Reserve.contract.UnpackLog(event, "DepositToken", log); err != nil {
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

// ReserveEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the Reserve contract.
type ReserveEtherWithdrawIterator struct {
	Event *ReserveEtherWithdraw // Event containing the contract specifics and raw log

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
func (it *ReserveEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveEtherWithdraw)
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
		it.Event = new(ReserveEtherWithdraw)
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
func (it *ReserveEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveEtherWithdraw represents a EtherWithdraw event raised by the Reserve contract.
type ReserveEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: e EtherWithdraw(amount uint256, sendTo address)
func (_Reserve *ReserveFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*ReserveEtherWithdrawIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &ReserveEtherWithdrawIterator{contract: _Reserve.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: e EtherWithdraw(amount uint256, sendTo address)
func (_Reserve *ReserveFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *ReserveEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveEtherWithdraw)
				if err := _Reserve.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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

// ReserveOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the Reserve contract.
type ReserveOperatorAddedIterator struct {
	Event *ReserveOperatorAdded // Event containing the contract specifics and raw log

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
func (it *ReserveOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveOperatorAdded)
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
		it.Event = new(ReserveOperatorAdded)
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
func (it *ReserveOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveOperatorAdded represents a OperatorAdded event raised by the Reserve contract.
type ReserveOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: e OperatorAdded(newOperator address, isAdd bool)
func (_Reserve *ReserveFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*ReserveOperatorAddedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &ReserveOperatorAddedIterator{contract: _Reserve.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: e OperatorAdded(newOperator address, isAdd bool)
func (_Reserve *ReserveFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *ReserveOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveOperatorAdded)
				if err := _Reserve.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ReserveSetContractAddressesIterator is returned from FilterSetContractAddresses and is used to iterate over the raw logs and unpacked data for SetContractAddresses events raised by the Reserve contract.
type ReserveSetContractAddressesIterator struct {
	Event *ReserveSetContractAddresses // Event containing the contract specifics and raw log

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
func (it *ReserveSetContractAddressesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveSetContractAddresses)
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
		it.Event = new(ReserveSetContractAddresses)
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
func (it *ReserveSetContractAddressesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveSetContractAddressesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveSetContractAddresses represents a SetContractAddresses event raised by the Reserve contract.
type ReserveSetContractAddresses struct {
	Network common.Address
	Rate    common.Address
	Sanity  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSetContractAddresses is a free log retrieval operation binding the contract event 0x7a85322644a4462d8ff5482d2a841a4d231f8cfb3c9f4a50f66f8b2bd568c31c.
//
// Solidity: e SetContractAddresses(network address, rate address, sanity address)
func (_Reserve *ReserveFilterer) FilterSetContractAddresses(opts *bind.FilterOpts) (*ReserveSetContractAddressesIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "SetContractAddresses")
	if err != nil {
		return nil, err
	}
	return &ReserveSetContractAddressesIterator{contract: _Reserve.contract, event: "SetContractAddresses", logs: logs, sub: sub}, nil
}

// WatchSetContractAddresses is a free log subscription operation binding the contract event 0x7a85322644a4462d8ff5482d2a841a4d231f8cfb3c9f4a50f66f8b2bd568c31c.
//
// Solidity: e SetContractAddresses(network address, rate address, sanity address)
func (_Reserve *ReserveFilterer) WatchSetContractAddresses(opts *bind.WatchOpts, sink chan<- *ReserveSetContractAddresses) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "SetContractAddresses")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveSetContractAddresses)
				if err := _Reserve.contract.UnpackLog(event, "SetContractAddresses", log); err != nil {
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

// ReserveTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the Reserve contract.
type ReserveTokenWithdrawIterator struct {
	Event *ReserveTokenWithdraw // Event containing the contract specifics and raw log

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
func (it *ReserveTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveTokenWithdraw)
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
		it.Event = new(ReserveTokenWithdraw)
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
func (it *ReserveTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveTokenWithdraw represents a TokenWithdraw event raised by the Reserve contract.
type ReserveTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: e TokenWithdraw(token address, amount uint256, sendTo address)
func (_Reserve *ReserveFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*ReserveTokenWithdrawIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &ReserveTokenWithdrawIterator{contract: _Reserve.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: e TokenWithdraw(token address, amount uint256, sendTo address)
func (_Reserve *ReserveFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *ReserveTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveTokenWithdraw)
				if err := _Reserve.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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

// ReserveTradeEnabledIterator is returned from FilterTradeEnabled and is used to iterate over the raw logs and unpacked data for TradeEnabled events raised by the Reserve contract.
type ReserveTradeEnabledIterator struct {
	Event *ReserveTradeEnabled // Event containing the contract specifics and raw log

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
func (it *ReserveTradeEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveTradeEnabled)
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
		it.Event = new(ReserveTradeEnabled)
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
func (it *ReserveTradeEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveTradeEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveTradeEnabled represents a TradeEnabled event raised by the Reserve contract.
type ReserveTradeEnabled struct {
	Enable bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTradeEnabled is a free log retrieval operation binding the contract event 0x7d7f00509dd73ac4449f698ae75ccc797895eff5fa9d446d3df387598a26e735.
//
// Solidity: e TradeEnabled(enable bool)
func (_Reserve *ReserveFilterer) FilterTradeEnabled(opts *bind.FilterOpts) (*ReserveTradeEnabledIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TradeEnabled")
	if err != nil {
		return nil, err
	}
	return &ReserveTradeEnabledIterator{contract: _Reserve.contract, event: "TradeEnabled", logs: logs, sub: sub}, nil
}

// WatchTradeEnabled is a free log subscription operation binding the contract event 0x7d7f00509dd73ac4449f698ae75ccc797895eff5fa9d446d3df387598a26e735.
//
// Solidity: e TradeEnabled(enable bool)
func (_Reserve *ReserveFilterer) WatchTradeEnabled(opts *bind.WatchOpts, sink chan<- *ReserveTradeEnabled) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "TradeEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveTradeEnabled)
				if err := _Reserve.contract.UnpackLog(event, "TradeEnabled", log); err != nil {
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

// ReserveTradeExecuteIterator is returned from FilterTradeExecute and is used to iterate over the raw logs and unpacked data for TradeExecute events raised by the Reserve contract.
type ReserveTradeExecuteIterator struct {
	Event *ReserveTradeExecute // Event containing the contract specifics and raw log

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
func (it *ReserveTradeExecuteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveTradeExecute)
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
		it.Event = new(ReserveTradeExecute)
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
func (it *ReserveTradeExecuteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveTradeExecuteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveTradeExecute represents a TradeExecute event raised by the Reserve contract.
type ReserveTradeExecute struct {
	Origin      common.Address
	Src         common.Address
	SrcAmount   *big.Int
	DestToken   common.Address
	DestAmount  *big.Int
	DestAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTradeExecute is a free log retrieval operation binding the contract event 0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de.
//
// Solidity: e TradeExecute(origin indexed address, src address, srcAmount uint256, destToken address, destAmount uint256, destAddress address)
func (_Reserve *ReserveFilterer) FilterTradeExecute(opts *bind.FilterOpts, origin []common.Address) (*ReserveTradeExecuteIterator, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TradeExecute", originRule)
	if err != nil {
		return nil, err
	}
	return &ReserveTradeExecuteIterator{contract: _Reserve.contract, event: "TradeExecute", logs: logs, sub: sub}, nil
}

// WatchTradeExecute is a free log subscription operation binding the contract event 0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de.
//
// Solidity: e TradeExecute(origin indexed address, src address, srcAmount uint256, destToken address, destAmount uint256, destAddress address)
func (_Reserve *ReserveFilterer) WatchTradeExecute(opts *bind.WatchOpts, sink chan<- *ReserveTradeExecute, origin []common.Address) (event.Subscription, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "TradeExecute", originRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveTradeExecute)
				if err := _Reserve.contract.UnpackLog(event, "TradeExecute", log); err != nil {
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

// ReserveTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the Reserve contract.
type ReserveTransferAdminPendingIterator struct {
	Event *ReserveTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *ReserveTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveTransferAdminPending)
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
		it.Event = new(ReserveTransferAdminPending)
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
func (it *ReserveTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveTransferAdminPending represents a TransferAdminPending event raised by the Reserve contract.
type ReserveTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: e TransferAdminPending(pendingAdmin address)
func (_Reserve *ReserveFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*ReserveTransferAdminPendingIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &ReserveTransferAdminPendingIterator{contract: _Reserve.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: e TransferAdminPending(pendingAdmin address)
func (_Reserve *ReserveFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *ReserveTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveTransferAdminPending)
				if err := _Reserve.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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

// ReserveWithdrawAddressApprovedIterator is returned from FilterWithdrawAddressApproved and is used to iterate over the raw logs and unpacked data for WithdrawAddressApproved events raised by the Reserve contract.
type ReserveWithdrawAddressApprovedIterator struct {
	Event *ReserveWithdrawAddressApproved // Event containing the contract specifics and raw log

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
func (it *ReserveWithdrawAddressApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveWithdrawAddressApproved)
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
		it.Event = new(ReserveWithdrawAddressApproved)
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
func (it *ReserveWithdrawAddressApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveWithdrawAddressApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveWithdrawAddressApproved represents a WithdrawAddressApproved event raised by the Reserve contract.
type ReserveWithdrawAddressApproved struct {
	Token   common.Address
	Addr    common.Address
	Approve bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawAddressApproved is a free log retrieval operation binding the contract event 0xd5fd5351efae1f4bb760079da9f0ff9589e2c3e216337ca9d39cdff573b245c4.
//
// Solidity: e WithdrawAddressApproved(token address, addr address, approve bool)
func (_Reserve *ReserveFilterer) FilterWithdrawAddressApproved(opts *bind.FilterOpts) (*ReserveWithdrawAddressApprovedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "WithdrawAddressApproved")
	if err != nil {
		return nil, err
	}
	return &ReserveWithdrawAddressApprovedIterator{contract: _Reserve.contract, event: "WithdrawAddressApproved", logs: logs, sub: sub}, nil
}

// WatchWithdrawAddressApproved is a free log subscription operation binding the contract event 0xd5fd5351efae1f4bb760079da9f0ff9589e2c3e216337ca9d39cdff573b245c4.
//
// Solidity: e WithdrawAddressApproved(token address, addr address, approve bool)
func (_Reserve *ReserveFilterer) WatchWithdrawAddressApproved(opts *bind.WatchOpts, sink chan<- *ReserveWithdrawAddressApproved) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "WithdrawAddressApproved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveWithdrawAddressApproved)
				if err := _Reserve.contract.UnpackLog(event, "WithdrawAddressApproved", log); err != nil {
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

// ReserveWithdrawFundsIterator is returned from FilterWithdrawFunds and is used to iterate over the raw logs and unpacked data for WithdrawFunds events raised by the Reserve contract.
type ReserveWithdrawFundsIterator struct {
	Event *ReserveWithdrawFunds // Event containing the contract specifics and raw log

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
func (it *ReserveWithdrawFundsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveWithdrawFunds)
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
		it.Event = new(ReserveWithdrawFunds)
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
func (it *ReserveWithdrawFundsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveWithdrawFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveWithdrawFunds represents a WithdrawFunds event raised by the Reserve contract.
type ReserveWithdrawFunds struct {
	Token       common.Address
	Amount      *big.Int
	Destination common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawFunds is a free log retrieval operation binding the contract event 0xb67719fc33c1f17d31bf3a698690d62066b1e0bae28fcd3c56cf2c015c2863d6.
//
// Solidity: e WithdrawFunds(token address, amount uint256, destination address)
func (_Reserve *ReserveFilterer) FilterWithdrawFunds(opts *bind.FilterOpts) (*ReserveWithdrawFundsIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "WithdrawFunds")
	if err != nil {
		return nil, err
	}
	return &ReserveWithdrawFundsIterator{contract: _Reserve.contract, event: "WithdrawFunds", logs: logs, sub: sub}, nil
}

// WatchWithdrawFunds is a free log subscription operation binding the contract event 0xb67719fc33c1f17d31bf3a698690d62066b1e0bae28fcd3c56cf2c015c2863d6.
//
// Solidity: e WithdrawFunds(token address, amount uint256, destination address)
func (_Reserve *ReserveFilterer) WatchWithdrawFunds(opts *bind.WatchOpts, sink chan<- *ReserveWithdrawFunds) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "WithdrawFunds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveWithdrawFunds)
				if err := _Reserve.contract.UnpackLog(event, "WithdrawFunds", log); err != nil {
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
