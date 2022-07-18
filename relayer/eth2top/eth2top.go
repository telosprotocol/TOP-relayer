package eth2top

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/contract/top/ethclient"
	"toprelayer/relayer/eth2top/ethashapp"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	FATALTIMEOUT int64 = 24 //hours
	SUCCESSDELAY int64 = 10
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 60

	CONFIRM_NUM uint64 = 25
	BATCH_NUM   uint64 = 5
)

type Eth2TopRelayer struct {
	context.Context
	contract   common.Address
	chainId    uint64
	wallet     wallet.IWallet
	topsdk     *topsdk.TopSdk
	ethsdk     *ethsdk.EthSdk
	transactor *ethclient.EthClientTransactor
	caller     *ethclient.EthClientCaller
}

type void struct{}

func (et *Eth2TopRelayer) Init(topUrl, ethUrl, keypath, pass string, chainid uint64, contract common.Address) error {
	topsdk, err := topsdk.NewTopSdk(topUrl)
	if err != nil {
		return err
	}
	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
	if err != nil {
		return err
	}

	et.topsdk = topsdk
	et.ethsdk = ethsdk
	et.contract = contract
	et.chainId = chainid

	w, err := wallet.NewWallet(topUrl, keypath, pass, chainid)
	if err != nil {
		return err
	}
	et.wallet = w

	et.transactor, err = ethclient.NewEthClientTransactor(contract, topsdk)
	if err != nil {
		return err
	}
	et.caller, err = ethclient.NewEthClientCaller(contract, topsdk)
	if err != nil {
		return err
	}

	return nil
}

func (et *Eth2TopRelayer) ChainId() uint64 {
	return et.chainId
}

func (et *Eth2TopRelayer) submitEthHeader(header []byte, nonce uint64) error {
	gaspric, err := et.wallet.GasPrice(context.Background())
	if err != nil {
		logger.Error("Eth2TopRelayer GasPrice:%v", err)
		return err
	}
	packHeader, err := ethclient.PackSyncParam(header)
	if err != nil {
		logger.Error("Eth2TopRelayer PackSyncParam:%v", err)
		return err
	}
	gaslimit, err := et.wallet.EstimateGas(context.Background(), &et.contract, gaspric, packHeader)
	if err != nil {
		logger.Error("EstimateGas error:", err)
		return err
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      et.wallet.CurrentAccount().Address,
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  gaslimit,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    et.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	sigTx, err := et.transactor.Sync(ops, header)
	if err != nil {
		logger.Error("Eth2TopRelayer sync:%v", err)
		return err
	}

	logger.Info("tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", et.wallet.CurrentAccount().Address, nonce, gaspric, sigTx.Hash(), len(header))
	return nil
}

//callback function to sign tx before send.
func (et *Eth2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := et.wallet.CurrentAccount()
	if strings.EqualFold(acc.Address.Hex(), addr.Hex()) {
		stx, err := et.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("address:%v not available", addr)
}

func (et *Eth2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayer relayer... chainid: %v, subBatch: %v certaintyBlocks: %v", et.chainId, BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Info("Eth2TopRelayer set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				opts := &bind.CallOpts{
					Pending:     false,
					From:        et.wallet.CurrentAccount().Address,
					BlockNumber: nil,
					Context:     context.Background(),
				}
				destHeight, err := et.caller.GetHeight(opts)
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check dest top Height: %v", destHeight)
				if destHeight == 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("top eth-client not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				srcHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check src eth Height: %v", srcHeight)

				if destHeight+1+CONFIRM_NUM > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}
				// check fork
				var checkError bool = false
				for {
					header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(destHeight))
					if err != nil {
						logger.Error("HeaderByNumber ", err)
						checkError = true
						break
					}
					// get known hashes with destHeight, mock now
					isKnown, err := et.caller.IsKnown(opts, header.Number, header.Hash())
					if err != nil {
						logger.Error("IsKnown ", err)
						checkError = true
						break
					}
					if isKnown {
						logger.Debug("%v hash is known", header.Number)
						break
					} else {
						logger.Debug("%v hash is not known", header.Number)
						destHeight -= 1
					}
				}
				if checkError {
					delay = time.Duration(ERRDELAY)
					break
				}

				syncStartHeight := destHeight + 1
				syncNum := srcHeight - CONFIRM_NUM - destHeight
				if syncNum > BATCH_NUM {
					syncNum = BATCH_NUM
				}
				syncEndHeight := syncStartHeight + syncNum - 1
				logger.Info("Eth2TopRelayer sync from %v to %v", syncStartHeight, syncEndHeight)

				err = et.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if err != nil {
					logger.Error("Eth2TopRelayer signAndSendTransactions failed:%v", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer sync round finish")
				delay = time.Duration(SUCCESSDELAY)
				// break
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout.", et.chainId)
	return nil
}

func (et *Eth2TopRelayer) signAndSendTransactions(lo, hi uint64) error {
	var batch []byte
	nonce, err := et.wallet.GetNonce(et.wallet.CurrentAccount().Address)
	if err != nil {
		logger.Error(err)
		return err
	}

	for h := lo; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			return err
		}
		ethashproof, err := ethashapp.EthashWithProofs(h, header)
		if err != nil {
			logger.Error(err)
			return err
		}
		rlp_bytes, err := rlp.EncodeToBytes(ethashproof)
		if err != nil {
			logger.Fatal("rlp encode error: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}

	// maybe verify block
	// if et.chainId == topChainId {
	// 	for _, header := range headers {
	// 		et.verifyBlocks(header)
	// 	}
	// }
	err = et.submitEthHeader(batch, nonce)
	if err != nil {
		logger.Error("Eth2TopRelayer submitHeaders failed:", err)
		return err
	}

	return nil
}
