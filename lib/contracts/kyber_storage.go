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

// KyberStorageABI is the input ABI used to generate the binding from.
const KyberStorageABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"contractIKyberHistory\",\"name\":\"_kyberNetworkHistory\",\"type\":\"address\"},{\"internalType\":\"contractIKyberHistory\",\"name\":\"_kyberFeeHandlerHistory\",\"type\":\"address\"},{\"internalType\":\"contractIKyberHistory\",\"name\":\"_kyberDaoHistory\",\"type\":\"address\"},{\"internalType\":\"contractIKyberHistory\",\"name\":\"_kyberMatchingEngineHistory\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumIKyberStorage.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"AddReserveToStorage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberNetwork\",\"name\":\"newKyberNetwork\",\"type\":\"address\"}],\"name\":\"KyberNetworkUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"ListReservePairs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"}],\"name\":\"RemoveReserveFromStorage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"}],\"name\":\"ReserveRebateWalletSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxApprovedProxies\",\"type\":\"uint256\"}],\"name\":\"addKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"internalType\":\"enumIKyberStorage.ReserveType\",\"name\":\"resType\",\"type\":\"uint8\"},{\"internalType\":\"addresspayable\",\"name\":\"rebateWallet\",\"type\":\"address\"}],\"name\":\"addReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"kyberDaoAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"kyberFeeHandlerAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"kyberMatchingEngineAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"kyberNetworkAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"name\":\"getEntitledRebateData\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"entitledRebateArr\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"name\":\"getFeeAccountedData\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"feeAccountedArr\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKyberProxies\",\"outputs\":[{\"internalType\":\"contractIKyberNetworkProxy[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"}],\"name\":\"getListedTokensByReserveId\",\"outputs\":[{\"internalType\":\"contractIERC20[]\",\"name\":\"srcTokens\",\"type\":\"address[]\"},{\"internalType\":\"contractIERC20[]\",\"name\":\"destTokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"name\":\"getRebateWalletsFromIds\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"rebateWallets\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"}],\"name\":\"getReserveAddressesByReserveId\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"reserveAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"name\":\"getReserveAddressesFromIds\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"reserveAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getReserveAddressesPerTokenSrc\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"reserveAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"}],\"name\":\"getReserveDetailsByAddress\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"},{\"internalType\":\"enumIKyberStorage.ReserveType\",\"name\":\"resType\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isFeeAccountedFlag\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isEntitledRebateFlag\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"}],\"name\":\"getReserveDetailsById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"reserveAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"},{\"internalType\":\"enumIKyberStorage.ReserveType\",\"name\":\"resType\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isFeeAccountedFlag\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isEntitledRebateFlag\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"}],\"name\":\"getReserveId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"reserveAddresses\",\"type\":\"address[]\"}],\"name\":\"getReserveIdsFromAddresses\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getReserveIdsPerTokenDest\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getReserveIdsPerTokenSrc\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"contractIKyberReserve[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"reserveIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"}],\"name\":\"getReservesData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"areAllReservesListed\",\"type\":\"bool\"},{\"internalType\":\"bool[]\",\"name\":\"feeAccountedArr\",\"type\":\"bool[]\"},{\"internalType\":\"bool[]\",\"name\":\"entitledRebateArr\",\"type\":\"bool[]\"},{\"internalType\":\"contractIKyberReserve[]\",\"name\":\"reserveAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIKyberStorage.ReserveType\",\"name\":\"resType\",\"type\":\"uint8\"}],\"name\":\"getReservesPerType\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isKyberProxyAdded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberDaoHistory\",\"outputs\":[{\"internalType\":\"contractIKyberHistory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberFeeHandlerHistory\",\"outputs\":[{\"internalType\":\"contractIKyberHistory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberMatchingEngineHistory\",\"outputs\":[{\"internalType\":\"contractIKyberHistory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberNetwork\",\"outputs\":[{\"internalType\":\"contractIKyberNetwork\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kyberNetworkHistory\",\"outputs\":[{\"internalType\":\"contractIKyberHistory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ethToToken\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"tokenToEth\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"listPairForReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"removeKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"}],\"name\":\"removeReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_kyberFeeHandler\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_kyberMatchingEngine\",\"type\":\"address\"}],\"name\":\"setContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"fpr\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"apr\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"bridge\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"utility\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"custom\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"orderbook\",\"type\":\"bool\"}],\"name\":\"setEntitledRebatePerReserveType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"fpr\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"apr\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"bridge\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"utility\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"custom\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"orderbook\",\"type\":\"bool\"}],\"name\":\"setFeeAccountedPerReserveType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_kyberDao\",\"type\":\"address\"}],\"name\":\"setKyberDaoContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberNetwork\",\"name\":\"_kyberNetwork\",\"type\":\"address\"}],\"name\":\"setNetworkContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reserveId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"rebateWallet\",\"type\":\"address\"}],\"name\":\"setRebateWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KyberStorage is an auto generated Go binding around an Ethereum contract.
type KyberStorage struct {
	KyberStorageCaller     // Read-only binding to the contract
	KyberStorageTransactor // Write-only binding to the contract
	KyberStorageFilterer   // Log filterer for contract events
}

// KyberStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type KyberStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KyberStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KyberStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KyberStorageSession struct {
	Contract     *KyberStorage     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KyberStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KyberStorageCallerSession struct {
	Contract *KyberStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// KyberStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KyberStorageTransactorSession struct {
	Contract     *KyberStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// KyberStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type KyberStorageRaw struct {
	Contract *KyberStorage // Generic contract binding to access the raw methods on
}

// KyberStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KyberStorageCallerRaw struct {
	Contract *KyberStorageCaller // Generic read-only contract binding to access the raw methods on
}

// KyberStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KyberStorageTransactorRaw struct {
	Contract *KyberStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKyberStorage creates a new instance of KyberStorage, bound to a specific deployed contract.
func NewKyberStorage(address common.Address, backend bind.ContractBackend) (*KyberStorage, error) {
	contract, err := bindKyberStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KyberStorage{KyberStorageCaller: KyberStorageCaller{contract: contract}, KyberStorageTransactor: KyberStorageTransactor{contract: contract}, KyberStorageFilterer: KyberStorageFilterer{contract: contract}}, nil
}

// NewKyberStorageCaller creates a new read-only instance of KyberStorage, bound to a specific deployed contract.
func NewKyberStorageCaller(address common.Address, caller bind.ContractCaller) (*KyberStorageCaller, error) {
	contract, err := bindKyberStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KyberStorageCaller{contract: contract}, nil
}

// NewKyberStorageTransactor creates a new write-only instance of KyberStorage, bound to a specific deployed contract.
func NewKyberStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*KyberStorageTransactor, error) {
	contract, err := bindKyberStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KyberStorageTransactor{contract: contract}, nil
}

// NewKyberStorageFilterer creates a new log filterer instance of KyberStorage, bound to a specific deployed contract.
func NewKyberStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*KyberStorageFilterer, error) {
	contract, err := bindKyberStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KyberStorageFilterer{contract: contract}, nil
}

