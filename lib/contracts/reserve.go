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
const ReserveABI = "[{\"inputs\":[{\"internalType\":\"contractIConversionRates\",\"name\":\"_ratesContract\",\"type\":\"address\"},{\"internalType\":\"contractIWeth\",\"name\":\"_weth\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Ext\",\"name\":\"_quoteToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"_maxGasPriceWei\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"newMaxGasPrice\",\"type\":\"uint128\"}],\"name\":\"MaxGasPriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"NewTokenWallet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIConversionRates\",\"name\":\"rate\",\"type\":\"address\"}],\"name\":\"SetConversionRateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIKyberSanity\",\"name\":\"sanity\",\"type\":\"address\"}],\"name\":\"SetSanityRateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIWeth\",\"name\":\"weth\",\"type\":\"address\"}],\"name\":\"SetWethAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"TradeEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"destToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"}],\"name\":\"TradeExecute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approve\",\"type\":\"bool\"}],\"name\":\"WithdrawAddressApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"WithdrawFunds\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BPS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approve\",\"type\":\"bool\"}],\"name\":\"approveWithdrawAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"approvedWithdrawAddresses\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"conversionRatesContract\",\"outputs\":[{\"internalType\":\"contractIConversionRates\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Ext\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcQty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getConversionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAddressApprovedForWithdrawal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxGasPriceWei\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quoteToken\",\"outputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sanityRatesContract\",\"outputs\":[{\"internalType\":\"contractIKyberSanity\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIConversionRates\",\"name\":\"_newConversionRate\",\"type\":\"address\"}],\"name\":\"setConversionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"newMaxGasPrice\",\"type\":\"uint128\"}],\"name\":\"setMaxGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberSanity\",\"name\":\"_newSanity\",\"type\":\"address\"}],\"name\":\"setSanityRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"setTokenWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIWeth\",\"name\":\"_newWeth\",\"type\":\"address\"}],\"name\":\"setWeth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"srcToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20Ext\",\"name\":\"destToken\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"conversionRate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"name\":\"trade\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tradeEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"weth\",\"outputs\":[{\"internalType\":\"contractIWeth\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Ext\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

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

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_Reserve *ReserveCaller) BPS(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "BPS")
	return *ret0, err
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_Reserve *ReserveSession) BPS() (*big.Int, error) {
	return _Reserve.Contract.BPS(&_Reserve.CallOpts)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_Reserve *ReserveCallerSession) BPS() (*big.Int, error) {
	return _Reserve.Contract.BPS(&_Reserve.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
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
// Solidity: function admin() view returns(address)
func (_Reserve *ReserveSession) Admin() (common.Address, error) {
	return _Reserve.Contract.Admin(&_Reserve.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Reserve *ReserveCallerSession) Admin() (common.Address, error) {
	return _Reserve.Contract.Admin(&_Reserve.CallOpts)
}

// ApprovedWithdrawAddresses is a free data retrieval call binding the contract method 0xd7b7024d.
//
// Solidity: function approvedWithdrawAddresses(bytes32 ) view returns(bool)
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
// Solidity: function approvedWithdrawAddresses(bytes32 ) view returns(bool)
func (_Reserve *ReserveSession) ApprovedWithdrawAddresses(arg0 [32]byte) (bool, error) {
	return _Reserve.Contract.ApprovedWithdrawAddresses(&_Reserve.CallOpts, arg0)
}

// ApprovedWithdrawAddresses is a free data retrieval call binding the contract method 0xd7b7024d.
//
// Solidity: function approvedWithdrawAddresses(bytes32 ) view returns(bool)
func (_Reserve *ReserveCallerSession) ApprovedWithdrawAddresses(arg0 [32]byte) (bool, error) {
	return _Reserve.Contract.ApprovedWithdrawAddresses(&_Reserve.CallOpts, arg0)
}

// ConversionRatesContract is a free data retrieval call binding the contract method 0xd5847d33.
//
// Solidity: function conversionRatesContract() view returns(address)
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
// Solidity: function conversionRatesContract() view returns(address)
func (_Reserve *ReserveSession) ConversionRatesContract() (common.Address, error) {
	return _Reserve.Contract.ConversionRatesContract(&_Reserve.CallOpts)
}

// ConversionRatesContract is a free data retrieval call binding the contract method 0xd5847d33.
//
// Solidity: function conversionRatesContract() view returns(address)
func (_Reserve *ReserveCallerSession) ConversionRatesContract() (common.Address, error) {
	return _Reserve.Contract.ConversionRatesContract(&_Reserve.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
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
// Solidity: function getAlerters() view returns(address[])
func (_Reserve *ReserveSession) GetAlerters() ([]common.Address, error) {
	return _Reserve.Contract.GetAlerters(&_Reserve.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_Reserve *ReserveCallerSession) GetAlerters() ([]common.Address, error) {
	return _Reserve.Contract.GetAlerters(&_Reserve.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address token) view returns(uint256)
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
// Solidity: function getBalance(address token) view returns(uint256)
func (_Reserve *ReserveSession) GetBalance(token common.Address) (*big.Int, error) {
	return _Reserve.Contract.GetBalance(&_Reserve.CallOpts, token)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address token) view returns(uint256)
func (_Reserve *ReserveCallerSession) GetBalance(token common.Address) (*big.Int, error) {
	return _Reserve.Contract.GetBalance(&_Reserve.CallOpts, token)
}

// GetConversionRate is a free data retrieval call binding the contract method 0x7cd44272.
//
// Solidity: function getConversionRate(address src, address dest, uint256 srcQty, uint256 blockNumber) view returns(uint256)
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
// Solidity: function getConversionRate(address src, address dest, uint256 srcQty, uint256 blockNumber) view returns(uint256)
func (_Reserve *ReserveSession) GetConversionRate(src common.Address, dest common.Address, srcQty *big.Int, blockNumber *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetConversionRate(&_Reserve.CallOpts, src, dest, srcQty, blockNumber)
}

// GetConversionRate is a free data retrieval call binding the contract method 0x7cd44272.
//
// Solidity: function getConversionRate(address src, address dest, uint256 srcQty, uint256 blockNumber) view returns(uint256)
func (_Reserve *ReserveCallerSession) GetConversionRate(src common.Address, dest common.Address, srcQty *big.Int, blockNumber *big.Int) (*big.Int, error) {
	return _Reserve.Contract.GetConversionRate(&_Reserve.CallOpts, src, dest, srcQty, blockNumber)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
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
// Solidity: function getOperators() view returns(address[])
func (_Reserve *ReserveSession) GetOperators() ([]common.Address, error) {
	return _Reserve.Contract.GetOperators(&_Reserve.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_Reserve *ReserveCallerSession) GetOperators() ([]common.Address, error) {
	return _Reserve.Contract.GetOperators(&_Reserve.CallOpts)
}

// GetTokenWallet is a free data retrieval call binding the contract method 0x85d75025.
//
// Solidity: function getTokenWallet(address token) view returns(address wallet)
func (_Reserve *ReserveCaller) GetTokenWallet(opts *bind.CallOpts, token common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "getTokenWallet", token)
	return *ret0, err
}

// GetTokenWallet is a free data retrieval call binding the contract method 0x85d75025.
//
// Solidity: function getTokenWallet(address token) view returns(address wallet)
func (_Reserve *ReserveSession) GetTokenWallet(token common.Address) (common.Address, error) {
	return _Reserve.Contract.GetTokenWallet(&_Reserve.CallOpts, token)
}

// GetTokenWallet is a free data retrieval call binding the contract method 0x85d75025.
//
// Solidity: function getTokenWallet(address token) view returns(address wallet)
func (_Reserve *ReserveCallerSession) GetTokenWallet(token common.Address) (common.Address, error) {
	return _Reserve.Contract.GetTokenWallet(&_Reserve.CallOpts, token)
}

// IsAddressApprovedForWithdrawal is a free data retrieval call binding the contract method 0xa56bb95b.
//
// Solidity: function isAddressApprovedForWithdrawal(address token, address addr) view returns(bool)
func (_Reserve *ReserveCaller) IsAddressApprovedForWithdrawal(opts *bind.CallOpts, token common.Address, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "isAddressApprovedForWithdrawal", token, addr)
	return *ret0, err
}

// IsAddressApprovedForWithdrawal is a free data retrieval call binding the contract method 0xa56bb95b.
//
// Solidity: function isAddressApprovedForWithdrawal(address token, address addr) view returns(bool)
func (_Reserve *ReserveSession) IsAddressApprovedForWithdrawal(token common.Address, addr common.Address) (bool, error) {
	return _Reserve.Contract.IsAddressApprovedForWithdrawal(&_Reserve.CallOpts, token, addr)
}

// IsAddressApprovedForWithdrawal is a free data retrieval call binding the contract method 0xa56bb95b.
//
// Solidity: function isAddressApprovedForWithdrawal(address token, address addr) view returns(bool)
func (_Reserve *ReserveCallerSession) IsAddressApprovedForWithdrawal(token common.Address, addr common.Address) (bool, error) {
	return _Reserve.Contract.IsAddressApprovedForWithdrawal(&_Reserve.CallOpts, token, addr)
}

// MaxGasPriceWei is a free data retrieval call binding the contract method 0xef3881c8.
//
// Solidity: function maxGasPriceWei() view returns(uint128)
func (_Reserve *ReserveCaller) MaxGasPriceWei(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "maxGasPriceWei")
	return *ret0, err
}

// MaxGasPriceWei is a free data retrieval call binding the contract method 0xef3881c8.
//
// Solidity: function maxGasPriceWei() view returns(uint128)
func (_Reserve *ReserveSession) MaxGasPriceWei() (*big.Int, error) {
	return _Reserve.Contract.MaxGasPriceWei(&_Reserve.CallOpts)
}

// MaxGasPriceWei is a free data retrieval call binding the contract method 0xef3881c8.
//
// Solidity: function maxGasPriceWei() view returns(uint128)
func (_Reserve *ReserveCallerSession) MaxGasPriceWei() (*big.Int, error) {
	return _Reserve.Contract.MaxGasPriceWei(&_Reserve.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
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
// Solidity: function pendingAdmin() view returns(address)
func (_Reserve *ReserveSession) PendingAdmin() (common.Address, error) {
	return _Reserve.Contract.PendingAdmin(&_Reserve.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Reserve *ReserveCallerSession) PendingAdmin() (common.Address, error) {
	return _Reserve.Contract.PendingAdmin(&_Reserve.CallOpts)
}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_Reserve *ReserveCaller) QuoteToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "quoteToken")
	return *ret0, err
}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_Reserve *ReserveSession) QuoteToken() (common.Address, error) {
	return _Reserve.Contract.QuoteToken(&_Reserve.CallOpts)
}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_Reserve *ReserveCallerSession) QuoteToken() (common.Address, error) {
	return _Reserve.Contract.QuoteToken(&_Reserve.CallOpts)
}

// SanityRatesContract is a free data retrieval call binding the contract method 0x47e6924f.
//
// Solidity: function sanityRatesContract() view returns(address)
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
// Solidity: function sanityRatesContract() view returns(address)
func (_Reserve *ReserveSession) SanityRatesContract() (common.Address, error) {
	return _Reserve.Contract.SanityRatesContract(&_Reserve.CallOpts)
}

// SanityRatesContract is a free data retrieval call binding the contract method 0x47e6924f.
//
// Solidity: function sanityRatesContract() view returns(address)
func (_Reserve *ReserveCallerSession) SanityRatesContract() (common.Address, error) {
	return _Reserve.Contract.SanityRatesContract(&_Reserve.CallOpts)
}

// TokenWallet is a free data retrieval call binding the contract method 0xa80cbac6.
//
// Solidity: function tokenWallet(address ) view returns(address)
func (_Reserve *ReserveCaller) TokenWallet(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "tokenWallet", arg0)
	return *ret0, err
}

// TokenWallet is a free data retrieval call binding the contract method 0xa80cbac6.
//
// Solidity: function tokenWallet(address ) view returns(address)
func (_Reserve *ReserveSession) TokenWallet(arg0 common.Address) (common.Address, error) {
	return _Reserve.Contract.TokenWallet(&_Reserve.CallOpts, arg0)
}

// TokenWallet is a free data retrieval call binding the contract method 0xa80cbac6.
//
// Solidity: function tokenWallet(address ) view returns(address)
func (_Reserve *ReserveCallerSession) TokenWallet(arg0 common.Address) (common.Address, error) {
	return _Reserve.Contract.TokenWallet(&_Reserve.CallOpts, arg0)
}

// TradeEnabled is a free data retrieval call binding the contract method 0xd621e813.
//
// Solidity: function tradeEnabled() view returns(bool)
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
// Solidity: function tradeEnabled() view returns(bool)
func (_Reserve *ReserveSession) TradeEnabled() (bool, error) {
	return _Reserve.Contract.TradeEnabled(&_Reserve.CallOpts)
}

// TradeEnabled is a free data retrieval call binding the contract method 0xd621e813.
//
// Solidity: function tradeEnabled() view returns(bool)
func (_Reserve *ReserveCallerSession) TradeEnabled() (bool, error) {
	return _Reserve.Contract.TradeEnabled(&_Reserve.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Reserve *ReserveCaller) Weth(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Reserve.contract.Call(opts, out, "weth")
	return *ret0, err
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Reserve *ReserveSession) Weth() (common.Address, error) {
	return _Reserve.Contract.Weth(&_Reserve.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Reserve *ReserveCallerSession) Weth() (common.Address, error) {
	return _Reserve.Contract.Weth(&_Reserve.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Reserve *ReserveTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Reserve *ReserveSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddAlerter(&_Reserve.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_Reserve *ReserveTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddAlerter(&_Reserve.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Reserve *ReserveTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Reserve *ReserveSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddOperator(&_Reserve.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_Reserve *ReserveTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.AddOperator(&_Reserve.TransactOpts, newOperator)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(address token, address addr, bool approve) returns()
func (_Reserve *ReserveTransactor) ApproveWithdrawAddress(opts *bind.TransactOpts, token common.Address, addr common.Address, approve bool) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "approveWithdrawAddress", token, addr, approve)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(address token, address addr, bool approve) returns()
func (_Reserve *ReserveSession) ApproveWithdrawAddress(token common.Address, addr common.Address, approve bool) (*types.Transaction, error) {
	return _Reserve.Contract.ApproveWithdrawAddress(&_Reserve.TransactOpts, token, addr, approve)
}

// ApproveWithdrawAddress is a paid mutator transaction binding the contract method 0x546dc71c.
//
// Solidity: function approveWithdrawAddress(address token, address addr, bool approve) returns()
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
// Solidity: function disableTrade() returns()
func (_Reserve *ReserveTransactor) DisableTrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "disableTrade")
}

// DisableTrade is a paid mutator transaction binding the contract method 0x6940030f.
//
// Solidity: function disableTrade() returns()
func (_Reserve *ReserveSession) DisableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.DisableTrade(&_Reserve.TransactOpts)
}

// DisableTrade is a paid mutator transaction binding the contract method 0x6940030f.
//
// Solidity: function disableTrade() returns()
func (_Reserve *ReserveTransactorSession) DisableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.DisableTrade(&_Reserve.TransactOpts)
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns()
func (_Reserve *ReserveTransactor) EnableTrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "enableTrade")
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns()
func (_Reserve *ReserveSession) EnableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.EnableTrade(&_Reserve.TransactOpts)
}

// EnableTrade is a paid mutator transaction binding the contract method 0x0099d386.
//
// Solidity: function enableTrade() returns()
func (_Reserve *ReserveTransactorSession) EnableTrade() (*types.Transaction, error) {
	return _Reserve.Contract.EnableTrade(&_Reserve.TransactOpts)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Reserve *ReserveTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Reserve *ReserveSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveAlerter(&_Reserve.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_Reserve *ReserveTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveAlerter(&_Reserve.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Reserve *ReserveTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Reserve *ReserveSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveOperator(&_Reserve.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_Reserve *ReserveTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.RemoveOperator(&_Reserve.TransactOpts, operator)
}

// SetConversionRate is a paid mutator transaction binding the contract method 0xfa307281.
//
// Solidity: function setConversionRate(address _newConversionRate) returns()
func (_Reserve *ReserveTransactor) SetConversionRate(opts *bind.TransactOpts, _newConversionRate common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setConversionRate", _newConversionRate)
}

// SetConversionRate is a paid mutator transaction binding the contract method 0xfa307281.
//
// Solidity: function setConversionRate(address _newConversionRate) returns()
func (_Reserve *ReserveSession) SetConversionRate(_newConversionRate common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetConversionRate(&_Reserve.TransactOpts, _newConversionRate)
}

// SetConversionRate is a paid mutator transaction binding the contract method 0xfa307281.
//
// Solidity: function setConversionRate(address _newConversionRate) returns()
func (_Reserve *ReserveTransactorSession) SetConversionRate(_newConversionRate common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetConversionRate(&_Reserve.TransactOpts, _newConversionRate)
}

// SetMaxGasPrice is a paid mutator transaction binding the contract method 0xcac1d649.
//
// Solidity: function setMaxGasPrice(uint128 newMaxGasPrice) returns()
func (_Reserve *ReserveTransactor) SetMaxGasPrice(opts *bind.TransactOpts, newMaxGasPrice *big.Int) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setMaxGasPrice", newMaxGasPrice)
}

// SetMaxGasPrice is a paid mutator transaction binding the contract method 0xcac1d649.
//
// Solidity: function setMaxGasPrice(uint128 newMaxGasPrice) returns()
func (_Reserve *ReserveSession) SetMaxGasPrice(newMaxGasPrice *big.Int) (*types.Transaction, error) {
	return _Reserve.Contract.SetMaxGasPrice(&_Reserve.TransactOpts, newMaxGasPrice)
}

// SetMaxGasPrice is a paid mutator transaction binding the contract method 0xcac1d649.
//
// Solidity: function setMaxGasPrice(uint128 newMaxGasPrice) returns()
func (_Reserve *ReserveTransactorSession) SetMaxGasPrice(newMaxGasPrice *big.Int) (*types.Transaction, error) {
	return _Reserve.Contract.SetMaxGasPrice(&_Reserve.TransactOpts, newMaxGasPrice)
}

// SetSanityRate is a paid mutator transaction binding the contract method 0x80f4da8b.
//
// Solidity: function setSanityRate(address _newSanity) returns()
func (_Reserve *ReserveTransactor) SetSanityRate(opts *bind.TransactOpts, _newSanity common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setSanityRate", _newSanity)
}

// SetSanityRate is a paid mutator transaction binding the contract method 0x80f4da8b.
//
// Solidity: function setSanityRate(address _newSanity) returns()
func (_Reserve *ReserveSession) SetSanityRate(_newSanity common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetSanityRate(&_Reserve.TransactOpts, _newSanity)
}

// SetSanityRate is a paid mutator transaction binding the contract method 0x80f4da8b.
//
// Solidity: function setSanityRate(address _newSanity) returns()
func (_Reserve *ReserveTransactorSession) SetSanityRate(_newSanity common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetSanityRate(&_Reserve.TransactOpts, _newSanity)
}

// SetTokenWallet is a paid mutator transaction binding the contract method 0x1bc7bfec.
//
// Solidity: function setTokenWallet(address token, address wallet) returns()
func (_Reserve *ReserveTransactor) SetTokenWallet(opts *bind.TransactOpts, token common.Address, wallet common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setTokenWallet", token, wallet)
}

// SetTokenWallet is a paid mutator transaction binding the contract method 0x1bc7bfec.
//
// Solidity: function setTokenWallet(address token, address wallet) returns()
func (_Reserve *ReserveSession) SetTokenWallet(token common.Address, wallet common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetTokenWallet(&_Reserve.TransactOpts, token, wallet)
}

// SetTokenWallet is a paid mutator transaction binding the contract method 0x1bc7bfec.
//
// Solidity: function setTokenWallet(address token, address wallet) returns()
func (_Reserve *ReserveTransactorSession) SetTokenWallet(token common.Address, wallet common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetTokenWallet(&_Reserve.TransactOpts, token, wallet)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _newWeth) returns()
func (_Reserve *ReserveTransactor) SetWeth(opts *bind.TransactOpts, _newWeth common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "setWeth", _newWeth)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _newWeth) returns()
func (_Reserve *ReserveSession) SetWeth(_newWeth common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetWeth(&_Reserve.TransactOpts, _newWeth)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _newWeth) returns()
func (_Reserve *ReserveTransactorSession) SetWeth(_newWeth common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.SetWeth(&_Reserve.TransactOpts, _newWeth)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(address srcToken, uint256 srcAmount, address destToken, address destAddress, uint256 conversionRate, bool ) payable returns(bool)
func (_Reserve *ReserveTransactor) Trade(opts *bind.TransactOpts, srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, arg5 bool) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "trade", srcToken, srcAmount, destToken, destAddress, conversionRate, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(address srcToken, uint256 srcAmount, address destToken, address destAddress, uint256 conversionRate, bool ) payable returns(bool)
func (_Reserve *ReserveSession) Trade(srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, arg5 bool) (*types.Transaction, error) {
	return _Reserve.Contract.Trade(&_Reserve.TransactOpts, srcToken, srcAmount, destToken, destAddress, conversionRate, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x6cf69811.
//
// Solidity: function trade(address srcToken, uint256 srcAmount, address destToken, address destAddress, uint256 conversionRate, bool ) payable returns(bool)
func (_Reserve *ReserveTransactorSession) Trade(srcToken common.Address, srcAmount *big.Int, destToken common.Address, destAddress common.Address, conversionRate *big.Int, arg5 bool) (*types.Transaction, error) {
	return _Reserve.Contract.Trade(&_Reserve.TransactOpts, srcToken, srcAmount, destToken, destAddress, conversionRate, arg5)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Reserve *ReserveTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Reserve *ReserveSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdmin(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_Reserve *ReserveTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdmin(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Reserve *ReserveTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Reserve *ReserveSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdminQuickly(&_Reserve.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_Reserve *ReserveTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.TransferAdminQuickly(&_Reserve.TransactOpts, newAdmin)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address token, uint256 amount, address destination) returns()
func (_Reserve *ReserveTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdraw", token, amount, destination)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address token, uint256 amount, address destination) returns()
func (_Reserve *ReserveSession) Withdraw(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.Withdraw(&_Reserve.TransactOpts, token, amount, destination)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address token, uint256 amount, address destination) returns()
func (_Reserve *ReserveTransactorSession) Withdraw(token common.Address, amount *big.Int, destination common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.Withdraw(&_Reserve.TransactOpts, token, amount, destination)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Reserve *ReserveTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Reserve *ReserveSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawEther(&_Reserve.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_Reserve *ReserveTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawEther(&_Reserve.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Reserve *ReserveTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Reserve *ReserveSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawToken(&_Reserve.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_Reserve *ReserveTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _Reserve.Contract.WithdrawToken(&_Reserve.TransactOpts, token, amount, sendTo)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reserve *ReserveTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reserve.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reserve *ReserveSession) Receive() (*types.Transaction, error) {
	return _Reserve.Contract.Receive(&_Reserve.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reserve *ReserveTransactorSession) Receive() (*types.Transaction, error) {
	return _Reserve.Contract.Receive(&_Reserve.TransactOpts)
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
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_Reserve *ReserveFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*ReserveAdminClaimedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &ReserveAdminClaimedIterator{contract: _Reserve.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
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

// ParseAdminClaimed is a log parse operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_Reserve *ReserveFilterer) ParseAdminClaimed(log types.Log) (*ReserveAdminClaimed, error) {
	event := new(ReserveAdminClaimed)
	if err := _Reserve.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_Reserve *ReserveFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*ReserveAlerterAddedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &ReserveAlerterAddedIterator{contract: _Reserve.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
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

// ParseAlerterAdded is a log parse operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_Reserve *ReserveFilterer) ParseAlerterAdded(log types.Log) (*ReserveAlerterAdded, error) {
	event := new(ReserveAlerterAdded)
	if err := _Reserve.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event DepositToken(address indexed token, uint256 amount)
func (_Reserve *ReserveFilterer) FilterDepositToken(opts *bind.FilterOpts, token []common.Address) (*ReserveDepositTokenIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "DepositToken", tokenRule)
	if err != nil {
		return nil, err
	}
	return &ReserveDepositTokenIterator{contract: _Reserve.contract, event: "DepositToken", logs: logs, sub: sub}, nil
}

// WatchDepositToken is a free log subscription operation binding the contract event 0x2d0c0a8842b9944ece1495eb61121621b5e36bd6af3bba0318c695f525aef79f.
//
// Solidity: event DepositToken(address indexed token, uint256 amount)
func (_Reserve *ReserveFilterer) WatchDepositToken(opts *bind.WatchOpts, sink chan<- *ReserveDepositToken, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "DepositToken", tokenRule)
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

// ParseDepositToken is a log parse operation binding the contract event 0x2d0c0a8842b9944ece1495eb61121621b5e36bd6af3bba0318c695f525aef79f.
//
// Solidity: event DepositToken(address indexed token, uint256 amount)
func (_Reserve *ReserveFilterer) ParseDepositToken(log types.Log) (*ReserveDepositToken, error) {
	event := new(ReserveDepositToken)
	if err := _Reserve.contract.UnpackLog(event, "DepositToken", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_Reserve *ReserveFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*ReserveEtherWithdrawIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &ReserveEtherWithdrawIterator{contract: _Reserve.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
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

// ParseEtherWithdraw is a log parse operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_Reserve *ReserveFilterer) ParseEtherWithdraw(log types.Log) (*ReserveEtherWithdraw, error) {
	event := new(ReserveEtherWithdraw)
	if err := _Reserve.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ReserveMaxGasPriceUpdatedIterator is returned from FilterMaxGasPriceUpdated and is used to iterate over the raw logs and unpacked data for MaxGasPriceUpdated events raised by the Reserve contract.
type ReserveMaxGasPriceUpdatedIterator struct {
	Event *ReserveMaxGasPriceUpdated // Event containing the contract specifics and raw log

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
func (it *ReserveMaxGasPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveMaxGasPriceUpdated)
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
		it.Event = new(ReserveMaxGasPriceUpdated)
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
func (it *ReserveMaxGasPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveMaxGasPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveMaxGasPriceUpdated represents a MaxGasPriceUpdated event raised by the Reserve contract.
type ReserveMaxGasPriceUpdated struct {
	NewMaxGasPrice *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMaxGasPriceUpdated is a free log retrieval operation binding the contract event 0x951ddb0e961044819fc0750a51148b147386ae97b72d9b6763f9b943de116e32.
//
// Solidity: event MaxGasPriceUpdated(uint128 newMaxGasPrice)
func (_Reserve *ReserveFilterer) FilterMaxGasPriceUpdated(opts *bind.FilterOpts) (*ReserveMaxGasPriceUpdatedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "MaxGasPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &ReserveMaxGasPriceUpdatedIterator{contract: _Reserve.contract, event: "MaxGasPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxGasPriceUpdated is a free log subscription operation binding the contract event 0x951ddb0e961044819fc0750a51148b147386ae97b72d9b6763f9b943de116e32.
//
// Solidity: event MaxGasPriceUpdated(uint128 newMaxGasPrice)
func (_Reserve *ReserveFilterer) WatchMaxGasPriceUpdated(opts *bind.WatchOpts, sink chan<- *ReserveMaxGasPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "MaxGasPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveMaxGasPriceUpdated)
				if err := _Reserve.contract.UnpackLog(event, "MaxGasPriceUpdated", log); err != nil {
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

// ParseMaxGasPriceUpdated is a log parse operation binding the contract event 0x951ddb0e961044819fc0750a51148b147386ae97b72d9b6763f9b943de116e32.
//
// Solidity: event MaxGasPriceUpdated(uint128 newMaxGasPrice)
func (_Reserve *ReserveFilterer) ParseMaxGasPriceUpdated(log types.Log) (*ReserveMaxGasPriceUpdated, error) {
	event := new(ReserveMaxGasPriceUpdated)
	if err := _Reserve.contract.UnpackLog(event, "MaxGasPriceUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ReserveNewTokenWalletIterator is returned from FilterNewTokenWallet and is used to iterate over the raw logs and unpacked data for NewTokenWallet events raised by the Reserve contract.
type ReserveNewTokenWalletIterator struct {
	Event *ReserveNewTokenWallet // Event containing the contract specifics and raw log

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
func (it *ReserveNewTokenWalletIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveNewTokenWallet)
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
		it.Event = new(ReserveNewTokenWallet)
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
func (it *ReserveNewTokenWalletIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveNewTokenWalletIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveNewTokenWallet represents a NewTokenWallet event raised by the Reserve contract.
type ReserveNewTokenWallet struct {
	Token  common.Address
	Wallet common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNewTokenWallet is a free log retrieval operation binding the contract event 0x81995c7b922889ac0a81e41866106d4046268ea3a9abaae9f9e080a6ce36ee7d.
//
// Solidity: event NewTokenWallet(address indexed token, address indexed wallet)
func (_Reserve *ReserveFilterer) FilterNewTokenWallet(opts *bind.FilterOpts, token []common.Address, wallet []common.Address) (*ReserveNewTokenWalletIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "NewTokenWallet", tokenRule, walletRule)
	if err != nil {
		return nil, err
	}
	return &ReserveNewTokenWalletIterator{contract: _Reserve.contract, event: "NewTokenWallet", logs: logs, sub: sub}, nil
}

// WatchNewTokenWallet is a free log subscription operation binding the contract event 0x81995c7b922889ac0a81e41866106d4046268ea3a9abaae9f9e080a6ce36ee7d.
//
// Solidity: event NewTokenWallet(address indexed token, address indexed wallet)
func (_Reserve *ReserveFilterer) WatchNewTokenWallet(opts *bind.WatchOpts, sink chan<- *ReserveNewTokenWallet, token []common.Address, wallet []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "NewTokenWallet", tokenRule, walletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveNewTokenWallet)
				if err := _Reserve.contract.UnpackLog(event, "NewTokenWallet", log); err != nil {
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

// ParseNewTokenWallet is a log parse operation binding the contract event 0x81995c7b922889ac0a81e41866106d4046268ea3a9abaae9f9e080a6ce36ee7d.
//
// Solidity: event NewTokenWallet(address indexed token, address indexed wallet)
func (_Reserve *ReserveFilterer) ParseNewTokenWallet(log types.Log) (*ReserveNewTokenWallet, error) {
	event := new(ReserveNewTokenWallet)
	if err := _Reserve.contract.UnpackLog(event, "NewTokenWallet", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_Reserve *ReserveFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*ReserveOperatorAddedIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &ReserveOperatorAddedIterator{contract: _Reserve.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
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

// ParseOperatorAdded is a log parse operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_Reserve *ReserveFilterer) ParseOperatorAdded(log types.Log) (*ReserveOperatorAdded, error) {
	event := new(ReserveOperatorAdded)
	if err := _Reserve.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ReserveSetConversionRateAddressIterator is returned from FilterSetConversionRateAddress and is used to iterate over the raw logs and unpacked data for SetConversionRateAddress events raised by the Reserve contract.
type ReserveSetConversionRateAddressIterator struct {
	Event *ReserveSetConversionRateAddress // Event containing the contract specifics and raw log

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
func (it *ReserveSetConversionRateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveSetConversionRateAddress)
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
		it.Event = new(ReserveSetConversionRateAddress)
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
func (it *ReserveSetConversionRateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveSetConversionRateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveSetConversionRateAddress represents a SetConversionRateAddress event raised by the Reserve contract.
type ReserveSetConversionRateAddress struct {
	Rate common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetConversionRateAddress is a free log retrieval operation binding the contract event 0x333c220e52469bea5ce17b670353e10868f0c96768325592128e44d06e5b99cc.
//
// Solidity: event SetConversionRateAddress(address indexed rate)
func (_Reserve *ReserveFilterer) FilterSetConversionRateAddress(opts *bind.FilterOpts, rate []common.Address) (*ReserveSetConversionRateAddressIterator, error) {

	var rateRule []interface{}
	for _, rateItem := range rate {
		rateRule = append(rateRule, rateItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "SetConversionRateAddress", rateRule)
	if err != nil {
		return nil, err
	}
	return &ReserveSetConversionRateAddressIterator{contract: _Reserve.contract, event: "SetConversionRateAddress", logs: logs, sub: sub}, nil
}

// WatchSetConversionRateAddress is a free log subscription operation binding the contract event 0x333c220e52469bea5ce17b670353e10868f0c96768325592128e44d06e5b99cc.
//
// Solidity: event SetConversionRateAddress(address indexed rate)
func (_Reserve *ReserveFilterer) WatchSetConversionRateAddress(opts *bind.WatchOpts, sink chan<- *ReserveSetConversionRateAddress, rate []common.Address) (event.Subscription, error) {

	var rateRule []interface{}
	for _, rateItem := range rate {
		rateRule = append(rateRule, rateItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "SetConversionRateAddress", rateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveSetConversionRateAddress)
				if err := _Reserve.contract.UnpackLog(event, "SetConversionRateAddress", log); err != nil {
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

// ParseSetConversionRateAddress is a log parse operation binding the contract event 0x333c220e52469bea5ce17b670353e10868f0c96768325592128e44d06e5b99cc.
//
// Solidity: event SetConversionRateAddress(address indexed rate)
func (_Reserve *ReserveFilterer) ParseSetConversionRateAddress(log types.Log) (*ReserveSetConversionRateAddress, error) {
	event := new(ReserveSetConversionRateAddress)
	if err := _Reserve.contract.UnpackLog(event, "SetConversionRateAddress", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ReserveSetSanityRateAddressIterator is returned from FilterSetSanityRateAddress and is used to iterate over the raw logs and unpacked data for SetSanityRateAddress events raised by the Reserve contract.
type ReserveSetSanityRateAddressIterator struct {
	Event *ReserveSetSanityRateAddress // Event containing the contract specifics and raw log

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
func (it *ReserveSetSanityRateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveSetSanityRateAddress)
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
		it.Event = new(ReserveSetSanityRateAddress)
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
func (it *ReserveSetSanityRateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveSetSanityRateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveSetSanityRateAddress represents a SetSanityRateAddress event raised by the Reserve contract.
type ReserveSetSanityRateAddress struct {
	Sanity common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSetSanityRateAddress is a free log retrieval operation binding the contract event 0xaa5a0552d43d8d9e7d64c0286bcfceeccb1875e7320c1ac63f71e9894f177166.
//
// Solidity: event SetSanityRateAddress(address indexed sanity)
func (_Reserve *ReserveFilterer) FilterSetSanityRateAddress(opts *bind.FilterOpts, sanity []common.Address) (*ReserveSetSanityRateAddressIterator, error) {

	var sanityRule []interface{}
	for _, sanityItem := range sanity {
		sanityRule = append(sanityRule, sanityItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "SetSanityRateAddress", sanityRule)
	if err != nil {
		return nil, err
	}
	return &ReserveSetSanityRateAddressIterator{contract: _Reserve.contract, event: "SetSanityRateAddress", logs: logs, sub: sub}, nil
}

// WatchSetSanityRateAddress is a free log subscription operation binding the contract event 0xaa5a0552d43d8d9e7d64c0286bcfceeccb1875e7320c1ac63f71e9894f177166.
//
// Solidity: event SetSanityRateAddress(address indexed sanity)
func (_Reserve *ReserveFilterer) WatchSetSanityRateAddress(opts *bind.WatchOpts, sink chan<- *ReserveSetSanityRateAddress, sanity []common.Address) (event.Subscription, error) {

	var sanityRule []interface{}
	for _, sanityItem := range sanity {
		sanityRule = append(sanityRule, sanityItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "SetSanityRateAddress", sanityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveSetSanityRateAddress)
				if err := _Reserve.contract.UnpackLog(event, "SetSanityRateAddress", log); err != nil {
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

// ParseSetSanityRateAddress is a log parse operation binding the contract event 0xaa5a0552d43d8d9e7d64c0286bcfceeccb1875e7320c1ac63f71e9894f177166.
//
// Solidity: event SetSanityRateAddress(address indexed sanity)
func (_Reserve *ReserveFilterer) ParseSetSanityRateAddress(log types.Log) (*ReserveSetSanityRateAddress, error) {
	event := new(ReserveSetSanityRateAddress)
	if err := _Reserve.contract.UnpackLog(event, "SetSanityRateAddress", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ReserveSetWethAddressIterator is returned from FilterSetWethAddress and is used to iterate over the raw logs and unpacked data for SetWethAddress events raised by the Reserve contract.
type ReserveSetWethAddressIterator struct {
	Event *ReserveSetWethAddress // Event containing the contract specifics and raw log

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
func (it *ReserveSetWethAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveSetWethAddress)
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
		it.Event = new(ReserveSetWethAddress)
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
func (it *ReserveSetWethAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveSetWethAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveSetWethAddress represents a SetWethAddress event raised by the Reserve contract.
type ReserveSetWethAddress struct {
	Weth common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetWethAddress is a free log retrieval operation binding the contract event 0xff8ab24f675c1eee431de04d5ba93b5d6e1e947359771788c5db3968d10c2e3e.
//
// Solidity: event SetWethAddress(address indexed weth)
func (_Reserve *ReserveFilterer) FilterSetWethAddress(opts *bind.FilterOpts, weth []common.Address) (*ReserveSetWethAddressIterator, error) {

	var wethRule []interface{}
	for _, wethItem := range weth {
		wethRule = append(wethRule, wethItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "SetWethAddress", wethRule)
	if err != nil {
		return nil, err
	}
	return &ReserveSetWethAddressIterator{contract: _Reserve.contract, event: "SetWethAddress", logs: logs, sub: sub}, nil
}

// WatchSetWethAddress is a free log subscription operation binding the contract event 0xff8ab24f675c1eee431de04d5ba93b5d6e1e947359771788c5db3968d10c2e3e.
//
// Solidity: event SetWethAddress(address indexed weth)
func (_Reserve *ReserveFilterer) WatchSetWethAddress(opts *bind.WatchOpts, sink chan<- *ReserveSetWethAddress, weth []common.Address) (event.Subscription, error) {

	var wethRule []interface{}
	for _, wethItem := range weth {
		wethRule = append(wethRule, wethItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "SetWethAddress", wethRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveSetWethAddress)
				if err := _Reserve.contract.UnpackLog(event, "SetWethAddress", log); err != nil {
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

// ParseSetWethAddress is a log parse operation binding the contract event 0xff8ab24f675c1eee431de04d5ba93b5d6e1e947359771788c5db3968d10c2e3e.
//
// Solidity: event SetWethAddress(address indexed weth)
func (_Reserve *ReserveFilterer) ParseSetWethAddress(log types.Log) (*ReserveSetWethAddress, error) {
	event := new(ReserveSetWethAddress)
	if err := _Reserve.contract.UnpackLog(event, "SetWethAddress", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_Reserve *ReserveFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*ReserveTokenWithdrawIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &ReserveTokenWithdrawIterator{contract: _Reserve.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
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

// ParseTokenWithdraw is a log parse operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_Reserve *ReserveFilterer) ParseTokenWithdraw(log types.Log) (*ReserveTokenWithdraw, error) {
	event := new(ReserveTokenWithdraw)
	if err := _Reserve.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event TradeEnabled(bool enable)
func (_Reserve *ReserveFilterer) FilterTradeEnabled(opts *bind.FilterOpts) (*ReserveTradeEnabledIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TradeEnabled")
	if err != nil {
		return nil, err
	}
	return &ReserveTradeEnabledIterator{contract: _Reserve.contract, event: "TradeEnabled", logs: logs, sub: sub}, nil
}

// WatchTradeEnabled is a free log subscription operation binding the contract event 0x7d7f00509dd73ac4449f698ae75ccc797895eff5fa9d446d3df387598a26e735.
//
// Solidity: event TradeEnabled(bool enable)
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

// ParseTradeEnabled is a log parse operation binding the contract event 0x7d7f00509dd73ac4449f698ae75ccc797895eff5fa9d446d3df387598a26e735.
//
// Solidity: event TradeEnabled(bool enable)
func (_Reserve *ReserveFilterer) ParseTradeEnabled(log types.Log) (*ReserveTradeEnabled, error) {
	event := new(ReserveTradeEnabled)
	if err := _Reserve.contract.UnpackLog(event, "TradeEnabled", log); err != nil {
		return nil, err
	}
	return event, nil
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
	Trader      common.Address
	Src         common.Address
	SrcAmount   *big.Int
	DestToken   common.Address
	DestAmount  *big.Int
	DestAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTradeExecute is a free log retrieval operation binding the contract event 0x4ee2afc3e9f9e97f558641bdc31ff31e4f34a1aaa2390cffbd64ee9ac18dfbec.
//
// Solidity: event TradeExecute(address origin, address indexed trader, address indexed src, uint256 srcAmount, address indexed destToken, uint256 destAmount, address destAddress)
func (_Reserve *ReserveFilterer) FilterTradeExecute(opts *bind.FilterOpts, trader []common.Address, src []common.Address, destToken []common.Address) (*ReserveTradeExecuteIterator, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}
	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	var destTokenRule []interface{}
	for _, destTokenItem := range destToken {
		destTokenRule = append(destTokenRule, destTokenItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TradeExecute", traderRule, srcRule, destTokenRule)
	if err != nil {
		return nil, err
	}
	return &ReserveTradeExecuteIterator{contract: _Reserve.contract, event: "TradeExecute", logs: logs, sub: sub}, nil
}

// WatchTradeExecute is a free log subscription operation binding the contract event 0x4ee2afc3e9f9e97f558641bdc31ff31e4f34a1aaa2390cffbd64ee9ac18dfbec.
//
// Solidity: event TradeExecute(address origin, address indexed trader, address indexed src, uint256 srcAmount, address indexed destToken, uint256 destAmount, address destAddress)
func (_Reserve *ReserveFilterer) WatchTradeExecute(opts *bind.WatchOpts, sink chan<- *ReserveTradeExecute, trader []common.Address, src []common.Address, destToken []common.Address) (event.Subscription, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}
	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	var destTokenRule []interface{}
	for _, destTokenItem := range destToken {
		destTokenRule = append(destTokenRule, destTokenItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "TradeExecute", traderRule, srcRule, destTokenRule)
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

// ParseTradeExecute is a log parse operation binding the contract event 0x4ee2afc3e9f9e97f558641bdc31ff31e4f34a1aaa2390cffbd64ee9ac18dfbec.
//
// Solidity: event TradeExecute(address origin, address indexed trader, address indexed src, uint256 srcAmount, address indexed destToken, uint256 destAmount, address destAddress)
func (_Reserve *ReserveFilterer) ParseTradeExecute(log types.Log) (*ReserveTradeExecute, error) {
	event := new(ReserveTradeExecute)
	if err := _Reserve.contract.UnpackLog(event, "TradeExecute", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_Reserve *ReserveFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*ReserveTransferAdminPendingIterator, error) {

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &ReserveTransferAdminPendingIterator{contract: _Reserve.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
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

// ParseTransferAdminPending is a log parse operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_Reserve *ReserveFilterer) ParseTransferAdminPending(log types.Log) (*ReserveTransferAdminPending, error) {
	event := new(ReserveTransferAdminPending)
	if err := _Reserve.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event WithdrawAddressApproved(address indexed token, address indexed addr, bool approve)
func (_Reserve *ReserveFilterer) FilterWithdrawAddressApproved(opts *bind.FilterOpts, token []common.Address, addr []common.Address) (*ReserveWithdrawAddressApprovedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "WithdrawAddressApproved", tokenRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &ReserveWithdrawAddressApprovedIterator{contract: _Reserve.contract, event: "WithdrawAddressApproved", logs: logs, sub: sub}, nil
}

// WatchWithdrawAddressApproved is a free log subscription operation binding the contract event 0xd5fd5351efae1f4bb760079da9f0ff9589e2c3e216337ca9d39cdff573b245c4.
//
// Solidity: event WithdrawAddressApproved(address indexed token, address indexed addr, bool approve)
func (_Reserve *ReserveFilterer) WatchWithdrawAddressApproved(opts *bind.WatchOpts, sink chan<- *ReserveWithdrawAddressApproved, token []common.Address, addr []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "WithdrawAddressApproved", tokenRule, addrRule)
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

// ParseWithdrawAddressApproved is a log parse operation binding the contract event 0xd5fd5351efae1f4bb760079da9f0ff9589e2c3e216337ca9d39cdff573b245c4.
//
// Solidity: event WithdrawAddressApproved(address indexed token, address indexed addr, bool approve)
func (_Reserve *ReserveFilterer) ParseWithdrawAddressApproved(log types.Log) (*ReserveWithdrawAddressApproved, error) {
	event := new(ReserveWithdrawAddressApproved)
	if err := _Reserve.contract.UnpackLog(event, "WithdrawAddressApproved", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event WithdrawFunds(address indexed token, uint256 amount, address indexed destination)
func (_Reserve *ReserveFilterer) FilterWithdrawFunds(opts *bind.FilterOpts, token []common.Address, destination []common.Address) (*ReserveWithdrawFundsIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Reserve.contract.FilterLogs(opts, "WithdrawFunds", tokenRule, destinationRule)
	if err != nil {
		return nil, err
	}
	return &ReserveWithdrawFundsIterator{contract: _Reserve.contract, event: "WithdrawFunds", logs: logs, sub: sub}, nil
}

// WatchWithdrawFunds is a free log subscription operation binding the contract event 0xb67719fc33c1f17d31bf3a698690d62066b1e0bae28fcd3c56cf2c015c2863d6.
//
// Solidity: event WithdrawFunds(address indexed token, uint256 amount, address indexed destination)
func (_Reserve *ReserveFilterer) WatchWithdrawFunds(opts *bind.WatchOpts, sink chan<- *ReserveWithdrawFunds, token []common.Address, destination []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Reserve.contract.WatchLogs(opts, "WithdrawFunds", tokenRule, destinationRule)
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

// ParseWithdrawFunds is a log parse operation binding the contract event 0xb67719fc33c1f17d31bf3a698690d62066b1e0bae28fcd3c56cf2c015c2863d6.
//
// Solidity: event WithdrawFunds(address indexed token, uint256 amount, address indexed destination)
func (_Reserve *ReserveFilterer) ParseWithdrawFunds(log types.Log) (*ReserveWithdrawFunds, error) {
	event := new(ReserveWithdrawFunds)
	if err := _Reserve.contract.UnpackLog(event, "WithdrawFunds", log); err != nil {
		return nil, err
	}
	return event, nil
}
