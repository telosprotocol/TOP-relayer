package ethtypes

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

const (
	// Bellatrix Fork Epoch for mainnet config.
	mainnetBellatrixForkEpoch = 144896 // Sept 6, 2022, 11:34:47am UTC
)

var (
	BellatrixForkVersion []byte = []byte{2, 0, 0, 0}
	BellatrixForkEpoch   uint64 = mainnetBellatrixForkEpoch
)

var (
	DomainSyncCommittee [4]byte = bytesutil.Uint32ToBytes4(0x07000000)
)

var (
	GenesisValidatorsRoot [32]byte = [32]byte{0x4b, 0x36, 0x3d, 0xb9, 0x4e, 0x28, 0x61, 0x20, 0xd7, 0x6e, 0xb9, 0x05, 0x34, 0x0f, 0xdd, 0x4e, 0x54, 0xbf, 0xe9, 0xf0, 0x6b, 0xf3, 0x3f, 0xf6, 0xcf, 0x5a, 0xd2, 0x7f, 0x51, 0x1b, 0xfe, 0x95}
)

type ExtendedBeaconBlockHeader struct {
	Header             *eth.BeaconBlockHeader
	BeaconBlockRoot    common.Hash
	ExecutionBlockHash common.Hash
}

type HeaderUpdate struct {
	BeaconHeader        *eth.BeaconBlockHeader
	ExecutionBlockHash  common.Hash
	ExecutionHashBranch []common.Hash
}

type FinalizedHeaderUpdate struct {
	HeaderUpdate   *HeaderUpdate
	FinalityBranch [][common.HashLength]byte
}

type SyncCommitteeUpdate struct {
	NextSyncCommittee       *eth.SyncCommittee
	NextSyncCommitteeBranch [][common.HashLength]byte
}

type LightClientUpdate struct {
	AttestedBeaconHeader    *eth.BeaconBlockHeader
	SyncAggregate           *eth.SyncAggregate
	SignatureSlot           uint64
	FinalizedUpdate         *FinalizedHeaderUpdate
	NextSyncCommitteeUpdate *SyncCommitteeUpdate
}

type LightClientState struct {
	FinalizedBeaconHeader *ExtendedBeaconBlockHeader
	CurrentSyncCommittee  *eth.SyncCommittee
	NextSyncCommittee     *eth.SyncCommittee
}

type ExecutionBlockProof struct {
	BlockHash common.Hash
	Proof     []common.Hash
}

func ConvertSliceHash2Bytes(hash []common.Hash) [][]byte {
	res := make([][]byte, len(hash))
	for i := range hash {
		res[i] = hash[i][:]
	}
	return res
}

func ConvertSliceHash2SliceBytes32(hash []common.Hash) [][common.HashLength]byte {
	res := make([][common.HashLength]byte, len(hash))
	for i := range hash {
		res[i] = hash[i]
	}
	return res
}

func ConvertSliceBytes2SliceBytes32(bytes [][]byte) ([][common.HashLength]byte, error) {
	res := make([][common.HashLength]byte, len(bytes))
	for i, b := range bytes {
		if len(b) != common.HashLength {
			err := fmt.Errorf("invalid hash length:%d, index:%d", len(b), i)
			return nil, err
		}
		res[i] = [common.HashLength]byte(b)
	}
	return res, nil
}

func ConvertSliceBytes2Hash(bytes [][32]byte) []common.Hash {
	res := make([]common.Hash, len(bytes))
	for i, b := range bytes {
		res[i] = b
	}
	return res
}
