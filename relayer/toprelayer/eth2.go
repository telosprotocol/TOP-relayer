package toprelayer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
	"toprelayer/config"
	eth2bridge "toprelayer/contract/top/eth2client"
	"toprelayer/relayer/toprelayer/ethashapp"
	"toprelayer/relayer/toprelayer/ethtypes"
	beaconrpc "toprelayer/rpc/ethbeacon_rpc"
	"toprelayer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	p2pType "github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/types"
	primitives "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
)

var (
	eth2ClientSystemContract = common.HexToAddress("0xff00000000000000000000000000000000000009")
)

type Eth2TopRelayerV2 struct {
	wallet          *wallet.Wallet
	ethrpcclient    *ethclient.Client
	beaconrpcclient *beaconrpc.BeaconGrpcClient
	transactor      *eth2bridge.Eth2ClientTransactor
	callerSession   *eth2bridge.Eth2ClientCallerSession
	lastSlot        uint64
}

func NewRelayerByRpcClient(prysm string) *Eth2TopRelayerV2 {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, "")
	if err != nil {
		logger.Error("NewRelayerByRpcClient NewBeaconGrpcClient error:", err)
		return nil
	}
	return &Eth2TopRelayerV2{
		beaconrpcclient: b,
	}
}

func (relayer *Eth2TopRelayerV2) Init(cfg *config.Relayer, listenUrl []string, pass string) error {
	w, err := wallet.NewTopWallet(cfg.Url[0], cfg.KeyPath, pass)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 NewWallet error:", err)
		return err
	}
	relayer.wallet = w

	if len(listenUrl) != 3 {
		err := errors.New("listenUrl num error")
		logger.Error("Eth2TopRelayerV2 listenUrl error:", err)
		return err
	}
	relayer.ethrpcclient, err = ethclient.Dial(listenUrl[0])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 ethclient.Dial error:", err)
		return err
	}
	relayer.beaconrpcclient, err = beaconrpc.NewBeaconGrpcClient(listenUrl[1], listenUrl[2])
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
	return nil
}

func (relayer *Eth2TopRelayerV2) blockKnownOnTop(slot uint64) (bool, error) {
	hash, err := relayer.beaconrpcclient.GetBlockHashForSlot(slot)
	//logger.Debug("blockKnownOnTop slot %v, hash %v", slot, hash)
	if err != nil {
		return false, err
	}
	return relayer.callerSession.IsKnownExecutionHeader(hash)
}

func (relayer *Eth2TopRelayerV2) findLeftNonErrorSlot(leftSlot, rightSlot uint64) (uint64, bool) {
	slot := leftSlot
	for slot != rightSlot {
		known, err := relayer.blockKnownOnTop(slot)
		if err != nil {
			slot += 1
		} else {
			return slot, known
		}
	}
	return slot, false
}

func (relayer *Eth2TopRelayerV2) linearSearchForward(slot, maxSlot uint64) (uint64, error) {
	for {
		if slot >= maxSlot {
			break
		}
		known, err := relayer.blockKnownOnTop(slot + 1)
		if err != nil {
			if beaconrpc.IsErrorNoBlockForSlot(err) {
				slot += 1
				continue
			} else {
				logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
				return 0, err
			}
		}
		if known {
			slot += 1
			logger.Debug("curr top known slot: %v,maxSlot: %v", slot, maxSlot)
		} else {
			break
		}
	}
	logger.Debug("linearSearchForward return slot: %v", slot)
	return slot, nil
}

func (relayer *Eth2TopRelayerV2) linearSearchBackward(startSlot, lastSlot uint64) (uint64, error) {
	slot := lastSlot
	lastFalseSlot := slot + 1
	for {
		if slot <= startSlot {
			break
		}
		known, err := relayer.blockKnownOnTop(slot)
		if err != nil {
			if beaconrpc.IsErrorNoBlockForSlot(err) {
				slot -= 1
				continue
			} else {
				logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
				return 0, err
			}
		}
		if known {
			break
		} else {
			lastFalseSlot = slot
			slot -= 1
		}
	}
	return lastFalseSlot - 1, nil
}

