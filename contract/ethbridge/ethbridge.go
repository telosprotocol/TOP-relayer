// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethbridge

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

// NearBridgeBridgeState is an auto generated low-level Go binding around an user-defined struct.
type NearBridgeBridgeState struct {
	CurrentHeight     *big.Int
	NextTimestamp     *big.Int
	NumBlockProducers *big.Int
}

// EthBridgeMetaData contains all meta data concerning the EthBridge contract.
var EthBridgeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashReverted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDBLOCK_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_BURN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_MINT_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTROLLED_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_initializing\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"addLightClientBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"flags\",\"type\":\"uint256\"}],\"name\":\"adminPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"blockHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"res\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"blockMerkleRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"res\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"currentHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numBlockProducers\",\"type\":\"uint256\"}],\"internalType\":\"structNearBridge.BridgeState\",\"name\":\"res\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashCode\",\"type\":\"bytes32\"}],\"name\":\"getHeightByHash\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"initWithBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractEd25519\",\"name\":\"ed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_lockEthAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockEthAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EthBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use EthBridgeMetaData.ABI instead.
var EthBridgeABI = EthBridgeMetaData.ABI

// EthBridge is an auto generated Go binding around an Ethereum contract.
type EthBridge struct {
	EthBridgeCaller     // Read-only binding to the contract
	EthBridgeTransactor // Write-only binding to the contract
	EthBridgeFilterer   // Log filterer for contract events
}

// EthBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthBridgeSession struct {
	Contract     *EthBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthBridgeCallerSession struct {
	Contract *EthBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// EthBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthBridgeTransactorSession struct {
	Contract     *EthBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// EthBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthBridgeRaw struct {
	Contract *EthBridge // Generic contract binding to access the raw methods on
}

// EthBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthBridgeCallerRaw struct {
	Contract *EthBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// EthBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthBridgeTransactorRaw struct {
	Contract *EthBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthBridge creates a new instance of EthBridge, bound to a specific deployed contract.
func NewEthBridge(address common.Address, backend bind.ContractBackend) (*EthBridge, error) {
	contract, err := bindEthBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthBridge{EthBridgeCaller: EthBridgeCaller{contract: contract}, EthBridgeTransactor: EthBridgeTransactor{contract: contract}, EthBridgeFilterer: EthBridgeFilterer{contract: contract}}, nil
}

// NewEthBridgeCaller creates a new read-only instance of EthBridge, bound to a specific deployed contract.
func NewEthBridgeCaller(address common.Address, caller bind.ContractCaller) (*EthBridgeCaller, error) {
	contract, err := bindEthBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthBridgeCaller{contract: contract}, nil
}

// NewEthBridgeTransactor creates a new write-only instance of EthBridge, bound to a specific deployed contract.
func NewEthBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*EthBridgeTransactor, error) {
	contract, err := bindEthBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthBridgeTransactor{contract: contract}, nil
}

// NewEthBridgeFilterer creates a new log filterer instance of EthBridge, bound to a specific deployed contract.
func NewEthBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*EthBridgeFilterer, error) {
	contract, err := bindEthBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthBridgeFilterer{contract: contract}, nil
}

