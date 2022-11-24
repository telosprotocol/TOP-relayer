package beaconrpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	pb "github.com/prysmaticlabs/prysm/v3/proto/eth/service"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/eth/v1"
	v2 "github.com/prysmaticlabs/prysm/v3/proto/eth/v2"
	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	SLOTS_PER_EPOCH   = 32
	EPOCHS_PER_PERIOD = 256

	ERROR_NO_BLOCK_FOR_SLOT = "not find requested block"
)

type BeaconGrpcClient struct {
	client      pb.BeaconChainClient
	debugclient pb.BeaconDebugClient

	httpclient *http.Client
	httpurl    string
}

func NewBeaconGrpcClient(grpcUrl, httpUrl string) (*BeaconGrpcClient, error) {
	grpc, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("grpc.Dial error:", err)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &BeaconGrpcClient{
		client:      pb.NewBeaconChainClient(grpc),
		debugclient: pb.NewBeaconDebugClient(grpc),
		httpclient:  &http.Client{Transport: tr},
		httpurl:     httpUrl,
	}
	return c, nil
}

func IsErrorNoBlockForSlot(err error) bool {
	return strings.Contains(err.Error(), ERROR_NO_BLOCK_FOR_SLOT)
}

func (c *BeaconGrpcClient) GetBeaconBlockBodyForBlockId(id string) (*v2.BeaconBlockBodyBellatrix, error) {
	resp, err := c.client.GetBlockV2(context.Background(), &v2.BlockRequestV2{BlockId: []byte(id)})
	if err != nil {
		logger.Error("GetBlockV2 id %v error %v", id, err)
		return nil, err
	}
	signedBlock, ok := resp.Data.Message.(*v2.SignedBeaconBlockContainer_BellatrixBlock)
	if !ok {
		return nil, errors.New("resp.data.message error")
	}
	return signedBlock.BellatrixBlock.GetBody(), nil
}

func (c *BeaconGrpcClient) GetBeaconBlockHeaderForBlockId(id string) (*eth.BeaconBlockHeader, error) {
	resp, err := c.client.GetBlockHeader(context.Background(), &v1.BlockRequest{BlockId: []byte(id)})
	if err != nil {
		logger.Error("GetBlockHeader error:", err)
		return nil, err
	}
	header := new(eth.BeaconBlockHeader)
	header.Slot = resp.Data.Header.Message.Slot
	header.ProposerIndex = resp.Data.Header.Message.ProposerIndex
	header.BodyRoot = resp.Data.Header.Message.BodyRoot
	header.ParentRoot = resp.Data.Header.Message.ParentRoot
	header.StateRoot = resp.Data.Header.Message.StateRoot
	return header, nil
}

func (c *BeaconGrpcClient) GetLastSlotNumber() (uint64, error) {
	h, err := c.GetBeaconBlockHeaderForBlockId("head")
	if err != nil {
		logger.Error("GetBeaconBlockHeaderForBlockId error:", err)
		return 0, err
	}
	return uint64(h.Slot), nil
}

func (c *BeaconGrpcClient) GetLastFinalizedSlotNumber() (uint64, error) {
	h, err := c.GetBeaconBlockHeaderForBlockId("finalized")
	if err != nil {
		logger.Error("GetBeaconBlockHeaderForBlockId error:", err)
		return 0, err
	}
	return uint64(h.Slot), nil
}

func (c *BeaconGrpcClient) GetBlockNumberForSlot(slot uint64) (uint64, error) {
	b, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(slot, 10))
	if err != nil {
		logger.Error("GetBeaconBlockBodyForBlockId error:", err)
		return 0, err
	}
	return b.GetExecutionPayload().BlockNumber, nil
}

func (c *BeaconGrpcClient) GetBlockHashForSlot(slot uint64) (common.Hash, error) {
	b, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(slot, 10))
	if err != nil {
		logger.Error("GetBeaconBlockBodyForBlockId slot %v error %v", slot, err)
		return common.Hash{}, err
	}
	return common.BytesToHash(b.GetExecutionPayload().BlockHash), nil
}

func GetPeriodForSlot(slot uint64) uint64 {
	return (slot / (SLOTS_PER_EPOCH * EPOCHS_PER_PERIOD))
}

