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

// KyberNetworkABI is the input ABI used to generate the binding from.
const KyberNetworkABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"contractIKyberStorage\",\"name\":\"_kyberStorage\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"}],\"name\":\"AdminClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"AlerterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EtherReceival\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"EtherWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIGasHelper\",\"name\":\"newGasHelper\",\"type\":\"address\"}],\"name\":\"GasHelperUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberDao\",\"name\":\"newKyberDao\",\"type\":\"address\"}],\"name\":\"KyberDaoUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberFeeHandler\",\"name\":\"newKyberFeeHandler\",\"type\":\"address\"}],\"name\":\"KyberFeeHandlerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIKyberMatchingEngine\",\"name\":\"newKyberMatchingEngine\",\"type\":\"address\"}],\"name\":\"KyberMatchingEngineUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxGasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"negligibleRateDiffBps\",\"type\":\"uint256\"}],\"name\":\"KyberNetworkParamsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"name\":\"KyberNetworkSetEnable\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"KyberProxyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"KyberProxyRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethWeiValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"networkFeeWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"customPlatformFeeWei\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"t2eIds\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"e2tIds\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"t2eSrcAmounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"e2tSrcAmounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"t2eRates\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"e2tRates\",\"type\":\"uint256[]\"}],\"name\":\"KyberTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"reserves\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"ListedReservesForToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isAdd\",\"type\":\"bool\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"TokenWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminPending\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAlerter\",\"type\":\"address\"}],\"name\":\"addAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"addKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAlerters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAndUpdateNetworkFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"networkFeeBps\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContracts\",\"outputs\":[{\"internalType\":\"contractIKyberFeeHandler\",\"name\":\"kyberFeeHandlerAddress\",\"type\":\"address\"},{\"internalType\":\"contractIKyberDao\",\"name\":\"kyberDaoAddress\",\"type\":\"address\"},{\"internalType\":\"contractIKyberMatchingEngine\",\"name\":\"kyberMatchingEngineAddress\",\"type\":\"address\"},{\"internalType\":\"contractIKyberStorage\",\"name\":\"kyberStorageAddress\",\"type\":\"address\"},{\"internalType\":\"contractIGasHelper\",\"name\":\"gasHelperAddress\",\"type\":\"address\"},{\"internalType\":\"contractIKyberNetworkProxy[]\",\"name\":\"kyberProxyAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcQty\",\"type\":\"uint256\"}],\"name\":\"getExpectedRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"worstRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcQty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"platformFeeBps\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"getExpectedRateWithHintAndFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rateWithNetworkFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rateWithAllFees\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNetworkData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"negligibleDiffBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"networkFeeBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"listReservesForToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"add\",\"type\":\"bool\"}],\"name\":\"listTokenForReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"alerter\",\"type\":\"address\"}],\"name\":\"removeAlerter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"kyberProxy\",\"type\":\"address\"}],\"name\":\"removeKyberProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberFeeHandler\",\"name\":\"_kyberFeeHandler\",\"type\":\"address\"},{\"internalType\":\"contractIKyberMatchingEngine\",\"name\":\"_kyberMatchingEngine\",\"type\":\"address\"},{\"internalType\":\"contractIGasHelper\",\"name\":\"_gasHelper\",\"type\":\"address\"}],\"name\":\"setContracts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"setEnable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIKyberDao\",\"name\":\"_kyberDao\",\"type\":\"address\"}],\"name\":\"setKyberDaoContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_negligibleRateDiffBps\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"trader\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"walletId\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"destAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"trader\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"srcAmount\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxDestAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minConversionRate\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"platformWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"platformFeeBps\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"hint\",\"type\":\"bytes\"}],\"name\":\"tradeWithHintAndFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"destAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminQuickly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawEther\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sendTo\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// KyberNetwork is an auto generated Go binding around an Ethereum contract.
type KyberNetwork struct {
	KyberNetworkCaller     // Read-only binding to the contract
	KyberNetworkTransactor // Write-only binding to the contract
	KyberNetworkFilterer   // Log filterer for contract events
}

// KyberNetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type KyberNetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KyberNetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KyberNetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KyberNetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KyberNetworkSession struct {
	Contract     *KyberNetwork     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KyberNetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KyberNetworkCallerSession struct {
	Contract *KyberNetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// KyberNetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KyberNetworkTransactorSession struct {
	Contract     *KyberNetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// KyberNetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type KyberNetworkRaw struct {
	Contract *KyberNetwork // Generic contract binding to access the raw methods on
}

// KyberNetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KyberNetworkCallerRaw struct {
	Contract *KyberNetworkCaller // Generic read-only contract binding to access the raw methods on
}

// KyberNetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KyberNetworkTransactorRaw struct {
	Contract *KyberNetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKyberNetwork creates a new instance of KyberNetwork, bound to a specific deployed contract.
func NewKyberNetwork(address common.Address, backend bind.ContractBackend) (*KyberNetwork, error) {
	contract, err := bindKyberNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KyberNetwork{KyberNetworkCaller: KyberNetworkCaller{contract: contract}, KyberNetworkTransactor: KyberNetworkTransactor{contract: contract}, KyberNetworkFilterer: KyberNetworkFilterer{contract: contract}}, nil
}

// NewKyberNetworkCaller creates a new read-only instance of KyberNetwork, bound to a specific deployed contract.
func NewKyberNetworkCaller(address common.Address, caller bind.ContractCaller) (*KyberNetworkCaller, error) {
	contract, err := bindKyberNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkCaller{contract: contract}, nil
}

// NewKyberNetworkTransactor creates a new write-only instance of KyberNetwork, bound to a specific deployed contract.
func NewKyberNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*KyberNetworkTransactor, error) {
	contract, err := bindKyberNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkTransactor{contract: contract}, nil
}

// NewKyberNetworkFilterer creates a new log filterer instance of KyberNetwork, bound to a specific deployed contract.
func NewKyberNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*KyberNetworkFilterer, error) {
	contract, err := bindKyberNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkFilterer{contract: contract}, nil
}