// bindEthBridge binds a generic wrapper to an already deployed contract.
func bindEthBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthBridge *EthBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthBridge.Contract.EthBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthBridge *EthBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBridge.Contract.EthBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthBridge *EthBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthBridge.Contract.EthBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthBridge *EthBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthBridge *EthBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthBridge *EthBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthBridge.Contract.contract.Transact(opts, method, params...)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) ADDBLOCKROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "ADDBLOCK_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) ADDBLOCKROLE() ([32]byte, error) {
	return _EthBridge.Contract.ADDBLOCKROLE(&_EthBridge.CallOpts)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) ADDBLOCKROLE() ([32]byte, error) {
	return _EthBridge.Contract.ADDBLOCKROLE(&_EthBridge.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) BLACKBURNROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "BLACK_BURN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) BLACKBURNROLE() ([32]byte, error) {
	return _EthBridge.Contract.BLACKBURNROLE(&_EthBridge.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) BLACKBURNROLE() ([32]byte, error) {
	return _EthBridge.Contract.BLACKBURNROLE(&_EthBridge.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) BLACKMINTROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "BLACK_MINT_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) BLACKMINTROLE() ([32]byte, error) {
	return _EthBridge.Contract.BLACKMINTROLE(&_EthBridge.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) BLACKMINTROLE() ([32]byte, error) {
	return _EthBridge.Contract.BLACKMINTROLE(&_EthBridge.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) CONTROLLEDROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "CONTROLLED_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) CONTROLLEDROLE() ([32]byte, error) {
	return _EthBridge.Contract.CONTROLLEDROLE(&_EthBridge.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) CONTROLLEDROLE() ([32]byte, error) {
	return _EthBridge.Contract.CONTROLLEDROLE(&_EthBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _EthBridge.Contract.DEFAULTADMINROLE(&_EthBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _EthBridge.Contract.DEFAULTADMINROLE(&_EthBridge.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) OWNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "OWNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) OWNERROLE() ([32]byte, error) {
	return _EthBridge.Contract.OWNERROLE(&_EthBridge.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) OWNERROLE() ([32]byte, error) {
	return _EthBridge.Contract.OWNERROLE(&_EthBridge.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCaller) WITHDRAWALROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "WITHDRAWAL_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeSession) WITHDRAWALROLE() ([32]byte, error) {
	return _EthBridge.Contract.WITHDRAWALROLE(&_EthBridge.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) WITHDRAWALROLE() ([32]byte, error) {
	return _EthBridge.Contract.WITHDRAWALROLE(&_EthBridge.CallOpts)
}

// Initialized1 is a free data retrieval call binding the contract method 0x3072cf60.
//
// Solidity: function _initialized() view returns(bool)
func (_EthBridge *EthBridgeCaller) Initialized1(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "_initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized1 is a free data retrieval call binding the contract method 0x3072cf60.
//
// Solidity: function _initialized() view returns(bool)
func (_EthBridge *EthBridgeSession) Initialized1() (bool, error) {
	return _EthBridge.Contract.Initialized1(&_EthBridge.CallOpts)
}

// Initialized1 is a free data retrieval call binding the contract method 0x3072cf60.
//
// Solidity: function _initialized() view returns(bool)
func (_EthBridge *EthBridgeCallerSession) Initialized1() (bool, error) {
	return _EthBridge.Contract.Initialized1(&_EthBridge.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_EthBridge *EthBridgeCaller) Initializing(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "_initializing")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_EthBridge *EthBridgeSession) Initializing() (bool, error) {
	return _EthBridge.Contract.Initializing(&_EthBridge.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_EthBridge *EthBridgeCallerSession) Initializing() (bool, error) {
	return _EthBridge.Contract.Initializing(&_EthBridge.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_EthBridge *EthBridgeCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_EthBridge *EthBridgeSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _EthBridge.Contract.BalanceOf(&_EthBridge.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_EthBridge *EthBridgeCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _EthBridge.Contract.BalanceOf(&_EthBridge.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x37da8ec5.
//
// Solidity: function blockHashes(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeCaller) BlockHashes(opts *bind.CallOpts, height uint64) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "blockHashes", height)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockHashes is a free data retrieval call binding the contract method 0x37da8ec5.
//
// Solidity: function blockHashes(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeSession) BlockHashes(height uint64) ([32]byte, error) {
	return _EthBridge.Contract.BlockHashes(&_EthBridge.CallOpts, height)
}

// BlockHashes is a free data retrieval call binding the contract method 0x37da8ec5.
//
// Solidity: function blockHashes(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeCallerSession) BlockHashes(height uint64) ([32]byte, error) {
	return _EthBridge.Contract.BlockHashes(&_EthBridge.CallOpts, height)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeCaller) BlockMerkleRoots(opts *bind.CallOpts, height uint64) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "blockMerkleRoots", height)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeSession) BlockMerkleRoots(height uint64) ([32]byte, error) {
	return _EthBridge.Contract.BlockMerkleRoots(&_EthBridge.CallOpts, height)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 height) view returns(bytes32 res)
func (_EthBridge *EthBridgeCallerSession) BlockMerkleRoots(height uint64) ([32]byte, error) {
	return _EthBridge.Contract.BlockMerkleRoots(&_EthBridge.CallOpts, height)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_EthBridge *EthBridgeCaller) BridgeState(opts *bind.CallOpts) (NearBridgeBridgeState, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "bridgeState")

	if err != nil {
		return *new(NearBridgeBridgeState), err
	}

	out0 := *abi.ConvertType(out[0], new(NearBridgeBridgeState)).(*NearBridgeBridgeState)

	return out0, err

}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_EthBridge *EthBridgeSession) BridgeState() (NearBridgeBridgeState, error) {
	return _EthBridge.Contract.BridgeState(&_EthBridge.CallOpts)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_EthBridge *EthBridgeCallerSession) BridgeState() (NearBridgeBridgeState, error) {
	return _EthBridge.Contract.BridgeState(&_EthBridge.CallOpts)
}

// GetHeightByHash is a free data retrieval call binding the contract method 0x83bfc629.
//
// Solidity: function getHeightByHash(bytes32 hashCode) view returns(uint64 height)
func (_EthBridge *EthBridgeCaller) GetHeightByHash(opts *bind.CallOpts, hashCode [32]byte) (uint64, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "getHeightByHash", hashCode)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetHeightByHash is a free data retrieval call binding the contract method 0x83bfc629.
//
// Solidity: function getHeightByHash(bytes32 hashCode) view returns(uint64 height)
func (_EthBridge *EthBridgeSession) GetHeightByHash(hashCode [32]byte) (uint64, error) {
	return _EthBridge.Contract.GetHeightByHash(&_EthBridge.CallOpts, hashCode)
}

// GetHeightByHash is a free data retrieval call binding the contract method 0x83bfc629.
//
// Solidity: function getHeightByHash(bytes32 hashCode) view returns(uint64 height)
func (_EthBridge *EthBridgeCallerSession) GetHeightByHash(hashCode [32]byte) (uint64, error) {
	return _EthBridge.Contract.GetHeightByHash(&_EthBridge.CallOpts, hashCode)
}

// GetMaxHeight is a free data retrieval call binding the contract method 0xcf22b577.
//
// Solidity: function getMaxHeight() view returns(uint64 height)
func (_EthBridge *EthBridgeCaller) GetMaxHeight(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "getMaxHeight")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetMaxHeight is a free data retrieval call binding the contract method 0xcf22b577.
//
// Solidity: function getMaxHeight() view returns(uint64 height)
func (_EthBridge *EthBridgeSession) GetMaxHeight() (uint64, error) {
	return _EthBridge.Contract.GetMaxHeight(&_EthBridge.CallOpts)
}

// GetMaxHeight is a free data retrieval call binding the contract method 0xcf22b577.
//
// Solidity: function getMaxHeight() view returns(uint64 height)
func (_EthBridge *EthBridgeCallerSession) GetMaxHeight() (uint64, error) {
	return _EthBridge.Contract.GetMaxHeight(&_EthBridge.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthBridge *EthBridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthBridge *EthBridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _EthBridge.Contract.GetRoleAdmin(&_EthBridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthBridge *EthBridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _EthBridge.Contract.GetRoleAdmin(&_EthBridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthBridge *EthBridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthBridge *EthBridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _EthBridge.Contract.HasRole(&_EthBridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthBridge *EthBridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _EthBridge.Contract.HasRole(&_EthBridge.CallOpts, role, account)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_EthBridge *EthBridgeCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_EthBridge *EthBridgeSession) Initialized() (bool, error) {
	return _EthBridge.Contract.Initialized(&_EthBridge.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_EthBridge *EthBridgeCallerSession) Initialized() (bool, error) {
	return _EthBridge.Contract.Initialized(&_EthBridge.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_EthBridge *EthBridgeCaller) LockEthAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "lockEthAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_EthBridge *EthBridgeSession) LockEthAmount() (*big.Int, error) {
	return _EthBridge.Contract.LockEthAmount(&_EthBridge.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_EthBridge *EthBridgeCallerSession) LockEthAmount() (*big.Int, error) {
	return _EthBridge.Contract.LockEthAmount(&_EthBridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EthBridge *EthBridgeCaller) Paused(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EthBridge *EthBridgeSession) Paused() (*big.Int, error) {
	return _EthBridge.Contract.Paused(&_EthBridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_EthBridge *EthBridgeCallerSession) Paused() (*big.Int, error) {
	return _EthBridge.Contract.Paused(&_EthBridge.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthBridge *EthBridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _EthBridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthBridge *EthBridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthBridge.Contract.SupportsInterface(&_EthBridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthBridge *EthBridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthBridge.Contract.SupportsInterface(&_EthBridge.CallOpts, interfaceId)
}

// AddLightClientBlock is a paid mutator transaction binding the contract method 0x6d2d6ae0.
//
// Solidity: function addLightClientBlock(bytes data) returns()
func (_EthBridge *EthBridgeTransactor) AddLightClientBlock(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "addLightClientBlock", data)
}

// AddLightClientBlock is a paid mutator transaction binding the contract method 0x6d2d6ae0.
//
// Solidity: function addLightClientBlock(bytes data) returns()
func (_EthBridge *EthBridgeSession) AddLightClientBlock(data []byte) (*types.Transaction, error) {
	return _EthBridge.Contract.AddLightClientBlock(&_EthBridge.TransactOpts, data)
}

// AddLightClientBlock is a paid mutator transaction binding the contract method 0x6d2d6ae0.
//
// Solidity: function addLightClientBlock(bytes data) returns()
func (_EthBridge *EthBridgeTransactorSession) AddLightClientBlock(data []byte) (*types.Transaction, error) {
	return _EthBridge.Contract.AddLightClientBlock(&_EthBridge.TransactOpts, data)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_EthBridge *EthBridgeTransactor) AdminPause(opts *bind.TransactOpts, flags *big.Int) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "adminPause", flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_EthBridge *EthBridgeSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _EthBridge.Contract.AdminPause(&_EthBridge.TransactOpts, flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_EthBridge *EthBridgeTransactorSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _EthBridge.Contract.AdminPause(&_EthBridge.TransactOpts, flags)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_EthBridge *EthBridgeTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_EthBridge *EthBridgeSession) Deposit() (*types.Transaction, error) {
	return _EthBridge.Contract.Deposit(&_EthBridge.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_EthBridge *EthBridgeTransactorSession) Deposit() (*types.Transaction, error) {
	return _EthBridge.Contract.Deposit(&_EthBridge.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.GrantRole(&_EthBridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.GrantRole(&_EthBridge.TransactOpts, role, account)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_EthBridge *EthBridgeTransactor) InitWithBlock(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "initWithBlock", data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_EthBridge *EthBridgeSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _EthBridge.Contract.InitWithBlock(&_EthBridge.TransactOpts, data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_EthBridge *EthBridgeTransactorSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _EthBridge.Contract.InitWithBlock(&_EthBridge.TransactOpts, data)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address ed, uint256 _lockEthAmount, address _owner) returns()
func (_EthBridge *EthBridgeTransactor) Initialize(opts *bind.TransactOpts, ed common.Address, _lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "initialize", ed, _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address ed, uint256 _lockEthAmount, address _owner) returns()
func (_EthBridge *EthBridgeSession) Initialize(ed common.Address, _lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.Initialize(&_EthBridge.TransactOpts, ed, _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address ed, uint256 _lockEthAmount, address _owner) returns()
func (_EthBridge *EthBridgeTransactorSession) Initialize(ed common.Address, _lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.Initialize(&_EthBridge.TransactOpts, ed, _lockEthAmount, _owner)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.RenounceRole(&_EthBridge.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.RenounceRole(&_EthBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.RevokeRole(&_EthBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthBridge *EthBridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthBridge.Contract.RevokeRole(&_EthBridge.TransactOpts, role, account)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_EthBridge *EthBridgeTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBridge.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_EthBridge *EthBridgeSession) Withdraw() (*types.Transaction, error) {
	return _EthBridge.Contract.Withdraw(&_EthBridge.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_EthBridge *EthBridgeTransactorSession) Withdraw() (*types.Transaction, error) {
	return _EthBridge.Contract.Withdraw(&_EthBridge.TransactOpts)
}

// EthBridgeBlockHashAddedIterator is returned from FilterBlockHashAdded and is used to iterate over the raw logs and unpacked data for BlockHashAdded events raised by the EthBridge contract.
type EthBridgeBlockHashAddedIterator struct {
	Event *EthBridgeBlockHashAdded // Event containing the contract specifics and raw log

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
func (it *EthBridgeBlockHashAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthBridgeBlockHashAdded)
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
		it.Event = new(EthBridgeBlockHashAdded)
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
func (it *EthBridgeBlockHashAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthBridgeBlockHashAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthBridgeBlockHashAdded represents a BlockHashAdded event raised by the EthBridge contract.
type EthBridgeBlockHashAdded struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashAdded is a free log retrieval operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_EthBridge *EthBridgeFilterer) FilterBlockHashAdded(opts *bind.FilterOpts, height []uint64) (*EthBridgeBlockHashAddedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _EthBridge.contract.FilterLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return &EthBridgeBlockHashAddedIterator{contract: _EthBridge.contract, event: "BlockHashAdded", logs: logs, sub: sub}, nil
}

// WatchBlockHashAdded is a free log subscription operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_EthBridge *EthBridgeFilterer) WatchBlockHashAdded(opts *bind.WatchOpts, sink chan<- *EthBridgeBlockHashAdded, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _EthBridge.contract.WatchLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthBridgeBlockHashAdded)
				if err := _EthBridge.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
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
func (_EthBridge *EthBridgeFilterer) ParseBlockHashAdded(log types.Log) (*EthBridgeBlockHashAdded, error) {
	event := new(EthBridgeBlockHashAdded)
	if err := _EthBridge.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthBridgeBlockHashRevertedIterator is returned from FilterBlockHashReverted and is used to iterate over the raw logs and unpacked data for BlockHashReverted events raised by the EthBridge contract.
type EthBridgeBlockHashRevertedIterator struct {
	Event *EthBridgeBlockHashReverted // Event containing the contract specifics and raw log

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
func (it *EthBridgeBlockHashRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthBridgeBlockHashReverted)
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
		it.Event = new(EthBridgeBlockHashReverted)
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
func (it *EthBridgeBlockHashRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthBridgeBlockHashRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthBridgeBlockHashReverted represents a BlockHashReverted event raised by the EthBridge contract.
type EthBridgeBlockHashReverted struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashReverted is a free log retrieval operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_EthBridge *EthBridgeFilterer) FilterBlockHashReverted(opts *bind.FilterOpts, height []uint64) (*EthBridgeBlockHashRevertedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _EthBridge.contract.FilterLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return &EthBridgeBlockHashRevertedIterator{contract: _EthBridge.contract, event: "BlockHashReverted", logs: logs, sub: sub}, nil
}

// WatchBlockHashReverted is a free log subscription operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_EthBridge *EthBridgeFilterer) WatchBlockHashReverted(opts *bind.WatchOpts, sink chan<- *EthBridgeBlockHashReverted, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _EthBridge.contract.WatchLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthBridgeBlockHashReverted)
				if err := _EthBridge.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
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
func (_EthBridge *EthBridgeFilterer) ParseBlockHashReverted(log types.Log) (*EthBridgeBlockHashReverted, error) {
	event := new(EthBridgeBlockHashReverted)
	if err := _EthBridge.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthBridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the EthBridge contract.
type EthBridgeRoleAdminChangedIterator struct {
	Event *EthBridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *EthBridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthBridgeRoleAdminChanged)
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
		it.Event = new(EthBridgeRoleAdminChanged)
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
func (it *EthBridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthBridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthBridgeRoleAdminChanged represents a RoleAdminChanged event raised by the EthBridge contract.
type EthBridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_EthBridge *EthBridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*EthBridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _EthBridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &EthBridgeRoleAdminChangedIterator{contract: _EthBridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_EthBridge *EthBridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *EthBridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _EthBridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthBridgeRoleAdminChanged)
				if err := _EthBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_EthBridge *EthBridgeFilterer) ParseRoleAdminChanged(log types.Log) (*EthBridgeRoleAdminChanged, error) {
	event := new(EthBridgeRoleAdminChanged)
	if err := _EthBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthBridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the EthBridge contract.
type EthBridgeRoleGrantedIterator struct {
	Event *EthBridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *EthBridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthBridgeRoleGranted)
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
		it.Event = new(EthBridgeRoleGranted)
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
func (it *EthBridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthBridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthBridgeRoleGranted represents a RoleGranted event raised by the EthBridge contract.
type EthBridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthBridge *EthBridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthBridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _EthBridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthBridgeRoleGrantedIterator{contract: _EthBridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthBridge *EthBridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *EthBridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _EthBridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthBridgeRoleGranted)
				if err := _EthBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_EthBridge *EthBridgeFilterer) ParseRoleGranted(log types.Log) (*EthBridgeRoleGranted, error) {
	event := new(EthBridgeRoleGranted)
	if err := _EthBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthBridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the EthBridge contract.
type EthBridgeRoleRevokedIterator struct {
	Event *EthBridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *EthBridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthBridgeRoleRevoked)
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
		it.Event = new(EthBridgeRoleRevoked)
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
func (it *EthBridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthBridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthBridgeRoleRevoked represents a RoleRevoked event raised by the EthBridge contract.
type EthBridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthBridge *EthBridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthBridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _EthBridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthBridgeRoleRevokedIterator{contract: _EthBridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthBridge *EthBridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *EthBridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _EthBridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthBridgeRoleRevoked)
				if err := _EthBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_EthBridge *EthBridgeFilterer) ParseRoleRevoked(log types.Log) (*EthBridgeRoleRevoked, error) {
	event := new(EthBridgeRoleRevoked)
	if err := _EthBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
