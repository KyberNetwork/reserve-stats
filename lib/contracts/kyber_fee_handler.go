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

// KyberFeeHandlerABI is the input ABI used to generate the binding from.
const KyberFeeHandlerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_daoSetter\",\"type\":\"address\"},{\"internalType\":\"contractIKyberProxy\",\"name\":\"_kyberProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_kyberNetwork\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_knc\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_burnBlockInterval\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_daoOperator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebateBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiryTimestamp\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"BRRUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractISanityRate\",\"name\":\"sanityRate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiToBurn\",\"type\":\"uint256\"}],\"name\":\"BurnConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"platformFeeWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebateWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"rebateWallets\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"rebatePercentBpsPerWallet\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmtWei\",\"type\":\"uint256\"}],\"name\":\"FeeDistributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"kncTWei\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"KncBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberDao\",\"name\":\"kyberDao\",\"type\":\"address\"}],\"name\":\"KyberDaoAddressSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"kyberNetwork\",\"type\":\"address\"}],\"name\":\"KyberNetworkUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberProxy\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"KyberProxyUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PlatformFeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RebatePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardsWei\",\"type\":\"uint256\"}],\"name\":\"RewardsRemovedToBurn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"brrAndEpochData\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"expiryTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"rewardBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"rebateBps\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnBlockInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnKnc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"kncBurnAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"}],\"name\":\"claimPlatformFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountWei\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"}],\"name\":\"claimReserveRebate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountWei\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"claimStakerReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountWei\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daoOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daoSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"feePerPlatformWallet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rewardBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rebateBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestSanityRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"kncToEthSanityRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSanityRateContracts\",\"outputs\":[{\"internalType\":\"contractISanityRate[]\",\"name\":\"sanityRates\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"rebateWallets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"rebateBpsPerWallet\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"platformFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"networkFee\",\"type\":\"uint256\"}],\"name\":\"handleFees\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hasClaimedReward\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"knc\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberDao\",\"outputs\":[{\"internalType\":\"contractIKyberDao\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberNetwork\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberProxy\",\"outputs\":[{\"internalType\":\"contractIKyberProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBurnBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"makeEpochRewardBurnable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"readBRRData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rewardBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rebateBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"rebatePerWallet\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rewardsPaidPerEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rewardsPerEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISanityRate\",\"name\":\"_sanityRate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_weiToBurn\",\"type\":\"uint256\"}],\"name\":\"setBurnConfigParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberDao\",\"name\":\"_kyberDao\",\"type\":\"address\"}],\"name\":\"setDaoContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberProxy\",\"name\":\"_newProxy\",\"type\":\"address\"}],\"name\":\"setKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_kyberNetwork\",\"type\":\"address\"}],\"name\":\"setNetworkContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPayoutBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"weiToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// KyberFeeHandler is an auto generated Go binding around an Ethereum contract.
type KyberFeeHandler struct {
	KyberFeeHandlerCaller     // Read-only binding to the contract
	KyberFeeHandlerTransactor // Write-only binding to the contract
	KyberFeeHandlerFilterer   // Log filterer for contract events
}

// KyberFeeHandlerCaller is an auto generated read-only Go binding around an Ethereum contract.
type KyberFeeHandlerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KyberFeeHandlerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KyberFeeHandlerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberFeeHandlerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KyberFeeHandlerSession struct {
	Contract     *KyberFeeHandler  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KyberFeeHandlerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KyberFeeHandlerCallerSession struct {
	Contract *KyberFeeHandlerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// KyberFeeHandlerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KyberFeeHandlerTransactorSession struct {
	Contract     *KyberFeeHandlerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// KyberFeeHandlerRaw is an auto generated low-level Go binding around an Ethereum contract.
type KyberFeeHandlerRaw struct {
	Contract *KyberFeeHandler // Generic contract binding to access the raw methods on
}

// KyberFeeHandlerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KyberFeeHandlerCallerRaw struct {
	Contract *KyberFeeHandlerCaller // Generic read-only contract binding to access the raw methods on
}

// KyberFeeHandlerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KyberFeeHandlerTransactorRaw struct {
	Contract *KyberFeeHandlerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKyberFeeHandler creates a new instance of KyberFeeHandler, bound to a specific deployed contract.
func NewKyberFeeHandler(address common.Address, backend bind.ContractBackend) (*KyberFeeHandler, error) {
	contract, err := bindKyberFeeHandler(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandler{KyberFeeHandlerCaller: KyberFeeHandlerCaller{contract: contract}, KyberFeeHandlerTransactor: KyberFeeHandlerTransactor{contract: contract}, KyberFeeHandlerFilterer: KyberFeeHandlerFilterer{contract: contract}}, nil
}

// NewKyberFeeHandlerCaller creates a new read-only instance of KyberFeeHandler, bound to a specific deployed contract.
func NewKyberFeeHandlerCaller(address common.Address, caller bind.ContractCaller) (*KyberFeeHandlerCaller, error) {
	contract, err := bindKyberFeeHandler(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerCaller{contract: contract}, nil
}

// NewKyberFeeHandlerTransactor creates a new write-only instance of KyberFeeHandler, bound to a specific deployed contract.
func NewKyberFeeHandlerTransactor(address common.Address, transactor bind.ContractTransactor) (*KyberFeeHandlerTransactor, error) {
	contract, err := bindKyberFeeHandler(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerTransactor{contract: contract}, nil
}

// NewKyberFeeHandlerFilterer creates a new log filterer instance of KyberFeeHandler, bound to a specific deployed contract.
func NewKyberFeeHandlerFilterer(address common.Address, filterer bind.ContractFilterer) (*KyberFeeHandlerFilterer, error) {
	contract, err := bindKyberFeeHandler(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerFilterer{contract: contract}, nil
}

// bindKyberFeeHandler binds a generic wrapper to an already deployed contract.
func bindKyberFeeHandler(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KyberFeeHandlerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberFeeHandler *KyberFeeHandlerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberFeeHandler.Contract.KyberFeeHandlerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberFeeHandler *KyberFeeHandlerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.KyberFeeHandlerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberFeeHandler *KyberFeeHandlerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.KyberFeeHandlerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberFeeHandler *KyberFeeHandlerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberFeeHandler.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberFeeHandler *KyberFeeHandlerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberFeeHandler *KyberFeeHandlerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.contract.Transact(opts, method, params...)
}

// BrrAndEpochData is a free data retrieval call binding the contract method 0xf392e218.
//
// Solidity: function brrAndEpochData() view returns(uint64 expiryTimestamp, uint32 epoch, uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandler *KyberFeeHandlerCaller) BrrAndEpochData(opts *bind.CallOpts) (struct {
	ExpiryTimestamp uint64
	Epoch           uint32
	RewardBps       uint16
	RebateBps       uint16
}, error) {
	ret := new(struct {
		ExpiryTimestamp uint64
		Epoch           uint32
		RewardBps       uint16
		RebateBps       uint16
	})
	out := ret
	err := _KyberFeeHandler.contract.Call(opts, out, "brrAndEpochData")
	return *ret, err
}

// BrrAndEpochData is a free data retrieval call binding the contract method 0xf392e218.
//
// Solidity: function brrAndEpochData() view returns(uint64 expiryTimestamp, uint32 epoch, uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandler *KyberFeeHandlerSession) BrrAndEpochData() (struct {
	ExpiryTimestamp uint64
	Epoch           uint32
	RewardBps       uint16
	RebateBps       uint16
}, error) {
	return _KyberFeeHandler.Contract.BrrAndEpochData(&_KyberFeeHandler.CallOpts)
}

// BrrAndEpochData is a free data retrieval call binding the contract method 0xf392e218.
//
// Solidity: function brrAndEpochData() view returns(uint64 expiryTimestamp, uint32 epoch, uint16 rewardBps, uint16 rebateBps)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) BrrAndEpochData() (struct {
	ExpiryTimestamp uint64
	Epoch           uint32
	RewardBps       uint16
	RebateBps       uint16
}, error) {
	return _KyberFeeHandler.Contract.BrrAndEpochData(&_KyberFeeHandler.CallOpts)
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) BurnBlockInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "burnBlockInterval")
	return *ret0, err
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) BurnBlockInterval() (*big.Int, error) {
	return _KyberFeeHandler.Contract.BurnBlockInterval(&_KyberFeeHandler.CallOpts)
}

