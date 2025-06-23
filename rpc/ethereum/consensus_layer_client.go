package ethereum

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"time"
	"toprelayer/relayer/toprelayer/ethtypes"
	"toprelayer/rpc/ethereum/light_client"

	"github.com/OffchainLabs/prysm/v6/api/client"
	"github.com/OffchainLabs/prysm/v6/api/client/beacon"
	"github.com/OffchainLabs/prysm/v6/api/server/structs"
	"github.com/OffchainLabs/prysm/v6/beacon-chain/state"
	"github.com/OffchainLabs/prysm/v6/beacon-chain/state/stateutil"
	fieldparams "github.com/OffchainLabs/prysm/v6/config/fieldparams"
	"github.com/OffchainLabs/prysm/v6/consensus-types/blocks"
	"github.com/OffchainLabs/prysm/v6/consensus-types/interfaces"
	"github.com/OffchainLabs/prysm/v6/consensus-types/primitives"
	"github.com/OffchainLabs/prysm/v6/container/trie"
	"github.com/OffchainLabs/prysm/v6/encoding/ssz/detect"
	eth "github.com/OffchainLabs/prysm/v6/proto/prysm/v1alpha1"
	"github.com/OffchainLabs/prysm/v6/runtime/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

type BeaconClient struct {
	*beacon.Client
	httpClient *http.Client
}

func NewBeaconClient(httpUrl string, opts ...client.ClientOpt) (*BeaconClient, error) {
	c, err := beacon.NewClient(httpUrl, opts...)
	if err != nil {
		return nil, err
	}
	return &BeaconClient{
		c,
		&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
	}, nil
}

func (c *BeaconClient) getBeconState(id primitives.Slot) (s state.BeaconState, err error) {
	start := time.Now()

	writer := httptest.NewRecorder()
	writer.Body = &bytes.Buffer{}

	sszBeaconState, err := c.GetState(context.Background(), beacon.StateOrBlockId(strconv.FormatUint(uint64(id), 10)))
	if err != nil {
		logger.Error("GetState error:", err)
		return nil, err
	}

	logger.Info("becon state len %d", len(sszBeaconState))
	unmarshaler, err := detect.FromState(sszBeaconState)
	if err != nil {
		logger.Error("From State error:", err)
		return nil, err
	}

	defer func() {
		logger.Info("Slot:%d, fork:%d, getBeaconState time:%v", id, unmarshaler.Fork, time.Since(start))
		sszBeaconState = nil
	}()

	beconstate, err := unmarshaler.UnmarshalBeaconState(sszBeaconState)
	if err != nil {
		logger.Error("UnmarshalSSZ error:", err)
		return nil, err
	}

	return beconstate, nil
}

func (c *BeaconClient) getNextSyncCommittee(beaconState state.BeaconState) (*ethtypes.SyncCommitteeUpdate, error) {
	if beaconState == nil {
		return nil, nil
	}

	nextSyncCommittee, err := beaconState.NextSyncCommittee()
	if err != nil {
		logger.Error("BeaconChainClient nextSyncCommittee error:", err)
		return nil, err
	}

	nextSyncCommitteeProofData, err := beaconState.NextSyncCommitteeProof(context.Background())
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
		NextSyncCommittee:       nextSyncCommittee,
		NextSyncCommitteeBranch: nextSyncCommitteeProof,
	}
	return update, nil
}

func (c *BeaconClient) getFinalityLightClientUpdateForState(attestedSlot, signatureSlot primitives.Slot, beaconState, finalityBeaconState state.BeaconState) (*ethtypes.LightClientUpdate, error) {
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

	finalityHash := beaconState.FinalizedCheckpoint().GetRoot()
	if err != nil {
		logger.Error("BeaconChainClient FinalizedCheckpoint error:", err)
		return nil, err
	}

	var hash32 [32]byte
	copy(hash32[:], finalityHash[:32])
	signedBeaconBlock, err := c.GetSignedBeaconBlock(beacon.IdFromRoot(hash32))
	if err != nil {
		logger.Error("BeaconChainClient GetSignedBeaconBlock error:", err)
		return nil, err
	}

	logger.Info("BeaconChainClient GetSignedBeaconBlock %x", hash32)

	finalityHeader, err := signedBeaconBlock.Header()
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockHeader error:", err)
		return nil, err
	}
	finalizedBlockBody := signedBeaconBlock.Block().Body()
	executionBlockProof, err := c.constructFromBeaconBlockBody(finalizedBlockBody)
	if err != nil {
		logger.Error("BeaconChainClient constructFromBeaconBlockBody hash error:", err)
		return nil, err
	}

	logger.Info("BeaconChainClient constructFromBeaconBlockBody")

	update := &ethtypes.LightClientUpdate{
		AttestedBeaconHeader: attestedHeader,
		SyncAggregate: &eth.SyncAggregate{
			SyncCommitteeBits:      syncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: syncAggregate.SyncCommitteeSignature,
		},
		SignatureSlot: uint64(signatureSlot),
	}

	proofData, err := beaconState.FinalizedRootProof(context.Background())
	if err != nil {
		logger.Error("BeaconChainClient FinalizedRootProof error:", err)
		return nil, err
	}

	logger.Info("BeaconChainClient FinalizedRootProof")
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

