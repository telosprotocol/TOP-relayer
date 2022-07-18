package wallet

import (
	"context"
	"math/big"
	"toprelayer/sdk"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet struct {
	chainId  uint64
	provider Provider
	account  accounts.Account // active account
	sdk      *sdk.SDK
}

//basic interface
type IWallet interface {
	GetBalance(address common.Address) (balance *big.Int, err error)
	CurrentAccount() accounts.Account
	GetNonce(address common.Address) (uint64, error)

	GasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, target *common.Address, gasPrice *big.Int, data []byte) (uint64, error)
	GasTipCap(ctx context.Context) (*big.Int, error)
	ChainID() *big.Int
	SignTx(tx *types.Transaction) (signedTx *types.Transaction, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}
