package ethbeacon_rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	state_native "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	v2 "github.com/prysmaticlabs/prysm/v4/proto/eth/v2"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"sort"
	"strconv"
	"toprelayer/relayer/toprelayer/ethtypes"
)

func (c *BeaconGrpcClient) GetLastFinalizedLightClientUpdateV2FinalizedSlot() (uint64, error) {
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		return 0, err
	}
	return getBeforeSlotInSamePeriod(finalizedSlot)
}

func (c *BeaconGrpcClient) GetLastFinalizedLightClientUpdateV2() (*LightClientUpdate, error) {
	finalizedSlot, err := c.GetLastFinalizedLightClientUpdateV2FinalizedSlot()
	if err != nil {
		return nil, err
	}
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, false)
}

func (c *BeaconGrpcClient) GetLastFinalizedLightClientUpdateV2WithNextSyncCommittee() (*LightClientUpdate, error) {
	finalizedSlot, err := c.GetLastFinalizedLightClientUpdateV2FinalizedSlot()
	if err != nil {
		return nil, err
	}
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, true)
}

func (c *BeaconGrpcClient) GetLightClientUpdateV2(period uint64) (*LightClientUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getLightClientUpdateByFinalizedSlot(currFinalizedSlot, true)
}

func (c *BeaconGrpcClient) GetNextSyncCommitteeUpdateV2(period uint64) (*SyncCommitteeUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getNextSyncCommitteeUpdateByFinalized(currFinalizedSlot)
}

func (c *BeaconGrpcClient) getLightClientUpdateByFinalizedSlot(finalizedSlot uint64, useNextSyncCommittee bool) (*LightClientUpdate, error) {
	attestedSlot, err := c.GetAttestedSlot(finalizedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return nil, err
	}
	if GetPeriodForSlot(attestedSlot) != GetPeriodForSlot(finalizedSlot) {
		return nil, fmt.Errorf("Eth2TopRelayerV2 attestedSlot(%d) and finalizedSlot(%d) not in same period", attestedSlot, finalizedSlot)
	}
	lcu, err := c.getFinalityLightClientUpdate(attestedSlot, useNextSyncCommittee)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getFinalityLightClientUpdate error:", err)
		return nil, err
	}
	logger.Info("LightClientUpdate FinalizedSlot:%d,AttestedSlot:%d",
		lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot, lcu.AttestedBeaconHeader.Slot)
	return convertEth2LightClientUpdate(lcu), nil
}

func (c *BeaconGrpcClient) getNextSyncCommitteeUpdateByFinalized(finalizedSlot uint64) (*SyncCommitteeUpdate, error) {
	attestedSlot, err := c.GetAttestedSlot(finalizedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	if GetPeriodForSlot(attestedSlot) != GetPeriodForSlot(finalizedSlot) {
		return nil, fmt.Errorf("Eth2TopRelayerV2 GetNextSyncCommitteeUpdateV2 attestedSlot(%d) and finalizedSlot(%d) not in same period", attestedSlot, finalizedSlot)
	}
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetNextSyncCommitteeUpdateV2 attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)
	beaconState, err := c.getBeaconState(strconv.FormatUint(attestedSlot, 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getBeaconState error:", err)
		return nil, err
	}
	cu, err := c.getNextSyncCommittee(beaconState)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getNextSyncCommittee error:", err)
		return nil, err
	}
	return &SyncCommitteeUpdate{
		NextSyncCommittee:       cu.NextSyncCommittee,
		NextSyncCommitteeBranch: cu.NextSyncCommitteeBranch,
	}, nil
}

func (c *BeaconGrpcClient) GetAttestedSlot(lastFinalizedSlot uint64) (uint64, error) {
	attestedSlot := getAttestationSlot(lastFinalizedSlot)
	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return 0, err
	}
	return uint64(header.Slot), nil
}

//func (c *BeaconGrpcClient) getAttestedSlotBeforeFinalizedSlot(finalizedSlot uint64) (uint64, error) {
//	// the range of attestedSlot is [finalizedSlot - 32, finalizedSlot]
//	attestedSlot := finalizedSlot - ONE_EPOCH_IN_SLOTS
//	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
//	if err != nil {
//		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
//		return 0, err
//	}
//	return uint64(header.Slot), nil
//}

func (c *BeaconGrpcClient) getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot uint64) (uint64, uint64, error) {
	currentAttestedSlot := attestedSlot
	for {
		h, err := c.GetNonEmptyBeaconBlockHeader(currentAttestedSlot + 1)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
			return 0, 0, err
		}
		signatureSlot := uint64(h.Slot)
		body, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(signatureSlot, 10))
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
			return 0, 0, err
		}
		syncCommitteeBitsSum := body.SyncAggregate.SyncCommitteeBits.Count()
		if syncCommitteeBitsSum*3 < (64 * 8 * 2) {
			currentAttestedSlot = signatureSlot
			continue
		}
		if len(body.Attestations) == 0 {
			currentAttestedSlot = signatureSlot
			continue
		}
		var attestedSlots []uint64
		for _, attestation := range body.Attestations {
			attestedSlots = append(attestedSlots, uint64(attestation.GetData().Slot))
		}
		sort.Slice(attestedSlots, func(i, j int) bool { return attestedSlots[i] > attestedSlots[j] })
		for i, v := range attestedSlots {
			if (i == 0 || v != attestedSlots[i-1]) && v >= attestedSlot {
				currentAttestedSlot = v
				_, err = c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(currentAttestedSlot, 10))
				if err != nil {
					continue
				}
				return currentAttestedSlot, signatureSlot, nil
			}
		}
		currentAttestedSlot = signatureSlot
	}
}

