package eth2top

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/base"
	"toprelayer/contract/top/bridge"
	"toprelayer/msg"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wonderivan/logger"
)

const (
	METHOD_GETBRIDGESTATE        = "getbridgestate"
	SYNCHEADERS                  = "addLightClientBlock"
	BRIDIGESTATESUCCESS   string = "0x1"

	SUCCESSDELAY int64 = 15 //mainnet 120
	FATALTIMEOUT int64 = 24 //hours
	FORKDELAY    int64 = 5  //mainnet 10000 seconds
	ERRDELAY     int64 = 10
)

type Eth2TopRelayer struct {
	context.Context
	contract        common.Address
	chainId         uint64
	wallet          wallet.IWallet
	topsdk          *topsdk.TopSdk
	ethsdk          *ethsdk.EthSdk
	certaintyBlocks int
	subBatch        int
	verifyBlock     bool
}

func (et *Eth2TopRelayer) Init(topUrl, ethUrl, keypath, pass string, chainid uint64, contract common.Address, batch, cert int, verify bool) error {
	topsdk, err := topsdk.NewTopSdk(topUrl)
	if err != nil {
		return err
	}
	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
	if err != nil {
		return err
	}

	et.topsdk = topsdk
	et.ethsdk = ethsdk
	et.contract = contract
	et.chainId = chainid
	et.subBatch = batch
	et.certaintyBlocks = cert
	et.verifyBlock = verify

	w, err := wallet.NewWallet(topUrl, keypath, pass, chainid)
	if err != nil {
		return err
	}
	et.wallet = w
	return nil
}
func (et *Eth2TopRelayer) ChainId() uint64 {
	return et.chainId
}

func (et *Eth2TopRelayer) submitEthHeader(header []byte, nonce uint64) (*types.Transaction, error) {
	logger.Info("submitEthHeader length:%v,chainid:%v", len(header), et.chainId)
	gaspric, err := et.wallet.GasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	/* gaslimit, err := et.estimateGas(header)
	if err != nil {
		return nil, err
	} */

	//test mock
	gaslimit := uint64(300000)
	capfee := big.NewInt(0).SetUint64(gaspric.Uint64() * gaslimit * 2)
	logger.Info("account[%v] nonce:%v", et.wallet.CurrentAccount().Address, nonce)

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:  et.wallet.CurrentAccount().Address,
		Nonce: big.NewInt(0).SetUint64(nonce),
		//GasPrice: gaspric,
		GasLimit:  gaslimit,
		GasFeeCap: capfee,
		GasTipCap: big.NewInt(0),
		Signer:    et.signTransaction,
		Context:   context.Background(),
		NoSend:    false, //false: Send the transaction to the target chain by default; true: don't send
	}

	contractcaller, err := bridge.NewBridgeTransactor(et.contract, et.topsdk)
	if err != nil {
		return nil, err
	}

	sigTx, err := contractcaller.AddLightClientBlock(ops, header)
	if err != nil {
		logger.Error("Eth2TopRelayer AddLightClientBlock:%v", err)
		return nil, err
	}

	logger.Debug("hash:%v", sigTx.Hash())
	return sigTx, nil
}

//callback function to sign tx before send.
func (et *Eth2TopRelayer) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := et.wallet.CurrentAccount()
	if strings.EqualFold(acc.Address.Hex(), addr.Hex()) {
		stx, err := et.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("address:%v not available", addr)
}

func (et *Eth2TopRelayer) getTopBridgeState() (*msg.BridgeState, error) {
	hscaller, err := bridge.NewBridgeCaller(et.contract, et.topsdk)
	if err != nil {
		return nil, err
	}

	hscRaw := bridge.BridgeCallerRaw{Contract: hscaller}
	result := make([]interface{}, 1)

	err = hscRaw.Call(nil, &result, METHOD_GETBRIDGESTATE, et.chainId)
	if err != nil {
		return nil, err
	}
	state, success := result[0].(msg.BridgeState)
	if !success {
		return nil, err
	}
	return &state, nil
}