// bindKyberStorage binds a generic wrapper to an already deployed contract.
func bindKyberStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KyberStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberStorage *KyberStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberStorage.Contract.KyberStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberStorage *KyberStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberStorage.Contract.KyberStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberStorage *KyberStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberStorage.Contract.KyberStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberStorage *KyberStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberStorage *KyberStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberStorage *KyberStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberStorage.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberStorage *KyberStorageCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberStorage *KyberStorageSession) Admin() (common.Address, error) {
	return _KyberStorage.Contract.Admin(&_KyberStorage.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) Admin() (common.Address, error) {
	return _KyberStorage.Contract.Admin(&_KyberStorage.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberStorage *KyberStorageCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberStorage *KyberStorageSession) GetAlerters() ([]common.Address, error) {
	return _KyberStorage.Contract.GetAlerters(&_KyberStorage.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberStorage *KyberStorageCallerSession) GetAlerters() ([]common.Address, error) {
	return _KyberStorage.Contract.GetAlerters(&_KyberStorage.CallOpts)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address[] kyberDaoAddresses, address[] kyberFeeHandlerAddresses, address[] kyberMatchingEngineAddresses, address[] kyberNetworkAddresses)
func (_KyberStorage *KyberStorageCaller) GetContracts(opts *bind.CallOpts) (struct {
	KyberDaoAddresses            []common.Address
	KyberFeeHandlerAddresses     []common.Address
	KyberMatchingEngineAddresses []common.Address
	KyberNetworkAddresses        []common.Address
}, error) {
	ret := new(struct {
		KyberDaoAddresses            []common.Address
		KyberFeeHandlerAddresses     []common.Address
		KyberMatchingEngineAddresses []common.Address
		KyberNetworkAddresses        []common.Address
	})
	out := ret
	err := _KyberStorage.contract.Call(opts, out, "getContracts")
	return *ret, err
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address[] kyberDaoAddresses, address[] kyberFeeHandlerAddresses, address[] kyberMatchingEngineAddresses, address[] kyberNetworkAddresses)
func (_KyberStorage *KyberStorageSession) GetContracts() (struct {
	KyberDaoAddresses            []common.Address
	KyberFeeHandlerAddresses     []common.Address
	KyberMatchingEngineAddresses []common.Address
	KyberNetworkAddresses        []common.Address
}, error) {
	return _KyberStorage.Contract.GetContracts(&_KyberStorage.CallOpts)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address[] kyberDaoAddresses, address[] kyberFeeHandlerAddresses, address[] kyberMatchingEngineAddresses, address[] kyberNetworkAddresses)
func (_KyberStorage *KyberStorageCallerSession) GetContracts() (struct {
	KyberDaoAddresses            []common.Address
	KyberFeeHandlerAddresses     []common.Address
	KyberMatchingEngineAddresses []common.Address
	KyberNetworkAddresses        []common.Address
}, error) {
	return _KyberStorage.Contract.GetContracts(&_KyberStorage.CallOpts)
}

// GetEntitledRebateData is a free data retrieval call binding the contract method 0x1dd4c3ef.
//
// Solidity: function getEntitledRebateData(bytes32[] reserveIds) view returns(bool[] entitledRebateArr)
func (_KyberStorage *KyberStorageCaller) GetEntitledRebateData(opts *bind.CallOpts, reserveIds [][32]byte) ([]bool, error) {
	var (
		ret0 = new([]bool)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getEntitledRebateData", reserveIds)
	return *ret0, err
}

// GetEntitledRebateData is a free data retrieval call binding the contract method 0x1dd4c3ef.
//
// Solidity: function getEntitledRebateData(bytes32[] reserveIds) view returns(bool[] entitledRebateArr)
func (_KyberStorage *KyberStorageSession) GetEntitledRebateData(reserveIds [][32]byte) ([]bool, error) {
	return _KyberStorage.Contract.GetEntitledRebateData(&_KyberStorage.CallOpts, reserveIds)
}

// GetEntitledRebateData is a free data retrieval call binding the contract method 0x1dd4c3ef.
//
// Solidity: function getEntitledRebateData(bytes32[] reserveIds) view returns(bool[] entitledRebateArr)
func (_KyberStorage *KyberStorageCallerSession) GetEntitledRebateData(reserveIds [][32]byte) ([]bool, error) {
	return _KyberStorage.Contract.GetEntitledRebateData(&_KyberStorage.CallOpts, reserveIds)
}

// GetFeeAccountedData is a free data retrieval call binding the contract method 0xe0fb2756.
//
// Solidity: function getFeeAccountedData(bytes32[] reserveIds) view returns(bool[] feeAccountedArr)
func (_KyberStorage *KyberStorageCaller) GetFeeAccountedData(opts *bind.CallOpts, reserveIds [][32]byte) ([]bool, error) {
	var (
		ret0 = new([]bool)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getFeeAccountedData", reserveIds)
	return *ret0, err
}

// GetFeeAccountedData is a free data retrieval call binding the contract method 0xe0fb2756.
//
// Solidity: function getFeeAccountedData(bytes32[] reserveIds) view returns(bool[] feeAccountedArr)
func (_KyberStorage *KyberStorageSession) GetFeeAccountedData(reserveIds [][32]byte) ([]bool, error) {
	return _KyberStorage.Contract.GetFeeAccountedData(&_KyberStorage.CallOpts, reserveIds)
}

// GetFeeAccountedData is a free data retrieval call binding the contract method 0xe0fb2756.
//
// Solidity: function getFeeAccountedData(bytes32[] reserveIds) view returns(bool[] feeAccountedArr)
func (_KyberStorage *KyberStorageCallerSession) GetFeeAccountedData(reserveIds [][32]byte) ([]bool, error) {
	return _KyberStorage.Contract.GetFeeAccountedData(&_KyberStorage.CallOpts, reserveIds)
}

// GetKyberProxies is a free data retrieval call binding the contract method 0xfa006d93.
//
// Solidity: function getKyberProxies() view returns(address[])
func (_KyberStorage *KyberStorageCaller) GetKyberProxies(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getKyberProxies")
	return *ret0, err
}

// GetKyberProxies is a free data retrieval call binding the contract method 0xfa006d93.
//
// Solidity: function getKyberProxies() view returns(address[])
func (_KyberStorage *KyberStorageSession) GetKyberProxies() ([]common.Address, error) {
	return _KyberStorage.Contract.GetKyberProxies(&_KyberStorage.CallOpts)
}

// GetKyberProxies is a free data retrieval call binding the contract method 0xfa006d93.
//
// Solidity: function getKyberProxies() view returns(address[])
func (_KyberStorage *KyberStorageCallerSession) GetKyberProxies() ([]common.Address, error) {
	return _KyberStorage.Contract.GetKyberProxies(&_KyberStorage.CallOpts)
}

// GetListedTokensByReserveId is a free data retrieval call binding the contract method 0x3d9036b4.
//
// Solidity: function getListedTokensByReserveId(bytes32 reserveId) view returns(address[] srcTokens, address[] destTokens)
func (_KyberStorage *KyberStorageCaller) GetListedTokensByReserveId(opts *bind.CallOpts, reserveId [32]byte) (struct {
	SrcTokens  []common.Address
	DestTokens []common.Address
}, error) {
	ret := new(struct {
		SrcTokens  []common.Address
		DestTokens []common.Address
	})
	out := ret
	err := _KyberStorage.contract.Call(opts, out, "getListedTokensByReserveId", reserveId)
	return *ret, err
}

// GetListedTokensByReserveId is a free data retrieval call binding the contract method 0x3d9036b4.
//
// Solidity: function getListedTokensByReserveId(bytes32 reserveId) view returns(address[] srcTokens, address[] destTokens)
func (_KyberStorage *KyberStorageSession) GetListedTokensByReserveId(reserveId [32]byte) (struct {
	SrcTokens  []common.Address
	DestTokens []common.Address
}, error) {
	return _KyberStorage.Contract.GetListedTokensByReserveId(&_KyberStorage.CallOpts, reserveId)
}

// GetListedTokensByReserveId is a free data retrieval call binding the contract method 0x3d9036b4.
//
// Solidity: function getListedTokensByReserveId(bytes32 reserveId) view returns(address[] srcTokens, address[] destTokens)
func (_KyberStorage *KyberStorageCallerSession) GetListedTokensByReserveId(reserveId [32]byte) (struct {
	SrcTokens  []common.Address
	DestTokens []common.Address
}, error) {
	return _KyberStorage.Contract.GetListedTokensByReserveId(&_KyberStorage.CallOpts, reserveId)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberStorage *KyberStorageCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberStorage *KyberStorageSession) GetOperators() ([]common.Address, error) {
	return _KyberStorage.Contract.GetOperators(&_KyberStorage.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberStorage *KyberStorageCallerSession) GetOperators() ([]common.Address, error) {
	return _KyberStorage.Contract.GetOperators(&_KyberStorage.CallOpts)
}

// GetRebateWalletsFromIds is a free data retrieval call binding the contract method 0x0a3cf98e.
//
// Solidity: function getRebateWalletsFromIds(bytes32[] reserveIds) view returns(address[] rebateWallets)
func (_KyberStorage *KyberStorageCaller) GetRebateWalletsFromIds(opts *bind.CallOpts, reserveIds [][32]byte) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getRebateWalletsFromIds", reserveIds)
	return *ret0, err
}

// GetRebateWalletsFromIds is a free data retrieval call binding the contract method 0x0a3cf98e.
//
// Solidity: function getRebateWalletsFromIds(bytes32[] reserveIds) view returns(address[] rebateWallets)
func (_KyberStorage *KyberStorageSession) GetRebateWalletsFromIds(reserveIds [][32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetRebateWalletsFromIds(&_KyberStorage.CallOpts, reserveIds)
}

// GetRebateWalletsFromIds is a free data retrieval call binding the contract method 0x0a3cf98e.
//
// Solidity: function getRebateWalletsFromIds(bytes32[] reserveIds) view returns(address[] rebateWallets)
func (_KyberStorage *KyberStorageCallerSession) GetRebateWalletsFromIds(reserveIds [][32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetRebateWalletsFromIds(&_KyberStorage.CallOpts, reserveIds)
}

// GetReserveAddressesByReserveId is a free data retrieval call binding the contract method 0xe4b80c4d.
//
// Solidity: function getReserveAddressesByReserveId(bytes32 reserveId) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCaller) GetReserveAddressesByReserveId(opts *bind.CallOpts, reserveId [32]byte) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveAddressesByReserveId", reserveId)
	return *ret0, err
}

// GetReserveAddressesByReserveId is a free data retrieval call binding the contract method 0xe4b80c4d.
//
// Solidity: function getReserveAddressesByReserveId(bytes32 reserveId) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageSession) GetReserveAddressesByReserveId(reserveId [32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesByReserveId(&_KyberStorage.CallOpts, reserveId)
}

// GetReserveAddressesByReserveId is a free data retrieval call binding the contract method 0xe4b80c4d.
//
// Solidity: function getReserveAddressesByReserveId(bytes32 reserveId) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCallerSession) GetReserveAddressesByReserveId(reserveId [32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesByReserveId(&_KyberStorage.CallOpts, reserveId)
}

// GetReserveAddressesFromIds is a free data retrieval call binding the contract method 0xd84c19c7.
//
// Solidity: function getReserveAddressesFromIds(bytes32[] reserveIds) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCaller) GetReserveAddressesFromIds(opts *bind.CallOpts, reserveIds [][32]byte) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveAddressesFromIds", reserveIds)
	return *ret0, err
}

// GetReserveAddressesFromIds is a free data retrieval call binding the contract method 0xd84c19c7.
//
// Solidity: function getReserveAddressesFromIds(bytes32[] reserveIds) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageSession) GetReserveAddressesFromIds(reserveIds [][32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesFromIds(&_KyberStorage.CallOpts, reserveIds)
}

// GetReserveAddressesFromIds is a free data retrieval call binding the contract method 0xd84c19c7.
//
// Solidity: function getReserveAddressesFromIds(bytes32[] reserveIds) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCallerSession) GetReserveAddressesFromIds(reserveIds [][32]byte) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesFromIds(&_KyberStorage.CallOpts, reserveIds)
}

// GetReserveAddressesPerTokenSrc is a free data retrieval call binding the contract method 0xd5891582.
//
// Solidity: function getReserveAddressesPerTokenSrc(address token, uint256 startIndex, uint256 endIndex) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCaller) GetReserveAddressesPerTokenSrc(opts *bind.CallOpts, token common.Address, startIndex *big.Int, endIndex *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveAddressesPerTokenSrc", token, startIndex, endIndex)
	return *ret0, err
}

