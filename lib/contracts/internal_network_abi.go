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

// InternalNetworkABI is the input ABI used to generate the binding from.
const InternalNetworkABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"destAddress\",\"type\":\"address\"},{\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"name\":\"walletId\",\"type\":\"address\"},{\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHint\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"infoFields\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"feeBurner\",\"type\":\"address\"}],\"name\":\"setFeeBurner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"reservesPerTokenSrc\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"whiteList\",\"type\":\"address\"}],\"name\":\"setWhiteList\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxGasPrice\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"negligibleRateDiff\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeBurnerContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"expectedRate\",\"type\":\"address\"}],\"name\":\"setExpectedRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"expectedRateContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"whiteListContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"field\",\"type\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserCapInWei\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_enable\",\"type\":\"bool\"}],\"name\":\"setEnable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kyberNetworkProxyContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isReserve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcQty\",\"type\":\"uint256\"}],\"name\":\"getExpectedRate\",\"outputs\":[{\"name\":\"expectedRate\",\"type\":\"uint256\"},{\"name\":\"slippageRate\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"reserves\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getUserCapInTokenWei\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"reservesPerTokenDest\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\"},{\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"addReserve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"}],\"name\":\"searchBestRate\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxGasPriceValue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"field\",\"type\":\"bytes32\"}],\"name\":\"info\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"srcAmount\",\"type\":\"uint256\"}],\"name\":\"findBestRate\",\"outputs\":[{\"name\":\"obsolete\",\"type\":\"uint256\"},{\"name\":\"rate\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_maxGasPrice\",\"type\":\"uint256\"},{\"name\":\"_negligibleRateDiff\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"networkProxy\",\"type\":\"address\"}],\"name\":\"setKyberProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNumReserves\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"ethToToken\",\"type\":\"bool\"},{\"name\":\"tokenToEth\",\"type\":\"bool\"},{\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"listPairForReserve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EtherReceival\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"AddReserveToNetwork\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"ListReservePairs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"KyberProxySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"srcAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"srcToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"destAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"destToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"destAmount\",\"type\":\"uint256\"}],\"name\":\"KyberTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"}]"

// InternalNetwork is an auto generated Go binding around an Ethereum contract.
type InternalNetwork struct {
	InternalNetworkCaller     // Read-only binding to the contract
	InternalNetworkTransactor // Write-only binding to the contract
	InternalNetworkFilterer   // Log filterer for contract events
}

// InternalNetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type InternalNetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InternalNetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InternalNetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InternalNetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InternalNetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InternalNetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InternalNetworkSession struct {
	Contract     *InternalNetwork  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InternalNetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InternalNetworkCallerSession struct {
	Contract *InternalNetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// InternalNetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InternalNetworkTransactorSession struct {
	Contract     *InternalNetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// InternalNetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type InternalNetworkRaw struct {
	Contract *InternalNetwork // Generic contract binding to access the raw methods on
}

// InternalNetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InternalNetworkCallerRaw struct {
	Contract *InternalNetworkCaller // Generic read-only contract binding to access the raw methods on
}

// InternalNetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InternalNetworkTransactorRaw struct {
	Contract *InternalNetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInternalNetwork creates a new instance of InternalNetwork, bound to a specific deployed contract.
func NewInternalNetwork(address common.Address, backend bind.ContractBackend) (*InternalNetwork, error) {
	contract, err := bindInternalNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InternalNetwork{InternalNetworkCaller: InternalNetworkCaller{contract: contract}, InternalNetworkTransactor: InternalNetworkTransactor{contract: contract}, InternalNetworkFilterer: InternalNetworkFilterer{contract: contract}}, nil
}

// NewInternalNetworkCaller creates a new read-only instance of InternalNetwork, bound to a specific deployed contract.
func NewInternalNetworkCaller(address common.Address, caller bind.ContractCaller) (*InternalNetworkCaller, error) {
	contract, err := bindInternalNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InternalNetworkCaller{contract: contract}, nil
}

// NewInternalNetworkTransactor creates a new write-only instance of InternalNetwork, bound to a specific deployed contract.
func NewInternalNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*InternalNetworkTransactor, error) {
	contract, err := bindInternalNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InternalNetworkTransactor{contract: contract}, nil
}

// NewInternalNetworkFilterer creates a new log filterer instance of InternalNetwork, bound to a specific deployed contract.
func NewInternalNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*InternalNetworkFilterer, error) {
	contract, err := bindInternalNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InternalNetworkFilterer{contract: contract}, nil
}

// bindInternalNetwork binds a generic wrapper to an already deployed contract.
func bindInternalNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(InternalNetworkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InternalNetwork *InternalNetworkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _InternalNetwork.Contract.InternalNetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InternalNetwork *InternalNetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InternalNetwork.Contract.InternalNetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InternalNetwork *InternalNetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InternalNetwork.Contract.InternalNetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InternalNetwork *InternalNetworkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _InternalNetwork.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InternalNetwork *InternalNetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InternalNetwork.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InternalNetwork *InternalNetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InternalNetwork.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) Admin() (common.Address, error) {
	return _InternalNetwork.Contract.Admin(&_InternalNetwork.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) Admin() (common.Address, error) {
	return _InternalNetwork.Contract.Admin(&_InternalNetwork.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkCaller) Enabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkSession) Enabled() (bool, error) {
	return _InternalNetwork.Contract.Enabled(&_InternalNetwork.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkCallerSession) Enabled() (bool, error) {
	return _InternalNetwork.Contract.Enabled(&_InternalNetwork.CallOpts)
}

// ExpectedRateContract is a free data retrieval call binding the contract method 0x5dada964.
//
// Solidity: function expectedRateContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) ExpectedRateContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "expectedRateContract")
	return *ret0, err
}

// ExpectedRateContract is a free data retrieval call binding the contract method 0x5dada964.
//
// Solidity: function expectedRateContract() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) ExpectedRateContract() (common.Address, error) {
	return _InternalNetwork.Contract.ExpectedRateContract(&_InternalNetwork.CallOpts)
}