// BurnBlockInterval is a free data retrieval call binding the contract method 0xb45782c7.
//
// Solidity: function burnBlockInterval() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) BurnBlockInterval() (*big.Int, error) {
	return _KyberFeeHandler.Contract.BurnBlockInterval(&_KyberFeeHandler.CallOpts)
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) DaoOperator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "daoOperator")
	return *ret0, err
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) DaoOperator() (common.Address, error) {
	return _KyberFeeHandler.Contract.DaoOperator(&_KyberFeeHandler.CallOpts)
}

// DaoOperator is a free data retrieval call binding the contract method 0x8c9bc208.
//
// Solidity: function daoOperator() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) DaoOperator() (common.Address, error) {
	return _KyberFeeHandler.Contract.DaoOperator(&_KyberFeeHandler.CallOpts)
}

// DaoSetter is a free data retrieval call binding the contract method 0xb6981e2c.
//
// Solidity: function daoSetter() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) DaoSetter(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "daoSetter")
	return *ret0, err
}

// DaoSetter is a free data retrieval call binding the contract method 0xb6981e2c.
//
// Solidity: function daoSetter() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) DaoSetter() (common.Address, error) {
	return _KyberFeeHandler.Contract.DaoSetter(&_KyberFeeHandler.CallOpts)
}

// DaoSetter is a free data retrieval call binding the contract method 0xb6981e2c.
//
// Solidity: function daoSetter() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) DaoSetter() (common.Address, error) {
	return _KyberFeeHandler.Contract.DaoSetter(&_KyberFeeHandler.CallOpts)
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) FeePerPlatformWallet(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "feePerPlatformWallet", arg0)
	return *ret0, err
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) FeePerPlatformWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandler.Contract.FeePerPlatformWallet(&_KyberFeeHandler.CallOpts, arg0)
}

// FeePerPlatformWallet is a free data retrieval call binding the contract method 0x03339513.
//
// Solidity: function feePerPlatformWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) FeePerPlatformWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandler.Contract.FeePerPlatformWallet(&_KyberFeeHandler.CallOpts, arg0)
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandler *KyberFeeHandlerCaller) GetLatestSanityRate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "getLatestSanityRate")
	return *ret0, err
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandler *KyberFeeHandlerSession) GetLatestSanityRate() (*big.Int, error) {
	return _KyberFeeHandler.Contract.GetLatestSanityRate(&_KyberFeeHandler.CallOpts)
}

// GetLatestSanityRate is a free data retrieval call binding the contract method 0xa840874f.
//
// Solidity: function getLatestSanityRate() view returns(uint256 kncToEthSanityRate)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) GetLatestSanityRate() (*big.Int, error) {
	return _KyberFeeHandler.Contract.GetLatestSanityRate(&_KyberFeeHandler.CallOpts)
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandler *KyberFeeHandlerCaller) GetSanityRateContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "getSanityRateContracts")
	return *ret0, err
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandler *KyberFeeHandlerSession) GetSanityRateContracts() ([]common.Address, error) {
	return _KyberFeeHandler.Contract.GetSanityRateContracts(&_KyberFeeHandler.CallOpts)
}

// GetSanityRateContracts is a free data retrieval call binding the contract method 0x66ab3fe6.
//
// Solidity: function getSanityRateContracts() view returns(address[] sanityRates)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) GetSanityRateContracts() ([]common.Address, error) {
	return _KyberFeeHandler.Contract.GetSanityRateContracts(&_KyberFeeHandler.CallOpts)
}

// HasClaimedReward is a free data retrieval call binding the contract method 0x7c360101.
//
// Solidity: function hasClaimedReward(address , uint256 ) view returns(bool)
func (_KyberFeeHandler *KyberFeeHandlerCaller) HasClaimedReward(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "hasClaimedReward", arg0, arg1)
	return *ret0, err
}

// HasClaimedReward is a free data retrieval call binding the contract method 0x7c360101.
//
// Solidity: function hasClaimedReward(address , uint256 ) view returns(bool)
func (_KyberFeeHandler *KyberFeeHandlerSession) HasClaimedReward(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _KyberFeeHandler.Contract.HasClaimedReward(&_KyberFeeHandler.CallOpts, arg0, arg1)
}

// HasClaimedReward is a free data retrieval call binding the contract method 0x7c360101.
//
// Solidity: function hasClaimedReward(address , uint256 ) view returns(bool)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) HasClaimedReward(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _KyberFeeHandler.Contract.HasClaimedReward(&_KyberFeeHandler.CallOpts, arg0, arg1)
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) Knc(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "knc")
	return *ret0, err
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) Knc() (common.Address, error) {
	return _KyberFeeHandler.Contract.Knc(&_KyberFeeHandler.CallOpts)
}

