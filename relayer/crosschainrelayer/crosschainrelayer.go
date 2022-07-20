package crosschainrelayer

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

	CONFIRM_NUM uint64 = 2
	BATCH_NUM   uint64 = 10

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

type CrossChainRelayer struct {
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

func (te *CrossChainRelayer) Init(chainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	te.name = chainName
	te.chainId = cfg.ChainId
	ethsdk, err := ethsdk.NewEthSdk(cfg.Url)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewEthSdk error:", err)
		return err
	}
	topsdk, err := topsdk.NewTopSdk(listenUrl)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewTopSdk error:", err)
		return err
	}
	te.topsdk = topsdk
	te.ethsdk = ethsdk
	te.contract = common.HexToAddress(cfg.Contract)

	w, err := wallet.NewWallet(cfg.Url, cfg.KeyPath, pass, cfg.ChainId)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewWallet error:", err)
		return err
	}
	te.wallet = w

	te.transactor, err = topclient.NewTopClientTransactor(te.contract, ethsdk)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewTopClientTransactor error:", err)
		return err
	}
	te.caller, err = topclient.NewTopClientCaller(te.contract, ethsdk)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewTopClientCaller error:", err)
		return err
	}
	return nil
}

func (te *CrossChainRelayer) ChainId() uint64 {
	return te.chainId
}

func (te *CrossChainRelayer) submitTopHeader(headers []byte, nonce uint64) error {
	logger.Info("CrossChainRelayer", te.name, "raw data:", common.Bytes2Hex(headers))
	gaspric, err := te.wallet.GasPrice(context.Background())
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "GasPrice error:", err)
		return err
	}
	packHeaders, err := topclient.PackSyncParam(headers)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "PackSyncParam error:", err)
		return err
	}
	gaslimit, err := te.wallet.EstimateGas(context.Background(), &te.contract, gaspric, packHeaders)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "EstimateGas error:", err)
		return err
	}
	//test mock
	//gaslimit := uint64(500000)

	balance, err := te.wallet.GetBalance(te.wallet.CurrentAccount().Address)
	if err != nil {
		return err
	}
	if balance.Uint64() <= gaspric.Uint64()*gaslimit {
		return fmt.Errorf("CrossChainRelayer %v account[%v] balance not enough:%v", te.name, te.wallet.CurrentAccount().Address, balance.Uint64())
	}

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:     te.wallet.CurrentAccount().Address,
		Nonce:    big.NewInt(0).SetUint64(nonce),
		GasPrice: gaspric,
		GasLimit: gaslimit,
		Signer:   te.signTransaction,
		Context:  context.Background(),
		NoSend:   false,
	}

	sigTx, err := te.transactor.AddLightClientBlocks(ops, headers)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "AddLightClientBlocks error:", err)
		return err
	}
	logger.Info("CrossChainRelayer", te.name, "tx info, account[%v] balance:%v,nonce:%v,gasprice:%v,gaslimit:%v,length:%v,chainid:%v,hash:%v", te.wallet.CurrentAccount().Address, balance.Uint64(), nonce, gaspric.Uint64(), gaslimit, len(headers), te.chainId, sigTx.Hash())
	return nil
}

//callback function to sign tx before send.
func (te *CrossChainRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
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

func (te *CrossChainRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start CrossChainRelayer %v... chainid: %v, subBatch: %v certaintyBlocks: %v", te.name, te.chainId, BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Info("CrossChainRelayer %v set timeout: %v hours", te.name, FATALTIMEOUT)
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
				logger.Info("CrossChainRelayer", te.name, "dest eth Height:", toHeight)
				fromHeight, err := te.topsdk.GetLatestTopElectBlockHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("CrossChainRelayer", te.name, "src top Height:", fromHeight)

				if lastSubHeight <= toHeight && toHeight < lastUnsubHeight {
					toHeight = lastUnsubHeight
				}
				if toHeight+1+CONFIRM_NUM > fromHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("CrossChainRelayer", te.name, "reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("CrossChainRelayer", te.name, "wait src top update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}
				syncStartHeight := toHeight + 1
				// syncNum := fromHeight - uint64(te.certaintyBlocks) - toHeight
				// if syncNum > uint64(te.subBatch) {
				// 	syncNum = uint64(te.subBatch)
				// }
				limitEndHeight := fromHeight - CONFIRM_NUM

				subHeight, unsubHeight, err := te.signAndSendTransactions(syncStartHeight, limitEndHeight, BATCH_NUM)
				if err != nil {
					logger.Error("CrossChainRelayer", te.name, "signAndSendTransactions failed:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if subHeight > lastSubHeight {
					logger.Info("CrossChainRelayer %v lastSubHeight: %v=>%v", te.name, lastSubHeight, subHeight)
					lastSubHeight = subHeight
				}
				if unsubHeight > lastUnsubHeight {
					logger.Info("CrossChainRelayer %v lastUnsubHeight: %v=>%v", te.name, lastUnsubHeight, unsubHeight)
					lastUnsubHeight = unsubHeight
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("CrossChainRelayer", te.name, "reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("CrossChainRelayer", te.name, "sync round finish")
				delay = time.Duration(SUCCESSDELAY)
				break
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout.", te.chainId)
	return nil
}

func (te *CrossChainRelayer) signAndSendTransactions(lo, hi, batchNum uint64) (uint64, uint64, error) {
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
			logger.Error("CrossChainRelayer", te.name, "GetTopElectBlockHeadByHeight error:", err)
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
			logger.Debug(">>>>> batch header")
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
			logger.Error("CrossChainRelayer", te.name, "EncodeHeaders failed:", err)
			return 0, 0, err
		}

		err = te.submitTopHeader(data, nonce)
		if err != nil {
			logger.Error("CrossChainRelayer", te.name, "submitHeaders failed:", err)
			return 0, 0, err
		}
	}

	return lastSubHeight, lastUnsubHeight, nil
}