// ExpectedRateContract is a free data retrieval call binding the contract method 0x5dada964.
//
// Solidity: function expectedRateContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) ExpectedRateContract() (common.Address, error) {
	return _InternalNetwork.Contract.ExpectedRateContract(&_InternalNetwork.CallOpts)
}

// FeeBurnerContract is a free data retrieval call binding the contract method 0x579425b7.
//
// Solidity: function feeBurnerContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) FeeBurnerContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "feeBurnerContract")
	return *ret0, err
}

// FeeBurnerContract is a free data retrieval call binding the contract method 0x579425b7.
//
// Solidity: function feeBurnerContract() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) FeeBurnerContract() (common.Address, error) {
	return _InternalNetwork.Contract.FeeBurnerContract(&_InternalNetwork.CallOpts)
}

// FeeBurnerContract is a free data retrieval call binding the contract method 0x579425b7.
//
// Solidity: function feeBurnerContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) FeeBurnerContract() (common.Address, error) {
	return _InternalNetwork.Contract.FeeBurnerContract(&_InternalNetwork.CallOpts)
}

// FindBestRate is a free data retrieval call binding the contract method 0xb8388aca.
//
// Solidity: function findBestRate(address src, address dest, uint256 srcAmount) constant returns(uint256 obsolete, uint256 rate)
func (_InternalNetwork *InternalNetworkCaller) FindBestRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcAmount *big.Int) (struct {
	Obsolete *big.Int
	Rate     *big.Int
}, error) {
	ret := new(struct {
		Obsolete *big.Int
		Rate     *big.Int
	})
	out := ret
	err := _InternalNetwork.contract.Call(opts, out, "findBestRate", src, dest, srcAmount)
	return *ret, err
}

// FindBestRate is a free data retrieval call binding the contract method 0xb8388aca.
//
// Solidity: function findBestRate(address src, address dest, uint256 srcAmount) constant returns(uint256 obsolete, uint256 rate)
func (_InternalNetwork *InternalNetworkSession) FindBestRate(src common.Address, dest common.Address, srcAmount *big.Int) (struct {
	Obsolete *big.Int
	Rate     *big.Int
}, error) {
	return _InternalNetwork.Contract.FindBestRate(&_InternalNetwork.CallOpts, src, dest, srcAmount)
}

// FindBestRate is a free data retrieval call binding the contract method 0xb8388aca.
//
// Solidity: function findBestRate(address src, address dest, uint256 srcAmount) constant returns(uint256 obsolete, uint256 rate)
func (_InternalNetwork *InternalNetworkCallerSession) FindBestRate(src common.Address, dest common.Address, srcAmount *big.Int) (struct {
	Obsolete *big.Int
	Rate     *big.Int
}, error) {
	return _InternalNetwork.Contract.FindBestRate(&_InternalNetwork.CallOpts, src, dest, srcAmount)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_InternalNetwork *InternalNetworkCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_InternalNetwork *InternalNetworkSession) GetAlerters() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetAlerters(&_InternalNetwork.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() constant returns(address[])
func (_InternalNetwork *InternalNetworkCallerSession) GetAlerters() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetAlerters(&_InternalNetwork.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) GetBalance(opts *bind.CallOpts, token common.Address, user common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getBalance", token, user)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) GetBalance(token common.Address, user common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetBalance(&_InternalNetwork.CallOpts, token, user)
}

// GetBalance is a free data retrieval call binding the contract method 0xd4fac45d.
//
// Solidity: function getBalance(address token, address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) GetBalance(token common.Address, user common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetBalance(&_InternalNetwork.CallOpts, token, user)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_InternalNetwork *InternalNetworkCaller) GetExpectedRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	ret := new(struct {
		ExpectedRate *big.Int
		SlippageRate *big.Int
	})
	out := ret
	err := _InternalNetwork.contract.Call(opts, out, "getExpectedRate", src, dest, srcQty)
	return *ret, err
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_InternalNetwork *InternalNetworkSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	return _InternalNetwork.Contract.GetExpectedRate(&_InternalNetwork.CallOpts, src, dest, srcQty)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) constant returns(uint256 expectedRate, uint256 slippageRate)
func (_InternalNetwork *InternalNetworkCallerSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	SlippageRate *big.Int
}, error) {
	return _InternalNetwork.Contract.GetExpectedRate(&_InternalNetwork.CallOpts, src, dest, srcQty)
}

// GetNumReserves is a free data retrieval call binding the contract method 0xcfff25bb.
//
// Solidity: function getNumReserves() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) GetNumReserves(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getNumReserves")
	return *ret0, err
}

// GetNumReserves is a free data retrieval call binding the contract method 0xcfff25bb.
//
// Solidity: function getNumReserves() constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) GetNumReserves() (*big.Int, error) {
	return _InternalNetwork.Contract.GetNumReserves(&_InternalNetwork.CallOpts)
}