// Knc is a free data retrieval call binding the contract method 0xe61387e0.
//
// Solidity: function knc() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) Knc() (common.Address, error) {
	return _KyberFeeHandler.Contract.Knc(&_KyberFeeHandler.CallOpts)
}

// KyberDao is a free data retrieval call binding the contract method 0x4d8f5105.
//
// Solidity: function kyberDao() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) KyberDao(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "kyberDao")
	return *ret0, err
}

// KyberDao is a free data retrieval call binding the contract method 0x4d8f5105.
//
// Solidity: function kyberDao() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) KyberDao() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberDao(&_KyberFeeHandler.CallOpts)
}

// KyberDao is a free data retrieval call binding the contract method 0x4d8f5105.
//
// Solidity: function kyberDao() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) KyberDao() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberDao(&_KyberFeeHandler.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) KyberNetwork(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "kyberNetwork")
	return *ret0, err
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) KyberNetwork() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberNetwork(&_KyberFeeHandler.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) KyberNetwork() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberNetwork(&_KyberFeeHandler.CallOpts)
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCaller) KyberProxy(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "kyberProxy")
	return *ret0, err
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerSession) KyberProxy() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberProxy(&_KyberFeeHandler.CallOpts)
}

// KyberProxy is a free data retrieval call binding the contract method 0xf0eeed81.
//
// Solidity: function kyberProxy() view returns(address)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) KyberProxy() (common.Address, error) {
	return _KyberFeeHandler.Contract.KyberProxy(&_KyberFeeHandler.CallOpts)
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) LastBurnBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "lastBurnBlock")
	return *ret0, err
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) LastBurnBlock() (*big.Int, error) {
	return _KyberFeeHandler.Contract.LastBurnBlock(&_KyberFeeHandler.CallOpts)
}

// LastBurnBlock is a free data retrieval call binding the contract method 0xc03e798c.
//
// Solidity: function lastBurnBlock() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) LastBurnBlock() (*big.Int, error) {
	return _KyberFeeHandler.Contract.LastBurnBlock(&_KyberFeeHandler.CallOpts)
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps, uint256 expiryTimestamp, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerCaller) ReadBRRData(opts *bind.CallOpts) (struct {
	RewardBps       *big.Int
	RebateBps       *big.Int
	ExpiryTimestamp *big.Int
	Epoch           *big.Int
}, error) {
	ret := new(struct {
		RewardBps       *big.Int
		RebateBps       *big.Int
		ExpiryTimestamp *big.Int
		Epoch           *big.Int
	})
	out := ret
	err := _KyberFeeHandler.contract.Call(opts, out, "readBRRData")
	return *ret, err
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps, uint256 expiryTimestamp, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerSession) ReadBRRData() (struct {
	RewardBps       *big.Int
	RebateBps       *big.Int
	ExpiryTimestamp *big.Int
	Epoch           *big.Int
}, error) {
	return _KyberFeeHandler.Contract.ReadBRRData(&_KyberFeeHandler.CallOpts)
}

// ReadBRRData is a free data retrieval call binding the contract method 0x770ba561.
//
// Solidity: function readBRRData() view returns(uint256 rewardBps, uint256 rebateBps, uint256 expiryTimestamp, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) ReadBRRData() (struct {
	RewardBps       *big.Int
	RebateBps       *big.Int
	ExpiryTimestamp *big.Int
	Epoch           *big.Int
}, error) {
	return _KyberFeeHandler.Contract.ReadBRRData(&_KyberFeeHandler.CallOpts)
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) RebatePerWallet(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "rebatePerWallet", arg0)
	return *ret0, err
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) RebatePerWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RebatePerWallet(&_KyberFeeHandler.CallOpts, arg0)
}

// RebatePerWallet is a free data retrieval call binding the contract method 0x579d6b74.
//
// Solidity: function rebatePerWallet(address ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) RebatePerWallet(arg0 common.Address) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RebatePerWallet(&_KyberFeeHandler.CallOpts, arg0)
}

// RewardsPaidPerEpoch is a free data retrieval call binding the contract method 0xf7ac3cbc.
//
// Solidity: function rewardsPaidPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) RewardsPaidPerEpoch(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "rewardsPaidPerEpoch", arg0)
	return *ret0, err
}

// RewardsPaidPerEpoch is a free data retrieval call binding the contract method 0xf7ac3cbc.
//
// Solidity: function rewardsPaidPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) RewardsPaidPerEpoch(arg0 *big.Int) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RewardsPaidPerEpoch(&_KyberFeeHandler.CallOpts, arg0)
}

// RewardsPaidPerEpoch is a free data retrieval call binding the contract method 0xf7ac3cbc.
//
// Solidity: function rewardsPaidPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) RewardsPaidPerEpoch(arg0 *big.Int) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RewardsPaidPerEpoch(&_KyberFeeHandler.CallOpts, arg0)
}

// RewardsPerEpoch is a free data retrieval call binding the contract method 0x94cee7b3.
//
// Solidity: function rewardsPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) RewardsPerEpoch(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "rewardsPerEpoch", arg0)
	return *ret0, err
}

// RewardsPerEpoch is a free data retrieval call binding the contract method 0x94cee7b3.
//
// Solidity: function rewardsPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) RewardsPerEpoch(arg0 *big.Int) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RewardsPerEpoch(&_KyberFeeHandler.CallOpts, arg0)
}

// RewardsPerEpoch is a free data retrieval call binding the contract method 0x94cee7b3.
//
// Solidity: function rewardsPerEpoch(uint256 ) view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) RewardsPerEpoch(arg0 *big.Int) (*big.Int, error) {
	return _KyberFeeHandler.Contract.RewardsPerEpoch(&_KyberFeeHandler.CallOpts, arg0)
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) TotalPayoutBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "totalPayoutBalance")
	return *ret0, err
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) TotalPayoutBalance() (*big.Int, error) {
	return _KyberFeeHandler.Contract.TotalPayoutBalance(&_KyberFeeHandler.CallOpts)
}

// TotalPayoutBalance is a free data retrieval call binding the contract method 0x12efe834.
//
// Solidity: function totalPayoutBalance() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) TotalPayoutBalance() (*big.Int, error) {
	return _KyberFeeHandler.Contract.TotalPayoutBalance(&_KyberFeeHandler.CallOpts)
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCaller) WeiToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberFeeHandler.contract.Call(opts, out, "weiToBurn")
	return *ret0, err
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerSession) WeiToBurn() (*big.Int, error) {
	return _KyberFeeHandler.Contract.WeiToBurn(&_KyberFeeHandler.CallOpts)
}

