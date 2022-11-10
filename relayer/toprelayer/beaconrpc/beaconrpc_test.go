package beaconrpc

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	primitives "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
)

const SEPOLIA_URL = "https://lodestar-sepolia.chainsafe.io"
const SEPOLIA_ETH1_URL = "https://rpc.sepolia.org"

func TestGetBeaconHeaderAndBlockForBlockId(t *testing.T) {
	var s uint64 = 969983
	ss := strconv.Itoa(969983)

	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	h, err := c.GetBeaconBlockHeaderForBlockId(ss)
	if err != nil {
		t.Fatal(err)
	}
	if uint64(h.Slot) != s {
		t.Fatal("slot not equal")
	}
	t.Log(h.ProposerIndex)
	t.Log(common.Bytes2Hex(h.BodyRoot))
	t.Log(common.Bytes2Hex(h.ParentRoot))
	t.Log(common.Bytes2Hex(h.StateRoot))

	time.Sleep(time.Duration(5) * time.Second)

	b, err := c.GetBeaconBlockBodyForBlockId(ss)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b.HashTreeRoot())
	t.Log(b.GetExecutionPayload().BlockHash)

	time.Sleep(time.Duration(5) * time.Second)

	n, err := c.GetBlockNumberForSlot(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
}

func TestGetBeaconBlockForBlockIdErrNoBlockForSlot(t *testing.T) {
	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBeaconBlockBodyForBlockId(strconv.Itoa(999572))
	if err != nil {
		if IsErrorNoBlockForSlot(err) {
			t.Log("catch error success", err)
		} else {
			t.Fatal("err not catch:", err)
		}
	} else {
		t.Fatal("err not catch:", err)
	}
}

func TestGetLightClientUpdate(t *testing.T) {
	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetLightClientUpdate(122)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFinalizedLightClientUpdateWithSyncCommityUpdate(t *testing.T) {
	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	update, err := c.GetFinalizedLightClientUpdateWithSyncCommityUpdate()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(update)
}

type ExtendedBeaconBlockHeader struct {
	Header             *BeaconBlockHeader
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
	b1, err := rlp.EncodeToBytes(init.FinalizedExecutionHeader)
	if err != nil {
		return nil, err
	}
	b2, err := init.FinalizedBeaconHeader.Encode()
	if err != nil {
		return nil, err
	}
	b3, err := rlp.EncodeToBytes(init.CurrentSyncCommittee)
	if err != nil {
		return nil, err
	}
	b4, err := rlp.EncodeToBytes(init.NextSyncCommittee)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
	rlpBytes = append(rlpBytes, b4...)
	list, err := rlp.EncodeToBytes(rlpBytes)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func TestEth1(t *testing.T) {
	eth1, err := ethclient.Dial(SEPOLIA_ETH1_URL)
	if err != nil {
		t.Fatal(err)
	}

	header, err := eth1.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(2252532))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(header)
}

func TestInitInputData(t *testing.T) {
	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}

	lastSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		t.Fatal(err)
	}
	lastPeriod := GetPeriodForSlot(lastSlot)
	lastUpdate, err := c.GetLightClientUpdate(lastPeriod)
	if err != nil {
		t.Fatal(err)
	}
	prevUpdate, err := c.GetLightClientUpdate(lastPeriod - 1)
	if err != nil {
		t.Fatal(err)
	}

	var beaconHeader eth.BeaconBlockHeader
	beaconHeader.Slot = primitives.Slot(lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot)
	beaconHeader.ProposerIndex = primitives.ValidatorIndex(lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.ProposerIndex)
	beaconHeader.BodyRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.BodyRoot
	beaconHeader.ParentRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.ParentRoot
	beaconHeader.StateRoot = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.StateRoot
	root, err := beaconHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	finalizedHeader := new(ExtendedBeaconBlockHeader)
	finalizedHeader.BeaconBlockRoot = root[:]
	finalizedHeader.Header = lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader
	finalizedHeader.ExecutionBlockHash = lastUpdate.FinalizedUpdate.HeaderUpdate.ExecutionBlockHash

	finalitySlot := lastUpdate.FinalizedUpdate.HeaderUpdate.BeaconHeader.Slot
	finalizeBody, err := c.GetBeaconBlockBodyForBlockId(strconv.FormatUint(finalitySlot, 10))
	if err != nil {
		t.Fatal(err)
	}
	number := finalizeBody.GetExecutionPayload().BlockNumber

	eth1, err := ethclient.Dial(SEPOLIA_ETH1_URL)
	if err != nil {
		t.Fatal(err)
	}

	header, err := eth1.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
	if err != nil {
		t.Fatal(err)
	}

	initParam := new(InitInput)
	initParam.FinalizedExecutionHeader = header
	initParam.FinalizedBeaconHeader = finalizedHeader
	initParam.NextSyncCommittee = lastUpdate.NextSyncCommitteeUpdate.NextSyncCommittee
	initParam.CurrentSyncCommittee = prevUpdate.NextSyncCommitteeUpdate.NextSyncCommittee

	bytes, err := initParam.Encode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bytes))
}

func TestLightClientUpdateData(t *testing.T) {
	c, err := NewBeaconGrpcClient(SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	lastUpdate, err := c.GetFinalizedLightClientUpdate()
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := lastUpdate.Encode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bytes))
}
