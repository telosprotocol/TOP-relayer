// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package topbridge

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

// TopBridgeMetaData contains all meta data concerning the TopBridge contract.
var TopBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"getBlockBashByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"hashcode\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"}],\"name\":\"getCurrentBlockHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"genesis\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"emitter\",\"type\":\"string\"}],\"name\":\"initGenesisHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"blockHeader\",\"type\":\"bytes\"}],\"name\":\"syncBlockHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridgeLight\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bridgeLight\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TopBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use TopBridgeMetaData.ABI instead.
var TopBridgeABI = TopBridgeMetaData.ABI

// TopBridge is an auto generated Go binding around an Ethereum contract.
type TopBridge struct {
	TopBridgeCaller     // Read-only binding to the contract
	TopBridgeTransactor // Write-only binding to the contract
	TopBridgeFilterer   // Log filterer for contract events
}

// TopBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TopBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TopBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TopBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TopBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TopBridgeSession struct {
	Contract     *TopBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TopBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TopBridgeCallerSession struct {
	Contract *TopBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TopBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TopBridgeTransactorSession struct {
	Contract     *TopBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TopBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TopBridgeRaw struct {
	Contract *TopBridge // Generic contract binding to access the raw methods on
}

// TopBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TopBridgeCallerRaw struct {
	Contract *TopBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// TopBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TopBridgeTransactorRaw struct {
	Contract *TopBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTopBridge creates a new instance of TopBridge, bound to a specific deployed contract.
func NewTopBridge(address common.Address, backend bind.ContractBackend) (*TopBridge, error) {
	contract, err := bindTopBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TopBridge{TopBridgeCaller: TopBridgeCaller{contract: contract}, TopBridgeTransactor: TopBridgeTransactor{contract: contract}, TopBridgeFilterer: TopBridgeFilterer{contract: contract}}, nil
}

// NewTopBridgeCaller creates a new read-only instance of TopBridge, bound to a specific deployed contract.
func NewTopBridgeCaller(address common.Address, caller bind.ContractCaller) (*TopBridgeCaller, error) {
	contract, err := bindTopBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TopBridgeCaller{contract: contract}, nil
}

// NewTopBridgeTransactor creates a new write-only instance of TopBridge, bound to a specific deployed contract.
func NewTopBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*TopBridgeTransactor, error) {
	contract, err := bindTopBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TopBridgeTransactor{contract: contract}, nil
}

// NewTopBridgeFilterer creates a new log filterer instance of TopBridge, bound to a specific deployed contract.
func NewTopBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*TopBridgeFilterer, error) {
	contract, err := bindTopBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TopBridgeFilterer{contract: contract}, nil
}

// bindTopBridge binds a generic wrapper to an already deployed contract.
func bindTopBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TopBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TopBridge *TopBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TopBridge.Contract.TopBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TopBridge *TopBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TopBridge.Contract.TopBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TopBridge *TopBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TopBridge.Contract.TopBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TopBridge *TopBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TopBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TopBridge *TopBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TopBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TopBridge *TopBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TopBridge.Contract.contract.Transact(opts, method, params...)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_TopBridge *TopBridgeCaller) BridgeLight(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TopBridge.contract.Call(opts, &out, "bridgeLight")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_TopBridge *TopBridgeSession) BridgeLight() (common.Address, error) {
	return _TopBridge.Contract.BridgeLight(&_TopBridge.CallOpts)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_TopBridge *TopBridgeCallerSession) BridgeLight() (common.Address, error) {
	return _TopBridge.Contract.BridgeLight(&_TopBridge.CallOpts)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_TopBridge *TopBridgeTransactor) GetBlockBashByHeight(opts *bind.TransactOpts, chainId uint64, height uint64) (*types.Transaction, error) {
	return _TopBridge.contract.Transact(opts, "getBlockBashByHeight", chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_TopBridge *TopBridgeSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _TopBridge.Contract.GetBlockBashByHeight(&_TopBridge.TransactOpts, chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_TopBridge *TopBridgeTransactorSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _TopBridge.Contract.GetBlockBashByHeight(&_TopBridge.TransactOpts, chainId, height)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_TopBridge *TopBridgeTransactor) GetCurrentBlockHeight(opts *bind.TransactOpts, chainId uint64) (*types.Transaction, error) {
	return _TopBridge.contract.Transact(opts, "getCurrentBlockHeight", chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_TopBridge *TopBridgeSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _TopBridge.Contract.GetCurrentBlockHeight(&_TopBridge.TransactOpts, chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_TopBridge *TopBridgeTransactorSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _TopBridge.Contract.GetCurrentBlockHeight(&_TopBridge.TransactOpts, chainId)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_TopBridge *TopBridgeTransactor) InitGenesisHeader(opts *bind.TransactOpts, genesis []byte, emitter string) (*types.Transaction, error) {
	return _TopBridge.contract.Transact(opts, "initGenesisHeader", genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_TopBridge *TopBridgeSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _TopBridge.Contract.InitGenesisHeader(&_TopBridge.TransactOpts, genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_TopBridge *TopBridgeTransactorSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _TopBridge.Contract.InitGenesisHeader(&_TopBridge.TransactOpts, genesis, emitter)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_TopBridge *TopBridgeTransactor) SyncBlockHeader(opts *bind.TransactOpts, blockHeader []byte) (*types.Transaction, error) {
	return _TopBridge.contract.Transact(opts, "syncBlockHeader", blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_TopBridge *TopBridgeSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _TopBridge.Contract.SyncBlockHeader(&_TopBridge.TransactOpts, blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_TopBridge *TopBridgeTransactorSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _TopBridge.Contract.SyncBlockHeader(&_TopBridge.TransactOpts, blockHeader)
}
