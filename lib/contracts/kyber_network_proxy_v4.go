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

// KyberNetworkProxyV4ABI is the input ABI used to generate the binding from.
const KyberNetworkProxyV4ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualSrcAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualDestAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"platformFeeBps\",\"type\":\"uint256\"}],\"name\":\"ExecuteTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberHint\",\"name\":\"kyberHintHandler\",\"type\":\"address\"}],\"name\":\"KyberHintHandlerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberNetwork\",\"name\":\"newKyberNetwork\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIKyberNetwork\",\"name\":\"previousKyberNetwork\",\"type\":\"address\"}],\"name\":\"KyberNetworkSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcQty\",\"type\":\"uint256\"}],\"name\":\"getExpectedRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"worstRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcQty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"platformFeeBps\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"getExpectedRateAfterFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberHintHandler\",\"outputs\":[{\"internalType\":\"contractIKyberHint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberNetwork\",\"outputs\":[{\"internalType\":\"contractIKyberNetwork\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberHint\",\"name\":\"_kyberHintHandler\",\"type\":\"address\"}],\"name\":\"setHintHandler\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberNetwork\",\"name\":\"_kyberNetwork\",\"type\":\"address\"}],\"name\":\"setKyberNetwork\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapEtherToToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapTokenToEther\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"}],\"name\":\"swapTokenToToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"platformWallet\",\"type\":\"address\"}],\"name\":\"trade\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"walletId\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"platformFeeBps\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHintAndFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"destAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KyberNetworkProxyV4 is an auto generated Go binding around an Ethereum contract.
type KyberNetworkProxyV4 struct {
	KyberNetworkProxyV4Caller     // Read-only binding to the contract
	KyberNetworkProxyV4Transactor // Write-only binding to the contract
	KyberNetworkProxyV4Filterer   // Log filterer for contract events
}

// KyberNetworkProxyV4Caller is an auto generated read-only Go binding around an Ethereum contract.
type KyberNetworkProxyV4Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkProxyV4Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KyberNetworkProxyV4Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkProxyV4Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KyberNetworkProxyV4Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkProxyV4Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KyberNetworkProxyV4Session struct {
	Contract     *KyberNetworkProxyV4 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// KyberNetworkProxyV4CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KyberNetworkProxyV4CallerSession struct {
	Contract *KyberNetworkProxyV4Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// KyberNetworkProxyV4TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KyberNetworkProxyV4TransactorSession struct {
	Contract     *KyberNetworkProxyV4Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// KyberNetworkProxyV4Raw is an auto generated low-level Go binding around an Ethereum contract.
type KyberNetworkProxyV4Raw struct {
	Contract *KyberNetworkProxyV4 // Generic contract binding to access the raw methods on
}

// KyberNetworkProxyV4CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KyberNetworkProxyV4CallerRaw struct {
	Contract *KyberNetworkProxyV4Caller // Generic read-only contract binding to access the raw methods on
}

// KyberNetworkProxyV4TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KyberNetworkProxyV4TransactorRaw struct {
	Contract *KyberNetworkProxyV4Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKyberNetworkProxyV4 creates a new instance of KyberNetworkProxyV4, bound to a specific deployed contract.
