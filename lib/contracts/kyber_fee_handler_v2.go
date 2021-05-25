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

// KyberFeeHandlerV2ABI is the input ABI used to generate the binding from.
const KyberFeeHandlerV2ABI = "[{\"inputs\":[{\"internalType\":\"contractIKyberProxy\",\"name\":\"_kyberProxy\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_knc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_burnBlockInterval\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_daoOperator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_feePool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rewardBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rebateBps\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebateBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnBps\",\"type\":\"uint256\"}],\"name\":\"BRRUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractISanityRate\",\"name\":\"sanityRate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiToBurn\",\"type\":\"uint256\"}],\"name\":\"BurnConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"daoOperator\",\"type\":\"address\"}],\"name\":\"DaoOperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"platformFeeWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebateWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"rebateWallets\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"rebatePercentBpsPerWallet\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmtWei\",\"type\":\"uint256\"}],\"name\":\"FeeDistributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"feePool\",\"type\":\"address\"}],\"name\":\"FeePoolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"kncTWei\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"KncBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberProxy\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"KyberProxyUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PlatformFeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RebatePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardPaid\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"brrData\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"rewardBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"rebateBps\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnBlockInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnKnc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"kncBurnAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"}],\"name\":\"claimPlatformFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountWei\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"}],\"name\":\"claimReserveRebate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountWei\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"claimStakerReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daoOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"feePerPlatformWallet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feePool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestSanityRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"kncToEthSanityRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSanityRateContracts\",\"outputs\":[{\"internalType\":\"contractISanityRate[]\",\"name\":\"sanityRates\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"rebateWallets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"rebateBpsPerWallet\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"platformFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"networkFee\",\"type\":\"uint256\"}],\"name\":\"handleFees\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"knc\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberProxy\",\"outputs\":[{\"internalType\":\"contractIKyberProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBurnBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"readBRRData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rewardBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rebateBps\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"rebatePerWallet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_burnBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rewardBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rebateBps\",\"type\":\"uint256\"}],\"name\":\"setBRRData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISanityRate\",\"name\":\"_sanityRate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_weiToBurn\",\"type\":\"uint256\"}],\"name\":\"setBurnConfigParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_daoOperator\",\"type\":\"address\"}],\"name\":\"setDaoOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feePool\",\"type\":\"address\"}],\"name\":\"setFeePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberProxy\",\"name\":\"_newProxy\",\"type\":\"address\"}],\"name\":\"setKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPayoutBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"weiToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// KyberFeeHandlerV2 is an auto generated Go binding around an Ethereum contract.
type KyberFeeHandlerV2 struct {
	KyberFeeHandlerV2Caller     // Read-only binding to the contract
	KyberFeeHandlerV2Transactor // Write-only binding to the contract
	KyberFeeHandlerV2Filterer   // Log filterer for contract events
}

// KyberFeeHandlerV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type KyberFeeHandlerV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KyberFeeHandlerV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KyberFeeHandlerV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KyberFeeHandlerV2Session struct {
	Contract     *KyberFeeHandlerV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// KyberFeeHandlerV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KyberFeeHandlerV2CallerSession struct {
	Contract *KyberFeeHandlerV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// KyberFeeHandlerV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KyberFeeHandlerV2TransactorSession struct {
	Contract     *KyberFeeHandlerV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// KyberFeeHandlerV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type KyberFeeHandlerV2Raw struct {
	Contract *KyberFeeHandlerV2 // Generic contract binding to access the raw methods on
}

// KyberFeeHandlerV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KyberFeeHandlerV2CallerRaw struct {
	Contract *KyberFeeHandlerV2Caller // Generic read-only contract binding to access the raw methods on
}

// KyberFeeHandlerV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KyberFeeHandlerV2TransactorRaw struct {
	Contract *KyberFeeHandlerV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKyberFeeHandlerV2 creates a new instance of KyberFeeHandlerV2, bound to a specific deployed contract.
func NewKyberFeeHandlerV2(address common.Address, backend bind.ContractBackend) (*KyberFeeHandlerV2, error) {
	contract, err := bindKyberFeeHandlerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2{KyberFeeHandlerV2Caller: KyberFeeHandlerV2Caller{contract: contract}, KyberFeeHandlerV2Transactor: KyberFeeHandlerV2Transactor{contract: contract}, KyberFeeHandlerV2Filterer: KyberFeeHandlerV2Filterer{contract: contract}}, nil
}

// NewKyberFeeHandlerV2Caller creates a new read-only instance of KyberFeeHandlerV2, bound to a specific deployed contract.
func NewKyberFeeHandlerV2Caller(address common.Address, caller bind.ContractCaller) (*KyberFeeHandlerV2Caller, error) {
	contract, err := bindKyberFeeHandlerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2Caller{contract: contract}, nil
}

// NewKyberFeeHandlerV2Transactor creates a new write-only instance of KyberFeeHandlerV2, bound to a specific deployed contract.
func NewKyberFeeHandlerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*KyberFeeHandlerV2Transactor, error) {
	contract, err := bindKyberFeeHandlerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2Transactor{contract: contract}, nil
}

