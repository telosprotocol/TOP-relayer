package eth2top

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"
	"toprelayer/contract/topbridge"
	"toprelayer/relayer/eth2top/ethashapp"
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
	METHOD_GETHEIGHT = "get_height"
	METHOD_SYNC      = "sync"

	ABI_PATH = "contract/topbridge/topbridge.abi"

	FATALTIMEOUT int64 = 24 //hours
	SUCCESSDELAY int64 = 10
	ERRDELAY     int64 = 10
	WAITDELAY    int64 = 60

	CONFIRM_NUM int = 25
	BATCH_NUM   int = 5

	BLOCKS_PER_EPOCH uint64 = 30000
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
	abi             abi.ABI
}

func (et *Eth2TopRelayer) Init(topUrl, ethUrl, keypath, pass string, chainid uint64, contract common.Address) error {
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
	et.subBatch = BATCH_NUM
	et.certaintyBlocks = CONFIRM_NUM

	w, err := wallet.NewWallet(topUrl, keypath, pass, chainid)
	if err != nil {
		return err
	}
	et.wallet = w
	a, err := initABI(ABI_PATH)
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

func (et *Eth2TopRelayer) submitEthHeader(header []byte, nonce uint64) error {
	gaspric, err := et.wallet.GasPrice(context.Background())
	if err != nil {
		return err
	}

	gaslimit, err := et.estimateSyncGas(gaspric, header)
	if err != nil {
		logger.Error("estimateGas error:", err)
		return err
	}

	capfee := big.NewInt(0).SetUint64(gaspric.Uint64())

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
		return err
	}

	sigTx, err := contractcaller.Sync(ops, header) //AddLightClientBlock(ops, header)
	if err != nil {
		logger.Error("Eth2TopRelayer AddLightClientBlock:%v", err)
		return err
	}

	if ops.NoSend {
		err = util.VerifyEthSignature(sigTx)
		if err != nil {
			logger.Error("Eth2TopRelayer VerifyEthSignature error:", err)
			return err
		}

		err := et.topsdk.SendTransaction(ops.Context, sigTx)
		if err != nil {
			logger.Error("Eth2TopRelayer SendTransaction error:", err)
			return err
		}
	}
	logger.Info("tx info, account[%v] nonce:%v,gaslimit:%v,capfee:%v,hash:%v,size:%v", et.wallet.CurrentAccount().Address, nonce, gaslimit, capfee, sigTx.Hash(), len(header))
	return nil
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

func (et *Eth2TopRelayer) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayer relayer... chainid: %v, subBatch: %v certaintyBlocks: %v", et.chainId, et.subBatch, et.certaintyBlocks)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Info("Eth2TopRelayer set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				destHeight, err := et.getTopBridgeCurrentHeight()
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check dest top Height: %v", destHeight)
				if destHeight == 0 {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("top eth-client not init yet")
					delay = time.Duration(ERRDELAY)
					break
				}
				srcHeight, err := et.ethsdk.BlockNumber(context.Background())
				if err != nil {
					logger.Error(err)
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer check src eth Height: %v", srcHeight)

				if destHeight+1+uint64(et.certaintyBlocks) > srcHeight {
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Debug("waiting src eth update, delay")
					delay = time.Duration(WAITDELAY)
					break
				}

				syncStartHeight := destHeight + 1
				syncNum := srcHeight - uint64(et.certaintyBlocks) - destHeight
				if syncNum > uint64(et.subBatch) {
					syncNum = uint64(et.subBatch)
				}
				syncEndHeight := syncStartHeight + syncNum - 1
				logger.Info("Eth2TopRelayer sync from %v to %v", syncStartHeight, syncEndHeight)

				err = et.signAndSendTransactions(syncStartHeight, syncEndHeight)
				if err != nil {
					logger.Error("Eth2TopRelayer signAndSendTransactions failed:%v", err)
					delay = time.Duration(ERRDELAY)
					break
				}
				if set := timeout.Reset(timeoutDuration); !set {
					logger.Error("reset timeout falied!")
					delay = time.Duration(ERRDELAY)
					break
				}
				logger.Info("Eth2TopRelayer sync round finish")
				delay = time.Duration(SUCCESSDELAY)
				break
			}
		}
	}(done)

	<-done
	logger.Error("relayer [%v] timeout.", et.chainId)
	return nil
}

func (et *Eth2TopRelayer) signAndSendTransactions(lo, hi uint64) error {
	var batch []byte
	nonce, err := et.wallet.GetNonce(et.wallet.CurrentAccount().Address)
	if err != nil {
		logger.Error(err)
		return err
	}

	for h := lo; h <= hi; h++ {
		header, err := et.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			logger.Error(err)
			return err
		}
		ethashproof, err := ethashapp.EthashWithProofs(h, header)
		if err != nil {
			logger.Error(err)
			return err
		}
		rlp_bytes, err := rlp.EncodeToBytes(ethashproof)
		if err != nil {
			logger.Fatal("rlp encode error: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}

	// maybe verify block
	// if et.chainId == topChainId {
	// 	for _, header := range headers {
	// 		et.verifyBlocks(header)
	// 	}
	// }
	err = et.submitEthHeader(batch, nonce)
	if err != nil {
		logger.Error("Eth2TopRelayer submitHeaders failed:", err)
		return err
	}

	return nil
}

func (et *Eth2TopRelayer) verifyBlocks(header *types.Header) error {
	return nil
}

func (et *Eth2TopRelayer) getTopBridgeCurrentHeight() (uint64, error) {
	input, err := et.abi.Pack(METHOD_GETHEIGHT)
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

func (et *Eth2TopRelayer) estimateSyncGas(gasprice *big.Int, data []byte) (uint64, error) {
	input, err := et.abi.Pack(METHOD_SYNC, data)
	if err != nil {
		return 0, err
	}

	callmsg := ethereum.CallMsg{
		From:      et.wallet.CurrentAccount().Address,
		To:        &et.contract,
		GasPrice:  gasprice,
		Gas:       0,
		GasFeeCap: nil,
		GasTipCap: nil,
		Data:      input,
	}

	return et.topsdk.EstimateGas(context.Background(), callmsg)
}
