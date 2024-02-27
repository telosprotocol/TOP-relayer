package toprelayer

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/prysmaticlabs/prysm/v4/api/client/beacon"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/container/slice"
	"github.com/prysmaticlabs/prysm/v4/math"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	eth2bridge "toprelayer/contract/top/eth2client"
	beaconrpc "toprelayer/rpc/ethbeacon_rpc"
	"toprelayer/rpc/ethereum"
	"toprelayer/rpc/ethereum/light_client"
	"toprelayer/wallet"
)

var (
	eth2ClientSystemContract = common.HexToAddress("0xff00000000000000000000000000000000000009")
)

type ClientModeEnum uint8

const (
	Invalid ClientModeEnum = iota
	SubmitLightClientUpdateMode
	SubmitHeaderMode
)

type Eth2TopRelayerV2 struct {
	wallet       *wallet.Wallet
	ethrpcclient *ethclient.Client
	// beaconrpcclient *beaconrpc.BeaconGrpcClient
	beaconClient  *ethereum.BeaconChainClient
	transactor    *eth2bridge.Eth2ClientTransactor
	callerSession *eth2bridge.Eth2ClientCallerSession
	lastSlot      uint64
}

func (relayer *Eth2TopRelayerV2) Init(cfg *config.Relayer, listenUrl []string, pass string) error {
	w, err := wallet.NewTopWallet(cfg.Url[0], cfg.KeyPath, pass)
	logger.Info("Eth2TopRelayerV2 TOP wallet url:", cfg.Url[0])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	if len(listenUrl) < 2 {
		err := errors.New("listenUrl num error")
		logger.Error("Eth2TopRelayerV2 listenUrl error:", err)
		return err
	}
	relayer.ethrpcclient, err = ethclient.Dial(listenUrl[0])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 ethclient.Dial error:", err)
		return err
	}
	relayer.beaconClient, err = ethereum.NewBeaconChainClient(listenUrl[1])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 NewBeaconGrpcClient error:", err)
		return err
	}
	topethlient, err := ethclient.Dial(cfg.Url[0])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 new topethlient error:", err)
		return err
	}

	relayer.transactor, err = eth2bridge.NewEth2ClientTransactor(eth2ClientSystemContract, topethlient)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 NewEthClientTransactor error:", err)
		return err
	}

	relayer.callerSession = new(eth2bridge.Eth2ClientCallerSession)
	relayer.callerSession.Contract, err = eth2bridge.NewEth2ClientCaller(eth2ClientSystemContract, topethlient)
	if err != nil {
		logger.Error("Eth2 NewEthClientCaller error:", err)
		return err
	}
	relayer.callerSession.CallOpts = bind.CallOpts{
		Pending:     false,
		From:        relayer.wallet.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}
	relayer.lastSlot = 0
	if err != nil {
		logger.Error("Eth2TopRelayerV2 New monitor error", err)
		return err
	}

	logger.Info("Eth2TopRelayerV2 init eth1 url:%s eth2 url:%s TOP eth1 url:%s ethHeaderSyncContractAddr:%s", listenUrl[0], listenUrl[1], cfg.Url[0], eth2ClientSystemContract.Hex())
	return nil
}

func (relayer *Eth2TopRelayerV2) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayerV2")
	go func() {
		defer wg.Done()
		sleepTime := time.Duration(10)
		for {
			logger.Info("Eth2TopRelayerV2 ===================== New Cycle Start =====================")
			time.Sleep(time.Second * sleepTime)
			if ok, err := relayer.callerSession.Initialized(); err != nil || !ok {
				logger.Error("Eth2TopRelayerV2 don't initialize the header,error:", err)
				sleepTime = time.Duration(30)
				continue
			}
			mode, err := relayer.callerSession.GetClientMode()
			if err != nil {
				logger.Error("Eth2TopRelayerV2 get client mode,error:", err.Error())
				sleepTime = time.Duration(30)
				continue
			}
			logger.Info("Eth2TopRelayerV2 ClientMode(%d) Start", mode)
			if ClientModeEnum(mode) == SubmitLightClientUpdateMode {
				err = relayer.sendLightClientUpdatesWithChecks()
			} else if ClientModeEnum(mode) == SubmitHeaderMode {
				err = relayer.submitHeaders()
			} else {
				logger.Error("Eth2TopRelayerV2 Invalid ClientMode(%v):", mode)
				sleepTime = time.Duration(30)
				continue
			}
			if err != nil {
				if strings.Contains(err.Error(), "no need to submit") {
					logger.Info("Eth2TopRelayerV2 ClientMode(%d), %s", mode, err)
				} else {
					logger.Error("Eth2TopRelayerV2 ClientMode(%d) Submit Fail,%s", mode, err)
				}
				sleepTime = time.Duration(30)
				continue
			}
			sleepTime = time.Duration(30)
			logger.Info("Eth2TopRelayerV2 ClientMode(%d) Success", mode)
		}
	}()
	return nil
}

