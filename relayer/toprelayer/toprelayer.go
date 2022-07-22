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
	"toprelayer/relayer/toprelayer/ethashapp"
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

	CONFIRM_NUM uint64 = 5
	BATCH_NUM   uint64 = 5
)

var (
	systemSyncContracts = map[string]common.Address{
		config.ETH_CHAIN:  common.HexToAddress("0xff00000000000000000000000000000000000002"),
		config.BSC_CHAIN:  common.HexToAddress("0xff00000000000000000000000000000000000003"),
		config.HECO_CHAIN: common.HexToAddress("0xff00000000000000000000000000000000000004")}
)

type TopRelayer struct {
	context.Context
	crossChainName string
	chainId        uint64
	wallet         wallet.IWallet
	topsdk         *topsdk.TopSdk
	contract       common.Address
	ethsdk         *ethsdk.EthSdk
	transactor     *ethclient.EthClientTransactor
	callerSession  *ethclient.EthClientCallerSession
}

type void struct{}

func (et *TopRelayer) Init(crossChainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	et.crossChainName = crossChainName
	et.chainId = cfg.ChainId
	topsdk, err := topsdk.NewTopSdk(cfg.Url)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "NewTopSdk error:", err)
		return err
	}
	et.topsdk = topsdk

	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass, cfg.ChainId)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "NewWallet error:", err)
		return err
	}
	et.wallet = w

	et.ethsdk, err = ethsdk.NewEthSdk(listenUrl)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "NewEthSdk error:", crossChainName, listenUrl)
		return err
	}
	et.contract = systemSyncContracts[crossChainName]
	et.transactor, err = ethclient.NewEthClientTransactor(et.contract, topsdk)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "NewEthClientTransactor error:", et.contract)
		return err
	}

	et.callerSession = new(ethclient.EthClientCallerSession)
	et.callerSession.Contract, err = ethclient.NewEthClientCaller(et.contract, topsdk)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "NewEthClientCaller error:", et.contract)
		return err
	}
	et.callerSession.CallOpts = bind.CallOpts{
		Pending:     false,
		From:        et.wallet.CurrentAccount().Address,
		BlockNumber: nil,
		Context:     context.Background(),
	}

	return nil
}

func (et *TopRelayer) ChainId() uint64 {
	return et.chainId
}

func (et *TopRelayer) submitEthHeader(header []byte) error {
	nonce, err := et.wallet.GetNonce(et.wallet.CurrentAccount().Address)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "GetNonce error:", err)
		return err
	}
	gaspric, err := et.wallet.GasPrice(context.Background())
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "GasPrice error:", err)
		return err
	}
	packHeader, err := ethclient.PackSyncParam(header)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "PackSyncParam error:", err)
		return err
	}
	gaslimit, err := et.wallet.EstimateGas(context.Background(), &et.contract, gaspric, packHeader)
	if err != nil {
		logger.Error("TopRelayer from", et.crossChainName, "EstimateGas error:", err)
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
		logger.Error("TopRelayer from", et.crossChainName, " sync error:", err)
		return err
	}

	logger.Info("TopRelayer from %v tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", et.crossChainName, et.wallet.CurrentAccount().Address, nonce, gaspric, sigTx.Hash(), len(header))
	return nil
}

//callback function to sign tx before send.
func (et *TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := et.wallet.CurrentAccount()
	if strings.EqualFold(acc.Address.Hex(), addr.Hex()) {
		stx, err := et.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("TopRelayer address:%v not available", addr)
}

func (et *TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
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
					logger.Error(err)
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
				var checkError bool = false
				for {
					header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(destHeight))
					if err != nil {
						logger.Debug("TopRelayer from", et.crossChainName, "HeaderByNumber error:", err)
						checkError = true
						break
					}
					// get known hashes with destHeight, mock now
					isKnown, err := et.callerSession.IsKnown(header.Number, header.Hash())
					if err != nil {
						logger.Debug("TopRelayer from", et.crossChainName, "IsKnown error:", err)
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

func (et *TopRelayer) signAndSendTransactions(lo, hi uint64) error {
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
			logger.Error("TopRelayer from", et.crossChainName, "submitHeaders failed:", err)
			return err
		}
	}

	return nil
}