// GetReserveAddressesPerTokenSrc is a free data retrieval call binding the contract method 0xd5891582.
//
// Solidity: function getReserveAddressesPerTokenSrc(address token, uint256 startIndex, uint256 endIndex) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageSession) GetReserveAddressesPerTokenSrc(token common.Address, startIndex *big.Int, endIndex *big.Int) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesPerTokenSrc(&_KyberStorage.CallOpts, token, startIndex, endIndex)
}

// GetReserveAddressesPerTokenSrc is a free data retrieval call binding the contract method 0xd5891582.
//
// Solidity: function getReserveAddressesPerTokenSrc(address token, uint256 startIndex, uint256 endIndex) view returns(address[] reserveAddresses)
func (_KyberStorage *KyberStorageCallerSession) GetReserveAddressesPerTokenSrc(token common.Address, startIndex *big.Int, endIndex *big.Int) ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserveAddressesPerTokenSrc(&_KyberStorage.CallOpts, token, startIndex, endIndex)
}

// GetReserveDetailsByAddress is a free data retrieval call binding the contract method 0xf16e429b.
//
// Solidity: function getReserveDetailsByAddress(address reserve) view returns(bytes32 reserveId, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageCaller) GetReserveDetailsByAddress(opts *bind.CallOpts, reserve common.Address) (struct {
	ReserveId            [32]byte
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	ret := new(struct {
		ReserveId            [32]byte
		RebateWallet         common.Address
		ResType              uint8
		IsFeeAccountedFlag   bool
		IsEntitledRebateFlag bool
	})
	out := ret
	err := _KyberStorage.contract.Call(opts, out, "getReserveDetailsByAddress", reserve)
	return *ret, err
}

// GetReserveDetailsByAddress is a free data retrieval call binding the contract method 0xf16e429b.
//
// Solidity: function getReserveDetailsByAddress(address reserve) view returns(bytes32 reserveId, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageSession) GetReserveDetailsByAddress(reserve common.Address) (struct {
	ReserveId            [32]byte
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	return _KyberStorage.Contract.GetReserveDetailsByAddress(&_KyberStorage.CallOpts, reserve)
}

// GetReserveDetailsByAddress is a free data retrieval call binding the contract method 0xf16e429b.
//
// Solidity: function getReserveDetailsByAddress(address reserve) view returns(bytes32 reserveId, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageCallerSession) GetReserveDetailsByAddress(reserve common.Address) (struct {
	ReserveId            [32]byte
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	return _KyberStorage.Contract.GetReserveDetailsByAddress(&_KyberStorage.CallOpts, reserve)
}

// GetReserveDetailsById is a free data retrieval call binding the contract method 0x073c4c65.
//
// Solidity: function getReserveDetailsById(bytes32 reserveId) view returns(address reserveAddress, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageCaller) GetReserveDetailsById(opts *bind.CallOpts, reserveId [32]byte) (struct {
	ReserveAddress       common.Address
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	ret := new(struct {
		ReserveAddress       common.Address
		RebateWallet         common.Address
		ResType              uint8
		IsFeeAccountedFlag   bool
		IsEntitledRebateFlag bool
	})
	out := ret
	err := _KyberStorage.contract.Call(opts, out, "getReserveDetailsById", reserveId)
	return *ret, err
}

// GetReserveDetailsById is a free data retrieval call binding the contract method 0x073c4c65.
//
// Solidity: function getReserveDetailsById(bytes32 reserveId) view returns(address reserveAddress, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageSession) GetReserveDetailsById(reserveId [32]byte) (struct {
	ReserveAddress       common.Address
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	return _KyberStorage.Contract.GetReserveDetailsById(&_KyberStorage.CallOpts, reserveId)
}

// GetReserveDetailsById is a free data retrieval call binding the contract method 0x073c4c65.
//
// Solidity: function getReserveDetailsById(bytes32 reserveId) view returns(address reserveAddress, address rebateWallet, uint8 resType, bool isFeeAccountedFlag, bool isEntitledRebateFlag)
func (_KyberStorage *KyberStorageCallerSession) GetReserveDetailsById(reserveId [32]byte) (struct {
	ReserveAddress       common.Address
	RebateWallet         common.Address
	ResType              uint8
	IsFeeAccountedFlag   bool
	IsEntitledRebateFlag bool
}, error) {
	return _KyberStorage.Contract.GetReserveDetailsById(&_KyberStorage.CallOpts, reserveId)
}

// GetReserveId is a free data retrieval call binding the contract method 0x106e9a4b.
//
// Solidity: function getReserveId(address reserve) view returns(bytes32)
func (_KyberStorage *KyberStorageCaller) GetReserveId(opts *bind.CallOpts, reserve common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveId", reserve)
	return *ret0, err
}

// GetReserveId is a free data retrieval call binding the contract method 0x106e9a4b.
//
// Solidity: function getReserveId(address reserve) view returns(bytes32)
func (_KyberStorage *KyberStorageSession) GetReserveId(reserve common.Address) ([32]byte, error) {
	return _KyberStorage.Contract.GetReserveId(&_KyberStorage.CallOpts, reserve)
}

// GetReserveId is a free data retrieval call binding the contract method 0x106e9a4b.
//
// Solidity: function getReserveId(address reserve) view returns(bytes32)
func (_KyberStorage *KyberStorageCallerSession) GetReserveId(reserve common.Address) ([32]byte, error) {
	return _KyberStorage.Contract.GetReserveId(&_KyberStorage.CallOpts, reserve)
}

// GetReserveIdsFromAddresses is a free data retrieval call binding the contract method 0x33825653.
//
// Solidity: function getReserveIdsFromAddresses(address[] reserveAddresses) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCaller) GetReserveIdsFromAddresses(opts *bind.CallOpts, reserveAddresses []common.Address) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveIdsFromAddresses", reserveAddresses)
	return *ret0, err
}

