package toprelayer

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/wonderivan/logger"
	"strconv"
	"strings"
	"testing"
	"time"
	beaconrpc "toprelayer/rpc/ethbeacon_rpc"
)

//  "http://69.194.0.7:8545","69.194.0.7:4000", "http://69.194.0.7:9596"
const ip = "69.194.0.7"

var eth1 = "http://" + ip + ":8545"
var prysm = ip + ":4000"
var lodestar = "http://" + ip + ":9596"

func TestGetFinalizedLightClientUpdate4Lodestar(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	//a, err := b.GetLightClientUpdate(755)
	//if err != nil {
	//	log.Error(err.Error())
	//	return
	//}
	//fmt.Printf("AttestedBeaconHeader:%+v\n", *a.AttestedBeaconHeader)
	//fmt.Printf("SyncAggregate:%+v\n", *a.SyncAggregate)
	//fmt.Printf("SignatureSlot:%+v\n", a.SignatureSlot)
	//fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *a.FinalizedUpdate.HeaderUpdate)
	//fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *a.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	//fmt.Printf("FinalizedUpdate:%+v\n", *a.FinalizedUpdate)
	//fmt.Printf("NextSyncCommitteeUpdate:%+v\n", *a.NextSyncCommitteeUpdate)
	lcu, err := b.GetFinalizedLightClientUpdate()
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("------------------")
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeBits)
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	//fmt.Printf("FinalizedUpdate->FinalityBranch:%v\n", *lcu.FinalizedUpdate.FinalityBranch)
	fmt.Printf("NextSyncCommitteeUpdate:%+v\n", lcu.NextSyncCommitteeUpdate)
}

func TestGetLightClientUpdate4Lodestar(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	lcu, err := b.GetLightClientUpdate(754)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("------------------")
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", *lcu.SyncAggregate)
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	//fmt.Printf("FinalizedUpdate->FinalityBranch:%v\n", *lcu.FinalizedUpdate.FinalityBranch)
	fmt.Printf("NextSyncCommitteeUpdate:%+v\n", *lcu.NextSyncCommitteeUpdate.NextSyncCommittee)
}

func TestGetGetNextSyncCommitteeUpdate4Lodestar(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	lastSlot, err := b.GetLastFinalizedSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return
	}
	lastPeriod := beaconrpc.GetPeriodForSlot(lastSlot)
	prevUpdate, err := b.GetNextSyncCommitteeUpdate(lastPeriod - 1)
	if err != nil {
		logger.Error("GetNextSyncCommitteeUpdate error:", err)
		return
	}
	fmt.Println("------------------")
	fmt.Printf("NextSyncCommittee:%+v\n", *prevUpdate.NextSyncCommittee)
	fmt.Printf("NextSyncCommittee->Pubkeys:%+v\n", prevUpdate.NextSyncCommittee.Pubkeys)
}

func TestPrysm(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	lcu, err := b.GetFinalityLightClientUpdate(6205891, false)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeBits)
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	//fmt.Printf("FinalizedUpdate->FinalityBranch:%v\n", *lcu.FinalizedUpdate.FinalityBranch)
	fmt.Printf("NextSyncCommitteeUpdate:%+v\n", lcu.NextSyncCommitteeUpdate)
}

func TestGetFinalizedLightClientUpdate4Prysm(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		t.Error(err.Error())
		return
	}
	lcu, err := b.GetFinalizedLightClientUpdateV2()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeBits)
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	if lcu.NextSyncCommitteeUpdate != nil {
		fmt.Printf("NextSyncCommitteeUpdate:%+v\n", *lcu.NextSyncCommitteeUpdate.NextSyncCommittee)
	}
}

func TestGetLightClientUpdateV24Prysm(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		t.Error(err.Error())
		return
	}
	lcu, err := b.GetLightClientUpdateV2(755)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeBits)
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	if lcu.NextSyncCommitteeUpdate != nil {
		fmt.Printf("NextSyncCommitteeUpdate:%+v\n", *lcu.NextSyncCommitteeUpdate.NextSyncCommittee)
	}
}

func TestGetGetNextSyncCommitteeUpdate4Prysm(t *testing.T) {
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	lastSlot, err := b.GetLastFinalizedSlotNumber()
	if err != nil {
		logger.Error("GetLastFinalizedSlotNumber error:", err)
		return
	}
	lastPeriod := beaconrpc.GetPeriodForSlot(lastSlot)
	prevUpdate, err := b.GetNextSyncCommitteeUpdateV2(lastPeriod - 1)
	if err != nil {
		logger.Error("GetNextSyncCommitteeUpdate error:", err)
		return
	}
	fmt.Println("------------------")
	fmt.Printf("NextSyncCommittee:%+v\n", *prevUpdate.NextSyncCommittee)
	fmt.Printf("NextSyncCommittee->Pubkeys:%+v\n", prevUpdate.NextSyncCommittee.Pubkeys)
}