// WeiToBurn is a free data retrieval call binding the contract method 0x80feeef3.
//
// Solidity: function weiToBurn() view returns(uint256)
func (_KyberFeeHandler *KyberFeeHandlerCallerSession) WeiToBurn() (*big.Int, error) {
	return _KyberFeeHandler.Contract.WeiToBurn(&_KyberFeeHandler.CallOpts)
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandler *KyberFeeHandlerTransactor) BurnKnc(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "burnKnc")
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandler *KyberFeeHandlerSession) BurnKnc() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.BurnKnc(&_KyberFeeHandler.TransactOpts)
}

// BurnKnc is a paid mutator transaction binding the contract method 0xa636a8a2.
//
// Solidity: function burnKnc() returns(uint256 kncBurnAmount)
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) BurnKnc() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.BurnKnc(&_KyberFeeHandler.TransactOpts)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactor) ClaimPlatformFee(opts *bind.TransactOpts, platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "claimPlatformFee", platformWallet)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerSession) ClaimPlatformFee(platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimPlatformFee(&_KyberFeeHandler.TransactOpts, platformWallet)
}

// ClaimPlatformFee is a paid mutator transaction binding the contract method 0x9907672a.
//
// Solidity: function claimPlatformFee(address platformWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) ClaimPlatformFee(platformWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimPlatformFee(&_KyberFeeHandler.TransactOpts, platformWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactor) ClaimReserveRebate(opts *bind.TransactOpts, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "claimReserveRebate", rebateWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerSession) ClaimReserveRebate(rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimReserveRebate(&_KyberFeeHandler.TransactOpts, rebateWallet)
}

// ClaimReserveRebate is a paid mutator transaction binding the contract method 0xc01bdf04.
//
// Solidity: function claimReserveRebate(address rebateWallet) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) ClaimReserveRebate(rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimReserveRebate(&_KyberFeeHandler.TransactOpts, rebateWallet)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address staker, uint256 epoch) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactor) ClaimStakerReward(opts *bind.TransactOpts, staker common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "claimStakerReward", staker, epoch)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address staker, uint256 epoch) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerSession) ClaimStakerReward(staker common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimStakerReward(&_KyberFeeHandler.TransactOpts, staker, epoch)
}

// ClaimStakerReward is a paid mutator transaction binding the contract method 0x53fa2eb7.
//
// Solidity: function claimStakerReward(address staker, uint256 epoch) returns(uint256 amountWei)
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) ClaimStakerReward(staker common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.ClaimStakerReward(&_KyberFeeHandler.TransactOpts, staker, epoch)
}

// GetBRR is a paid mutator transaction binding the contract method 0xb3613f11.
//
// Solidity: function getBRR() returns(uint256 rewardBps, uint256 rebateBps, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerTransactor) GetBRR(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "getBRR")
}

// GetBRR is a paid mutator transaction binding the contract method 0xb3613f11.
//
// Solidity: function getBRR() returns(uint256 rewardBps, uint256 rebateBps, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerSession) GetBRR() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.GetBRR(&_KyberFeeHandler.TransactOpts)
}

// GetBRR is a paid mutator transaction binding the contract method 0xb3613f11.
//
// Solidity: function getBRR() returns(uint256 rewardBps, uint256 rebateBps, uint256 epoch)
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) GetBRR() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.GetBRR(&_KyberFeeHandler.TransactOpts)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) HandleFees(opts *bind.TransactOpts, token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "handleFees", token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) HandleFees(token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.HandleFees(&_KyberFeeHandler.TransactOpts, token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// HandleFees is a paid mutator transaction binding the contract method 0xb7c5ab41.
//
// Solidity: function handleFees(address token, address[] rebateWallets, uint256[] rebateBpsPerWallet, address platformWallet, uint256 platformFee, uint256 networkFee) payable returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) HandleFees(token common.Address, rebateWallets []common.Address, rebateBpsPerWallet []*big.Int, platformWallet common.Address, platformFee *big.Int, networkFee *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.HandleFees(&_KyberFeeHandler.TransactOpts, token, rebateWallets, rebateBpsPerWallet, platformWallet, platformFee, networkFee)
}

// MakeEpochRewardBurnable is a paid mutator transaction binding the contract method 0x8bca3efe.
//
// Solidity: function makeEpochRewardBurnable(uint256 epoch) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) MakeEpochRewardBurnable(opts *bind.TransactOpts, epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "makeEpochRewardBurnable", epoch)
}

// MakeEpochRewardBurnable is a paid mutator transaction binding the contract method 0x8bca3efe.
//
// Solidity: function makeEpochRewardBurnable(uint256 epoch) returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) MakeEpochRewardBurnable(epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.MakeEpochRewardBurnable(&_KyberFeeHandler.TransactOpts, epoch)
}

// MakeEpochRewardBurnable is a paid mutator transaction binding the contract method 0x8bca3efe.
//
// Solidity: function makeEpochRewardBurnable(uint256 epoch) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) MakeEpochRewardBurnable(epoch *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.MakeEpochRewardBurnable(&_KyberFeeHandler.TransactOpts, epoch)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) SetBurnConfigParams(opts *bind.TransactOpts, _sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "setBurnConfigParams", _sanityRate, _weiToBurn)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) SetBurnConfigParams(_sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetBurnConfigParams(&_KyberFeeHandler.TransactOpts, _sanityRate, _weiToBurn)
}

// SetBurnConfigParams is a paid mutator transaction binding the contract method 0x692bdfd5.
//
// Solidity: function setBurnConfigParams(address _sanityRate, uint256 _weiToBurn) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) SetBurnConfigParams(_sanityRate common.Address, _weiToBurn *big.Int) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetBurnConfigParams(&_KyberFeeHandler.TransactOpts, _sanityRate, _weiToBurn)
}

// SetDaoContract is a paid mutator transaction binding the contract method 0x8fb58285.
//
// Solidity: function setDaoContract(address _kyberDao) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) SetDaoContract(opts *bind.TransactOpts, _kyberDao common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "setDaoContract", _kyberDao)
}

// SetDaoContract is a paid mutator transaction binding the contract method 0x8fb58285.
//
// Solidity: function setDaoContract(address _kyberDao) returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) SetDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetDaoContract(&_KyberFeeHandler.TransactOpts, _kyberDao)
}

