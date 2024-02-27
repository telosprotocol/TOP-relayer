package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/api/client/beacon"
	state_native "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"google.golang.org/protobuf/proto"
	"net/http/httptest"
	"sort"
	"strconv"
	"time"
	"toprelayer/relayer/toprelayer/ethtypes"
	"toprelayer/rpc/ethereum/light_client"
)

type BeaconChainClient struct {
	*beacon.Client
}

func NewBeaconChainClient(httpUrl string) (*BeaconChainClient, error) {
	c, err := beacon.NewClient(httpUrl)
	if err != nil {
		return nil, err
	}
	return &BeaconChainClient{c}, nil
}

func (c *BeaconChainClient) GetBlindedSignedBeaconBlock(blockId beacon.StateOrBlockId) (interfaces.ReadOnlySignedBeaconBlock, error) {
	signedBeaconBlockSsz, err := c.GetBlock(context.Background(), blockId)
	if err != nil {
		logger.Error("GetBlock error:%s", err.Error())
		return nil, err
	}

	//logger.Info("GetSignedBeaconBlock blockId:%s,signedBeaconBlockSsz:%s", blockId, common.Bytes2Hex(signedBeaconBlockSsz))

	var signedBeaconBlockPb eth.SignedBeaconBlockDeneb
	err = signedBeaconBlockPb.UnmarshalSSZ(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("UnmarshalSSZ error:%s", err.Error())
		return nil, err
	}

	signedBeaconBlock, err := blocks.NewSignedBeaconBlock(&signedBeaconBlockPb)
	if err != nil {
		logger.Error("NewSignedBeaconBlock error:%s", err.Error())
		return nil, err
	}

	return signedBeaconBlock.ToBlinded()
}

func (c *BeaconChainClient) GetSignedBeaconBlock(blockId beacon.StateOrBlockId) (interfaces.SignedBeaconBlock, error) {
	signedBeaconBlockSsz, err := c.GetBlock(context.Background(), blockId)
	if err != nil {
		logger.Error("GetBlock error:%s", err.Error())
		return nil, err
	}

	//logger.Info("GetSignedBeaconBlock blockId:%s,signedBeaconBlockSsz:%s", blockId, common.Bytes2Hex(signedBeaconBlockSsz))

	var signedBeaconBlockPb eth.SignedBeaconBlockDeneb
	err = signedBeaconBlockPb.UnmarshalSSZ(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("UnmarshalSSZ error:%s", err.Error())
		return nil, err
	}

	signedBeaconBlock, err := blocks.NewSignedBeaconBlock(&signedBeaconBlockPb)
	if err != nil {
		logger.Error("NewSignedBeaconBlock error:%s", err.Error())
		return nil, err
	}

	return signedBeaconBlock, nil
}

func (c *BeaconChainClient) GetBeaconBlockBody(blockId beacon.StateOrBlockId) (interfaces.ReadOnlyBeaconBlockBody, error) {
	signedBeaconBlock, err := c.GetBlindedSignedBeaconBlock(blockId)
	if err != nil {
		return nil, err
	}

	return signedBeaconBlock.Block().Body(), nil
}

func (c *BeaconChainClient) GetBeaconBlockBodyFromSignedBeaconBlock(signedBeaconBlock interfaces.ReadOnlySignedBeaconBlock) interfaces.ReadOnlyBeaconBlockBody {
	return signedBeaconBlock.Block().Body()
}

func (c *BeaconChainClient) GetBeaconBlockBodyFromBeaconBlock(beaconBlock interfaces.ReadOnlyBeaconBlock) interfaces.ReadOnlyBeaconBlockBody {
	return beaconBlock.Body()
}

func (c *BeaconChainClient) GetBeaconBlockHeader(blockId beacon.StateOrBlockId) (*eth.BeaconBlockHeader, error) {
	signedBeaconBlock, err := c.GetBlindedSignedBeaconBlock(blockId)
	if err != nil {
		return nil, err
	}

	signedBeaconBlockHeader, err := signedBeaconBlock.Header()
	if err != nil {
		return nil, err
	}

	return signedBeaconBlockHeader.GetHeader(), nil
}