//func (relayer *Eth2TopRelayerV2) blockKnownOnTop(slot uint64) (bool, error) {
//	height, err := relayer.beaconClient.GetBlockNumberForSlot(primitives.Slot(slot))
//	//logger.Debug("blockKnownOnTop slot %v, height %v", slot, height)
//	if err != nil {
//		return false, err
//	}
//	return relayer.callerSession.IsKnownExecutionHeader(height)
//}

//func (relayer *Eth2TopRelayerV2) findLeftNonErrorSlot(leftSlot, rightSlot uint64) (uint64, bool) {
//	slot := leftSlot
//	for slot != rightSlot {
//		known, err := relayer.blockKnownOnTop(slot)
//		if err != nil {
//			slot += 1
//		} else {
//			return slot, known
//		}
//	}
//	return slot, false
//}

//func (relayer *Eth2TopRelayerV2) linearSearchForward(slot, maxSlot uint64) (uint64, error) {
//	for {
//		if slot >= maxSlot {
//			break
//		}
//		known, err := relayer.blockKnownOnTop(slot + 1)
//		if err != nil {
//			if ethereum.IsErrorNoBlockForSlot(err) {
//				slot += 1
//				continue
//			} else {
//				logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
//				return 0, err
//			}
//		}
//		if known {
//			slot += 1
//			logger.Debug("curr top known slot: %v,maxSlot: %v", slot, maxSlot)
//		} else {
//			break
//		}
//	}
//	logger.Debug("linearSearchForward return slot: %v", slot)
//	return slot, nil
//}

//func (relayer *Eth2TopRelayerV2) linearSearchBackward(startSlot, lastSlot uint64) (uint64, error) {
//	slot := lastSlot
//	lastFalseSlot := slot + 1
//	for {
//		if slot <= startSlot {
//			break
//		}
//		known, err := relayer.blockKnownOnTop(slot)
//		if err != nil {
//			if beaconrpc.IsErrorNoBlockForSlot(err) {
//				slot -= 1
//				continue
//			} else {
//				logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
//				return 0, err
//			}
//		}
//		if known {
//			break
//		} else {
//			lastFalseSlot = slot
//			slot -= 1
//		}
//	}
//	return lastFalseSlot - 1, nil
//}

//func (relayer *Eth2TopRelayerV2) linerSlotSearch(slot, finalizedSlot, lastEthSlot uint64) (uint64, error) {
//	if slot == finalizedSlot {
//		logger.Info("slot equal finalizedSlot %v go forward,greater than lastEthSlot:%v", slot, lastEthSlot)
//		return relayer.linearSearchForward(slot, lastEthSlot)
//	}
//	known, err := relayer.blockKnownOnTop(slot)
//	if err != nil {
//		if beaconrpc.IsErrorNoBlockForSlot(err) {
//			leftSlot, known := relayer.findLeftNonErrorSlot(slot+1, lastEthSlot+1)
//			if known {
//				return relayer.linearSearchForward(leftSlot, lastEthSlot)
//			} else {
//				return relayer.linearSearchForward(finalizedSlot, leftSlot-1)
//			}
//		} else {
//			logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
//			return 0, err
//		}
//	}
//	if known {
//		logger.Debug("slot %v known, go forward", slot)
//		return relayer.linearSearchForward(slot, lastEthSlot)
//	} else {
//		logger.Debug("slot %v unknown, go backward", slot)
//		return relayer.linearSearchBackward(finalizedSlot, slot)
//	}
//}

