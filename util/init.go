package util

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/prysmaticlabs/prysm/v4/api/client/beacon"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wonderivan/logger"
	"math/big"
	"strconv"
	"toprelayer/relayer/toprelayer/ethtypes"
	"toprelayer/rpc/ethereum"
	beaconrpc "toprelayer/rpc/ethereum"
	lightclient "toprelayer/rpc/ethereum/light_client"
)

type ExtendedBeaconBlockHeader struct {
	Header             *lightclient.BeaconBlockHeader
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

func getEthInitData(eth1, prysm string) ([]byte, error) {
	beaconrpcclient, err := ethereum.NewBeaconChainClient(prysm)
	if err != nil {
		logger.Error("getEthInitData NewBeaconGrpcClient error:", err)
		return nil, err
	}
	ethrpcclient, err := ethclient.Dial(eth1)
	if err != nil {
		logger.Error("getEthInitData ethclient.Dial error:", err)
		return nil, err
	}
	lastUpdate, err := beaconrpcclient.GetLastFinalizedLightClientUpdateV2WithNextSyncCommittee()
	if err != nil {
		logger.Error("getEthInitData GetLightClientUpdate error:", err)
		return nil, err
	}
	lastSlot := lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	lastPeriod := ethereum.GetPeriodForSlot(lastSlot)
	prevUpdate, err := beaconrpcclient.GetNextSyncCommitteeUpdateV2(lastPeriod - 1)
	if err != nil {
		logger.Error(fmt.Sprintf("getEthInitData GetNextSyncCommitteeUpdate lastSlot:%d ，err：%s", lastSlot, err.Error()))
		return nil, err
	}

	var beaconHeader eth.BeaconBlockHeader
	beaconHeader.Slot = primitives.Slot(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot)
	beaconHeader.ProposerIndex = primitives.ValidatorIndex(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ProposerIndex)
	beaconHeader.BodyRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.BodyRoot)
	beaconHeader.ParentRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ParentRoot)
	beaconHeader.StateRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.StateRoot)
	root, err := beaconHeader.HashTreeRoot()
	if err != nil {
		logger.Error("getEthInitData HashTreeRoot error:", err)
		return nil, err
	}
	finalizedHeader := new(ExtendedBeaconBlockHeader)
	finalizedHeader.BeaconBlockRoot = root[:]
	finalizedHeader.Header = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader // 2381600
	finalizedHeader.ExecutionBlockHash = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.ExecutionBlockHash)

	finalitySlot := lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	finalizeBody, err := beaconrpcclient.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(finalitySlot), 10)))
	if err != nil {
		logger.Error("getEthInitData GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	executionPayload, err := finalizeBody.Execution()
	if err != nil {
		logger.Error("getEthInitData GetBeaconBlockBodyForBlockId error:", err)
		return nil, err
	}
	number := executionPayload.BlockNumber()

	header, err := ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		logger.Error("getEthInitData HeaderByNumber error:", err)
		return nil, err
	}

	initParam := new(InitInput)
	initParam.FinalizedExecutionHeader = header
	initParam.FinalizedBeaconHeader = finalizedHeader
	initParam.NextSyncCommittee = lastUpdate.NextSyncCommitteeUpdate.NextSyncCommittee
	initParam.CurrentSyncCommittee = prevUpdate.NextSyncCommittee

	printEthHeaderInfo(initParam)

	h := common.Bytes2Hex(initParam.FinalizedExecutionHeader.Hash().Bytes())
	fmt.Println("eth header hash:", h)
	ExecutionBlockHash := common.Bytes2Hex(initParam.FinalizedBeaconHeader.ExecutionBlockHash)
	fmt.Println("ExecutionBlockHash hash:", ExecutionBlockHash)
	bytes, err := initParam.Encode()
	if err != nil {
		logger.Error("getEthInitData initParam.Encode error:", err)
		return nil, err
	}
	return bytes, nil
}