func NewKyberNetworkProxyV4(address common.Address, backend bind.ContractBackend) (*KyberNetworkProxyV4, error) {
	contract, err := bindKyberNetworkProxyV4(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4{KyberNetworkProxyV4Caller: KyberNetworkProxyV4Caller{contract: contract}, KyberNetworkProxyV4Transactor: KyberNetworkProxyV4Transactor{contract: contract}, KyberNetworkProxyV4Filterer: KyberNetworkProxyV4Filterer{contract: contract}}, nil
}

// NewKyberNetworkProxyV4Caller creates a new read-only instance of KyberNetworkProxyV4, bound to a specific deployed contract.
func NewKyberNetworkProxyV4Caller(address common.Address, caller bind.ContractCaller) (*KyberNetworkProxyV4Caller, error) {
	contract, err := bindKyberNetworkProxyV4(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4Caller{contract: contract}, nil
}

// NewKyberNetworkProxyV4Transactor creates a new write-only instance of KyberNetworkProxyV4, bound to a specific deployed contract.
func NewKyberNetworkProxyV4Transactor(address common.Address, transactor bind.ContractTransactor) (*KyberNetworkProxyV4Transactor, error) {
	contract, err := bindKyberNetworkProxyV4(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4Transactor{contract: contract}, nil
}

// NewKyberNetworkProxyV4Filterer creates a new log filterer instance of KyberNetworkProxyV4, bound to a specific deployed contract.
func NewKyberNetworkProxyV4Filterer(address common.Address, filterer bind.ContractFilterer) (*KyberNetworkProxyV4Filterer, error) {
	contract, err := bindKyberNetworkProxyV4(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4Filterer{contract: contract}, nil
}

// bindKyberNetworkProxyV4 binds a generic wrapper to an already deployed contract.
func bindKyberNetworkProxyV4(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KyberNetworkProxyV4ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberNetworkProxyV4.Contract.KyberNetworkProxyV4Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.KyberNetworkProxyV4Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.KyberNetworkProxyV4Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberNetworkProxyV4.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) Admin() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.Admin(&_KyberNetworkProxyV4.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) Admin() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.Admin(&_KyberNetworkProxyV4.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) Enabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) Enabled() (bool, error) {
	return _KyberNetworkProxyV4.Contract.Enabled(&_KyberNetworkProxyV4.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) Enabled() (bool, error) {
	return _KyberNetworkProxyV4.Contract.Enabled(&_KyberNetworkProxyV4.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) GetAlerters() ([]common.Address, error) {
	return _KyberNetworkProxyV4.Contract.GetAlerters(&_KyberNetworkProxyV4.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) GetAlerters() ([]common.Address, error) {
	return _KyberNetworkProxyV4.Contract.GetAlerters(&_KyberNetworkProxyV4.CallOpts)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) GetExpectedRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	ret := new(struct {
		ExpectedRate *big.Int
		WorstRate    *big.Int
	})
	out := ret
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "getExpectedRate", src, dest, srcQty)
	return *ret, err
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	return _KyberNetworkProxyV4.Contract.GetExpectedRate(&_KyberNetworkProxyV4.CallOpts, src, dest, srcQty)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	return _KyberNetworkProxyV4.Contract.GetExpectedRate(&_KyberNetworkProxyV4.CallOpts, src, dest, srcQty)
}

// GetExpectedRateAfterFee is a free data retrieval call binding the contract method 0x418436bc.
//
// Solidity: function getExpectedRateAfterFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 expectedRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) GetExpectedRateAfterFee(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "getExpectedRateAfterFee", src, dest, srcQty, platformFeeBps, hint)
	return *ret0, err
}

// GetExpectedRateAfterFee is a free data retrieval call binding the contract method 0x418436bc.
//
// Solidity: function getExpectedRateAfterFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 expectedRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) GetExpectedRateAfterFee(src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (*big.Int, error) {
	return _KyberNetworkProxyV4.Contract.GetExpectedRateAfterFee(&_KyberNetworkProxyV4.CallOpts, src, dest, srcQty, platformFeeBps, hint)
}

// GetExpectedRateAfterFee is a free data retrieval call binding the contract method 0x418436bc.
//
// Solidity: function getExpectedRateAfterFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 expectedRate)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) GetExpectedRateAfterFee(src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (*big.Int, error) {
	return _KyberNetworkProxyV4.Contract.GetExpectedRateAfterFee(&_KyberNetworkProxyV4.CallOpts, src, dest, srcQty, platformFeeBps, hint)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) GetOperators() ([]common.Address, error) {
	return _KyberNetworkProxyV4.Contract.GetOperators(&_KyberNetworkProxyV4.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) GetOperators() ([]common.Address, error) {
	return _KyberNetworkProxyV4.Contract.GetOperators(&_KyberNetworkProxyV4.CallOpts)
}

// KyberHintHandler is a free data retrieval call binding the contract method 0x13c213b7.
//
// Solidity: function kyberHintHandler() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) KyberHintHandler(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "kyberHintHandler")
	return *ret0, err
}

// KyberHintHandler is a free data retrieval call binding the contract method 0x13c213b7.
//
// Solidity: function kyberHintHandler() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) KyberHintHandler() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.KyberHintHandler(&_KyberNetworkProxyV4.CallOpts)
}