// NewKyberFeeHandlerV2Filterer creates a new log filterer instance of KyberFeeHandlerV2, bound to a specific deployed contract.
func NewKyberFeeHandlerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*KyberFeeHandlerV2Filterer, error) {
	contract, err := bindKyberFeeHandlerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2Filterer{contract: contract}, nil
}

// bindKyberFeeHandlerV2 binds a generic wrapper to an already deployed contract.
func bindKyberFeeHandlerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KyberFeeHandlerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberFeeHandlerV2.Contract.KyberFeeHandlerV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.KyberFeeHandlerV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.KyberFeeHandlerV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberFeeHandlerV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.contract.Transact(opts, method, params...)
}

// BrrData is a free data retrieval call binding the contract method 0x2e4da25c.
//
// Solidity: function brrData() view returns(uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) BrrData(opts *bind.CallOpts) (struct {
	RewardBps uint16
	RebateBps uint16
}, error) {
	ret := new(struct {
		RewardBps uint16
		RebateBps uint16
	})
	out := ret
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "brrData")
	return *ret, err
}

// BrrData is a free data retrieval call binding the contract method 0x2e4da25c.
//
// Solidity: function brrData() view returns(uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) BrrData() (struct {
	RewardBps uint16
	RebateBps uint16
}, error) {
	return _KyberFeeHandlerV2.Contract.BrrData(&_KyberFeeHandlerV2.CallOpts)
}

// BrrData is a free data retrieval call binding the contract method 0x2e4da25c.
//
// Solidity: function brrData() view returns(uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) BrrData() (struct {
	RewardBps uint16
	RebateBps uint16
}, error) {
	return _KyberFeeHandlerV2.Contract.BrrData(&_KyberFeeHandlerV2.CallOpts)
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) BurnBlockInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "burnBlockInterval")
	return *ret0, err
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) BurnBlockInterval() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.BurnBlockInterval(&_KyberFeeHandlerV2.CallOpts)
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) BurnBlockInterval() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.BurnBlockInterval(&_KyberFeeHandlerV2.CallOpts)
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) DaoOperator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "daoOperator")
	return *ret0, err
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) DaoOperator() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.DaoOperator(&_KyberFeeHandlerV2.CallOpts)
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) DaoOperator() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.DaoOperator(&_KyberFeeHandlerV2.CallOpts)
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) FeePerPlatformWallet(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "feePerPlatformWallet", arg0)
	return *ret0, err
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) FeePerPlatformWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.FeePerPlatformWallet(&_KyberFeeHandlerV2.CallOpts, arg0)
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) FeePerPlatformWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.FeePerPlatformWallet(&_KyberFeeHandlerV2.CallOpts, arg0)
}

// FeePool is a free data retrieval call binding the contract method 0xae2e933b.
//
// Solidity: function feePool() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) FeePool(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "feePool")
	return *ret0, err
}

// FeePool is a free data retrieval call binding the contract method 0xae2e933b.
//
// Solidity: function feePool() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) FeePool() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.FeePool(&_KyberFeeHandlerV2.CallOpts)
}

