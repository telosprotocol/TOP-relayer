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
	"github.com/wonderivan/logger"
)

var (
	hecoClientContract = common.HexToAddress("0xff00000000000000000000000000000000000004")
)

type Heco2TopRelayer struct {
	context.Context
	crossChainName string
	wallet         *wallet.Wallet
	contract       common.Address
	ethsdk         *ethclient.Client
	transactor     *ethbridge.EthClientTransactor
	callerSession  *ethbridge.EthClientCallerSession
	congress       *congress.Congress
}

func (relayer *Heco2TopRelayer) Init(crossChainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	relayer.crossChainName = crossChainName

	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	relayer.ethsdk, err = ethclient.Dial(listenUrl)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "ethsdk create error:", crossChainName, listenUrl)
		return err
	}
	relayer.contract = hecoClientContract

	topethlient, err := ethclient.Dial(cfg.Url)
	if err != nil {
		logger.Error("TopRelayer new topethlient error:", err)
		return err
	}
	relayer.transactor, err = ethbridge.NewEthClientTransactor(relayer.contract, topethlient)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewEthClientTransactor error:", relayer.contract)
		return err
	}

	relayer.callerSession = new(ethbridge.EthClientCallerSession)
	relayer.callerSession.Contract, err = ethbridge.NewEthClientCaller(relayer.contract, topethlient)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewEthClientCaller error:", relayer.contract)
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
		logger.Error("TopRelayer from", et.crossChainName, "GetNonce error:", err)
		return err
	}
	gaspric, err := et.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "GasPrice error:", err)
		return err
	}
	packHeader, err := ethbridge.PackSyncParam(header)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "PackSyncParam error:", err)
		return err
	}
	gaslimit, err := et.wallet.EstimateGas(context.Background(), &et.contract, packHeader)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "EstimateGas error:", err)
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
		logger.Error("TopRelayer from", et.crossChainName, " sync error:", err)
		return err
	}

	logger.Info("TopRelayer from %v tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", et.crossChainName, et.wallet.Address(), nonce, gaspric, sigTx.Hash(), len(header))
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
	logger.Info("Start TopRelayer from %v... subBatch: %v certaintyBlocks: %v", et.crossChainName, BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Debug("TopRelayer from %v set timeout: %v hours", et.crossChainName, FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		for {
			destHeight, err := et.callerSession.GetHeight()
			if err != nil {
				logger.Error("TopRelayer from ", et.crossChainName, " get height error:", err)
				time.Sleep(time.Second * time.Duration(ERRDELAY))
				continue
			}
			logger.Info("TopRelayer from", et.crossChainName, "check dest top Height:", destHeight)
			if destHeight != 0 {
				err = et.congress.Init(destHeight)
				if err == nil {
					break
				} else {
					logger.Error("TopRelayer from ", et.crossChainName, " congress init error:", err)
				}
			} else {
				logger.Info("TopRelayer from ", et.crossChainName, " not init yet")
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
					logger.Error("TopRelayer from ", et.crossChainName, " get height error:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("TopRelayer from", et.crossChainName, "check dest top Height:", destHeight)
				if destHeight == 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("TopRelayer from", et.crossChainName, "reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("TopRelayer from ", et.crossChainName, " not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				srcHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error("TopRelayer from ", et.crossChainName, " get number error:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("TopRelayer from", et.crossChainName, "check src eth Height:", srcHeight)

				if destHeight+1+CONFIRM_NUM > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("TopRelayer from", et.crossChainName, "reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("TopRelayer from", et.crossChainName, "waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}

				// check fork
				checkError := false
				for {
					header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(destHeight))
					if err != nil {
						logger.Error("TopRelayer from", et.crossChainName, "HeaderByNumber error:", err)
						checkError = true
						break
					}
					// get known hashes with destHeight, mock now
					isKnown, err := et.callerSession.IsKnown(header.Number, header.Hash())
					if err != nil {
						logger.Error("TopRelayer from", et.crossChainName, "IsKnown error:", err)
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
				logger.Info("TopRelayer from %v sync from %v to %v", et.crossChainName, syncStartHeight, syncEndHeight)

				err = et.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if err != nil {
					logger.Error("TopRelayer from", et.crossChainName, "signAndSendTransactions failed:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("TopRelayer from", et.crossChainName, "reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("TopRelayer from", et.crossChainName, "sync round finish")
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
	logger.Error("relayer [%v] timeout.", et.crossChainName)
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

	// maybe verify block
	// if et.chainId == topChainId {
	// 	for _, header := range headers {
	// 		et.verifyBlocks(header)
	// 	}
	// }
	if len(batch) > 0 {
		err := et.submitEthHeader(batch)
		if err != nil {
			logger.Error("TopRelayer from", et.crossChainName, "submitHeaders failed:", err)
			return err
		}
	}

	return nil
}