// SetDaoContract is a paid mutator transaction binding the contract method 0x8fb58285.
//
// Solidity: function setDaoContract(address _kyberDao) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) SetDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetDaoContract(&_KyberFeeHandler.TransactOpts, _kyberDao)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) SetKyberProxy(opts *bind.TransactOpts, _newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "setKyberProxy", _newProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) SetKyberProxy(_newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetKyberProxy(&_KyberFeeHandler.TransactOpts, _newProxy)
}

// SetKyberProxy is a paid mutator transaction binding the contract method 0xc6c3f3f9.
//
// Solidity: function setKyberProxy(address _newProxy) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) SetKyberProxy(_newProxy common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetKyberProxy(&_KyberFeeHandler.TransactOpts, _newProxy)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) SetNetworkContract(opts *bind.TransactOpts, _kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.Transact(opts, "setNetworkContract", _kyberNetwork)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) SetNetworkContract(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetNetworkContract(&_KyberFeeHandler.TransactOpts, _kyberNetwork)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) SetNetworkContract(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.SetNetworkContract(&_KyberFeeHandler.TransactOpts, _kyberNetwork)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberFeeHandler.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandler *KyberFeeHandlerSession) Receive() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.Receive(&_KyberFeeHandler.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberFeeHandler *KyberFeeHandlerTransactorSession) Receive() (*types.Transaction, error) {
	return _KyberFeeHandler.Contract.Receive(&_KyberFeeHandler.TransactOpts)
}

