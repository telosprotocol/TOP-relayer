package top2eth

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	"toprelayer/contract/eth/topclient"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/util"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	FATALTIMEOUT int64 = 24 //hours
	SUCCESSDELAY int64 = 60
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 60

	CONFIRM_NUM int = 2
	BATCH_NUM   int = 10

	ELECTION_BLOCK    = "election"
	AGGREGATE_BLOCK   = "aggregate"
	TRANSACTION_BLOCK = "transactions"
)

var (
	sendFlag = map[string]uint64{
		config.ETH_CHAIN:  0x1,
		config.BSC_CHAIN:  0x2,
		config.HECO_CHAIN: 0x3}
)

type Top2EthRelayer struct {
	context.Context
	name       string
	chainId    uint64
	contract   common.Address
	wallet     wallet.IWallet
	ethsdk     *ethsdk.EthSdk
	topsdk     *topsdk.TopSdk
	transactor *topclient.TopClientTransactor
	caller     *topclient.TopClientCaller
}

func (te *Top2EthRelayer) Init(chainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	te.name = chainName
	te.chainId = cfg.ChainId
	ethsdk, err := ethsdk.NewEthSdk(cfg.Url)
	if err != nil {
		return err
	}
	topsdk, err := topsdk.NewTopSdk(listenUrl)
	if err != nil {
		logger.Error("Top2EthRelayer NewTopSdk error:", err)
		return err
	}
	te.topsdk = topsdk
	te.ethsdk = ethsdk
	te.contract = common.HexToAddress(cfg.Contract)

	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass, cfg.ChainId)
	if err != nil {
		logger.Error("Top2EthRelayer NewWallet error:", err)
		return err
	}
	te.wallet = w

	te.transactor, err = topclient.NewTopClientTransactor(te.contract, ethsdk)
	if err != nil {
		logger.Error("Top2EthRelayer NewTopClientTransactor error:", err)
		return err
	}
	te.caller, err = topclient.NewTopClientCaller(te.contract, ethsdk)
	if err != nil {
		logger.Error("Top2EthRelayer NewTopClientCaller error:", err)
		return err
	}
	return nil
}

func (te *Top2EthRelayer) ChainId() uint64 {
	return te.chainId
}

func (te *Top2EthRelayer) submitTopHeader(headers []byte, nonce uint64) error {
	logger.Info("Top2EthRelayer submitTopHeader length: %v,chainid: %v", len(headers), te.chainId)
	logger.Info("Top2EthRelayer raw data: %v", common.Bytes2Hex(headers))
	gaspric, err := te.wallet.GasPrice(context.Background())
	if err != nil {
		logger.Error("Top2EthRelayer GasPrice:", err)
		return err
	}
	packHeaders, err := topclient.PackSyncParam(headers)
	if err != nil {
		logger.Error("Eth2TopRelayer PackSyncParam:%v", err)
		return err
	}
	gaslimit, err := te.wallet.EstimateGas(context.Background(), &te.contract, gaspric, packHeaders)
	if err != nil {
		return err
	}
	//test mock
	//gaslimit := uint64(500000)

	balance, err := te.wallet.GetBalance(te.wallet.CurrentAccount().Address)
	if err != nil {
		return err
	}
	logger.Info("account[%v] balance:%v,nonce:%v,gasprice:%v,gaslimit:%v", te.wallet.CurrentAccount().Address, balance.Uint64(), nonce, gaspric.Uint64(), gaslimit)
	if balance.Uint64() <= gaspric.Uint64()*gaslimit {
		return fmt.Errorf("account[%v] not sufficient funds,balance:%v", te.wallet.CurrentAccount().Address, balance.Uint64())
	}

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:     te.wallet.CurrentAccount().Address,
		Nonce:    big.NewInt(0).SetUint64(nonce),
		GasPrice: gaspric,
		GasLimit: gaslimit,
		Signer:   te.signTransaction,
		Context:  context.Background(),
		NoSend:   true,
	}

	contractcaller, err := topclient.NewTopClientTransactor(te.contract, te.ethsdk)
	if err != nil {
		logger.Error("Top2EthRelayer NewBridgeTransactor:", err)
		return err
	}

	sigTx, err := contractcaller.AddLightClientBlocks(ops, headers)
	if err != nil {
		logger.Error("Top2EthRelayer AddLightClientBlocks error:", err)
		return err
	}

	if ops.NoSend {
		err = util.VerifyEthSignature(sigTx)
		if err != nil {
			logger.Error("Top2EthRelayer VerifyEthSignature error:", err)
			return err
		}

		err := te.ethsdk.SendTransaction(ops.Context, sigTx)
		if err != nil {
			logger.Error("Top2EthRelayer SendTransaction error:", err)
			return err
		}
	}
	logger.Info("hash:%v", sigTx.Hash())
	return nil
}