func unmarshalSignedBlock(signedBeaconBlockSsz []byte) (interfaces.SignedBeaconBlock, error) {
	unmarshaler, err := detect.FromBlock(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("From Block error:", err)
		return nil, err
	}

	signedBeaconBlock, err := unmarshaler.UnmarshalBeaconBlock(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("unmarshal beacon block error:", err)
		return nil, err
	}

	return signedBeaconBlock, nil
}

func (c *BeaconClient) GetBlindedSignedBeaconBlock(blockId beacon.StateOrBlockId) (interfaces.ReadOnlySignedBeaconBlock, error) {
	signedBeaconBlockSsz, err := c.GetBlock(context.Background(), blockId)
	if err != nil {
		logger.Error("GetBlock error:%s", err.Error())
		return nil, err
	}

	defer func() {
		signedBeaconBlockSsz = nil
	}()

	//logger.Info("GetSignedBeaconBlock blockId:%s,signedBeaconBlockSsz:%s", blockId, common.Bytes2Hex(signedBeaconBlockSsz))

	signedBeaconBlock, err := unmarshalSignedBlock(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("unmarshalSignedBlock error:%s", err.Error())
		return nil, err
	}

	return signedBeaconBlock.ToBlinded()
}

func (c *BeaconClient) GetSignedBeaconBlock(blockId beacon.StateOrBlockId) (interfaces.SignedBeaconBlock, error) {
	signedBeaconBlockSsz, err := c.GetBlock(context.Background(), blockId)
	if err != nil {
		logger.Error("GetBlock error:%s", err.Error())
		return nil, err
	}

	defer func() {
		signedBeaconBlockSsz = nil
	}()

	//logger.Info("GetSignedBeaconBlock blockId:%s,signedBeaconBlockSsz:%s", blockId, common.Bytes2Hex(signedBeaconBlockSsz))

	signedBeaconBlock, err := unmarshalSignedBlock(signedBeaconBlockSsz)
	if err != nil {
		logger.Error("unmarshalSignedBlock error:%s", err.Error())
		return nil, err
	}

	return signedBeaconBlock, nil
}

func (c *BeaconClient) GetBeaconBlockBody(blockId beacon.StateOrBlockId) (interfaces.ReadOnlyBeaconBlockBody, error) {
	signedBeaconBlock, err := c.GetSignedBeaconBlock(blockId)
	if err != nil {
		return nil, err
	}

	return signedBeaconBlock.Block().Body(), nil
}

func (c *BeaconClient) GetBeaconBlockBodyFromSignedBeaconBlock(signedBeaconBlock interfaces.ReadOnlySignedBeaconBlock) interfaces.ReadOnlyBeaconBlockBody {
	return signedBeaconBlock.Block().Body()
}

func (c *BeaconClient) GetBeaconBlockBodyFromBeaconBlock(beaconBlock interfaces.ReadOnlyBeaconBlock) interfaces.ReadOnlyBeaconBlockBody {
	return beaconBlock.Body()
}

