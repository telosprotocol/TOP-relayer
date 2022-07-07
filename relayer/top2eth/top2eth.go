package top2eth

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/contract/ethbridge"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/util"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	METHOD_GETHEIGHT = "maxMainHeight"
	METHOD_SYNC      = "addLightClientBlocks"

	ABI_PATH = "contract/ethbridge/ethbridge.abi"

	FATALTIMEOUT int64 = 24 //hours
	SUCCESSDELAY int64 = 60
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 1000

	CONFIRM_NUM int = 2
	BATCH_NUM   int = 20
)

type Top2EthRelayer struct {
	context.Context
	contract        common.Address
	chainId         uint64
	wallet          wallet.IWallet
	ethsdk          *ethsdk.EthSdk
	topsdk          *topsdk.TopSdk
	certaintyBlocks int
	subBatch        int
	abi             abi.ABI
}

func (te *Top2EthRelayer) Init(ethUrl, topUrl, keypath, pass string, chainid uint64, contract common.Address) error {
	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
	if err != nil {
		return err
	}
	topsdk, err := topsdk.NewTopSdk(topUrl)
	if err != nil {
		return err
	}
	te.topsdk = topsdk
	te.ethsdk = ethsdk
	te.contract = contract
	te.chainId = chainid
	te.subBatch = BATCH_NUM
	te.certaintyBlocks = CONFIRM_NUM

	w, err := wallet.NewWallet(ethUrl, keypath, pass, chainid)
	if err != nil {
		return err
	}
	te.wallet = w
	a, err := initABI(ABI_PATH)
	if err != nil {
		return err
	}
	te.abi = a
	return nil
}

func initABI(abifile string) (abi.ABI, error) {
	abidata, err := ioutil.ReadFile(abifile)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(strings.NewReader(string(abidata)))
}

func (te *Top2EthRelayer) ChainId() uint64 {
	return te.chainId
}

func (te *Top2EthRelayer) submitTopHeader(headers []byte, nonce uint64) error {
	logger.Info("Top2EthRelayer submitTopHeader length: %v,chainid: %v", len(headers), te.chainId)
	logger.Info("Top2EthRelayer raw data: %v", common.Bytes2Hex(headers))
	gaspric, err := te.wallet.GasPrice(context.Background())
	if err != nil {
		return err
	}

	gaslimit, err := te.estimateSyncGas(gaspric, headers)
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

	contractcaller, err := ethbridge.NewEthBridgeTransactor(te.contract, te.ethsdk)
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
	logger.Info("Start Top2EthRelayer relayer... chainid: %v, subBatch: %v certaintyBlocks: %v", te.chainId, te.subBatch, te.certaintyBlocks)
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
				toHeight, err := te.getEthBridgeCurrentHeight()
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
				if toHeight+1+uint64(te.certaintyBlocks) > fromHeight {
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
				limitEndHeight := fromHeight - uint64(te.certaintyBlocks)

				subHeight, unsubHeight, err := te.signAndSendTransactions(syncStartHeight, limitEndHeight, uint64(te.subBatch))
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

	var batch uint64 = 0
	h := lo
	for ; h <= hi; h++ {
		var flag bool
		header, flag, err := te.topsdk.GetTopElectBlockHeadByHeight(h)
		if err != nil {
			logger.Error(err)
			return 0, 0, err
		}
		if flag {
			batchHeaders = append(batchHeaders, header)
			lastSubHeight = h
			batch += 1
			if batch >= batchNum {
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
		// {
		// 	hd := common.Bytes2Hex(data)
		// 	logger.Debug("hex: ", hd)
		// }

		err = te.submitTopHeader(data, nonce)
		if err != nil {
			logger.Error("Top2EthRelayer submitHeaders failed:", err)
			return 0, 0, err
		}
	}

	return lastSubHeight, lastUnsubHeight, nil
}

func (te *Top2EthRelayer) getEthBridgeCurrentHeight() (uint64, error) {
	input, err := te.abi.Pack(METHOD_GETHEIGHT)
	if err != nil {
		logger.Error("Pack:", err)
		return 0, err
	}

	msg := ethereum.CallMsg{
		From: te.wallet.CurrentAccount().Address,
		To:   &te.contract,
		Data: input,
	}

	ret, err := te.ethsdk.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error("CallContract:", err)
		return 0, err
	}

	return big.NewInt(0).SetBytes(ret).Uint64(), nil
}

func (te *Top2EthRelayer) estimateSyncGas(price *big.Int, data []byte) (uint64, error) {
	input, err := te.abi.Pack(METHOD_SYNC, data)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{
		From:     te.wallet.CurrentAccount().Address,
		To:       &te.contract,
		GasPrice: price,
		Value:    big.NewInt(0),
		Data:     input,
	}

	return te.wallet.EstimateGas(context.Background(), msg)
}