// KyberHintHandler is a free data retrieval call binding the contract method 0x13c213b7.
//
// Solidity: function kyberHintHandler() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) KyberHintHandler() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.KyberHintHandler(&_KyberNetworkProxyV4.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) KyberNetwork(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "kyberNetwork")
	return *ret0, err
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) KyberNetwork() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.KyberNetwork(&_KyberNetworkProxyV4.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) KyberNetwork() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.KyberNetwork(&_KyberNetworkProxyV4.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) MaxGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "maxGasPrice")
	return *ret0, err
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) MaxGasPrice() (*big.Int, error) {
	return _KyberNetworkProxyV4.Contract.MaxGasPrice(&_KyberNetworkProxyV4.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) MaxGasPrice() (*big.Int, error) {
	return _KyberNetworkProxyV4.Contract.MaxGasPrice(&_KyberNetworkProxyV4.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Caller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetworkProxyV4.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) PendingAdmin() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.PendingAdmin(&_KyberNetworkProxyV4.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4CallerSession) PendingAdmin() (common.Address, error) {
	return _KyberNetworkProxyV4.Contract.PendingAdmin(&_KyberNetworkProxyV4.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.AddAlerter(&_KyberNetworkProxyV4.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.AddAlerter(&_KyberNetworkProxyV4.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.AddOperator(&_KyberNetworkProxyV4.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.AddOperator(&_KyberNetworkProxyV4.TransactOpts, newOperator)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) ClaimAdmin() (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.ClaimAdmin(&_KyberNetworkProxyV4.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.ClaimAdmin(&_KyberNetworkProxyV4.TransactOpts)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.RemoveAlerter(&_KyberNetworkProxyV4.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.RemoveAlerter(&_KyberNetworkProxyV4.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.RemoveOperator(&_KyberNetworkProxyV4.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.RemoveOperator(&_KyberNetworkProxyV4.TransactOpts, operator)
}

// SetHintHandler is a paid mutator transaction binding the contract method 0xb85d950f.
//
// Solidity: function setHintHandler(address _kyberHintHandler) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) SetHintHandler(opts *bind.TransactOpts, _kyberHintHandler common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "setHintHandler", _kyberHintHandler)
}

// SetHintHandler is a paid mutator transaction binding the contract method 0xb85d950f.
//
// Solidity: function setHintHandler(address _kyberHintHandler) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) SetHintHandler(_kyberHintHandler common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SetHintHandler(&_KyberNetworkProxyV4.TransactOpts, _kyberHintHandler)
}

// SetHintHandler is a paid mutator transaction binding the contract method 0xb85d950f.
//
// Solidity: function setHintHandler(address _kyberHintHandler) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) SetHintHandler(_kyberHintHandler common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SetHintHandler(&_KyberNetworkProxyV4.TransactOpts, _kyberHintHandler)
}

// SetKyberNetwork is a paid mutator transaction binding the contract method 0x54a325a6.
//
// Solidity: function setKyberNetwork(address _kyberNetwork) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) SetKyberNetwork(opts *bind.TransactOpts, _kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "setKyberNetwork", _kyberNetwork)
}

// SetKyberNetwork is a paid mutator transaction binding the contract method 0x54a325a6.
//
// Solidity: function setKyberNetwork(address _kyberNetwork) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) SetKyberNetwork(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SetKyberNetwork(&_KyberNetworkProxyV4.TransactOpts, _kyberNetwork)
}

// SetKyberNetwork is a paid mutator transaction binding the contract method 0x54a325a6.
//
// Solidity: function setKyberNetwork(address _kyberNetwork) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) SetKyberNetwork(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SetKyberNetwork(&_KyberNetworkProxyV4.TransactOpts, _kyberNetwork)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) SwapEtherToToken(opts *bind.TransactOpts, token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "swapEtherToToken", token, minConversionRate)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) SwapEtherToToken(token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapEtherToToken(&_KyberNetworkProxyV4.TransactOpts, token, minConversionRate)
}

// SwapEtherToToken is a paid mutator transaction binding the contract method 0x7a2a0456.
//
// Solidity: function swapEtherToToken(address token, uint256 minConversionRate) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) SwapEtherToToken(token common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapEtherToToken(&_KyberNetworkProxyV4.TransactOpts, token, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) SwapTokenToEther(opts *bind.TransactOpts, token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "swapTokenToEther", token, srcAmount, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) SwapTokenToEther(token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapTokenToEther(&_KyberNetworkProxyV4.TransactOpts, token, srcAmount, minConversionRate)
}

// SwapTokenToEther is a paid mutator transaction binding the contract method 0x3bba21dc.
//
// Solidity: function swapTokenToEther(address token, uint256 srcAmount, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) SwapTokenToEther(token common.Address, srcAmount *big.Int, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapTokenToEther(&_KyberNetworkProxyV4.TransactOpts, token, srcAmount, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) SwapTokenToToken(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "swapTokenToToken", src, srcAmount, dest, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) SwapTokenToToken(src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapTokenToToken(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, minConversionRate)
}

// SwapTokenToToken is a paid mutator transaction binding the contract method 0x7409e2eb.
//
// Solidity: function swapTokenToToken(address src, uint256 srcAmount, address dest, uint256 minConversionRate) returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) SwapTokenToToken(src common.Address, srcAmount *big.Int, dest common.Address, minConversionRate *big.Int) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.SwapTokenToToken(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, minConversionRate)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) Trade(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "trade", src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) Trade(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.Trade(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet)
}

// Trade is a paid mutator transaction binding the contract method 0xcb3c28c7.
//
// Solidity: function trade(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) Trade(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.Trade(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) TradeWithHint(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "tradeWithHint", src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) TradeWithHint(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TradeWithHint(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x29589f61.
//
// Solidity: function tradeWithHint(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) TradeWithHint(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TradeWithHint(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xae591d54.
//
// Solidity: function tradeWithHintAndFee(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) TradeWithHintAndFee(opts *bind.TransactOpts, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "tradeWithHintAndFee", src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xae591d54.
//
// Solidity: function tradeWithHintAndFee(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) TradeWithHintAndFee(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TradeWithHintAndFee(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xae591d54.
//
// Solidity: function tradeWithHintAndFee(address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) TradeWithHintAndFee(src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TradeWithHintAndFee(&_KyberNetworkProxyV4.TransactOpts, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TransferAdmin(&_KyberNetworkProxyV4.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TransferAdmin(&_KyberNetworkProxyV4.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TransferAdminQuickly(&_KyberNetworkProxyV4.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.TransferAdminQuickly(&_KyberNetworkProxyV4.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.WithdrawEther(&_KyberNetworkProxyV4.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.WithdrawEther(&_KyberNetworkProxyV4.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Transactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Session) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.WithdrawToken(&_KyberNetworkProxyV4.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4TransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetworkProxyV4.Contract.WithdrawToken(&_KyberNetworkProxyV4.TransactOpts, token, amount, sendTo)
}