func (c *BeaconClient) GetBeaconBlockHeader(blockId beacon.StateOrBlockId) (*eth.BeaconBlockHeader, error) {
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

func (c *BeaconClient) GetLastSlotNumber() (primitives.Slot, error) {
	beaconBlockHeader, err := c.GetBeaconBlockHeader("head")
	if err != nil {
		logger.Error("GetBeaconBlockHeader error:", err)
		return 0, err
	}

	return beaconBlockHeader.GetSlot(), nil
}

func (c *BeaconClient) GetLastFinalizedSlotNumber() (primitives.Slot, error) {
	h, err := c.GetBeaconBlockHeader("finalized")
	if err != nil {
		logger.Error("GetBeaconBlockHeader error:", err)
		return 0, err
	}
	return h.GetSlot(), nil
}

func (c *BeaconClient) GetBlockNumberForSlot(slot primitives.Slot) (uint64, error) {
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

func (c *BeaconClient) GetNonEmptyBeaconBlockHeader(startSlot primitives.Slot) (*eth.BeaconBlockHeader, error) {
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

func (c *BeaconClient) BeaconHeaderConvert(data *light_client.BeaconBlockHeaderJson) (*light_client.BeaconBlockHeader, error) {
	slotVal, err := strconv.ParseUint(data.BeaconJson.Slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	indexVal, err := strconv.ParseUint(data.BeaconJson.ProposerIndex, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	slot := primitives.Slot(slotVal)
	index := primitives.ValidatorIndex(indexVal)

	h := new(light_client.BeaconBlockHeader)
	h.Slot = slot
	h.ProposerIndex = index
	h.BodyRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.BeaconJson.BodyRoot[2:]))
	h.ParentRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.BeaconJson.ParentRoot[2:]))
	h.StateRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.BeaconJson.StateRoot[2:]))
	return h, nil
}

func (c *BeaconClient) PrysmBeaconBlockHeaderConvert(data *light_client.PrysmBeaconBlockHeaderJson) (*light_client.BeaconBlockHeader, error) {
	slotVal, err := strconv.ParseUint(data.Slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	indexVal, err := strconv.ParseUint(data.ProposerIndex, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}

	slot := primitives.Slot(slotVal)
	index := primitives.ValidatorIndex(indexVal)

	h := new(light_client.BeaconBlockHeader)
	h.Slot = slot
	h.ProposerIndex = index
	h.BodyRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.BodyRoot[2:]))
	h.ParentRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.ParentRoot[2:]))
	h.StateRoot = [fieldparams.RootLength]byte(common.Hex2Bytes(data.StateRoot[2:]))
	return h, nil
}

func (c *BeaconClient) SyncAggregateConvert(data *light_client.SyncAggregateJson) (*light_client.SyncAggregate, error) {
	aggregate := new(light_client.SyncAggregate)
	aggregate.SyncCommitteeBits = [fieldparams.SyncAggregateSyncCommitteeBytesLength]byte(common.Hex2Bytes(data.SyncCommitteeBits[2:]))
	aggregate.SyncCommitteeSignature = [fieldparams.BLSSignatureLength]byte(common.Hex2Bytes(data.SyncCommitteeSignature[2:]))
	return aggregate, nil
}

func (c *BeaconClient) CommitteeConvert(committee *light_client.SyncCommitteeJson, branch []string) (*light_client.SyncCommitteeUpdate, error) {
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

func (c *BeaconClient) FinalizedUpdateConvert(header *light_client.BeaconBlockHeaderJson, branch []string) (*light_client.FinalizedHeaderUpdate, error) {
	if len(header.ExecutionJson.BlockHash) != len("0x")+fieldparams.RootLength*2 {
		err := fmt.Errorf("invalid execution block hash. hash:%s", header.ExecutionJson.BlockHash)
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
	headerUpdate.ExecutionBlockHash = [fieldparams.RootLength]byte(common.Hex2Bytes(header.ExecutionJson.BlockHash[2:]))
	headerUpdate.ExecutionHashBranch = make([][fieldparams.RootLength]byte, len(header.ExecutionBranch))
	for i, s := range header.ExecutionBranch {
		headerUpdate.ExecutionHashBranch[i] = [fieldparams.RootLength]byte(common.Hex2Bytes(s[2:]))
	}

	update.HeaderUpdate = headerUpdate
	return update, nil
}

func (c *BeaconClient) PrysmFianlityUpdateConvert(header *light_client.PrysmBeaconBlockHeaderJson, branch []string) (*light_client.FinalizedHeaderUpdate, error) {
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
	finalizedHeader, err := c.PrysmBeaconBlockHeaderConvert(header)
	if err != nil {
		logger.Error("BeaconHeaderConvert error:", err)
		return nil, err
	}

	finalizedBlockBody, err := c.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(finalizedHeader.Slot), 10)))
	if err != nil {
		logger.Error("GetBeaconBlockBody error:", err)
		return nil, err
	}

	finalizedBlockExeuctionDataProof, err := c.constructFromBeaconBlockBody(finalizedBlockBody)
	if err != nil {
		logger.Error("constructFromBeaconBlockBody error:", err)
		return nil, err
	}

	headerUpdate.BeaconHeader = finalizedHeader
	headerUpdate.ExecutionBlockHash = finalizedBlockExeuctionDataProof.BlockHash
	headerUpdate.ExecutionHashBranch = ethtypes.ConvertSliceHash2SliceBytes32(finalizedBlockExeuctionDataProof.Proof)

	update.HeaderUpdate = headerUpdate
	return update, nil
}