//func (relayer *Eth2TopRelayerV2) getMaxSlotForSubmission() (uint64, error) {
//	slot, err := relayer.beaconClient.GetLastSlotNumber()
//	if err != nil {
//		return 0, err
//	}
//	return uint64(slot), nil
//}

//func (relayer *Eth2TopRelayerV2) getLastEth2SlotOnTop(lastEthSlot uint64) (uint64, error) {
//	finalizedSlot, err := relayer.callerSession.FinalizedBeaconBlockSlot()
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 FinalizedBeaconBlockSlot error", err)
//		return 0, nil
//	}
//	lastSubmittedSlot := relayer.lastSlot
//	slot := finalizedSlot
//	if lastSubmittedSlot > finalizedSlot {
//		slot = lastSubmittedSlot
//	}
//	logger.Debug("getLastEth2SlotOnTop finalizedSlot: %v, lastSubmittedSlot: %v, slot: %v", finalizedSlot, lastSubmittedSlot, slot)
//	return relayer.linerSlotSearch(slot, finalizedSlot, lastEthSlot)
//}

func (relayer *Eth2TopRelayerV2) getLastFinalizedSlotOnTop() (primitives.Slot, error) {
	slotVal, err := relayer.callerSession.FinalizedBeaconBlockSlot()
	if err != nil {
		return 0, err
	}
	return primitives.Slot(slotVal), nil
}

func (relayer *Eth2TopRelayerV2) getLastFinalizedSlotOnEth() (primitives.Slot, error) {
	slotVal, err := relayer.beaconClient.GetLastFinalizedSlotNumber()
	if err != nil {
		return 0, err
	}
	return slotVal, nil
}

func (relayer *Eth2TopRelayerV2) sendRegularLightClientUpdate(lastFinalizedTopSlot, lastFinalizedEthSlot primitives.Slot) error {
	lastPeriodOnTOP, lastPeriodOnEth := ethereum.GetPeriodForSlot(lastFinalizedTopSlot), ethereum.GetPeriodForSlot(lastFinalizedEthSlot)
	var data *light_client.LightClientUpdate
	var err error
	if lastPeriodOnTOP == lastPeriodOnEth {
		data, err = relayer.beaconClient.GetFinalizedLightClientUpdateByEthSlot(lastFinalizedEthSlot)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetLightClientUpdate at same period error:", err)
			return err
		}
	} else if lastPeriodOnTOP+1 == lastPeriodOnEth {
		data, err = relayer.beaconClient.GetLastFinalizedLightClientUpdateV2WithNextSyncCommitteeByEthSlot(lastFinalizedEthSlot)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetLightClientUpdate at next period error:", err)
			return err
		}
	} else {
		data, err = relayer.beaconClient.GetLightClientUpdateV2(lastPeriodOnTOP + 1)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetLightClientUpdate at near period error:", err)
			return err
		}
	}
	bytes, err := data.Encode()
	if err != nil {
		logger.Error("EncodeToBytes error:", err)
		return err
	}
	return relayer.submitLightClientUpdate(bytes)
}

func (relayer *Eth2TopRelayerV2) sendLightClientUpdatesWithChecks() error {
	lastFinalizedSlotOnTop, err := relayer.getLastFinalizedSlotOnTop()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getLastFinalizedSlotOnTop error:", err)
		return err
	}
	lastFinalizedSlotOnEth, err := relayer.getLastFinalizedSlotOnEth()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getLastFinalizedSlotOnEth error:", err)
		return err
	}
	lastPeriodOnTOP, lastPeriodOnEth := ethereum.GetPeriodForSlot(lastFinalizedSlotOnTop), ethereum.GetPeriodForSlot(lastFinalizedSlotOnEth)
	logger.Info("Eth2TopRelayerV2 lastFinalizedSlot TOP(Period:%d,slot:%d), ETH(period:%d,slot:%d)", lastPeriodOnTOP, lastFinalizedSlotOnTop, lastPeriodOnEth, lastFinalizedSlotOnEth)
	if !relayer.isEnoughBlocksForLightClientUpdate(lastFinalizedSlotOnTop, lastFinalizedSlotOnEth) {
		return errors.New("no need to submit LightClientUpdate")
	}
	if err = relayer.sendRegularLightClientUpdate(lastFinalizedSlotOnTop, lastFinalizedSlotOnEth); err != nil {
		logger.Error("Eth2TopRelayerV2 sendLightClientUpdates error:", err)
		return err
	}
	return nil
}