// bindKyberNetwork binds a generic wrapper to an already deployed contract.
func bindKyberNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KyberNetworkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberNetwork *KyberNetworkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberNetwork.Contract.KyberNetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberNetwork *KyberNetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetwork.Contract.KyberNetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberNetwork *KyberNetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberNetwork.Contract.KyberNetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KyberNetwork *KyberNetworkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KyberNetwork.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KyberNetwork *KyberNetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetwork.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KyberNetwork *KyberNetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KyberNetwork.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetwork *KyberNetworkCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetwork *KyberNetworkSession) Admin() (common.Address, error) {
	return _KyberNetwork.Contract.Admin(&_KyberNetwork.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_KyberNetwork *KyberNetworkCallerSession) Admin() (common.Address, error) {
	return _KyberNetwork.Contract.Admin(&_KyberNetwork.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetwork *KyberNetworkCaller) Enabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetwork *KyberNetworkSession) Enabled() (bool, error) {
	return _KyberNetwork.Contract.Enabled(&_KyberNetwork.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() view returns(bool)
func (_KyberNetwork *KyberNetworkCallerSession) Enabled() (bool, error) {
	return _KyberNetwork.Contract.Enabled(&_KyberNetwork.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetwork *KyberNetworkCaller) GetAlerters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "getAlerters")
	return *ret0, err
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetwork *KyberNetworkSession) GetAlerters() ([]common.Address, error) {
	return _KyberNetwork.Contract.GetAlerters(&_KyberNetwork.CallOpts)
}

// GetAlerters is a free data retrieval call binding the contract method 0x7c423f54.
//
// Solidity: function getAlerters() view returns(address[])
func (_KyberNetwork *KyberNetworkCallerSession) GetAlerters() ([]common.Address, error) {
	return _KyberNetwork.Contract.GetAlerters(&_KyberNetwork.CallOpts)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address kyberFeeHandlerAddress, address kyberDaoAddress, address kyberMatchingEngineAddress, address kyberStorageAddress, address gasHelperAddress, address[] kyberProxyAddresses)
func (_KyberNetwork *KyberNetworkCaller) GetContracts(opts *bind.CallOpts) (struct {
	KyberFeeHandlerAddress     common.Address
	KyberDaoAddress            common.Address
	KyberMatchingEngineAddress common.Address
	KyberStorageAddress        common.Address
	GasHelperAddress           common.Address
	KyberProxyAddresses        []common.Address
}, error) {
	ret := new(struct {
		KyberFeeHandlerAddress     common.Address
		KyberDaoAddress            common.Address
		KyberMatchingEngineAddress common.Address
		KyberStorageAddress        common.Address
		GasHelperAddress           common.Address
		KyberProxyAddresses        []common.Address
	})
	out := ret
	err := _KyberNetwork.contract.Call(opts, out, "getContracts")
	return *ret, err
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address kyberFeeHandlerAddress, address kyberDaoAddress, address kyberMatchingEngineAddress, address kyberStorageAddress, address gasHelperAddress, address[] kyberProxyAddresses)
func (_KyberNetwork *KyberNetworkSession) GetContracts() (struct {
	KyberFeeHandlerAddress     common.Address
	KyberDaoAddress            common.Address
	KyberMatchingEngineAddress common.Address
	KyberStorageAddress        common.Address
	GasHelperAddress           common.Address
	KyberProxyAddresses        []common.Address
}, error) {
	return _KyberNetwork.Contract.GetContracts(&_KyberNetwork.CallOpts)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns(address kyberFeeHandlerAddress, address kyberDaoAddress, address kyberMatchingEngineAddress, address kyberStorageAddress, address gasHelperAddress, address[] kyberProxyAddresses)
func (_KyberNetwork *KyberNetworkCallerSession) GetContracts() (struct {
	KyberFeeHandlerAddress     common.Address
	KyberDaoAddress            common.Address
	KyberMatchingEngineAddress common.Address
	KyberStorageAddress        common.Address
	GasHelperAddress           common.Address
	KyberProxyAddresses        []common.Address
}, error) {
	return _KyberNetwork.Contract.GetContracts(&_KyberNetwork.CallOpts)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetwork *KyberNetworkCaller) GetExpectedRate(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	ret := new(struct {
		ExpectedRate *big.Int
		WorstRate    *big.Int
	})
	out := ret
	err := _KyberNetwork.contract.Call(opts, out, "getExpectedRate", src, dest, srcQty)
	return *ret, err
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetwork *KyberNetworkSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	return _KyberNetwork.Contract.GetExpectedRate(&_KyberNetwork.CallOpts, src, dest, srcQty)
}

// GetExpectedRate is a free data retrieval call binding the contract method 0x809a9e55.
//
// Solidity: function getExpectedRate(address src, address dest, uint256 srcQty) view returns(uint256 expectedRate, uint256 worstRate)
func (_KyberNetwork *KyberNetworkCallerSession) GetExpectedRate(src common.Address, dest common.Address, srcQty *big.Int) (struct {
	ExpectedRate *big.Int
	WorstRate    *big.Int
}, error) {
	return _KyberNetwork.Contract.GetExpectedRate(&_KyberNetwork.CallOpts, src, dest, srcQty)
}

// GetExpectedRateWithHintAndFee is a free data retrieval call binding the contract method 0x8ff68a80.
//
// Solidity: function getExpectedRateWithHintAndFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 rateWithNetworkFee, uint256 rateWithAllFees)
func (_KyberNetwork *KyberNetworkCaller) GetExpectedRateWithHintAndFee(opts *bind.CallOpts, src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (struct {
	RateWithNetworkFee *big.Int
	RateWithAllFees    *big.Int
}, error) {
	ret := new(struct {
		RateWithNetworkFee *big.Int
		RateWithAllFees    *big.Int
	})
	out := ret
	err := _KyberNetwork.contract.Call(opts, out, "getExpectedRateWithHintAndFee", src, dest, srcQty, platformFeeBps, hint)
	return *ret, err
}

// GetExpectedRateWithHintAndFee is a free data retrieval call binding the contract method 0x8ff68a80.
//
// Solidity: function getExpectedRateWithHintAndFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 rateWithNetworkFee, uint256 rateWithAllFees)
func (_KyberNetwork *KyberNetworkSession) GetExpectedRateWithHintAndFee(src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (struct {
	RateWithNetworkFee *big.Int
	RateWithAllFees    *big.Int
}, error) {
	return _KyberNetwork.Contract.GetExpectedRateWithHintAndFee(&_KyberNetwork.CallOpts, src, dest, srcQty, platformFeeBps, hint)
}

// GetExpectedRateWithHintAndFee is a free data retrieval call binding the contract method 0x8ff68a80.
//
// Solidity: function getExpectedRateWithHintAndFee(address src, address dest, uint256 srcQty, uint256 platformFeeBps, bytes hint) view returns(uint256 rateWithNetworkFee, uint256 rateWithAllFees)
func (_KyberNetwork *KyberNetworkCallerSession) GetExpectedRateWithHintAndFee(src common.Address, dest common.Address, srcQty *big.Int, platformFeeBps *big.Int, hint []byte) (struct {
	RateWithNetworkFee *big.Int
	RateWithAllFees    *big.Int
}, error) {
	return _KyberNetwork.Contract.GetExpectedRateWithHintAndFee(&_KyberNetwork.CallOpts, src, dest, srcQty, platformFeeBps, hint)
}

// GetNetworkData is a free data retrieval call binding the contract method 0x8881654e.
//
// Solidity: function getNetworkData() view returns(uint256 negligibleDiffBps, uint256 networkFeeBps, uint256 expiryTimestamp)
func (_KyberNetwork *KyberNetworkCaller) GetNetworkData(opts *bind.CallOpts) (struct {
	NegligibleDiffBps *big.Int
	NetworkFeeBps     *big.Int
	ExpiryTimestamp   *big.Int
}, error) {
	ret := new(struct {
		NegligibleDiffBps *big.Int
		NetworkFeeBps     *big.Int
		ExpiryTimestamp   *big.Int
	})
	out := ret
	err := _KyberNetwork.contract.Call(opts, out, "getNetworkData")
	return *ret, err
}

// GetNetworkData is a free data retrieval call binding the contract method 0x8881654e.
//
// Solidity: function getNetworkData() view returns(uint256 negligibleDiffBps, uint256 networkFeeBps, uint256 expiryTimestamp)
func (_KyberNetwork *KyberNetworkSession) GetNetworkData() (struct {
	NegligibleDiffBps *big.Int
	NetworkFeeBps     *big.Int
	ExpiryTimestamp   *big.Int
}, error) {
	return _KyberNetwork.Contract.GetNetworkData(&_KyberNetwork.CallOpts)
}

// GetNetworkData is a free data retrieval call binding the contract method 0x8881654e.
//
// Solidity: function getNetworkData() view returns(uint256 negligibleDiffBps, uint256 networkFeeBps, uint256 expiryTimestamp)
func (_KyberNetwork *KyberNetworkCallerSession) GetNetworkData() (struct {
	NegligibleDiffBps *big.Int
	NetworkFeeBps     *big.Int
	ExpiryTimestamp   *big.Int
}, error) {
	return _KyberNetwork.Contract.GetNetworkData(&_KyberNetwork.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetwork *KyberNetworkCaller) GetOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "getOperators")
	return *ret0, err
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetwork *KyberNetworkSession) GetOperators() ([]common.Address, error) {
	return _KyberNetwork.Contract.GetOperators(&_KyberNetwork.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0x27a099d8.
//
// Solidity: function getOperators() view returns(address[])
func (_KyberNetwork *KyberNetworkCallerSession) GetOperators() ([]common.Address, error) {
	return _KyberNetwork.Contract.GetOperators(&_KyberNetwork.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetwork *KyberNetworkCaller) MaxGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "maxGasPrice")
	return *ret0, err
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetwork *KyberNetworkSession) MaxGasPrice() (*big.Int, error) {
	return _KyberNetwork.Contract.MaxGasPrice(&_KyberNetwork.CallOpts)
}

// MaxGasPrice is a free data retrieval call binding the contract method 0x3de39c11.
//
// Solidity: function maxGasPrice() view returns(uint256)
func (_KyberNetwork *KyberNetworkCallerSession) MaxGasPrice() (*big.Int, error) {
	return _KyberNetwork.Contract.MaxGasPrice(&_KyberNetwork.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetwork *KyberNetworkCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KyberNetwork.contract.Call(opts, out, "pendingAdmin")
	return *ret0, err
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetwork *KyberNetworkSession) PendingAdmin() (common.Address, error) {
	return _KyberNetwork.Contract.PendingAdmin(&_KyberNetwork.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_KyberNetwork *KyberNetworkCallerSession) PendingAdmin() (common.Address, error) {
	return _KyberNetwork.Contract.PendingAdmin(&_KyberNetwork.CallOpts)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetwork *KyberNetworkTransactor) AddAlerter(opts *bind.TransactOpts, newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "addAlerter", newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetwork *KyberNetworkSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddAlerter(&_KyberNetwork.TransactOpts, newAlerter)
}

// AddAlerter is a paid mutator transaction binding the contract method 0x408ee7fe.
//
// Solidity: function addAlerter(address newAlerter) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) AddAlerter(newAlerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddAlerter(&_KyberNetwork.TransactOpts, newAlerter)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x52dd35b9.
//
// Solidity: function addKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkTransactor) AddKyberProxy(opts *bind.TransactOpts, kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "addKyberProxy", kyberProxy)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x52dd35b9.
//
// Solidity: function addKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkSession) AddKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddKyberProxy(&_KyberNetwork.TransactOpts, kyberProxy)
}

// AddKyberProxy is a paid mutator transaction binding the contract method 0x52dd35b9.
//
// Solidity: function addKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) AddKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddKyberProxy(&_KyberNetwork.TransactOpts, kyberProxy)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetwork *KyberNetworkTransactor) AddOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "addOperator", newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetwork *KyberNetworkSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddOperator(&_KyberNetwork.TransactOpts, newOperator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address newOperator) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) AddOperator(newOperator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.AddOperator(&_KyberNetwork.TransactOpts, newOperator)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetwork *KyberNetworkTransactor) ClaimAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "claimAdmin")
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetwork *KyberNetworkSession) ClaimAdmin() (*types.Transaction, error) {
	return _KyberNetwork.Contract.ClaimAdmin(&_KyberNetwork.TransactOpts)
}

// ClaimAdmin is a paid mutator transaction binding the contract method 0x77f50f97.
//
// Solidity: function claimAdmin() returns()
func (_KyberNetwork *KyberNetworkTransactorSession) ClaimAdmin() (*types.Transaction, error) {
	return _KyberNetwork.Contract.ClaimAdmin(&_KyberNetwork.TransactOpts)
}

// GetAndUpdateNetworkFee is a paid mutator transaction binding the contract method 0xdbe2dc9b.
//
// Solidity: function getAndUpdateNetworkFee() returns(uint256 networkFeeBps)
func (_KyberNetwork *KyberNetworkTransactor) GetAndUpdateNetworkFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "getAndUpdateNetworkFee")
}

// GetAndUpdateNetworkFee is a paid mutator transaction binding the contract method 0xdbe2dc9b.
//
// Solidity: function getAndUpdateNetworkFee() returns(uint256 networkFeeBps)
func (_KyberNetwork *KyberNetworkSession) GetAndUpdateNetworkFee() (*types.Transaction, error) {
	return _KyberNetwork.Contract.GetAndUpdateNetworkFee(&_KyberNetwork.TransactOpts)
}

// GetAndUpdateNetworkFee is a paid mutator transaction binding the contract method 0xdbe2dc9b.
//
// Solidity: function getAndUpdateNetworkFee() returns(uint256 networkFeeBps)
func (_KyberNetwork *KyberNetworkTransactorSession) GetAndUpdateNetworkFee() (*types.Transaction, error) {
	return _KyberNetwork.Contract.GetAndUpdateNetworkFee(&_KyberNetwork.TransactOpts)
}

// ListReservesForToken is a paid mutator transaction binding the contract method 0x7bb0ea82.
//
// Solidity: function listReservesForToken(address token, uint256 startIndex, uint256 endIndex, bool add) returns()
func (_KyberNetwork *KyberNetworkTransactor) ListReservesForToken(opts *bind.TransactOpts, token common.Address, startIndex *big.Int, endIndex *big.Int, add bool) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "listReservesForToken", token, startIndex, endIndex, add)
}

// ListReservesForToken is a paid mutator transaction binding the contract method 0x7bb0ea82.
//
// Solidity: function listReservesForToken(address token, uint256 startIndex, uint256 endIndex, bool add) returns()
func (_KyberNetwork *KyberNetworkSession) ListReservesForToken(token common.Address, startIndex *big.Int, endIndex *big.Int, add bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.ListReservesForToken(&_KyberNetwork.TransactOpts, token, startIndex, endIndex, add)
}

// ListReservesForToken is a paid mutator transaction binding the contract method 0x7bb0ea82.
//
// Solidity: function listReservesForToken(address token, uint256 startIndex, uint256 endIndex, bool add) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) ListReservesForToken(token common.Address, startIndex *big.Int, endIndex *big.Int, add bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.ListReservesForToken(&_KyberNetwork.TransactOpts, token, startIndex, endIndex, add)
}