func (c *BeaconClient) PrysmLightClientUpdateDataConvert(data *light_client.PrysmLightClientUpdateDataJson, hasNextSyncCommittee bool) (*light_client.LightClientUpdate, error) {
	attestedHeader, err := c.PrysmBeaconBlockHeaderConvert(data.AttestedHeader)
	if err != nil {
		logger.Error("BeaconHeaderConvert error:", err)
		return nil, err
	}
	aggregate, err := c.SyncAggregateConvert(data.SyncAggregate)
	if err != nil {
		logger.Error("SyncAggregateConvert error:", err)
		return nil, err
	}

	var committeeUpdate *light_client.SyncCommitteeUpdate = nil
	if hasNextSyncCommittee {
		committeeUpdate, err = c.CommitteeConvert(data.NextSyncCommittee, data.NextSyncCommitteeBranch)
		if err != nil {
			logger.Error("CommitteeConvert error:", err)
			return nil, err
		}
	}
	finalizedUpdate, err := c.PrysmFianlityUpdateConvert(data.FinalizedHeader, data.FinalityBranch)
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

func (c *BeaconClient) GetPrysmLightClientUpdate(period uint64) (*light_client.LightClientUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%d&count=1", c.Client.NodeURL(), period)
	resp, err := c.httpClient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("close http error:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("io.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}

	var result []structs.LightClientUpdateResponse
	if err = json.Unmarshal(body, &result); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}
	if len(result) != 1 {
		err = fmt.Errorf("LightClientUpdateMsg size is not equal to 1")
		logger.Error("Unmarshal error:", err)
		return nil, err
	}

	var lightClient light_client.PrysmLightClientUpdateDataJson

	var attestedHeader structs.LightClientHeader
	if err = json.Unmarshal(result[0].Data.AttestedHeader, &attestedHeader); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}

	var finalizedHeader structs.LightClientHeader
	if err = json.Unmarshal(result[0].Data.FinalizedHeader, &finalizedHeader); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}

	lightClient.AttestedHeader = (*light_client.PrysmBeaconBlockHeaderJson)(attestedHeader.Beacon)
	lightClient.FinalizedHeader = (*light_client.PrysmBeaconBlockHeaderJson)(finalizedHeader.Beacon)
	lightClient.FinalityBranch = result[0].Data.FinalityBranch
	lightClient.SyncAggregate = (*light_client.SyncAggregateJson)(result[0].Data.SyncAggregate)
	lightClient.NextSyncCommittee = (*light_client.SyncCommitteeJson)(result[0].Data.NextSyncCommittee)
	lightClient.NextSyncCommitteeBranch = result[0].Data.NextSyncCommitteeBranch
	lightClient.SignatureSlot = result[0].Data.SignatureSlot

	return c.PrysmLightClientUpdateDataConvert(&lightClient, true)
}

func (c *BeaconClient) GetPrysmFinalityLightClientUpdate() (*light_client.LightClientUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/finality_update", c.Client.NodeURL())
	resp, err := c.httpClient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("close http error:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("io.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}

	var result structs.LightClientFinalityUpdateResponse
	if err = json.Unmarshal(body, &result); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}

	var lightClient light_client.PrysmLightClientUpdateDataJson

	var attestedHeader structs.LightClientHeader
	if err = json.Unmarshal(result.Data.AttestedHeader, &attestedHeader); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}

	var finalizedHeader structs.LightClientHeader
	if err = json.Unmarshal(result.Data.FinalizedHeader, &finalizedHeader); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err)
		return nil, err
	}

	lightClient.AttestedHeader = (*light_client.PrysmBeaconBlockHeaderJson)(attestedHeader.Beacon)
	lightClient.FinalizedHeader = (*light_client.PrysmBeaconBlockHeaderJson)(finalizedHeader.Beacon)
	lightClient.FinalityBranch = result.Data.FinalityBranch
	lightClient.SyncAggregate = (*light_client.SyncAggregateJson)(result.Data.SyncAggregate)
	lightClient.SignatureSlot = result.Data.SignatureSlot
	return c.PrysmLightClientUpdateDataConvert(&lightClient, false)
}