func (et *Eth2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayer... chainid:%v", et.chainId)
	defer wg.Done()

	timeoutDur := time.Duration(time.Second * 300) //test mock
	//timeoutDur := time.Duration(time.Hour * FATALTIMEOUT)
	timeout := time.NewTimer(timeoutDur)
	defer timeout.Stop()

	go func(timeoutDur time.Duration, timeout *time.Timer) {
		var syncStartHeight uint64 = 6000
		var delay time.Duration = time.Duration(1)
		for {
			time.Sleep(time.Second * delay)
			/* bridgeState, err := et.getTopBridgeState()
			if err != nil {
				logger.Error(err)
				delay = time.Duration(ERRDELAY)
				continue
			}
			if bridgeState.ConfirmState == CONFIRMSUCCESS {
				syncStartHeight = bridgeState.LatestSyncedHeight.Uint64() + 1
			} else {
				logger.Warn("top bridge confirm eth header failed,height:%v.", bridgeState.LatestConfirmedHeight.Uint64())
				syncStartHeight = bridgeState.LatestConfirmedHeight.Uint64()
			} */

			ethCurrentHeight, err := et.ethsdk.BlockNumber(context.Background())
			if err != nil {
				logger.Error(err)
				delay = time.Duration(ERRDELAY)
				continue
			}
			ethConfirmedBlockHeight := ethCurrentHeight - uint64(et.certaintyBlocks)
			if syncStartHeight <= ethConfirmedBlockHeight {
				hashes, err := et.signAndSendTransactions(syncStartHeight, ethConfirmedBlockHeight)
				if len(hashes) > 0 {
					logger.Info("Eth2TopRelayer sent block header from %v to :%v", syncStartHeight, ethConfirmedBlockHeight)
					delay = time.Duration(SUCCESSDELAY * int64(len(hashes)))
					timeout.Reset(timeoutDur)
					logger.Debug("timeout.Reset:%v", timeoutDur)
					syncStartHeight = ethConfirmedBlockHeight + 1 //test mock
					continue
				}
				if err != nil {
					logger.Error("Eth2TopRelayer signAndSendTransactions failed:%v", err)
					delay = time.Duration(ERRDELAY)
					continue
				}
			}
			//eth fork?
			logger.Warn("eth chain reverted?,syncStartHeight[%v] > ethConfirmedBlockHeight[%v]", syncStartHeight, ethConfirmedBlockHeight)
			delay = time.Duration(FORKDELAY)
		}
	}(timeoutDur, timeout)

	<-timeout.C
	logger.Error("relayer [%v] timeout.", et.chainId)
	return nil
}

func (et *Eth2TopRelayer) batch(headers []*types.Header, nonce uint64) (common.Hash, error) {
	logger.Info("batch headers number:", len(headers))
	if et.chainId == base.TOP && et.verifyBlock {
		for _, header := range headers {
			et.verifyBlocks(header)
		}
	}
	data, err := msg.EncodeHeaders(headers)
	if err != nil {
		logger.Error("Eth2TopRelayer EncodeHeaders failed:", err)
		return common.Hash{}, err
	}
	tx, err := et.submitEthHeader(data, nonce)
	if err != nil {
		logger.Error("Eth2TopRelayer submitHeaders failed:", err)
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

//test mock
//var nonce uint64 = 1

func (et *Eth2TopRelayer) signAndSendTransactions(lo, hi uint64) ([]common.Hash, error) {
	logger.Info("signAndSendTransactions height from:%v,to:%v", lo, hi)
	var batchHeaders []*types.Header
	var hashes []common.Hash
	nonce, err := et.wallet.GetNonce(et.wallet.CurrentAccount().Address)
	if err != nil {
		logger.Error(err)
		return hashes, err
	}
	h := lo
	for ; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			return hashes, err
		}
		batchHeaders = append(batchHeaders, header)
		if (h-lo+1)%uint64(et.subBatch) == 0 {
			hash, err := et.batch(batchHeaders, nonce)
			if err != nil {
				return hashes, err
			}
			batchHeaders = []*types.Header{}
			hashes = append(hashes, hash)
			nonce++
		}
	}
	if h > hi {
		if len(batchHeaders) > 0 {
			hash, err := et.batch(batchHeaders, nonce)
			if err != nil {
				return hashes, err
			}
			batchHeaders = []*types.Header{}
			hashes = append(hashes, hash)

			//test mock
			nonce++
		}
	}
	return hashes, nil
}

func (et *Eth2TopRelayer) verifyBlocks(header *types.Header) error {
	return nil
}

func (et *Eth2TopRelayer) estimateGas(gasprice *big.Int, headers []byte) (uint64, error) {
	/* rd := []byte{}
	ABI, err := abi.JSON(rd)
	if err != nil {
		return 0, err
	}

	data, err := ABI.Pack(SYNCHEADERS, headers)
	if err != nil {
		return 0, err
	} */

	callmsg := ethereum.CallMsg{
		From:     et.wallet.CurrentAccount().Address,
		To:       &et.contract,
		GasPrice: gasprice,
		Data:     []byte{},
	}

	return et.topsdk.EstimateGas(context.Background(), callmsg)
}