func (relayer *Eth2TopRelayerV2) txOption(packData []byte) (*bind.TransactOpts, error) {
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonce error:", err)
		return nil, err
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GasPrice error:", err)
		return nil, err
	}
	gaslimit, err := relayer.wallet.EstimateGas(context.Background(), &eth2ClientSystemContract, packData)
	if err != nil {
		logger.Error(fmt.Sprintf("Eth2TopRelayer EstimateGas error:%s, data len:%v", err, len(packData)))
		return nil, err
	}
	logger.Info("Eth2TopRelayer tx option info, account[%v] nonce:%v,capfee:%v", relayer.wallet.Address(), nonce, gaspric)
	return &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  gaslimit,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}, nil
}

func (relayer *Eth2TopRelayerV2) submitHeaders() error {
	lastFinalizedBlockNumber, err := relayer.callerSession.LastBlockNumber()
	if err != nil {
		return err
	}
	currentBlockNumber, err := relayer.getMaxBlockNumber()
	if err != nil {
		return err
	}
	logger.Info("Eth2TopRelayerV2 lastBlockNumberOnTop Finalized:%v, Tail:%v", lastFinalizedBlockNumber, currentBlockNumber+1)
	minBlockNumberInBatch := math.Max(lastFinalizedBlockNumber+1, currentBlockNumber-beaconrpc.HEADER_BATCH_SIZE+1)
	headers, err := relayer.getExecutionBlocksBetweenByNumber(minBlockNumberInBatch, currentBlockNumber)
	if err != nil {
		return err
	}
	logger.Info("Eth2TopRelayerV2 submitHeaders len: %v", len(headers))
	slice.Reverse(headers)
	if err = relayer.submitEthHeader(headers); err != nil {
		return err
	}
	return nil
}

func (relayer *Eth2TopRelayerV2) getMaxBlockNumber() (uint64, error) {
	number, err := relayer.callerSession.GetUnfinalizedTailBlockNumber()
	if err != nil {
		return 0, err
	}
	if number > 0 {
		return number - 1, nil
	}
	if slot, err := relayer.getLastFinalizedSlotOnTop(); err != nil {
		return 0, err
	} else {
		return relayer.beaconClient.GetBlockNumberForSlot(slot)
	}
}

func (relayer *Eth2TopRelayerV2) submitEthHeader(headers []*types.Header) error {
	if len(headers) == 0 {
		return errors.New("submitEthHeader headers is nil")
	}
	encodeHeaders, err := relayer.encodeEthHeaders(headers)
	if err != nil {
		return err
	}
	packHeader, err := eth2bridge.PackSubmitExecutionHeaderParam(encodeHeaders)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 PackSubmitExecutionHeaderParam error:", err)
		return err
	}
	ops, err := relayer.txOption(packHeader)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 txOption error:", err)
		return err
	}
	sigTx, err := relayer.transactor.SubmitExecutionHeaders(ops, encodeHeaders)
	if err != nil {
		logger.Error("Eth2TopRelayer sync error:", err)
		return err
	}
	logger.Info("Eth2TopRelayer submitEthHeader tx info, account[%v] txHash:%v,size:%v", relayer.wallet.Address(), sigTx.Hash(), len(headers))
	return nil
}

func (relayer *Eth2TopRelayerV2) submitLightClientUpdate(update []byte) error {
	packUpdate, err := eth2bridge.PackSubmitBeaconChainLightClientUpdateParam(update)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 PackSubmitBeaconChainLightClientUpdateParam error:", err)
		return err
	}
	ops, err := relayer.txOption(packUpdate)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 txOption error:", err)
		return err
	}
	logger.Info("Eth2TopRelayer submitLightClientUpdate len:", len(update))
	sigTx, err := relayer.transactor.SubmitBeaconChainLightClientUpdate(ops, update)
	if err != nil {
		logger.Error("Eth2TopRelayer SubmitBeaconChainLightClientUpdate error:", err)
		return err
	}
	logger.Info("Eth2TopRelayer submitLightClientUpdate tx info, account[%v] hash:%v,size:%v", relayer.wallet.Address(), sigTx.Hash(), len(update))
	return nil
}

