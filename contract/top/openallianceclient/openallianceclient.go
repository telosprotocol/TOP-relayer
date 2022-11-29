// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package openallianceclient

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
	_ = abi.ConvertType
)

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	CurrentHeight     *big.Int
	NextTimestamp     *big.Int
	NumBlockProducers *big.Int
}

// OpenAllianceClientMetaData contains all meta data concerning the OpenAllianceClient contract.
var OpenAllianceClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"addLightClientBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"flags\",\"type\":\"uint256\"}],\"name\":\"adminPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"BlockHashReverted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lockEthAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"initWithBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_initializing\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ADDBLOCK_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_BURN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_MINT_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"blockHashes\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"blockHeights\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"blockMerkleRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"currentHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numBlockProducers\",\"type\":\"uint256\"}],\"internalType\":\"structTopBridge.BridgeState\",\"name\":\"res\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTROLLED_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSubmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockEthAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxMainHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OpenAllianceClientABI is the input ABI used to generate the binding from.
// Deprecated: Use OpenAllianceClientMetaData.ABI instead.
var OpenAllianceClientABI = OpenAllianceClientMetaData.ABI

// OpenAllianceClient is an auto generated Go binding around an Ethereum contract.
type OpenAllianceClient struct {
	OpenAllianceClientCaller     // Read-only binding to the contract
	OpenAllianceClientTransactor // Write-only binding to the contract
	OpenAllianceClientFilterer   // Log filterer for contract events
}

// OpenAllianceClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpenAllianceClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenAllianceClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpenAllianceClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenAllianceClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpenAllianceClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenAllianceClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpenAllianceClientSession struct {
	Contract     *OpenAllianceClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OpenAllianceClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpenAllianceClientCallerSession struct {
	Contract *OpenAllianceClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OpenAllianceClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpenAllianceClientTransactorSession struct {
	Contract     *OpenAllianceClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OpenAllianceClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpenAllianceClientRaw struct {
	Contract *OpenAllianceClient // Generic contract binding to access the raw methods on
}

// OpenAllianceClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpenAllianceClientCallerRaw struct {
	Contract *OpenAllianceClientCaller // Generic read-only contract binding to access the raw methods on
}

// OpenAllianceClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpenAllianceClientTransactorRaw struct {
	Contract *OpenAllianceClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpenAllianceClient creates a new instance of OpenAllianceClient, bound to a specific deployed contract.
func NewOpenAllianceClient(address common.Address, backend bind.ContractBackend) (*OpenAllianceClient, error) {
	contract, err := bindOpenAllianceClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClient{OpenAllianceClientCaller: OpenAllianceClientCaller{contract: contract}, OpenAllianceClientTransactor: OpenAllianceClientTransactor{contract: contract}, OpenAllianceClientFilterer: OpenAllianceClientFilterer{contract: contract}}, nil
}

// NewOpenAllianceClientCaller creates a new read-only instance of OpenAllianceClient, bound to a specific deployed contract.
func NewOpenAllianceClientCaller(address common.Address, caller bind.ContractCaller) (*OpenAllianceClientCaller, error) {
	contract, err := bindOpenAllianceClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientCaller{contract: contract}, nil
}

// NewOpenAllianceClientTransactor creates a new write-only instance of OpenAllianceClient, bound to a specific deployed contract.
func NewOpenAllianceClientTransactor(address common.Address, transactor bind.ContractTransactor) (*OpenAllianceClientTransactor, error) {
	contract, err := bindOpenAllianceClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientTransactor{contract: contract}, nil
}

// NewOpenAllianceClientFilterer creates a new log filterer instance of OpenAllianceClient, bound to a specific deployed contract.
func NewOpenAllianceClientFilterer(address common.Address, filterer bind.ContractFilterer) (*OpenAllianceClientFilterer, error) {
	contract, err := bindOpenAllianceClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientFilterer{contract: contract}, nil
}

// bindOpenAllianceClient binds a generic wrapper to an already deployed contract.
func bindOpenAllianceClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OpenAllianceClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenAllianceClient *OpenAllianceClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpenAllianceClient.Contract.OpenAllianceClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenAllianceClient *OpenAllianceClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.OpenAllianceClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenAllianceClient *OpenAllianceClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.OpenAllianceClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenAllianceClient *OpenAllianceClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpenAllianceClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenAllianceClient *OpenAllianceClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenAllianceClient *OpenAllianceClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.contract.Transact(opts, method, params...)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) ADDBLOCKROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "ADDBLOCK_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) ADDBLOCKROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.ADDBLOCKROLE(&_OpenAllianceClient.CallOpts)
}

