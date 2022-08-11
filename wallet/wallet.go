package wallet

import (
	"context"
	"fmt"
	"math/big"
	"toprelayer/sdk"

	"github.com/wonderivan/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewWallet(url, path, pass string, chainid uint64) (IWallet, error) {
	if path == "" {
		return nil, fmt.Errorf("empty keypath")
	}

	w := new(Wallet)
	w.chainId = chainid

	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}
	w.sdk = sdk

	store := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	w.provider = store

	acc, err := loadAccount(store, path, pass)
	if err != nil {
		return nil, err
	}
	w.account = acc

	logger.Info("wallet loads chain[%v] account:%v", w.chainId, acc.Address)

	// unlock
	err = store.Unlock(w.account, pass)
	if err != nil {
		return nil, err
	}
	// verify chainId
	id, err := w.sdk.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	if w.chainId != id.Uint64() {
		return nil, fmt.Errorf("ChainID does not match %v, want: %v", w.chainId, id.Uint64())
	}

	return w, nil
}

func (w *Wallet) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return w.sdk.NonceAt(ctx, account, blockNumber)
}

func (w *Wallet) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	return w.sdk.BalanceAt(ctx, account, nil)
}

func (w *Wallet) Address() common.Address {
	return w.account.Address
}

func (w *Wallet) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0).SetUint64(w.chainId), nil
}

func (w *Wallet) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return w.sdk.SuggestGasPrice(ctx)
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
	return w.sdk.EstimateGas(ctx, msg)
}

func (w *Wallet) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return w.sdk.SuggestGasTipCap(ctx)
}

//sign tx
func (w *Wallet) SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error) {
	return w.provider.SignTx(w.account, tx, big.NewInt(0).SetUint64(w.chainId))
}

//send signed tx
func (w *Wallet) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return w.sdk.SendTransaction(ctx, tx)
}

func (w *Wallet) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return w.sdk.TransactionReceipt(ctx, hash)
}