func (c *BeaconChainClient) GetLastSlotNumber() (primitives.Slot, error) {
	beaconBlockHeader, err := c.GetBeaconBlockHeader("head")
	if err != nil {
		logger.Error("GetBeaconBlockHeader error:", err)
		return 0, err
	}

	return beaconBlockHeader.GetSlot(), nil
}

func (c *BeaconChainClient) GetLastFinalizedSlotNumber() (primitives.Slot, error) {
	h, err := c.GetBeaconBlockHeader("finalized")
	if err != nil {
		logger.Error("GetBeaconBlockHeader error:", err)
		return 0, err
	}
	return h.GetSlot(), nil
}

func (c *BeaconChainClient) GetBlockNumberForSlot(slot primitives.Slot) (uint64, error) {
	beaconBlockBody, err := c.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(slot), 10)))
	if err != nil {
		return 0, err
	}
	executionPayload, err := beaconBlockBody.Execution()
	if err != nil {
		return 0, err
	}
	return executionPayload.BlockNumber(), nil
}

func (c *BeaconChainClient) getBeaconState(id primitives.Slot) (*eth.BeaconStateDeneb, error) {
	start := time.Now()
	defer func() {
		logger.Info("Slot:%s,getBeaconState time:%v", id, time.Since(start))
	}()

	writer := httptest.NewRecorder()
	writer.Body = &bytes.Buffer{}

	sszBeaconState, err := c.GetState(context.Background(), beacon.StateOrBlockId(strconv.FormatUint(uint64(id), 10)))
	if err != nil {
		logger.Error("GetState error:", err)
		return nil, err
	}

	var state eth.BeaconStateDeneb
	if err = state.UnmarshalSSZ(sszBeaconState); err != nil {
		logger.Error("UnmarshalSSZ error:", err)
		return nil, err
	}
	//return proto.Clone(&state).(*eth.BeaconState), nil
	return &state, nil
}

func (c *BeaconChainClient) GetNonEmptyBeaconBlockHeader(startSlot primitives.Slot) (*eth.BeaconBlockHeader, error) {
	lastSlot, err := c.GetLastSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return nil, err
	}

	for slot := startSlot; slot <= lastSlot; slot++ {
		if h, err := c.GetBeaconBlockHeader(beacon.StateOrBlockId(strconv.FormatUint(uint64(slot), 10))); err != nil {
			if IsErrorNoBlockForSlot(err) {
				logger.Info("GetBeaconBlockHeaderForBlockId slot(%d) error:%s", slot, err.Error())
				continue
			} else {
				logger.Error("GetBeaconBlockBodyForBlockId error:", err)
				return nil, err
			}
		} else {
			return h, nil
		}
	}
	return nil, fmt.Errorf("unable to get non empty beacon block in range [%d, %d)", startSlot, lastSlot)
}

func (c *BeaconChainClient) BeaconHeaderConvert(data *light_client.BeaconBlockHeaderData) (*light_client.BeaconBlockHeader, error) {
	slotVal, err := strconv.ParseUint(data.Beacon.Slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	indexVal, err := strconv.ParseUint(data.Beacon.ProposerIndex, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	slot := primitives.Slot(slotVal)
	index := primitives.ValidatorIndex(indexVal)

	h := new(light_client.BeaconBlockHeader)
	h.Slot = slot
	h.ProposerIndex = index
	h.BodyRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.Beacon.BodyRoot[2:]))
	h.ParentRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.Beacon.ParentRoot[2:]))
	h.StateRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.Beacon.StateRoot[2:]))
	return h, nil
}

func (c *BeaconChainClient) SyncAggregateConvert(data *light_client.SyncAggregateData) (*light_client.SyncAggregate, error) {
	aggregate := new(light_client.SyncAggregate)
	//aggregate.SyncCommitteeBits = data.SyncCommitteeBits
	aggregate.SyncCommitteeSignature = [fieldparams.BLSSignatureLength]byte(common.Hex2Bytes(data.SyncCommitteeSignature[2:]))
	return aggregate, nil
}