func (c *BeaconGrpcClient) GetBeaconState(id string) (*eth.BeaconStateBellatrix, error) {
	resp, err := c.debugclient.GetBeaconStateSSZV2(context.Background(), &v2.BeaconStateRequestV2{StateId: []byte(id)})
	if err != nil {
		logger.Error("GetBeaconStateV2 error:", err)
		return nil, err
	}
	var state eth.BeaconStateBellatrix
	err = state.UnmarshalSSZ(resp.Data)
	if err != nil {
		logger.Error("UnmarshalSSZ error:", err)
		return nil, err
	}
	return &state, nil
}

func (c *BeaconGrpcClient) GetCheckpointRoot(id string) (*v1.Checkpoint, error) {
	resp, err := c.client.GetFinalityCheckpoints(context.Background(), &v1.StateRequest{StateId: []byte(id)})
	if err != nil {
		logger.Error("GetFinalityCheckpoints error:", err)
		return nil, err
	}

	return resp.Data.GetFinalized(), nil
}

func (c *BeaconGrpcClient) GetNonEmptyBeaconBlockHeader(startSlot uint64) (*eth.BeaconBlockHeader, error) {
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return nil, err
	}
	for slot := startSlot; slot < finalizedSlot; slot += 1 {
		h, err := c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(slot, 10))
		if err != nil {
			logger.Error("GetBeaconBlockBodyForBlockId error:", err)
			return nil, err
		}
		return h, nil
	}
	return nil, fmt.Errorf("unable to get non empty beacon block in range [%d, %d)", startSlot, finalizedSlot)
}

func (c *BeaconGrpcClient) GetLightClientUpdate(period uint64) (*LightClientUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%d&count=1", c.httpurl, period)
	resp, err := c.httpclient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result LightClientUpdateMsg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return nil, err
	}
	return c.LightClientUpdateConvert(&result.Data[0])
}

func (c *BeaconGrpcClient) GetNextSyncCommitteeUpdate(period uint64) (*SyncCommitteeUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%d&count=1", c.httpurl, period)
	resp, err := c.httpclient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result LightClientUpdateMsg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return nil, err
	}
	committeeUpdate, err := c.CommitteeConvert(result.Data[0].NextSyncCommittee, result.Data[0].NextSyncCommitteeBranch)
	if err != nil {
		logger.Error("CommitteeConvert error:", err)
		return nil, err
	}
	return committeeUpdate, nil
}

func (c *BeaconGrpcClient) GetFinalizedLightClientUpdate() (*LightClientUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/finality_update", c.httpurl)
	resp, err := c.httpclient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result LightClientUpdateNoCommitteeMsg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Unmarshal error:", err)
		return nil, err
	}
	return c.LightClientUpdateConvertNoCommitteeConvert(&result.Data)
}

type BeaconBlockHeaderData struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type SyncAggregateData struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}

type SyncCommitteeData struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
}

type LightClientUpdateDataNoCommittee struct {
	AttestedHeader  *BeaconBlockHeaderData `json:"attested_header"`
	FinalizedHeader *BeaconBlockHeaderData `json:"finalized_header"`
	FinalityBranch  []string               `json:"finality_branch"`
	SyncAggregate   *SyncAggregateData     `json:"sync_aggregate"`
	SignatureSlot   string                 `json:"signature_slot"`
}

type LightClientUpdateData struct {
	AttestedHeader          *BeaconBlockHeaderData `json:"attested_header"`
	FinalizedHeader         *BeaconBlockHeaderData `json:"finalized_header"`
	FinalityBranch          []string               `json:"finality_branch"`
	SyncAggregate           *SyncAggregateData     `json:"sync_aggregate"`
	NextSyncCommittee       *SyncCommitteeData     `json:"next_sync_committee"`
	NextSyncCommitteeBranch []string               `json:"next_sync_committee_branch"`
	SignatureSlot           string                 `json:"signature_slot"`
}

type LightClientUpdateNoCommitteeMsg struct {
	Data LightClientUpdateDataNoCommittee `json:"data"`
}

type LightClientUpdateMsg struct {
	Data []LightClientUpdateData `json:"data"`
}

type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    []byte
	StateRoot     []byte
	BodyRoot      []byte
}

