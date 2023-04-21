package toprelayer

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/wonderivan/logger"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
	"toprelayer/config"
	beaconrpc "toprelayer/rpc/ethbeacon_rpc"
)

const ip = "128.199.183.143"

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
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeSignature)
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
	lcu, err := b.GetLightClientUpdate(267)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("------------------")
	fmt.Printf("AttestedBeaconHeader:%+v\n", *lcu.AttestedBeaconHeader)
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeSignature)
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
		t.Error(err.Error())
		return
	}
	lcu, err := b.GetFinalityLightClientUpdate(6205891, false)
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
	fmt.Printf("SyncAggregate:%+v\n", lcu.SyncAggregate.SyncCommitteeSignature)
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

func TestPrysmInit(t *testing.T) {
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer := &Eth2TopRelayerV2{}
	err := relayer.Init(cfg, []string{hecoUrl}, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Error(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error(err)
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  50000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	tx, err := relayer.transactor.Init(ops, []byte("0x00000000"))
	if err != nil {
		t.Error(err)
	}
	t.Log(tx.Hash())
}

func TestPrysmDeposit(t *testing.T) {
	//FinalizedUpdate->HeaderUpdate:{BeaconHeader:0xc000306180 ExecutionBlockHash:[86 239 6 57 229 213 76 133 39 5 55 24 212 46 219 181 2 206 114 138 118 128 184 174 71 100 235 120 65 100 237 2]}
	//FinalizedUpdate->HeaderUpdate->BeaconHeader:{Slot:6209408 ProposerIndex:520505 ParentRoot:[76 239 217 180 58 159 189 129 95 58 48 136 21 158 220 121 31 87 160 240 251 246 131 188 176 107 247 204 61 192 30 111] StateRoot:[193 169 204 212 96 108 16 57 250 19 74 117 176 150 194 70 167 116 8 156 26 53 199 7 187 16 164 35 122 244 47 250] BodyRoot:[93 1 170 150 106 141 99 252 62 51 211 159 218 171 33 194 176 182 170 213 154 59 178 235 177 125 68 178 54 74 96 4]}
	//FinalizedUpdate:{HeaderUpdate:0xc000334100 FinalityBranch:[[252 245 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [74 232 20 64 233 229 244 123 49 138 69 133 203 221 172 178 68 191 187 200 213 195 117 64 167 173 226 29 217 23 47 92] [90 142 147 213 147 3 201 56 253 43 102 205 176 214 20 7 104 204 137 153 65 229 10 6 145 19 87 41 99 105 4 247] [16 8 20 29 137 111 167 198 77 245 248 231 85 44 15 158 146 90 255 238 169 136 235 239 155 51 111 240 245 137 154 23] [42 145 154 121 245 31 139 147 37 239 80 228 159 131 76 244 157 134 82 33 135 36 210 136 165 185 240 182 60 167 117 12] [89 68 212 174 147 244 167 100 223 229 134 101 255 73 206 162 90 162 140 241 80 8 7 226 158 226 78 96 196 228 38 172]]}
	f := &beaconrpc.FinalizedHeaderUpdate{
		HeaderUpdate: &beaconrpc.HeaderUpdate{
			BeaconHeader: &beaconrpc.BeaconBlockHeader{
				Slot:          6209408,
				ProposerIndex: 520505,
				ParentRoot:    []byte{76, 239, 217, 180, 58, 159, 189, 129, 95, 58, 48, 136, 21, 158, 220, 121, 31, 87, 160, 240, 251, 246, 131, 188, 176, 107, 247, 204, 61, 192, 30, 111},
				StateRoot:     []byte{193, 169, 204, 212, 96, 108, 16, 57, 250, 19, 74, 117, 176, 150, 194, 70, 167, 116, 8, 156, 26, 53, 199, 7, 187, 16, 164, 35, 122, 244, 47, 250},
				BodyRoot:      []byte{93, 1, 170, 150, 106, 141, 99, 252, 62, 51, 211, 159, 218, 171, 33, 194, 176, 182, 170, 213, 154, 59, 178, 235, 177, 125, 68, 178, 54, 74, 96, 4},
			},
			ExecutionBlockHash: []byte{86, 239, 6, 57, 229, 213, 76, 133, 39, 5, 55, 24, 212, 46, 219, 181, 2, 206, 114, 138, 118, 128, 184, 174, 71, 100, 235, 120, 65, 100, 237, 2},
		},
		FinalityBranch: [][]byte{
			{252, 245, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{74, 232, 20, 64, 233, 229, 244, 123, 49, 138, 69, 133, 203, 221, 172, 178, 68, 191, 187, 200, 213, 195, 117, 64, 167, 173, 226, 29, 217, 23, 47, 92},
			{90, 142, 147, 213, 147, 3, 201, 56, 253, 43, 102, 205, 176, 214, 20, 7, 104, 204, 137, 153, 65, 229, 10, 6, 145, 19, 87, 41, 99, 105, 4, 247},
			{16, 8, 20, 29, 137, 111, 167, 198, 77, 245, 248, 231, 85, 44, 15, 158, 146, 90, 255, 238, 169, 136, 235, 239, 155, 51, 111, 240, 245, 137, 154, 23},
			{42, 145, 154, 121, 245, 31, 139, 147, 37, 239, 80, 228, 159, 131, 76, 244, 157, 134, 82, 33, 135, 36, 210, 136, 165, 185, 240, 182, 60, 167, 117, 12},
			{89, 68, 212, 174, 147, 244, 167, 100, 223, 229, 134, 101, 255, 73, 206, 162, 90, 162, 140, 241, 80, 8, 7, 226, 158, 226, 78, 96, 196, 228, 38, 172},
		},
	}
	leaf, _ := f.Encode()
	root := make([]byte, 32)
	var branch []byte
	for i := 0; i < len(f.FinalityBranch); i++ {
		branch = append(branch, f.FinalityBranch[i]...)
	}
	ret := UnsafeMerkleProof(leaf, branch, 6, 41, root)
	fmt.Println("ret:", ret)
	fmt.Println("root:", root)
}

func TestD(t *testing.T) {
	s := (2193857 / 32) * 32
	fmt.Println(s)
}