//func (c *BeaconClient) LightClientUpdateConvertNoCommitteeConvert(data *light_client.LightClientUpdateDataNoCommittee) (*light_client.LightClientUpdate, error) {
//	attestedHeader, err := c.BeaconHeaderConvert(data.AttestedHeader)
//	if err != nil {
//		logger.Error("BeaconHeaderConvert error:", err)
//		return nil, err
//	}
//	aggregate, err := c.SyncAggregateConvert(data.SyncAggregate)
//	if err != nil {
//		logger.Error("SyncAggregateConvert error:", err)
//		return nil, err
//	}
//	finalizedUpdate, err := c.FinalizedUpdateConvert(data.FinalizedHeader, data.FinalityBranch)
//	if err != nil {
//		logger.Error("FinalizedUpdateConvert error:", err)
//		return nil, err
//	}
//	slotVal, err := strconv.ParseUint(data.SignatureSlot, 0, 64)
//	if err != nil {
//		logger.Error("ParseUint error:", err)
//		return nil, err
//	}
//	update := new(light_client.LightClientUpdate)
//	update.AttestedBeaconHeader = attestedHeader
//	update.SyncAggregate = aggregate
//	update.NextSyncCommitteeUpdate = nil
//	update.FinalityUpdate = finalizedUpdate
//	update.SignatureSlot = primitives.Slot(slotVal)
//	return update, nil
//}

func (c *BeaconClient) GetAttestedSlot(lastFinalizedTopSlot primitives.Slot) (primitives.Slot, error) {
	attestedSlot := getAttestationSlot(lastFinalizedTopSlot)
	header, err := c.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient GetNonEmptyBeaconBlockHeader error:", err)
		return 0, err
	}
	return header.Slot, nil
}