func (relayer *Eth2TopRelayerV2) linerSlotSearch(slot, finalizedSlot, lastEthSlot uint64) (uint64, error) {
	if slot == finalizedSlot {
		logger.Info("slot equal finalizedSlot %v go forward,greater than lastEthSlot:%v", slot, lastEthSlot)
		return relayer.linearSearchForward(slot, lastEthSlot)
	}
	known, err := relayer.blockKnownOnTop(slot)
	if err != nil {
		if beaconrpc.IsErrorNoBlockForSlot(err) {
			leftSlot, known := relayer.findLeftNonErrorSlot(slot+1, lastEthSlot+1)
			if known {
				return relayer.linearSearchForward(leftSlot, lastEthSlot)
			} else {
				return relayer.linearSearchForward(finalizedSlot, leftSlot-1)
			}
		} else {
			logger.Error("Eth2TopRelayerV2 blockKnownOnTop error", err)
			return 0, err
		}
	}
	if known {
		logger.Debug("slot %v known, go forward", slot)
		return relayer.linearSearchForward(slot, lastEthSlot)
	} else {
		logger.Debug("slot %v unknown, go backward", slot)
		return relayer.linearSearchBackward(finalizedSlot, slot)
	}
}

func (relayer *Eth2TopRelayerV2) getMaxSlotForSubmission() (uint64, error) {
	return relayer.beaconrpcclient.GetLastSlotNumber()
}

func (relayer *Eth2TopRelayerV2) getLastEth2SlotOnTop(lastEthSlot uint64) (uint64, error) {
	finalizedSlot, err := relayer.callerSession.FinalizedBeaconBlockSlot()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 FinalizedBeaconBlockSlot error", err)
		return 0, nil
	}
	lastSubmittedSlot := relayer.lastSlot
	slot := finalizedSlot
	if lastSubmittedSlot > finalizedSlot {
		slot = lastSubmittedSlot
	}
	logger.Debug("getLastEth2SlotOnTop finalizedSlot: %v, lastSubmittedSlot: %v, slot: %v", finalizedSlot, lastSubmittedSlot, slot)
	return relayer.linerSlotSearch(slot, finalizedSlot, lastEthSlot)
}

func (relayer *Eth2TopRelayerV2) getLastFinalizedSlotOnTop() (uint64, error) {
	return relayer.callerSession.FinalizedBeaconBlockSlot()
}

func (relayer *Eth2TopRelayerV2) getLastFinalizedSlotOnEth() (uint64, error) {
	return relayer.beaconrpcclient.GetLastFinalizedSlotNumber()
}

func (relayer *Eth2TopRelayerV2) submitExecutionBlocks(headers []byte, curSlot uint64) error {
	if len(headers) > 0 {
		err := relayer.submitEthHeader(headers)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 submitHeaders failed:", err)
			return err
		}
	}
	relayer.lastSlot = curSlot
	return nil
}