func (c *BeaconChainClient) CommitteeConvert(committee *light_client.SyncCommitteeData, branch []string) (*light_client.SyncCommitteeUpdate, error) {
	committeeUpdate := new(light_client.SyncCommitteeUpdate)

	nextCommittee := new(eth.SyncCommittee)
	nextCommittee.AggregatePubkey = common.Hex2Bytes(committee.AggregatePubkey[2:])
	for _, s := range committee.Pubkeys {
		nextCommittee.Pubkeys = append(nextCommittee.Pubkeys, common.Hex2Bytes(s[2:]))
	}
	committeeUpdate.NextSyncCommittee = nextCommittee

	for _, s := range branch {
		committeeUpdate.NextSyncCommitteeBranch = append(committeeUpdate.NextSyncCommitteeBranch, [fieldparams.RootLength]byte(common.Hex2Bytes(s[2:])))
	}
	return committeeUpdate, nil
}

func (c *BeaconChainClient) FinalizedUpdateConvert(header *light_client.BeaconBlockHeaderData, branch []string) (*light_client.FinalizedHeaderUpdate, error) {
	if len(header.ExecutionData.BlockHash) != len("0x")+fieldparams.RootLength*2 {
		err := fmt.Errorf("invalid execution block hash. hash:%s", header.ExecutionData.BlockHash)
		logger.Error("invalid execution hash:", err)
		return nil, err
	}

	for i, s := range header.ExecutionBranch {
		if len(s) != len("0x")+fieldparams.RootLength*2 {
			err := fmt.Errorf("invalid execution branch hash. index:%d hash:%s", i, s)
			logger.Error("invalid execution branch hash:", err)
			return nil, err
		}
	}

	for i, s := range branch {
		if len(s) != len("0x")+fieldparams.RootLength*2 {
			err := fmt.Errorf("invalid finality branch hash. index:%d hash:%s", i, s)
			logger.Error("invalid finality branch hash:", err)
			return nil, err
		}
	}

	update := new(light_client.FinalizedHeaderUpdate)

	for _, s := range branch {
		update.FinalityBranch = append(update.FinalityBranch, [fieldparams.RootLength]byte(common.Hex2Bytes(s[2:])))
	}

	headerUpdate := new(light_client.HeaderUpdate)
	h, err := c.BeaconHeaderConvert(header)
	if err != nil {
		logger.Error("BeaconHeaderConvert error:", err)
		return nil, err
	}
	//body, err := c.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(h.Slot), 10)))
	//if err != nil {
	//	logger.Error("GetBeaconBlockBodyForBlockId error:", err)
	//	return nil, err
	//}
	//
	//executionPayload, err := body.Execution()
	//if err != nil {
	//	logger.Error("GetBeaconBlockBodyForBlockId error:", err)
	//	return nil, err
	//}

	headerUpdate.BeaconHeader = h
	headerUpdate.ExecutionBlockHash = [fieldparams.RootLength]byte(common.Hex2Bytes(header.ExecutionData.BlockHash[2:]))
	headerUpdate.ExecutionHashBranch = make([][fieldparams.RootLength]byte, len(header.ExecutionBranch))
	for i, s := range header.ExecutionBranch {
		headerUpdate.ExecutionHashBranch[i] = [fieldparams.RootLength]byte(common.Hex2Bytes(s[2:]))
	}

	update.HeaderUpdate = headerUpdate
	return update, nil
}

