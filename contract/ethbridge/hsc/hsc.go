// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hsc

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

// HscMetaData contains all meta data concerning the Hsc contract.
var HscMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridgeLight\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bridgeLight\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"name\":\"getBlockBashByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"hashcode\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"}],\"name\":\"getCurrentBlockHeight\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"genesis\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"emitter\",\"type\":\"string\"}],\"name\":\"initGenesisHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"blockHeader\",\"type\":\"bytes\"}],\"name\":\"syncBlockHeader\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610def380380610def833981810160405281019061003291906100db565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610108565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100a88261007d565b9050919050565b6100b88161009d565b81146100c357600080fd5b50565b6000815190506100d5816100af565b92915050565b6000602082840312156100f1576100f0610078565b5b60006100ff848285016100c6565b91505092915050565b610cd8806101176000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806301402e0c1461005c5780631e0906261461007a578063abd70ea6146100aa578063ccadba86146100da578063e9a9c4cd1461010a575b600080fd5b61006461013a565b604051610071919061069f565b60405180910390f35b610094600480360381019061008f9190610814565b61015e565b6040516100a19190610878565b60405180910390f35b6100c460048036038101906100bf91906108d3565b61028e565b6040516100d1919061090f565b60405180910390f35b6100f460048036038101906100ef91906109cb565b6103f2565b6040516101019190610878565b60405180910390f35b610124600480360381019061011f9190610a43565b610525565b6040516101319190610b0b565b60405180910390f35b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080826040516024016101729190610b0b565b6040516020818303038152906040527f70243237000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16816040516102369190610b69565b6000604051808303816000865af19150503d8060008114610273576040519150601f19603f3d011682016040523d82523d6000602084013e610278565b606091505b5050809250508161028857600080fd5b50919050565b600080826040516024016102a2919061090f565b6040516020818303038152906040527fdaf0c99a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168360405161036a9190610b69565b6000604051808303816000865af19150503d80600081146103a7576040519150601f19603f3d011682016040523d82523d6000602084013e6103ac565b606091505b5080925081935050506020846103c29190610baf565b67ffffffffffffffff16815110156103d957600080fd5b60208101519350816103ea57600080fd5b505050919050565b6000808383604051602401610408929190610c42565b6040516020818303038152906040527f19bc024a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16816040516104cc9190610b69565b6000604051808303816000865af19150503d8060008114610509576040519150601f19603f3d011682016040523d82523d6000602084013e61050e565b606091505b5050809250508161051e57600080fd5b5092915050565b60606000838360405160240161053c929190610c79565b6040516020818303038152906040527fbec5610f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000808054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16826040516106019190610b69565b6000604051808303816000865af19150503d806000811461063e576040519150601f19603f3d011682016040523d82523d6000602084013e610643565b606091505b5080945081925050508061065657600080fd5b505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106898261065e565b9050919050565b6106998161067e565b82525050565b60006020820190506106b46000830184610690565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610721826106d8565b810181811067ffffffffffffffff821117156107405761073f6106e9565b5b80604052505050565b60006107536106ba565b905061075f8282610718565b919050565b600067ffffffffffffffff82111561077f5761077e6106e9565b5b610788826106d8565b9050602081019050919050565b82818337600083830152505050565b60006107b76107b284610764565b610749565b9050828152602081018484840111156107d3576107d26106d3565b5b6107de848285610795565b509392505050565b600082601f8301126107fb576107fa6106ce565b5b813561080b8482602086016107a4565b91505092915050565b60006020828403121561082a576108296106c4565b5b600082013567ffffffffffffffff811115610848576108476106c9565b5b610854848285016107e6565b91505092915050565b60008115159050919050565b6108728161085d565b82525050565b600060208201905061088d6000830184610869565b92915050565b600067ffffffffffffffff82169050919050565b6108b081610893565b81146108bb57600080fd5b50565b6000813590506108cd816108a7565b92915050565b6000602082840312156108e9576108e86106c4565b5b60006108f7848285016108be565b91505092915050565b61090981610893565b82525050565b60006020820190506109246000830184610900565b92915050565b600067ffffffffffffffff821115610945576109446106e9565b5b61094e826106d8565b9050602081019050919050565b600061096e6109698461092a565b610749565b90508281526020810184848401111561098a576109896106d3565b5b610995848285610795565b509392505050565b600082601f8301126109b2576109b16106ce565b5b81356109c284826020860161095b565b91505092915050565b600080604083850312156109e2576109e16106c4565b5b600083013567ffffffffffffffff811115610a00576109ff6106c9565b5b610a0c858286016107e6565b925050602083013567ffffffffffffffff811115610a2d57610a2c6106c9565b5b610a398582860161099d565b9150509250929050565b60008060408385031215610a5a57610a596106c4565b5b6000610a68858286016108be565b9250506020610a79858286016108be565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610abd578082015181840152602081019050610aa2565b83811115610acc576000848401525b50505050565b6000610add82610a83565b610ae78185610a8e565b9350610af7818560208601610a9f565b610b00816106d8565b840191505092915050565b60006020820190508181036000830152610b258184610ad2565b905092915050565b600081905092915050565b6000610b4382610a83565b610b4d8185610b2d565b9350610b5d818560208601610a9f565b80840191505092915050565b6000610b758284610b38565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610bba82610893565b9150610bc583610893565b92508267ffffffffffffffff03821115610be257610be1610b80565b5b828201905092915050565b600081519050919050565b600082825260208201905092915050565b6000610c1482610bed565b610c1e8185610bf8565b9350610c2e818560208601610a9f565b610c37816106d8565b840191505092915050565b60006040820190508181036000830152610c5c8185610ad2565b90508181036020830152610c708184610c09565b90509392505050565b6000604082019050610c8e6000830185610900565b610c9b6020830184610900565b939250505056fea26469706673582212203a20bfd8d24930411e86982486f075fa25c89c138ce4c351fa353fce8b46fd1d64736f6c634300080d0033",
}