func TestGetFinalizedLightClientUpdate4Lodestar4Prysm(t *testing.T) {
	//client, err := beaconrpc.NewBeaconGrpcClient(prysm, "")
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	b, err := beaconrpc.NewBeaconGrpcClient(prysm, lodestar)
	if err != nil {
		log.Error(err.Error())
		return
	}
	relayer := NewRelayerByRpcClient(prysm)

	// 6209472 当前的slot
	header, err := relayer.beaconrpcclient.GetNonEmptyBeaconBlockHeader(6188252)
	if err != nil {
		logger.Error("Eth2TopRelayerV2 GetNonEmptyBeaconBlockHeader error:", err)
		return
	}
	fmt.Printf("header:%d \n", uint64(header.Slot))
	lcu, err := b.GetFinalityLightClientUpdate(uint64(header.Slot), false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeBits.Bytes())
	fmt.Printf("SignatureSlot:%+v\n", lcu.SignatureSlot)
	fmt.Printf("FinalizedUpdate->HeaderUpdate:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate)
	fmt.Printf("FinalizedUpdate->HeaderUpdate->BeaconHeader:%+v\n", *lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader)
	fmt.Printf("FinalizedUpdate:%+v\n", *lcu.FinalizedUpdate)
	if lcu.NextSyncCommitteeUpdate != nil {
		fmt.Printf("NextSyncCommitteeUpdate:%+v\n", *lcu.NextSyncCommitteeUpdate.NextSyncCommittee)
	}
	//if err = relayer.isCorrectFinalityUpdate(lcu, lcu.NextSyncCommitteeUpdate.NextSyncCommittee); err != nil {
	//	t.Error(err.Error())
	//	return
	//}
}
func TestABI(t *testing.T) {
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("store(uint256)"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("store_require(uint256)"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("store_assert(uint256)"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("store_revert(uint256)"))))

	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("retrieve()"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("retrieve_revert()"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("retrieve_revert(uint256)"))))
	fmt.Println(common.Bytes2Hex(crypto.Keccak256([]byte("retrieve_assert()")))[:8])
}
func TestGetBlockState(t *testing.T) {
	client, err := beaconrpc.NewBeaconGrpcClient(prysm, "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	start := time.Now()
	beaconState, err := client.GetBeaconState(strconv.FormatUint(6187835, 10))
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("time:", time.Since(start))
	fmt.Printf("beaconState:%+v\n", *beaconState)
}

func TestGetAttestedSlot(t *testing.T) {
	client, err := beaconrpc.NewBeaconGrpcClient(prysm, "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 测试的 slot  6209184
	lastFinalizedSlot := uint64(6209184)
	fmt.Println("lastFinalizedSlot:", lastFinalizedSlot)
	nextFinalizedSlot := lastFinalizedSlot + beaconrpc.ONE_EPOCH_IN_SLOTS
	attestedSlot := uint64(nextFinalizedSlot + 2*beaconrpc.ONE_EPOCH_IN_SLOTS)
	// 6209280
	fmt.Println("attestedSlot:", attestedSlot)
	beaconBlockHeader, err := client.GetNonEmptyBeaconBlockHeader(attestedSlot)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("beaconBlockHeader.Slot: %+v \n", beaconBlockHeader.Slot)
}

func TestBeaconGrpcApi(t *testing.T) {
	c, err := beaconrpc.NewBeaconGrpcClient(prysm, "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	finalizedSlot, err := c.GetLastFinalizedSlotNumber()
	if err != nil {
		t.Error("GetLastFinalizedSlotNumber error:", err)
		return
	}
	fmt.Println("finalizedSlot", finalizedSlot)
	for slot := uint64(6209072); slot < finalizedSlot; slot++ {
		h, err := c.GetBeaconBlockHeaderForBlockId(strconv.FormatUint(slot, 10))
		if err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				fmt.Println(err.Error())
				continue
			} else {
				t.Error("GetBeaconBlockBodyForBlockId error:", err)
				return
			}
		}
		fmt.Printf("slot: %d, h.Slot: %+v \n", slot, h.Slot)
	}

}

func TestDemo(t *testing.T) {
	b, _ := hex.DecodeString("b467f3f58c2ef26d0342948a3afd0299ea70b615")
	//num := 10095634316
	//hexStr := fmt.Sprintf("%x", num)
	fmt.Println(string(b))
}