// GetNumReserves is a free data retrieval call binding the contract method 0xcfff25bb.
//
// Solidity: function getNumReserves() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) GetNumReserves() (*big.Int, error) {
	return _InternalNetwork.Contract.GetNumReserves(&_InternalNetwork.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_InternalNetwork *InternalNetworkCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_InternalNetwork *InternalNetworkSession) GetOperators() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetOperators(&_InternalNetwork.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() constant returns(address[])
func (_InternalNetwork *InternalNetworkCallerSession) GetOperators() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetOperators(&_InternalNetwork.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() constant returns(address[])
func (_InternalNetwork *InternalNetworkCaller) GetReserves(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getReserves")
	return *ret0, err
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() constant returns(address[])
func (_InternalNetwork *InternalNetworkSession) GetReserves() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetReserves(&_InternalNetwork.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() constant returns(address[])
func (_InternalNetwork *InternalNetworkCallerSession) GetReserves() ([]common.Address, error) {
	return _InternalNetwork.Contract.GetReserves(&_InternalNetwork.CallOpts)
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) GetUserCapInTokenWei(opts *bind.CallOpts, user common.Address, token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getUserCapInTokenWei", user, token)
	return *ret0, err
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) GetUserCapInTokenWei(user common.Address, token common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetUserCapInTokenWei(&_InternalNetwork.CallOpts, user, token)
}

// GetUserCapInTokenWei is a free data retrieval call binding the contract method 0x8eaaeecf.
//
// Solidity: function getUserCapInTokenWei(address user, address token) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) GetUserCapInTokenWei(user common.Address, token common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetUserCapInTokenWei(&_InternalNetwork.CallOpts, user, token)
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) GetUserCapInWei(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "getUserCapInWei", user)
	return *ret0, err
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) GetUserCapInWei(user common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetUserCapInWei(&_InternalNetwork.CallOpts, user)
}

// GetUserCapInWei is a free data retrieval call binding the contract method 0x6432679f.
//
// Solidity: function getUserCapInWei(address user) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) GetUserCapInWei(user common.Address) (*big.Int, error) {
	return _InternalNetwork.Contract.GetUserCapInWei(&_InternalNetwork.CallOpts, user)
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) Info(opts *bind.CallOpts, field [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "info", field)
	return *ret0, err
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) Info(field [32]byte) (*big.Int, error) {
	return _InternalNetwork.Contract.Info(&_InternalNetwork.CallOpts, field)
}

// Info is a free data retrieval call binding the contract method 0xb64a097e.
//
// Solidity: function info(bytes32 field) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) Info(field [32]byte) (*big.Int, error) {
	return _InternalNetwork.Contract.Info(&_InternalNetwork.CallOpts, field)
}

// InfoFields is a free data retrieval call binding the contract method 0x1610b59b.
//
// Solidity: function infoFields(bytes32 ) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) InfoFields(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "infoFields", arg0)
	return *ret0, err
}

// InfoFields is a free data retrieval call binding the contract method 0x1610b59b.
//
// Solidity: function infoFields(bytes32 ) constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) InfoFields(arg0 [32]byte) (*big.Int, error) {
	return _InternalNetwork.Contract.InfoFields(&_InternalNetwork.CallOpts, arg0)
}

// InfoFields is a free data retrieval call binding the contract method 0x1610b59b.
//
// Solidity: function infoFields(bytes32 ) constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) InfoFields(arg0 [32]byte) (*big.Int, error) {
	return _InternalNetwork.Contract.InfoFields(&_InternalNetwork.CallOpts, arg0)
}

// IsEnabled is a free data retrieval call binding the contract method 0x6aa633b6.
//
// Solidity: function isEnabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkCaller) IsEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "isEnabled")
	return *ret0, err
}

// IsEnabled is a free data retrieval call binding the contract method 0x6aa633b6.
//
// Solidity: function isEnabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkSession) IsEnabled() (bool, error) {
	return _InternalNetwork.Contract.IsEnabled(&_InternalNetwork.CallOpts)
}

// IsEnabled is a free data retrieval call binding the contract method 0x6aa633b6.
//
// Solidity: function isEnabled() constant returns(bool)
func (_InternalNetwork *InternalNetworkCallerSession) IsEnabled() (bool, error) {
	return _InternalNetwork.Contract.IsEnabled(&_InternalNetwork.CallOpts)
}

// IsReserve is a free data retrieval call binding the contract method 0x7a2b0587.
//
// Solidity: function isReserve(address ) constant returns(bool)
func (_InternalNetwork *InternalNetworkCaller) IsReserve(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "isReserve", arg0)
	return *ret0, err
}

// IsReserve is a free data retrieval call binding the contract method 0x7a2b0587.
//
// Solidity: function isReserve(address ) constant returns(bool)
func (_InternalNetwork *InternalNetworkSession) IsReserve(arg0 common.Address) (bool, error) {
	return _InternalNetwork.Contract.IsReserve(&_InternalNetwork.CallOpts, arg0)
}

// IsReserve is a free data retrieval call binding the contract method 0x7a2b0587.
//
// Solidity: function isReserve(address ) constant returns(bool)
func (_InternalNetwork *InternalNetworkCallerSession) IsReserve(arg0 common.Address) (bool, error) {
	return _InternalNetwork.Contract.IsReserve(&_InternalNetwork.CallOpts, arg0)
}

// KyberNetworkProxyContract is a free data retrieval call binding the contract method 0x785250da.
//
// Solidity: function kyberNetworkProxyContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) KyberNetworkProxyContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "kyberNetworkProxyContract")
	return *ret0, err
}

// KyberNetworkProxyContract is a free data retrieval call binding the contract method 0x785250da.
//
// Solidity: function kyberNetworkProxyContract() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) KyberNetworkProxyContract() (common.Address, error) {
	return _InternalNetwork.Contract.KyberNetworkProxyContract(&_InternalNetwork.CallOpts)
}

// KyberNetworkProxyContract is a free data retrieval call binding the contract method 0x785250da.
//
// Solidity: function kyberNetworkProxyContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) KyberNetworkProxyContract() (common.Address, error) {
	return _InternalNetwork.Contract.KyberNetworkProxyContract(&_InternalNetwork.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) MaxGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "maxGasPrice")
	return *ret0, err
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) MaxGasPrice() (*big.Int, error) {
	return _InternalNetwork.Contract.MaxGasPrice(&_InternalNetwork.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) MaxGasPrice() (*big.Int, error) {
	return _InternalNetwork.Contract.MaxGasPrice(&_InternalNetwork.CallOpts)
}

// MaxGasPriceValue is a free data retrieval call binding the contract method 0xb2d111f6.
//
// Solidity: function maxGasPriceValue() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) MaxGasPriceValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "maxGasPriceValue")
	return *ret0, err
}