// FeePool is a free data retrieval call binding the contract method 0xae2e933b.
//
// Solidity: function feePool() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) FeePool() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.FeePool(&_KyberFeeHandlerV2.CallOpts)
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) GetLatestSanityRate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "getLatestSanityRate")
	return *ret0, err
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) GetLatestSanityRate() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.GetLatestSanityRate(&_KyberFeeHandlerV2.CallOpts)
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) GetLatestSanityRate() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.GetLatestSanityRate(&_KyberFeeHandlerV2.CallOpts)
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) GetSanityRateContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "getSanityRateContracts")
	return *ret0, err
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) GetSanityRateContracts() ([]common.Address, error) {
	return _KyberFeeHandlerV2.Contract.GetSanityRateContracts(&_KyberFeeHandlerV2.CallOpts)
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) GetSanityRateContracts() ([]common.Address, error) {
	return _KyberFeeHandlerV2.Contract.GetSanityRateContracts(&_KyberFeeHandlerV2.CallOpts)
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) Knc(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "knc")
	return *ret0, err
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) Knc() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.Knc(&_KyberFeeHandlerV2.CallOpts)
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) Knc() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.Knc(&_KyberFeeHandlerV2.CallOpts)
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) KyberProxy(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "kyberProxy")
	return *ret0, err
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) KyberProxy() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.KyberProxy(&_KyberFeeHandlerV2.CallOpts)
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) KyberProxy() (common.Address, error) {
	return _KyberFeeHandlerV2.Contract.KyberProxy(&_KyberFeeHandlerV2.CallOpts)
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) LastBurnBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "lastBurnBlock")
	return *ret0, err
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) LastBurnBlock() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.LastBurnBlock(&_KyberFeeHandlerV2.CallOpts)
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) LastBurnBlock() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.LastBurnBlock(&_KyberFeeHandlerV2.CallOpts)
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) ReadBRRData(opts *bind.CallOpts) (struct {
	RewardBps *big.Int
	RebateBps *big.Int
}, error) {
	ret := new(struct {
		RewardBps *big.Int
		RebateBps *big.Int
	})
	out := ret
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "readBRRData")
	return *ret, err
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) ReadBRRData() (struct {
	RewardBps *big.Int
	RebateBps *big.Int
}, error) {
	return _KyberFeeHandlerV2.Contract.ReadBRRData(&_KyberFeeHandlerV2.CallOpts)
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) ReadBRRData() (struct {
	RewardBps *big.Int
	RebateBps *big.Int
}, error) {
	return _KyberFeeHandlerV2.Contract.ReadBRRData(&_KyberFeeHandlerV2.CallOpts)
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) RebatePerWallet(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "rebatePerWallet", arg0)
	return *ret0, err
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) RebatePerWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.RebatePerWallet(&_KyberFeeHandlerV2.CallOpts, arg0)
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) RebatePerWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.RebatePerWallet(&_KyberFeeHandlerV2.CallOpts, arg0)
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) TotalPayoutBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "totalPayoutBalance")
	return *ret0, err
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) TotalPayoutBalance() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.TotalPayoutBalance(&_KyberFeeHandlerV2.CallOpts)
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) TotalPayoutBalance() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.TotalPayoutBalance(&_KyberFeeHandlerV2.CallOpts)
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Caller) WeiToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandlerV2.contract.Call(opts, out, "weiToBurn")
	return *ret0, err
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) WeiToBurn() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.WeiToBurn(&_KyberFeeHandlerV2.CallOpts)
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2CallerSession) WeiToBurn() (*big.Int, error) {
	return _KyberFeeHandlerV2.Contract.WeiToBurn(&_KyberFeeHandlerV2.CallOpts)
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) BurnKnc(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "burnKnc")
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) BurnKnc() (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.BurnKnc(&_KyberFeeHandlerV2.TransactOpts)
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) BurnKnc() (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.BurnKnc(&_KyberFeeHandlerV2.TransactOpts)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) ClaimPlatformFee(opts *bind.TransactOpts, platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "claimPlatformFee", platformWallet)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) ClaimPlatformFee(platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimPlatformFee(&_KyberFeeHandlerV2.TransactOpts, platformWallet)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) ClaimPlatformFee(platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimPlatformFee(&_KyberFeeHandlerV2.TransactOpts, platformWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) ClaimReserveRebate(opts *bind.TransactOpts, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "claimReserveRebate", rebateWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) ClaimReserveRebate(rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimReserveRebate(&_KyberFeeHandlerV2.TransactOpts, rebateWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) ClaimReserveRebate(rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimReserveRebate(&_KyberFeeHandlerV2.TransactOpts, rebateWallet)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address , uint256 ) returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) ClaimStakerReward(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "claimStakerReward", arg0, arg1)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address , uint256 ) returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) ClaimStakerReward(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimStakerReward(&_KyberFeeHandlerV2.TransactOpts, arg0, arg1)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address , uint256 ) returns(uint256)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) ClaimStakerReward(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.ClaimStakerReward(&_KyberFeeHandlerV2.TransactOpts, arg0, arg1)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) HandleFees(opts *bind.TransactOpts, token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "handleFees", token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) HandleFees(token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.HandleFees(&_KyberFeeHandlerV2.TransactOpts, token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) HandleFees(token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.HandleFees(&_KyberFeeHandlerV2.TransactOpts, token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// SetBRRData is a paid mutator transaction binding the contract method 0x94427da9.
//
// Solidity: function setBRRData(uint256 _burnBps, uint256 _rewardBps, uint256 _rebateBps) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) SetBRRData(opts *bind.TransactOpts, _burnBps *big.Int, _rewardBps *big.Int, _rebateBps *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "setBRRData", _burnBps, _rewardBps, _rebateBps)
}

// SetBRRData is a paid mutator transaction binding the contract method 0x94427da9.
//
// Solidity: function setBRRData(uint256 _burnBps, uint256 _rewardBps, uint256 _rebateBps) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) SetBRRData(_burnBps *big.Int, _rewardBps *big.Int, _rebateBps *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetBRRData(&_KyberFeeHandlerV2.TransactOpts, _burnBps, _rewardBps, _rebateBps)
}

// SetBRRData is a paid mutator transaction binding the contract method 0x94427da9.
//
// Solidity: function setBRRData(uint256 _burnBps, uint256 _rewardBps, uint256 _rebateBps) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) SetBRRData(_burnBps *big.Int, _rewardBps *big.Int, _rebateBps *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetBRRData(&_KyberFeeHandlerV2.TransactOpts, _burnBps, _rewardBps, _rebateBps)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) SetBurnConfigParams(opts *bind.TransactOpts, _sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "setBurnConfigParams", _sanityRate, _weiToBurn)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) SetBurnConfigParams(_sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetBurnConfigParams(&_KyberFeeHandlerV2.TransactOpts, _sanityRate, _weiToBurn)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) SetBurnConfigParams(_sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetBurnConfigParams(&_KyberFeeHandlerV2.TransactOpts, _sanityRate, _weiToBurn)
}