func (c *BeaconChainClient) LightClientUpdateConvert(data *light_client.LightClientUpdateData) (*light_client.LightClientUpdate, error) {
	attestedHeader, err := c.BeaconHeaderConvert(data.AttestedHeader)
	if err != nil {
		logger.Error("BeaconHeaderConvert error:", err)
		return nil, err
	}
	aggregate, err := c.SyncAggregateConvert(data.SyncAggregate)
	if err != nil {
		logger.Error("SyncAggregateConvert error:", err)
		return nil, err
	}
	committeeUpdate, err := c.CommitteeConvert(data.NextSyncCommittee, data.NextSyncCommitteeBranch)
	if err != nil {
		logger.Error("CommitteeConvert error:", err)
		return nil, err
	}
	finalizedUpdate, err := c.FinalizedUpdateConvert(data.FinalizedHeader, data.FinalityBranch)
	if err != nil {
		logger.Error("FinalizedUpdateConvert error:", err)
		return nil, err
	}
	slotVal, err := strconv.ParseUint(data.SignatureSlot, 0, 64)
	if err != nil {
		logger.Error("ParseUint error:", err)
		return nil, err
	}
	update := new(light_client.LightClientUpdate)
	update.AttestedBeaconHeader = attestedHeader
	update.SyncAggregate = aggregate
	update.NextSyncCommitteeUpdate = committeeUpdate
	update.FinalityUpdate = finalizedUpdate
	update.SignatureSlot = primitives.Slot(slotVal)
	return update, nil
}

func (c *BeaconChainClient) GetLightClientUpdate(period uint64) (*light_client.LightClientUpdate, error) {
	str := fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=1", period)
	resp, err := c.Get(context.Background(), str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}

	var result []light_client.LightClientUpdateMsg
	if len(resp) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
	if err = json.Unmarshal(resp, &result); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(resp))
		logger.Error(err)
		return nil, err
	}
	if len(result) != 1 {
		err = fmt.Errorf("LightClientUpdateMsg size is not equal to 1")
		logger.Error("Unmarshal error:", err)
		return nil, err
	}
	return c.LightClientUpdateConvert(&result[0].Data)
}

func (c *BeaconChainClient) LightClientUpdateConvertNoCommitteeConvert(data *light_client.LightClientUpdateDataNoCommittee) (*light_client.LightClientUpdate, error) {
	attestedHeader, err := c.BeaconHeaderConvert(data.AttestedHeader)
	if err != nil {
		logger.Error("BeaconHeaderConvert error:", err)
		return nil, err
	}
	aggregate, err := c.SyncAggregateConvert(data.SyncAggregate)
	if err != nil {
		logger.Error("SyncAggregateConvert error:", err)
		return nil, err
	}
	finalizedUpdate, err := c.FinalizedUpdateConvert(data.FinalizedHeader, data.FinalityBranch)
	if err != nil {
		logger.Error("FinalizedUpdateConvert error:", err)
		return nil, err
	}
	slotVal, err := strconv.ParseUint(data.SignatureSlot, 0, 64)
	if err != nil {
		logger.Error("ParseUint error:", err)
		return nil, err
	}
	update := new(light_client.LightClientUpdate)
	update.AttestedBeaconHeader = attestedHeader
	update.SyncAggregate = aggregate
	update.NextSyncCommitteeUpdate = nil
	update.FinalityUpdate = finalizedUpdate
	update.SignatureSlot = primitives.Slot(slotVal)
	return update, nil
}

func (c *BeaconChainClient) GetAttestedSlot(lastFinalizedTopSlot primitives.Slot) (primitives.Slot, error) {
	attestedSlot := getAttestationSlot(lastFinalizedTopSlot)
	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient GetNonEmptyBeaconBlockHeader error:", err)
		return 0, err
	}
	return header.Slot, nil
}