// MaxGasPriceValue is a free data retrieval call binding the contract method 0xb2d111f6.
//
// Solidity: function maxGasPriceValue() constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) MaxGasPriceValue() (*big.Int, error) {
	return _InternalNetwork.Contract.MaxGasPriceValue(&_InternalNetwork.CallOpts)
}

// MaxGasPriceValue is a free data retrieval call binding the contract method 0xb2d111f6.
//
// Solidity: function maxGasPriceValue() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) MaxGasPriceValue() (*big.Int, error) {
	return _InternalNetwork.Contract.MaxGasPriceValue(&_InternalNetwork.CallOpts)
}

// NegligibleRateDiff is a free data retrieval call binding the contract method 0x4cef5a5c.
//
// Solidity: function negligibleRateDiff() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCaller) NegligibleRateDiff(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "negligibleRateDiff")
	return *ret0, err
}

// NegligibleRateDiff is a free data retrieval call binding the contract method 0x4cef5a5c.
//
// Solidity: function negligibleRateDiff() constant returns(uint256)
func (_InternalNetwork *InternalNetworkSession) NegligibleRateDiff() (*big.Int, error) {
	return _InternalNetwork.Contract.NegligibleRateDiff(&_InternalNetwork.CallOpts)
}

// NegligibleRateDiff is a free data retrieval call binding the contract method 0x4cef5a5c.
//
// Solidity: function negligibleRateDiff() constant returns(uint256)
func (_InternalNetwork *InternalNetworkCallerSession) NegligibleRateDiff() (*big.Int, error) {
	return _InternalNetwork.Contract.NegligibleRateDiff(&_InternalNetwork.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) PendingAdmin() (common.Address, error) {
	return _InternalNetwork.Contract.PendingAdmin(&_InternalNetwork.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) PendingAdmin() (common.Address, error) {
	return _InternalNetwork.Contract.PendingAdmin(&_InternalNetwork.CallOpts)
}

// Reserves is a free data retrieval call binding the contract method 0x8334278d.
//
// Solidity: function reserves(uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) Reserves(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "reserves", arg0)
	return *ret0, err
}

// Reserves is a free data retrieval call binding the contract method 0x8334278d.
//
// Solidity: function reserves(uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkSession) Reserves(arg0 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.Reserves(&_InternalNetwork.CallOpts, arg0)
}

// Reserves is a free data retrieval call binding the contract method 0x8334278d.
//
// Solidity: function reserves(uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) Reserves(arg0 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.Reserves(&_InternalNetwork.CallOpts, arg0)
}

// ReservesPerTokenDest is a free data retrieval call binding the contract method 0x937e909b.
//
// Solidity: function reservesPerTokenDest(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) ReservesPerTokenDest(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "reservesPerTokenDest", arg0, arg1)
	return *ret0, err
}

// ReservesPerTokenDest is a free data retrieval call binding the contract method 0x937e909b.
//
// Solidity: function reservesPerTokenDest(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkSession) ReservesPerTokenDest(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.ReservesPerTokenDest(&_InternalNetwork.CallOpts, arg0, arg1)
}

// ReservesPerTokenDest is a free data retrieval call binding the contract method 0x937e909b.
//
// Solidity: function reservesPerTokenDest(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) ReservesPerTokenDest(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.ReservesPerTokenDest(&_InternalNetwork.CallOpts, arg0, arg1)
}

// ReservesPerTokenSrc is a free data retrieval call binding the contract method 0x2ab8fc2d.
//
// Solidity: function reservesPerTokenSrc(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) ReservesPerTokenSrc(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "reservesPerTokenSrc", arg0, arg1)
	return *ret0, err
}

// ReservesPerTokenSrc is a free data retrieval call binding the contract method 0x2ab8fc2d.
//
// Solidity: function reservesPerTokenSrc(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkSession) ReservesPerTokenSrc(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.ReservesPerTokenSrc(&_InternalNetwork.CallOpts, arg0, arg1)
}

// ReservesPerTokenSrc is a free data retrieval call binding the contract method 0x2ab8fc2d.
//
// Solidity: function reservesPerTokenSrc(address , uint256 ) constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) ReservesPerTokenSrc(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _InternalNetwork.Contract.ReservesPerTokenSrc(&_InternalNetwork.CallOpts, arg0, arg1)
}

// SearchBestRate is a free data retrieval call binding the contract method 0xab127a0c.
//
// Solidity: function searchBestRate(address src, address dest, uint256 srcAmount) constant returns(address, uint256)
func (_InternalNetwork *InternalNetworkCaller) SearchBestRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcAmount *big.Int) (common.Address, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _InternalNetwork.contract.Call(opts, out, "searchBestRate", src, dest, srcAmount)
	return *ret0, *ret1, err
}

// SearchBestRate is a free data retrieval call binding the contract method 0xab127a0c.
//
// Solidity: function searchBestRate(address src, address dest, uint256 srcAmount) constant returns(address, uint256)
func (_InternalNetwork *InternalNetworkSession) SearchBestRate(src common.Address, dest common.Address, srcAmount *big.Int) (common.Address, *big.Int, error) {
	return _InternalNetwork.Contract.SearchBestRate(&_InternalNetwork.CallOpts, src, dest, srcAmount)
}

// SearchBestRate is a free data retrieval call binding the contract method 0xab127a0c.
//
// Solidity: function searchBestRate(address src, address dest, uint256 srcAmount) constant returns(address, uint256)
func (_InternalNetwork *InternalNetworkCallerSession) SearchBestRate(src common.Address, dest common.Address, srcAmount *big.Int) (common.Address, *big.Int, error) {
	return _InternalNetwork.Contract.SearchBestRate(&_InternalNetwork.CallOpts, src, dest, srcAmount)
}

// WhiteListContract is a free data retrieval call binding the contract method 0x5ed5ea28.
//
// Solidity: function whiteListContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCaller) WhiteListContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _InternalNetwork.contract.Call(opts, out, "whiteListContract")
	return *ret0, err
}