// ListTokenForReserve is a paid mutator transaction binding the contract method 0x32c8bd9f.
//
// Solidity: function listTokenForReserve(address reserve, address token, bool add) returns()
func (_KyberNetwork *KyberNetworkTransactor) ListTokenForReserve(opts *bind.TransactOpts, reserve common.Address, token common.Address, add bool) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "listTokenForReserve", reserve, token, add)
}

// ListTokenForReserve is a paid mutator transaction binding the contract method 0x32c8bd9f.
//
// Solidity: function listTokenForReserve(address reserve, address token, bool add) returns()
func (_KyberNetwork *KyberNetworkSession) ListTokenForReserve(reserve common.Address, token common.Address, add bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.ListTokenForReserve(&_KyberNetwork.TransactOpts, reserve, token, add)
}

// ListTokenForReserve is a paid mutator transaction binding the contract method 0x32c8bd9f.
//
// Solidity: function listTokenForReserve(address reserve, address token, bool add) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) ListTokenForReserve(reserve common.Address, token common.Address, add bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.ListTokenForReserve(&_KyberNetwork.TransactOpts, reserve, token, add)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetwork *KyberNetworkTransactor) RemoveAlerter(opts *bind.TransactOpts, alerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "removeAlerter", alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetwork *KyberNetworkSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveAlerter(&_KyberNetwork.TransactOpts, alerter)
}

// RemoveAlerter is a paid mutator transaction binding the contract method 0x01a12fd3.
//
// Solidity: function removeAlerter(address alerter) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) RemoveAlerter(alerter common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveAlerter(&_KyberNetwork.TransactOpts, alerter)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkTransactor) RemoveKyberProxy(opts *bind.TransactOpts, kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "removeKyberProxy", kyberProxy)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkSession) RemoveKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveKyberProxy(&_KyberNetwork.TransactOpts, kyberProxy)
}

// RemoveKyberProxy is a paid mutator transaction binding the contract method 0x803d58c8.
//
// Solidity: function removeKyberProxy(address kyberProxy) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) RemoveKyberProxy(kyberProxy common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveKyberProxy(&_KyberNetwork.TransactOpts, kyberProxy)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetwork *KyberNetworkTransactor) RemoveOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "removeOperator", operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetwork *KyberNetworkSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveOperator(&_KyberNetwork.TransactOpts, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address operator) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) RemoveOperator(operator common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.RemoveOperator(&_KyberNetwork.TransactOpts, operator)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine, address _gasHelper) returns()
func (_KyberNetwork *KyberNetworkTransactor) SetContracts(opts *bind.TransactOpts, _kyberFeeHandler common.Address, _kyberMatchingEngine common.Address, _gasHelper common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "setContracts", _kyberFeeHandler, _kyberMatchingEngine, _gasHelper)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine, address _gasHelper) returns()
func (_KyberNetwork *KyberNetworkSession) SetContracts(_kyberFeeHandler common.Address, _kyberMatchingEngine common.Address, _gasHelper common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetContracts(&_KyberNetwork.TransactOpts, _kyberFeeHandler, _kyberMatchingEngine, _gasHelper)
}

// SetContracts is a paid mutator transaction binding the contract method 0xb3066d49.
//
// Solidity: function setContracts(address _kyberFeeHandler, address _kyberMatchingEngine, address _gasHelper) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) SetContracts(_kyberFeeHandler common.Address, _kyberMatchingEngine common.Address, _gasHelper common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetContracts(&_KyberNetwork.TransactOpts, _kyberFeeHandler, _kyberMatchingEngine, _gasHelper)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool enable) returns()
func (_KyberNetwork *KyberNetworkTransactor) SetEnable(opts *bind.TransactOpts, enable bool) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "setEnable", enable)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool enable) returns()
func (_KyberNetwork *KyberNetworkSession) SetEnable(enable bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetEnable(&_KyberNetwork.TransactOpts, enable)
}

