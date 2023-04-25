package ethbeacon_rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	state_native "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"sort"
	"strconv"
	"toprelayer/relayer/toprelayer/ethtypes"
)

func (c *BeaconGrpcClient) GetFinalizedLightClientUpdateV2() (*LightClientUpdate, error) {
	return c.getFinalizedLightClientUpdate(false)
}

func (c *BeaconGrpcClient) GetFinalizedLightClientUpdateV2WithNextSyncCommittee() (*LightClientUpdate, error) {
	return c.getFinalizedLightClientUpdate(true)
}

func (c *BeaconGrpcClient) getFinalizedLightClientUpdate(useNextSyncCommittee bool) (*LightClientUpdate, error) {
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		return nil, err
	}
	attestedSlot, err := c.getAttestedSlotBeforeFinalizedSlot(finalizedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return nil, err
	}
	if GetPeriodForSlot(attestedSlot) != GetPeriodForSlot(finalizedSlot) {
		return nil, fmt.Errorf("Eth2TopRelayerV2 attestedSlot(%d) and finalizedSlot(%d) not in same period", attestedSlot, finalizedSlot)
	}
	lcu, err := c.GetFinalityLightClientUpdate(attestedSlot, useNextSyncCommittee)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getFinalityLightClientUpdate error:", err)
		return nil, err
	}
	return ConvertEth2LightClientUpdate(lcu), nil
}

func (c *BeaconGrpcClient) GetLightClientUpdateV2(period uint64) (*LightClientUpdate, error) {
	currFinalizedSlot := GetFinalizedForPeriod(period)
	attestedSlot, err := c.GetAttestedSlot(currFinalizedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return nil, err
	}
	lcu, err := c.GetFinalityLightClientUpdate(attestedSlot, true)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getFinalityLightClientUpdate error:", err)
		return nil, err
	}
	return ConvertEth2LightClientUpdate(lcu), nil
}

func (c *BeaconGrpcClient) GetNextSyncCommitteeUpdateV2(period uint64) (*SyncCommitteeUpdate, error) {
	currFinalizedSlot := GetFinalizedForPeriod(period)
	attestedSlot, err := c.GetAttestedSlot(currFinalizedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	if GetPeriodForSlot(attestedSlot) != GetPeriodForSlot(currFinalizedSlot) {
		return nil, fmt.Errorf("Eth2TopRelayerV2 GetNextSyncCommitteeUpdateV2 attestedSlot(%d) and finalizedSlot(%d) not in same period", attestedSlot, currFinalizedSlot)
	}
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetNextSyncCommitteeUpdateV2 attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)
	beaconState, err := c.GetBeaconState(strconv.FormatUint(attestedSlot, 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconState error:", err)
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

func (c *BeaconGrpcClient) GetAttestedSlot(lastFinalizedSlotOnNear uint64) (uint64, error) {
	nextFinalizedSlot := lastFinalizedSlotOnNear + ONE_EPOCH_IN_SLOTS
	attestedSlot := nextFinalizedSlot + 2*ONE_EPOCH_IN_SLOTS
	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return 0, err
	}
	return uint64(header.Slot), nil
}

func (c *BeaconGrpcClient) getAttestedSlotBeforeFinalizedSlot(finalizedSlot uint64) (uint64, error) {
	// the range of attestedSlot is [finalizedSlot - 32, finalizedSlot]
	attestedSlot := finalizedSlot - ONE_EPOCH_IN_SLOTS
	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return 0, err
	}
	return uint64(header.Slot), nil
}

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
	var state, err = state_native.InitializeFromProtoUnsafeCapella(beaconState)
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
	//finalizedBlockBodyHash, err := finalizedBlockBody.HashTreeRoot()
	blockHash := finalizedBlockBody.ExecutionPayload.GetBlockHash()
	var finalizedBlockBodyHash common.Hash
	copy(finalizedBlockBodyHash[:], blockHash[:])
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
			BeaconHeader:       finalityHeader,
			ExecutionBlockHash: finalizedBlockBodyHash,
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

func (c *BeaconGrpcClient) GetFinalityLightClientUpdate(attestedSlot uint64, useNextSyncCommittee bool) (*ethtypes.LightClientUpdate, error) {
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetFinalityLightClientUpdate attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)
	beaconState, err := c.GetBeaconState(strconv.FormatUint(attestedSlot, 10))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconState error:", err)
		return nil, err
	}
	finalityHash := beaconState.GetFinalizedCheckpoint().Root

	finalityHeader, err := c.GetBeaconBlockHeaderForBlockId(string(finalityHash))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetBeaconBlockHeaderForBlockId error:", err)
		return nil, err
	}

	finalitySlot := finalityHeader.Slot
	var finalityBeaconState *eth.BeaconStateCapella = nil
	if useNextSyncCommittee == true {
		finalityBeaconState, err = c.GetBeaconState(strconv.FormatUint(uint64(finalitySlot), 10))
		if err != nil {
			logger.Error("Eth2TopRelayerV2 GetBeaconState error:", err)
			return nil, err
		}
	}
	return c.getFinalityLightClientUpdateForState(attestedSlot, signatureSlot, beaconState, finalityBeaconState)
}