// WhiteListContract is a free data retrieval call binding the contract method 0x5ed5ea28.
//
// Solidity: function whiteListContract() constant returns(address)
func (_InternalNetwork *InternalNetworkSession) WhiteListContract() (common.Address, error) {
	return _InternalNetwork.Contract.WhiteListContract(&_InternalNetwork.CallOpts)
}

// WhiteListContract is a free data retrieval call binding the contract method 0x5ed5ea28.
//
// Solidity: function whiteListContract() constant returns(address)
func (_InternalNetwork *InternalNetworkCallerSession) WhiteListContract() (common.Address, error) {
	return _InternalNetwork.Contract.WhiteListContract(&_InternalNetwork.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_InternalNetwork *InternalNetworkTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_InternalNetwork *InternalNetworkSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddAlerter(&_InternalNetwork.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddAlerter(&_InternalNetwork.TransactOpts, newAlerter)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_InternalNetwork *InternalNetworkTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_InternalNetwork *InternalNetworkSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddOperator(&_InternalNetwork.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddOperator(&_InternalNetwork.TransactOpts, newOperator)
}

// AddReserve is a paid mutator transaction binding the contract method 0xa0d7bb1b.
//
// Solidity: function addReserve(address reserve, bool add) returns()
func (_InternalNetwork *InternalNetworkTransactor) AddReserve(opts *bind.TransactOpts, reserve common.Address, add bool) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "addReserve", reserve, add)
}

// AddReserve is a paid mutator transaction binding the contract method 0xa0d7bb1b.
//
// Solidity: function addReserve(address reserve, bool add) returns()
func (_InternalNetwork *InternalNetworkSession) AddReserve(reserve common.Address, add bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddReserve(&_InternalNetwork.TransactOpts, reserve, add)
}

// AddReserve is a paid mutator transaction binding the contract method 0xa0d7bb1b.
//
// Solidity: function addReserve(address reserve, bool add) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) AddReserve(reserve common.Address, add bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.AddReserve(&_InternalNetwork.TransactOpts, reserve, add)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_InternalNetwork *InternalNetworkTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_InternalNetwork *InternalNetworkSession) ClaimAdmin() (*types.Transaction, error) {
	return _InternalNetwork.Contract.ClaimAdmin(&_InternalNetwork.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_InternalNetwork *InternalNetworkTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _InternalNetwork.Contract.ClaimAdmin(&_InternalNetwork.TransactOpts)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0xe02584bf.
//
// Solidity: function listPairForReserve(address reserve, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_InternalNetwork *InternalNetworkTransactor) ListPairForReserve(opts *bind.TransactOpts, reserve common.Address, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "listPairForReserve", reserve, token, ethToToken, tokenToEth, add)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0xe02584bf.
//
// Solidity: function listPairForReserve(address reserve, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_InternalNetwork *InternalNetworkSession) ListPairForReserve(reserve common.Address, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.ListPairForReserve(&_InternalNetwork.TransactOpts, reserve, token, ethToToken, tokenToEth, add)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0xe02584bf.
//
// Solidity: function listPairForReserve(address reserve, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) ListPairForReserve(reserve common.Address, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.ListPairForReserve(&_InternalNetwork.TransactOpts, reserve, token, ethToToken, tokenToEth, add)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_InternalNetwork *InternalNetworkTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_InternalNetwork *InternalNetworkSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.RemoveAlerter(&_InternalNetwork.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.RemoveAlerter(&_InternalNetwork.TransactOpts, alerter)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_InternalNetwork *InternalNetworkTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_InternalNetwork *InternalNetworkSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.RemoveOperator(&_InternalNetwork.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.RemoveOperator(&_InternalNetwork.TransactOpts, operator)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool _enable) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetEnable(opts *bind.TransactOpts, _enable bool) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setEnable", _enable)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool _enable) returns()
func (_InternalNetwork *InternalNetworkSession) SetEnable(_enable bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetEnable(&_InternalNetwork.TransactOpts, _enable)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool _enable) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetEnable(_enable bool) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetEnable(&_InternalNetwork.TransactOpts, _enable)
}

// SetExpectedRate is a paid mutator transaction binding the contract method 0x5d270cdc.
//
// Solidity: function setExpectedRate(address expectedRate) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetExpectedRate(opts *bind.TransactOpts, expectedRate common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setExpectedRate", expectedRate)
}

// SetExpectedRate is a paid mutator transaction binding the contract method 0x5d270cdc.
//
// Solidity: function setExpectedRate(address expectedRate) returns()
func (_InternalNetwork *InternalNetworkSession) SetExpectedRate(expectedRate common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetExpectedRate(&_InternalNetwork.TransactOpts, expectedRate)
}

// SetExpectedRate is a paid mutator transaction binding the contract method 0x5d270cdc.
//
// Solidity: function setExpectedRate(address expectedRate) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetExpectedRate(expectedRate common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetExpectedRate(&_InternalNetwork.TransactOpts, expectedRate)
}

// SetFeeBurner is a paid mutator transaction binding the contract method 0x1a79464e.
//
// Solidity: function setFeeBurner(address feeBurner) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetFeeBurner(opts *bind.TransactOpts, feeBurner common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setFeeBurner", feeBurner)
}

// SetFeeBurner is a paid mutator transaction binding the contract method 0x1a79464e.
//
// Solidity: function setFeeBurner(address feeBurner) returns()
func (_InternalNetwork *InternalNetworkSession) SetFeeBurner(feeBurner common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetFeeBurner(&_InternalNetwork.TransactOpts, feeBurner)
}

// SetFeeBurner is a paid mutator transaction binding the contract method 0x1a79464e.
//
// Solidity: function setFeeBurner(address feeBurner) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetFeeBurner(feeBurner common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetFeeBurner(&_InternalNetwork.TransactOpts, feeBurner)
}

// SetInfo is a paid mutator transaction binding the contract method 0x5f65d703.
//
// Solidity: function setInfo(bytes32 field, uint256 value) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetInfo(opts *bind.TransactOpts, field [32]byte, value *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setInfo", field, value)
}