func (h *BeaconBlockHeader) Encode() ([]byte, error) {
	b1, err := rlp.EncodeToBytes(h.Slot)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(h.ProposerIndex)
	if err != nil {
		return nil, err
	}
	b3, err := rlp.EncodeToBytes(h.ParentRoot)
	if err != nil {
		return nil, err
	}
	b4, err := rlp.EncodeToBytes(h.StateRoot)
	if err != nil {
		return nil, err
	}
	b5, err := rlp.EncodeToBytes(h.BodyRoot)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
	rlpBytes = append(rlpBytes, b4...)
	rlpBytes = append(rlpBytes, b5...)
	return rlpBytes, nil
}

type SyncAggregate struct {
	SyncCommitteeBits      string
	SyncCommitteeSignature []byte
}

func (s *SyncAggregate) Encode() ([]byte, error) {
	b1, err := rlp.EncodeToBytes(s.SyncCommitteeBits)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(s.SyncCommitteeSignature)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	return rlpBytes, nil
}

type HeaderUpdate struct {
	BeaconHeader       *BeaconBlockHeader
	ExecutionBlockHash []byte
}

func (update *HeaderUpdate) Encode() ([]byte, error) {
	h, err := update.BeaconHeader.Encode()
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(h)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(update.ExecutionBlockHash)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	return rlpBytes, nil
}

type FinalizedHeaderUpdate struct {
	HeaderUpdate   *HeaderUpdate
	FinalityBranch [][]byte
}

func (update *FinalizedHeaderUpdate) Encode() ([]byte, error) {
	headerUpdate, err := update.HeaderUpdate.Encode()
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(headerUpdate)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(update.FinalityBranch)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	return rlpBytes, nil
}

type SyncCommitteeUpdate struct {
	NextSyncCommittee       *eth.SyncCommittee
	NextSyncCommitteeBranch [][]byte
}

func (update *SyncCommitteeUpdate) Encode() ([]byte, error) {
	committee, err := rlp.EncodeToBytes(update.NextSyncCommittee)
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(committee)
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(update.NextSyncCommitteeBranch)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	return rlpBytes, nil
}

type LightClientUpdate struct {
	AttestedBeaconHeader    *BeaconBlockHeader
	SyncAggregate           *SyncAggregate
	SignatureSlot           uint64
	FinalizedUpdate         *FinalizedHeaderUpdate
	NextSyncCommitteeUpdate *SyncCommitteeUpdate
}

