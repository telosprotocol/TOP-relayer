package crosschainrelayer

import (
	"container/list"
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	"toprelayer/contract/eth/topclient"
	"toprelayer/relayer/monitor"
	top "toprelayer/util"
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
	SUCCESSDELAY int64 = 60
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 60

	ELECTION_BLOCK    = "election"
	AGGREGATE_BLOCK   = "aggregate"
	TRANSACTION_BLOCK = "transactions"
)

var (
	sendFlag = map[string]uint64{
		config.ETH_CHAIN:     0x1,
		config.BSC_CHAIN:     0x2,
		config.HECO_CHAIN:    0x4,
		config.OPEN_ALLIANCE: 0x8}
)

type VerifyInfo struct {
	Block      *top.TopHeader
	VerifyList []string
}

type VerifyResp struct {
	Code       int32  `json:"code"`
	Logno      string `json:"logno"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	Result     bool   `json:"result"`
	Servertime string `json:"servertime"`
}

type CrossChainRelayer struct {
	name       string
	contract   common.Address
	wallet     *wallet.Wallet
	transactor *topclient.TopClientTransactor
	caller     *topclient.TopClientCaller
	monitor    *monitor.Monitor
	blockList  *list.List
}

func (te *CrossChainRelayer) Init(chainName string, cfg *config.Relayer, listenUrl string, pass string) error {
	te.name = chainName

	if cfg.Contract == "" {
		logger.Error("CrossChainRelayer", te.name, "contract nil:", cfg.Contract)
		return fmt.Errorf("contract error")
	}
	te.contract = common.HexToAddress(cfg.Contract)

	w, err := wallet.NewEthWallet(cfg.Url[0], listenUrl, cfg.KeyPath, pass)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "NewWallet error:", err)
		return err
	}
	te.wallet = w

	ethsdk, err := ethclient.Dial(cfg.Url[0])
	if err != nil {
		return err
	}
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
	te.monitor, err = monitor.New(te.wallet.Address(), cfg.Url[0])
	if err != nil {
		logger.Error("TopRelayer from", te.name, "New monitor error:", err)
		return err
	}

	te.blockList = list.New()

	logger.Info(te)
	return nil
}

func (te *CrossChainRelayer) submitTopHeader(headers []byte) error {
	logger.Info("CrossChainRelayer", te.name, "raw data:", common.Bytes2Hex(headers))
	nonce, err := te.wallet.NonceAt(context.Background(), te.wallet.Address(), nil)
	if err != nil {
		return err
	}
	gaspric, err := te.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "GasPrice error:", err)
		return err
	}
	packHeaders, err := topclient.PackSyncParam(headers)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "PackSyncParam error:", err)
		return err
	}
	gaslimit, err := te.wallet.EstimateGas(context.Background(), &te.contract, packHeaders)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "EstimateGas error:", err)
		return err
	}
	//test mock
	//gaslimit := uint64(500000)

	balance, err := te.wallet.BalanceAt(context.Background(), te.wallet.Address(), nil)
	if err != nil {
		return err
	}
	if balance.Uint64() <= gaspric.Uint64()*gaslimit {
		return fmt.Errorf("CrossChainRelayer %v account[%v] balance not enough:%v", te.name, te.wallet.Address(), balance.Uint64())
	}

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:     te.wallet.Address(),
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
	te.monitor.AddTx(sigTx.Hash())
	logger.Info("CrossChainRelayer %v tx info, account[%v] balance:%v,nonce:%v,gasprice:%v,gaslimit:%v,length:%v,hash:%v", te.name, te.wallet.Address(), balance.Uint64(), nonce, gaspric.Uint64(), gaslimit, len(headers), sigTx.Hash())
	return nil
}

// callback function to sign tx before send.
func (te *CrossChainRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := te.wallet.Address()
	if strings.EqualFold(acc.Hex(), addr.Hex()) {
		stx, err := te.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("address:%v not available", addr)
}

func (te *CrossChainRelayer) queryBlocks(lo, hi uint64) (uint64, uint64, error) {
	var lastSubHeight uint64 = 0
	var lastUnsubHeight uint64 = 0

	flag := sendFlag[te.name]
	for h := lo; h <= hi; h++ {
		block, err := te.wallet.TopHeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error("CrossChainRelayer", te.name, "GetTopElectBlockHeadByHeight error:", err)
			break
		}
		logger.Debug("Top block, height: %v, type: %v, chainbits: %v", block.Number, block.BlockType, block.ChainBits)
		submit := false
		if block.BlockType == ELECTION_BLOCK {
			submit = true
		} else if block.BlockType == AGGREGATE_BLOCK {
			blockFlag, err := strconv.ParseInt(block.ChainBits, 0, 64)
			if err != nil {
				logger.Error("ParseInt error:", err)
				break
			}
			if blockFlag == 0 || int64(flag)&blockFlag > 0 {
				submit = true
			}
		}
		if submit {
			logger.Debug(">>>>> submit header")
			lastSubHeight = h
			te.blockList.PushBack(*block)
			break
		} else {
			lastUnsubHeight = h
		}
	}

	return lastSubHeight, lastUnsubHeight, nil
}

func (te *CrossChainRelayer) verifyAndSendTransaction() {
	if te.blockList.Len() == 0 {
		return
	}
	element := te.blockList.Front()
	if element == nil {
		logger.Error("txList get front nil")
		return
	}
	header, ok := element.Value.(top.TopHeader)
	if !ok {
		logger.Error("txList get front error")
		return
	}
	if res := doWithHeader(header); res == false {
		logger.Info("do with header not ok:", header.Hash)
		return
	}
	logger.Info("do with header ok:", header.Hash)

	var batchHeaders [][]byte
	batchHeaders = append(batchHeaders, common.Hex2Bytes(header.Header[2:]))
	data, err := rlp.EncodeToBytes(batchHeaders)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "EncodeHeaders failed:", err)
		return
	}
	err = te.submitTopHeader(data)
	if err != nil {
		logger.Error("CrossChainRelayer", te.name, "submitHeaders failed:", err)
		return
	}

	// clear list
	te.blockList.Remove(element)
}

func (te *CrossChainRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start CrossChainRelayer %v...", te.name)
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
					From:        te.wallet.Address(),
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
				if te.blockList.Len() > 0 {
					logger.Debug("CrossChainRelayer", te.name, "find block to verify")
					te.verifyAndSendTransaction()
					delay = time.Duration(WAITDELAY)
					break
				}
				fromHeight, err := te.wallet.TopBlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("CrossChainRelayer", te.name, "src top Height:", fromHeight)

				if lastSubHeight <= toHeight && toHeight < lastUnsubHeight {
					toHeight = lastUnsubHeight
				}
				if toHeight+1 > fromHeight {
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
				limitEndHeight := fromHeight

				subHeight, unsubHeight, err := te.queryBlocks(syncStartHeight, limitEndHeight)
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
				delay = time.Duration(SUCCESSDELAY)
				break
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout", te.name)
	return nil
}