// SetDaoOperator is a paid mutator transaction binding the contract method 0x64354d65.
//
// Solidity: function setDaoOperator(address _daoOperator) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) SetDaoOperator(opts *bind.TransactOpts, _daoOperator common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "setDaoOperator", _daoOperator)
}

// SetDaoOperator is a paid mutator transaction binding the contract method 0x64354d65.
//
// Solidity: function setDaoOperator(address _daoOperator) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) SetDaoOperator(_daoOperator common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetDaoOperator(&_KyberFeeHandlerV2.TransactOpts, _daoOperator)
}

// SetDaoOperator is a paid mutator transaction binding the contract method 0x64354d65.
//
// Solidity: function setDaoOperator(address _daoOperator) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) SetDaoOperator(_daoOperator common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetDaoOperator(&_KyberFeeHandlerV2.TransactOpts, _daoOperator)
}

// SetFeePool is a paid mutator transaction binding the contract method 0x19db2228.
//
// Solidity: function setFeePool(address _feePool) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) SetFeePool(opts *bind.TransactOpts, _feePool common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "setFeePool", _feePool)
}

// SetFeePool is a paid mutator transaction binding the contract method 0x19db2228.
//
// Solidity: function setFeePool(address _feePool) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) SetFeePool(_feePool common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetFeePool(&_KyberFeeHandlerV2.TransactOpts, _feePool)
}

// SetFeePool is a paid mutator transaction binding the contract method 0x19db2228.
//
// Solidity: function setFeePool(address _feePool) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) SetFeePool(_feePool common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetFeePool(&_KyberFeeHandlerV2.TransactOpts, _feePool)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) SetKyberProxy(opts *bind.TransactOpts, _newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.Transact(opts, "setKyberProxy", _newProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) SetKyberProxy(_newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetKyberProxy(&_KyberFeeHandlerV2.TransactOpts, _newProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) SetKyberProxy(_newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.SetKyberProxy(&_KyberFeeHandlerV2.TransactOpts, _newProxy)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandlerV2.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Session) Receive() (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.Receive(&_KyberFeeHandlerV2.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2TransactorSession) Receive() (*types.Transaction, error) {
	return _KyberFeeHandlerV2.Contract.Receive(&_KyberFeeHandlerV2.TransactOpts)
}

