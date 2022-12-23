package toprelayer

import (
	"context"
	"math/big"
	"sync"
	"time"
	"toprelayer/config"
	"toprelayer/contract/top/openallianceclient"
	rpc "toprelayer/rpc/openalliance_rpc"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	openAllianceBatchNum = 10

	ELECTION_BLOCK    = "election"
	AGGREGATE_BLOCK   = "aggregate"
	TRANSACTION_BLOCK = "transactions"

	// TODO
	sendFlag = 0x1
)

var (
	openAllianceClientSystemContract = common.HexToAddress("0x2769D9a843ba3830DEcCb15222d09559B441709C")
)

type OpenAlliance2TopRelayer struct {
	wallet          *wallet.Wallet
	openAllianceRpc *rpc.OpenAllianceRpc
	transactor      *openallianceclient.OpenAllianceClientTransactor
	callerSession   *openallianceclient.OpenAllianceClientCallerSession
}

func (relayer *OpenAlliance2TopRelayer) Init(cfg *config.Relayer, listenUrl []string, pass string) error {
	w, err := wallet.NewTopWallet(cfg.Url[0], cfg.KeyPath, pass)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	relayer.openAllianceRpc, err = rpc.NewOpenAllianceRpc(listenUrl[0])
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer ethclient.Dial error:", err)
		return err
	}

	topethlient, err := ethclient.Dial(cfg.Url[0])
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer new topethlient error:", err)
		return err
	}

	relayer.transactor, err = openallianceclient.NewOpenAllianceClientTransactor(openAllianceClientSystemContract, topethlient)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer NewEthClientTransactor error:", err)
		return err
	}

	relayer.callerSession = new(openallianceclient.OpenAllianceClientCallerSession)
	relayer.callerSession.Contract, err = openallianceclient.NewOpenAllianceClientCaller(openAllianceClientSystemContract, topethlient)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer NewEthClientCaller error:", err)
		return err
	}
	relayer.callerSession.CallOpts = bind.CallOpts{
		Pending:     false,
		From:        relayer.wallet.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}
	return nil
}

func (relayer *OpenAlliance2TopRelayer) signAndSendTransactions(lo, hi uint64) (uint64, uint64, error) {
	var lastSubHeight uint64 = 0
	var lastUnsubHeight uint64 = 0

	var batchHeaders [][]byte
	for h := lo; h <= hi; h++ {
		header, err := relayer.openAllianceRpc.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error("OpenAlliance2TopRelayer HeaderByNumber error:", err)
			break
		}
		logger.Debug("Top block, height: %v, type: %v, chainbits: %v", header.Number, header.BlockType, header.ChainBits)
		verify := false
		if header.BlockType == ELECTION_BLOCK {
			verify = true
		} else if header.BlockType == AGGREGATE_BLOCK {
			verify = true
		}
		if verify {
			logger.Debug(">>>>> send header")
			lastSubHeight = h
			batchHeaders = append(batchHeaders, common.Hex2Bytes(header.Header[2:]))
		} else {
			lastUnsubHeight = h
		}
	}
	if len(batchHeaders) > 0 {
		data, err := rlp.EncodeToBytes(batchHeaders)
		if err != nil {
			logger.Error("OpenAlliance2TopRelayer EncodeHeaders failed:", err)
			return 0, 0, err
		}
		err = relayer.submitEthHeader(data)
		if err != nil {
			logger.Error("OpenAlliance2TopRelayer submitHeaders failed:", err)
			return 0, 0, err
		}
	}
	return lastSubHeight, lastUnsubHeight, nil
}

func (relayer *OpenAlliance2TopRelayer) submitEthHeader(header []byte) error {
	logger.Info("CrossChainRelayer OpenAlliance2TopRelayer raw data:", common.Bytes2Hex(header))
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer GetNonce error:", err)
		return err
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer GasPrice error:", err)
		return err
	}
	packHeader, err := openallianceclient.PackSyncParam(header)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer PackSyncParam error:", err)
		return err
	}
	gaslimit, err := relayer.wallet.EstimateGas(context.Background(), &openAllianceClientSystemContract, packHeader)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer EstimateGas error:", err)
		return err
	}

	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  gaslimit,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	sigTx, err := relayer.transactor.AddLightClientBlocks(ops, header)
	if err != nil {
		logger.Error("OpenAlliance2TopRelayer sync error:", err)
		return err
	}
	logger.Info("OpenAlliance2TopRelayer tx info, account[%v] nonce:%v,capfee:%v,hash:%v,size:%v", relayer.wallet.Address(), nonce, gaspric, sigTx.Hash(), len(header))
	return nil
}

func (relayer *OpenAlliance2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	return relayer.wallet.SignTx(tx)
}

func (relayer *OpenAlliance2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start OpenAlliance2TopRelayer, subBatch: %v", openAllianceBatchNum)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Debug("OpenAlliance2TopRelayer set timeout: %v hours", FATALTIMEOUT)
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
				init, err := relayer.callerSession.Initialized()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if !init {
					logger.Info("OpenAlliance2TopRelayer not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				destHeight, err := relayer.callerSession.MaxMainHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("OpenAlliance2TopRelayer check dest top Height:", destHeight)
				// if destHeight == 0 {
				// 	if set := timeout.Reset(timeoutDuration); !set {
				// 		logger.Error("OpenAlliance2TopRelayer reset timeout falied!")
				// 		delay = time.Duration(ERRDELAY)
				// 		break
				// 	}
				// 	logger.Info("OpenAlliance2TopRelayer not init yet")
				// 	delay = time.Duration(ERRDELAY)
				// 	break
				// }
				srcHeight, err := relayer.openAllianceRpc.BlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("OpenAlliance2TopRelayer check src open alliance Height:", srcHeight)
				if lastSubHeight <= destHeight && destHeight < lastUnsubHeight {
					destHeight = lastUnsubHeight
				}
				if destHeight+1 > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("OpenAlliance2TopRelayer reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("OpenAlliance2TopRelayer waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}

				syncStartHeight := destHeight + 1
				syncNum := srcHeight - destHeight
				if syncNum > openAllianceBatchNum {
					syncNum = openAllianceBatchNum
				}
				syncEndHeight := syncStartHeight + syncNum - 1
				logger.Info("OpenAlliance2TopRelayer sync from %v to %v", syncStartHeight, syncEndHeight)

				subHeight, unsubHeight, err := relayer.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if err != nil {
					logger.Error("OpenAlliance2TopRelayer signAndSendTransactions failed:", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if subHeight > lastSubHeight {
					logger.Info("OpenAlliance2TopRelayer lastSubHeight: %v=>%v", lastSubHeight, subHeight)
					lastSubHeight = subHeight
				}
				if unsubHeight > lastUnsubHeight {
					logger.Info("OpenAlliance2TopRelayer lastUnsubHeight: %v=>%v", lastUnsubHeight, unsubHeight)
					lastUnsubHeight = unsubHeight
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("OpenAlliance2TopRelayer reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("OpenAlliance2TopRelayer sync round finish")
				if syncNum == openAllianceBatchNum {
					delay = time.Duration(SUCCESSDELAY)
				} else {
					delay = time.Duration(WAITDELAY)
				}
			}
		}
	}(done)

	<-done
	logger.Error("OpenAlliance2TopRelayer timeout")
	return nil
}

func (relayer *OpenAlliance2TopRelayer) GetInitData() ([]byte, error) {
	return nil, nil
}