func (c *BeaconGrpcClient) getNextSyncCommittee(beaconState *eth.BeaconStateCapella) (*ethtypes.SyncCommitteeUpdate, error) {
	if beaconState.NextSyncCommittee == nil {
		logger.Error("Eth2TopRelayerV2 NextSyncCommittee nil")
		return nil, errors.New("NextSyncCommittee nil")
	}
	var state, err = state_native.InitializeFromProtoCapella(beaconState)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 InitializeFromProtoUnsafeBellatrix error:", err)
		return nil, err
	}
	nscp, err := state.NextSyncCommitteeProof(context.Background())
	if err != nil {
		logger.Error("Eth2TopRelayerV2 NextSyncCommitteeProof error:", err)
		return nil, err
	}
	update := &ethtypes.SyncCommitteeUpdate{
		NextSyncCommittee:       beaconState.NextSyncCommittee,
		NextSyncCommitteeBranch: nscp,
	}
	return update, nil
}

func (c *BeaconGrpcClient) constructFromBeaconBlockBody(beaconBlockBody *v2.BeaconBlockBodyCapella) (*ethtypes.ExecutionBlockProof, error) {
	blockHash := beaconBlockBody.ExecutionPayload.GetBlockHash()
	var finalizedBlockBodyHash common.Hash
	copy(finalizedBlockBodyHash[:], blockHash[:])
	var err error
	var beaconBlockMerkleTree, executionPayloadMerkleTree MerkleTreeNode
	if beaconBlockMerkleTree, err = BeaconBlockBodyMerkleTreeNew(beaconBlockBody); err != nil {
		return nil, err
	}
	if executionPayloadMerkleTree, err = ExecutionPayloadMerkleTreeNew(beaconBlockBody.GetExecutionPayload()); err != nil {
		return nil, err
	}
	_, proof1 := generateProof(beaconBlockMerkleTree, L1BeaconBlockBodyTreeExecutionPayloadIndex, L1BeaconBlockBodyProofSize)
	_, proof2 := generateProof(executionPayloadMerkleTree, L2ExecutionPayloadTreeExecutionBlockIndex, L2ExecutionPayloadProofSize)
	proof2 = append(proof2, proof1...)
	return &ethtypes.ExecutionBlockProof{
		BlockHash: finalizedBlockBodyHash,
		Proof:     ethtypes.ConvertSliceBytes2Hash(proof2),
	}, nil
}

func (c *BeaconGrpcClient) getFinalityLightClientUpdateForState(attestedSlot, signatureSlot uint64, beaconState, finalityBeaconState *eth.BeaconStateCapella) (*ethtypes.LightClientUpdate, error) {
	signatureBeaconBody, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(signatureSlot, 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	if signatureBeaconBody.SyncAggregate == nil {
		logger.Error("Eth2TopRelayerV2 syncAggregate nil")
		return nil, errors.New("syncAggregate nil")
	}
	attestedHeader, err := c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(uint64(attestedSlot), 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconBlockHeaderForBlockId error:", err)
		return nil, err
	}
	finalityHash := beaconState.FinalizedCheckpoint.Root
	finalityHeader, err := c.GetBeaconBlockHeaderForBlockId(string(finalityHash))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconBlockHeaderForBlockId error:", err)
		return nil, err
	}
	finalizedBlockBody, err := c.GetBeaconBlockBodyForBlockId(string(finalityHash))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	executionBlockProof, err := c.constructFromBeaconBlockBody(finalizedBlockBody)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 constructFromBeaconBlockBody hash error:", err)
		return nil, err
	}
	if err != nil {
		logger.Error("Eth2TopRelayerV2 finalizedBlockBody hash error:", err)
		return nil, err
	}
	state, err := state_native.InitializeFromProtoUnsafeCapella(beaconState)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 InitializeFromProtoUnsafeBellatrix error:", err)
		return nil, err
	}
	proof, err := state.FinalizedRootProof(context.Background())
	if err != nil {
		logger.Error("Eth2TopRelayerV2 FinalizedRootProof error:", err)
		return nil, err
	}
	update := &ethtypes.LightClientUpdate{
		AttestedBeaconHeader: attestedHeader,
		SyncAggregate: &eth.SyncAggregate{
			SyncCommitteeBits:      signatureBeaconBody.SyncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: signatureBeaconBody.SyncAggregate.SyncCommitteeSignature,
		},
		SignatureSlot: signatureSlot,
	}
	update.FinalizedUpdate = &ethtypes.FinalizedHeaderUpdate{
		HeaderUpdate: &ethtypes.HeaderUpdate{
			BeaconHeader:        finalityHeader,
			ExecutionBlockHash:  executionBlockProof.BlockHash,
			ExecutionHashBranch: executionBlockProof.Proof,
		},
		FinalityBranch: proof,
	}
	if finalityBeaconState != nil {
		update.NextSyncCommitteeUpdate, err = c.getNextSyncCommittee(finalityBeaconState)
		if err != nil {
			logger.Error("Eth2TopRelayerV2 getNextSyncCommittee error:", err)
			return nil, err
		}
	}
	return update, nil
}

func (c *BeaconGrpcClient) getFinalityLightClientUpdate(attestedSlot uint64, useNextSyncCommittee bool) (*ethtypes.LightClientUpdate, error) {
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetFinalityLightClientUpdate attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)
	beaconState, err := c.getBeaconState(strconv.FormatUint(attestedSlot, 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getBeaconState error:", err)
		return nil, err
	}
	var finalityBeaconState *eth.BeaconStateCapella = nil
	if useNextSyncCommittee == true {
		finalityBeaconState = beaconState
	}
	return c.getFinalityLightClientUpdateForState(attestedSlot, signatureSlot, beaconState, finalityBeaconState)
}