// GetReserveIdsFromAddresses is a free data retrieval call binding the contract method 0x33825653.
//
// Solidity: function getReserveIdsFromAddresses(address[] reserveAddresses) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageSession) GetReserveIdsFromAddresses(reserveAddresses []common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsFromAddresses(&_KyberStorage.CallOpts, reserveAddresses)
}

// GetReserveIdsFromAddresses is a free data retrieval call binding the contract method 0x33825653.
//
// Solidity: function getReserveIdsFromAddresses(address[] reserveAddresses) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCallerSession) GetReserveIdsFromAddresses(reserveAddresses []common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsFromAddresses(&_KyberStorage.CallOpts, reserveAddresses)
}

// GetReserveIdsPerTokenDest is a free data retrieval call binding the contract method 0xa59b60e4.
//
// Solidity: function getReserveIdsPerTokenDest(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCaller) GetReserveIdsPerTokenDest(opts *bind.CallOpts, token common.Address) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveIdsPerTokenDest", token)
	return *ret0, err
}

// GetReserveIdsPerTokenDest is a free data retrieval call binding the contract method 0xa59b60e4.
//
// Solidity: function getReserveIdsPerTokenDest(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageSession) GetReserveIdsPerTokenDest(token common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsPerTokenDest(&_KyberStorage.CallOpts, token)
}

// GetReserveIdsPerTokenDest is a free data retrieval call binding the contract method 0xa59b60e4.
//
// Solidity: function getReserveIdsPerTokenDest(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCallerSession) GetReserveIdsPerTokenDest(token common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsPerTokenDest(&_KyberStorage.CallOpts, token)
}

// GetReserveIdsPerTokenSrc is a free data retrieval call binding the contract method 0x3d3dc52c.
//
// Solidity: function getReserveIdsPerTokenSrc(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCaller) GetReserveIdsPerTokenSrc(opts *bind.CallOpts, token common.Address) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserveIdsPerTokenSrc", token)
	return *ret0, err
}

// GetReserveIdsPerTokenSrc is a free data retrieval call binding the contract method 0x3d3dc52c.
//
// Solidity: function getReserveIdsPerTokenSrc(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageSession) GetReserveIdsPerTokenSrc(token common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsPerTokenSrc(&_KyberStorage.CallOpts, token)
}

// GetReserveIdsPerTokenSrc is a free data retrieval call binding the contract method 0x3d3dc52c.
//
// Solidity: function getReserveIdsPerTokenSrc(address token) view returns(bytes32[] reserveIds)
func (_KyberStorage *KyberStorageCallerSession) GetReserveIdsPerTokenSrc(token common.Address) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReserveIdsPerTokenSrc(&_KyberStorage.CallOpts, token)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(address[])
func (_KyberStorage *KyberStorageCaller) GetReserves(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReserves")
	return *ret0, err
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(address[])
func (_KyberStorage *KyberStorageSession) GetReserves() ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserves(&_KyberStorage.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(address[])
func (_KyberStorage *KyberStorageCallerSession) GetReserves() ([]common.Address, error) {
	return _KyberStorage.Contract.GetReserves(&_KyberStorage.CallOpts)
}

// GetReservesData is a free data retrieval call binding the contract method 0x50dceb74.
//
// Solidity: function getReservesData(bytes32[] reserveIds, address src, address dest) view returns(bool areAllReservesListed, bool[] feeAccountedArr, bool[] entitledRebateArr, address[] reserveAddresses)
func (_KyberStorage *KyberStorageCaller) GetReservesData(opts *bind.CallOpts, reserveIds [][32]byte, src common.Address, dest common.Address) (struct {
	AreAllReservesListed bool
	FeeAccountedArr      []bool
	EntitledRebateArr    []bool
	ReserveAddresses     []common.Address
}, error) {
	ret := new(struct {
		AreAllReservesListed bool
		FeeAccountedArr      []bool
		EntitledRebateArr    []bool
		ReserveAddresses     []common.Address
	})
	out := ret
	err := _KyberStorage.contract.Call(opts, out, "getReservesData", reserveIds, src, dest)
	return *ret, err
}

// GetReservesData is a free data retrieval call binding the contract method 0x50dceb74.
//
// Solidity: function getReservesData(bytes32[] reserveIds, address src, address dest) view returns(bool areAllReservesListed, bool[] feeAccountedArr, bool[] entitledRebateArr, address[] reserveAddresses)
func (_KyberStorage *KyberStorageSession) GetReservesData(reserveIds [][32]byte, src common.Address, dest common.Address) (struct {
	AreAllReservesListed bool
	FeeAccountedArr      []bool
	EntitledRebateArr    []bool
	ReserveAddresses     []common.Address
}, error) {
	return _KyberStorage.Contract.GetReservesData(&_KyberStorage.CallOpts, reserveIds, src, dest)
}

// GetReservesData is a free data retrieval call binding the contract method 0x50dceb74.
//
// Solidity: function getReservesData(bytes32[] reserveIds, address src, address dest) view returns(bool areAllReservesListed, bool[] feeAccountedArr, bool[] entitledRebateArr, address[] reserveAddresses)
func (_KyberStorage *KyberStorageCallerSession) GetReservesData(reserveIds [][32]byte, src common.Address, dest common.Address) (struct {
	AreAllReservesListed bool
	FeeAccountedArr      []bool
	EntitledRebateArr    []bool
	ReserveAddresses     []common.Address
}, error) {
	return _KyberStorage.Contract.GetReservesData(&_KyberStorage.CallOpts, reserveIds, src, dest)
}

// GetReservesPerType is a free data retrieval call binding the contract method 0xa03bd93a.
//
// Solidity: function getReservesPerType(uint8 resType) view returns(bytes32[])
func (_KyberStorage *KyberStorageCaller) GetReservesPerType(opts *bind.CallOpts, resType uint8) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "getReservesPerType", resType)
	return *ret0, err
}

// GetReservesPerType is a free data retrieval call binding the contract method 0xa03bd93a.
//
// Solidity: function getReservesPerType(uint8 resType) view returns(bytes32[])
func (_KyberStorage *KyberStorageSession) GetReservesPerType(resType uint8) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReservesPerType(&_KyberStorage.CallOpts, resType)
}

// GetReservesPerType is a free data retrieval call binding the contract method 0xa03bd93a.
//
// Solidity: function getReservesPerType(uint8 resType) view returns(bytes32[])
func (_KyberStorage *KyberStorageCallerSession) GetReservesPerType(resType uint8) ([][32]byte, error) {
	return _KyberStorage.Contract.GetReservesPerType(&_KyberStorage.CallOpts, resType)
}

// IsKyberProxyAdded is a free data retrieval call binding the contract method 0xaa1da48a.
//
// Solidity: function isKyberProxyAdded() view returns(bool)
func (_KyberStorage *KyberStorageCaller) IsKyberProxyAdded(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "isKyberProxyAdded")
	return *ret0, err
}

// IsKyberProxyAdded is a free data retrieval call binding the contract method 0xaa1da48a.
//
// Solidity: function isKyberProxyAdded() view returns(bool)
func (_KyberStorage *KyberStorageSession) IsKyberProxyAdded() (bool, error) {
	return _KyberStorage.Contract.IsKyberProxyAdded(&_KyberStorage.CallOpts)
}

// IsKyberProxyAdded is a free data retrieval call binding the contract method 0xaa1da48a.
//
// Solidity: function isKyberProxyAdded() view returns(bool)
func (_KyberStorage *KyberStorageCallerSession) IsKyberProxyAdded() (bool, error) {
	return _KyberStorage.Contract.IsKyberProxyAdded(&_KyberStorage.CallOpts)
}

// KyberDaoHistory is a free data retrieval call binding the contract method 0xaaaaa6c4.
//
// Solidity: function kyberDaoHistory() view returns(address)
func (_KyberStorage *KyberStorageCaller) KyberDaoHistory(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "kyberDaoHistory")
	return *ret0, err
}

// KyberDaoHistory is a free data retrieval call binding the contract method 0xaaaaa6c4.
//
// Solidity: function kyberDaoHistory() view returns(address)
func (_KyberStorage *KyberStorageSession) KyberDaoHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberDaoHistory(&_KyberStorage.CallOpts)
}

// KyberDaoHistory is a free data retrieval call binding the contract method 0xaaaaa6c4.
//
// Solidity: function kyberDaoHistory() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) KyberDaoHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberDaoHistory(&_KyberStorage.CallOpts)
}

