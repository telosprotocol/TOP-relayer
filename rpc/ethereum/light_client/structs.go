package light_client

import (
	"github.com/ethereum/go-ethereum/rlp"
	field_params "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

type BeaconBlockHeader struct {
	Slot          primitives.Slot
	ProposerIndex primitives.ValidatorIndex
	ParentRoot    [field_params.RootLength]byte
	StateRoot     [field_params.RootLength]byte
	BodyRoot      [field_params.RootLength]byte
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
	SyncCommitteeBits      []byte
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
	BeaconHeader        *BeaconBlockHeader
	ExecutionBlockHash  []byte
	ExecutionHashBranch [][]byte
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
	b3, err := rlp.EncodeToBytes(update.ExecutionHashBranch)
	if err != nil {
		return nil, err
	}
	var rlpBytes []byte
	rlpBytes = append(rlpBytes, b1...)
	rlpBytes = append(rlpBytes, b2...)
	rlpBytes = append(rlpBytes, b3...)
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
	SignatureSlot           primitives.Slot
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

type BeaconBlockHeaderData struct {
	Beacon struct {
		Slot          string `json:"slot"`
		ProposerIndex string `json:"proposer_index"`
		ParentRoot    string `json:"parent_root"`
		StateRoot     string `json:"state_root"`
		BodyRoot      string `json:"body_root"`
	} `json:"beacon"`
}

type SyncAggregateData struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}

type SyncCommitteeData struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
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

type LightClientUpdateMsg struct {
	Data LightClientUpdateData `json:"data"`
}

type LightClientUpdateDataNoCommittee struct {
	AttestedHeader  *BeaconBlockHeaderData `json:"attested_header"`
	FinalizedHeader *BeaconBlockHeaderData `json:"finalized_header"`
	FinalityBranch  []string               `json:"finality_branch"`
	SyncAggregate   *SyncAggregateData     `json:"sync_aggregate"`
	SignatureSlot   string                 `json:"signature_slot"`
}