// KyberFeeHandlerBRRUpdatedIterator is returned from FilterBRRUpdated and is used to iterate over the raw logs and unpacked data for BRRUpdated events raised by the KyberFeeHandler contract.
type KyberFeeHandlerBRRUpdatedIterator struct {
	Event *KyberFeeHandlerBRRUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerBRRUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerBRRUpdated)
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
		it.Event = new(KyberFeeHandlerBRRUpdated)
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
func (it *KyberFeeHandlerBRRUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerBRRUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerBRRUpdated represents a BRRUpdated event raised by the KyberFeeHandler contract.
type KyberFeeHandlerBRRUpdated struct {
	RewardBps       *big.Int
	RebateBps       *big.Int
	BurnBps         *big.Int
	ExpiryTimestamp *big.Int
	Epoch           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBRRUpdated is a free log retrieval operation binding the contract event 0x4b3150a36b957ed95a132721c7412af319174861da7c8c7a55ef6e1a2794528d.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps, uint256 expiryTimestamp, uint256 indexed epoch)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterBRRUpdated(opts *bind.FilterOpts, epoch []*big.Int) (*KyberFeeHandlerBRRUpdatedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "BRRUpdated", epochRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerBRRUpdatedIterator{contract: _KyberFeeHandler.contract, event: "BRRUpdated", logs: logs, sub: sub}, nil
}

// WatchBRRUpdated is a free log subscription operation binding the contract event 0x4b3150a36b957ed95a132721c7412af319174861da7c8c7a55ef6e1a2794528d.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps, uint256 expiryTimestamp, uint256 indexed epoch)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchBRRUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerBRRUpdated, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "BRRUpdated", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerBRRUpdated)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "BRRUpdated", log); err != nil {
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

// ParseBRRUpdated is a log parse operation binding the contract event 0x4b3150a36b957ed95a132721c7412af319174861da7c8c7a55ef6e1a2794528d.
//
// Solidity: event BRRUpdated(uint256 rewardBps, uint256 rebateBps, uint256 burnBps, uint256 expiryTimestamp, uint256 indexed epoch)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseBRRUpdated(log types.Log) (*KyberFeeHandlerBRRUpdated, error) {
	event := new(KyberFeeHandlerBRRUpdated)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "BRRUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerBurnConfigSetIterator is returned from FilterBurnConfigSet and is used to iterate over the raw logs and unpacked data for BurnConfigSet events raised by the KyberFeeHandler contract.
type KyberFeeHandlerBurnConfigSetIterator struct {
	Event *KyberFeeHandlerBurnConfigSet // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerBurnConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerBurnConfigSet)
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
		it.Event = new(KyberFeeHandlerBurnConfigSet)
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
func (it *KyberFeeHandlerBurnConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerBurnConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerBurnConfigSet represents a BurnConfigSet event raised by the KyberFeeHandler contract.
type KyberFeeHandlerBurnConfigSet struct {
	SanityRate common.Address
	WeiToBurn  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBurnConfigSet is a free log retrieval operation binding the contract event 0xe40f97f23269c4682610e9b2522d6d4272ee56f115906d71fcb3da82a860f755.
//
// Solidity: event BurnConfigSet(address sanityRate, uint256 weiToBurn)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterBurnConfigSet(opts *bind.FilterOpts) (*KyberFeeHandlerBurnConfigSetIterator, error) {

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "BurnConfigSet")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerBurnConfigSetIterator{contract: _KyberFeeHandler.contract, event: "BurnConfigSet", logs: logs, sub: sub}, nil
}

// WatchBurnConfigSet is a free log subscription operation binding the contract event 0xe40f97f23269c4682610e9b2522d6d4272ee56f115906d71fcb3da82a860f755.
//
// Solidity: event BurnConfigSet(address sanityRate, uint256 weiToBurn)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchBurnConfigSet(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerBurnConfigSet) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "BurnConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerBurnConfigSet)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "BurnConfigSet", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseBurnConfigSet(log types.Log) (*KyberFeeHandlerBurnConfigSet, error) {
	event := new(KyberFeeHandlerBurnConfigSet)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "BurnConfigSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerEthReceivedIterator is returned from FilterEthReceived and is used to iterate over the raw logs and unpacked data for EthReceived events raised by the KyberFeeHandler contract.
type KyberFeeHandlerEthReceivedIterator struct {
	Event *KyberFeeHandlerEthReceived // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerEthReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerEthReceived)
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
		it.Event = new(KyberFeeHandlerEthReceived)
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
func (it *KyberFeeHandlerEthReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerEthReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerEthReceived represents a EthReceived event raised by the KyberFeeHandler contract.
type KyberFeeHandlerEthReceived struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEthReceived is a free log retrieval operation binding the contract event 0x353bcaaf167a6add95a753d39727e3d3beb865129a69a10ed774b0b899671403.
//
// Solidity: event EthReceived(uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterEthReceived(opts *bind.FilterOpts) (*KyberFeeHandlerEthReceivedIterator, error) {

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "EthReceived")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerEthReceivedIterator{contract: _KyberFeeHandler.contract, event: "EthReceived", logs: logs, sub: sub}, nil
}

// WatchEthReceived is a free log subscription operation binding the contract event 0x353bcaaf167a6add95a753d39727e3d3beb865129a69a10ed774b0b899671403.
//
// Solidity: event EthReceived(uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchEthReceived(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerEthReceived) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "EthReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerEthReceived)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "EthReceived", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseEthReceived(log types.Log) (*KyberFeeHandlerEthReceived, error) {
	event := new(KyberFeeHandlerEthReceived)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "EthReceived", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerFeeDistributedIterator is returned from FilterFeeDistributed and is used to iterate over the raw logs and unpacked data for FeeDistributed events raised by the KyberFeeHandler contract.
type KyberFeeHandlerFeeDistributedIterator struct {
	Event *KyberFeeHandlerFeeDistributed // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerFeeDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerFeeDistributed)
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
		it.Event = new(KyberFeeHandlerFeeDistributed)
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
func (it *KyberFeeHandlerFeeDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerFeeDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerFeeDistributed represents a FeeDistributed event raised by the KyberFeeHandler contract.
type KyberFeeHandlerFeeDistributed struct {
	Token                     common.Address
	PlatformWallet            common.Address
	PlatformFeeWei            *big.Int
	RewardWei                 *big.Int
	RebateWei                 *big.Int
	RebateWallets             []common.Address
	RebatePercentBpsPerWallet []*big.Int
	BurnAmtWei                *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterFeeDistributed is a free log retrieval operation binding the contract event 0x53e2e1b5ab64e0a76fcc6a932558eba265d4e58c512401a7d776ae0f8fc08994.
//
// Solidity: event FeeDistributed(address indexed token, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterFeeDistributed(opts *bind.FilterOpts, token []common.Address, platformWallet []common.Address) (*KyberFeeHandlerFeeDistributedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "FeeDistributed", tokenRule, platformWalletRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerFeeDistributedIterator{contract: _KyberFeeHandler.contract, event: "FeeDistributed", logs: logs, sub: sub}, nil
}

// WatchFeeDistributed is a free log subscription operation binding the contract event 0x53e2e1b5ab64e0a76fcc6a932558eba265d4e58c512401a7d776ae0f8fc08994.
//
// Solidity: event FeeDistributed(address indexed token, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchFeeDistributed(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerFeeDistributed, token []common.Address, platformWallet []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "FeeDistributed", tokenRule, platformWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerFeeDistributed)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "FeeDistributed", log); err != nil {
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

// ParseFeeDistributed is a log parse operation binding the contract event 0x53e2e1b5ab64e0a76fcc6a932558eba265d4e58c512401a7d776ae0f8fc08994.
//
// Solidity: event FeeDistributed(address indexed token, address indexed platformWallet, uint256 platformFeeWei, uint256 rewardWei, uint256 rebateWei, address[] rebateWallets, uint256[] rebatePercentBpsPerWallet, uint256 burnAmtWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseFeeDistributed(log types.Log) (*KyberFeeHandlerFeeDistributed, error) {
	event := new(KyberFeeHandlerFeeDistributed)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "FeeDistributed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerKncBurnedIterator is returned from FilterKncBurned and is used to iterate over the raw logs and unpacked data for KncBurned events raised by the KyberFeeHandler contract.
type KyberFeeHandlerKncBurnedIterator struct {
	Event *KyberFeeHandlerKncBurned // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerKncBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerKncBurned)
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
		it.Event = new(KyberFeeHandlerKncBurned)
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
func (it *KyberFeeHandlerKncBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerKncBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerKncBurned represents a KncBurned event raised by the KyberFeeHandler contract.
type KyberFeeHandlerKncBurned struct {
	KncTWei *big.Int
	Token   common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKncBurned is a free log retrieval operation binding the contract event 0xa0fcef56e2b45fcbeb91d5e629ef6b2b6e982d0768f02d1232610315cd23ea10.
//
// Solidity: event KncBurned(uint256 kncTWei, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterKncBurned(opts *bind.FilterOpts, token []common.Address) (*KyberFeeHandlerKncBurnedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "KncBurned", tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerKncBurnedIterator{contract: _KyberFeeHandler.contract, event: "KncBurned", logs: logs, sub: sub}, nil
}

// WatchKncBurned is a free log subscription operation binding the contract event 0xa0fcef56e2b45fcbeb91d5e629ef6b2b6e982d0768f02d1232610315cd23ea10.
//
// Solidity: event KncBurned(uint256 kncTWei, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchKncBurned(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerKncBurned, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "KncBurned", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerKncBurned)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "KncBurned", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseKncBurned(log types.Log) (*KyberFeeHandlerKncBurned, error) {
	event := new(KyberFeeHandlerKncBurned)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "KncBurned", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerKyberDaoAddressSetIterator is returned from FilterKyberDaoAddressSet and is used to iterate over the raw logs and unpacked data for KyberDaoAddressSet events raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberDaoAddressSetIterator struct {
	Event *KyberFeeHandlerKyberDaoAddressSet // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerKyberDaoAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerKyberDaoAddressSet)
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
		it.Event = new(KyberFeeHandlerKyberDaoAddressSet)
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
func (it *KyberFeeHandlerKyberDaoAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerKyberDaoAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerKyberDaoAddressSet represents a KyberDaoAddressSet event raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberDaoAddressSet struct {
	KyberDao common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterKyberDaoAddressSet is a free log retrieval operation binding the contract event 0x6792eb8fe9de88d4eaaee7128e99aede17da98cd391520d3ec51a365804722c4.
//
// Solidity: event KyberDaoAddressSet(address kyberDao)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterKyberDaoAddressSet(opts *bind.FilterOpts) (*KyberFeeHandlerKyberDaoAddressSetIterator, error) {

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "KyberDaoAddressSet")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerKyberDaoAddressSetIterator{contract: _KyberFeeHandler.contract, event: "KyberDaoAddressSet", logs: logs, sub: sub}, nil
}

// WatchKyberDaoAddressSet is a free log subscription operation binding the contract event 0x6792eb8fe9de88d4eaaee7128e99aede17da98cd391520d3ec51a365804722c4.
//
// Solidity: event KyberDaoAddressSet(address kyberDao)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchKyberDaoAddressSet(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerKyberDaoAddressSet) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "KyberDaoAddressSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerKyberDaoAddressSet)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberDaoAddressSet", log); err != nil {
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

// ParseKyberDaoAddressSet is a log parse operation binding the contract event 0x6792eb8fe9de88d4eaaee7128e99aede17da98cd391520d3ec51a365804722c4.
//
// Solidity: event KyberDaoAddressSet(address kyberDao)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseKyberDaoAddressSet(log types.Log) (*KyberFeeHandlerKyberDaoAddressSet, error) {
	event := new(KyberFeeHandlerKyberDaoAddressSet)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberDaoAddressSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerKyberNetworkUpdatedIterator is returned from FilterKyberNetworkUpdated and is used to iterate over the raw logs and unpacked data for KyberNetworkUpdated events raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberNetworkUpdatedIterator struct {
	Event *KyberFeeHandlerKyberNetworkUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerKyberNetworkUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerKyberNetworkUpdated)
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
		it.Event = new(KyberFeeHandlerKyberNetworkUpdated)
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
func (it *KyberFeeHandlerKyberNetworkUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerKyberNetworkUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerKyberNetworkUpdated represents a KyberNetworkUpdated event raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberNetworkUpdated struct {
	KyberNetwork common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkUpdated is a free log retrieval operation binding the contract event 0x18970d46ac8a7d7e0da90e1bebb0be3e87ffc7705fc09d3bba5373d59b7a12aa.
//
// Solidity: event KyberNetworkUpdated(address kyberNetwork)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterKyberNetworkUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerKyberNetworkUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "KyberNetworkUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerKyberNetworkUpdatedIterator{contract: _KyberFeeHandler.contract, event: "KyberNetworkUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkUpdated is a free log subscription operation binding the contract event 0x18970d46ac8a7d7e0da90e1bebb0be3e87ffc7705fc09d3bba5373d59b7a12aa.
//
// Solidity: event KyberNetworkUpdated(address kyberNetwork)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchKyberNetworkUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerKyberNetworkUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "KyberNetworkUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerKyberNetworkUpdated)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberNetworkUpdated", log); err != nil {
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

// ParseKyberNetworkUpdated is a log parse operation binding the contract event 0x18970d46ac8a7d7e0da90e1bebb0be3e87ffc7705fc09d3bba5373d59b7a12aa.
//
// Solidity: event KyberNetworkUpdated(address kyberNetwork)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseKyberNetworkUpdated(log types.Log) (*KyberFeeHandlerKyberNetworkUpdated, error) {
	event := new(KyberFeeHandlerKyberNetworkUpdated)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberNetworkUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerKyberProxyUpdatedIterator is returned from FilterKyberProxyUpdated and is used to iterate over the raw logs and unpacked data for KyberProxyUpdated events raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberProxyUpdatedIterator struct {
	Event *KyberFeeHandlerKyberProxyUpdated // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerKyberProxyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerKyberProxyUpdated)
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
		it.Event = new(KyberFeeHandlerKyberProxyUpdated)
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
func (it *KyberFeeHandlerKyberProxyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerKyberProxyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerKyberProxyUpdated represents a KyberProxyUpdated event raised by the KyberFeeHandler contract.
type KyberFeeHandlerKyberProxyUpdated struct {
	KyberProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterKyberProxyUpdated is a free log retrieval operation binding the contract event 0x8457f9bd0d13488a6c265af376d291f3c6bd2311d9e8dee5671d4169ca6e0ae0.
//
// Solidity: event KyberProxyUpdated(address kyberProxy)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterKyberProxyUpdated(opts *bind.FilterOpts) (*KyberFeeHandlerKyberProxyUpdatedIterator, error) {

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "KyberProxyUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerKyberProxyUpdatedIterator{contract: _KyberFeeHandler.contract, event: "KyberProxyUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberProxyUpdated is a free log subscription operation binding the contract event 0x8457f9bd0d13488a6c265af376d291f3c6bd2311d9e8dee5671d4169ca6e0ae0.
//
// Solidity: event KyberProxyUpdated(address kyberProxy)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchKyberProxyUpdated(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerKyberProxyUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "KyberProxyUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerKyberProxyUpdated)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberProxyUpdated", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseKyberProxyUpdated(log types.Log) (*KyberFeeHandlerKyberProxyUpdated, error) {
	event := new(KyberFeeHandlerKyberProxyUpdated)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "KyberProxyUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerPlatformFeePaidIterator is returned from FilterPlatformFeePaid and is used to iterate over the raw logs and unpacked data for PlatformFeePaid events raised by the KyberFeeHandler contract.
type KyberFeeHandlerPlatformFeePaidIterator struct {
	Event *KyberFeeHandlerPlatformFeePaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerPlatformFeePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerPlatformFeePaid)
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
		it.Event = new(KyberFeeHandlerPlatformFeePaid)
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
func (it *KyberFeeHandlerPlatformFeePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerPlatformFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerPlatformFeePaid represents a PlatformFeePaid event raised by the KyberFeeHandler contract.
type KyberFeeHandlerPlatformFeePaid struct {
	PlatformWallet common.Address
	Token          common.Address
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPlatformFeePaid is a free log retrieval operation binding the contract event 0xebe3db09f5650582b4782506e0d272262129183570e55fcf8768dd6e91f8c0f6.
//
// Solidity: event PlatformFeePaid(address indexed platformWallet, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterPlatformFeePaid(opts *bind.FilterOpts, platformWallet []common.Address, token []common.Address) (*KyberFeeHandlerPlatformFeePaidIterator, error) {

	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "PlatformFeePaid", platformWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerPlatformFeePaidIterator{contract: _KyberFeeHandler.contract, event: "PlatformFeePaid", logs: logs, sub: sub}, nil
}

// WatchPlatformFeePaid is a free log subscription operation binding the contract event 0xebe3db09f5650582b4782506e0d272262129183570e55fcf8768dd6e91f8c0f6.
//
// Solidity: event PlatformFeePaid(address indexed platformWallet, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchPlatformFeePaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerPlatformFeePaid, platformWallet []common.Address, token []common.Address) (event.Subscription, error) {

	var platformWalletRule []interface{}
	for _, platformWalletItem := range platformWallet {
		platformWalletRule = append(platformWalletRule, platformWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "PlatformFeePaid", platformWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerPlatformFeePaid)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "PlatformFeePaid", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParsePlatformFeePaid(log types.Log) (*KyberFeeHandlerPlatformFeePaid, error) {
	event := new(KyberFeeHandlerPlatformFeePaid)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "PlatformFeePaid", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerRebatePaidIterator is returned from FilterRebatePaid and is used to iterate over the raw logs and unpacked data for RebatePaid events raised by the KyberFeeHandler contract.
type KyberFeeHandlerRebatePaidIterator struct {
	Event *KyberFeeHandlerRebatePaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerRebatePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerRebatePaid)
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
		it.Event = new(KyberFeeHandlerRebatePaid)
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
func (it *KyberFeeHandlerRebatePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerRebatePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerRebatePaid represents a RebatePaid event raised by the KyberFeeHandler contract.
type KyberFeeHandlerRebatePaid struct {
	RebateWallet common.Address
	Token        common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRebatePaid is a free log retrieval operation binding the contract event 0xb5ec5e03662403108373ab6431d3e834cb1011fca164541aef315fc7dea7b3b6.
//
// Solidity: event RebatePaid(address indexed rebateWallet, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterRebatePaid(opts *bind.FilterOpts, rebateWallet []common.Address, token []common.Address) (*KyberFeeHandlerRebatePaidIterator, error) {

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "RebatePaid", rebateWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerRebatePaidIterator{contract: _KyberFeeHandler.contract, event: "RebatePaid", logs: logs, sub: sub}, nil
}

// WatchRebatePaid is a free log subscription operation binding the contract event 0xb5ec5e03662403108373ab6431d3e834cb1011fca164541aef315fc7dea7b3b6.
//
// Solidity: event RebatePaid(address indexed rebateWallet, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchRebatePaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerRebatePaid, rebateWallet []common.Address, token []common.Address) (event.Subscription, error) {

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "RebatePaid", rebateWalletRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerRebatePaid)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "RebatePaid", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseRebatePaid(log types.Log) (*KyberFeeHandlerRebatePaid, error) {
	event := new(KyberFeeHandlerRebatePaid)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "RebatePaid", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerRewardPaidIterator is returned from FilterRewardPaid and is used to iterate over the raw logs and unpacked data for RewardPaid events raised by the KyberFeeHandler contract.
type KyberFeeHandlerRewardPaidIterator struct {
	Event *KyberFeeHandlerRewardPaid // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerRewardPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerRewardPaid)
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
		it.Event = new(KyberFeeHandlerRewardPaid)
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
func (it *KyberFeeHandlerRewardPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerRewardPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerRewardPaid represents a RewardPaid event raised by the KyberFeeHandler contract.
type KyberFeeHandlerRewardPaid struct {
	Staker common.Address
	Epoch  *big.Int
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardPaid is a free log retrieval operation binding the contract event 0xaf206e736916d38b56e2d559931a189bc3119b8fc6d6850bd34e382f09030587.
//
// Solidity: event RewardPaid(address indexed staker, uint256 indexed epoch, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterRewardPaid(opts *bind.FilterOpts, staker []common.Address, epoch []*big.Int, token []common.Address) (*KyberFeeHandlerRewardPaidIterator, error) {

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

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "RewardPaid", stakerRule, epochRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerRewardPaidIterator{contract: _KyberFeeHandler.contract, event: "RewardPaid", logs: logs, sub: sub}, nil
}

// WatchRewardPaid is a free log subscription operation binding the contract event 0xaf206e736916d38b56e2d559931a189bc3119b8fc6d6850bd34e382f09030587.
//
// Solidity: event RewardPaid(address indexed staker, uint256 indexed epoch, address indexed token, uint256 amount)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchRewardPaid(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerRewardPaid, staker []common.Address, epoch []*big.Int, token []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "RewardPaid", stakerRule, epochRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerRewardPaid)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "RewardPaid", log); err != nil {
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
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseRewardPaid(log types.Log) (*KyberFeeHandlerRewardPaid, error) {
	event := new(KyberFeeHandlerRewardPaid)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "RewardPaid", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberFeeHandlerRewardsRemovedToBurnIterator is returned from FilterRewardsRemovedToBurn and is used to iterate over the raw logs and unpacked data for RewardsRemovedToBurn events raised by the KyberFeeHandler contract.
type KyberFeeHandlerRewardsRemovedToBurnIterator struct {
	Event *KyberFeeHandlerRewardsRemovedToBurn // Event containing the contract specifics and raw log

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
func (it *KyberFeeHandlerRewardsRemovedToBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberFeeHandlerRewardsRemovedToBurn)
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
		it.Event = new(KyberFeeHandlerRewardsRemovedToBurn)
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
func (it *KyberFeeHandlerRewardsRemovedToBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberFeeHandlerRewardsRemovedToBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberFeeHandlerRewardsRemovedToBurn represents a RewardsRemovedToBurn event raised by the KyberFeeHandler contract.
type KyberFeeHandlerRewardsRemovedToBurn struct {
	Epoch      *big.Int
	RewardsWei *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRewardsRemovedToBurn is a free log retrieval operation binding the contract event 0x11c852d8be537f120b8d4b4d5c3c211870522fd96a8bd9fa51d102774077a51b.
//
// Solidity: event RewardsRemovedToBurn(uint256 indexed epoch, uint256 rewardsWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) FilterRewardsRemovedToBurn(opts *bind.FilterOpts, epoch []*big.Int) (*KyberFeeHandlerRewardsRemovedToBurnIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.FilterLogs(opts, "RewardsRemovedToBurn", epochRule)
	if err != nil {
		return nil, err
	}
	return &KyberFeeHandlerRewardsRemovedToBurnIterator{contract: _KyberFeeHandler.contract, event: "RewardsRemovedToBurn", logs: logs, sub: sub}, nil
}

// WatchRewardsRemovedToBurn is a free log subscription operation binding the contract event 0x11c852d8be537f120b8d4b4d5c3c211870522fd96a8bd9fa51d102774077a51b.
//
// Solidity: event RewardsRemovedToBurn(uint256 indexed epoch, uint256 rewardsWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) WatchRewardsRemovedToBurn(opts *bind.WatchOpts, sink chan<- *KyberFeeHandlerRewardsRemovedToBurn, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _KyberFeeHandler.contract.WatchLogs(opts, "RewardsRemovedToBurn", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberFeeHandlerRewardsRemovedToBurn)
				if err := _KyberFeeHandler.contract.UnpackLog(event, "RewardsRemovedToBurn", log); err != nil {
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

// ParseRewardsRemovedToBurn is a log parse operation binding the contract event 0x11c852d8be537f120b8d4b4d5c3c211870522fd96a8bd9fa51d102774077a51b.
//
// Solidity: event RewardsRemovedToBurn(uint256 indexed epoch, uint256 rewardsWei)
func (_KyberFeeHandler *KyberFeeHandlerFilterer) ParseRewardsRemovedToBurn(log types.Log) (*KyberFeeHandlerRewardsRemovedToBurn, error) {
	event := new(KyberFeeHandlerRewardsRemovedToBurn)
	if err := _KyberFeeHandler.contract.UnpackLog(event, "RewardsRemovedToBurn", log); err != nil {
		return nil, err
	}
	return event, nil
}