// SetEnable is a paid mutator transaction binding the contract method 0x7726bed3.
//
// Solidity: function setEnable(bool enable) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) SetEnable(enable bool) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetEnable(&_KyberNetwork.TransactOpts, enable)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberNetwork *KyberNetworkTransactor) SetKyberDaoContract(opts *bind.TransactOpts, _kyberDao common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "setKyberDaoContract", _kyberDao)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberNetwork *KyberNetworkSession) SetKyberDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetKyberDaoContract(&_KyberNetwork.TransactOpts, _kyberDao)
}

// SetKyberDaoContract is a paid mutator transaction binding the contract method 0x6ff277de.
//
// Solidity: function setKyberDaoContract(address _kyberDao) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) SetKyberDaoContract(_kyberDao common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetKyberDaoContract(&_KyberNetwork.TransactOpts, _kyberDao)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiffBps) returns()
func (_KyberNetwork *KyberNetworkTransactor) SetParams(opts *bind.TransactOpts, _maxGasPrice *big.Int, _negligibleRateDiffBps *big.Int) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "setParams", _maxGasPrice, _negligibleRateDiffBps)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiffBps) returns()
func (_KyberNetwork *KyberNetworkSession) SetParams(_maxGasPrice *big.Int, _negligibleRateDiffBps *big.Int) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetParams(&_KyberNetwork.TransactOpts, _maxGasPrice, _negligibleRateDiffBps)
}

