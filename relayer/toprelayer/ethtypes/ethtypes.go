package ethtypes

import (
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
	BeaconHeader       *eth.BeaconBlockHeader
	ExecutionBlockHash common.Hash
}

type FinalizedHeaderUpdate struct {
	HeaderUpdate   *HeaderUpdate
	FinalityBranch [][]byte
}

type SyncCommitteeUpdate struct {
	NextSyncCommittee       *eth.SyncCommittee
	NextSyncCommitteeBranch [][]byte
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