// SetInfo is a paid mutator transaction binding the contract method 0x5f65d703.
//
// Solidity: function setInfo(bytes32 field, uint256 value) returns()
func (_InternalNetwork *InternalNetworkSession) SetInfo(field [32]byte, value *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetInfo(&_InternalNetwork.TransactOpts, field, value)
}

// SetInfo is a paid mutator transaction binding the contract method 0x5f65d703.
//
// Solidity: function setInfo(bytes32 field, uint256 value) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetInfo(field [32]byte, value *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetInfo(&_InternalNetwork.TransactOpts, field, value)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address networkProxy) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetKyberProxy(opts *bind.TransactOpts, networkProxy common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setKyberProxy", networkProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address networkProxy) returns()
func (_InternalNetwork *InternalNetworkSession) SetKyberProxy(networkProxy common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetKyberProxy(&_InternalNetwork.TransactOpts, networkProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address networkProxy) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetKyberProxy(networkProxy common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetKyberProxy(&_InternalNetwork.TransactOpts, networkProxy)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiff) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetParams(opts *bind.TransactOpts, _maxGasPrice *big.Int, _negligibleRateDiff *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setParams", _maxGasPrice, _negligibleRateDiff)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiff) returns()
func (_InternalNetwork *InternalNetworkSession) SetParams(_maxGasPrice *big.Int, _negligibleRateDiff *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetParams(&_InternalNetwork.TransactOpts, _maxGasPrice, _negligibleRateDiff)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiff) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetParams(_maxGasPrice *big.Int, _negligibleRateDiff *big.Int) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetParams(&_InternalNetwork.TransactOpts, _maxGasPrice, _negligibleRateDiff)
}

// SetWhiteList is a paid mutator transaction binding the contract method 0x39e899ee.
//
// Solidity: function setWhiteList(address whiteList) returns()
func (_InternalNetwork *InternalNetworkTransactor) SetWhiteList(opts *bind.TransactOpts, whiteList common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "setWhiteList", whiteList)
}

// SetWhiteList is a paid mutator transaction binding the contract method 0x39e899ee.
//
// Solidity: function setWhiteList(address whiteList) returns()
func (_InternalNetwork *InternalNetworkSession) SetWhiteList(whiteList common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetWhiteList(&_InternalNetwork.TransactOpts, whiteList)
}

