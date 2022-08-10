package toprelayer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	ethbridge "toprelayer/contract/top/ethclient"
	"toprelayer/relayer/monitor"
	"toprelayer/relayer/toprelayer/ethashapp"
	"toprelayer/sdk/topsdk"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	FATALTIMEOUT int64 = 24 //hours
	SUCCESSDELAY int64 = 10
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 60

	CONFIRM_NUM uint64 = 5
	BATCH_NUM   uint64 = 5
)

var (
	ethClientSystemContract = common.HexToAddress("0xff00000000000000000000000000000000000002")
)

type Eth2TopRelayer struct {
	context.Context
	wallet        wallet.IWallet
	ethsdk        *ethclient.Client
	transactor    *ethbridge.EthClientTransactor
	callerSession *ethbridge.EthClientCallerSession
	monitor       *monitor.Monitor
}

type void struct{}

func (relayer *Eth2TopRelayer) Init(crossChainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass, cfg.ChainId)
	if err != nil {
		logger.Error("Eth2TopRelayer NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	relayer.ethsdk, err = ethclient.Dial(listenUrl)
	if err != nil {
		logger.Error("Eth2TopRelayer ethclient.Dial error:", err)
		return err
	}
	topsdk, err := topsdk.NewTopSdk(cfg.Url)
	if err != nil {
		logger.Error("Eth2TopRelayer NewTopSdk error:", err)
		return err
	}

	relayer.transactor, err = ethbridge.NewEthClientTransactor(ethClientSystemContract, topsdk)
	if err != nil {
		logger.Error("Eth2TopRelayer NewEthClientTransactor error:", err)
		return err
	}

	relayer.callerSession = new(ethbridge.EthClientCallerSession)
	relayer.callerSession.Contract, err = ethbridge.NewEthClientCaller(ethClientSystemContract, topsdk)
	if err != nil {
		logger.Error("Eth2TopRelayer NewEthClientCaller error:", err)
		return err
	}
	relayer.callerSession.CallOpts = bind.CallOpts{
		Pending:     false,
		From:        relayer.wallet.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}
	relayer.monitor, err = monitor.New(relayer.wallet.Address(), cfg.Url)
	if err != nil {
		logger.Error("Eth2TopRelayer New monitor error", err)
		return err
	}
	return nil
}

func (et *Eth2TopRelayer) submitEthHeader(header []byte) error {
	nonce, err := et.wallet.NonceAt(et.wallet.Address())
	if err != nil {
		logger.Error("Eth2TopRelayer GetNonce error:", err)
		return err
	}
	gaspric, err := et.wallet.SuggestGasPrice()
	if err != nil {
		logger.Error("Eth2TopRelayer GasPrice error:", err)
		return err
	}
	packHeader, err := ethbridge.PackSyncParam(header)
	if err != nil {
		logger.Error("Eth2TopRelayer PackSyncParam error:", err)
		return err
	}
	gaslimit, err := et.wallet.EstimateGas(context.Background(), &ethClientSystemContract, packHeader)
	if err != nil {
		logger.Error("Eth2TopRelayer EstimateGas error:", err)
		return err
	}

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      et.wallet.Address(),
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
		logger.Error("Eth2TopRelayer sync error:", err)
		return err
	}
	et.monitor.AddTx(sigTx.Hash())
	logger.Info("Eth2TopRelayer tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", et.wallet.Address(), nonce, gaspric, sigTx.Hash(), len(header))
	return nil
}

//callback function to sign tx before send.
func (et *Eth2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := et.wallet.Address()
	if strings.EqualFold(acc.Hex(), addr.Hex()) {
		stx, err := et.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("Eth2TopRelayer address:%v not available", addr)
}

func (et *Eth2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayer, subBatch: %v certaintyBlocks: %v", BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Debug("Eth2TopRelayer set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				destHeight, err := et.callerSession.GetHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check dest top Height:", destHeight)
				if destHeight == 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("Eth2TopRelayer reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("Eth2TopRelayer not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				srcHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check src eth Height:", srcHeight)

				if destHeight+1+CONFIRM_NUM > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("Eth2TopRelayer reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("Eth2TopRelayer waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}
				// check fork
				var checkError bool = false
				for {
					header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(destHeight))
					if err != nil {
						logger.Debug("Eth2TopRelayer HeaderByNumber error:", err)
						checkError = true
						break
					}
					// get known hashes with destHeight, mock now
					isKnown, err := et.callerSession.IsKnown(header.Number, header.Hash())
					if err != nil {
						logger.Error("Eth2TopRelayer IsKnown error:", err)
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
					logger.Error("Eth2TopRelayer signAndSendTransactions failed:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("Eth2TopRelayer reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer sync round finish")
				if syncNum == BATCH_NUM {
					delay = time.Duration(SUCCESSDELAY)
				} else {
					delay = time.Duration(WAITDELAY)
				}
				// break
			}
		}
	}(done)

	<-done
	logger.Error("Eth2TopRelayer timeout")
	return nil
}

func (et *Eth2TopRelayer) signAndSendTransactions(lo, hi uint64) error {
	var batch []byte
	for h := lo; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			break
		}
		ethashproof, err := ethashapp.EthashWithProofs(h, header)
		if err != nil {
			logger.Error(err)
			return err
		}
		rlp_bytes, err := rlp.EncodeToBytes(ethashproof)
		if err != nil {
			logger.Error("rlp encode error: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}

	// maybe verify block
	// if et.chainId == topChainId {
	// 	for _, header := range headers {
	// 		et.verifyBlocks(header)
	// 	}
	// }
	if len(batch) > 0 {
		err := et.submitEthHeader(batch)
		if err != nil {
			logger.Error("Eth2TopRelayer submitHeaders failed:", err)
			return err
		}
	}

	return nil
}