// ADDBLOCKROLE is a free data retrieval call binding the contract method 0x7af8fc10.
//
// Solidity: function ADDBLOCK_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) ADDBLOCKROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.ADDBLOCKROLE(&_OpenAllianceClient.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) BLACKBURNROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "BLACK_BURN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) BLACKBURNROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.BLACKBURNROLE(&_OpenAllianceClient.CallOpts)
}

// BLACKBURNROLE is a free data retrieval call binding the contract method 0x13eb7c0b.
//
// Solidity: function BLACK_BURN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BLACKBURNROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.BLACKBURNROLE(&_OpenAllianceClient.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) BLACKMINTROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "BLACK_MINT_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) BLACKMINTROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.BLACKMINTROLE(&_OpenAllianceClient.CallOpts)
}

// BLACKMINTROLE is a free data retrieval call binding the contract method 0xb3a6e108.
//
// Solidity: function BLACK_MINT_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BLACKMINTROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.BLACKMINTROLE(&_OpenAllianceClient.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) CONTROLLEDROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "CONTROLLED_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) CONTROLLEDROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.CONTROLLEDROLE(&_OpenAllianceClient.CallOpts)
}

// CONTROLLEDROLE is a free data retrieval call binding the contract method 0xef84e7af.
//
// Solidity: function CONTROLLED_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) CONTROLLEDROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.CONTROLLEDROLE(&_OpenAllianceClient.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.DEFAULTADMINROLE(&_OpenAllianceClient.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.DEFAULTADMINROLE(&_OpenAllianceClient.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) OWNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "OWNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) OWNERROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.OWNERROLE(&_OpenAllianceClient.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) OWNERROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.OWNERROLE(&_OpenAllianceClient.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) WITHDRAWALROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "WITHDRAWAL_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) WITHDRAWALROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.WITHDRAWALROLE(&_OpenAllianceClient.CallOpts)
}

// WITHDRAWALROLE is a free data retrieval call binding the contract method 0x67db90c2.
//
// Solidity: function WITHDRAWAL_ROLE() view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) WITHDRAWALROLE() ([32]byte, error) {
	return _OpenAllianceClient.Contract.WITHDRAWALROLE(&_OpenAllianceClient.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) Initializing(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "_initializing")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) Initializing() (bool, error) {
	return _OpenAllianceClient.Contract.Initializing(&_OpenAllianceClient.CallOpts)
}

