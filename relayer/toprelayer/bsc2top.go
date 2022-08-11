package toprelayer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	"toprelayer/contract/top/ethclient"
	"toprelayer/relayer/toprelayer/parlia"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wonderivan/logger"
)

var (
	bscClientContract = common.HexToAddress("0xff00000000000000000000000000000000000003")
)

type Bsc2TopRelayer struct {
	context.Context
	crossChainName string
	chainId        uint64
	wallet         wallet.IWallet
	topsdk         *topsdk.TopSdk
	contract       common.Address
	ethsdk         *ethsdk.EthSdk
	transactor     *ethclient.EthClientTransactor
	callerSession  *ethclient.EthClientCallerSession
	parlia         *parlia.Parlia
}

func (relayer *Bsc2TopRelayer) Init(crossChainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	relayer.crossChainName = crossChainName
	relayer.chainId = cfg.ChainId
	topsdk, err := topsdk.NewTopSdk(cfg.Url)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewTopSdk error:", err)
		return err
	}
	relayer.topsdk = topsdk

	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass, cfg.ChainId)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	relayer.ethsdk, err = ethsdk.NewEthSdk(listenUrl)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewEthSdk error:", crossChainName, listenUrl)
		return err
	}
	relayer.contract = bscClientContract
	relayer.transactor, err = ethclient.NewEthClientTransactor(relayer.contract, topsdk)
	if err != nil {
		logger.Error("TopRelayer from", relayer.crossChainName, "NewEthClientTransactor error:", relayer.contract)
		return err
	}

	relayer.callerSession = new(ethclient.EthClientCallerSession)
	relayer.callerSession.Contract, err = ethclient.NewEthClientCaller(relayer.contract, topsdk)
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

	relayer.parlia = parlia.New(relayer.ethsdk)

	return nil
}

func (relayer *Bsc2TopRelayer) ChainId() uint64 {
	return relayer.chainId
}

func (et *Bsc2TopRelayer) submitEthHeader(header []byte) error {
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
	packHeader, err := ethclient.PackSyncParam(header)
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
func (et *Bsc2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
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

func (et *Bsc2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start TopRelayer from %v... chainid: %v, subBatch: %v certaintyBlocks: %v", et.crossChainName, et.chainId, BATCH_NUM, CONFIRM_NUM)
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
				continue
			}
			logger.Info("TopRelayer from", et.crossChainName, "check dest top Height:", destHeight)
			if destHeight != 0 {
				err = et.parlia.Init(destHeight)
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
	logger.Error("relayer [%v] timeout.", et.chainId)
	return nil
}

func (et *Bsc2TopRelayer) signAndSendTransactions(lo, hi uint64) error {
	var batch []byte
	for h := lo; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			break
		}
		rlp_bytes, err := et.parlia.GetLastSnapBytes(header)
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
