package wallet

import (
	"context"
	"math/big"
	"toprelayer/sdk"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet struct {
	chainId  uint64
	provider *keystore.KeyStore
	account  accounts.Account // active account
	sdk      *sdk.SDK
}

//basic interface
type IWallet interface {
	Address() common.Address
	ChainID(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)

	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, target *common.Address, data []byte) (uint64, error)

	SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
}