func (h *LightClientUpdate) Encode() ([]byte, error) {
	attestedHeader, err := h.AttestedBeaconHeader.Encode()
	if err != nil {
		return nil, err
	}
	b1, err := rlp.EncodeToBytes(attestedHeader)
	if err != nil {
		return nil, err
	}
	sig, err := h.SyncAggregate.Encode()
	if err != nil {
		return nil, err
	}
	b2, err := rlp.EncodeToBytes(sig)
	if err != nil {
		return nil, err
	}
	b3, err := rlp.EncodeToBytes(h.SignatureSlot)
	if err != nil {
		return nil, err
	}
	finalizedHeader, err := h.FinalizedUpdate.Encode()
	if err != nil {
		return nil, err
	}
	b4, err := rlp.EncodeToBytes(finalizedHeader)
	if err != nil {
		return nil, err
	}
	var b5 []byte
	if h.NextSyncCommitteeUpdate != nil {
		committee, err := h.NextSyncCommitteeUpdate.Encode()
		if err != nil {
			return nil, err
		}
		b5, err = rlp.EncodeToBytes(committee)
		if err != nil {
			return nil, err
		}
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
	rlpBytes = append(rlpBytes, b4...)
	rlpBytes = append(rlpBytes, b5...)
	return rlpBytes, nil
}

func (c *BeaconGrpcClient) BeaconHeaderconvert(data *BeaconBlockHeaderData) (*BeaconBlockHeader, error) {
	slot, err := strconv.ParseUint(data.Slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}
	index, err := strconv.ParseUint(data.ProposerIndex, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}
	h := new(BeaconBlockHeader)
	h.Slot = slot
	h.ProposerIndex = index
	h.BodyRoot = common.Hex2Bytes(data.BodyRoot[2:])
	h.ParentRoot = common.Hex2Bytes(data.ParentRoot[2:])
	h.StateRoot = common.Hex2Bytes(data.StateRoot[2:])
	return h, nil
}

func (c *BeaconGrpcClient) SyncAggregateconvert(data *SyncAggregateData) (*SyncAggregate, error) {
	aggregate := new(SyncAggregate)
	aggregate.SyncCommitteeBits = data.SyncCommitteeBits
	aggregate.SyncCommitteeSignature = common.Hex2Bytes(data.SyncCommitteeSignature[2:])
	return aggregate, nil
}

func (c *BeaconGrpcClient) CommitteeConvert(committee *SyncCommitteeData, branch []string) (*SyncCommitteeUpdate, error) {
	committeeUpdate := new(SyncCommitteeUpdate)

	nextCommittee := new(eth.SyncCommittee)
	nextCommittee.AggregatePubkey = common.Hex2Bytes(committee.AggregatePubkey[2:])
	for _, s := range committee.Pubkeys {
		nextCommittee.Pubkeys = append(nextCommittee.Pubkeys, common.Hex2Bytes(s[2:]))
	}
	committeeUpdate.NextSyncCommittee = nextCommittee

	for _, s := range branch {
		committeeUpdate.NextSyncCommitteeBranch = append(committeeUpdate.NextSyncCommitteeBranch, common.Hex2Bytes(s[2:]))
	}
	return committeeUpdate, nil
}

func (c *BeaconGrpcClient) FinalizedUpdateConvert(header *BeaconBlockHeaderData, branch []string) (*FinalizedHeaderUpdate, error) {
	update := new(FinalizedHeaderUpdate)

	for _, s := range branch {
		update.FinalityBranch = append(update.FinalityBranch, common.Hex2Bytes(s[2:]))
	}

	headerUpdate := new(HeaderUpdate)
	h, err := c.BeaconHeaderconvert(header)
	if err != nil {
		logger.Error("BeaconHeaderconvert error:", err)
		return nil, err
	}
	body, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(h.Slot, 10))
	if err != nil {
		logger.Error("GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	hash := body.GetExecutionPayload().BlockHash

	headerUpdate.BeaconHeader = h
	headerUpdate.ExecutionBlockHash = hash

	update.HeaderUpdate = headerUpdate
	return update, nil
}

func (c *BeaconGrpcClient) LightClientUpdateConvertNoCommitteeConvert(data *LightClientUpdateDataNoCommittee) (*LightClientUpdate, error) {
	attestedHeader, err := c.BeaconHeaderconvert(data.AttestedHeader)
	if err != nil {
		logger.Error("BeaconHeaderconvert error:", err)
		return nil, err
	}
	aggregate, err := c.SyncAggregateconvert(data.SyncAggregate)
	if err != nil {
		logger.Error("SyncAggregateconvert error:", err)
		return nil, err
	}
	finalizedUpdate, err := c.FinalizedUpdateConvert(data.FinalizedHeader, data.FinalityBranch)
	if err != nil {
		logger.Error("FinalizedUpdateConvert error:", err)
		return nil, err
	}
	slot, err := strconv.ParseUint(data.SignatureSlot, 0, 64)
	if err != nil {
		logger.Error("ParseUint error:", err)
		return nil, err
	}
	update := new(LightClientUpdate)
	update.AttestedBeaconHeader = attestedHeader
	update.SyncAggregate = aggregate
	update.NextSyncCommitteeUpdate = nil
	update.FinalizedUpdate = finalizedUpdate
	update.SignatureSlot = slot
	return update, nil
}

func (c *BeaconGrpcClient) LightClientUpdateConvert(data *LightClientUpdateData) (*LightClientUpdate, error) {
	attestedHeader, err := c.BeaconHeaderconvert(data.AttestedHeader)
	if err != nil {
		logger.Error("BeaconHeaderconvert error:", err)
		return nil, err
	}
	aggregate, err := c.SyncAggregateconvert(data.SyncAggregate)
	if err != nil {
		logger.Error("SyncAggregateconvert error:", err)
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
	slot, err := strconv.ParseUint(data.SignatureSlot, 0, 64)
	if err != nil {
		logger.Error("ParseUint error:", err)
		return nil, err
	}
	update := new(LightClientUpdate)
	update.AttestedBeaconHeader = attestedHeader
	update.SyncAggregate = aggregate
	update.NextSyncCommitteeUpdate = committeeUpdate
	update.FinalizedUpdate = finalizedUpdate
	update.SignatureSlot = slot
	return update, nil
}