// KyberFeeHandlerHistory is a free data retrieval call binding the contract method 0xb6cb51c2.
//
// Solidity: function kyberFeeHandlerHistory() view returns(address)
func (_KyberStorage *KyberStorageCaller) KyberFeeHandlerHistory(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "kyberFeeHandlerHistory")
	return *ret0, err
}

// KyberFeeHandlerHistory is a free data retrieval call binding the contract method 0xb6cb51c2.
//
// Solidity: function kyberFeeHandlerHistory() view returns(address)
func (_KyberStorage *KyberStorageSession) KyberFeeHandlerHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberFeeHandlerHistory(&_KyberStorage.CallOpts)
}

// KyberFeeHandlerHistory is a free data retrieval call binding the contract method 0xb6cb51c2.
//
// Solidity: function kyberFeeHandlerHistory() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) KyberFeeHandlerHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberFeeHandlerHistory(&_KyberStorage.CallOpts)
}

// KyberMatchingEngineHistory is a free data retrieval call binding the contract method 0x4821ccfc.
//
// Solidity: function kyberMatchingEngineHistory() view returns(address)
func (_KyberStorage *KyberStorageCaller) KyberMatchingEngineHistory(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "kyberMatchingEngineHistory")
	return *ret0, err
}

// KyberMatchingEngineHistory is a free data retrieval call binding the contract method 0x4821ccfc.
//
// Solidity: function kyberMatchingEngineHistory() view returns(address)
func (_KyberStorage *KyberStorageSession) KyberMatchingEngineHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberMatchingEngineHistory(&_KyberStorage.CallOpts)
}

// KyberMatchingEngineHistory is a free data retrieval call binding the contract method 0x4821ccfc.
//
// Solidity: function kyberMatchingEngineHistory() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) KyberMatchingEngineHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberMatchingEngineHistory(&_KyberStorage.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberStorage *KyberStorageCaller) KyberNetwork(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "kyberNetwork")
	return *ret0, err
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberStorage *KyberStorageSession) KyberNetwork() (common.Address, error) {
	return _KyberStorage.Contract.KyberNetwork(&_KyberStorage.CallOpts)
}

// KyberNetwork is a free data retrieval call binding the contract method 0xb78b842d.
//
// Solidity: function kyberNetwork() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) KyberNetwork() (common.Address, error) {
	return _KyberStorage.Contract.KyberNetwork(&_KyberStorage.CallOpts)
}

// KyberNetworkHistory is a free data retrieval call binding the contract method 0x6f795fa3.
//
// Solidity: function kyberNetworkHistory() view returns(address)
func (_KyberStorage *KyberStorageCaller) KyberNetworkHistory(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "kyberNetworkHistory")
	return *ret0, err
}

// KyberNetworkHistory is a free data retrieval call binding the contract method 0x6f795fa3.
//
// Solidity: function kyberNetworkHistory() view returns(address)
func (_KyberStorage *KyberStorageSession) KyberNetworkHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberNetworkHistory(&_KyberStorage.CallOpts)
}

// KyberNetworkHistory is a free data retrieval call binding the contract method 0x6f795fa3.
//
// Solidity: function kyberNetworkHistory() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) KyberNetworkHistory() (common.Address, error) {
	return _KyberStorage.Contract.KyberNetworkHistory(&_KyberStorage.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberStorage *KyberStorageCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberStorage.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberStorage *KyberStorageSession) PendingAdmin() (common.Address, error) {
	return _KyberStorage.Contract.PendingAdmin(&_KyberStorage.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberStorage *KyberStorageCallerSession) PendingAdmin() (common.Address, error) {
	return _KyberStorage.Contract.PendingAdmin(&_KyberStorage.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberStorage *KyberStorageTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberStorage *KyberStorageSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddAlerter(&_KyberStorage.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberStorage *KyberStorageTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddAlerter(&_KyberStorage.TransactOpts, newAlerter)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x07f442e1.
//
// Solidity: function addKyberProxy(address kyberProxy, uint256 maxApprovedProxies) returns()
func (_KyberStorage *KyberStorageTransactor) AddKyberProxy(opts *bind.TransactOpts, kyberProxy common.Address, maxApprovedProxies *big.Int) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "addKyberProxy", kyberProxy, maxApprovedProxies)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x07f442e1.
//
// Solidity: function addKyberProxy(address kyberProxy, uint256 maxApprovedProxies) returns()
func (_KyberStorage *KyberStorageSession) AddKyberProxy(kyberProxy common.Address, maxApprovedProxies *big.Int) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddKyberProxy(&_KyberStorage.TransactOpts, kyberProxy, maxApprovedProxies)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x07f442e1.
//
// Solidity: function addKyberProxy(address kyberProxy, uint256 maxApprovedProxies) returns()
func (_KyberStorage *KyberStorageTransactorSession) AddKyberProxy(kyberProxy common.Address, maxApprovedProxies *big.Int) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddKyberProxy(&_KyberStorage.TransactOpts, kyberProxy, maxApprovedProxies)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberStorage *KyberStorageTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberStorage *KyberStorageSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddOperator(&_KyberStorage.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberStorage *KyberStorageTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddOperator(&_KyberStorage.TransactOpts, newOperator)
}

// AddReserve is a paid mutator transaction binding the contract method 0x3dfce895.
//
// Solidity: function addReserve(address reserve, bytes32 reserveId, uint8 resType, address rebateWallet) returns()
func (_KyberStorage *KyberStorageTransactor) AddReserve(opts *bind.TransactOpts, reserve common.Address, reserveId [32]byte, resType uint8, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "addReserve", reserve, reserveId, resType, rebateWallet)
}

// AddReserve is a paid mutator transaction binding the contract method 0x3dfce895.
//
// Solidity: function addReserve(address reserve, bytes32 reserveId, uint8 resType, address rebateWallet) returns()
func (_KyberStorage *KyberStorageSession) AddReserve(reserve common.Address, reserveId [32]byte, resType uint8, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddReserve(&_KyberStorage.TransactOpts, reserve, reserveId, resType, rebateWallet)
}

// AddReserve is a paid mutator transaction binding the contract method 0x3dfce895.
//
// Solidity: function addReserve(address reserve, bytes32 reserveId, uint8 resType, address rebateWallet) returns()
func (_KyberStorage *KyberStorageTransactorSession) AddReserve(reserve common.Address, reserveId [32]byte, resType uint8, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.AddReserve(&_KyberStorage.TransactOpts, reserve, reserveId, resType, rebateWallet)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberStorage *KyberStorageTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberStorage *KyberStorageSession) ClaimAdmin() (*types.Transaction, error) {
	return _KyberStorage.Contract.ClaimAdmin(&_KyberStorage.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberStorage *KyberStorageTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _KyberStorage.Contract.ClaimAdmin(&_KyberStorage.TransactOpts)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0x65a04d6d.
//
// Solidity: function listPairForReserve(bytes32 reserveId, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_KyberStorage *KyberStorageTransactor) ListPairForReserve(opts *bind.TransactOpts, reserveId [32]byte, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "listPairForReserve", reserveId, token, ethToToken, tokenToEth, add)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0x65a04d6d.
//
// Solidity: function listPairForReserve(bytes32 reserveId, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_KyberStorage *KyberStorageSession) ListPairForReserve(reserveId [32]byte, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.ListPairForReserve(&_KyberStorage.TransactOpts, reserveId, token, ethToToken, tokenToEth, add)
}

// ListPairForReserve is a paid mutator transaction binding the contract method 0x65a04d6d.
//
// Solidity: function listPairForReserve(bytes32 reserveId, address token, bool ethToToken, bool tokenToEth, bool add) returns()
func (_KyberStorage *KyberStorageTransactorSession) ListPairForReserve(reserveId [32]byte, token common.Address, ethToToken bool, tokenToEth bool, add bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.ListPairForReserve(&_KyberStorage.TransactOpts, reserveId, token, ethToToken, tokenToEth, add)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberStorage *KyberStorageTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberStorage *KyberStorageSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveAlerter(&_KyberStorage.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberStorage *KyberStorageTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveAlerter(&_KyberStorage.TransactOpts, alerter)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberStorage *KyberStorageTransactor) RemoveKyberProxy(opts *bind.TransactOpts, kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "removeKyberProxy", kyberProxy)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberStorage *KyberStorageSession) RemoveKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveKyberProxy(&_KyberStorage.TransactOpts, kyberProxy)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberStorage *KyberStorageTransactorSession) RemoveKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveKyberProxy(&_KyberStorage.TransactOpts, kyberProxy)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberStorage *KyberStorageTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberStorage *KyberStorageSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveOperator(&_KyberStorage.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberStorage *KyberStorageTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveOperator(&_KyberStorage.TransactOpts, operator)
}

// RemoveReserve is a paid mutator transaction binding the contract method 0x2185340b.
//
// Solidity: function removeReserve(bytes32 reserveId, uint256 startIndex) returns()
func (_KyberStorage *KyberStorageTransactor) RemoveReserve(opts *bind.TransactOpts, reserveId [32]byte, startIndex *big.Int) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "removeReserve", reserveId, startIndex)
}