// Initializing is a free data retrieval call binding the contract method 0xf8d27e23.
//
// Solidity: function _initializing() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) Initializing() (bool, error) {
	return _OpenAllianceClient.Contract.Initializing(&_OpenAllianceClient.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _OpenAllianceClient.Contract.BalanceOf(&_OpenAllianceClient.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _OpenAllianceClient.Contract.BalanceOf(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) BlockHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "blockHashes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) BlockHashes(arg0 [32]byte) (bool, error) {
	return _OpenAllianceClient.Contract.BlockHashes(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockHashes is a free data retrieval call binding the contract method 0x2b8a6d16.
//
// Solidity: function blockHashes(bytes32 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BlockHashes(arg0 [32]byte) (bool, error) {
	return _OpenAllianceClient.Contract.BlockHashes(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) BlockHeights(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "blockHeights", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) BlockHeights(arg0 uint64) (bool, error) {
	return _OpenAllianceClient.Contract.BlockHeights(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockHeights is a free data retrieval call binding the contract method 0xb995ac08.
//
// Solidity: function blockHeights(uint64 ) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BlockHeights(arg0 uint64) (bool, error) {
	return _OpenAllianceClient.Contract.BlockHeights(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) BlockMerkleRoots(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "blockMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) BlockMerkleRoots(arg0 uint64) ([32]byte, error) {
	return _OpenAllianceClient.Contract.BlockMerkleRoots(&_OpenAllianceClient.CallOpts, arg0)
}

// BlockMerkleRoots is a free data retrieval call binding the contract method 0x1e703806.
//
// Solidity: function blockMerkleRoots(uint64 ) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BlockMerkleRoots(arg0 uint64) ([32]byte, error) {
	return _OpenAllianceClient.Contract.BlockMerkleRoots(&_OpenAllianceClient.CallOpts, arg0)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_OpenAllianceClient *OpenAllianceClientCaller) BridgeState(opts *bind.CallOpts) (Struct0, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "bridgeState")

	if err != nil {
		return *new(Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct0)).(*Struct0)

	return out0, err

}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_OpenAllianceClient *OpenAllianceClientSession) BridgeState() (Struct0, error) {
	return _OpenAllianceClient.Contract.BridgeState(&_OpenAllianceClient.CallOpts)
}

// BridgeState is a free data retrieval call binding the contract method 0x4466ec2c.
//
// Solidity: function bridgeState() view returns((uint256,uint256,uint256) res)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) BridgeState() (Struct0, error) {
	return _OpenAllianceClient.Contract.BridgeState(&_OpenAllianceClient.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _OpenAllianceClient.Contract.GetRoleAdmin(&_OpenAllianceClient.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _OpenAllianceClient.Contract.GetRoleAdmin(&_OpenAllianceClient.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _OpenAllianceClient.Contract.HasRole(&_OpenAllianceClient.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _OpenAllianceClient.Contract.HasRole(&_OpenAllianceClient.CallOpts, role, account)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) Initialized() (bool, error) {
	return _OpenAllianceClient.Contract.Initialized(&_OpenAllianceClient.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) Initialized() (bool, error) {
	return _OpenAllianceClient.Contract.Initialized(&_OpenAllianceClient.CallOpts)
}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_OpenAllianceClient *OpenAllianceClientCaller) LastSubmitter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "lastSubmitter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_OpenAllianceClient *OpenAllianceClientSession) LastSubmitter() (common.Address, error) {
	return _OpenAllianceClient.Contract.LastSubmitter(&_OpenAllianceClient.CallOpts)
}

// LastSubmitter is a free data retrieval call binding the contract method 0x2d7dd574.
//
// Solidity: function lastSubmitter() view returns(address)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) LastSubmitter() (common.Address, error) {
	return _OpenAllianceClient.Contract.LastSubmitter(&_OpenAllianceClient.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCaller) LockEthAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "lockEthAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientSession) LockEthAmount() (*big.Int, error) {
	return _OpenAllianceClient.Contract.LockEthAmount(&_OpenAllianceClient.CallOpts)
}

// LockEthAmount is a free data retrieval call binding the contract method 0x7875a55c.
//
// Solidity: function lockEthAmount() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) LockEthAmount() (*big.Int, error) {
	return _OpenAllianceClient.Contract.LockEthAmount(&_OpenAllianceClient.CallOpts)
}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_OpenAllianceClient *OpenAllianceClientCaller) MaxMainHeight(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "maxMainHeight")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_OpenAllianceClient *OpenAllianceClientSession) MaxMainHeight() (uint64, error) {
	return _OpenAllianceClient.Contract.MaxMainHeight(&_OpenAllianceClient.CallOpts)
}

// MaxMainHeight is a free data retrieval call binding the contract method 0x966a6023.
//
// Solidity: function maxMainHeight() view returns(uint64)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) MaxMainHeight() (uint64, error) {
	return _OpenAllianceClient.Contract.MaxMainHeight(&_OpenAllianceClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCaller) Paused(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientSession) Paused() (*big.Int, error) {
	return _OpenAllianceClient.Contract.Paused(&_OpenAllianceClient.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) Paused() (*big.Int, error) {
	return _OpenAllianceClient.Contract.Paused(&_OpenAllianceClient.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _OpenAllianceClient.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _OpenAllianceClient.Contract.SupportsInterface(&_OpenAllianceClient.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OpenAllianceClient *OpenAllianceClientCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _OpenAllianceClient.Contract.SupportsInterface(&_OpenAllianceClient.CallOpts, interfaceId)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) AddLightClientBlocks(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "addLightClientBlocks", data)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) AddLightClientBlocks(data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.AddLightClientBlocks(&_OpenAllianceClient.TransactOpts, data)
}

// AddLightClientBlocks is a paid mutator transaction binding the contract method 0x7908f846.
//
// Solidity: function addLightClientBlocks(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) AddLightClientBlocks(data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.AddLightClientBlocks(&_OpenAllianceClient.TransactOpts, data)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) AdminPause(opts *bind.TransactOpts, flags *big.Int) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "adminPause", flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.AdminPause(&_OpenAllianceClient.TransactOpts, flags)
}

// AdminPause is a paid mutator transaction binding the contract method 0x2692c59f.
//
// Solidity: function adminPause(uint256 flags) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) AdminPause(flags *big.Int) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.AdminPause(&_OpenAllianceClient.TransactOpts, flags)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.GrantRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.GrantRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) InitWithBlock(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "initWithBlock", data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.InitWithBlock(&_OpenAllianceClient.TransactOpts, data)
}

