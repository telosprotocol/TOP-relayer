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
	METHOD_GETBRIDGESTATE        = "maxMainHeight"
	SYNCHEADERS                  = "addLightClientBlocks"
	CONFIRMSUCCESS        string = "0x1"

	SUCCESSDELAY int64 = 60 //mainnet 1000
	FATALTIMEOUT int64 = 24 //hours
	FORKDELAY    int64 = 5  //mainnet 3000 seconds
	ERRDELAY     int64 = 10
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

func (te *Top2EthRelayer) Init(ethUrl, topUrl, keypath, pass, abipath string, chainid uint64, contract common.Address, batch, cert int, verify bool) error {
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
	te.subBatch = batch
	te.certaintyBlocks = cert

	w, err := wallet.NewWallet(ethUrl, keypath, pass, chainid)
	if err != nil {
		return err
	}
	te.wallet = w
	a, err := initABI(abipath)
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

func (te *Top2EthRelayer) submitTopHeader(headers []byte, nonce uint64) (*types.Transaction, error) {
	logger.Info("Top2EthRelayer submitTopHeader length: %v,chainid: %v", len(headers), te.chainId)
	logger.Info("Top2EthRelayer raw data: %v", common.Bytes2Hex(headers))
	gaspric, err := te.wallet.GasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gaslimit, err := te.estimateGas(gaspric, headers)
	if err != nil {
		return nil, err
	}
	//test mock
	//gaslimit := uint64(500000)

	balance, err := te.wallet.GetBalance(te.wallet.CurrentAccount().Address)
	if err != nil {
		return nil, err
	}
	logger.Info("account[%v] balance:%v,nonce:%v,gasprice:%v,gaslimit:%v", te.wallet.CurrentAccount().Address, balance.Uint64(), nonce, gaspric.Uint64(), gaslimit)
	if balance.Uint64() <= gaspric.Uint64()*gaslimit {
		return nil, fmt.Errorf("account[%v] not sufficient funds,balance:%v", te.wallet.CurrentAccount().Address, balance.Uint64())
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
		return nil, err
	}

	sigTx, err := contractcaller.AddLightClientBlocks(ops, headers)
	if err != nil {
		logger.Error("Top2EthRelayer AddLightClientBlocks error:", err)
		return nil, err
	}

	if ops.NoSend {
		err = util.VerifyEthSignature(sigTx)
		if err != nil {
			logger.Error("Top2EthRelayer VerifyEthSignature error:", err)
			return nil, err
		}

		err := te.ethsdk.SendTransaction(ops.Context, sigTx)
		if err != nil {
			logger.Error("Top2EthRelayer SendTransaction error:", err)
			return nil, err
		}
	}
	logger.Info("hash:%v", sigTx.Hash())
	return sigTx, nil
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

func (te *Top2EthRelayer) getEthBridgeCurrentHeight() (uint64, error) {
	/* hscaller, err := hsc.NewHscCaller(te.contract, te.ethsdk)
	if err != nil {
		return nil, err
	}

	hscRaw := hsc.HscCallerRaw{Contract: hscaller}
	result := make([]interface{}, 1)
	err = hscRaw.Call(nil, &result, METHOD_GETBRIDGESTATE)
	if err != nil {
		return nil, err
	}

	state, success := result[0].(base.BridgeState)
	if !success {
		return nil, err
	} */

	input, err := te.abi.Pack(METHOD_GETBRIDGESTATE)
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

	// logger.Debug("getEthBridgeCurrentHeight height:", ret, common.Bytes2Hex(ret))

	return big.NewInt(0).SetBytes(ret).Uint64(), nil
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
				logger.Info("Top2EthRelayer to ethHeight: %v", toHeight)
				fromHeight, err := te.topsdk.GetLatestTopElectBlockHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Top2EthRelayer from topHeight: %v", fromHeight)

				if toHeight+1+uint64(te.certaintyBlocks) > fromHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("height not satisfied, delay")
					delay = time.Duration(ERRDELAY)
					break
				}
				syncStartHeight := toHeight + 1
				syncNum := fromHeight - uint64(te.certaintyBlocks) - toHeight
				if syncNum > uint64(te.subBatch) {
					syncNum = uint64(te.subBatch)
				}
				syncEndHeight := syncStartHeight + syncNum - 1

				hashes, err := te.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if len(hashes) > 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("Top2EthRelayer sent block header from %v to %v success", syncStartHeight, syncEndHeight)
					delay = time.Duration(SUCCESSDELAY)
					break
				}
				if err != nil {
					logger.Error("Top2EthRelayer signAndSendTransactions failed: %v", err)
					delay = time.Duration(ERRDELAY)
				}
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout.", te.chainId)
	return nil
}

func (te *Top2EthRelayer) batch(headers [][]byte, nonce uint64) (common.Hash, error) {
	// logger.Info("batch headers number:", len(headers))
	// data := bytes.Join(headers, []byte{})
	data, err := rlp.EncodeToBytes(headers)
	if err != nil {
		logger.Error("Eth2TopRelayer EncodeHeaders failed:", err)
		return common.Hash{}, err
	}
	// {
	// 	hd := common.Bytes2Hex(data)
	// 	logger.Debug("hex: ", hd)
	// }

	tx, err := te.submitTopHeader(data, nonce)
	if err != nil {
		logger.Error("Top2EthRelayer submitHeaders failed:", err)
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (te *Top2EthRelayer) signAndSendTransactions(lo, hi uint64) ([]common.Hash, error) {
	var batchHeaders [][]byte
	var hashes []common.Hash
	nonce, err := te.wallet.GetNonce(te.wallet.CurrentAccount().Address)
	if err != nil {
		return hashes, err
	}

	h := lo
	for ; h <= hi; h++ {
		header, err := te.topsdk.GetTopElectBlockHeadByHeight(h)
		if err != nil {
			logger.Error(err)
			return hashes, err
		}

		batchHeaders = append(batchHeaders, header)
		if (h-lo+1)%uint64(te.subBatch) == 0 {
			hash, err := te.batch(batchHeaders, nonce)
			if err != nil {
				return hashes, err
			}
			batchHeaders = [][]byte{}
			hashes = append(hashes, hash)
			nonce++
		}
		time.Sleep(time.Second * 1)
	}
	if h > hi {
		if len(batchHeaders) > 0 {
			hash, err := te.batch(batchHeaders, nonce)
			if err != nil {
				return hashes, err
			}
			batchHeaders = [][]byte{}
			hashes = append(hashes, hash)
		}
	}

	return hashes, nil
}

func (te *Top2EthRelayer) estimateGas(price *big.Int, data []byte) (uint64, error) {
	input, err := te.abi.Pack(SYNCHEADERS, data)
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