// RemoveReserve is a paid mutator transaction binding the contract method 0x2185340b.
//
// Solidity: function removeReserve(bytes32 reserveId, uint256 startIndex) returns()
func (_KyberStorage *KyberStorageSession) RemoveReserve(reserveId [32]byte, startIndex *big.Int) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveReserve(&_KyberStorage.TransactOpts, reserveId, startIndex)
}

// RemoveReserve is a paid mutator transaction binding the contract method 0x2185340b.
//
// Solidity: function removeReserve(bytes32 reserveId, uint256 startIndex) returns()
func (_KyberStorage *KyberStorageTransactorSession) RemoveReserve(reserveId [32]byte, startIndex *big.Int) (*types.Transaction, error) {
	return _KyberStorage.Contract.RemoveReserve(&_KyberStorage.TransactOpts, reserveId, startIndex)
}

// SetContracts is a paid mutator transaction binding the contract method 0xd8952a49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine) returns()
func (_KyberStorage *KyberStorageTransactor) SetContracts(opts *bind.TransactOpts, _kyberFeeHandler common.Address, _kyberMatchingEngine common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setContracts", _kyberFeeHandler, _kyberMatchingEngine)
}

// SetContracts is a paid mutator transaction binding the contract method 0xd8952a49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine) returns()
func (_KyberStorage *KyberStorageSession) SetContracts(_kyberFeeHandler common.Address, _kyberMatchingEngine common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetContracts(&_KyberStorage.TransactOpts, _kyberFeeHandler, _kyberMatchingEngine)
}

// SetContracts is a paid mutator transaction binding the contract method 0xd8952a49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetContracts(_kyberFeeHandler common.Address, _kyberMatchingEngine common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetContracts(&_KyberStorage.TransactOpts, _kyberFeeHandler, _kyberMatchingEngine)
}

// SetEntitledRebatePerReserveType is a paid mutator transaction binding the contract method 0x13d47a8d.
//
// Solidity: function setEntitledRebatePerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageTransactor) SetEntitledRebatePerReserveType(opts *bind.TransactOpts, fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setEntitledRebatePerReserveType", fpr, apr, bridge, utility, custom, orderbook)
}

// SetEntitledRebatePerReserveType is a paid mutator transaction binding the contract method 0x13d47a8d.
//
// Solidity: function setEntitledRebatePerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageSession) SetEntitledRebatePerReserveType(fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetEntitledRebatePerReserveType(&_KyberStorage.TransactOpts, fpr, apr, bridge, utility, custom, orderbook)
}

// SetEntitledRebatePerReserveType is a paid mutator transaction binding the contract method 0x13d47a8d.
//
// Solidity: function setEntitledRebatePerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetEntitledRebatePerReserveType(fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetEntitledRebatePerReserveType(&_KyberStorage.TransactOpts, fpr, apr, bridge, utility, custom, orderbook)
}

// SetFeeAccountedPerReserveType is a paid mutator transaction binding the contract method 0x45bd576c.
//
// Solidity: function setFeeAccountedPerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageTransactor) SetFeeAccountedPerReserveType(opts *bind.TransactOpts, fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setFeeAccountedPerReserveType", fpr, apr, bridge, utility, custom, orderbook)
}

// SetFeeAccountedPerReserveType is a paid mutator transaction binding the contract method 0x45bd576c.
//
// Solidity: function setFeeAccountedPerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageSession) SetFeeAccountedPerReserveType(fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetFeeAccountedPerReserveType(&_KyberStorage.TransactOpts, fpr, apr, bridge, utility, custom, orderbook)
}

// SetFeeAccountedPerReserveType is a paid mutator transaction binding the contract method 0x45bd576c.
//
// Solidity: function setFeeAccountedPerReserveType(bool fpr, bool apr, bool bridge, bool utility, bool custom, bool orderbook) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetFeeAccountedPerReserveType(fpr bool, apr bool, bridge bool, utility bool, custom bool, orderbook bool) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetFeeAccountedPerReserveType(&_KyberStorage.TransactOpts, fpr, apr, bridge, utility, custom, orderbook)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberStorage *KyberStorageTransactor) SetKyberDaoContract(opts *bind.TransactOpts, _kyberDao common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setKyberDaoContract", _kyberDao)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberStorage *KyberStorageSession) SetKyberDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetKyberDaoContract(&_KyberStorage.TransactOpts, _kyberDao)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetKyberDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetKyberDaoContract(&_KyberStorage.TransactOpts, _kyberDao)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberStorage *KyberStorageTransactor) SetNetworkContract(opts *bind.TransactOpts, _kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setNetworkContract", _kyberNetwork)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberStorage *KyberStorageSession) SetNetworkContract(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetNetworkContract(&_KyberStorage.TransactOpts, _kyberNetwork)
}

// SetNetworkContract is a paid mutator transaction binding the contract method 0x599b9348.
//
// Solidity: function setNetworkContract(address _kyberNetwork) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetNetworkContract(_kyberNetwork common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetNetworkContract(&_KyberStorage.TransactOpts, _kyberNetwork)
}

// SetRebateWallet is a paid mutator transaction binding the contract method 0xd95c2730.
//
// Solidity: function setRebateWallet(bytes32 reserveId, address rebateWallet) returns()
func (_KyberStorage *KyberStorageTransactor) SetRebateWallet(opts *bind.TransactOpts, reserveId [32]byte, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "setRebateWallet", reserveId, rebateWallet)
}

// SetRebateWallet is a paid mutator transaction binding the contract method 0xd95c2730.
//
// Solidity: function setRebateWallet(bytes32 reserveId, address rebateWallet) returns()
func (_KyberStorage *KyberStorageSession) SetRebateWallet(reserveId [32]byte, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetRebateWallet(&_KyberStorage.TransactOpts, reserveId, rebateWallet)
}