func (relayer *Eth2TopRelayerV2) signTransaction(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	acc := relayer.wallet.Address()
	if strings.EqualFold(acc.Hex(), addr.Hex()) {
		stx, err := relayer.wallet.SignTx(tx)
		if err != nil {
			return nil, err
		}
		return stx, nil
	}
	return nil, fmt.Errorf("Eth2TopRelayer address:%v not available", addr)
}

//func (relayer *Eth2TopRelayerV2) getExecutionBlocksBetween(start, end primitives.Slot) ([]*types.Header, error) {
//	curSlot := start
//	logger.Info("Eth2TopRelayerV2 getExecutionBlocksBetween start:%v,end:%v", start, end)
//	var headers []*types.Header
//	for curSlot <= end {
//		header, err := relayer.getExecutionBlockBySlot(curSlot)
//		curSlot += 1
//		if err != nil {
//			if beaconrpc.IsErrorNoBlockForSlot(err) {
//				continue
//			}
//			logger.Error("Eth2TopRelayerV2 getExecutionBlockBySlot error", err)
//			return nil, err
//		}
//		headers = append(headers, header)
//	}
//	return headers, nil
//}

func (relayer *Eth2TopRelayerV2) getExecutionBlocksBetweenByNumber(low, height uint64) ([]*types.Header, error) {
	curNumber := low
	logger.Info("Eth2TopRelayerV2 SubmitBlocksNumber start:%v,end:%v", low, height)
	var headers []*types.Header
	for curNumber <= height {
		header, err := relayer.getExecutionBlockByNumber(curNumber)
		curNumber += 1
		if err != nil {
			if beaconrpc.IsErrorNoBlockForSlot(err) {
				continue
			}
			logger.Error("Eth2TopRelayerV2 getExecutionBlockBySlot error", err)
			return nil, err
		}
		headers = append(headers, header)
	}
	return headers, nil
}

type Output struct {
	HeaderRLP    string   `json:"header_rlp"`
	MerkleRoot   string   `json:"merkle_root"`
	Elements     []string `json:"elements"`
	MerkleProofs []string `json:"merkle_proofs"`
	ProofLength  uint64   `json:"proof_length"`
}

func (relayer *Eth2TopRelayerV2) encodeEthHeaders(headers []*types.Header) ([]byte, error) {
	var encodedHeaders []byte
	for _, header := range headers {
		rlpBytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			logger.Error("rlp encode error: ", err)
			return nil, err
		}
		if outBytes, err := rlp.EncodeToBytes(Output{
			HeaderRLP: string(rlpBytes),
		}); err != nil {
			logger.Error("Eth2TopRelayerV2 Output rlp encode error: ", err)
			return nil, err
		} else {
			encodedHeaders = append(encodedHeaders, outBytes...)
		}
		//encodedHeaders = append(encodedHeaders, rlpBytes...)
	}
	return encodedHeaders, nil
}

//func (relayer *Eth2TopRelayerV2) getExecutionBlockBySlot(slot primitives.Slot) (*types.Header, error) {
//	number, err := relayer.beaconClient.GetBlockNumberForSlot(slot)
//	if err != nil {
//		return nil, err
//	}
//	header, err := relayer.ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 HeaderByNumber error:", err)
//		return nil, err
//	}
//	return header, nil
//}

func (relayer *Eth2TopRelayerV2) getExecutionBlockByNumber(number uint64) (*types.Header, error) {
	header, err := relayer.ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 HeaderByNumber error:", err)
		return nil, err
	}
	return header, nil
}