// HscABI is the input ABI used to generate the binding from.
// Deprecated: Use HscMetaData.ABI instead.
var HscABI = HscMetaData.ABI

// HscBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HscMetaData.Bin instead.
var HscBin = HscMetaData.Bin

// DeployHsc deploys a new Ethereum contract, binding an instance of Hsc to it.
func DeployHsc(auth *bind.TransactOpts, backend bind.ContractBackend, _bridgeLight common.Address) (common.Address, *types.Transaction, *Hsc, error) {
	parsed, err := HscMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HscBin), backend, _bridgeLight)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Hsc{HscCaller: HscCaller{contract: contract}, HscTransactor: HscTransactor{contract: contract}, HscFilterer: HscFilterer{contract: contract}}, nil
}

// Hsc is an auto generated Go binding around an Ethereum contract.
type Hsc struct {
	HscCaller     // Read-only binding to the contract
	HscTransactor // Write-only binding to the contract
	HscFilterer   // Log filterer for contract events
}

// HscCaller is an auto generated read-only Go binding around an Ethereum contract.
type HscCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HscTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HscTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HscFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HscFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HscSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HscSession struct {
	Contract     *Hsc              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HscCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HscCallerSession struct {
	Contract *HscCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HscTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HscTransactorSession struct {
	Contract     *HscTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HscRaw is an auto generated low-level Go binding around an Ethereum contract.
type HscRaw struct {
	Contract *Hsc // Generic contract binding to access the raw methods on
}

// HscCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HscCallerRaw struct {
	Contract *HscCaller // Generic read-only contract binding to access the raw methods on
}

// HscTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HscTransactorRaw struct {
	Contract *HscTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHsc creates a new instance of Hsc, bound to a specific deployed contract.
func NewHsc(address common.Address, backend bind.ContractBackend) (*Hsc, error) {
	contract, err := bindHsc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Hsc{HscCaller: HscCaller{contract: contract}, HscTransactor: HscTransactor{contract: contract}, HscFilterer: HscFilterer{contract: contract}}, nil
}

// NewHscCaller creates a new read-only instance of Hsc, bound to a specific deployed contract.
func NewHscCaller(address common.Address, caller bind.ContractCaller) (*HscCaller, error) {
	contract, err := bindHsc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HscCaller{contract: contract}, nil
}

// NewHscTransactor creates a new write-only instance of Hsc, bound to a specific deployed contract.
func NewHscTransactor(address common.Address, transactor bind.ContractTransactor) (*HscTransactor, error) {
	contract, err := bindHsc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HscTransactor{contract: contract}, nil
}

// NewHscFilterer creates a new log filterer instance of Hsc, bound to a specific deployed contract.
func NewHscFilterer(address common.Address, filterer bind.ContractFilterer) (*HscFilterer, error) {
	contract, err := bindHsc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HscFilterer{contract: contract}, nil
}

// bindHsc binds a generic wrapper to an already deployed contract.
func bindHsc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HscABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hsc *HscRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hsc.Contract.HscCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hsc *HscRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hsc.Contract.HscTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hsc *HscRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hsc.Contract.HscTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Hsc *HscCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Hsc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Hsc *HscTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Hsc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Hsc *HscTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Hsc.Contract.contract.Transact(opts, method, params...)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Hsc *HscCaller) BridgeLight(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Hsc.contract.Call(opts, &out, "bridgeLight")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Hsc *HscSession) BridgeLight() (common.Address, error) {
	return _Hsc.Contract.BridgeLight(&_Hsc.CallOpts)
}

// BridgeLight is a free data retrieval call binding the contract method 0x01402e0c.
//
// Solidity: function bridgeLight() view returns(address)
func (_Hsc *HscCallerSession) BridgeLight() (common.Address, error) {
	return _Hsc.Contract.BridgeLight(&_Hsc.CallOpts)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Hsc *HscTransactor) GetBlockBashByHeight(opts *bind.TransactOpts, chainId uint64, height uint64) (*types.Transaction, error) {
	return _Hsc.contract.Transact(opts, "getBlockBashByHeight", chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Hsc *HscSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _Hsc.Contract.GetBlockBashByHeight(&_Hsc.TransactOpts, chainId, height)
}

// GetBlockBashByHeight is a paid mutator transaction binding the contract method 0xe9a9c4cd.
//
// Solidity: function getBlockBashByHeight(uint64 chainId, uint64 height) returns(bytes hashcode)
func (_Hsc *HscTransactorSession) GetBlockBashByHeight(chainId uint64, height uint64) (*types.Transaction, error) {
	return _Hsc.Contract.GetBlockBashByHeight(&_Hsc.TransactOpts, chainId, height)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Hsc *HscTransactor) GetCurrentBlockHeight(opts *bind.TransactOpts, chainId uint64) (*types.Transaction, error) {
	return _Hsc.contract.Transact(opts, "getCurrentBlockHeight", chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Hsc *HscSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _Hsc.Contract.GetCurrentBlockHeight(&_Hsc.TransactOpts, chainId)
}

// GetCurrentBlockHeight is a paid mutator transaction binding the contract method 0xabd70ea6.
//
// Solidity: function getCurrentBlockHeight(uint64 chainId) returns(uint64 height)
func (_Hsc *HscTransactorSession) GetCurrentBlockHeight(chainId uint64) (*types.Transaction, error) {
	return _Hsc.Contract.GetCurrentBlockHeight(&_Hsc.TransactOpts, chainId)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Hsc *HscTransactor) InitGenesisHeader(opts *bind.TransactOpts, genesis []byte, emitter string) (*types.Transaction, error) {
	return _Hsc.contract.Transact(opts, "initGenesisHeader", genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Hsc *HscSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _Hsc.Contract.InitGenesisHeader(&_Hsc.TransactOpts, genesis, emitter)
}

// InitGenesisHeader is a paid mutator transaction binding the contract method 0xccadba86.
//
// Solidity: function initGenesisHeader(bytes genesis, string emitter) returns(bool success)
func (_Hsc *HscTransactorSession) InitGenesisHeader(genesis []byte, emitter string) (*types.Transaction, error) {
	return _Hsc.Contract.InitGenesisHeader(&_Hsc.TransactOpts, genesis, emitter)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Hsc *HscTransactor) SyncBlockHeader(opts *bind.TransactOpts, blockHeader []byte) (*types.Transaction, error) {
	return _Hsc.contract.Transact(opts, "syncBlockHeader", blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Hsc *HscSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _Hsc.Contract.SyncBlockHeader(&_Hsc.TransactOpts, blockHeader)
}

// SyncBlockHeader is a paid mutator transaction binding the contract method 0x1e090626.
//
// Solidity: function syncBlockHeader(bytes blockHeader) returns(bool success)
func (_Hsc *HscTransactorSession) SyncBlockHeader(blockHeader []byte) (*types.Transaction, error) {
	return _Hsc.Contract.SyncBlockHeader(&_Hsc.TransactOpts, blockHeader)
}
