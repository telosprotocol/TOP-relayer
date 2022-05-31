package eth2top

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/base"
	"toprelayer/contract/topbridge"
	"toprelayer/sdk/ethsdk"
	"toprelayer/sdk/topsdk"
	"toprelayer/util"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wonderivan/logger"
)

const (
	METHOD_GETBRIDGESTATE = "getCurrentBlockHeight"
	SYNCHEADERS           = "syncBlockHeader"

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
	abi             abi.ABI
}

func (et *Eth2TopRelayer) Init(topUrl, ethUrl, keypath, pass, abipath string, chainid uint64, contract common.Address, batch, cert int, verify bool) error {
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
	a, err := initABI(abipath)
	if err != nil {
		return err
	}
	et.abi = a
	return nil
}

func initABI(abifile string) (abi.ABI, error) {
	abidata, err := ioutil.ReadFile(abifile)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(strings.NewReader(string(abidata)))
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

	/* gaslimit, err := et.estimateGas(gaspric, header)
	if err != nil {
		logger.Error("estimateGas error:", err)
		return nil, err
	} */

	gaslimit := uint64(50000) //test mock
	capfee := big.NewInt(0).SetUint64(gaspric.Uint64() * gaslimit * 2)
	logger.Info("account[%v] nonce:%v,gaslimit:%v,capfee:%v", et.wallet.CurrentAccount().Address, nonce, gaslimit, capfee)

	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      et.wallet.CurrentAccount().Address,
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  gaslimit,
		GasFeeCap: capfee,
		GasTipCap: big.NewInt(0),
		Signer:    et.signTransaction,
		Context:   context.Background(),
		NoSend:    true, //false: Send the transaction to the target chain by default; true: don't send
	}

	contractcaller, err := topbridge.NewTopBridgeTransactor(et.contract, et.topsdk)
	if err != nil {
		return nil, err
	}

	sigTx, err := contractcaller.SyncBlockHeader(ops, header) //AddLightClientBlock(ops, header)
	if err != nil {
		logger.Error("Eth2TopRelayer AddLightClientBlock:%v", err)
		return nil, err
	}
	{
		byt, err := sigTx.MarshalBinary()
		if err != nil {
			logger.Error("MarshalBinary error:", err)
		}
		logger.Debug("rawtx:", hexutil.Encode(byt))
	}

	if ops.NoSend {
		err = util.VerifyEthSignature(sigTx)
		if err != nil {
			logger.Error("Eth2TopRelayer VerifyEthSignature error:", err)
			return nil, err
		}

		err := et.topsdk.SendTransaction(ops.Context, sigTx)
		if err != nil {
			logger.Error("Eth2TopRelayer SendTransaction error:", err)
			return nil, err
		}
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

func (et *Eth2TopRelayer) getTopBridgeCurrentHeight() (uint64, error) {
	input, err := et.abi.Pack(METHOD_GETBRIDGESTATE, et.chainId)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{
		From: et.wallet.CurrentAccount().Address,
		To:   &et.contract,
		Data: input,
	}
	ret, err := et.topsdk.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	return big.NewInt(0).SetBytes(ret).Uint64(), nil
}

func (et *Eth2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayer... chainid:%v", et.chainId)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDur := time.Duration(time.Second * 300) //test mock
		//timeoutDur := time.Duration(time.Hour * FATALTIMEOUT)
		timeout := time.NewTimer(timeoutDur)
		defer timeout.Stop()

		//var syncStartHeight uint64 = 6000
		var delay time.Duration = time.Duration(1)
		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				bridgeCurrentHeight, err := et.getTopBridgeCurrentHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				syncStartHeight := bridgeCurrentHeight + 1
				/* ethCurrentHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				ethConfirmedBlockHeight := ethCurrentHeight - uint64(et.certaintyBlocks) */
				var ethConfirmedBlockHeight uint64 = 1000
				if syncStartHeight <= ethConfirmedBlockHeight {
					hashes, err := et.signAndSendTransactions(syncStartHeight, ethConfirmedBlockHeight)
					if len(hashes) > 0 {
						if set := timeout.Reset(timeoutDur); !set {
							logger.Error("reset timeout falied!")
							delay = time.Duration(ERRDELAY)
							break
						}
						//syncStartHeight = ethConfirmedBlockHeight + 1 //test mock

						logger.Debug("timeout.Reset:%v", timeoutDur)
						logger.Info("Eth2TopRelayer sent block header from %v to :%v", syncStartHeight, ethConfirmedBlockHeight)

						delay = time.Duration(SUCCESSDELAY * int64(len(hashes)))
						break
					}
					if err != nil {
						logger.Error("Eth2TopRelayer signAndSendTransactions failed:%v", err)
						delay = time.Duration(ERRDELAY)
						break
					}
				}
				//eth fork?
				logger.Warn("eth chain reverted?,syncStartHeight[%v] > ethConfirmedBlockHeight[%v]", syncStartHeight, ethConfirmedBlockHeight)
				delay = time.Duration(FORKDELAY)
			}
		}
	}(done)

	<-done
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
	data, err := base.EncodeHeaders(headers)
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
		}
	}
	return hashes, nil
}

func (et *Eth2TopRelayer) verifyBlocks(header *types.Header) error {
	return nil
}

func (et *Eth2TopRelayer) estimateGas(gasprice *big.Int, data []byte) (uint64, error) {
	input, err := et.abi.Pack(SYNCHEADERS, data)
	if err != nil {
		return 0, err
	}

	capfee := big.NewInt(0).SetUint64(base.GetChainGasCapFee(et.chainId))
	callmsg := ethereum.CallMsg{
		From:      et.wallet.CurrentAccount().Address,
		To:        &et.contract,
		GasPrice:  gasprice,
		Gas:       500000,
		GasFeeCap: capfee,
		GasTipCap: nil,
		Data:      input,
	}

	return et.topsdk.EstimateGas(context.Background(), callmsg)
}