func (relayer *Eth2TopRelayerV2) isEnoughBlocksForLightClientUpdate(lastFinalizedTopSlot, lastFinalizedEthSlot primitives.Slot) bool {
	if lastFinalizedTopSlot >= lastFinalizedEthSlot {
		return false
	}
	lastPeriodOnEth := ethereum.GetPeriodForSlot(lastFinalizedEthSlot)
	lastPeriodOnTop := ethereum.GetPeriodForSlot(lastFinalizedTopSlot)
	//if lastPeriodOnEth == lastPeriodOnTop+1 {
	//	if beaconrpc.GetFinalizedSlotForPeriod(lastPeriodOnEth) >= lastFinalizedEthSlot {
	//		return false
	//	}
	//}
	if lastPeriodOnEth == lastPeriodOnTop {
		if (lastFinalizedEthSlot - lastFinalizedTopSlot) < ethereum.ONE_EPOCH_IN_SLOTS*3 {
			return false
		}
		//submitFinalizedSlot, err := relayer.beaconrpcclient.GetLastFinalizedLightClientUpdateV2FinalizedSlot()
		//if err != nil {
		//	return false
		//}
		//if lastFinalizedTopSlot >= submitFinalizedSlot+32 {
		//	logger.Warn("must : lastFinalizedTopSlot(%d) < submitFinalizedSlot:(%d)", lastFinalizedTopSlot, submitFinalizedSlot)
		//	return false
		//}
	}
	return true
}

//func FilterSyncCommitteeVotes(committeeKeys [][]byte, sync *eth.SyncAggregate) ([]bls.PublicKey, error) {
//	if sync.SyncCommitteeBits.Len() > uint64(len(committeeKeys)) {
//		return nil, errors.New("bits length exceeds committee length")
//	}
//	votedKeys := make([]bls.PublicKey, 0, len(committeeKeys))
//	for i := uint64(0); i < sync.SyncCommitteeBits.Len(); i++ {
//		if sync.SyncCommitteeBits.BitAt(i) {
//			pubKey, err := bls.PublicKeyFromBytes(committeeKeys[i])
//			if err != nil {
//				return nil, err
//			}
//			votedKeys = append(votedKeys, pubKey)
//		}
//	}
//	return votedKeys, nil
//}

//func (relayer *Eth2TopRelayerV2) isCorrectFinalityUpdate(update *ethtypes.LightClientUpdate, committee *eth.SyncCommittee) error {
//	pubKeys, err := FilterSyncCommitteeVotes(committee.Pubkeys, update.SyncAggregate)
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 FilterSyncCommitteeVotes error:", err)
//		return err
//	}
//
//	domain, err := signing.ComputeDomain(ethtypes.DomainSyncCommittee, ethtypes.BellatrixForkVersion, ethtypes.GenesisValidatorsRoot[:])
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 ComputeDomain error:", err)
//		return err
//	}
//	pbr, err := update.AttestedBeaconHeader.HashTreeRoot()
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 HashTreeRoot error:", err)
//		return err
//	}
//	sszBytes := p2pType.SSZBytes(pbr[:])
//	signingRoot, err := signing.ComputeSigningRoot(&sszBytes, domain)
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 ComputeSigningRoot error:", err)
//		return err
//	}
//
//	aggregateSign, err := bls.SignatureFromBytes(update.SyncAggregate.SyncCommitteeSignature)
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 SignatureFromBytes error:", err)
//		return err
//	}
//	if !aggregateSign.Eth2FastAggregateVerify(pubKeys, signingRoot) {
//		return errors.New("invalid sync committee signature")
//	}
//	return nil
//}

//func (relayer *Eth2TopRelayerV2) verify_bls_signature_for_finality_update(update *ethtypes.LightClientUpdate) error {
//	signatureSlotPeriod := beaconrpc.GetPeriodForSlot(update.SignatureSlot)
//	topFinalizedBeaconBlockSlot, err := relayer.callerSession.FinalizedBeaconBlockSlot()
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 FinalizedBeaconBlockSlot error:", err)
//		return err
//	}
//	finalizedSlotPeriod := beaconrpc.GetPeriodForSlot(topFinalizedBeaconBlockSlot)
//	stateBytes, err := relayer.callerSession.GetLightClientState()
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 GetLightClientState error:", err)
//		return err
//	}
//	var state ethtypes.LightClientState
//	err = rlp.DecodeBytes(stateBytes, state)
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 DecodeBytes error:", err)
//		return err
//	}
//	var committee *eth.SyncCommittee
//	if signatureSlotPeriod == finalizedSlotPeriod {
//		committee = state.CurrentSyncCommittee
//	} else {
//		committee = state.NextSyncCommittee
//	}
//	return relayer.isCorrectFinalityUpdate(update, committee)
//}

type ExtendedBeaconBlockHeader struct {
	Header             *light_client.BeaconBlockHeader
	BeaconBlockRoot    [fieldparams.RootLength]byte
	ExecutionBlockHash [fieldparams.RootLength]byte
}

