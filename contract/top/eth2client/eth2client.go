// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth2client

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

// Eth2ClientMetaData contains all meta data concerning the Eth2Client contract.
var Eth2ClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"block_hash_safe\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalized_beacon_block_header\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalized_beacon_block_root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalized_beacon_block_slot\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"slot\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_light_client_state\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"state\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"genesis\",\"type\":\"bytes\"}],\"name\":\"init\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"inited\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"is_known_execution_header\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"known\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_block_number\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"number\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"name\":\"submit_beacon_chain_light_client_update\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"submit_execution_header\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Eth2ClientABI is the input ABI used to generate the binding from.
// Deprecated: Use Eth2ClientMetaData.ABI instead.
var Eth2ClientABI = Eth2ClientMetaData.ABI

// Eth2Client is an auto generated Go binding around an Ethereum contract.
type Eth2Client struct {
	Eth2ClientCaller     // Read-only binding to the contract
	Eth2ClientTransactor // Write-only binding to the contract
	Eth2ClientFilterer   // Log filterer for contract events
}

// Eth2ClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type Eth2ClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Eth2ClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Eth2ClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Eth2ClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Eth2ClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Eth2ClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Eth2ClientSession struct {
	Contract     *Eth2Client       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Eth2ClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Eth2ClientCallerSession struct {
	Contract *Eth2ClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// Eth2ClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Eth2ClientTransactorSession struct {
	Contract     *Eth2ClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// Eth2ClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type Eth2ClientRaw struct {
	Contract *Eth2Client // Generic contract binding to access the raw methods on
}

// Eth2ClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Eth2ClientCallerRaw struct {
	Contract *Eth2ClientCaller // Generic read-only contract binding to access the raw methods on
}

// Eth2ClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Eth2ClientTransactorRaw struct {
	Contract *Eth2ClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEth2Client creates a new instance of Eth2Client, bound to a specific deployed contract.