// SetParams is a paid mutator transaction binding the contract method 0xc0324c77.
//
// Solidity: function setParams(uint256 _maxGasPrice, uint256 _negligibleRateDiffBps) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) SetParams(_maxGasPrice *big.Int, _negligibleRateDiffBps *big.Int) (*types.Transaction, error) {
	return _KyberNetwork.Contract.SetParams(&_KyberNetwork.TransactOpts, _maxGasPrice, _negligibleRateDiffBps)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkTransactor) TradeWithHint(opts *bind.TransactOpts, trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "tradeWithHint", trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkSession) TradeWithHint(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TradeWithHint(&_KyberNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHint is a paid mutator transaction binding the contract method 0x088322ef.
//
// Solidity: function tradeWithHint(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address walletId, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkTransactorSession) TradeWithHint(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, walletId common.Address, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TradeWithHint(&_KyberNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, walletId, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xc43190f5.
//
// Solidity: function tradeWithHintAndFee(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkTransactor) TradeWithHintAndFee(opts *bind.TransactOpts, trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "tradeWithHintAndFee", trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xc43190f5.
//
// Solidity: function tradeWithHintAndFee(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkSession) TradeWithHintAndFee(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TradeWithHintAndFee(&_KyberNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TradeWithHintAndFee is a paid mutator transaction binding the contract method 0xc43190f5.
//
// Solidity: function tradeWithHintAndFee(address trader, address src, uint256 srcAmount, address dest, address destAddress, uint256 maxDestAmount, uint256 minConversionRate, address platformWallet, uint256 platformFeeBps, bytes hint) payable returns(uint256 destAmount)
func (_KyberNetwork *KyberNetworkTransactorSession) TradeWithHintAndFee(trader common.Address, src common.Address, srcAmount *big.Int, dest common.Address, destAddress common.Address, maxDestAmount *big.Int, minConversionRate *big.Int, platformWallet common.Address, platformFeeBps *big.Int, hint []byte) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TradeWithHintAndFee(&_KyberNetwork.TransactOpts, trader, src, srcAmount, dest, destAddress, maxDestAmount, minConversionRate, platformWallet, platformFeeBps, hint)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TransferAdmin(&_KyberNetwork.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TransferAdmin(&_KyberNetwork.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkTransactor) TransferAdminQuickly(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "transferAdminQuickly", newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TransferAdminQuickly(&_KyberNetwork.TransactOpts, newAdmin)
}

// TransferAdminQuickly is a paid mutator transaction binding the contract method 0x7acc8678.
//
// Solidity: function transferAdminQuickly(address newAdmin) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) TransferAdminQuickly(newAdmin common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.TransferAdminQuickly(&_KyberNetwork.TransactOpts, newAdmin)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkTransactor) WithdrawEther(opts *bind.TransactOpts, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "withdrawEther", amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.WithdrawEther(&_KyberNetwork.TransactOpts, amount, sendTo)
}

// WithdrawEther is a paid mutator transaction binding the contract method 0xce56c454.
//
// Solidity: function withdrawEther(uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) WithdrawEther(amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.WithdrawEther(&_KyberNetwork.TransactOpts, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.contract.Transact(opts, "withdrawToken", token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.WithdrawToken(&_KyberNetwork.TransactOpts, token, amount, sendTo)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address token, uint256 amount, address sendTo) returns()
func (_KyberNetwork *KyberNetworkTransactorSession) WithdrawToken(token common.Address, amount *big.Int, sendTo common.Address) (*types.Transaction, error) {
	return _KyberNetwork.Contract.WithdrawToken(&_KyberNetwork.TransactOpts, token, amount, sendTo)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberNetwork *KyberNetworkTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KyberNetwork.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberNetwork *KyberNetworkSession) Receive() (*types.Transaction, error) {
	return _KyberNetwork.Contract.Receive(&_KyberNetwork.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KyberNetwork *KyberNetworkTransactorSession) Receive() (*types.Transaction, error) {
	return _KyberNetwork.Contract.Receive(&_KyberNetwork.TransactOpts)
}

// KyberNetworkAdminClaimedIterator is returned from FilterAdminClaimed and is used to iterate over the raw logs and unpacked data for AdminClaimed events raised by the KyberNetwork contract.
type KyberNetworkAdminClaimedIterator struct {
	Event *KyberNetworkAdminClaimed // Event containing the contract specifics and raw log

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
func (it *KyberNetworkAdminClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkAdminClaimed)
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
		it.Event = new(KyberNetworkAdminClaimed)
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
func (it *KyberNetworkAdminClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkAdminClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkAdminClaimed represents a AdminClaimed event raised by the KyberNetwork contract.
type KyberNetworkAdminClaimed struct {
	NewAdmin      common.Address
	PreviousAdmin common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminClaimed is a free log retrieval operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberNetwork *KyberNetworkFilterer) FilterAdminClaimed(opts *bind.FilterOpts) (*KyberNetworkAdminClaimedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkAdminClaimedIterator{contract: _KyberNetwork.contract, event: "AdminClaimed", logs: logs, sub: sub}, nil
}

// WatchAdminClaimed is a free log subscription operation binding the contract event 0x65da1cfc2c2e81576ad96afb24a581f8e109b7a403b35cbd3243a1c99efdb9ed.
//
// Solidity: event AdminClaimed(address newAdmin, address previousAdmin)
func (_KyberNetwork *KyberNetworkFilterer) WatchAdminClaimed(opts *bind.WatchOpts, sink chan<- *KyberNetworkAdminClaimed) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "AdminClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkAdminClaimed)
				if err := _KyberNetwork.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseAdminClaimed(log types.Log) (*KyberNetworkAdminClaimed, error) {
	event := new(KyberNetworkAdminClaimed)
	if err := _KyberNetwork.contract.UnpackLog(event, "AdminClaimed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkAlerterAddedIterator is returned from FilterAlerterAdded and is used to iterate over the raw logs and unpacked data for AlerterAdded events raised by the KyberNetwork contract.
type KyberNetworkAlerterAddedIterator struct {
	Event *KyberNetworkAlerterAdded // Event containing the contract specifics and raw log

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
func (it *KyberNetworkAlerterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkAlerterAdded)
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
		it.Event = new(KyberNetworkAlerterAdded)
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
func (it *KyberNetworkAlerterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkAlerterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkAlerterAdded represents a AlerterAdded event raised by the KyberNetwork contract.
type KyberNetworkAlerterAdded struct {
	NewAlerter common.Address
	IsAdd      bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlerterAdded is a free log retrieval operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberNetwork *KyberNetworkFilterer) FilterAlerterAdded(opts *bind.FilterOpts) (*KyberNetworkAlerterAddedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkAlerterAddedIterator{contract: _KyberNetwork.contract, event: "AlerterAdded", logs: logs, sub: sub}, nil
}

// WatchAlerterAdded is a free log subscription operation binding the contract event 0x5611bf3e417d124f97bf2c788843ea8bb502b66079fbee02158ef30b172cb762.
//
// Solidity: event AlerterAdded(address newAlerter, bool isAdd)
func (_KyberNetwork *KyberNetworkFilterer) WatchAlerterAdded(opts *bind.WatchOpts, sink chan<- *KyberNetworkAlerterAdded) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "AlerterAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkAlerterAdded)
				if err := _KyberNetwork.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseAlerterAdded(log types.Log) (*KyberNetworkAlerterAdded, error) {
	event := new(KyberNetworkAlerterAdded)
	if err := _KyberNetwork.contract.UnpackLog(event, "AlerterAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkEtherReceivalIterator is returned from FilterEtherReceival and is used to iterate over the raw logs and unpacked data for EtherReceival events raised by the KyberNetwork contract.
type KyberNetworkEtherReceivalIterator struct {
	Event *KyberNetworkEtherReceival // Event containing the contract specifics and raw log

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
func (it *KyberNetworkEtherReceivalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkEtherReceival)
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
		it.Event = new(KyberNetworkEtherReceival)
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
func (it *KyberNetworkEtherReceivalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkEtherReceivalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkEtherReceival represents a EtherReceival event raised by the KyberNetwork contract.
type KyberNetworkEtherReceival struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherReceival is a free log retrieval operation binding the contract event 0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619.
//
// Solidity: event EtherReceival(address indexed sender, uint256 amount)
func (_KyberNetwork *KyberNetworkFilterer) FilterEtherReceival(opts *bind.FilterOpts, sender []common.Address) (*KyberNetworkEtherReceivalIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "EtherReceival", senderRule)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkEtherReceivalIterator{contract: _KyberNetwork.contract, event: "EtherReceival", logs: logs, sub: sub}, nil
}

// WatchEtherReceival is a free log subscription operation binding the contract event 0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619.
//
// Solidity: event EtherReceival(address indexed sender, uint256 amount)
func (_KyberNetwork *KyberNetworkFilterer) WatchEtherReceival(opts *bind.WatchOpts, sink chan<- *KyberNetworkEtherReceival, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "EtherReceival", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkEtherReceival)
				if err := _KyberNetwork.contract.UnpackLog(event, "EtherReceival", log); err != nil {
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

// ParseEtherReceival is a log parse operation binding the contract event 0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619.
//
// Solidity: event EtherReceival(address indexed sender, uint256 amount)
func (_KyberNetwork *KyberNetworkFilterer) ParseEtherReceival(log types.Log) (*KyberNetworkEtherReceival, error) {
	event := new(KyberNetworkEtherReceival)
	if err := _KyberNetwork.contract.UnpackLog(event, "EtherReceival", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkEtherWithdrawIterator is returned from FilterEtherWithdraw and is used to iterate over the raw logs and unpacked data for EtherWithdraw events raised by the KyberNetwork contract.
type KyberNetworkEtherWithdrawIterator struct {
	Event *KyberNetworkEtherWithdraw // Event containing the contract specifics and raw log

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
func (it *KyberNetworkEtherWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkEtherWithdraw)
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
		it.Event = new(KyberNetworkEtherWithdraw)
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
func (it *KyberNetworkEtherWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkEtherWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkEtherWithdraw represents a EtherWithdraw event raised by the KyberNetwork contract.
type KyberNetworkEtherWithdraw struct {
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdraw is a free log retrieval operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_KyberNetwork *KyberNetworkFilterer) FilterEtherWithdraw(opts *bind.FilterOpts) (*KyberNetworkEtherWithdrawIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkEtherWithdrawIterator{contract: _KyberNetwork.contract, event: "EtherWithdraw", logs: logs, sub: sub}, nil
}

// WatchEtherWithdraw is a free log subscription operation binding the contract event 0xec47e7ed86c86774d1a72c19f35c639911393fe7c1a34031fdbd260890da90de.
//
// Solidity: event EtherWithdraw(uint256 amount, address sendTo)
func (_KyberNetwork *KyberNetworkFilterer) WatchEtherWithdraw(opts *bind.WatchOpts, sink chan<- *KyberNetworkEtherWithdraw) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "EtherWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkEtherWithdraw)
				if err := _KyberNetwork.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseEtherWithdraw(log types.Log) (*KyberNetworkEtherWithdraw, error) {
	event := new(KyberNetworkEtherWithdraw)
	if err := _KyberNetwork.contract.UnpackLog(event, "EtherWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkGasHelperUpdatedIterator is returned from FilterGasHelperUpdated and is used to iterate over the raw logs and unpacked data for GasHelperUpdated events raised by the KyberNetwork contract.
type KyberNetworkGasHelperUpdatedIterator struct {
	Event *KyberNetworkGasHelperUpdated // Event containing the contract specifics and raw log

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
func (it *KyberNetworkGasHelperUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkGasHelperUpdated)
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
		it.Event = new(KyberNetworkGasHelperUpdated)
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
func (it *KyberNetworkGasHelperUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkGasHelperUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkGasHelperUpdated represents a GasHelperUpdated event raised by the KyberNetwork contract.
type KyberNetworkGasHelperUpdated struct {
	NewGasHelper common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterGasHelperUpdated is a free log retrieval operation binding the contract event 0x95ba6becebde78de944071b522d81414292f67e3d95db7d9df46bb8e8b3da8b8.
//
// Solidity: event GasHelperUpdated(address newGasHelper)
func (_KyberNetwork *KyberNetworkFilterer) FilterGasHelperUpdated(opts *bind.FilterOpts) (*KyberNetworkGasHelperUpdatedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "GasHelperUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkGasHelperUpdatedIterator{contract: _KyberNetwork.contract, event: "GasHelperUpdated", logs: logs, sub: sub}, nil
}

// WatchGasHelperUpdated is a free log subscription operation binding the contract event 0x95ba6becebde78de944071b522d81414292f67e3d95db7d9df46bb8e8b3da8b8.
//
// Solidity: event GasHelperUpdated(address newGasHelper)
func (_KyberNetwork *KyberNetworkFilterer) WatchGasHelperUpdated(opts *bind.WatchOpts, sink chan<- *KyberNetworkGasHelperUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "GasHelperUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkGasHelperUpdated)
				if err := _KyberNetwork.contract.UnpackLog(event, "GasHelperUpdated", log); err != nil {
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

// ParseGasHelperUpdated is a log parse operation binding the contract event 0x95ba6becebde78de944071b522d81414292f67e3d95db7d9df46bb8e8b3da8b8.
//
// Solidity: event GasHelperUpdated(address newGasHelper)
func (_KyberNetwork *KyberNetworkFilterer) ParseGasHelperUpdated(log types.Log) (*KyberNetworkGasHelperUpdated, error) {
	event := new(KyberNetworkGasHelperUpdated)
	if err := _KyberNetwork.contract.UnpackLog(event, "GasHelperUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberDaoUpdatedIterator is returned from FilterKyberDaoUpdated and is used to iterate over the raw logs and unpacked data for KyberDaoUpdated events raised by the KyberNetwork contract.
type KyberNetworkKyberDaoUpdatedIterator struct {
	Event *KyberNetworkKyberDaoUpdated // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberDaoUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberDaoUpdated)
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
		it.Event = new(KyberNetworkKyberDaoUpdated)
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
func (it *KyberNetworkKyberDaoUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberDaoUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberDaoUpdated represents a KyberDaoUpdated event raised by the KyberNetwork contract.
type KyberNetworkKyberDaoUpdated struct {
	NewKyberDao common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterKyberDaoUpdated is a free log retrieval operation binding the contract event 0x16a2e1af8449067f38aa765b54d479785c94d8ebdfbba7b410e3488b0877c1e4.
//
// Solidity: event KyberDaoUpdated(address newKyberDao)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberDaoUpdated(opts *bind.FilterOpts) (*KyberNetworkKyberDaoUpdatedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberDaoUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberDaoUpdatedIterator{contract: _KyberNetwork.contract, event: "KyberDaoUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberDaoUpdated is a free log subscription operation binding the contract event 0x16a2e1af8449067f38aa765b54d479785c94d8ebdfbba7b410e3488b0877c1e4.
//
// Solidity: event KyberDaoUpdated(address newKyberDao)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberDaoUpdated(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberDaoUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberDaoUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberDaoUpdated)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberDaoUpdated", log); err != nil {
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

// ParseKyberDaoUpdated is a log parse operation binding the contract event 0x16a2e1af8449067f38aa765b54d479785c94d8ebdfbba7b410e3488b0877c1e4.
//
// Solidity: event KyberDaoUpdated(address newKyberDao)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberDaoUpdated(log types.Log) (*KyberNetworkKyberDaoUpdated, error) {
	event := new(KyberNetworkKyberDaoUpdated)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberDaoUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberFeeHandlerUpdatedIterator is returned from FilterKyberFeeHandlerUpdated and is used to iterate over the raw logs and unpacked data for KyberFeeHandlerUpdated events raised by the KyberNetwork contract.
type KyberNetworkKyberFeeHandlerUpdatedIterator struct {
	Event *KyberNetworkKyberFeeHandlerUpdated // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberFeeHandlerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberFeeHandlerUpdated)
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
		it.Event = new(KyberNetworkKyberFeeHandlerUpdated)
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
func (it *KyberNetworkKyberFeeHandlerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberFeeHandlerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberFeeHandlerUpdated represents a KyberFeeHandlerUpdated event raised by the KyberNetwork contract.
type KyberNetworkKyberFeeHandlerUpdated struct {
	NewKyberFeeHandler common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterKyberFeeHandlerUpdated is a free log retrieval operation binding the contract event 0x5128fc9be01065f3cabe4c8b72796eb6b8a00284f39a2390cd71e91b509f90b6.
//
// Solidity: event KyberFeeHandlerUpdated(address newKyberFeeHandler)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberFeeHandlerUpdated(opts *bind.FilterOpts) (*KyberNetworkKyberFeeHandlerUpdatedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberFeeHandlerUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberFeeHandlerUpdatedIterator{contract: _KyberNetwork.contract, event: "KyberFeeHandlerUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberFeeHandlerUpdated is a free log subscription operation binding the contract event 0x5128fc9be01065f3cabe4c8b72796eb6b8a00284f39a2390cd71e91b509f90b6.
//
// Solidity: event KyberFeeHandlerUpdated(address newKyberFeeHandler)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberFeeHandlerUpdated(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberFeeHandlerUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberFeeHandlerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberFeeHandlerUpdated)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberFeeHandlerUpdated", log); err != nil {
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

// ParseKyberFeeHandlerUpdated is a log parse operation binding the contract event 0x5128fc9be01065f3cabe4c8b72796eb6b8a00284f39a2390cd71e91b509f90b6.
//
// Solidity: event KyberFeeHandlerUpdated(address newKyberFeeHandler)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberFeeHandlerUpdated(log types.Log) (*KyberNetworkKyberFeeHandlerUpdated, error) {
	event := new(KyberNetworkKyberFeeHandlerUpdated)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberFeeHandlerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberMatchingEngineUpdatedIterator is returned from FilterKyberMatchingEngineUpdated and is used to iterate over the raw logs and unpacked data for KyberMatchingEngineUpdated events raised by the KyberNetwork contract.
type KyberNetworkKyberMatchingEngineUpdatedIterator struct {
	Event *KyberNetworkKyberMatchingEngineUpdated // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberMatchingEngineUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberMatchingEngineUpdated)
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
		it.Event = new(KyberNetworkKyberMatchingEngineUpdated)
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
func (it *KyberNetworkKyberMatchingEngineUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberMatchingEngineUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberMatchingEngineUpdated represents a KyberMatchingEngineUpdated event raised by the KyberNetwork contract.
type KyberNetworkKyberMatchingEngineUpdated struct {
	NewKyberMatchingEngine common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterKyberMatchingEngineUpdated is a free log retrieval operation binding the contract event 0x92b5317eb7846d2e62df8ff23b97a564e012e63defef32777f017249f76bf264.
//
// Solidity: event KyberMatchingEngineUpdated(address newKyberMatchingEngine)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberMatchingEngineUpdated(opts *bind.FilterOpts) (*KyberNetworkKyberMatchingEngineUpdatedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberMatchingEngineUpdated")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberMatchingEngineUpdatedIterator{contract: _KyberNetwork.contract, event: "KyberMatchingEngineUpdated", logs: logs, sub: sub}, nil
}

// WatchKyberMatchingEngineUpdated is a free log subscription operation binding the contract event 0x92b5317eb7846d2e62df8ff23b97a564e012e63defef32777f017249f76bf264.
//
// Solidity: event KyberMatchingEngineUpdated(address newKyberMatchingEngine)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberMatchingEngineUpdated(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberMatchingEngineUpdated) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberMatchingEngineUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberMatchingEngineUpdated)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberMatchingEngineUpdated", log); err != nil {
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

// ParseKyberMatchingEngineUpdated is a log parse operation binding the contract event 0x92b5317eb7846d2e62df8ff23b97a564e012e63defef32777f017249f76bf264.
//
// Solidity: event KyberMatchingEngineUpdated(address newKyberMatchingEngine)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberMatchingEngineUpdated(log types.Log) (*KyberNetworkKyberMatchingEngineUpdated, error) {
	event := new(KyberNetworkKyberMatchingEngineUpdated)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberMatchingEngineUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberNetworkParamsSetIterator is returned from FilterKyberNetworkParamsSet and is used to iterate over the raw logs and unpacked data for KyberNetworkParamsSet events raised by the KyberNetwork contract.
type KyberNetworkKyberNetworkParamsSetIterator struct {
	Event *KyberNetworkKyberNetworkParamsSet // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberNetworkParamsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberNetworkParamsSet)
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
		it.Event = new(KyberNetworkKyberNetworkParamsSet)
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
func (it *KyberNetworkKyberNetworkParamsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberNetworkParamsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberNetworkParamsSet represents a KyberNetworkParamsSet event raised by the KyberNetwork contract.
type KyberNetworkKyberNetworkParamsSet struct {
	MaxGasPrice           *big.Int
	NegligibleRateDiffBps *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkParamsSet is a free log retrieval operation binding the contract event 0xc1e6729d7fd9a615adc03ebe7d8ff15649d8eed7516bf6c30538a1e722bb1975.
//
// Solidity: event KyberNetworkParamsSet(uint256 maxGasPrice, uint256 negligibleRateDiffBps)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberNetworkParamsSet(opts *bind.FilterOpts) (*KyberNetworkKyberNetworkParamsSetIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberNetworkParamsSet")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberNetworkParamsSetIterator{contract: _KyberNetwork.contract, event: "KyberNetworkParamsSet", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkParamsSet is a free log subscription operation binding the contract event 0xc1e6729d7fd9a615adc03ebe7d8ff15649d8eed7516bf6c30538a1e722bb1975.
//
// Solidity: event KyberNetworkParamsSet(uint256 maxGasPrice, uint256 negligibleRateDiffBps)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberNetworkParamsSet(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberNetworkParamsSet) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberNetworkParamsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberNetworkParamsSet)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberNetworkParamsSet", log); err != nil {
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

// ParseKyberNetworkParamsSet is a log parse operation binding the contract event 0xc1e6729d7fd9a615adc03ebe7d8ff15649d8eed7516bf6c30538a1e722bb1975.
//
// Solidity: event KyberNetworkParamsSet(uint256 maxGasPrice, uint256 negligibleRateDiffBps)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberNetworkParamsSet(log types.Log) (*KyberNetworkKyberNetworkParamsSet, error) {
	event := new(KyberNetworkKyberNetworkParamsSet)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberNetworkParamsSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberNetworkSetEnableIterator is returned from FilterKyberNetworkSetEnable and is used to iterate over the raw logs and unpacked data for KyberNetworkSetEnable events raised by the KyberNetwork contract.
type KyberNetworkKyberNetworkSetEnableIterator struct {
	Event *KyberNetworkKyberNetworkSetEnable // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberNetworkSetEnableIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberNetworkSetEnable)
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
		it.Event = new(KyberNetworkKyberNetworkSetEnable)
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
func (it *KyberNetworkKyberNetworkSetEnableIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberNetworkSetEnableIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberNetworkSetEnable represents a KyberNetworkSetEnable event raised by the KyberNetwork contract.
type KyberNetworkKyberNetworkSetEnable struct {
	IsEnabled bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterKyberNetworkSetEnable is a free log retrieval operation binding the contract event 0x8a846a525e22497042ee2f99423a8ff8bbb831d3ae5384692bf6040f591c1eba.
//
// Solidity: event KyberNetworkSetEnable(bool isEnabled)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberNetworkSetEnable(opts *bind.FilterOpts) (*KyberNetworkKyberNetworkSetEnableIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberNetworkSetEnable")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberNetworkSetEnableIterator{contract: _KyberNetwork.contract, event: "KyberNetworkSetEnable", logs: logs, sub: sub}, nil
}

// WatchKyberNetworkSetEnable is a free log subscription operation binding the contract event 0x8a846a525e22497042ee2f99423a8ff8bbb831d3ae5384692bf6040f591c1eba.
//
// Solidity: event KyberNetworkSetEnable(bool isEnabled)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberNetworkSetEnable(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberNetworkSetEnable) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberNetworkSetEnable")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberNetworkSetEnable)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberNetworkSetEnable", log); err != nil {
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

// ParseKyberNetworkSetEnable is a log parse operation binding the contract event 0x8a846a525e22497042ee2f99423a8ff8bbb831d3ae5384692bf6040f591c1eba.
//
// Solidity: event KyberNetworkSetEnable(bool isEnabled)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberNetworkSetEnable(log types.Log) (*KyberNetworkKyberNetworkSetEnable, error) {
	event := new(KyberNetworkKyberNetworkSetEnable)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberNetworkSetEnable", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberProxyAddedIterator is returned from FilterKyberProxyAdded and is used to iterate over the raw logs and unpacked data for KyberProxyAdded events raised by the KyberNetwork contract.
type KyberNetworkKyberProxyAddedIterator struct {
	Event *KyberNetworkKyberProxyAdded // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberProxyAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberProxyAdded)
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
		it.Event = new(KyberNetworkKyberProxyAdded)
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
func (it *KyberNetworkKyberProxyAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberProxyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberProxyAdded represents a KyberProxyAdded event raised by the KyberNetwork contract.
type KyberNetworkKyberProxyAdded struct {
	KyberProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterKyberProxyAdded is a free log retrieval operation binding the contract event 0x0b008ff10c7e378a96d6566635e1aa748886d16fb87659faee2aa20608fec815.
//
// Solidity: event KyberProxyAdded(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberProxyAdded(opts *bind.FilterOpts) (*KyberNetworkKyberProxyAddedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberProxyAdded")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberProxyAddedIterator{contract: _KyberNetwork.contract, event: "KyberProxyAdded", logs: logs, sub: sub}, nil
}

// WatchKyberProxyAdded is a free log subscription operation binding the contract event 0x0b008ff10c7e378a96d6566635e1aa748886d16fb87659faee2aa20608fec815.
//
// Solidity: event KyberProxyAdded(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberProxyAdded(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberProxyAdded) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberProxyAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberProxyAdded)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberProxyAdded", log); err != nil {
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

// ParseKyberProxyAdded is a log parse operation binding the contract event 0x0b008ff10c7e378a96d6566635e1aa748886d16fb87659faee2aa20608fec815.
//
// Solidity: event KyberProxyAdded(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberProxyAdded(log types.Log) (*KyberNetworkKyberProxyAdded, error) {
	event := new(KyberNetworkKyberProxyAdded)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberProxyAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberProxyRemovedIterator is returned from FilterKyberProxyRemoved and is used to iterate over the raw logs and unpacked data for KyberProxyRemoved events raised by the KyberNetwork contract.
type KyberNetworkKyberProxyRemovedIterator struct {
	Event *KyberNetworkKyberProxyRemoved // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberProxyRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberProxyRemoved)
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
		it.Event = new(KyberNetworkKyberProxyRemoved)
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
func (it *KyberNetworkKyberProxyRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberProxyRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberProxyRemoved represents a KyberProxyRemoved event raised by the KyberNetwork contract.
type KyberNetworkKyberProxyRemoved struct {
	KyberProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterKyberProxyRemoved is a free log retrieval operation binding the contract event 0xbb9ee888852ae070b75270fa50ea2845ba32102d3a96842c7c416d12aad2f487.
//
// Solidity: event KyberProxyRemoved(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberProxyRemoved(opts *bind.FilterOpts) (*KyberNetworkKyberProxyRemovedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberProxyRemoved")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberProxyRemovedIterator{contract: _KyberNetwork.contract, event: "KyberProxyRemoved", logs: logs, sub: sub}, nil
}

// WatchKyberProxyRemoved is a free log subscription operation binding the contract event 0xbb9ee888852ae070b75270fa50ea2845ba32102d3a96842c7c416d12aad2f487.
//
// Solidity: event KyberProxyRemoved(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberProxyRemoved(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberProxyRemoved) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberProxyRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberProxyRemoved)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberProxyRemoved", log); err != nil {
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

// ParseKyberProxyRemoved is a log parse operation binding the contract event 0xbb9ee888852ae070b75270fa50ea2845ba32102d3a96842c7c416d12aad2f487.
//
// Solidity: event KyberProxyRemoved(address kyberProxy)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberProxyRemoved(log types.Log) (*KyberNetworkKyberProxyRemoved, error) {
	event := new(KyberNetworkKyberProxyRemoved)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberProxyRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkKyberTradeIterator is returned from FilterKyberTrade and is used to iterate over the raw logs and unpacked data for KyberTrade events raised by the KyberNetwork contract.
type KyberNetworkKyberTradeIterator struct {
	Event *KyberNetworkKyberTrade // Event containing the contract specifics and raw log

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
func (it *KyberNetworkKyberTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkKyberTrade)
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
		it.Event = new(KyberNetworkKyberTrade)
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
func (it *KyberNetworkKyberTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkKyberTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkKyberTrade represents a KyberTrade event raised by the KyberNetwork contract.
type KyberNetworkKyberTrade struct {
	Src                  common.Address
	Dest                 common.Address
	EthWeiValue          *big.Int
	NetworkFeeWei        *big.Int
	CustomPlatformFeeWei *big.Int
	T2eIds               [][32]byte
	E2tIds               [][32]byte
	T2eSrcAmounts        []*big.Int
	E2tSrcAmounts        []*big.Int
	T2eRates             []*big.Int
	E2tRates             []*big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterKyberTrade is a free log retrieval operation binding the contract event 0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d.
//
// Solidity: event KyberTrade(address indexed src, address indexed dest, uint256 ethWeiValue, uint256 networkFeeWei, uint256 customPlatformFeeWei, bytes32[] t2eIds, bytes32[] e2tIds, uint256[] t2eSrcAmounts, uint256[] e2tSrcAmounts, uint256[] t2eRates, uint256[] e2tRates)
func (_KyberNetwork *KyberNetworkFilterer) FilterKyberTrade(opts *bind.FilterOpts, src []common.Address, dest []common.Address) (*KyberNetworkKyberTradeIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "KyberTrade", srcRule, destRule)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkKyberTradeIterator{contract: _KyberNetwork.contract, event: "KyberTrade", logs: logs, sub: sub}, nil
}

// WatchKyberTrade is a free log subscription operation binding the contract event 0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d.
//
// Solidity: event KyberTrade(address indexed src, address indexed dest, uint256 ethWeiValue, uint256 networkFeeWei, uint256 customPlatformFeeWei, bytes32[] t2eIds, bytes32[] e2tIds, uint256[] t2eSrcAmounts, uint256[] e2tSrcAmounts, uint256[] t2eRates, uint256[] e2tRates)
func (_KyberNetwork *KyberNetworkFilterer) WatchKyberTrade(opts *bind.WatchOpts, sink chan<- *KyberNetworkKyberTrade, src []common.Address, dest []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "KyberTrade", srcRule, destRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkKyberTrade)
				if err := _KyberNetwork.contract.UnpackLog(event, "KyberTrade", log); err != nil {
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

// ParseKyberTrade is a log parse operation binding the contract event 0x30bbea603a7b36858fe5e3ec6ba5ff59dde039d02120d758eacfaed01520577d.
//
// Solidity: event KyberTrade(address indexed src, address indexed dest, uint256 ethWeiValue, uint256 networkFeeWei, uint256 customPlatformFeeWei, bytes32[] t2eIds, bytes32[] e2tIds, uint256[] t2eSrcAmounts, uint256[] e2tSrcAmounts, uint256[] t2eRates, uint256[] e2tRates)
func (_KyberNetwork *KyberNetworkFilterer) ParseKyberTrade(log types.Log) (*KyberNetworkKyberTrade, error) {
	event := new(KyberNetworkKyberTrade)
	if err := _KyberNetwork.contract.UnpackLog(event, "KyberTrade", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkListedReservesForTokenIterator is returned from FilterListedReservesForToken and is used to iterate over the raw logs and unpacked data for ListedReservesForToken events raised by the KyberNetwork contract.
type KyberNetworkListedReservesForTokenIterator struct {
	Event *KyberNetworkListedReservesForToken // Event containing the contract specifics and raw log

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
func (it *KyberNetworkListedReservesForTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkListedReservesForToken)
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
		it.Event = new(KyberNetworkListedReservesForToken)
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
func (it *KyberNetworkListedReservesForTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkListedReservesForTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkListedReservesForToken represents a ListedReservesForToken event raised by the KyberNetwork contract.
type KyberNetworkListedReservesForToken struct {
	Token    common.Address
	Reserves []common.Address
	Add      bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterListedReservesForToken is a free log retrieval operation binding the contract event 0xd4b0877e3beef91cd767680ac04114217ec7c9cb3a4705c03fc8061de81168fc.
//
// Solidity: event ListedReservesForToken(address indexed token, address[] reserves, bool add)
func (_KyberNetwork *KyberNetworkFilterer) FilterListedReservesForToken(opts *bind.FilterOpts, token []common.Address) (*KyberNetworkListedReservesForTokenIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "ListedReservesForToken", tokenRule)
	if err != nil {
		return nil, err
	}
	return &KyberNetworkListedReservesForTokenIterator{contract: _KyberNetwork.contract, event: "ListedReservesForToken", logs: logs, sub: sub}, nil
}

// WatchListedReservesForToken is a free log subscription operation binding the contract event 0xd4b0877e3beef91cd767680ac04114217ec7c9cb3a4705c03fc8061de81168fc.
//
// Solidity: event ListedReservesForToken(address indexed token, address[] reserves, bool add)
func (_KyberNetwork *KyberNetworkFilterer) WatchListedReservesForToken(opts *bind.WatchOpts, sink chan<- *KyberNetworkListedReservesForToken, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "ListedReservesForToken", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkListedReservesForToken)
				if err := _KyberNetwork.contract.UnpackLog(event, "ListedReservesForToken", log); err != nil {
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

// ParseListedReservesForToken is a log parse operation binding the contract event 0xd4b0877e3beef91cd767680ac04114217ec7c9cb3a4705c03fc8061de81168fc.
//
// Solidity: event ListedReservesForToken(address indexed token, address[] reserves, bool add)
func (_KyberNetwork *KyberNetworkFilterer) ParseListedReservesForToken(log types.Log) (*KyberNetworkListedReservesForToken, error) {
	event := new(KyberNetworkListedReservesForToken)
	if err := _KyberNetwork.contract.UnpackLog(event, "ListedReservesForToken", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the KyberNetwork contract.
type KyberNetworkOperatorAddedIterator struct {
	Event *KyberNetworkOperatorAdded // Event containing the contract specifics and raw log

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
func (it *KyberNetworkOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkOperatorAdded)
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
		it.Event = new(KyberNetworkOperatorAdded)
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
func (it *KyberNetworkOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkOperatorAdded represents a OperatorAdded event raised by the KyberNetwork contract.
type KyberNetworkOperatorAdded struct {
	NewOperator common.Address
	IsAdd       bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberNetwork *KyberNetworkFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*KyberNetworkOperatorAddedIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkOperatorAddedIterator{contract: _KyberNetwork.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x091a7a4b85135fdd7e8dbc18b12fabe5cc191ea867aa3c2e1a24a102af61d58b.
//
// Solidity: event OperatorAdded(address newOperator, bool isAdd)
func (_KyberNetwork *KyberNetworkFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *KyberNetworkOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkOperatorAdded)
				if err := _KyberNetwork.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseOperatorAdded(log types.Log) (*KyberNetworkOperatorAdded, error) {
	event := new(KyberNetworkOperatorAdded)
	if err := _KyberNetwork.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkTokenWithdrawIterator is returned from FilterTokenWithdraw and is used to iterate over the raw logs and unpacked data for TokenWithdraw events raised by the KyberNetwork contract.
type KyberNetworkTokenWithdrawIterator struct {
	Event *KyberNetworkTokenWithdraw // Event containing the contract specifics and raw log

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
func (it *KyberNetworkTokenWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkTokenWithdraw)
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
		it.Event = new(KyberNetworkTokenWithdraw)
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
func (it *KyberNetworkTokenWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkTokenWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkTokenWithdraw represents a TokenWithdraw event raised by the KyberNetwork contract.
type KyberNetworkTokenWithdraw struct {
	Token  common.Address
	Amount *big.Int
	SendTo common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenWithdraw is a free log retrieval operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_KyberNetwork *KyberNetworkFilterer) FilterTokenWithdraw(opts *bind.FilterOpts) (*KyberNetworkTokenWithdrawIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkTokenWithdrawIterator{contract: _KyberNetwork.contract, event: "TokenWithdraw", logs: logs, sub: sub}, nil
}

// WatchTokenWithdraw is a free log subscription operation binding the contract event 0x72cb8a894ddb372ceec3d2a7648d86f17d5a15caae0e986c53109b8a9a9385e6.
//
// Solidity: event TokenWithdraw(address token, uint256 amount, address sendTo)
func (_KyberNetwork *KyberNetworkFilterer) WatchTokenWithdraw(opts *bind.WatchOpts, sink chan<- *KyberNetworkTokenWithdraw) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "TokenWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkTokenWithdraw)
				if err := _KyberNetwork.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseTokenWithdraw(log types.Log) (*KyberNetworkTokenWithdraw, error) {
	event := new(KyberNetworkTokenWithdraw)
	if err := _KyberNetwork.contract.UnpackLog(event, "TokenWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KyberNetworkTransferAdminPendingIterator is returned from FilterTransferAdminPending and is used to iterate over the raw logs and unpacked data for TransferAdminPending events raised by the KyberNetwork contract.
type KyberNetworkTransferAdminPendingIterator struct {
	Event *KyberNetworkTransferAdminPending // Event containing the contract specifics and raw log

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
func (it *KyberNetworkTransferAdminPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KyberNetworkTransferAdminPending)
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
		it.Event = new(KyberNetworkTransferAdminPending)
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
func (it *KyberNetworkTransferAdminPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KyberNetworkTransferAdminPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KyberNetworkTransferAdminPending represents a TransferAdminPending event raised by the KyberNetwork contract.
type KyberNetworkTransferAdminPending struct {
	PendingAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminPending is a free log retrieval operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberNetwork *KyberNetworkFilterer) FilterTransferAdminPending(opts *bind.FilterOpts) (*KyberNetworkTransferAdminPendingIterator, error) {

	logs, sub, err := _KyberNetwork.contract.FilterLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return &KyberNetworkTransferAdminPendingIterator{contract: _KyberNetwork.contract, event: "TransferAdminPending", logs: logs, sub: sub}, nil
}

// WatchTransferAdminPending is a free log subscription operation binding the contract event 0x3b81caf78fa51ecbc8acb482fd7012a277b428d9b80f9d156e8a54107496cc40.
//
// Solidity: event TransferAdminPending(address pendingAdmin)
func (_KyberNetwork *KyberNetworkFilterer) WatchTransferAdminPending(opts *bind.WatchOpts, sink chan<- *KyberNetworkTransferAdminPending) (event.Subscription, error) {

	logs, sub, err := _KyberNetwork.contract.WatchLogs(opts, "TransferAdminPending")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KyberNetworkTransferAdminPending)
				if err := _KyberNetwork.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
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
func (_KyberNetwork *KyberNetworkFilterer) ParseTransferAdminPending(log types.Log) (*KyberNetworkTransferAdminPending, error) {
	event := new(KyberNetworkTransferAdminPending)
	if err := _KyberNetwork.contract.UnpackLog(event, "TransferAdminPending", log); err != nil {
		return nil, err
	}
	return event, nil
}