func (h *ExtendedBeaconBlockHeader) Encode() ([]byte, error) {
	headerBytes, err := h.Header.Encode()
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(headerBytes)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(h.BeaconBlockRoot)
	if err != nil {
		return nil, err
	}
	b3, err := rlp.EncodeToBytes(h.ExecutionBlockHash)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
	return rlpBytes, nil
}

type InitInput struct {
	FinalizedExecutionHeader *types.Header
	FinalizedBeaconHeader    *ExtendedBeaconBlockHeader
	CurrentSyncCommittee     *eth.SyncCommittee
	NextSyncCommittee        *eth.SyncCommittee
}

func (init *InitInput) Encode() ([]byte, error) {
	exeHeader, err := rlp.EncodeToBytes(init.FinalizedExecutionHeader)
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(exeHeader)
	if err != nil {
		return nil, err
	}
	finHeader, err := init.FinalizedBeaconHeader.Encode()
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(finHeader)
	if err != nil {
		return nil, err
	}
	cur, err := rlp.EncodeToBytes(init.CurrentSyncCommittee)
	if err != nil {
		return nil, err
	}
	b3, err := rlp.EncodeToBytes(cur)
	if err != nil {
		return nil, err
	}
	next, err := rlp.EncodeToBytes(init.NextSyncCommittee)
	if err != nil {
		return nil, err
	}
	b4, err := rlp.EncodeToBytes(next)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
	rlpBytes = append(rlpBytes, b4...)
	return rlpBytes, nil
}

func (relayer *Eth2TopRelayerV2) GetInitData() ([]byte, error) {
	lastSlot, err := relayer.beaconClient.GetLastFinalizedSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return nil, err
	}
	lastPeriod := ethereum.GetPeriodForSlot(lastSlot)
	lastUpdate, err := relayer.beaconClient.GetLightClientUpdate(lastPeriod)
	if err != nil {
		logger.Error("GetLightClientUpdate error:", err)
		return nil, err
	}
	// prevUpdate, err := relayer.beaconrpcclient.GetLightClientUpdate(lastPeriod - 1)
	// if err != nil {
	// 	logger.Error("GetLightClientUpdate error:", err)
	// 	return nil, err
	// }
	prevUpdate, err := relayer.beaconClient.GetNextSyncCommitteeUpdate(lastPeriod - 1)
	if err != nil {
		logger.Error("GetNextSyncCommitteeUpdate error:", err)
		return nil, err
	}

	var beaconHeader eth.BeaconBlockHeader
	beaconHeader.Slot = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	beaconHeader.ProposerIndex = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ProposerIndex
	beaconHeader.BodyRoot = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.BodyRoot[:]
	beaconHeader.ParentRoot = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ParentRoot[:]
	beaconHeader.StateRoot = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.StateRoot[:]
	root, err := beaconHeader.HashTreeRoot()
	if err != nil {
		logger.Error("HashTreeRoot error:", err)
		return nil, err
	}
	finalizedHeader := new(ExtendedBeaconBlockHeader)
	finalizedHeader.BeaconBlockRoot = root
	finalizedHeader.Header = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader
	finalizedHeader.ExecutionBlockHash = lastUpdate.FinalityUpdate.HeaderUpdate.ExecutionBlockHash

	finalitySlot := lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	finalizeBody, err := relayer.beaconClient.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(finalitySlot), 10)))
	if err != nil {
		logger.Error("GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	executionData, err := finalizeBody.Execution()
	if err != nil {
		return nil, err
	}
	number := executionData.BlockNumber()

	header, err := relayer.ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		logger.Error("HeaderByNumber error:", err)
		return nil, err
	}

	initParam := new(InitInput)
	initParam.FinalizedExecutionHeader = header
	initParam.FinalizedBeaconHeader = finalizedHeader
	initParam.NextSyncCommittee = lastUpdate.NextSyncCommitteeUpdate.NextSyncCommittee
	initParam.CurrentSyncCommittee = prevUpdate.NextSyncCommittee
	bytes, err := initParam.Encode()
	if err != nil {
		logger.Error("initParam.Encode error:", err)
		return nil, err
	}
	return bytes, nil
}