func (relayer *Eth2TopRelayerV2) sendRegularLightClientUpdate(lastFinalizedTopSlot, lastFinalizedEthSlot uint64) error {
	lastEth2PeriodOnTopChain := beaconrpc.GetPeriodForSlot(lastFinalizedTopSlot)
	endPeriod := beaconrpc.GetPeriodForSlot(lastFinalizedEthSlot)
	logger.Info("Eth2TopRelayerV2 sendRegularLightClientUpdate period: %d, %d", lastEth2PeriodOnTopChain, endPeriod)

	var data *beaconrpc.LightClientUpdate
	var err error
	if lastEth2PeriodOnTopChain == endPeriod {
		data, err = relayer.beaconrpcclient.GetFinalizedLightClientUpdateV2()
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetLightClientUpdate error:", err)
			return err
		}
	} else {
		data, err = relayer.beaconrpcclient.GetLightClientUpdateV2(lastEth2PeriodOnTopChain + 1)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetLightClientUpdate error:", err)
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

func (relayer *Eth2TopRelayerV2) sendLightClientUpdatesWithChecks(slot uint64) (bool, error) {
	topFinalizedSlot, err := relayer.getLastFinalizedSlotOnTop()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getLastFinalizedSlotOnTop error:", err)
		return false, err
	}
	ethSlot, err := relayer.getLastFinalizedSlotOnEth()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getLastFinalizedSlotOnEth error:", err)
		return false, err
	}
	if relayer.isEnoughBlocksForLightClientUpdate(slot, topFinalizedSlot, ethSlot) {
		err = relayer.sendRegularLightClientUpdate(topFinalizedSlot, ethSlot)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 sendLightClientUpdates error:", err)
			return false, err
		}
		return true, nil
	} else {
		logger.Debug("Eth2TopRelayerV2 no need to light client update")
	}
	return false, nil
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
		logger.Error(fmt.Sprintf("Eth2TopRelayer EstimateGas error:%s, data:%v", err, common.Bytes2Hex(packData)))
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

func (relayer *Eth2TopRelayerV2) submitEthHeader(headers []byte) error {
	packHeader, err := eth2bridge.PackSubmitExecutionHeaderParam(headers)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 PackSubmitExecutionHeaderParam error:", err)
		return err
	}
	ops, err := relayer.txOption(packHeader)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 txOption error:", err)
		return err
	}
	sigTx, err := relayer.transactor.SubmitExecutionHeader(ops, headers)
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
	logger.Info("Eth2TopRelayer submitLightClientUpdate data:", common.Bytes2Hex(update))
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