func (c *BeaconClient) getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot primitives.Slot) (primitives.Slot, primitives.Slot, error) {
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

		if syncCommitteeBitsSum*3 < (syncAggregate.SyncCommitteeBits.Len() * 2) {
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

func (c *BeaconClient) constructFromBeaconBlockBody(beaconBlockBody interfaces.ReadOnlyBeaconBlockBody) (*ethtypes.ExecutionBlockProof, error) {
	var l2ExecutionPayloadProofSize uint64 = 4
	if beaconBlockBody.Version() >= version.Deneb {
		l2ExecutionPayloadProofSize = 5
	}

	executionPayload, err := beaconBlockBody.Execution()
	if err != nil {
		logger.Error("BeaconChainClient GetBeaconBlockBody error:", err)
		return nil, err
	}

	blockHash := executionPayload.BlockHash()
	var finalizedBlockBodyHash common.Hash
	copy(finalizedBlockBodyHash[:], blockHash[:])

	var executionPayloadMerkleTree MerkleTreeNode
	if beaconBlockBody.Version() >= version.Deneb {
		if executionPayloadMerkleTree, err = ExecutionPayloadMerkleTreeCancun(executionPayload); err != nil {
			logger.Error("BeaconClient ExecutionPayloadMerkleTreeNew error: ", err)
			return nil, err
		}
	} else {
		if executionPayloadMerkleTree, err = ExecutionPayloadMerkleTreeShanghai(executionPayload); err != nil {
			return nil, err
		}
	}

	//proof
	fieldRoots, err := blocks.ComputeBlockBodyFieldRoots(context.Background(), beaconBlockBody.(*blocks.BeaconBlockBody))
	if err != nil {
		return nil, err
	}

	fieldRootsTrie := stateutil.Merkleize(fieldRoots)
	l1ExecutionPayloadProof := trie.ProofFromMerkleLayers(fieldRootsTrie, int(L1BeaconBlockBodyTreeExecutionPayloadIndex))
	l1ExecutionPayloadProofByte32, err := ethtypes.ConvertSliceBytes2SliceBytes32(l1ExecutionPayloadProof)
	if err != nil {
		return nil, err
	}

	_, blockProof := generateProof(executionPayloadMerkleTree, L2ExecutionPayloadTreeExecutionBlockIndex, l2ExecutionPayloadProofSize)

	blockProof = append(blockProof, l1ExecutionPayloadProofByte32...)
	return &ethtypes.ExecutionBlockProof{
		BlockHash: finalizedBlockBodyHash,
		Proof:     ethtypes.ConvertSliceBytes2Hash(blockProof),
	}, nil
}

func (c *BeaconClient) getFinalityLightClientUpdate(attestedSlot primitives.Slot, useNextSyncCommittee bool) (*ethtypes.LightClientUpdate, error) {
	attestedSlot, signatureSlot, err := c.getAttestedSlotWithEnoughSyncCommitteeBitsSum(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient getAttestedSlotWithEnoughSyncCommitteeBitsSum error:", err)
		return nil, err
	}
	logger.Info("GetFinalityLightClientUpdate attestedSlot:%d, signatureSlot:%d", attestedSlot, signatureSlot)

	beaconState, err := c.getBeconState(attestedSlot)
	if err != nil {
		logger.Error("BeaconChainClient getBeaconState error:", err)
		return nil, err
	}
	var finalityBeaconState state.BeaconState = nil
	if useNextSyncCommittee == true {
		finalityBeaconState = beaconState
	}

	return c.getFinalityLightClientUpdateForState(attestedSlot, signatureSlot, beaconState, finalityBeaconState)
}

func (c *BeaconClient) getLightClientUpdateByFinalizedSlot(finalizedSlot primitives.Slot, useNextSyncCommittee bool) (*light_client.LightClientUpdate, error) {
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

func (c *BeaconClient) GetLightClientUpdateV2(period uint64) (*light_client.LightClientUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getLightClientUpdateByFinalizedSlot(currFinalizedSlot, true)
}

func (c *BeaconClient) GetNextSyncCommitteeUpdate(period uint64) (*light_client.SyncCommitteeUpdate, error) {
	str := fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=1", period)
	resp, err := c.Get(context.Background(), str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}

	var result []structs.LightClientUpdateResponse
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

	nextSyncCommittee := (*light_client.SyncCommitteeJson)(result[0].Data.NextSyncCommittee)
	nextSyncCommitteeBranch := result[0].Data.NextSyncCommitteeBranch
	committeeUpdate, err := c.CommitteeConvert(nextSyncCommittee, nextSyncCommitteeBranch)
	if err != nil {
		logger.Error("CommitteeConvert error:", err)
		return nil, err
	}
	return committeeUpdate, nil
}

func (c *BeaconClient) GetLastFinalizedLightClientUpdateV2FinalizedSlot() (primitives.Slot, error) {
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		return 0, err
	}
	return getBeforeSlotInSamePeriod(finalizedSlot)
}

func (c *BeaconClient) GetLastFinalizedLightClientUpdateV2WithNextSyncCommittee() (*light_client.LightClientUpdate, error) {
	finalizedSlot, err := c.GetLastFinalizedLightClientUpdateV2FinalizedSlot()
	if err != nil {
		return nil, err
	}
	return c.getLightClientUpdateByFinalizedSlot(finalizedSlot, true)
}

func (c *BeaconClient) getNextSyncCommitteeUpdateByFinalized(finalizedSlot primitives.Slot) (*light_client.SyncCommitteeUpdate, error) {
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
	beaconState, err := c.getBeconState(primitives.Slot(attestedSlot))
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getBeaconState error:", err)
		return nil, err
	}

	committee, err := c.getNextSyncCommittee(beaconState)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 getNextSyncCommittee error:", err)
		return nil, err
	}

	return &light_client.SyncCommitteeUpdate{
		NextSyncCommittee:       committee.NextSyncCommittee,
		NextSyncCommitteeBranch: committee.NextSyncCommitteeBranch,
	}, nil
}

func (c *BeaconClient) GetNextSyncCommitteeUpdateV2(period uint64) (*light_client.SyncCommitteeUpdate, error) {
	currFinalizedSlot := GetFinalizedSlotForPeriod(period)
	return c.getNextSyncCommitteeUpdateByFinalized(currFinalizedSlot)
}