// SetRebateWallet is a paid mutator transaction binding the contract method 0xd95c2730.
//
// Solidity: function setRebateWallet(bytes32 reserveId, address rebateWallet) returns()
func (_KyberStorage *KyberStorageTransactorSession) SetRebateWallet(reserveId [32]byte, rebateWallet common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.SetRebateWallet(&_KyberStorage.TransactOpts, reserveId, rebateWallet)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberStorage *KyberStorageTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberStorage *KyberStorageSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.TransferAdmin(&_KyberStorage.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberStorage *KyberStorageTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.TransferAdmin(&_KyberStorage.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberStorage *KyberStorageTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberStorage *KyberStorageSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.TransferAdminQuickly(&_KyberStorage.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberStorage *KyberStorageTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberStorage.Contract.TransferAdminQuickly(&_KyberStorage.TransactOpts, newAdmin)
}

// KyberStorageAddReserveToStorageIterator is returned from FilterAddReserveToStorage and is used to iterate over the raw logs and unpacked data for AddReserveToStorage events raised by the KyberStorage contract.
type KyberStorageAddReserveToStorageIterator struct {
	Event *KyberStorageAddReserveToStorage // Event containing the contract specifics and raw log

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
func (it *KyberStorageAddReserveToStorageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageAddReserveToStorage)
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
		it.Event = new(KyberStorageAddReserveToStorage)
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
func (it *KyberStorageAddReserveToStorageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageAddReserveToStorageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageAddReserveToStorage represents a AddReserveToStorage event raised by the KyberStorage contract.
type KyberStorageAddReserveToStorage struct {
	Reserve      common.Address
	ReserveId    [32]byte
	ReserveType  uint8
	RebateWallet common.Address
	Add          bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAddReserveToStorage is a free log retrieval operation binding the contract event 0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004.
//
// Solidity: event AddReserveToStorage(address indexed reserve, bytes32 indexed reserveId, uint8 reserveType, address indexed rebateWallet, bool add)
func (_KyberStorage *KyberStorageFilterer) FilterAddReserveToStorage(opts *bind.FilterOpts, reserve []common.Address, reserveId [][32]byte, rebateWallet []common.Address) (*KyberStorageAddReserveToStorageIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "AddReserveToStorage", reserveRule, reserveIdRule, rebateWalletRule)
	if err != nil {
		return nil, err
	}
	return &KyberStorageAddReserveToStorageIterator{contract: _KyberStorage.contract, event: "AddReserveToStorage", logs: logs, sub: sub}, nil
}

// WatchAddReserveToStorage is a free log subscription operation binding the contract event 0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004.
//
// Solidity: event AddReserveToStorage(address indexed reserve, bytes32 indexed reserveId, uint8 reserveType, address indexed rebateWallet, bool add)
func (_KyberStorage *KyberStorageFilterer) WatchAddReserveToStorage(opts *bind.WatchOpts, sink chan<- *KyberStorageAddReserveToStorage, reserve []common.Address, reserveId [][32]byte, rebateWallet []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "AddReserveToStorage", reserveRule, reserveIdRule, rebateWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageAddReserveToStorage)
				if err := _KyberStorage.contract.UnpackLog(event, "AddReserveToStorage", log); err != nil {
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

// ParseAddReserveToStorage is a log parse operation binding the contract event 0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004.
//
// Solidity: event AddReserveToStorage(address indexed reserve, bytes32 indexed reserveId, uint8 reserveType, address indexed rebateWallet, bool add)
func (_KyberStorage *KyberStorageFilterer) ParseAddReserveToStorage(log types.Log) (*KyberStorageAddReserveToStorage, error) {
	event := new(KyberStorageAddReserveToStorage)
	if err := _KyberStorage.contract.UnpackLog(event, "AddReserveToStorage", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the KyberStorage contract.
type KyberStorageAdminClaimedIterator struct {
	Event *KyberStorageAdminClaimed // Event containing the contract specifics and raw log

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
func (it *KyberStorageAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageAdminClaimed)
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
		it.Event = new(KyberStorageAdminClaimed)
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
func (it *KyberStorageAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageAdminClaimed represents a AdminClaimed event raised by the KyberStorage contract.
type KyberStorageAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberStorage *KyberStorageFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*KyberStorageAdminClaimedIterator, error) {

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &KyberStorageAdminClaimedIterator{contract: _KyberStorage.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberStorage *KyberStorageFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *KyberStorageAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageAdminClaimed)
				if err := _KyberStorage.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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
func (_KyberStorage *KyberStorageFilterer) ParseAdminClaimed(log types.Log) (*KyberStorageAdminClaimed, error) {
	event := new(KyberStorageAdminClaimed)
	if err := _KyberStorage.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the KyberStorage contract.
type KyberStorageAlerterAddedIterator struct {
	Event *KyberStorageAlerterAdded // Event containing the contract specifics and raw log

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
func (it *KyberStorageAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageAlerterAdded)
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
		it.Event = new(KyberStorageAlerterAdded)
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
func (it *KyberStorageAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageAlerterAdded represents a AlerterAdded event raised by the KyberStorage contract.
type KyberStorageAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberStorage *KyberStorageFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*KyberStorageAlerterAddedIterator, error) {

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &KyberStorageAlerterAddedIterator{contract: _KyberStorage.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberStorage *KyberStorageFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *KyberStorageAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageAlerterAdded)
				if err := _KyberStorage.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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
func (_KyberStorage *KyberStorageFilterer) ParseAlerterAdded(log types.Log) (*KyberStorageAlerterAdded, error) {
	event := new(KyberStorageAlerterAdded)
	if err := _KyberStorage.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageKyberNetworkUpdatedIterator is returned from FilterKyberNetworkUpdated and is used to iterate over the raw logs and unpacked data for KyberNetworkUpdated events raised by the KyberStorage contract.
type KyberStorageKyberNetworkUpdatedIterator struct {
	Event *KyberStorageKyberNetworkUpdated // Event containing the contract specifics and raw log

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
func (it *KyberStorageKyberNetworkUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageKyberNetworkUpdated)
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
		it.Event = new(KyberStorageKyberNetworkUpdated)
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
func (it *KyberStorageKyberNetworkUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageKyberNetworkUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageKyberNetworkUpdated represents a KyberNetworkUpdated event raised by the KyberStorage contract.
type KyberStorageKyberNetworkUpdated struct {
	NewKyberNetwork common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkUpdated is a free log retrieval operation binding the contract event 0x18970d46ac8a7d7e0da90e1bebb0be3e87ffc7705fc09d3bba5373d59b7a12aa.
//
// Solidity: event KyberNetworkUpdated(address newKyberNetwork)
func (_KyberStorage *KyberStorageFilterer) FilterKyberNetworkUpdated(opts *bind.FilterOpts) (*KyberStorageKyberNetworkUpdatedIterator, error) {

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "KyberNetworkUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberStorageKyberNetworkUpdatedIterator{contract: _KyberStorage.contract, event: "KyberNetworkUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkUpdated is a free log subscription operation binding the contract event 0x18970d46ac8a7d7e0da90e1bebb0be3e87ffc7705fc09d3bba5373d59b7a12aa.
//
// Solidity: event KyberNetworkUpdated(address newKyberNetwork)
func (_KyberStorage *KyberStorageFilterer) WatchKyberNetworkUpdated(opts *bind.WatchOpts, sink chan<- *KyberStorageKyberNetworkUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "KyberNetworkUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageKyberNetworkUpdated)
				if err := _KyberStorage.contract.UnpackLog(event, "KyberNetworkUpdated", log); err != nil {
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
// Solidity: event KyberNetworkUpdated(address newKyberNetwork)
func (_KyberStorage *KyberStorageFilterer) ParseKyberNetworkUpdated(log types.Log) (*KyberStorageKyberNetworkUpdated, error) {
	event := new(KyberStorageKyberNetworkUpdated)
	if err := _KyberStorage.contract.UnpackLog(event, "KyberNetworkUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageListReservePairsIterator is returned from FilterListReservePairs and is used to iterate over the raw logs and unpacked data for ListReservePairs events raised by the KyberStorage contract.
type KyberStorageListReservePairsIterator struct {
	Event *KyberStorageListReservePairs // Event containing the contract specifics and raw log

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
func (it *KyberStorageListReservePairsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageListReservePairs)
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
		it.Event = new(KyberStorageListReservePairs)
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
func (it *KyberStorageListReservePairsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageListReservePairsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageListReservePairs represents a ListReservePairs event raised by the KyberStorage contract.
type KyberStorageListReservePairs struct {
	ReserveId [32]byte
	Reserve   common.Address
	Src       common.Address
	Dest      common.Address
	Add       bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterListReservePairs is a free log retrieval operation binding the contract event 0xfcdbd685961328a43b4aa133de257d3769cc01891b4ee00fd5058e5aa3564ca5.
//
// Solidity: event ListReservePairs(bytes32 indexed reserveId, address reserve, address indexed src, address indexed dest, bool add)
func (_KyberStorage *KyberStorageFilterer) FilterListReservePairs(opts *bind.FilterOpts, reserveId [][32]byte, src []common.Address, dest []common.Address) (*KyberStorageListReservePairsIterator, error) {

	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "ListReservePairs", reserveIdRule, srcRule, destRule)
	if err != nil {
		return nil, err
	}
	return &KyberStorageListReservePairsIterator{contract: _KyberStorage.contract, event: "ListReservePairs", logs: logs, sub: sub}, nil
}

// WatchListReservePairs is a free log subscription operation binding the contract event 0xfcdbd685961328a43b4aa133de257d3769cc01891b4ee00fd5058e5aa3564ca5.
//
// Solidity: event ListReservePairs(bytes32 indexed reserveId, address reserve, address indexed src, address indexed dest, bool add)
func (_KyberStorage *KyberStorageFilterer) WatchListReservePairs(opts *bind.WatchOpts, sink chan<- *KyberStorageListReservePairs, reserveId [][32]byte, src []common.Address, dest []common.Address) (event.Subscription, error) {

	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "ListReservePairs", reserveIdRule, srcRule, destRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageListReservePairs)
				if err := _KyberStorage.contract.UnpackLog(event, "ListReservePairs", log); err != nil {
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

// ParseListReservePairs is a log parse operation binding the contract event 0xfcdbd685961328a43b4aa133de257d3769cc01891b4ee00fd5058e5aa3564ca5.
//
// Solidity: event ListReservePairs(bytes32 indexed reserveId, address reserve, address indexed src, address indexed dest, bool add)
func (_KyberStorage *KyberStorageFilterer) ParseListReservePairs(log types.Log) (*KyberStorageListReservePairs, error) {
	event := new(KyberStorageListReservePairs)
	if err := _KyberStorage.contract.UnpackLog(event, "ListReservePairs", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the KyberStorage contract.
type KyberStorageOperatorAddedIterator struct {
	Event *KyberStorageOperatorAdded // Event containing the contract specifics and raw log

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
func (it *KyberStorageOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageOperatorAdded)
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
		it.Event = new(KyberStorageOperatorAdded)
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
func (it *KyberStorageOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageOperatorAdded represents a OperatorAdded event raised by the KyberStorage contract.
type KyberStorageOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberStorage *KyberStorageFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*KyberStorageOperatorAddedIterator, error) {

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &KyberStorageOperatorAddedIterator{contract: _KyberStorage.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberStorage *KyberStorageFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *KyberStorageOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageOperatorAdded)
				if err := _KyberStorage.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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
func (_KyberStorage *KyberStorageFilterer) ParseOperatorAdded(log types.Log) (*KyberStorageOperatorAdded, error) {
	event := new(KyberStorageOperatorAdded)
	if err := _KyberStorage.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageRemoveReserveFromStorageIterator is returned from FilterRemoveReserveFromStorage and is used to iterate over the raw logs and unpacked data for RemoveReserveFromStorage events raised by the KyberStorage contract.
type KyberStorageRemoveReserveFromStorageIterator struct {
	Event *KyberStorageRemoveReserveFromStorage // Event containing the contract specifics and raw log

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
func (it *KyberStorageRemoveReserveFromStorageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageRemoveReserveFromStorage)
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
		it.Event = new(KyberStorageRemoveReserveFromStorage)
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
func (it *KyberStorageRemoveReserveFromStorageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageRemoveReserveFromStorageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageRemoveReserveFromStorage represents a RemoveReserveFromStorage event raised by the KyberStorage contract.
type KyberStorageRemoveReserveFromStorage struct {
	Reserve   common.Address
	ReserveId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRemoveReserveFromStorage is a free log retrieval operation binding the contract event 0xa5cd88a226efb041d6bdc0ac32964affd749b8a7c4d9e0c4ffba575e7180b1c9.
//
// Solidity: event RemoveReserveFromStorage(address indexed reserve, bytes32 indexed reserveId)
func (_KyberStorage *KyberStorageFilterer) FilterRemoveReserveFromStorage(opts *bind.FilterOpts, reserve []common.Address, reserveId [][32]byte) (*KyberStorageRemoveReserveFromStorageIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "RemoveReserveFromStorage", reserveRule, reserveIdRule)
	if err != nil {
		return nil, err
	}
	return &KyberStorageRemoveReserveFromStorageIterator{contract: _KyberStorage.contract, event: "RemoveReserveFromStorage", logs: logs, sub: sub}, nil
}

// WatchRemoveReserveFromStorage is a free log subscription operation binding the contract event 0xa5cd88a226efb041d6bdc0ac32964affd749b8a7c4d9e0c4ffba575e7180b1c9.
//
// Solidity: event RemoveReserveFromStorage(address indexed reserve, bytes32 indexed reserveId)
func (_KyberStorage *KyberStorageFilterer) WatchRemoveReserveFromStorage(opts *bind.WatchOpts, sink chan<- *KyberStorageRemoveReserveFromStorage, reserve []common.Address, reserveId [][32]byte) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "RemoveReserveFromStorage", reserveRule, reserveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageRemoveReserveFromStorage)
				if err := _KyberStorage.contract.UnpackLog(event, "RemoveReserveFromStorage", log); err != nil {
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

// ParseRemoveReserveFromStorage is a log parse operation binding the contract event 0xa5cd88a226efb041d6bdc0ac32964affd749b8a7c4d9e0c4ffba575e7180b1c9.
//
// Solidity: event RemoveReserveFromStorage(address indexed reserve, bytes32 indexed reserveId)
func (_KyberStorage *KyberStorageFilterer) ParseRemoveReserveFromStorage(log types.Log) (*KyberStorageRemoveReserveFromStorage, error) {
	event := new(KyberStorageRemoveReserveFromStorage)
	if err := _KyberStorage.contract.UnpackLog(event, "RemoveReserveFromStorage", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageReserveRebateWalletSetIterator is returned from FilterReserveRebateWalletSet and is used to iterate over the raw logs and unpacked data for ReserveRebateWalletSet events raised by the KyberStorage contract.
type KyberStorageReserveRebateWalletSetIterator struct {
	Event *KyberStorageReserveRebateWalletSet // Event containing the contract specifics and raw log

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
func (it *KyberStorageReserveRebateWalletSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageReserveRebateWalletSet)
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
		it.Event = new(KyberStorageReserveRebateWalletSet)
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
func (it *KyberStorageReserveRebateWalletSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageReserveRebateWalletSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageReserveRebateWalletSet represents a ReserveRebateWalletSet event raised by the KyberStorage contract.
type KyberStorageReserveRebateWalletSet struct {
	ReserveId    [32]byte
	RebateWallet common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterReserveRebateWalletSet is a free log retrieval operation binding the contract event 0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3.
//
// Solidity: event ReserveRebateWalletSet(bytes32 indexed reserveId, address indexed rebateWallet)
func (_KyberStorage *KyberStorageFilterer) FilterReserveRebateWalletSet(opts *bind.FilterOpts, reserveId [][32]byte, rebateWallet []common.Address) (*KyberStorageReserveRebateWalletSetIterator, error) {

	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}
	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "ReserveRebateWalletSet", reserveIdRule, rebateWalletRule)
	if err != nil {
		return nil, err
	}
	return &KyberStorageReserveRebateWalletSetIterator{contract: _KyberStorage.contract, event: "ReserveRebateWalletSet", logs: logs, sub: sub}, nil
}

// WatchReserveRebateWalletSet is a free log subscription operation binding the contract event 0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3.
//
// Solidity: event ReserveRebateWalletSet(bytes32 indexed reserveId, address indexed rebateWallet)
func (_KyberStorage *KyberStorageFilterer) WatchReserveRebateWalletSet(opts *bind.WatchOpts, sink chan<- *KyberStorageReserveRebateWalletSet, reserveId [][32]byte, rebateWallet []common.Address) (event.Subscription, error) {

	var reserveIdRule []interface{}
	for _, reserveIdItem := range reserveId {
		reserveIdRule = append(reserveIdRule, reserveIdItem)
	}
	var rebateWalletRule []interface{}
	for _, rebateWalletItem := range rebateWallet {
		rebateWalletRule = append(rebateWalletRule, rebateWalletItem)
	}

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "ReserveRebateWalletSet", reserveIdRule, rebateWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageReserveRebateWalletSet)
				if err := _KyberStorage.contract.UnpackLog(event, "ReserveRebateWalletSet", log); err != nil {
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

// ParseReserveRebateWalletSet is a log parse operation binding the contract event 0x42cac9e63e37f62d5689493d04887a67fe3c68e1d3763c3f0890e1620a0465b3.
//
// Solidity: event ReserveRebateWalletSet(bytes32 indexed reserveId, address indexed rebateWallet)
func (_KyberStorage *KyberStorageFilterer) ParseReserveRebateWalletSet(log types.Log) (*KyberStorageReserveRebateWalletSet, error) {
	event := new(KyberStorageReserveRebateWalletSet)
	if err := _KyberStorage.contract.UnpackLog(event, "ReserveRebateWalletSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberStorageTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the KyberStorage contract.
type KyberStorageTransferAdminPendingIterator struct {
	Event *KyberStorageTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *KyberStorageTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberStorageTransferAdminPending)
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
		it.Event = new(KyberStorageTransferAdminPending)
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
func (it *KyberStorageTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberStorageTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberStorageTransferAdminPending represents a TransferAdminPending event raised by the KyberStorage contract.
type KyberStorageTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberStorage *KyberStorageFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*KyberStorageTransferAdminPendingIterator, error) {

	logs, sub, err := _KyberStorage.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &KyberStorageTransferAdminPendingIterator{contract: _KyberStorage.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberStorage *KyberStorageFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *KyberStorageTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _KyberStorage.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberStorageTransferAdminPending)
				if err := _KyberStorage.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
func (_KyberStorage *KyberStorageFilterer) ParseTransferAdminPending(log types.Log) (*KyberStorageTransferAdminPending, error) {
	event := new(KyberStorageTransferAdminPending)
	if err := _KyberStorage.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
		return nil, err
	}
	return event, nil
}