func NewEth2Client(address common.Address, backend bind.ContractBackend) (*Eth2Client, error) {
	contract, err := bindEth2Client(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Eth2Client{Eth2ClientCaller: Eth2ClientCaller{contract: contract}, Eth2ClientTransactor: Eth2ClientTransactor{contract: contract}, Eth2ClientFilterer: Eth2ClientFilterer{contract: contract}}, nil
}

// NewEth2ClientCaller creates a new read-only instance of Eth2Client, bound to a specific deployed contract.
func NewEth2ClientCaller(address common.Address, caller bind.ContractCaller) (*Eth2ClientCaller, error) {
	contract, err := bindEth2Client(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Eth2ClientCaller{contract: contract}, nil
}

// NewEth2ClientTransactor creates a new write-only instance of Eth2Client, bound to a specific deployed contract.
func NewEth2ClientTransactor(address common.Address, transactor bind.ContractTransactor) (*Eth2ClientTransactor, error) {
	contract, err := bindEth2Client(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Eth2ClientTransactor{contract: contract}, nil
}

// NewEth2ClientFilterer creates a new log filterer instance of Eth2Client, bound to a specific deployed contract.
func NewEth2ClientFilterer(address common.Address, filterer bind.ContractFilterer) (*Eth2ClientFilterer, error) {
	contract, err := bindEth2Client(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Eth2ClientFilterer{contract: contract}, nil
}

// bindEth2Client binds a generic wrapper to an already deployed contract.
func bindEth2Client(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Eth2ClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth2Client *Eth2ClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth2Client.Contract.Eth2ClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth2Client *Eth2ClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth2Client.Contract.Eth2ClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth2Client *Eth2ClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth2Client.Contract.Eth2ClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth2Client *Eth2ClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth2Client.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth2Client *Eth2ClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth2Client.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth2Client *Eth2ClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth2Client.Contract.contract.Transact(opts, method, params...)
}

// BlockHashSafe is a free data retrieval call binding the contract method 0x3bcdaaab.
//
// Solidity: function block_hash_safe(uint64 height) view returns(bytes32 hash)
func (_Eth2Client *Eth2ClientCaller) BlockHashSafe(opts *bind.CallOpts, height uint64) ([32]byte, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "block_hash_safe", height)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockHashSafe is a free data retrieval call binding the contract method 0x3bcdaaab.
//
// Solidity: function block_hash_safe(uint64 height) view returns(bytes32 hash)
func (_Eth2Client *Eth2ClientSession) BlockHashSafe(height uint64) ([32]byte, error) {
	return _Eth2Client.Contract.BlockHashSafe(&_Eth2Client.CallOpts, height)
}

// BlockHashSafe is a free data retrieval call binding the contract method 0x3bcdaaab.
//
// Solidity: function block_hash_safe(uint64 height) view returns(bytes32 hash)
func (_Eth2Client *Eth2ClientCallerSession) BlockHashSafe(height uint64) ([32]byte, error) {
	return _Eth2Client.Contract.BlockHashSafe(&_Eth2Client.CallOpts, height)
}

// FinalizedBeaconBlockHeader is a free data retrieval call binding the contract method 0x55b39f6e.
//
// Solidity: function finalized_beacon_block_header() view returns(bytes header)
func (_Eth2Client *Eth2ClientCaller) FinalizedBeaconBlockHeader(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "finalized_beacon_block_header")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// FinalizedBeaconBlockHeader is a free data retrieval call binding the contract method 0x55b39f6e.
//
// Solidity: function finalized_beacon_block_header() view returns(bytes header)
func (_Eth2Client *Eth2ClientSession) FinalizedBeaconBlockHeader() ([]byte, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockHeader(&_Eth2Client.CallOpts)
}

// FinalizedBeaconBlockHeader is a free data retrieval call binding the contract method 0x55b39f6e.
//
// Solidity: function finalized_beacon_block_header() view returns(bytes header)
func (_Eth2Client *Eth2ClientCallerSession) FinalizedBeaconBlockHeader() ([]byte, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockHeader(&_Eth2Client.CallOpts)
}

// FinalizedBeaconBlockRoot is a free data retrieval call binding the contract method 0x4b469132.
//
// Solidity: function finalized_beacon_block_root() view returns(bytes32 root)
func (_Eth2Client *Eth2ClientCaller) FinalizedBeaconBlockRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "finalized_beacon_block_root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FinalizedBeaconBlockRoot is a free data retrieval call binding the contract method 0x4b469132.
//
// Solidity: function finalized_beacon_block_root() view returns(bytes32 root)
func (_Eth2Client *Eth2ClientSession) FinalizedBeaconBlockRoot() ([32]byte, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockRoot(&_Eth2Client.CallOpts)
}

// FinalizedBeaconBlockRoot is a free data retrieval call binding the contract method 0x4b469132.
//
// Solidity: function finalized_beacon_block_root() view returns(bytes32 root)
func (_Eth2Client *Eth2ClientCallerSession) FinalizedBeaconBlockRoot() ([32]byte, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockRoot(&_Eth2Client.CallOpts)
}

// FinalizedBeaconBlockSlot is a free data retrieval call binding the contract method 0x074b1681.
//
// Solidity: function finalized_beacon_block_slot() view returns(uint64 slot)
func (_Eth2Client *Eth2ClientCaller) FinalizedBeaconBlockSlot(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "finalized_beacon_block_slot")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// FinalizedBeaconBlockSlot is a free data retrieval call binding the contract method 0x074b1681.
//
// Solidity: function finalized_beacon_block_slot() view returns(uint64 slot)
func (_Eth2Client *Eth2ClientSession) FinalizedBeaconBlockSlot() (uint64, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockSlot(&_Eth2Client.CallOpts)
}

// FinalizedBeaconBlockSlot is a free data retrieval call binding the contract method 0x074b1681.
//
// Solidity: function finalized_beacon_block_slot() view returns(uint64 slot)
func (_Eth2Client *Eth2ClientCallerSession) FinalizedBeaconBlockSlot() (uint64, error) {
	return _Eth2Client.Contract.FinalizedBeaconBlockSlot(&_Eth2Client.CallOpts)
}

// GetLightClientState is a free data retrieval call binding the contract method 0x3ae8d743.
//
// Solidity: function get_light_client_state() view returns(bytes state)
func (_Eth2Client *Eth2ClientCaller) GetLightClientState(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "get_light_client_state")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetLightClientState is a free data retrieval call binding the contract method 0x3ae8d743.
//
// Solidity: function get_light_client_state() view returns(bytes state)
func (_Eth2Client *Eth2ClientSession) GetLightClientState() ([]byte, error) {
	return _Eth2Client.Contract.GetLightClientState(&_Eth2Client.CallOpts)
}

// GetLightClientState is a free data retrieval call binding the contract method 0x3ae8d743.
//
// Solidity: function get_light_client_state() view returns(bytes state)
func (_Eth2Client *Eth2ClientCallerSession) GetLightClientState() ([]byte, error) {
	return _Eth2Client.Contract.GetLightClientState(&_Eth2Client.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool inited)
func (_Eth2Client *Eth2ClientCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool inited)
func (_Eth2Client *Eth2ClientSession) Initialized() (bool, error) {
	return _Eth2Client.Contract.Initialized(&_Eth2Client.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool inited)
func (_Eth2Client *Eth2ClientCallerSession) Initialized() (bool, error) {
	return _Eth2Client.Contract.Initialized(&_Eth2Client.CallOpts)
}

// IsKnownExecutionHeader is a free data retrieval call binding the contract method 0x43b1378b.
//
// Solidity: function is_known_execution_header(bytes32 data) view returns(bool known)
func (_Eth2Client *Eth2ClientCaller) IsKnownExecutionHeader(opts *bind.CallOpts, data [32]byte) (bool, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "is_known_execution_header", data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsKnownExecutionHeader is a free data retrieval call binding the contract method 0x43b1378b.
//
// Solidity: function is_known_execution_header(bytes32 data) view returns(bool known)
func (_Eth2Client *Eth2ClientSession) IsKnownExecutionHeader(data [32]byte) (bool, error) {
	return _Eth2Client.Contract.IsKnownExecutionHeader(&_Eth2Client.CallOpts, data)
}

// IsKnownExecutionHeader is a free data retrieval call binding the contract method 0x43b1378b.
//
// Solidity: function is_known_execution_header(bytes32 data) view returns(bool known)
func (_Eth2Client *Eth2ClientCallerSession) IsKnownExecutionHeader(data [32]byte) (bool, error) {
	return _Eth2Client.Contract.IsKnownExecutionHeader(&_Eth2Client.CallOpts, data)
}

// LastBlockNumber is a free data retrieval call binding the contract method 0x1eeaebb2.
//
// Solidity: function last_block_number() view returns(uint64 number)
func (_Eth2Client *Eth2ClientCaller) LastBlockNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Eth2Client.contract.Call(opts, &out, "last_block_number")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastBlockNumber is a free data retrieval call binding the contract method 0x1eeaebb2.
//
// Solidity: function last_block_number() view returns(uint64 number)
func (_Eth2Client *Eth2ClientSession) LastBlockNumber() (uint64, error) {
	return _Eth2Client.Contract.LastBlockNumber(&_Eth2Client.CallOpts)
}

// LastBlockNumber is a free data retrieval call binding the contract method 0x1eeaebb2.
//
// Solidity: function last_block_number() view returns(uint64 number)
func (_Eth2Client *Eth2ClientCallerSession) LastBlockNumber() (uint64, error) {
	return _Eth2Client.Contract.LastBlockNumber(&_Eth2Client.CallOpts)
}

// Init is a paid mutator transaction binding the contract method 0x4ddf47d4.
//
// Solidity: function init(bytes genesis) returns(bool success)
func (_Eth2Client *Eth2ClientTransactor) Init(opts *bind.TransactOpts, genesis []byte) (*types.Transaction, error) {
	return _Eth2Client.contract.Transact(opts, "init", genesis)
}

// Init is a paid mutator transaction binding the contract method 0x4ddf47d4.
//
// Solidity: function init(bytes genesis) returns(bool success)
func (_Eth2Client *Eth2ClientSession) Init(genesis []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.Init(&_Eth2Client.TransactOpts, genesis)
}

// Init is a paid mutator transaction binding the contract method 0x4ddf47d4.
//
// Solidity: function init(bytes genesis) returns(bool success)
func (_Eth2Client *Eth2ClientTransactorSession) Init(genesis []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.Init(&_Eth2Client.TransactOpts, genesis)
}

// SubmitBeaconChainLightClientUpdate is a paid mutator transaction binding the contract method 0x2e139f0c.
//
// Solidity: function submit_beacon_chain_light_client_update(bytes update) returns(bool success)
func (_Eth2Client *Eth2ClientTransactor) SubmitBeaconChainLightClientUpdate(opts *bind.TransactOpts, update []byte) (*types.Transaction, error) {
	return _Eth2Client.contract.Transact(opts, "submit_beacon_chain_light_client_update", update)
}

// SubmitBeaconChainLightClientUpdate is a paid mutator transaction binding the contract method 0x2e139f0c.
//
// Solidity: function submit_beacon_chain_light_client_update(bytes update) returns(bool success)
func (_Eth2Client *Eth2ClientSession) SubmitBeaconChainLightClientUpdate(update []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.SubmitBeaconChainLightClientUpdate(&_Eth2Client.TransactOpts, update)
}

// SubmitBeaconChainLightClientUpdate is a paid mutator transaction binding the contract method 0x2e139f0c.
//
// Solidity: function submit_beacon_chain_light_client_update(bytes update) returns(bool success)
func (_Eth2Client *Eth2ClientTransactorSession) SubmitBeaconChainLightClientUpdate(update []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.SubmitBeaconChainLightClientUpdate(&_Eth2Client.TransactOpts, update)
}

// SubmitExecutionHeader is a paid mutator transaction binding the contract method 0x3c1a38b6.
//
// Solidity: function submit_execution_header(bytes header) returns(bool success)
func (_Eth2Client *Eth2ClientTransactor) SubmitExecutionHeader(opts *bind.TransactOpts, header []byte) (*types.Transaction, error) {
	return _Eth2Client.contract.Transact(opts, "submit_execution_header", header)
}

// SubmitExecutionHeader is a paid mutator transaction binding the contract method 0x3c1a38b6.
//
// Solidity: function submit_execution_header(bytes header) returns(bool success)
func (_Eth2Client *Eth2ClientSession) SubmitExecutionHeader(header []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.SubmitExecutionHeader(&_Eth2Client.TransactOpts, header)
}

// SubmitExecutionHeader is a paid mutator transaction binding the contract method 0x3c1a38b6.
//
// Solidity: function submit_execution_header(bytes header) returns(bool success)
func (_Eth2Client *Eth2ClientTransactorSession) SubmitExecutionHeader(header []byte) (*types.Transaction, error) {
	return _Eth2Client.Contract.SubmitExecutionHeader(&_Eth2Client.TransactOpts, header)
}