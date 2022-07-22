// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package topclient

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TopBridgeBridgeState is an auto generated low-level Go binding around an user-defined struct.
type TopBridgeBridgeState struct {
	CurrentHeight     *big.Int
	NextTimestamp     *big.Int
	NumBlockProducers *big.Int
}

// TopClientMetaData contains all meta data concerning the TopClient contract.
var TopClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"addLightClientBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"flags\",\"type\":\"uint256\"}],\"name\":\"adminPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashReverted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lockEthAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"initWithBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_initializing\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ADDBLOCK_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_BURN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_MINT_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"blockHashes\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"blockHeights\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"blockMerkleRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"currentHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numBlockProducers\",\"type\":\"uint256\"}],\"internalType\":\"structTopBridge.BridgeState\",\"name\":\"res\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTROLLED_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSubmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockEthAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxMainHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TopClientABI is the input ABI used to generate the binding from.
// Deprecated: Use TopClientMetaData.ABI instead.
var TopClientABI = TopClientMetaData.ABI

// TopClient is an auto generated Go binding around an Ethereum contract.
type TopClient struct {
	TopClientCaller     // Read-only binding to the contract
	TopClientTransactor // Write-only binding to the contract
	TopClientFilterer   // Log filterer for contract events
}

// TopClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type TopClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TopClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TopClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TopClientSession struct {
	Contract     *TopClient        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TopClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TopClientCallerSession struct {
	Contract *TopClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TopClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TopClientTransactorSession struct {
	Contract     *TopClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TopClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type TopClientRaw struct {
	Contract *TopClient // Generic contract binding to access the raw methods on
}

// TopClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TopClientCallerRaw struct {
	Contract *TopClientCaller // Generic read-only contract binding to access the raw methods on
}

// TopClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TopClientTransactorRaw struct {
	Contract *TopClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTopClient creates a new instance of TopClient, bound to a specific deployed contract.
func NewTopClient(address common.Address, backend bind.ContractBackend) (*TopClient, error) {
	contract, err := bindTopClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TopClient{TopClientCaller: TopClientCaller{contract: contract}, TopClientTransactor: TopClientTransactor{contract: contract}, TopClientFilterer: TopClientFilterer{contract: contract}}, nil
}

// NewTopClientCaller creates a new read-only instance of TopClient, bound to a specific deployed contract.
func NewTopClientCaller(address common.Address, caller bind.ContractCaller) (*TopClientCaller, error) {
	contract, err := bindTopClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TopClientCaller{contract: contract}, nil
}

// NewTopClientTransactor creates a new write-only instance of TopClient, bound to a specific deployed contract.
func NewTopClientTransactor(address common.Address, transactor bind.ContractTransactor) (*TopClientTransactor, error) {
	contract, err := bindTopClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TopClientTransactor{contract: contract}, nil
}

// NewTopClientFilterer creates a new log filterer instance of TopClient, bound to a specific deployed contract.
func NewTopClientFilterer(address common.Address, filterer bind.ContractFilterer) (*TopClientFilterer, error) {
	contract, err := bindTopClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TopClientFilterer{contract: contract}, nil
}

// bindTopClient binds a generic wrapper to an already deployed contract.
func bindTopClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TopClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TopClient *TopClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TopClient.Contract.TopClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TopClient *TopClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TopClient.Contract.TopClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TopClient *TopClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TopClient.Contract.TopClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TopClient *TopClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TopClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TopClient *TopClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TopClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TopClient *TopClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TopClient.Contract.contract.Transact(opts, method, params...)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) ADDBLOCKROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "ADDBLOCK_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) ADDBLOCKROLE() ([32]byte, error) {
	return _TopClient.Contract.ADDBLOCKROLE(&_TopClient.CallOpts)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) ADDBLOCKROLE() ([32]byte, error) {
	return _TopClient.Contract.ADDBLOCKROLE(&_TopClient.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) BLACKBURNROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "BLACK_BURN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) BLACKBURNROLE() ([32]byte, error) {
	return _TopClient.Contract.BLACKBURNROLE(&_TopClient.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) BLACKBURNROLE() ([32]byte, error) {
	return _TopClient.Contract.BLACKBURNROLE(&_TopClient.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) BLACKMINTROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "BLACK_MINT_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) BLACKMINTROLE() ([32]byte, error) {
	return _TopClient.Contract.BLACKMINTROLE(&_TopClient.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) BLACKMINTROLE() ([32]byte, error) {
	return _TopClient.Contract.BLACKMINTROLE(&_TopClient.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) CONTROLLEDROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "CONTROLLED_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) CONTROLLEDROLE() ([32]byte, error) {
	return _TopClient.Contract.CONTROLLEDROLE(&_TopClient.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) CONTROLLEDROLE() ([32]byte, error) {
	return _TopClient.Contract.CONTROLLEDROLE(&_TopClient.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TopClient.Contract.DEFAULTADMINROLE(&_TopClient.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TopClient.Contract.DEFAULTADMINROLE(&_TopClient.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) OWNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "OWNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) OWNERROLE() ([32]byte, error) {
	return _TopClient.Contract.OWNERROLE(&_TopClient.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) OWNERROLE() ([32]byte, error) {
	return _TopClient.Contract.OWNERROLE(&_TopClient.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_TopClient *TopClientCaller) WITHDRAWALROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "WITHDRAWAL_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_TopClient *TopClientSession) WITHDRAWALROLE() ([32]byte, error) {
	return _TopClient.Contract.WITHDRAWALROLE(&_TopClient.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_TopClient *TopClientCallerSession) WITHDRAWALROLE() ([32]byte, error) {
	return _TopClient.Contract.WITHDRAWALROLE(&_TopClient.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_TopClient *TopClientCaller) Initializing(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "_initializing")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_TopClient *TopClientSession) Initializing() (bool, error) {
	return _TopClient.Contract.Initializing(&_TopClient.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_TopClient *TopClientCallerSession) Initializing() (bool, error) {
	return _TopClient.Contract.Initializing(&_TopClient.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_TopClient *TopClientCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_TopClient *TopClientSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _TopClient.Contract.BalanceOf(&_TopClient.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_TopClient *TopClientCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _TopClient.Contract.BalanceOf(&_TopClient.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_TopClient *TopClientCaller) BlockHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "blockHashes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_TopClient *TopClientSession) BlockHashes(arg0 [32]byte) (bool, error) {
	return _TopClient.Contract.BlockHashes(&_TopClient.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_TopClient *TopClientCallerSession) BlockHashes(arg0 [32]byte) (bool, error) {
	return _TopClient.Contract.BlockHashes(&_TopClient.CallOpts, arg0)
}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_TopClient *TopClientCaller) BlockHeights(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "blockHeights", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_TopClient *TopClientSession) BlockHeights(arg0 uint64) (bool, error) {
	return _TopClient.Contract.BlockHeights(&_TopClient.CallOpts, arg0)
}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_TopClient *TopClientCallerSession) BlockHeights(arg0 uint64) (bool, error) {
	return _TopClient.Contract.BlockHeights(&_TopClient.CallOpts, arg0)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_TopClient *TopClientCaller) BlockMerkleRoots(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "blockMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_TopClient *TopClientSession) BlockMerkleRoots(arg0 uint64) ([32]byte, error) {
	return _TopClient.Contract.BlockMerkleRoots(&_TopClient.CallOpts, arg0)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_TopClient *TopClientCallerSession) BlockMerkleRoots(arg0 uint64) ([32]byte, error) {
	return _TopClient.Contract.BlockMerkleRoots(&_TopClient.CallOpts, arg0)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_TopClient *TopClientCaller) BridgeState(opts *bind.CallOpts) (TopBridgeBridgeState, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "bridgeState")

	if err != nil {
		return *new(TopBridgeBridgeState), err
	}

	out0 := *abi.ConvertType(out[0], new(TopBridgeBridgeState)).(*TopBridgeBridgeState)

	return out0, err

}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_TopClient *TopClientSession) BridgeState() (TopBridgeBridgeState, error) {
	return _TopClient.Contract.BridgeState(&_TopClient.CallOpts)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_TopClient *TopClientCallerSession) BridgeState() (TopBridgeBridgeState, error) {
	return _TopClient.Contract.BridgeState(&_TopClient.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TopClient *TopClientCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TopClient *TopClientSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TopClient.Contract.GetRoleAdmin(&_TopClient.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TopClient *TopClientCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TopClient.Contract.GetRoleAdmin(&_TopClient.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TopClient *TopClientCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TopClient *TopClientSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TopClient.Contract.HasRole(&_TopClient.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TopClient *TopClientCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TopClient.Contract.HasRole(&_TopClient.CallOpts, role, account)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_TopClient *TopClientCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_TopClient *TopClientSession) Initialized() (bool, error) {
	return _TopClient.Contract.Initialized(&_TopClient.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_TopClient *TopClientCallerSession) Initialized() (bool, error) {
	return _TopClient.Contract.Initialized(&_TopClient.CallOpts)
}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_TopClient *TopClientCaller) LastSubmitter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "lastSubmitter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_TopClient *TopClientSession) LastSubmitter() (common.Address, error) {
	return _TopClient.Contract.LastSubmitter(&_TopClient.CallOpts)
}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_TopClient *TopClientCallerSession) LastSubmitter() (common.Address, error) {
	return _TopClient.Contract.LastSubmitter(&_TopClient.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_TopClient *TopClientCaller) LockEthAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "lockEthAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_TopClient *TopClientSession) LockEthAmount() (*big.Int, error) {
	return _TopClient.Contract.LockEthAmount(&_TopClient.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_TopClient *TopClientCallerSession) LockEthAmount() (*big.Int, error) {
	return _TopClient.Contract.LockEthAmount(&_TopClient.CallOpts)
}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_TopClient *TopClientCaller) MaxMainHeight(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "maxMainHeight")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_TopClient *TopClientSession) MaxMainHeight() (uint64, error) {
	return _TopClient.Contract.MaxMainHeight(&_TopClient.CallOpts)
}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_TopClient *TopClientCallerSession) MaxMainHeight() (uint64, error) {
	return _TopClient.Contract.MaxMainHeight(&_TopClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_TopClient *TopClientCaller) Paused(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_TopClient *TopClientSession) Paused() (*big.Int, error) {
	return _TopClient.Contract.Paused(&_TopClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_TopClient *TopClientCallerSession) Paused() (*big.Int, error) {
	return _TopClient.Contract.Paused(&_TopClient.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TopClient *TopClientCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TopClient.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TopClient *TopClientSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TopClient.Contract.SupportsInterface(&_TopClient.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TopClient *TopClientCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TopClient.Contract.SupportsInterface(&_TopClient.CallOpts, interfaceId)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_TopClient *TopClientTransactor) AddLightClientBlocks(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "addLightClientBlocks", data)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_TopClient *TopClientSession) AddLightClientBlocks(data []byte) (*types.Transaction, error) {
	return _TopClient.Contract.AddLightClientBlocks(&_TopClient.TransactOpts, data)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_TopClient *TopClientTransactorSession) AddLightClientBlocks(data []byte) (*types.Transaction, error) {
	return _TopClient.Contract.AddLightClientBlocks(&_TopClient.TransactOpts, data)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_TopClient *TopClientTransactor) AdminPause(opts *bind.TransactOpts, flags *big.Int) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "adminPause", flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_TopClient *TopClientSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _TopClient.Contract.AdminPause(&_TopClient.TransactOpts, flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_TopClient *TopClientTransactorSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _TopClient.Contract.AdminPause(&_TopClient.TransactOpts, flags)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TopClient *TopClientSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.GrantRole(&_TopClient.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.GrantRole(&_TopClient.TransactOpts, role, account)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_TopClient *TopClientTransactor) InitWithBlock(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "initWithBlock", data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_TopClient *TopClientSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _TopClient.Contract.InitWithBlock(&_TopClient.TransactOpts, data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_TopClient *TopClientTransactorSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _TopClient.Contract.InitWithBlock(&_TopClient.TransactOpts, data)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_TopClient *TopClientTransactor) Initialize(opts *bind.TransactOpts, _lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "initialize", _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_TopClient *TopClientSession) Initialize(_lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.Initialize(&_TopClient.TransactOpts, _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_TopClient *TopClientTransactorSession) Initialize(_lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.Initialize(&_TopClient.TransactOpts, _lockEthAmount, _owner)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TopClient *TopClientSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.RenounceRole(&_TopClient.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.RenounceRole(&_TopClient.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TopClient *TopClientSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.RevokeRole(&_TopClient.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TopClient *TopClientTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TopClient.Contract.RevokeRole(&_TopClient.TransactOpts, role, account)
}

// TopClientBlockHashAddedIterator is returned from FilterBlockHashAdded and is used to iterate over the raw logs and unpacked data for BlockHashAdded events raised by the TopClient contract.
type TopClientBlockHashAddedIterator struct {
	Event *TopClientBlockHashAdded // Event containing the contract specifics and raw log

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
func (it *TopClientBlockHashAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TopClientBlockHashAdded)
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
		it.Event = new(TopClientBlockHashAdded)
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
func (it *TopClientBlockHashAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TopClientBlockHashAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TopClientBlockHashAdded represents a BlockHashAdded event raised by the TopClient contract.
type TopClientBlockHashAdded struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashAdded is a free log retrieval operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) FilterBlockHashAdded(opts *bind.FilterOpts, height []uint64) (*TopClientBlockHashAddedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _TopClient.contract.FilterLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return &TopClientBlockHashAddedIterator{contract: _TopClient.contract, event: "BlockHashAdded", logs: logs, sub: sub}, nil
}

// WatchBlockHashAdded is a free log subscription operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) WatchBlockHashAdded(opts *bind.WatchOpts, sink chan<- *TopClientBlockHashAdded, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _TopClient.contract.WatchLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TopClientBlockHashAdded)
				if err := _TopClient.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
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

// ParseBlockHashAdded is a log parse operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) ParseBlockHashAdded(log types.Log) (*TopClientBlockHashAdded, error) {
	event := new(TopClientBlockHashAdded)
	if err := _TopClient.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TopClientBlockHashRevertedIterator is returned from FilterBlockHashReverted and is used to iterate over the raw logs and unpacked data for BlockHashReverted events raised by the TopClient contract.
type TopClientBlockHashRevertedIterator struct {
	Event *TopClientBlockHashReverted // Event containing the contract specifics and raw log

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
func (it *TopClientBlockHashRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TopClientBlockHashReverted)
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
		it.Event = new(TopClientBlockHashReverted)
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
func (it *TopClientBlockHashRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TopClientBlockHashRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TopClientBlockHashReverted represents a BlockHashReverted event raised by the TopClient contract.
type TopClientBlockHashReverted struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashReverted is a free log retrieval operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) FilterBlockHashReverted(opts *bind.FilterOpts, height []uint64) (*TopClientBlockHashRevertedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _TopClient.contract.FilterLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return &TopClientBlockHashRevertedIterator{contract: _TopClient.contract, event: "BlockHashReverted", logs: logs, sub: sub}, nil
}

// WatchBlockHashReverted is a free log subscription operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) WatchBlockHashReverted(opts *bind.WatchOpts, sink chan<- *TopClientBlockHashReverted, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _TopClient.contract.WatchLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TopClientBlockHashReverted)
				if err := _TopClient.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
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

// ParseBlockHashReverted is a log parse operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_TopClient *TopClientFilterer) ParseBlockHashReverted(log types.Log) (*TopClientBlockHashReverted, error) {
	event := new(TopClientBlockHashReverted)
	if err := _TopClient.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TopClientRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TopClient contract.
type TopClientRoleAdminChangedIterator struct {
	Event *TopClientRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TopClientRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TopClientRoleAdminChanged)
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
		it.Event = new(TopClientRoleAdminChanged)
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
func (it *TopClientRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TopClientRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TopClientRoleAdminChanged represents a RoleAdminChanged event raised by the TopClient contract.
type TopClientRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TopClient *TopClientFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TopClientRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TopClient.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TopClientRoleAdminChangedIterator{contract: _TopClient.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TopClient *TopClientFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TopClientRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TopClient.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TopClientRoleAdminChanged)
				if err := _TopClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TopClient *TopClientFilterer) ParseRoleAdminChanged(log types.Log) (*TopClientRoleAdminChanged, error) {
	event := new(TopClientRoleAdminChanged)
	if err := _TopClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TopClientRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TopClient contract.
type TopClientRoleGrantedIterator struct {
	Event *TopClientRoleGranted // Event containing the contract specifics and raw log

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
func (it *TopClientRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TopClientRoleGranted)
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
		it.Event = new(TopClientRoleGranted)
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
func (it *TopClientRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TopClientRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TopClientRoleGranted represents a RoleGranted event raised by the TopClient contract.
type TopClientRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TopClientRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TopClient.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TopClientRoleGrantedIterator{contract: _TopClient.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TopClientRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TopClient.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TopClientRoleGranted)
				if err := _TopClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) ParseRoleGranted(log types.Log) (*TopClientRoleGranted, error) {
	event := new(TopClientRoleGranted)
	if err := _TopClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TopClientRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TopClient contract.
type TopClientRoleRevokedIterator struct {
	Event *TopClientRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TopClientRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TopClientRoleRevoked)
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
		it.Event = new(TopClientRoleRevoked)
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
func (it *TopClientRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TopClientRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TopClientRoleRevoked represents a RoleRevoked event raised by the TopClient contract.
type TopClientRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TopClientRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TopClient.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TopClientRoleRevokedIterator{contract: _TopClient.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TopClientRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TopClient.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TopClientRoleRevoked)
				if err := _TopClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TopClient *TopClientFilterer) ParseRoleRevoked(log types.Log) (*TopClientRoleRevoked, error) {
	event := new(TopClientRoleRevoked)
	if err := _TopClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
