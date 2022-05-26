// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bridge

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

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"getBlockBashByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"hashcode\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"}],\"name\":\"getCurrentBlockHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"genesis\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"emitter\",\"type\":\"string\"}],\"name\":\"initGenesisHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"blockHeader\",\"type\":\"bytes\"}],\"name\":\"syncBlockHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridgeLight\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bridgeLight\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Bridge *BridgeCaller) BridgeLight(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "bridgeLight")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Bridge *BridgeSession) BridgeLight() (common.Address, error) {
	return _Bridge.Contract.BridgeLight(&_Bridge.CallOpts)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Bridge *BridgeCallerSession) BridgeLight() (common.Address, error) {
	return _Bridge.Contract.BridgeLight(&_Bridge.CallOpts)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Bridge *BridgeTransactor) GetBlockBashByHeight(opts *bind.TransactOpts, chainId uint64, height uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "getBlockBashByHeight", chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Bridge *BridgeSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _Bridge.Contract.GetBlockBashByHeight(&_Bridge.TransactOpts, chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Bridge *BridgeTransactorSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _Bridge.Contract.GetBlockBashByHeight(&_Bridge.TransactOpts, chainId, height)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Bridge *BridgeTransactor) GetCurrentBlockHeight(opts *bind.TransactOpts, chainId uint64) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "getCurrentBlockHeight", chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Bridge *BridgeSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _Bridge.Contract.GetCurrentBlockHeight(&_Bridge.TransactOpts, chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Bridge *BridgeTransactorSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _Bridge.Contract.GetCurrentBlockHeight(&_Bridge.TransactOpts, chainId)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Bridge *BridgeTransactor) InitGenesisHeader(opts *bind.TransactOpts, genesis []byte, emitter string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "initGenesisHeader", genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Bridge *BridgeSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _Bridge.Contract.InitGenesisHeader(&_Bridge.TransactOpts, genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Bridge *BridgeTransactorSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _Bridge.Contract.InitGenesisHeader(&_Bridge.TransactOpts, genesis, emitter)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Bridge *BridgeTransactor) SyncBlockHeader(opts *bind.TransactOpts, blockHeader []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "syncBlockHeader", blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Bridge *BridgeSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _Bridge.Contract.SyncBlockHeader(&_Bridge.TransactOpts, blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Bridge *BridgeTransactorSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _Bridge.Contract.SyncBlockHeader(&_Bridge.TransactOpts, blockHeader)
}
