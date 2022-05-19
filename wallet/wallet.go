package wallet

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"toprelayer/base"
	"toprelayer/sdk"

	"github.com/wonderivan/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
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

func newkeystore(path, pass string) (*keystore.KeyStore, error) {
	store := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	_, err := os.Stat(path)
	if err != nil {
		_, err := createAccount(store, pass)
		if err != nil {
			return nil, err
		}
	}
	return store, nil
}
func (w *Wallet) initWallet(path, pass string) error {
	store, err := newkeystore(path, pass)
	if err != nil {
		return err
	}

	acc, err := loadAccount(store, path, pass)
	if err != nil {
		return err
	}
	logger.Info("wallet loads chain[%v] account:%v", w.chainId, acc.Address)
	w.account = acc

	p := newKeyStoreProvider(store, pass)
	if p == nil {
		return fmt.Errorf("keystore provider is nil")
	}
	w.provider = p
	err = p.UnlockAccount(w.account)
	if err != nil {
		return err
	}

	if w.sdk == nil {
		return fmt.Errorf("fatal error: nil sdk!")
	}
	nc, err := w.getNonceAt()
	if err != nil {
		fmt.Println("check err2:", err)
		return err
	}
	w.nonce = nc
	return w.verifyChainId()
}
func (w *Wallet) getNonceAt() (uint64, error) {
	return w.sdk.NonceAt(context.Background(), w.account.Address, nil)
}
func (w *Wallet) verifyChainId() error {
	chainId, err := w.sdk.ChainID(context.Background())
	if err != nil {
		fmt.Println("check err3:", err)
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
func (w *Wallet) GetBalance() (balance *big.Int, err error) {
	return w.sdk.BalanceAt(context.Background(), w.account.Address, nil)
}
func (w *Wallet) GetNonce(address common.Address) (uint64, error) {
	return w.getNonceAt()
}

func (w *Wallet) CurrentAccount() accounts.Account {
	return w.currentAccount()
}
func (w *Wallet) currentAccount() accounts.Account {
	return w.account
}
func (w *Wallet) ChainID() *big.Int {
	return big.NewInt(0).SetUint64(w.chainId)
}
func (w *Wallet) GasPrice(ctx context.Context) (*big.Int, error) {
	return w.sdk.SuggestGasPrice(ctx)
}
func (w *Wallet) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return w.sdk.EstimateGas(ctx, msg)
}

func (w *Wallet) GasTipCap(ctx context.Context) (*big.Int, error) {
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

/*to create chains wallet...*/
func (w *Wallet) createChainWallet(ID uint64) (IWallet, error) {
	switch ID {
	case base.ETH:
		return w, nil
	case base.TOP:
		return w, nil
	}
	return nil, fmt.Errorf("unsupport chain:%v", ID)
}