// KyberFeeHandlerV2BRRUpdatedIterator is returned from FilterBRRUpdated and is used to iterate over the raw logs and unpacked data for BRRUpdated events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2BRRUpdatedIterator struct {
	Event *KyberFeeHandlerV2BRRUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2BRRUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2BRRUpdated)
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
		it.Event = new(KyberFeeHandlerV2BRRUpdated)
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
func (it *KyberFeeHandlerV2BRRUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2BRRUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2BRRUpdated represents a BRRUpdated event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2BRRUpdated struct {
	RewardBps *big.Int
	RebateBps *big.Int
	BurnBps   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBRRUpdated is a free log retrieval operation binding the contract event 0x7806a23da6d7c83c8caf5e3fc8ec22c30900a5fcf10e266ac158ad3c9e3384c4.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterBRRUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerV2BRRUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "BRRUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2BRRUpdatedIterator{contract: _KyberFeeHandlerV2.contract, event: "BRRUpdated", logs: logs, sub: sub}, nil
}

// WatchBRRUpdated is a free log subscription operation binding the contract event 0x7806a23da6d7c83c8caf5e3fc8ec22c30900a5fcf10e266ac158ad3c9e3384c4.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchBRRUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2BRRUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "BRRUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2BRRUpdated)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "BRRUpdated", log); err != nil {
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

// ParseBRRUpdated is a log parse operation binding the contract event 0x7806a23da6d7c83c8caf5e3fc8ec22c30900a5fcf10e266ac158ad3c9e3384c4.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseBRRUpdated(log types.Log) (*KyberFeeHandlerV2BRRUpdated, error) {
	event := new(KyberFeeHandlerV2BRRUpdated)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "BRRUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2BurnConfigSetIterator is returned from FilterBurnConfigSet and is used to iterate over the raw logs and unpacked data for BurnConfigSet events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2BurnConfigSetIterator struct {
	Event *KyberFeeHandlerV2BurnConfigSet // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2BurnConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2BurnConfigSet)
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
		it.Event = new(KyberFeeHandlerV2BurnConfigSet)
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
func (it *KyberFeeHandlerV2BurnConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2BurnConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2BurnConfigSet represents a BurnConfigSet event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2BurnConfigSet struct {
	SanityRate common.Address
	WeiToBurn  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBurnConfigSet is a free log retrieval operation binding the contract event 0xe40f97f23269c4682610e9b2522d6d4272ee56f115906d71fcb3da82a860f755.
//
// Solidity: event BurnConfigSet(address sanityRate, uint256 weiToBurn)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterBurnConfigSet(opts *bind.FilterOpts) (*KyberFeeHandlerV2BurnConfigSetIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "BurnConfigSet")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2BurnConfigSetIterator{contract: _KyberFeeHandlerV2.contract, event: "BurnConfigSet", logs: logs, sub: sub}, nil
}

// WatchBurnConfigSet is a free log subscription operation binding the contract event 0xe40f97f23269c4682610e9b2522d6d4272ee56f115906d71fcb3da82a860f755.
//
// Solidity: event BurnConfigSet(address sanityRate, uint256 weiToBurn)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchBurnConfigSet(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2BurnConfigSet) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "BurnConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2BurnConfigSet)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "BurnConfigSet", log); err != nil {
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

// ParseBurnConfigSet is a log parse operation binding the contract event 0xe40f97f23269c4682610e9b2522d6d4272ee56f115906d71fcb3da82a860f755.
//
// Solidity: event BurnConfigSet(address sanityRate, uint256 weiToBurn)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseBurnConfigSet(log types.Log) (*KyberFeeHandlerV2BurnConfigSet, error) {
	event := new(KyberFeeHandlerV2BurnConfigSet)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "BurnConfigSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2DaoOperatorUpdatedIterator is returned from FilterDaoOperatorUpdated and is used to iterate over the raw logs and unpacked data for DaoOperatorUpdated events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2DaoOperatorUpdatedIterator struct {
	Event *KyberFeeHandlerV2DaoOperatorUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2DaoOperatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2DaoOperatorUpdated)
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
		it.Event = new(KyberFeeHandlerV2DaoOperatorUpdated)
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
func (it *KyberFeeHandlerV2DaoOperatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2DaoOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2DaoOperatorUpdated represents a DaoOperatorUpdated event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2DaoOperatorUpdated struct {
	DaoOperator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDaoOperatorUpdated is a free log retrieval operation binding the contract event 0xadf9b628cb2a4e665382961f42205fce0577c0bb2c0e31ef9f87f4a35033c691.
//
// Solidity: event DaoOperatorUpdated(address daoOperator)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterDaoOperatorUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerV2DaoOperatorUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "DaoOperatorUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2DaoOperatorUpdatedIterator{contract: _KyberFeeHandlerV2.contract, event: "DaoOperatorUpdated", logs: logs, sub: sub}, nil
}

// WatchDaoOperatorUpdated is a free log subscription operation binding the contract event 0xadf9b628cb2a4e665382961f42205fce0577c0bb2c0e31ef9f87f4a35033c691.
//
// Solidity: event DaoOperatorUpdated(address daoOperator)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchDaoOperatorUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2DaoOperatorUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "DaoOperatorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2DaoOperatorUpdated)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "DaoOperatorUpdated", log); err != nil {
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

// ParseDaoOperatorUpdated is a log parse operation binding the contract event 0xadf9b628cb2a4e665382961f42205fce0577c0bb2c0e31ef9f87f4a35033c691.
//
// Solidity: event DaoOperatorUpdated(address daoOperator)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseDaoOperatorUpdated(log types.Log) (*KyberFeeHandlerV2DaoOperatorUpdated, error) {
	event := new(KyberFeeHandlerV2DaoOperatorUpdated)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "DaoOperatorUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2EthReceivedIterator is returned from FilterEthReceived and is used to iterate over the raw logs and unpacked data for EthReceived events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2EthReceivedIterator struct {
	Event *KyberFeeHandlerV2EthReceived // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2EthReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2EthReceived)
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
		it.Event = new(KyberFeeHandlerV2EthReceived)
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
func (it *KyberFeeHandlerV2EthReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2EthReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2EthReceived represents a EthReceived event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2EthReceived struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEthReceived is a free log retrieval operation binding the contract event 0x353bcaaf167a6add95a753d39727e3d3beb865129a69a10ed774b0b899671403.
//
// Solidity: event EthReceived(uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterEthReceived(opts *bind.FilterOpts) (*KyberFeeHandlerV2EthReceivedIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "EthReceived")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2EthReceivedIterator{contract: _KyberFeeHandlerV2.contract, event: "EthReceived", logs: logs, sub: sub}, nil
}

// WatchEthReceived is a free log subscription operation binding the contract event 0x353bcaaf167a6add95a753d39727e3d3beb865129a69a10ed774b0b899671403.
//
// Solidity: event EthReceived(uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchEthReceived(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2EthReceived) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "EthReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2EthReceived)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "EthReceived", log); err != nil {
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

// ParseEthReceived is a log parse operation binding the contract event 0x353bcaaf167a6add95a753d39727e3d3beb865129a69a10ed774b0b899671403.
//
// Solidity: event EthReceived(uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseEthReceived(log types.Log) (*KyberFeeHandlerV2EthReceived, error) {
	event := new(KyberFeeHandlerV2EthReceived)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "EthReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2FeeDistributedIterator is returned from FilterFeeDistributed and is used to iterate over the raw logs and unpacked data for FeeDistributed events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2FeeDistributedIterator struct {
	Event *KyberFeeHandlerV2FeeDistributed // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2FeeDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2FeeDistributed)
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
		it.Event = new(KyberFeeHandlerV2FeeDistributed)
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
func (it *KyberFeeHandlerV2FeeDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2FeeDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2FeeDistributed represents a FeeDistributed event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2FeeDistributed struct {
	Token                     common.Address
	Sender                    common.Address
	PlatformWallet            common.Address
	PlatformFeeWei            *big.Int
	RewardWei                 *big.Int
	RebateWei                 *big.Int
	RebateWallets             []common.Address
	RebatePercentBpsPerWallet []*big.Int
	BurnAmtWei                *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterFeeDistributed is a free log retrieval operation binding the contract event 0xc207a63c18c4070ce1e33e5fcc02efb09ac984caa6a2046e2b1d2811723846f1.
//
// Solidity: event FeeDistributed(address indexed token, address indexed sender, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterFeeDistributed(opts *bind.FilterOpts, token []common.Address, sender []common.Address, platformWallet []common.Address) (*KyberFeeHandlerV2FeeDistributedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "FeeDistributed", tokenRule, senderRule, platformWalletRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2FeeDistributedIterator{contract: _KyberFeeHandlerV2.contract, event: "FeeDistributed", logs: logs, sub: sub}, nil
}

// WatchFeeDistributed is a free log subscription operation binding the contract event 0xc207a63c18c4070ce1e33e5fcc02efb09ac984caa6a2046e2b1d2811723846f1.
//
// Solidity: event FeeDistributed(address indexed token, address indexed sender, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchFeeDistributed(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2FeeDistributed, token []common.Address, sender []common.Address, platformWallet []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "FeeDistributed", tokenRule, senderRule, platformWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2FeeDistributed)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "FeeDistributed", log); err != nil {
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

// ParseFeeDistributed is a log parse operation binding the contract event 0xc207a63c18c4070ce1e33e5fcc02efb09ac984caa6a2046e2b1d2811723846f1.
//
// Solidity: event FeeDistributed(address indexed token, address indexed sender, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseFeeDistributed(log types.Log) (*KyberFeeHandlerV2FeeDistributed, error) {
	event := new(KyberFeeHandlerV2FeeDistributed)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "FeeDistributed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2FeePoolUpdatedIterator is returned from FilterFeePoolUpdated and is used to iterate over the raw logs and unpacked data for FeePoolUpdated events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2FeePoolUpdatedIterator struct {
	Event *KyberFeeHandlerV2FeePoolUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2FeePoolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2FeePoolUpdated)
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
		it.Event = new(KyberFeeHandlerV2FeePoolUpdated)
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
func (it *KyberFeeHandlerV2FeePoolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2FeePoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2FeePoolUpdated represents a FeePoolUpdated event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2FeePoolUpdated struct {
	FeePool common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFeePoolUpdated is a free log retrieval operation binding the contract event 0x6d1d088acfe4f30d6014f6f693c61c16258f9784a6ed8439b2c59213eecb6295.
//
// Solidity: event FeePoolUpdated(address feePool)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterFeePoolUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerV2FeePoolUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "FeePoolUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2FeePoolUpdatedIterator{contract: _KyberFeeHandlerV2.contract, event: "FeePoolUpdated", logs: logs, sub: sub}, nil
}

// WatchFeePoolUpdated is a free log subscription operation binding the contract event 0x6d1d088acfe4f30d6014f6f693c61c16258f9784a6ed8439b2c59213eecb6295.
//
// Solidity: event FeePoolUpdated(address feePool)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchFeePoolUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2FeePoolUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "FeePoolUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2FeePoolUpdated)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "FeePoolUpdated", log); err != nil {
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

// ParseFeePoolUpdated is a log parse operation binding the contract event 0x6d1d088acfe4f30d6014f6f693c61c16258f9784a6ed8439b2c59213eecb6295.
//
// Solidity: event FeePoolUpdated(address feePool)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseFeePoolUpdated(log types.Log) (*KyberFeeHandlerV2FeePoolUpdated, error) {
	event := new(KyberFeeHandlerV2FeePoolUpdated)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "FeePoolUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2KncBurnedIterator is returned from FilterKncBurned and is used to iterate over the raw logs and unpacked data for KncBurned events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2KncBurnedIterator struct {
	Event *KyberFeeHandlerV2KncBurned // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2KncBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2KncBurned)
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
		it.Event = new(KyberFeeHandlerV2KncBurned)
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
func (it *KyberFeeHandlerV2KncBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2KncBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2KncBurned represents a KncBurned event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2KncBurned struct {
	KncTWei *big.Int
	Token   common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKncBurned is a free log retrieval operation binding the contract event 0xa0fcef56e2b45fcbeb91d5e629ef6b2b6e982d0768f02d1232610315cd23ea10.
//
// Solidity: event KncBurned(uint256 kncTWei, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterKncBurned(opts *bind.FilterOpts, token []common.Address) (*KyberFeeHandlerV2KncBurnedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "KncBurned", tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2KncBurnedIterator{contract: _KyberFeeHandlerV2.contract, event: "KncBurned", logs: logs, sub: sub}, nil
}

// WatchKncBurned is a free log subscription operation binding the contract event 0xa0fcef56e2b45fcbeb91d5e629ef6b2b6e982d0768f02d1232610315cd23ea10.
//
// Solidity: event KncBurned(uint256 kncTWei, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchKncBurned(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2KncBurned, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "KncBurned", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2KncBurned)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "KncBurned", log); err != nil {
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

// ParseKncBurned is a log parse operation binding the contract event 0xa0fcef56e2b45fcbeb91d5e629ef6b2b6e982d0768f02d1232610315cd23ea10.
//
// Solidity: event KncBurned(uint256 kncTWei, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseKncBurned(log types.Log) (*KyberFeeHandlerV2KncBurned, error) {
	event := new(KyberFeeHandlerV2KncBurned)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "KncBurned", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2KyberProxyUpdatedIterator is returned from FilterKyberProxyUpdated and is used to iterate over the raw logs and unpacked data for KyberProxyUpdated events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2KyberProxyUpdatedIterator struct {
	Event *KyberFeeHandlerV2KyberProxyUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2KyberProxyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2KyberProxyUpdated)
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
		it.Event = new(KyberFeeHandlerV2KyberProxyUpdated)
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
func (it *KyberFeeHandlerV2KyberProxyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2KyberProxyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2KyberProxyUpdated represents a KyberProxyUpdated event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2KyberProxyUpdated struct {
	KyberProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterKyberProxyUpdated is a free log retrieval operation binding the contract event 0x8457f9bd0d13488a6c265af376d291f3c6bd2311d9e8dee5671d4169ca6e0ae0.
//
// Solidity: event KyberProxyUpdated(address kyberProxy)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterKyberProxyUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerV2KyberProxyUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "KyberProxyUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2KyberProxyUpdatedIterator{contract: _KyberFeeHandlerV2.contract, event: "KyberProxyUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberProxyUpdated is a free log subscription operation binding the contract event 0x8457f9bd0d13488a6c265af376d291f3c6bd2311d9e8dee5671d4169ca6e0ae0.
//
// Solidity: event KyberProxyUpdated(address kyberProxy)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchKyberProxyUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2KyberProxyUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "KyberProxyUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2KyberProxyUpdated)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "KyberProxyUpdated", log); err != nil {
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

// ParseKyberProxyUpdated is a log parse operation binding the contract event 0x8457f9bd0d13488a6c265af376d291f3c6bd2311d9e8dee5671d4169ca6e0ae0.
//
// Solidity: event KyberProxyUpdated(address kyberProxy)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseKyberProxyUpdated(log types.Log) (*KyberFeeHandlerV2KyberProxyUpdated, error) {
	event := new(KyberFeeHandlerV2KyberProxyUpdated)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "KyberProxyUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2PlatformFeePaidIterator is returned from FilterPlatformFeePaid and is used to iterate over the raw logs and unpacked data for PlatformFeePaid events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2PlatformFeePaidIterator struct {
	Event *KyberFeeHandlerV2PlatformFeePaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2PlatformFeePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2PlatformFeePaid)
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
		it.Event = new(KyberFeeHandlerV2PlatformFeePaid)
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
func (it *KyberFeeHandlerV2PlatformFeePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2PlatformFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2PlatformFeePaid represents a PlatformFeePaid event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2PlatformFeePaid struct {
	PlatformWallet common.Address
	Token          common.Address
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPlatformFeePaid is a free log retrieval operation binding the contract event 0xebe3db09f5650582b4782506e0d272262129183570e55fcf8768dd6e91f8c0f6.
//
// Solidity: event PlatformFeePaid(address indexed platformWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterPlatformFeePaid(opts *bind.FilterOpts, platformWallet []common.Address, token []common.Address) (*KyberFeeHandlerV2PlatformFeePaidIterator, error) {

	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "PlatformFeePaid", platformWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2PlatformFeePaidIterator{contract: _KyberFeeHandlerV2.contract, event: "PlatformFeePaid", logs: logs, sub: sub}, nil
}

// WatchPlatformFeePaid is a free log subscription operation binding the contract event 0xebe3db09f5650582b4782506e0d272262129183570e55fcf8768dd6e91f8c0f6.
//
// Solidity: event PlatformFeePaid(address indexed platformWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchPlatformFeePaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2PlatformFeePaid, platformWallet []common.Address, token []common.Address) (event.Subscription, error) {

	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "PlatformFeePaid", platformWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2PlatformFeePaid)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "PlatformFeePaid", log); err != nil {
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

// ParsePlatformFeePaid is a log parse operation binding the contract event 0xebe3db09f5650582b4782506e0d272262129183570e55fcf8768dd6e91f8c0f6.
//
// Solidity: event PlatformFeePaid(address indexed platformWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParsePlatformFeePaid(log types.Log) (*KyberFeeHandlerV2PlatformFeePaid, error) {
	event := new(KyberFeeHandlerV2PlatformFeePaid)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "PlatformFeePaid", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2RebatePaidIterator is returned from FilterRebatePaid and is used to iterate over the raw logs and unpacked data for RebatePaid events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2RebatePaidIterator struct {
	Event *KyberFeeHandlerV2RebatePaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2RebatePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2RebatePaid)
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
		it.Event = new(KyberFeeHandlerV2RebatePaid)
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
func (it *KyberFeeHandlerV2RebatePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2RebatePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2RebatePaid represents a RebatePaid event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2RebatePaid struct {
	RebateWallet common.Address
	Token        common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRebatePaid is a free log retrieval operation binding the contract event 0xb5ec5e03662403108373ab6431d3e834cb1011fca164541aef315fc7dea7b3b6.
//
// Solidity: event RebatePaid(address indexed rebateWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterRebatePaid(opts *bind.FilterOpts, rebateWallet []common.Address, token []common.Address) (*KyberFeeHandlerV2RebatePaidIterator, error) {

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "RebatePaid", rebateWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2RebatePaidIterator{contract: _KyberFeeHandlerV2.contract, event: "RebatePaid", logs: logs, sub: sub}, nil
}

// WatchRebatePaid is a free log subscription operation binding the contract event 0xb5ec5e03662403108373ab6431d3e834cb1011fca164541aef315fc7dea7b3b6.
//
// Solidity: event RebatePaid(address indexed rebateWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchRebatePaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2RebatePaid, rebateWallet []common.Address, token []common.Address) (event.Subscription, error) {

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "RebatePaid", rebateWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2RebatePaid)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "RebatePaid", log); err != nil {
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

// ParseRebatePaid is a log parse operation binding the contract event 0xb5ec5e03662403108373ab6431d3e834cb1011fca164541aef315fc7dea7b3b6.
//
// Solidity: event RebatePaid(address indexed rebateWallet, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseRebatePaid(log types.Log) (*KyberFeeHandlerV2RebatePaid, error) {
	event := new(KyberFeeHandlerV2RebatePaid)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "RebatePaid", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerV2RewardPaidIterator is returned from FilterRewardPaid and is used to iterate over the raw logs and unpacked data for RewardPaid events raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2RewardPaidIterator struct {
	Event *KyberFeeHandlerV2RewardPaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerV2RewardPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerV2RewardPaid)
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
		it.Event = new(KyberFeeHandlerV2RewardPaid)
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
func (it *KyberFeeHandlerV2RewardPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerV2RewardPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerV2RewardPaid represents a RewardPaid event raised by the KyberFeeHandlerV2 contract.
type KyberFeeHandlerV2RewardPaid struct {
	Staker common.Address
	Epoch  *big.Int
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardPaid is a free log retrieval operation binding the contract event 0xaf206e736916d38b56e2d559931a189bc3119b8fc6d6850bd34e382f09030587.
//
// Solidity: event RewardPaid(address indexed staker, uint256 indexed epoch, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) FilterRewardPaid(opts *bind.FilterOpts, staker []common.Address, epoch []*big.Int, token []common.Address) (*KyberFeeHandlerV2RewardPaidIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.FilterLogs(opts, "RewardPaid", stakerRule, epochRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerV2RewardPaidIterator{contract: _KyberFeeHandlerV2.contract, event: "RewardPaid", logs: logs, sub: sub}, nil
}

// WatchRewardPaid is a free log subscription operation binding the contract event 0xaf206e736916d38b56e2d559931a189bc3119b8fc6d6850bd34e382f09030587.
//
// Solidity: event RewardPaid(address indexed staker, uint256 indexed epoch, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) WatchRewardPaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerV2RewardPaid, staker []common.Address, epoch []*big.Int, token []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandlerV2.contract.WatchLogs(opts, "RewardPaid", stakerRule, epochRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerV2RewardPaid)
				if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "RewardPaid", log); err != nil {
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

// ParseRewardPaid is a log parse operation binding the contract event 0xaf206e736916d38b56e2d559931a189bc3119b8fc6d6850bd34e382f09030587.
//
// Solidity: event RewardPaid(address indexed staker, uint256 indexed epoch, address indexed token, uint256 amount)
func (_KyberFeeHandlerV2 *KyberFeeHandlerV2Filterer) ParseRewardPaid(log types.Log) (*KyberFeeHandlerV2RewardPaid, error) {
	event := new(KyberFeeHandlerV2RewardPaid)
	if err := _KyberFeeHandlerV2.contract.UnpackLog(event, "RewardPaid", log); err != nil {
		return nil, err
	}
	return event, nil
}