func getEthInitDataWithHeight(eth1, prysm, slot string) ([]byte, error) {
	beaconrpcclient, err := beaconrpc.NewBeaconChainClient(prysm)
	if err != nil {
		logger.Error("getEthInitData NewBeaconGrpcClient error:", err)
		return nil, err
	}
	ethrpcclient, err := ethclient.Dial(eth1)
	if err != nil {
		logger.Error("getEthInitData ethclient.Dial error:", err)
		return nil, err
	}
	lastSlot, err := strconv.ParseUint(slot, 0, 64)
	if err != nil {
		logger.Error("ParseInt error:", err)
		return nil, err
	}
	lastPeriod := beaconrpc.GetPeriodForSlot(primitives.Slot(lastSlot))
	// 269 2203865
	lastUpdate, err := beaconrpcclient.GetLightClientUpdateV2(lastPeriod)
	if err != nil {
		logger.Error("getEthInitData GetLightClientUpdate error:", err)
		return nil, err
	}
	prevUpdate, err := beaconrpcclient.GetNextSyncCommitteeUpdateV2(lastPeriod - 1)
	if err != nil {
		logger.Error("getEthInitData GetNextSyncCommitteeUpdate error:", err)
		return nil, err
	}

	var beaconHeader eth.BeaconBlockHeader
	beaconHeader.Slot = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	beaconHeader.ProposerIndex = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ProposerIndex
	beaconHeader.BodyRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.BodyRoot)
	beaconHeader.ParentRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.ParentRoot)
	beaconHeader.StateRoot = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.StateRoot)
	root, err := beaconHeader.HashTreeRoot()
	if err != nil {
		logger.Error("getEthInitData HashTreeRoot error:", err)
		return nil, err
	}
	finalizedHeader := new(ExtendedBeaconBlockHeader)
	finalizedHeader.BeaconBlockRoot = root[:]
	finalizedHeader.Header = lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader
	finalizedHeader.ExecutionBlockHash = ethtypes.ConvertBytes32ToBytesSlice(lastUpdate.FinalityUpdate.HeaderUpdate.ExecutionBlockHash)

	finalitySlot := lastUpdate.FinalityUpdate.HeaderUpdate.BeaconHeader.Slot
	finalizeBody, err := beaconrpcclient.GetBeaconBlockBody(beacon.StateOrBlockId(strconv.FormatUint(uint64(finalitySlot), 10)))
	if err != nil {
		logger.Error("getEthInitData GetBeaconBlockBody error:", err)
		return nil, err
	}
	executionData, err := finalizeBody.Execution()
	if err != nil {
		logger.Error("getEthInitData GetBeaconBlockBody error:", err)
		return nil, err
	}
	number := executionData.BlockNumber()

	header, err := ethrpcclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		logger.Error("getEthInitData HeaderByNumber error:", err)
		return nil, err
	}

	initParam := new(InitInput)
	initParam.FinalizedExecutionHeader = header
	initParam.FinalizedBeaconHeader = finalizedHeader
	initParam.NextSyncCommittee = lastUpdate.NextSyncCommitteeUpdate.NextSyncCommittee
	initParam.CurrentSyncCommittee = prevUpdate.NextSyncCommittee

	printEthHeaderInfo(initParam)

	bytes, err := initParam.Encode()
	if err != nil {
		logger.Error("getEthInitData initParam.Encode error:", err)
		return nil, err
	}
	return bytes, nil
}

func printEthHeaderInfo(initParam *InitInput) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ETH Data <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Printf("FinalizedExecutionHeader: %+v \n", *initParam.FinalizedExecutionHeader)
	fmt.Printf("FinalizedExecutionHeader.Bloom: %+v \n", common.Bytes2Hex(initParam.FinalizedExecutionHeader.Bloom[:]))
	fmt.Printf("FinalizedExecutionHeader.Extra: %+v \n", common.Bytes2Hex(initParam.FinalizedExecutionHeader.Extra))
	fmt.Printf("FinalizedBeaconHeader.Header.Slot: %+v \n", initParam.FinalizedBeaconHeader.Header.Slot)
	fmt.Printf("FinalizedBeaconHeader.Header.ProposerIndex: %+v \n", initParam.FinalizedBeaconHeader.Header.ProposerIndex)
	fmt.Printf("FinalizedBeaconHeader.Header.ParentRoot: %+v \n", common.Bytes2Hex(initParam.FinalizedBeaconHeader.Header.ParentRoot[:]))
	fmt.Printf("FinalizedBeaconHeader.Header.StateRoot: %+v \n", common.Bytes2Hex(initParam.FinalizedBeaconHeader.Header.StateRoot[:]))
	fmt.Printf("FinalizedBeaconHeader.Header.BodyRoot: %+v \n", common.Bytes2Hex(initParam.FinalizedBeaconHeader.Header.BodyRoot[:]))
	fmt.Printf("FinalizedBeaconHeader.BeaconBlockRoot: %+v \n", common.Bytes2Hex(initParam.FinalizedBeaconHeader.BeaconBlockRoot))
	fmt.Printf("FinalizedBeaconHeader.ExecutionBlockHash: %+v \n", common.Bytes2Hex(initParam.FinalizedBeaconHeader.ExecutionBlockHash))
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ETH Data <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