func (c *BeaconChainClient) getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot primitives.Slot) (primitives.Slot, primitives.Slot, error) {
	currentAttestedSlot := attestedSlot
	for {
		h, err := c.GetNonEmptyBeaconBlockHeader(currentAttestedSlot + 1)
		if err != nil {
			logger.Error("BeaconChainClient GetNonEmptyBeaconBlockHeader error:", err)
			return 0, 0, err
		}

		signedBeaconBlock, err := c.GetBlindedSignedBeaconBlock(beacon.StateOrBlockId(strconv.FormatUint(uint64(h.Slot), 10)))
		if err != nil {
			logger.Error("BeaconChainClient GetBeaconBlock error:", err)
			return 0, 0, err
		}

		body := c.GetBeaconBlockBodyFromSignedBeaconBlock(signedBeaconBlock)
		syncAggregate, err := body.SyncAggregate()
		if err != nil {
			logger.Error("BeaconChainClient GetBeaconBlockBody error:", err)
			return 0, 0, err
		}
		syncCommitteeBitsSum := syncAggregate.SyncCommitteeBits.Count()

		if syncCommitteeBitsSum*3 < (64 * 8 * 2) {
			currentAttestedSlot = h.GetSlot()
			continue
		}
		if len(body.Attestations()) == 0 {
			currentAttestedSlot = h.GetSlot()
			continue
		}
		var attestedSlots []primitives.Slot
		for _, attestation := range body.Attestations() {
			attestedSlots = append(attestedSlots, attestation.GetData().Slot)
		}
		sort.Slice(attestedSlots, func(i, j int) bool { return attestedSlots[i] > attestedSlots[j] })
		for i, v := range attestedSlots {
			if (i == 0 || v != attestedSlots[i-1]) && v >= attestedSlot {
				currentAttestedSlot = v
				_, err = c.GetBeaconBlockHeader(beacon.StateOrBlockId(strconv.FormatUint(uint64(currentAttestedSlot), 10)))
				if err != nil {
					continue
				}
				return currentAttestedSlot, h.Slot, nil
			}
		}
		currentAttestedSlot = h.GetSlot()
	}
}

