package wallet

import (
	"context"
	"fmt"
	"math/big"
	top "toprelayer/util"

	"github.com/wonderivan/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Wallet struct {
	chainId   uint64
	provider  *keystore.KeyStore
	account   accounts.Account // active account
	ethclient *ethclient.Client
	rpc       *rpc.Client
}

func NewTopWallet(topurl, path, pass string) (*Wallet, error) {
	return newWallet(topurl, path, pass)
}

func NewEthWallet(ethurl, topurl, path, pass string) (*Wallet, error) {
	w, err := newWallet(ethurl, path, pass)
	if err != nil {
		return nil, err
	}
	rpcclient, err := rpc.Dial(topurl)
	if err != nil {
		return nil, err
	}
	w.rpc = rpcclient
	return w, nil
}

func newWallet(url, path, pass string) (*Wallet, error) {
	if path == "" {
		return nil, fmt.Errorf("empty keypath")
	}

	w := new(Wallet)

	ethclient, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	w.ethclient = ethclient

	store := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	w.provider = store

	// account
	acc, err := loadAccount(store, path, pass)
	if err != nil {
		return nil, err
	}
	w.account = acc

	// chainId
	id, err := w.ethclient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	w.chainId = id.Uint64()

	logger.Info("wallet loads chain[%v] account:%v", w.chainId, acc.Address)

	// unlock
	err = store.Unlock(w.account, pass)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return w.ethclient.NonceAt(ctx, account, blockNumber)
}

func (w *Wallet) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	return w.ethclient.BalanceAt(ctx, account, nil)
}

func (w *Wallet) Address() common.Address {
	return w.account.Address
}

func (w *Wallet) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0).SetUint64(w.chainId), nil
}

func (w *Wallet) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return w.ethclient.SuggestGasPrice(ctx)
}

func (w *Wallet) EstimateGas(ctx context.Context, target *common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From:      w.Address(),
		To:        target,
		GasPrice:  nil,
		Gas:       0,
		GasFeeCap: nil,
		GasTipCap: nil,
		Data:      data,
	}
	return w.ethclient.EstimateGas(ctx, msg)
}

func (w *Wallet) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return w.ethclient.SuggestGasTipCap(ctx)
}

// sign tx
func (w *Wallet) SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error) {
	return w.provider.SignTx(w.account, tx, big.NewInt(0).SetUint64(w.chainId))
}

// send signed tx
func (w *Wallet) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return w.ethclient.SendTransaction(ctx, tx)
}

func (w *Wallet) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return w.ethclient.TransactionReceipt(ctx, hash)
}

func (w *Wallet) TopBalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	var result hexutil.Big
	err = w.rpc.CallContext(ctx, &result, "top_getBalance", account)
	return (*big.Int)(&result), err
}

func (w *Wallet) TopBlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := w.rpc.CallContext(ctx, &result, "topRelay_blockNumber")
	return uint64(result), err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func (w *Wallet) TopHeaderByNumber(ctx context.Context, number *big.Int) (*top.TopHeader, error) {
	var head *top.TopHeader
	err := w.rpc.CallContext(ctx, &head, "topRelay_getBlockByNumber", toBlockNumArg(number), false, "transaction")
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}
