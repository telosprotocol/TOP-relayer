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
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}

	w := new(Wallet)
	w.chainId = chainid
	w.sdk = sdk

	err = w.initWallet(path, pass)
	if err != nil {
		return nil, err
	}
	return w.createChainWallet(chainid)
}

func (w *Wallet) initWallet(path, pass string) error {
	if path == "" {
		return fmt.Errorf("empty keypath")
	}
	store := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := loadAccount(store, path, pass)
	if err != nil {
		return err
	}
	logger.Info("wallet loads chain[%v] account:%v", w.chainId, acc.Address)
	w.account = acc
	w.provider = store
	err = store.Unlock(w.account, pass)
	if err != nil {
		return err
	}

	return w.verifyChainId()
}

func (w *Wallet) verifyChainId() error {
	chainId, err := w.sdk.ChainID(context.Background())
	if err != nil {
		return err
	}
	if chainId == nil {
		return fmt.Errorf("wallet get chanid nil")
	}

	if w.chainId != chainId.Uint64() {
		return fmt.Errorf("ChainID does not match %v, want: %v", w.chainId, chainId.Uint64())
	}
	return nil
}

func (w *Wallet) NonceAt(address common.Address) (uint64, error) {
	return w.sdk.NonceAt(context.Background(), w.account.Address, nil)
}

func (w *Wallet) BalanceAt(address common.Address) (balance *big.Int, err error) {
	return w.sdk.BalanceAt(context.Background(), address, nil)
}

func (w *Wallet) Address() common.Address {
	return w.account.Address
}

func (w *Wallet) ChainID() *big.Int {
	return big.NewInt(0).SetUint64(w.chainId)
}

func (w *Wallet) SuggestGasPrice() (*big.Int, error) {
	return w.sdk.SuggestGasPrice(context.Background())
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

func (w *Wallet) SuggestGasTipCap() (*big.Int, error) {
	return w.sdk.SuggestGasTipCap(context.Background())
}

//sign tx
func (w *Wallet) SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error) {
	return w.provider.SignTx(w.account, tx, big.NewInt(0).SetUint64(w.chainId))
}

//send signed tx
func (w *Wallet) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return w.sdk.SendTransaction(ctx, tx)
}

/*to create chains wallet...*/
func (w *Wallet) createChainWallet(ID uint64) (IWallet, error) {
	return w, nil
}

func (w *Wallet) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return w.sdk.TransactionReceipt(context.Background(), hash)
}