// InitWithBlock is a paid mutator transaction binding the contract method 0x160bc0ba.
//
// Solidity: function initWithBlock(bytes data) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) InitWithBlock(data []byte) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.InitWithBlock(&_OpenAllianceClient.TransactOpts, data)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) Initialize(opts *bind.TransactOpts, _lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "initialize", _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) Initialize(_lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.Initialize(&_OpenAllianceClient.TransactOpts, _lockEthAmount, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 _lockEthAmount, address _owner) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) Initialize(_lockEthAmount *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.Initialize(&_OpenAllianceClient.TransactOpts, _lockEthAmount, _owner)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.RenounceRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.RenounceRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.RevokeRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_OpenAllianceClient *OpenAllianceClientTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _OpenAllianceClient.Contract.RevokeRole(&_OpenAllianceClient.TransactOpts, role, account)
}

// OpenAllianceClientBlockHashAddedIterator is returned from FilterBlockHashAdded and is used to iterate over the raw logs and unpacked data for BlockHashAdded events raised by the OpenAllianceClient contract.
type OpenAllianceClientBlockHashAddedIterator struct {
	Event *OpenAllianceClientBlockHashAdded // Event containing the contract specifics and raw log

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
func (it *OpenAllianceClientBlockHashAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpenAllianceClientBlockHashAdded)
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
		it.Event = new(OpenAllianceClientBlockHashAdded)
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
func (it *OpenAllianceClientBlockHashAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpenAllianceClientBlockHashAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpenAllianceClientBlockHashAdded represents a BlockHashAdded event raised by the OpenAllianceClient contract.
type OpenAllianceClientBlockHashAdded struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashAdded is a free log retrieval operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_OpenAllianceClient *OpenAllianceClientFilterer) FilterBlockHashAdded(opts *bind.FilterOpts, height []uint64) (*OpenAllianceClientBlockHashAddedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _OpenAllianceClient.contract.FilterLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientBlockHashAddedIterator{contract: _OpenAllianceClient.contract, event: "BlockHashAdded", logs: logs, sub: sub}, nil
}

// WatchBlockHashAdded is a free log subscription operation binding the contract event 0x5d45c22c440038a3aaf9f8134e7aa1fa59aa2a7fa411d7e818d7701c63827d7e.
//
// Solidity: event BlockHashAdded(uint64 indexed height, bytes32 blockHash)
func (_OpenAllianceClient *OpenAllianceClientFilterer) WatchBlockHashAdded(opts *bind.WatchOpts, sink chan<- *OpenAllianceClientBlockHashAdded, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _OpenAllianceClient.contract.WatchLogs(opts, "BlockHashAdded", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpenAllianceClientBlockHashAdded)
				if err := _OpenAllianceClient.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
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
func (_OpenAllianceClient *OpenAllianceClientFilterer) ParseBlockHashAdded(log types.Log) (*OpenAllianceClientBlockHashAdded, error) {
	event := new(OpenAllianceClientBlockHashAdded)
	if err := _OpenAllianceClient.contract.UnpackLog(event, "BlockHashAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OpenAllianceClientBlockHashRevertedIterator is returned from FilterBlockHashReverted and is used to iterate over the raw logs and unpacked data for BlockHashReverted events raised by the OpenAllianceClient contract.
type OpenAllianceClientBlockHashRevertedIterator struct {
	Event *OpenAllianceClientBlockHashReverted // Event containing the contract specifics and raw log

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
func (it *OpenAllianceClientBlockHashRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpenAllianceClientBlockHashReverted)
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
		it.Event = new(OpenAllianceClientBlockHashReverted)
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
func (it *OpenAllianceClientBlockHashRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpenAllianceClientBlockHashRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpenAllianceClientBlockHashReverted represents a BlockHashReverted event raised by the OpenAllianceClient contract.
type OpenAllianceClientBlockHashReverted struct {
	Height    uint64
	BlockHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHashReverted is a free log retrieval operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_OpenAllianceClient *OpenAllianceClientFilterer) FilterBlockHashReverted(opts *bind.FilterOpts, height []uint64) (*OpenAllianceClientBlockHashRevertedIterator, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _OpenAllianceClient.contract.FilterLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientBlockHashRevertedIterator{contract: _OpenAllianceClient.contract, event: "BlockHashReverted", logs: logs, sub: sub}, nil
}

// WatchBlockHashReverted is a free log subscription operation binding the contract event 0x4e9ddd5df7d5ac983348809fe8a0617e2e53415abf6f504c73ee2b2b22076ef6.
//
// Solidity: event BlockHashReverted(uint64 indexed height, bytes32 blockHash)
func (_OpenAllianceClient *OpenAllianceClientFilterer) WatchBlockHashReverted(opts *bind.WatchOpts, sink chan<- *OpenAllianceClientBlockHashReverted, height []uint64) (event.Subscription, error) {

	var heightRule []interface{}
	for _, heightItem := range height {
		heightRule = append(heightRule, heightItem)
	}

	logs, sub, err := _OpenAllianceClient.contract.WatchLogs(opts, "BlockHashReverted", heightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpenAllianceClientBlockHashReverted)
				if err := _OpenAllianceClient.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
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
func (_OpenAllianceClient *OpenAllianceClientFilterer) ParseBlockHashReverted(log types.Log) (*OpenAllianceClientBlockHashReverted, error) {
	event := new(OpenAllianceClientBlockHashReverted)
	if err := _OpenAllianceClient.contract.UnpackLog(event, "BlockHashReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OpenAllianceClientRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleAdminChangedIterator struct {
	Event *OpenAllianceClientRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *OpenAllianceClientRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpenAllianceClientRoleAdminChanged)
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
		it.Event = new(OpenAllianceClientRoleAdminChanged)
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
func (it *OpenAllianceClientRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpenAllianceClientRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpenAllianceClientRoleAdminChanged represents a RoleAdminChanged event raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_OpenAllianceClient *OpenAllianceClientFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*OpenAllianceClientRoleAdminChangedIterator, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientRoleAdminChangedIterator{contract: _OpenAllianceClient.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_OpenAllianceClient *OpenAllianceClientFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *OpenAllianceClientRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpenAllianceClientRoleAdminChanged)
				if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_OpenAllianceClient *OpenAllianceClientFilterer) ParseRoleAdminChanged(log types.Log) (*OpenAllianceClientRoleAdminChanged, error) {
	event := new(OpenAllianceClientRoleAdminChanged)
	if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OpenAllianceClientRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleGrantedIterator struct {
	Event *OpenAllianceClientRoleGranted // Event containing the contract specifics and raw log

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
func (it *OpenAllianceClientRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpenAllianceClientRoleGranted)
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
		it.Event = new(OpenAllianceClientRoleGranted)
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
func (it *OpenAllianceClientRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpenAllianceClientRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpenAllianceClientRoleGranted represents a RoleGranted event raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_OpenAllianceClient *OpenAllianceClientFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*OpenAllianceClientRoleGrantedIterator, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientRoleGrantedIterator{contract: _OpenAllianceClient.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_OpenAllianceClient *OpenAllianceClientFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *OpenAllianceClientRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpenAllianceClientRoleGranted)
				if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_OpenAllianceClient *OpenAllianceClientFilterer) ParseRoleGranted(log types.Log) (*OpenAllianceClientRoleGranted, error) {
	event := new(OpenAllianceClientRoleGranted)
	if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OpenAllianceClientRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleRevokedIterator struct {
	Event *OpenAllianceClientRoleRevoked // Event containing the contract specifics and raw log

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
func (it *OpenAllianceClientRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpenAllianceClientRoleRevoked)
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
		it.Event = new(OpenAllianceClientRoleRevoked)
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
func (it *OpenAllianceClientRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpenAllianceClientRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpenAllianceClientRoleRevoked represents a RoleRevoked event raised by the OpenAllianceClient contract.
type OpenAllianceClientRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_OpenAllianceClient *OpenAllianceClientFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*OpenAllianceClientRoleRevokedIterator, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceClientRoleRevokedIterator{contract: _OpenAllianceClient.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_OpenAllianceClient *OpenAllianceClientFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *OpenAllianceClientRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OpenAllianceClient.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpenAllianceClientRoleRevoked)
				if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_OpenAllianceClient *OpenAllianceClientFilterer) ParseRoleRevoked(log types.Log) (*OpenAllianceClientRoleRevoked, error) {
	event := new(OpenAllianceClientRoleRevoked)
	if err := _OpenAllianceClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
