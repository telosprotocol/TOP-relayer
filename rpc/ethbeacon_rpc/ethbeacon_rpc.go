package ethbeacon_rpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	pb "github.com/prysmaticlabs/prysm/v4/proto/eth/service"
	v1 "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	v2 "github.com/prysmaticlabs/prysm/v4/proto/eth/v2"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewBeaconGrpcClient(grpcUrl string) (*BeaconGrpcClient, error) {
	grpcDefault, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("create grpcDefault error:", err)
		return nil, err
	}

	grpcBigData, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(200000000)))
	if err != nil {
		logger.Error("create grpcBigData error:", err)
		return nil, err
	}
	c := &BeaconGrpcClient{
		client:      pb.NewBeaconChainClient(grpcDefault),
		debugclient: pb.NewBeaconDebugClient(grpcBigData),
		httpclient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
		httpurl: "",
	}
	return c, nil
}

func (c *BeaconGrpcClient) GetBeaconBlockBodyForBlockId(id string) (*v2.BeaconBlockBodyCapella, error) {
	resp, err := c.client.GetBlockV2(context.Background(), &v2.BlockRequestV2{BlockId: []byte(id)})
	if err != nil {
		return nil, err
	}
	signedBlock, ok := resp.Data.Message.(*v2.SignedBeaconBlockContainer_CapellaBlock)
	if !ok {
		return nil, errors.New("resp.data.message error")
	}
	return signedBlock.CapellaBlock.GetBody(), nil
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
		return 0, err
	}
	return b.GetExecutionPayload().BlockNumber, nil
}

func (c *BeaconGrpcClient) GetBlockHashForSlot(slot uint64) (common.Hash, error) {
	b, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(slot, 10))
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(b.GetExecutionPayload().BlockHash), nil
}

func GetPeriodForSlot(slot uint64) uint64 {
	return (slot / (SLOTS_PER_EPOCH * EPOCHS_PER_PERIOD))
}

func (c *BeaconGrpcClient) getBeaconState(id string) (*eth.BeaconStateCapella, error) {
	start := time.Now()
	defer func() {
		logger.Info("Slot:%s,getBeaconState time:%v", id, time.Since(start))
	}()
	var data []byte
	//attestedSlot, err := strconv.Atoi(id)
	//period, epoch, slot := SplitSlot(uint64(attestedSlot))
	//fileName := fmt.Sprintf("state_%d_%d_%d_%s", period, epoch, slot, id)
	//data, err = os.ReadFile(fileName)
	//if err != nil {
	resp, err := c.debugclient.GetBeaconStateSSZV2(context.Background(), &v2.BeaconStateRequestV2{StateId: []byte(id)})
	if err != nil {
		logger.Error("GetBeaconStateV2 error:", err)
		return nil, err
	}
	data = resp.Data
	//	if err = os.WriteFile(fileName, resp.Data, 0666); err != nil {
	//		return nil, err
	//	}
	//}

	var state eth.BeaconStateCapella
	if err = state.UnmarshalSSZ(data); err != nil {
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
	lastSlot, err := c.GetLastSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return nil, err
	}

	for slot := startSlot; slot <= lastSlot; slot++ {
		if h, err := c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(slot, 10)); err != nil {
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

func (c *BeaconGrpcClient) GetNonEmptyBeaconBlockHeaderLimitRange(startSlot, finalizedSlot uint64) (*eth.BeaconBlockHeader, error) {
	for slot := startSlot; slot < finalizedSlot; slot++ {
		if h, err := c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(slot, 10)); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				logger.Error("GetBeaconBlockHeaderForBlockId slot(%d) error:", slot, err)
				continue
			} else {
				logger.Error("GetBeaconBlockBodyForBlockId error:", err)
				return nil, err
			}
		} else {
			return h, nil
		}
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

	var result []LightClientUpdateMsg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
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
	return c.LightClientUpdateConvert(&result[0].Data)
}

func (c *BeaconGrpcClient) GetNextSyncCommitteeUpdate(period uint64) (*SyncCommitteeUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%d&count=1", c.httpurl, period)
	resp, err := c.httpclient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result []LightClientUpdateMsg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	if len(body) == 0 {
		logger.Error("body empty")
		return nil, errors.New("http body empty")
	}
	if err = json.Unmarshal(body, &result); err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		logger.Error(err.Error())
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

func (c *BeaconGrpcClient) GetFinalizedLightClientUpdate() (*LightClientUpdate, error) {
	str := fmt.Sprintf("%s/eth/v1/beacon/light_client/finality_update", c.httpurl)
	resp, err := c.httpclient.Get(str)
	if err != nil {
		logger.Error("http Get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result LightClientUpdateNoCommitteeMsg
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("outil.ReadAll error:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		err = fmt.Errorf("unmarshal error:%s body: %s", err.Error(), string(body))
		return nil, err
	}
	return c.LightClientUpdateConvertNoCommitteeConvert(&result.Data)
}

func (c *BeaconGrpcClient) BeaconHeaderconvert(data *BeaconBlockHeaderData) (*BeaconBlockHeader, error) {
	slot, err := strconv.ParseUint(data.Beacon.Slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}
	index, err := strconv.ParseUint(data.Beacon.ProposerIndex, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}
	h := new(BeaconBlockHeader)
	h.Slot = slot
	h.ProposerIndex = index
	h.BodyRoot = common.Hex2Bytes(data.Beacon.BodyRoot[2:])
	h.ParentRoot = common.Hex2Bytes(data.Beacon.ParentRoot[2:])
	h.StateRoot = common.Hex2Bytes(data.Beacon.StateRoot[2:])
	return h, nil
}

func (c *BeaconGrpcClient) SyncAggregateconvert(data *SyncAggregateData) (*SyncAggregate, error) {
	aggregate := new(SyncAggregate)
	//aggregate.SyncCommitteeBits = data.SyncCommitteeBits
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