//callback function to sign tx before send.
func (te *Top2EthRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := te.wallet.CurrentAccount()
	if strings.EqualFold(acc.Address.Hex(), addr.Hex()) {
		stx, err := te.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("address:%v not available", addr)
}

func (te *Top2EthRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Top2EthRelayer relayer... chainid: %v, subBatch: %v certaintyBlocks: %v", te.chainId, BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Info("Top2EthRelayer set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		var lastSubHeight uint64 = 0
		var lastUnsubHeight uint64 = 0

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				opts := &bind.CallOpts{
					Pending:     false,
					From:        te.wallet.CurrentAccount().Address,
					BlockNumber: nil,
					Context:     context.Background(),
				}
				toHeight, err := te.caller.MaxMainHeight(opts)
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Top2EthRelayer dest eth Height: %v", toHeight)
				fromHeight, err := te.topsdk.GetLatestTopElectBlockHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Top2EthRelayer src top Height: %v", fromHeight)

				if lastSubHeight <= toHeight && toHeight < lastUnsubHeight {
					toHeight = lastUnsubHeight
				}
				if toHeight+1+uint64(CONFIRM_NUM) > fromHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("wait src top update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}
				syncStartHeight := toHeight + 1
				// syncNum := fromHeight - uint64(te.certaintyBlocks) - toHeight
				// if syncNum > uint64(te.subBatch) {
				// 	syncNum = uint64(te.subBatch)
				// }
				limitEndHeight := fromHeight - uint64(CONFIRM_NUM)

				subHeight, unsubHeight, err := te.signAndSendTransactions(syncStartHeight, limitEndHeight, uint64(BATCH_NUM))
				if err != nil {
					logger.Error("Top2EthRelayer signAndSendTransactions failed: %v", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if subHeight > lastSubHeight {
					logger.Info("Top2EthRelayer lastSubHeight: %v=>%v", lastSubHeight, subHeight)
					lastSubHeight = subHeight
				}
				if unsubHeight > lastUnsubHeight {
					logger.Info("Top2EthRelayer lastUnsubHeight: %v=>%v", lastUnsubHeight, unsubHeight)
					lastUnsubHeight = unsubHeight
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Top2EthRelayer sync round finish")
				delay = time.Duration(SUCCESSDELAY)
				break
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout.", te.chainId)
	return nil
}

func (te *Top2EthRelayer) signAndSendTransactions(lo, hi, batchNum uint64) (uint64, uint64, error) {
	var lastSubHeight uint64 = 0
	var lastUnsubHeight uint64 = 0
	var batchHeaders [][]byte
	nonce, err := te.wallet.GetNonce(te.wallet.CurrentAccount().Address)
	if err != nil {
		return 0, 0, err
	}

	num := uint64(0)
	flag := sendFlag[te.name]
	for h := lo; h <= hi; h++ {
		block, err := te.topsdk.GetTopElectBlockHeadByHeight(h)
		if err != nil {
			logger.Error("GetTopElectBlockHeadByHeight error:", err)
			break
		}
		logger.Debug("Top block, height: %v, type: %v, chainbits: %v", block.Number, block.BlockType, block.ChainBits)
		batch := false
		if block.BlockType == ELECTION_BLOCK {
			batch = true
		} else if block.BlockType == AGGREGATE_BLOCK {
			blockFlag, err := strconv.ParseInt(block.ChainBits, 0, 64)
			if err != nil {
				logger.Error("ParseInt error:", err)
				break
			}
			if int64(flag)&blockFlag > 0 {
				batch = true
			}
		}
		if batch {
			logger.Debug("batch header")
			bytes := common.Hex2Bytes(block.Header[2:])
			batchHeaders = append(batchHeaders, bytes)
			lastSubHeight = h
			num += 1
			if num >= batchNum {
				break
			}
		} else {
			lastUnsubHeight = h
		}
	}
	if len(batchHeaders) > 0 {
		data, err := rlp.EncodeToBytes(batchHeaders)
		if err != nil {
			logger.Error("Eth2TopRelayer EncodeHeaders failed:", err)
			return 0, 0, err
		}

		err = te.submitTopHeader(data, nonce)
		if err != nil {
			logger.Error("Top2EthRelayer submitHeaders failed:", err)
			return 0, 0, err
		}
	}

	return lastSubHeight, lastUnsubHeight, nil
}
