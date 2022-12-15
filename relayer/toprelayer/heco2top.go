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
	"toprelayer/relayer/toprelayer/congress"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

var (
	hecoClientContract = common.HexToAddress("0xff00000000000000000000000000000000000004")
)

type Heco2TopRelayer struct {
	wallet        *wallet.Wallet
	ethsdk        *ethclient.Client
	transactor    *ethbridge.EthClientTransactor
	callerSession *ethbridge.EthClientCallerSession
	congress      *congress.Congress
}

func (relayer *Heco2TopRelayer) Init(cfg *config.Relayer, listenUrl []string, pass string) error {
	w, err := wallet.NewEthWallet(cfg.Url[0], listenUrl[0], cfg.KeyPath, pass)
	if err != nil {
		logger.Error("Heco2TopRelayer NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	relayer.ethsdk, err = ethclient.Dial(listenUrl[0])
	if err != nil {
		logger.Error("Heco2TopRelayer ethsdk create error:", err)
		return err
	}

	topethlient, err := ethclient.Dial(cfg.Url[0])
	if err != nil {
		logger.Error("Heco2TopRelayer new topethlient error:", err)
		return err
	}
	relayer.transactor, err = ethbridge.NewEthClientTransactor(hecoClientContract, topethlient)
	if err != nil {
		logger.Error("Heco2TopRelayer NewEthClientTransactor error:", err)
		return err
	}

	relayer.callerSession = new(ethbridge.EthClientCallerSession)
	relayer.callerSession.Contract, err = ethbridge.NewEthClientCaller(hecoClientContract, topethlient)
	if err != nil {
		logger.Error("Heco2TopRelayer NewEthClientCaller error:", err)
		return err
	}
	relayer.callerSession.CallOpts = bind.CallOpts{
		Pending:     false,
		From:        relayer.wallet.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}

	relayer.congress = congress.New(relayer.ethsdk)

	return nil
}

func (et *Heco2TopRelayer) submitEthHeader(header []byte) error {
	nonce, err := et.wallet.NonceAt(context.Background(), et.wallet.Address(), nil)
	if err != nil {
		logger.Error("Heco2TopRelayer NonceAt error:", err)
		return err
	}
	gaspric, err := et.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error("Heco2TopRelayer SuggestGasPrice error:", err)
		return err
	}
	packHeader, err := ethbridge.PackSyncParam(header)
	if err != nil {
		logger.Error("Heco2TopRelayer PackSyncParam error:", err)
		return err
	}
	gaslimit, err := et.wallet.EstimateGas(context.Background(), &hecoClientContract, packHeader)
	if err != nil {
		logger.Error("Heco2TopRelayer EstimateGas error:", err)
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
		logger.Error("Heco2TopRelayer sync error:", err)
		return err
	}

	logger.Info("Heco2TopRelayer tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", et.wallet.Address(), nonce, gaspric, sigTx.Hash(), len(header))
	return nil
}

//callback function to sign tx before send.
func (et *Heco2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := et.wallet.Address()
	if strings.EqualFold(acc.Hex(), addr.Hex()) {
		stx, err := et.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("TopRelayer address:%v not available", addr)
}

func (et *Heco2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Heco2TopRelayer start... subBatch: %v certaintyBlocks: %v", BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Debug("Heco2TopRelayer set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		for {
			destHeight, err := et.callerSession.GetHeight()
			if err != nil {
				logger.Error("Heco2TopRelayer get height error:", err)
				time.Sleep(time.Second * time.Duration(ERRDELAY))
				continue
			}
			logger.Info("Heco2TopRelayer check dest top Height:", destHeight)
			if destHeight != 0 {
				err = et.congress.Init(destHeight)
				if err == nil {
					break
				} else {
					logger.Error("Heco2TopRelayer congress init error:", err)
				}
			} else {
				logger.Info("Heco2TopRelayer not init yet")
			}
			time.Sleep(time.Second * time.Duration(ERRDELAY))
		}

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				destHeight, err := et.callerSession.GetHeight()
				if err != nil {
					logger.Error("Heco2TopRelayer get height error:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Heco2TopRelayer check dest top Height:", destHeight)
				if destHeight == 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("Heco2TopRelayer reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("Heco2TopRelayer not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				srcHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error("Heco2TopRelayer get number error:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Heco2TopRelayer check src eth Height:", srcHeight)

				if destHeight+1+CONFIRM_NUM > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("Heco2TopRelayer reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("Heco2TopRelayer waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}

				// check fork
				checkError := false
				for {
					header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(destHeight))
					if err != nil {
						logger.Error("Heco2TopRelayer HeaderByNumber error:", err)
						checkError = true
						break
					}
					// get known hashes with destHeight, mock now
					isKnown, err := et.callerSession.IsKnown(header.Number, header.Hash())
					if err != nil {
						logger.Error("Heco2TopRelayer IsKnown error:", err)
						checkError = true
						break
					}
					if isKnown {
						logger.Debug("%v hash is known", header.Number)
						break
					} else {
						logger.Warn("%v hash is not known", header.Number)
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
				logger.Info("Heco2TopRelayer sync from %v to %v", syncStartHeight, syncEndHeight)

				err = et.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if err != nil {
					logger.Error("Heco2TopRelayer signAndSendTransactions failed:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("Heco2TopRelayer reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Heco2TopRelayer sync round finish")
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
	logger.Error("Heco2TopRelayer timeout")
	return nil
}

func (et *Heco2TopRelayer) signAndSendTransactions(lo, hi uint64) error {
	var batch []byte
	for h := lo; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			break
		}
		rlp_bytes, err := et.congress.GetLastSnapBytes(header)
		if err != nil {
			logger.Error(err)
			return err
		}
		batch = append(batch, rlp_bytes...)
	}

	if len(batch) > 0 {
		err := et.submitEthHeader(batch)
		if err != nil {
			logger.Error("Heco2TopRelayer submitHeaders failed:", err)
			return err
		}
	}

	return nil
}

func (et *Heco2TopRelayer) GetInitData() ([]byte, error) {
	destHeight, err := et.callerSession.GetHeight()
	if err != nil {
		logger.Error("Heco2TopRelayer get height error:", err)
		return nil, err
	}
	height := (destHeight - 11) / 200 * 200

	logger.Error("heco init with height: %v - %v", height, height+11)
	var batch []byte
	for i := height; i <= height+11; i++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(i))
		if err != nil {
			logger.Error("Heco2TopRelayer HeaderByNumber error:", err)
			return nil, err
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			logger.Error("Heco2TopRelayer EncodeToBytes error:", err)
			return nil, err
		}
		batch = append(batch, rlp_bytes...)
	}

	return batch, nil
}