func (relayer *Eth2TopRelayerV2) StartRelayer(wg *sync.WaitGroup) error {
	logger.Info("Start Eth2TopRelayerV2, subBatch: %v certaintyBlocks: %v", BATCH_NUM, CONFIRM_NUM)
	defer wg.Done()

	done := make(chan struct{})
	defer close(done)

	go func(done chan struct{}) {
		timeoutDuration := time.Duration(FATALTIMEOUT) * time.Hour
		timeout := time.NewTimer(timeoutDuration)
		defer timeout.Stop()
		logger.Debug("Eth2TopRelayerV2 set timeout: %v hours", FATALTIMEOUT)
		var delay time.Duration = time.Duration(1)

		prevPeriod := uint64(0)
		curPeriod := uint64(0)

		for {
			time.Sleep(time.Second * delay)
			select {
			case <-timeout.C:
				done <- struct{}{}
				return
			default:
				for {
					time.Sleep(time.Second * delay)
					// step1: eth slot
					eth2Slot, err := relayer.getMaxSlotForSubmission()
					if err != nil {
						logger.Error(err)
						delay = time.Duration(ERRDELAY)
						break
					}
					if eth2Slot == 0 {
						logger.Info("Eth2TopRelayerV2 beacon endpoint slot 0")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("Eth2TopRelayerV2 check src eth2 slot:", eth2Slot)
					// step2: top slot
					topSlot, err := relayer.getLastEth2SlotOnTop(eth2Slot)
					if err != nil {
						logger.Error(err)
						delay = time.Duration(ERRDELAY)
						break
					}
					if topSlot == 0 {
						if set := timeout.Reset(timeoutDuration); !set {
							logger.Error("Eth2TopRelayerV2 reset timeout falied!")
							delay = time.Duration(ERRDELAY)
						} else {
							logger.Info("Eth2TopRelayerV2 not init yet")
							delay = time.Duration(ERRDELAY)
						}
						break
					}
					logger.Info("Eth2TopRelayerV2 check dest top slot:", topSlot)
					// step3: submit headers
					if topSlot < eth2Slot {
						headers, curSlot, err := relayer.getExecutionBlocksBetween(topSlot+1, eth2Slot)
						if err != nil {
							logger.Error("Eth2TopRelayerV2 GetExecutionBlocksBetween failed:", err)
							delay = time.Duration(ERRDELAY)
							break
						}
						if curSlot+8 > eth2Slot {
							logger.Info("Eth2TopRelayerV2 headers update not finish, continue update headers next round")
							delay = time.Duration(SUCCESSDELAY)
							break
						}
						if err = relayer.submitExecutionBlocks(headers, curSlot); err != nil {
							logger.Error("Eth2TopRelayerV2 submitExecutionBlocks failed:", err)
							delay = time.Duration(ERRDELAY)
							break
						} else {
							topSlot = curSlot
						}
						if prevPeriod == 0 {
							if topLastSlot, err := relayer.getLastFinalizedSlotOnTop(); err != nil {
								logger.Error("Eth2TopRelayerV2 getLastFinalizedSlotOnTop error:", err)
							} else {
								prevPeriod = beaconrpc.GetPeriodForSlot(topLastSlot)
							}
						}
						curPeriod = beaconrpc.GetPeriodForSlot(curSlot)
						logger.Info("Eth2TopRelayerV2 topFinalized(period:%v),topUnFinalized(period:%v,slot:%v)", prevPeriod, curPeriod, curSlot)
					}

					logger.Info("Eth2TopRelayerV2 headers update finish, update light client update for a while")
					time.Sleep(time.Second * time.Duration(SUCCESSDELAY))
					if ret, err := relayer.sendLightClientUpdatesWithChecks(topSlot); err != nil {
						logger.Error("Eth2TopRelayerV2 sendLightClientUpdatesWithChecks error:", err)
					} else if ret == true {
						prevPeriod = curPeriod
					}
					if set := timeout.Reset(timeoutDuration); !set {
						logger.Error("Eth2TopRelayerV2 reset timeout falied!")
						delay = time.Duration(ERRDELAY)
						break
					}
					logger.Info("Eth2TopRelayerV2 sync round finish")
					delay = time.Duration(SUCCESSDELAY)
				}
			}
		}
	}(done)

	<-done
	logger.Error("Eth2TopRelayerV2 timeout")
	return nil
}

func (relayer *Eth2TopRelayerV2) getExecutionBlocksBetween(start, end uint64) ([]byte, uint64, error) {
	curSlot := start
	headersCnt := 0
	var batchHeaders []byte
	for (headersCnt < beaconrpc.HEADER_BATCH_SIZE) && (curSlot <= end) {
		header, err := relayer.getExecutionBlockBySlot(curSlot)
		if err != nil {
			if beaconrpc.IsErrorNoBlockForSlot(err) {
				curSlot += 1
				continue
			}
			logger.Error("Eth2TopRelayerV2 getExecutionBlockBySlot error", err)
			return nil, 0, err
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			logger.Error("rlp encode error: ", err)
		}
		var out ethashapp.Output
		out.HeaderRLP = string(rlp_bytes)
		outBytes, err := rlp.EncodeToBytes(out)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 Output rlp encode error: ", err)
			return nil, 0, err
		}
		batchHeaders = append(batchHeaders, outBytes...)
		curSlot += 1
		headersCnt += 1
	}
	curSlot -= 1
	return batchHeaders, curSlot, nil
}

func (relayer *Eth2TopRelayerV2) getExecutionBlockBySlot(slot uint64) (*types.Header, error) {
	number, err := relayer.beaconrpcclient.GetBlockNumberForSlot(slot)
	if err != nil {
		return nil, err
	}
	header, err := relayer.ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 HeaderByNumber error:", err)
		return nil, err
	}
	return header, nil
}

func (relayer *Eth2TopRelayerV2) isEnoughBlocksForLightClientUpdate(lastSubmittedSlot, lastFinalizedTopSlot, lastFinalizedEthSlot uint64) bool {
	if lastSubmittedSlot < lastFinalizedTopSlot {
		return false
	}
	if (lastSubmittedSlot - lastFinalizedTopSlot) < beaconrpc.ONE_EPOCH_IN_SLOTS {
		return false
	}
	if lastFinalizedEthSlot <= lastFinalizedTopSlot {
		return false
	}
	// period of different LightClientUpdate should be submitted,Unless it's in the same period as eth
	if beaconrpc.GetPeriodForSlot(lastSubmittedSlot) == beaconrpc.GetPeriodForSlot(lastFinalizedTopSlot) {
		if beaconrpc.GetPeriodForSlot(lastSubmittedSlot) != beaconrpc.GetPeriodForSlot(lastFinalizedEthSlot) {
			return false
		}
	}
	return true
}