func (c *BeaconChainClient) constructFromBeaconBlockBody(beaconBlockBody interfaces.ReadOnlyBeaconBlockBody) (*ethtypes.ExecutionBlockProof, error) {
	executionPayload, err := beaconBlockBody.Execution()
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockBody error:", err)
		return nil, err
	}

	blockHash := executionPayload.BlockHash()
	var finalizedBlockBodyHash common.Hash
	copy(finalizedBlockBodyHash[:], blockHash[:])

	var beaconBlockMerkleTree, executionPayloadMerkleTree MerkleTreeNode
	if beaconBlockMerkleTree, err = BeaconBlockBodyMerkleTreeNew(beaconBlockBody); err != nil {
		return nil, err
	}
	if executionPayloadMerkleTree, err = ExecutionPayloadMerkleTreeNew(executionPayload); err != nil {
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

func (c *BeaconChainClient) getNextSyncCommittee(beaconState *eth.BeaconStateDeneb) (*ethtypes.SyncCommitteeUpdate, error) {
	beaconStateDeneb := proto.Clone(beaconState).(*eth.BeaconStateDeneb)

	if beaconStateDeneb == nil {
		return nil, nil
	}

	if beaconStateDeneb.GetNextSyncCommittee() == nil {
		logger.Error("BeaconChainClient NextSyncCommittee nil")
		return nil, errors.New("NextSyncCommittee nil")
	}
	var state, err = state_native.InitializeFromProtoDeneb(beaconStateDeneb)
	if err != nil {
		logger.Error("BeaconChainClient InitializeFromProtoUnsafeBellatrix error:", err)
		return nil, err
	}
	nextSyncCommitteeProofData, err := state.NextSyncCommitteeProof(context.Background())
	if err != nil {
		logger.Error("BeaconChainClient NextSyncCommitteeProof error:", err)
		return nil, err
	}

	nextSyncCommitteeProof, err := ethtypes.ConvertSliceBytes2SliceBytes32(nextSyncCommitteeProofData)
	if err != nil {
		logger.Error("BeaconChainClient ConvertSliceBytes2SliceBytes32 error:", err)
		return nil, err
	}

	update := &ethtypes.SyncCommitteeUpdate{
		NextSyncCommittee:       beaconStateDeneb.NextSyncCommittee,
		NextSyncCommitteeBranch: nextSyncCommitteeProof,
	}
	return update, nil
}

func (c *BeaconChainClient) getFinalityLightClientUpdateForState(attestedSlot, signatureSlot primitives.Slot, beaconState, finalityBeaconState *eth.BeaconStateDeneb) (*ethtypes.LightClientUpdate, error) {
	beaconBody, err := c.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(signatureSlot), 10)))
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}

	syncAggregate, err := beaconBody.SyncAggregate()
	if err != nil {
		logger.Error("BeaconChainClient SyncAggregate error:", err)
		return nil, err
	}

	attestedHeader, err := c.GetBeaconBlockHeader(beacon.StateOrBlockId(strconv.FormatUint(uint64(attestedSlot), 10)))
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockHeader error:", err)
		return nil, err
	}
	finalityHash := beaconState.FinalizedCheckpoint.Root
	signedBeaconBlock, err := c.GetSignedBeaconBlock(beacon.StateOrBlockId(finalityHash))
	if err != nil {
		logger.Error("BeaconChainClient GetSignedBeaconBlock error:", err)
		return nil, err
	}
	finalityHeader, err := signedBeaconBlock.Header()
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockHeader error:", err)
		return nil, err
	}
	finalizedBlockBody := signedBeaconBlock.Block().Body()
	//if err != nil {
	//	logger.Error("BeaconChainClient GetBeaconBlockBody error:", err)
	//	return nil, err
	//}
	executionBlockProof, err := c.constructFromBeaconBlockBody(finalizedBlockBody)
	if err != nil {
		logger.Error("BeaconChainClient constructFromBeaconBlockBody hash error:", err)
		return nil, err
	}
	if err != nil {
		logger.Error("BeaconChainClient finalizedBlockBody hash error:", err)
		return nil, err
	}
	state, err := state_native.InitializeFromProtoUnsafeDeneb(proto.Clone(beaconState).(*eth.BeaconStateDeneb))
	if err != nil {
		logger.Error("BeaconChainClient InitializeFromProtoUnsafeBellatrix error:", err)
		return nil, err
	}

	update := &ethtypes.LightClientUpdate{
		AttestedBeaconHeader: attestedHeader,
		SyncAggregate: &eth.SyncAggregate{
			SyncCommitteeBits:      syncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: syncAggregate.SyncCommitteeSignature,
		},
		SignatureSlot: uint64(signatureSlot),
	}
	proofData, err := state.FinalizedRootProof(context.Background())
	if err != nil {
		logger.Error("BeaconChainClient FinalizedRootProof error:", err)
		return nil, err
	}
	proof, err := ethtypes.ConvertSliceBytes2SliceBytes32(proofData)
	if err != nil {
		logger.Error("BeaconChainClient ConvertSliceBytes2SliceBytes32 error:", err)
		return nil, err
	}
	update.FinalizedUpdate = &ethtypes.FinalizedHeaderUpdate{
		HeaderUpdate: &ethtypes.HeaderUpdate{
			BeaconHeader:        finalityHeader.GetHeader(),
			ExecutionBlockHash:  executionBlockProof.BlockHash,
			ExecutionHashBranch: executionBlockProof.Proof,
		},
		FinalityBranch: proof,
	}
	if finalityBeaconState != nil {
		update.NextSyncCommitteeUpdate, err = c.getNextSyncCommittee(finalityBeaconState)
		if err != nil {
			logger.Error("BeaconChainClient getNextSyncCommittee error:", err)
			return nil, err
		}
	}
	return update, nil
}

func (c *BeaconChainClient) getFinalityLightClientUpdate(attestedSlot primitives.Slot, useNextSyncCommittee bool) (*ethtypes.LightClientUpdate, error) {
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetFinalityLightClientUpdate attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)
	beaconState, err := c.getBeaconState(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient getBeaconState error:", err)
		return nil, err
	}
	var finalityBeaconState *eth.BeaconStateDeneb = nil
	if useNextSyncCommittee == true {
		finalityBeaconState = beaconState
	}
	return c.getFinalityLightClientUpdateForState(attestedSlot, signatureSlot, beaconState, finalityBeaconState)
}

