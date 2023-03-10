package ethbeacon_rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"math/big"
	"strconv"
	"testing"
	"time"
)

const LOCAL_GRPC_URL = "localhost:4000"
const SEPOLIA_URL = "https://lodestar-sepolia.chainsafe.io"
const SEPOLIA_ETH1_URL = "https://rpc.sepolia.org"

func TestGetBeaconHeaderAndBlockForBlockId(t *testing.T) {
	var s uint64 = 969983
	ss := strconv.Itoa(969983)

	c, err := NewBeaconGrpcClient(LOCAL_GRPC_URL, SEPOLIA_URL)
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
	c, err := NewBeaconGrpcClient(LOCAL_GRPC_URL, SEPOLIA_URL)
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

func TestEth1(t *testing.T) {
	// config
	height := big.NewInt(0).SetUint64(2256927)
	hash := common.HexToHash("d26a8a468987d1ea34406ba622a4ae44eb67922d4166784cc84496a8b04be874")

	// test
	eth1, err := ethclient.Dial(SEPOLIA_ETH1_URL)
	if err != nil {
		t.Fatal(err)
	}

	h1, err := eth1.HeaderByNumber(context.Background(), height)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("HeaderByNumber:", h1)

	h2, err := eth1.HeaderByHash(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("HeaderByHash:", h2)
}

func TestGetEth1Data(t *testing.T) {
	// config
	start := uint64(2260231)
	end := uint64(2264210)
	// test
	eth1, err := ethclient.Dial(SEPOLIA_ETH1_URL)
	if err != nil {
		t.Fatal(err)
	}
	for h := start; h <= end; h += 1 {
		header, err := eth1.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			t.Fatal(err)
		}
		bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(h, common.Bytes2Hex(bytes))
		time.Sleep(time.Millisecond * 500)
	}
}

func TestGetLightClientUpdateData(t *testing.T) {
	c, err := NewBeaconGrpcClient(LOCAL_GRPC_URL, SEPOLIA_URL)
	if err != nil {
		t.Fatal(err)
	}
	lastSlot, err := c.GetLastSlotNumber()
	if err != nil {
		t.Fatal(err)
	}
	lastPeriod := GetPeriodForSlot(lastSlot)
	update, err := c.GetLightClientUpdate(lastPeriod)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := update.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(common.Bytes2Hex(bytes))
}

func TestGetFinalizedLightClientUpdateData(t *testing.T) {
	c, err := NewBeaconGrpcClient(LOCAL_GRPC_URL, SEPOLIA_URL)
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
	t.Log(common.Bytes2Hex(bytes))
}

func TestRLP4BeaconBlockHeader(t *testing.T) {
	data := &LightClientUpdate{
		AttestedBeaconHeader: &BeaconBlockHeader{
			Slot:          1,
			ProposerIndex: 2,
			ParentRoot:    []byte("ParentRoot"),
			StateRoot:     []byte("StateRoot"),
			BodyRoot:      []byte("BodyRoot"),
		},
		SyncAggregate: &SyncAggregate{
			SyncCommitteeBits:      "1111111111000000011111",
			SyncCommitteeSignature: []byte("SyncCommitteeSignature"),
		},
		SignatureSlot: 0,
		FinalizedUpdate: &FinalizedHeaderUpdate{
			HeaderUpdate: &HeaderUpdate{
				BeaconHeader: &BeaconBlockHeader{
					Slot:          1,
					ProposerIndex: 2,
					ParentRoot:    []byte("ParentRoot"),
					StateRoot:     []byte("StateRoot"),
					BodyRoot:      []byte("BodyRoot"),
				},
				ExecutionBlockHash: []byte("ExecutionBlockHash"),
			},
			FinalityBranch: [][]byte{
				[]byte("aaa"),
				[]byte("bbb"),
				[]byte("ccc"),
			},
		},
		NextSyncCommitteeUpdate: &SyncCommitteeUpdate{
			NextSyncCommittee: &eth.SyncCommittee{
				Pubkeys:         nil,
				AggregatePubkey: nil,
			},
			NextSyncCommitteeBranch: nil,
		},
	}
	
	ret, err := rlp.EncodeToBytes(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("ret:%+v", ret)

	//abiData, err := eth2bridge.PackSubmitBeaconChainLightClientUpdateParam(ret)
	//if err != nil {
	//	t.Fatal(err.Error())
	//}
	//t.Logf("abiData:%+x", abiData)
	//if len(abiData)%32 != 0 {
	//	t.Fatal(len(abiData))
	//}
	//abi, err := eth2bridge.Eth2ClientMetaData.GetAbi()
	//if err != nil {
	//	t.Fatal(err.Error())
	//}
	//
	//var lightClientUpdate LightClientUpdate
	//err = abi.UnpackIntoInterface(&lightClientUpdate, "submit_beacon_chain_light_client_update", abiData)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("lightClientUpdate:%v", lightClientUpdate)

}