func (relayer *Eth2TopRelayerV2) sendLightClientUpdates(lastFinalizedTopSlot, lastFinalizedEthSlot uint64) error {
	attestedSlot, err := relayer.beaconrpcclient.GetAttestedSlot(lastFinalizedTopSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlot error:", err)
		return err
	}
	lastTopPeriod := beaconrpc.GetPeriodForSlot(lastFinalizedTopSlot)
	endPeriod := beaconrpc.GetPeriodForSlot(lastFinalizedEthSlot)
	useNextSyncCommittee := lastTopPeriod == endPeriod
	for {
		update, err := relayer.beaconrpcclient.GetFinalityLightClientUpdate(attestedSlot, useNextSyncCommittee)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 getFinalityLightClientUpdate error:", err)
			return err
		}
		finalityUpdateSlot := uint64(update.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot)
		if finalityUpdateSlot <= lastFinalizedTopSlot {
			attestedSlot, err = relayer.beaconrpcclient.GetAttestedSlot(lastFinalizedTopSlot + beaconrpc.SLOTS_PER_EPOCH)
			if err != nil {
				logger.Error("Eth2TopRelayerV2 getAttestedSlot error:", err)
				return err
			}
			continue
		}
		return relayer.sendSpecificLightClientUpdate(update)
	}
}

func (relayer *Eth2TopRelayerV2) sendSpecificLightClientUpdate(update *ethtypes.LightClientUpdate) error {
	isKnown, err := relayer.callerSession.IsKnownExecutionHeader(update.FinalizedUpdate.HeaderUpdate.ExecutionBlockHash)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 IsKnownExecutionHeader error:", err)
		return err
	}
	if !isKnown {
		logger.Error("Eth2TopRelayerV2 IsKnownExecutionHeader not known block")
		return nil
	}
	err = relayer.verify_bls_signature_for_finality_update(update)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 verify_bls_signature_for_finality_update error:", err)
		return nil
	}
	upateBytes, err := rlp.EncodeToBytes(update)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 EncodeToBytes error:", err)
		return nil
	}
	err = relayer.submitLightClientUpdate(upateBytes)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 submitLightClientUpdate error:", err)
		return err
	}

	return nil
}

func FilterSyncCommitteeVotes(committeeKeys [][]byte, sync *eth.SyncAggregate) ([]bls.PublicKey, error) {
	if sync.SyncCommitteeBits.Len() > uint64(len(committeeKeys)) {
		return nil, errors.New("bits length exceeds committee length")
	}
	votedKeys := make([]bls.PublicKey, 0, len(committeeKeys))
	for i := uint64(0); i < sync.SyncCommitteeBits.Len(); i++ {
		if sync.SyncCommitteeBits.BitAt(i) {
			pubKey, err := bls.PublicKeyFromBytes(committeeKeys[i])
			if err != nil {
				return nil, err
			}
			votedKeys = append(votedKeys, pubKey)
		}
	}
	return votedKeys, nil
}