func (c *BeaconChainClient) getLightClientUpdateByFinalizedSlot(finalizedSlot primitives.Slot, useNextSyncCommittee bool) (*light_client.LightClientUpdate, error) {
	attestedSlot, err := c.GetAttestedSlot(finalizedSlot)
	if err != nil {
		logger.Error("BeaconChainClient GetNonEmptyBeaconBlockHeader error:", err)
		return nil, err
	}
	if GetPeriodForSlot(attestedSlot) != GetPeriodForSlot(finalizedSlot) {
		return nil, fmt.Errorf("BeaconChainClient attestedSlot(%d) and finalizedSlot(%d) not in same period", attestedSlot, finalizedSlot)
	}
	lcu, err := c.getFinalityLightClientUpdate(attestedSlot, useNextSyncCommittee)
	if err != nil {
		logger.Error("BeaconChainClient getFinalityLightClientUpdate error:", err)
		return nil, err
	}
	logger.Info("LightClientUpdate FinalizedSlot:%d,AttestedSlot:%d",
		lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot, lcu.AttestedBeaconHeader.Slot)
	return convertEth2LightClientUpdate(lcu), nil
}

func (c *BeaconChainClient) GetFinalizedLightClientUpdateByEthSlot(lastFinalizedEthSlot primitives.Slot) (*light_client.LightClientUpdate, error) {
	finalizedSlot, err := getBeforeSlotInSamePeriod(lastFinalizedEthSlot)
	if err != nil {
		return nil, err
	}
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, false)
}

func (c *BeaconChainClient) GetLastFinalizedLightClientUpdateV2WithNextSyncCommitteeByEthSlot(lastFinalizedEthSlot primitives.Slot) (*light_client.LightClientUpdate, error) {
	finalizedSlot, _ := getBeforeSlotInSamePeriod(lastFinalizedEthSlot)
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, true)
}

func (c *BeaconChainClient) GetLightClientUpdateV2(period uint64) (*light_client.LightClientUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getLightClientUpdateByFinalizedSlot(currFinalizedSlot, true)
}

func (c *BeaconChainClient) GetNextSyncCommitteeUpdate(period uint64) (*light_client.SyncCommitteeUpdate, error) {
	str := fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=1", period)
	resp, err := c.Get(context.Background(), str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}

	var result []light_client.LightClientUpdateMsg
	if len(resp) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
	if err = json.Unmarshal(resp, &result); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(resp))
		logger.Error(err)
		return nil, err
	}
	if len(result) != 1 {
		err = fmt.Errorf("LightClientUpdateMsg size is not equal to 1")
		logger.Error("Unmarshal error:", err)
		return nil, err
	}
	committeeUpdate, err := c.CommitteeConvert(result[0].Data.NextSyncCommittee, result[0].Data.NextSyncCommitteeBranch)
	if err != nil {
		logger.Error("CommitteeConvert error:", err)
		return nil, err
	}
	return committeeUpdate, nil
}

func (c *BeaconChainClient) GetLastFinalizedLightClientUpdateV2FinalizedSlot() (primitives.Slot, error) {
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		return 0, err
	}
	return getBeforeSlotInSamePeriod(finalizedSlot)
}

func (c *BeaconChainClient) GetLastFinalizedLightClientUpdateV2WithNextSyncCommittee() (*light_client.LightClientUpdate, error) {
	finalizedSlot, err := c.GetLastFinalizedLightClientUpdateV2FinalizedSlot()
	if err != nil {
		return nil, err
	}
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, true)
}

func (c *BeaconChainClient) getNextSyncCommitteeUpdateByFinalized(finalizedSlot primitives.Slot) (*light_client.SyncCommitteeUpdate, error) {
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
	beaconState, err := c.getBeaconState(primitives.Slot(attestedSlot))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getBeaconState error:", err)
		return nil, err
	}
	cu, err := c.getNextSyncCommittee(beaconState)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getNextSyncCommittee error:", err)
		return nil, err
	}
	return &light_client.SyncCommitteeUpdate{
		NextSyncCommittee:       cu.NextSyncCommittee,
		NextSyncCommitteeBranch: cu.NextSyncCommitteeBranch,
	}, nil
}

func (c *BeaconChainClient) GetNextSyncCommitteeUpdateV2(period uint64) (*light_client.SyncCommitteeUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getNextSyncCommitteeUpdateByFinalized(currFinalizedSlot)
}
