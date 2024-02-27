// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethclient

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

// EthClientMetaData contains all meta data concerning the EthClient contract.
var EthClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"get_height\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"genesis\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"emitter\",\"type\":\"string\"}],\"name\":\"init\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"data\",\"type\":\"bytes32\"}],\"name\":\"is_known\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"ret\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"blockHeader\",\"type\":\"bytes\"}],\"name\":\"sync\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EthClientABI is the input ABI used to generate the binding from.
// Deprecated: Use EthClientMetaData.ABI instead.
var EthClientABI = EthClientMetaData.ABI

// EthClient is an auto generated Go binding around an Ethereum contract.
type EthClient struct {
	EthClientCaller     // Read-only binding to the contract
	EthClientTransactor // Write-only binding to the contract
	EthClientFilterer   // Log filterer for contract events
}

// EthClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthClientSession struct {
	Contract     *EthClient        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthClientCallerSession struct {
	Contract *EthClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// EthClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthClientTransactorSession struct {
	Contract     *EthClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// EthClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthClientRaw struct {
	Contract *EthClient // Generic contract binding to access the raw methods on
}

// EthClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthClientCallerRaw struct {
	Contract *EthClientCaller // Generic read-only contract binding to access the raw methods on
}

// EthClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthClientTransactorRaw struct {
	Contract *EthClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthClient creates a new instance of EthClient, bound to a specific deployed contract.
func NewEthClient(address common.Address, backend bind.ContractBackend) (*EthClient, error) {
	contract, err := bindEthClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthClient{EthClientCaller: EthClientCaller{contract: contract}, EthClientTransactor: EthClientTransactor{contract: contract}, EthClientFilterer: EthClientFilterer{contract: contract}}, nil
}

// NewEthClientCaller creates a new read-only instance of EthClient, bound to a specific deployed contract.
func NewEthClientCaller(address common.Address, caller bind.ContractCaller) (*EthClientCaller, error) {
	contract, err := bindEthClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthClientCaller{contract: contract}, nil
}

// NewEthClientTransactor creates a new write-only instance of EthClient, bound to a specific deployed contract.
func NewEthClientTransactor(address common.Address, transactor bind.ContractTransactor) (*EthClientTransactor, error) {
	contract, err := bindEthClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthClientTransactor{contract: contract}, nil
}

// NewEthClientFilterer creates a new log filterer instance of EthClient, bound to a specific deployed contract.
func NewEthClientFilterer(address common.Address, filterer bind.ContractFilterer) (*EthClientFilterer, error) {
	contract, err := bindEthClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthClientFilterer{contract: contract}, nil
}

// bindEthClient binds a generic wrapper to an already deployed contract.
func bindEthClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthClient *EthClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthClient.Contract.EthClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthClient *EthClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthClient.Contract.EthClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthClient *EthClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthClient.Contract.EthClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthClient *EthClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthClient *EthClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthClient *EthClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthClient.Contract.contract.Transact(opts, method, params...)
}

// GetHeight is a free data retrieval call binding the contract method 0xb15ad2e8.
//
// Solidity: function get_height() view returns(uint64 height)
func (_EthClient *EthClientCaller) GetHeight(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EthClient.contract.Call(opts, &out, "get_height")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetHeight is a free data retrieval call binding the contract method 0xb15ad2e8.
//
// Solidity: function get_height() view returns(uint64 height)
func (_EthClient *EthClientSession) GetHeight() (uint64, error) {
	return _EthClient.Contract.GetHeight(&_EthClient.CallOpts)
}

// GetHeight is a free data retrieval call binding the contract method 0xb15ad2e8.
//
// Solidity: function get_height() view returns(uint64 height)
func (_EthClient *EthClientCallerSession) GetHeight() (uint64, error) {
	return _EthClient.Contract.GetHeight(&_EthClient.CallOpts)
}

// IsKnown is a free data retrieval call binding the contract method 0x6d571daf.
//
// Solidity: function is_known(uint256 height, bytes32 data) view returns(bool ret)
func (_EthClient *EthClientCaller) IsKnown(opts *bind.CallOpts, height *big.Int, data [32]byte) (bool, error) {
	var out []interface{}
	err := _EthClient.contract.Call(opts, &out, "is_known", height, data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsKnown is a free data retrieval call binding the contract method 0x6d571daf.
//
// Solidity: function is_known(uint256 height, bytes32 data) view returns(bool ret)
func (_EthClient *EthClientSession) IsKnown(height *big.Int, data [32]byte) (bool, error) {
	return _EthClient.Contract.IsKnown(&_EthClient.CallOpts, height, data)
}

// IsKnown is a free data retrieval call binding the contract method 0x6d571daf.
//
// Solidity: function is_known(uint256 height, bytes32 data) view returns(bool ret)
func (_EthClient *EthClientCallerSession) IsKnown(height *big.Int, data [32]byte) (bool, error) {
	return _EthClient.Contract.IsKnown(&_EthClient.CallOpts, height, data)
}

// Init is a paid mutator transaction binding the contract method 0x6158600d.
//
// Solidity: function init(bytes genesis, string emitter) returns(bool success)
func (_EthClient *EthClientTransactor) Init(opts *bind.TransactOpts, genesis []byte, emitter string) (*types.Transaction, error) {
	return _EthClient.contract.Transact(opts, "init", genesis, emitter)
}

// Init is a paid mutator transaction binding the contract method 0x6158600d.
//
// Solidity: function init(bytes genesis, string emitter) returns(bool success)
func (_EthClient *EthClientSession) Init(genesis []byte, emitter string) (*types.Transaction, error) {
	return _EthClient.Contract.Init(&_EthClient.TransactOpts, genesis, emitter)
}

// Init is a paid mutator transaction binding the contract method 0x6158600d.
//
// Solidity: function init(bytes genesis, string emitter) returns(bool success)
func (_EthClient *EthClientTransactorSession) Init(genesis []byte, emitter string) (*types.Transaction, error) {
	return _EthClient.Contract.Init(&_EthClient.TransactOpts, genesis, emitter)
}

// Sync is a paid mutator transaction binding the contract method 0x7eefcfa2.
//
// Solidity: function sync(bytes blockHeader) returns(bool success)
func (_EthClient *EthClientTransactor) Sync(opts *bind.TransactOpts, blockHeader []byte) (*types.Transaction, error) {
	return _EthClient.contract.Transact(opts, "sync", blockHeader)
}

// Sync is a paid mutator transaction binding the contract method 0x7eefcfa2.
//
// Solidity: function sync(bytes blockHeader) returns(bool success)
func (_EthClient *EthClientSession) Sync(blockHeader []byte) (*types.Transaction, error) {
	return _EthClient.Contract.Sync(&_EthClient.TransactOpts, blockHeader)
}

// Sync is a paid mutator transaction binding the contract method 0x7eefcfa2.
//
// Solidity: function sync(bytes blockHeader) returns(bool success)
func (_EthClient *EthClientTransactorSession) Sync(blockHeader []byte) (*types.Transaction, error) {
	return _EthClient.Contract.Sync(&_EthClient.TransactOpts, blockHeader)
}