func (relayer *Eth2TopRelayerV2) isCorrectFinalityUpdate(update *ethtypes.LightClientUpdate, committee *eth.SyncCommittee) error {
	pubKeys, err := FilterSyncCommitteeVotes(committee.Pubkeys, update.SyncAggregate)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 FilterSyncCommitteeVotes error:", err)
		return err
	}

	domain, err := signing.ComputeDomain(ethtypes.DomainSyncCommittee, ethtypes.BellatrixForkVersion, ethtypes.GenesisValidatorsRoot[:])
	if err != nil {
		logger.Error("Eth2TopRelayerV2 ComputeDomain error:", err)
		return err
	}
	pbr, err := update.AttestedBeaconHeader.HashTreeRoot()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 HashTreeRoot error:", err)
		return err
	}
	sszBytes := p2pType.SSZBytes(pbr[:])
	signingRoot, err := signing.ComputeSigningRoot(&sszBytes, domain)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 ComputeSigningRoot error:", err)
		return err
	}

	aggregateSign, err := bls.SignatureFromBytes(update.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 SignatureFromBytes error:", err)
		return err
	}
	if !aggregateSign.Eth2FastAggregateVerify(pubKeys, signingRoot) {
		return errors.New("invalid sync committee signature")
	}
	return nil
}

func (relayer *Eth2TopRelayerV2) verify_bls_signature_for_finality_update(update *ethtypes.LightClientUpdate) error {
	signatureSlotPeriod := beaconrpc.GetPeriodForSlot(update.SignatureSlot)
	topFinalizedBeaconBlockSlot, err := relayer.callerSession.FinalizedBeaconBlockSlot()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 FinalizedBeaconBlockSlot error:", err)
		return err
	}
	finalizedSlotPeriod := beaconrpc.GetPeriodForSlot(topFinalizedBeaconBlockSlot)
	stateBytes, err := relayer.callerSession.GetLightClientState()
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetLightClientState error:", err)
		return err
	}
	var state ethtypes.LightClientState
	rlp.DecodeBytes(stateBytes, state)
	var committee *eth.SyncCommittee
	if signatureSlotPeriod == finalizedSlotPeriod {
		committee = state.CurrentSyncCommittee
	} else {
		committee = state.NextSyncCommittee
	}
	return relayer.isCorrectFinalityUpdate(update, committee)
}

type ExtendedBeaconBlockHeader struct {
	Header             *beaconrpc.BeaconBlockHeader
	BeaconBlockRoot    []byte
	ExecutionBlockHash []byte
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
	lastSlot, err := relayer.beaconrpcclient.GetLastFinalizedSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return nil, err
	}
	lastPeriod := beaconrpc.GetPeriodForSlot(lastSlot)
	lastUpdate, err := relayer.beaconrpcclient.GetLightClientUpdate(lastPeriod)
	if err != nil {
		logger.Error("GetLightClientUpdate error:", err)
		return nil, err
	}
	// prevUpdate, err := relayer.beaconrpcclient.GetLightClientUpdate(lastPeriod - 1)
	// if err != nil {
	// 	logger.Error("GetLightClientUpdate error:", err)
	// 	return nil, err
	// }
	prevUpdate, err := relayer.beaconrpcclient.GetNextSyncCommitteeUpdate(lastPeriod - 1)
	if err != nil {
		logger.Error("GetNextSyncCommitteeUpdate error:", err)
		return nil, err
	}

	var beaconHeader eth.BeaconBlockHeader
	beaconHeader.Slot = primitives.Slot(lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot)
	beaconHeader.ProposerIndex = primitives.ValidatorIndex(lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.ProposerIndex)
	beaconHeader.BodyRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.BodyRoot
	beaconHeader.ParentRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.ParentRoot
	beaconHeader.StateRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.StateRoot
	root, err := beaconHeader.HashTreeRoot()
	if err != nil {
		logger.Error("HashTreeRoot error:", err)
		return nil, err
	}
	finalizedHeader := new(ExtendedBeaconBlockHeader)
	finalizedHeader.BeaconBlockRoot = root[:]
	finalizedHeader.Header = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader
	finalizedHeader.ExecutionBlockHash = lastUpdate.FinalizedUpdate.HeaderUpdate.ExecutionBlockHash

	finalitySlot := lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot
	finalizeBody, err := relayer.beaconrpcclient.GetBeaconBlockBodyForBlockId(strconv.FormatUint(finalitySlot, 10))
	if err != nil {
		logger.Error("GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	number := finalizeBody.GetExecutionPayload().BlockNumber

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