// KyberNetworkProxyV4AdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4AdminClaimedIterator struct {
	Event *KyberNetworkProxyV4AdminClaimed // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4AdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4AdminClaimed)
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
		it.Event = new(KyberNetworkProxyV4AdminClaimed)
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
func (it *KyberNetworkProxyV4AdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4AdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4AdminClaimed represents a AdminClaimed event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4AdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterAdminClaimed(opts *bind.FilterOpts) (*KyberNetworkProxyV4AdminClaimedIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4AdminClaimedIterator{contract: _KyberNetworkProxyV4.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4AdminClaimed) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4AdminClaimed)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseAdminClaimed(log types.Log) (*KyberNetworkProxyV4AdminClaimed, error) {
	event := new(KyberNetworkProxyV4AdminClaimed)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4AlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4AlerterAddedIterator struct {
	Event *KyberNetworkProxyV4AlerterAdded // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4AlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4AlerterAdded)
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
		it.Event = new(KyberNetworkProxyV4AlerterAdded)
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
func (it *KyberNetworkProxyV4AlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4AlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4AlerterAdded represents a AlerterAdded event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4AlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterAlerterAdded(opts *bind.FilterOpts) (*KyberNetworkProxyV4AlerterAddedIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4AlerterAddedIterator{contract: _KyberNetworkProxyV4.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4AlerterAdded) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4AlerterAdded)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseAlerterAdded(log types.Log) (*KyberNetworkProxyV4AlerterAdded, error) {
	event := new(KyberNetworkProxyV4AlerterAdded)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4EtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4EtherWithdrawIterator struct {
	Event *KyberNetworkProxyV4EtherWithdraw // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4EtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4EtherWithdraw)
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
		it.Event = new(KyberNetworkProxyV4EtherWithdraw)
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
func (it *KyberNetworkProxyV4EtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4EtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4EtherWithdraw represents a EtherWithdraw event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4EtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*KyberNetworkProxyV4EtherWithdrawIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4EtherWithdrawIterator{contract: _KyberNetworkProxyV4.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4EtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4EtherWithdraw)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseEtherWithdraw(log types.Log) (*KyberNetworkProxyV4EtherWithdraw, error) {
	event := new(KyberNetworkProxyV4EtherWithdraw)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4ExecuteTradeIterator is returned from FilterExecuteTrade and is used to iterate over the raw logs and unpacked data for ExecuteTrade events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4ExecuteTradeIterator struct {
	Event *KyberNetworkProxyV4ExecuteTrade // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4ExecuteTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4ExecuteTrade)
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
		it.Event = new(KyberNetworkProxyV4ExecuteTrade)
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
func (it *KyberNetworkProxyV4ExecuteTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4ExecuteTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4ExecuteTrade represents a ExecuteTrade event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4ExecuteTrade struct {
	Trader           common.Address
	Src              common.Address
	Dest             common.Address
	DestAddress      common.Address
	ActualSrcAmount  *big.Int
	ActualDestAmount *big.Int
	PlatformWallet   common.Address
	PlatformFeeBps   *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterExecuteTrade is a free log retrieval operation binding the contract event 0xf724b4df6617473612b53d7f88ecc6ea983074b30960a049fcd0657ffe808083.
//
// Solidity: event ExecuteTrade(address indexed trader, address src, address dest, address destAddress, uint256 actualSrcAmount, uint256 actualDestAmount, address platformWallet, uint256 platformFeeBps)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterExecuteTrade(opts *bind.FilterOpts, trader []common.Address) (*KyberNetworkProxyV4ExecuteTradeIterator, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "ExecuteTrade", traderRule)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4ExecuteTradeIterator{contract: _KyberNetworkProxyV4.contract, event: "ExecuteTrade", logs: logs, sub: sub}, nil
}

// WatchExecuteTrade is a free log subscription operation binding the contract event 0xf724b4df6617473612b53d7f88ecc6ea983074b30960a049fcd0657ffe808083.
//
// Solidity: event ExecuteTrade(address indexed trader, address src, address dest, address destAddress, uint256 actualSrcAmount, uint256 actualDestAmount, address platformWallet, uint256 platformFeeBps)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchExecuteTrade(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4ExecuteTrade, trader []common.Address) (event.Subscription, error) {

	var traderRule []interface{}
	for _, traderItem := range trader {
		traderRule = append(traderRule, traderItem)
	}

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "ExecuteTrade", traderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4ExecuteTrade)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "ExecuteTrade", log); err != nil {
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

// ParseExecuteTrade is a log parse operation binding the contract event 0xf724b4df6617473612b53d7f88ecc6ea983074b30960a049fcd0657ffe808083.
//
// Solidity: event ExecuteTrade(address indexed trader, address src, address dest, address destAddress, uint256 actualSrcAmount, uint256 actualDestAmount, address platformWallet, uint256 platformFeeBps)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseExecuteTrade(log types.Log) (*KyberNetworkProxyV4ExecuteTrade, error) {
	event := new(KyberNetworkProxyV4ExecuteTrade)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "ExecuteTrade", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4KyberHintHandlerSetIterator is returned from FilterKyberHintHandlerSet and is used to iterate over the raw logs and unpacked data for KyberHintHandlerSet events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4KyberHintHandlerSetIterator struct {
	Event *KyberNetworkProxyV4KyberHintHandlerSet // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4KyberHintHandlerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4KyberHintHandlerSet)
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
		it.Event = new(KyberNetworkProxyV4KyberHintHandlerSet)
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
func (it *KyberNetworkProxyV4KyberHintHandlerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4KyberHintHandlerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4KyberHintHandlerSet represents a KyberHintHandlerSet event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4KyberHintHandlerSet struct {
	KyberHintHandler common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterKyberHintHandlerSet is a free log retrieval operation binding the contract event 0x6deb3a98fd141d661e9c0fb2d847541cc0c629cfb100c61011a76f57cb3b3a9b.
//
// Solidity: event KyberHintHandlerSet(address kyberHintHandler)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterKyberHintHandlerSet(opts *bind.FilterOpts) (*KyberNetworkProxyV4KyberHintHandlerSetIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "KyberHintHandlerSet")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4KyberHintHandlerSetIterator{contract: _KyberNetworkProxyV4.contract, event: "KyberHintHandlerSet", logs: logs, sub: sub}, nil
}

// WatchKyberHintHandlerSet is a free log subscription operation binding the contract event 0x6deb3a98fd141d661e9c0fb2d847541cc0c629cfb100c61011a76f57cb3b3a9b.
//
// Solidity: event KyberHintHandlerSet(address kyberHintHandler)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchKyberHintHandlerSet(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4KyberHintHandlerSet) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "KyberHintHandlerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4KyberHintHandlerSet)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "KyberHintHandlerSet", log); err != nil {
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

// ParseKyberHintHandlerSet is a log parse operation binding the contract event 0x6deb3a98fd141d661e9c0fb2d847541cc0c629cfb100c61011a76f57cb3b3a9b.
//
// Solidity: event KyberHintHandlerSet(address kyberHintHandler)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseKyberHintHandlerSet(log types.Log) (*KyberNetworkProxyV4KyberHintHandlerSet, error) {
	event := new(KyberNetworkProxyV4KyberHintHandlerSet)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "KyberHintHandlerSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4KyberNetworkSetIterator is returned from FilterKyberNetworkSet and is used to iterate over the raw logs and unpacked data for KyberNetworkSet events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4KyberNetworkSetIterator struct {
	Event *KyberNetworkProxyV4KyberNetworkSet // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4KyberNetworkSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4KyberNetworkSet)
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
		it.Event = new(KyberNetworkProxyV4KyberNetworkSet)
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
func (it *KyberNetworkProxyV4KyberNetworkSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4KyberNetworkSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4KyberNetworkSet represents a KyberNetworkSet event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4KyberNetworkSet struct {
	NewKyberNetwork      common.Address
	PreviousKyberNetwork common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkSet is a free log retrieval operation binding the contract event 0x8936e1f096bf0a8c9df862b3d1d5b82774cad78116200175f00b5b7ba3010b02.
//
// Solidity: event KyberNetworkSet(address newKyberNetwork, address previousKyberNetwork)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterKyberNetworkSet(opts *bind.FilterOpts) (*KyberNetworkProxyV4KyberNetworkSetIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "KyberNetworkSet")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4KyberNetworkSetIterator{contract: _KyberNetworkProxyV4.contract, event: "KyberNetworkSet", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkSet is a free log subscription operation binding the contract event 0x8936e1f096bf0a8c9df862b3d1d5b82774cad78116200175f00b5b7ba3010b02.
//
// Solidity: event KyberNetworkSet(address newKyberNetwork, address previousKyberNetwork)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchKyberNetworkSet(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4KyberNetworkSet) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "KyberNetworkSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4KyberNetworkSet)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "KyberNetworkSet", log); err != nil {
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

// ParseKyberNetworkSet is a log parse operation binding the contract event 0x8936e1f096bf0a8c9df862b3d1d5b82774cad78116200175f00b5b7ba3010b02.
//
// Solidity: event KyberNetworkSet(address newKyberNetwork, address previousKyberNetwork)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseKyberNetworkSet(log types.Log) (*KyberNetworkProxyV4KyberNetworkSet, error) {
	event := new(KyberNetworkProxyV4KyberNetworkSet)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "KyberNetworkSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4OperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4OperatorAddedIterator struct {
	Event *KyberNetworkProxyV4OperatorAdded // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4OperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4OperatorAdded)
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
		it.Event = new(KyberNetworkProxyV4OperatorAdded)
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
func (it *KyberNetworkProxyV4OperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4OperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4OperatorAdded represents a OperatorAdded event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4OperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterOperatorAdded(opts *bind.FilterOpts) (*KyberNetworkProxyV4OperatorAddedIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4OperatorAddedIterator{contract: _KyberNetworkProxyV4.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4OperatorAdded) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4OperatorAdded)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseOperatorAdded(log types.Log) (*KyberNetworkProxyV4OperatorAdded, error) {
	event := new(KyberNetworkProxyV4OperatorAdded)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4TokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4TokenWithdrawIterator struct {
	Event *KyberNetworkProxyV4TokenWithdraw // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4TokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4TokenWithdraw)
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
		it.Event = new(KyberNetworkProxyV4TokenWithdraw)
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
func (it *KyberNetworkProxyV4TokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4TokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4TokenWithdraw represents a TokenWithdraw event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4TokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*KyberNetworkProxyV4TokenWithdrawIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4TokenWithdrawIterator{contract: _KyberNetworkProxyV4.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4TokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4TokenWithdraw)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseTokenWithdraw(log types.Log) (*KyberNetworkProxyV4TokenWithdraw, error) {
	event := new(KyberNetworkProxyV4TokenWithdraw)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkProxyV4TransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4TransferAdminPendingIterator struct {
	Event *KyberNetworkProxyV4TransferAdminPending // Event containing the contract specifics and raw log

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
func (it *KyberNetworkProxyV4TransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkProxyV4TransferAdminPending)
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
		it.Event = new(KyberNetworkProxyV4TransferAdminPending)
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
func (it *KyberNetworkProxyV4TransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkProxyV4TransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkProxyV4TransferAdminPending represents a TransferAdminPending event raised by the KyberNetworkProxyV4 contract.
type KyberNetworkProxyV4TransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*KyberNetworkProxyV4TransferAdminPendingIterator, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkProxyV4TransferAdminPendingIterator{contract: _KyberNetworkProxyV4.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *KyberNetworkProxyV4TransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _KyberNetworkProxyV4.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkProxyV4TransferAdminPending)
				if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
func (_KyberNetworkProxyV4 *KyberNetworkProxyV4Filterer) ParseTransferAdminPending(log types.Log) (*KyberNetworkProxyV4TransferAdminPending, error) {
	event := new(KyberNetworkProxyV4TransferAdminPending)
	if err := _KyberNetworkProxyV4.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
		return nil, err
	}
	return event, nil
}