// SetWhiteList is a paid mutator transaction binding the contract method 0x39e899ee.
//
// Solidity: function setWhiteList(address whiteList) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) SetWhiteList(whiteList common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.SetWhiteList(&_InternalNetwork.TransactOpts, whiteList)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_InternalNetwork *InternalNetworkTransactor) TradeWithHint(opts *bind.TransactOpts, trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "tradeWithHint", trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_InternalNetwork *InternalNetworkSession) TradeWithHint(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TradeWithHint(&_InternalNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) returns(uint256)
func (_InternalNetwork *InternalNetworkTransactorSession) TradeWithHint(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TradeWithHint(&_InternalNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TransferAdmin(&_InternalNetwork.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TransferAdmin(&_InternalNetwork.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TransferAdminQuickly(&_InternalNetwork.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.TransferAdminQuickly(&_InternalNetwork.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.WithdrawEther(&_InternalNetwork.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.WithdrawEther(&_InternalNetwork.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.WithdrawToken(&_InternalNetwork.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_InternalNetwork *InternalNetworkTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _InternalNetwork.Contract.WithdrawToken(&_InternalNetwork.TransactOpts, token, amount, sendTo)
}

// InternalNetworkAddReserveToNetworkIterator is returned from FilterAddReserveToNetwork and is used to iterate over the raw logs and unpacked data for AddReserveToNetwork events raised by the InternalNetwork contract.
type InternalNetworkAddReserveToNetworkIterator struct {
	Event *InternalNetworkAddReserveToNetwork // Event containing the contract specifics and raw log

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
func (it *InternalNetworkAddReserveToNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkAddReserveToNetwork)
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
		it.Event = new(InternalNetworkAddReserveToNetwork)
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
func (it *InternalNetworkAddReserveToNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkAddReserveToNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkAddReserveToNetwork represents a AddReserveToNetwork event raised by the InternalNetwork contract.
type InternalNetworkAddReserveToNetwork struct {
	Reserve common.Address
	Add     bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddReserveToNetwork is a free log retrieval operation binding the contract event 0x7752182b29e356eb432239f464340b4481e1b0bfad97f06aa2ff8cdc74611449.
//
// Solidity: event AddReserveToNetwork(address reserve, bool add)
func (_InternalNetwork *InternalNetworkFilterer) FilterAddReserveToNetwork(opts *bind.FilterOpts) (*InternalNetworkAddReserveToNetworkIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "AddReserveToNetwork")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkAddReserveToNetworkIterator{contract: _InternalNetwork.contract, event: "AddReserveToNetwork", logs: logs, sub: sub}, nil
}

// WatchAddReserveToNetwork is a free log subscription operation binding the contract event 0x7752182b29e356eb432239f464340b4481e1b0bfad97f06aa2ff8cdc74611449.
//
// Solidity: event AddReserveToNetwork(address reserve, bool add)
func (_InternalNetwork *InternalNetworkFilterer) WatchAddReserveToNetwork(opts *bind.WatchOpts, sink chan<- *InternalNetworkAddReserveToNetwork) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "AddReserveToNetwork")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkAddReserveToNetwork)
				if err := _InternalNetwork.contract.UnpackLog(event, "AddReserveToNetwork", log); err != nil {
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

// InternalNetworkAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the InternalNetwork contract.
type InternalNetworkAdminClaimedIterator struct {
	Event *InternalNetworkAdminClaimed // Event containing the contract specifics and raw log

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
func (it *InternalNetworkAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkAdminClaimed)
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
		it.Event = new(InternalNetworkAdminClaimed)
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
func (it *InternalNetworkAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkAdminClaimed represents a AdminClaimed event raised by the InternalNetwork contract.
type InternalNetworkAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_InternalNetwork *InternalNetworkFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*InternalNetworkAdminClaimedIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkAdminClaimedIterator{contract: _InternalNetwork.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_InternalNetwork *InternalNetworkFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *InternalNetworkAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkAdminClaimed)
				if err := _InternalNetwork.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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

// InternalNetworkAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the InternalNetwork contract.
type InternalNetworkAlerterAddedIterator struct {
	Event *InternalNetworkAlerterAdded // Event containing the contract specifics and raw log

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
func (it *InternalNetworkAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkAlerterAdded)
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
		it.Event = new(InternalNetworkAlerterAdded)
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
func (it *InternalNetworkAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkAlerterAdded represents a AlerterAdded event raised by the InternalNetwork contract.
type InternalNetworkAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_InternalNetwork *InternalNetworkFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*InternalNetworkAlerterAddedIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkAlerterAddedIterator{contract: _InternalNetwork.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_InternalNetwork *InternalNetworkFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *InternalNetworkAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkAlerterAdded)
				if err := _InternalNetwork.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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

// InternalNetworkEtherReceivalIterator is returned from FilterEtherReceival and is used to iterate over the raw logs and unpacked data for EtherReceival events raised by the InternalNetwork contract.
type InternalNetworkEtherReceivalIterator struct {
	Event *InternalNetworkEtherReceival // Event containing the contract specifics and raw log

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
func (it *InternalNetworkEtherReceivalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkEtherReceival)
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
		it.Event = new(InternalNetworkEtherReceival)
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
func (it *InternalNetworkEtherReceivalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkEtherReceivalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkEtherReceival represents a EtherReceival event raised by the InternalNetwork contract.
type InternalNetworkEtherReceival struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherReceival is a free log retrieval operation binding the contract event 0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619.
//
// Solidity: event EtherReceival(address indexed sender, uint256 amount)
func (_InternalNetwork *InternalNetworkFilterer) FilterEtherReceival(opts *bind.FilterOpts, sender []common.Address) (*InternalNetworkEtherReceivalIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "EtherReceival", senderRule)
	if err != nil {
		return nil, err
	}
	return &InternalNetworkEtherReceivalIterator{contract: _InternalNetwork.contract, event: "EtherReceival", logs: logs, sub: sub}, nil
}

// WatchEtherReceival is a free log subscription operation binding the contract event 0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619.
//
// Solidity: event EtherReceival(address indexed sender, uint256 amount)
func (_InternalNetwork *InternalNetworkFilterer) WatchEtherReceival(opts *bind.WatchOpts, sink chan<- *InternalNetworkEtherReceival, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "EtherReceival", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkEtherReceival)
				if err := _InternalNetwork.contract.UnpackLog(event, "EtherReceival", log); err != nil {
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

// InternalNetworkEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the InternalNetwork contract.
type InternalNetworkEtherWithdrawIterator struct {
	Event *InternalNetworkEtherWithdraw // Event containing the contract specifics and raw log

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
func (it *InternalNetworkEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkEtherWithdraw)
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
		it.Event = new(InternalNetworkEtherWithdraw)
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
func (it *InternalNetworkEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkEtherWithdraw represents a EtherWithdraw event raised by the InternalNetwork contract.
type InternalNetworkEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_InternalNetwork *InternalNetworkFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*InternalNetworkEtherWithdrawIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkEtherWithdrawIterator{contract: _InternalNetwork.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_InternalNetwork *InternalNetworkFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *InternalNetworkEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkEtherWithdraw)
				if err := _InternalNetwork.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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

// InternalNetworkKyberProxySetIterator is returned from FilterKyberProxySet and is used to iterate over the raw logs and unpacked data for KyberProxySet events raised by the InternalNetwork contract.
type InternalNetworkKyberProxySetIterator struct {
	Event *InternalNetworkKyberProxySet // Event containing the contract specifics and raw log

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
func (it *InternalNetworkKyberProxySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkKyberProxySet)
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
		it.Event = new(InternalNetworkKyberProxySet)
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
func (it *InternalNetworkKyberProxySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkKyberProxySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkKyberProxySet represents a KyberProxySet event raised by the InternalNetwork contract.
type InternalNetworkKyberProxySet struct {
	Proxy  common.Address
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterKyberProxySet is a free log retrieval operation binding the contract event 0xfdd305502f7797ff3390aa08825f7f6aec92c27a94e103bfaf45452b4cf1d4f4.
//
// Solidity: event KyberProxySet(address proxy, address sender)
func (_InternalNetwork *InternalNetworkFilterer) FilterKyberProxySet(opts *bind.FilterOpts) (*InternalNetworkKyberProxySetIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "KyberProxySet")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkKyberProxySetIterator{contract: _InternalNetwork.contract, event: "KyberProxySet", logs: logs, sub: sub}, nil
}

// WatchKyberProxySet is a free log subscription operation binding the contract event 0xfdd305502f7797ff3390aa08825f7f6aec92c27a94e103bfaf45452b4cf1d4f4.
//
// Solidity: event KyberProxySet(address proxy, address sender)
func (_InternalNetwork *InternalNetworkFilterer) WatchKyberProxySet(opts *bind.WatchOpts, sink chan<- *InternalNetworkKyberProxySet) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "KyberProxySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkKyberProxySet)
				if err := _InternalNetwork.contract.UnpackLog(event, "KyberProxySet", log); err != nil {
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

// InternalNetworkKyberTradeIterator is returned from FilterKyberTrade and is used to iterate over the raw logs and unpacked data for KyberTrade events raised by the InternalNetwork contract.
type InternalNetworkKyberTradeIterator struct {
	Event *InternalNetworkKyberTrade // Event containing the contract specifics and raw log

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
func (it *InternalNetworkKyberTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkKyberTrade)
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
		it.Event = new(InternalNetworkKyberTrade)
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
func (it *InternalNetworkKyberTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkKyberTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkKyberTrade represents a KyberTrade event raised by the InternalNetwork contract.
type InternalNetworkKyberTrade struct {
	SrcAddress  common.Address
	SrcToken    common.Address
	SrcAmount   *big.Int
	DestAddress common.Address
	DestToken   common.Address
	DestAmount  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterKyberTrade is a free log retrieval operation binding the contract event 0x1c8399ecc5c956b9cb18c820248b10b634cca4af308755e07cd467655e8ec3c7.
//
// Solidity: event KyberTrade(address srcAddress, address srcToken, uint256 srcAmount, address destAddress, address destToken, uint256 destAmount)
func (_InternalNetwork *InternalNetworkFilterer) FilterKyberTrade(opts *bind.FilterOpts) (*InternalNetworkKyberTradeIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "KyberTrade")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkKyberTradeIterator{contract: _InternalNetwork.contract, event: "KyberTrade", logs: logs, sub: sub}, nil
}

// WatchKyberTrade is a free log subscription operation binding the contract event 0x1c8399ecc5c956b9cb18c820248b10b634cca4af308755e07cd467655e8ec3c7.
//
// Solidity: event KyberTrade(address srcAddress, address srcToken, uint256 srcAmount, address destAddress, address destToken, uint256 destAmount)
func (_InternalNetwork *InternalNetworkFilterer) WatchKyberTrade(opts *bind.WatchOpts, sink chan<- *InternalNetworkKyberTrade) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "KyberTrade")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkKyberTrade)
				if err := _InternalNetwork.contract.UnpackLog(event, "KyberTrade", log); err != nil {
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

// InternalNetworkListReservePairsIterator is returned from FilterListReservePairs and is used to iterate over the raw logs and unpacked data for ListReservePairs events raised by the InternalNetwork contract.
type InternalNetworkListReservePairsIterator struct {
	Event *InternalNetworkListReservePairs // Event containing the contract specifics and raw log

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
func (it *InternalNetworkListReservePairsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkListReservePairs)
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
		it.Event = new(InternalNetworkListReservePairs)
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
func (it *InternalNetworkListReservePairsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkListReservePairsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkListReservePairs represents a ListReservePairs event raised by the InternalNetwork contract.
type InternalNetworkListReservePairs struct {
	Reserve common.Address
	Src     common.Address
	Dest    common.Address
	Add     bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterListReservePairs is a free log retrieval operation binding the contract event 0xadb5a4f14d89b3a5ffb3900ac1ea4574d991f93887f6199fabaf25393644e01c.
//
// Solidity: event ListReservePairs(address reserve, address src, address dest, bool add)
func (_InternalNetwork *InternalNetworkFilterer) FilterListReservePairs(opts *bind.FilterOpts) (*InternalNetworkListReservePairsIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "ListReservePairs")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkListReservePairsIterator{contract: _InternalNetwork.contract, event: "ListReservePairs", logs: logs, sub: sub}, nil
}

// WatchListReservePairs is a free log subscription operation binding the contract event 0xadb5a4f14d89b3a5ffb3900ac1ea4574d991f93887f6199fabaf25393644e01c.
//
// Solidity: event ListReservePairs(address reserve, address src, address dest, bool add)
func (_InternalNetwork *InternalNetworkFilterer) WatchListReservePairs(opts *bind.WatchOpts, sink chan<- *InternalNetworkListReservePairs) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "ListReservePairs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkListReservePairs)
				if err := _InternalNetwork.contract.UnpackLog(event, "ListReservePairs", log); err != nil {
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

// InternalNetworkOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the InternalNetwork contract.
type InternalNetworkOperatorAddedIterator struct {
	Event *InternalNetworkOperatorAdded // Event containing the contract specifics and raw log

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
func (it *InternalNetworkOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkOperatorAdded)
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
		it.Event = new(InternalNetworkOperatorAdded)
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
func (it *InternalNetworkOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkOperatorAdded represents a OperatorAdded event raised by the InternalNetwork contract.
type InternalNetworkOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_InternalNetwork *InternalNetworkFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*InternalNetworkOperatorAddedIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkOperatorAddedIterator{contract: _InternalNetwork.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_InternalNetwork *InternalNetworkFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *InternalNetworkOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkOperatorAdded)
				if err := _InternalNetwork.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// InternalNetworkTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the InternalNetwork contract.
type InternalNetworkTokenWithdrawIterator struct {
	Event *InternalNetworkTokenWithdraw // Event containing the contract specifics and raw log

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
func (it *InternalNetworkTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkTokenWithdraw)
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
		it.Event = new(InternalNetworkTokenWithdraw)
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
func (it *InternalNetworkTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkTokenWithdraw represents a TokenWithdraw event raised by the InternalNetwork contract.
type InternalNetworkTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_InternalNetwork *InternalNetworkFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*InternalNetworkTokenWithdrawIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkTokenWithdrawIterator{contract: _InternalNetwork.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_InternalNetwork *InternalNetworkFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *InternalNetworkTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkTokenWithdraw)
				if err := _InternalNetwork.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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

// InternalNetworkTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the InternalNetwork contract.
type InternalNetworkTransferAdminPendingIterator struct {
	Event *InternalNetworkTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *InternalNetworkTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InternalNetworkTransferAdminPending)
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
		it.Event = new(InternalNetworkTransferAdminPending)
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
func (it *InternalNetworkTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InternalNetworkTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InternalNetworkTransferAdminPending represents a TransferAdminPending event raised by the InternalNetwork contract.
type InternalNetworkTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_InternalNetwork *InternalNetworkFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*InternalNetworkTransferAdminPendingIterator, error) {

	logs, sub, err := _InternalNetwork.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &InternalNetworkTransferAdminPendingIterator{contract: _InternalNetwork.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_InternalNetwork *InternalNetworkFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *InternalNetworkTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _InternalNetwork.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InternalNetworkTransferAdminPending)
				if err := _InternalNetwork.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
