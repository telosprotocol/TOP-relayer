package wallet

import (
	"context"
	"math/big"
	"toprelayer/sdk"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type KeyStore struct {
	*keystore.KeyStore
}

type Wallet struct {
	chainId  uint64
	provider Provider
	account  accounts.Account // active account
	nonce    uint64           // account nonces
	sdk      *sdk.SDK
}

//basic interface
type IWallet interface {
	GetBalance() (balance *big.Int, err error)
	CurrentAccount() accounts.Account
	GetNonce(address common.Address) (uint64, error)

	GasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	GasTipCap(ctx context.Context) (*big.Int, error)
	ChainID() *big.Int
	SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}
