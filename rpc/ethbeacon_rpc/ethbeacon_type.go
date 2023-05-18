package ethbeacon_rpc

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	ssz "github.com/prysmaticlabs/fastssz"
	v11 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	pb "github.com/prysmaticlabs/prysm/v4/proto/eth/service"
	v2 "github.com/prysmaticlabs/prysm/v4/proto/eth/v2"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"net/http"
	"strings"
	"toprelayer/relayer/toprelayer/ethtypes"
)

const (
	ONE_EPOCH_IN_SLOTS = 32
	HEADER_BATCH_SIZE  = 128

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

func IsErrorNoBlockForSlot(err error) bool {
	return strings.Contains(err.Error(), ERROR_NO_BLOCK_FOR_SLOT)
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
	Data LightClientUpdateData `json:"data"`
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

func beaconBlockHeaderConvert(header *eth.BeaconBlockHeader) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:          uint64(header.Slot),
		ProposerIndex: uint64(header.ProposerIndex),
		ParentRoot:    header.ParentRoot,
		StateRoot:     header.StateRoot,
		BodyRoot:      header.BodyRoot,
	}
}

func convertEth2LightClientUpdate(lcu *ethtypes.LightClientUpdate) *LightClientUpdate {
	ret := &LightClientUpdate{
		AttestedBeaconHeader: beaconBlockHeaderConvert(lcu.AttestedBeaconHeader),
		SyncAggregate: &SyncAggregate{
			SyncCommitteeBits:      lcu.SyncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: lcu.SyncAggregate.SyncCommitteeSignature,
		},
		SignatureSlot: lcu.SignatureSlot,
		FinalizedUpdate: &FinalizedHeaderUpdate{
			HeaderUpdate: &HeaderUpdate{
				BeaconHeader:        beaconBlockHeaderConvert(lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader),
				ExecutionBlockHash:  lcu.FinalizedUpdate.HeaderUpdate.ExecutionBlockHash[:],
				ExecutionHashBranch: ethtypes.ConvertSliceHash2Bytes(lcu.FinalizedUpdate.HeaderUpdate.ExecutionHashBranch),
			},
			FinalityBranch: lcu.FinalizedUpdate.FinalityBranch,
		},
	}
	if lcu.NextSyncCommitteeUpdate != nil {
		ret.NextSyncCommitteeUpdate = &SyncCommitteeUpdate{
			NextSyncCommittee:       lcu.NextSyncCommitteeUpdate.NextSyncCommittee,
			NextSyncCommitteeBranch: lcu.NextSyncCommitteeUpdate.NextSyncCommitteeBranch,
		}
	}
	return ret
}

func SplitSlot(slot uint64) (period, epochInPeriod, slotInEpoch uint64) {
	period = GetPeriodForSlot(slot)
	currPeriodSlot := slot - (period * SLOTS_PER_EPOCH * EPOCHS_PER_PERIOD)
	epochInPeriod = currPeriodSlot / SLOTS_PER_EPOCH
	currPeriodEpochSlot := currPeriodSlot - epochInPeriod*SLOTS_PER_EPOCH
	return period, epochInPeriod, currPeriodEpochSlot % SLOTS_PER_EPOCH
}

func epochInPeriodForPeriod(period uint64) uint64 {
	batch := period * EPOCHS_PER_PERIOD / 154
	return (batch+1)*154 - (period * EPOCHS_PER_PERIOD)
}

func GetFinalizedSlotForPeriod(period uint64) uint64 {
	epoch := epochInPeriodForPeriod(period)
	return period*EPOCHS_PER_PERIOD*SLOTS_PER_EPOCH + epoch*ONE_EPOCH_IN_SLOTS
}

func addHexPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex
	}
	return "0x" + hex
}

func getBeforeSlotInSamePeriod(finalizedSlot uint64) (uint64, error) {
	slot := finalizedSlot - 4*ONE_EPOCH_IN_SLOTS
	period, epoch, _ := SplitSlot(slot)
	if epoch > 245 {
		slot = period*EPOCHS_PER_PERIOD*SLOTS_PER_EPOCH + 245*SLOTS_PER_EPOCH
	}
	if GetPeriodForSlot(slot) != GetPeriodForSlot(finalizedSlot) {
		return slot, fmt.Errorf("not an available slot:%d,it should be bigger", finalizedSlot)
	}
	return slot, nil
}

func getAttestationSlot(lastFinalizedSlotOnNear uint64) uint64 {
	nextFinalizedSlot := lastFinalizedSlotOnNear + ONE_EPOCH_IN_SLOTS
	return nextFinalizedSlot + 2*ONE_EPOCH_IN_SLOTS
}

const (
	BeaconBlockBodyTreeDepth  uint64 = 4
	ExecutionPayloadTreeDepth uint64 = 4

	L1BeaconBlockBodyTreeExecutionPayloadIndex uint64 = 9
	L1BeaconBlockBodyProofSize                 uint64 = BeaconBlockBodyTreeDepth

	L2ExecutionPayloadTreeExecutionBlockIndex uint64 = 12
	L2ExecutionPayloadProofSize               uint64 = ExecutionPayloadTreeDepth
)

func VecObjectHashTreeRoot(data []ssz.HashRoot, lenLimit uint64) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := vecObjectHashTreeRootWith(hh, data, lenLimit); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func vecObjectHashTreeRootWith(hh *ssz.Hasher, data []ssz.HashRoot, lenLimit uint64) (err error) {
	subIdx := hh.Index()
	num := uint64(len(data))
	if num > lenLimit {
		err = ssz.ErrIncorrectListSize
		return
	}
	for _, elem := range data {
		if err = elem.HashTreeRootWith(hh); err != nil {
			return
		}
	}
	if ssz.EnableVectorizedHTR {
		hh.MerkleizeWithMixinVectorizedHTR(subIdx, num, lenLimit)
	} else {
		hh.MerkleizeWithMixin(subIdx, num, lenLimit)
	}
	return nil
}

func BytesHashTreeRoot(data []byte, lenLimit int, remark string) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if size := len(data); size != lenLimit {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, ssz.ErrBytesLengthFn("--."+remark, size, lenLimit)
	}
	hh.PutBytes(data)
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func Uint64HashTreeRoot(data uint64) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	hh.PutUint64(data)
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func specialFieldExtraDataHashTreeRoot(extraData []byte) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	elemIdx := hh.Index()
	byteLen := uint64(len(extraData))
	if byteLen > 32 {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, ssz.ErrIncorrectListSize
	}
	hh.PutBytes(extraData)
	if ssz.EnableVectorizedHTR {
		hh.MerkleizeWithMixinVectorizedHTR(elemIdx, byteLen, (32+31)/32)
	} else {
		hh.MerkleizeWithMixin(elemIdx, byteLen, (32+31)/32)
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func specialFieldTransactionsHashTreeRoot(transactions [][]byte) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	subIdx := hh.Index()
	num := uint64(len(transactions))
	if num > 1048576 {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, ssz.ErrIncorrectListSize
	}
	for _, elem := range transactions {
		{
			elemIdx := hh.Index()
			byteLen := uint64(len(elem))
			if byteLen > 1073741824 {
				ssz.DefaultHasherPool.Put(hh)
				return [32]byte{}, ssz.ErrIncorrectListSize
			}
			hh.AppendBytes32(elem)
			if ssz.EnableVectorizedHTR {
				hh.MerkleizeWithMixinVectorizedHTR(elemIdx, byteLen, (1073741824+31)/32)
			} else {
				hh.MerkleizeWithMixin(elemIdx, byteLen, (1073741824+31)/32)
			}
		}
	}
	if ssz.EnableVectorizedHTR {
		hh.MerkleizeWithMixinVectorizedHTR(subIdx, num, 1048576)
	} else {
		hh.MerkleizeWithMixin(subIdx, num, 1048576)
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

//func BeaconBlockBodyMerkleTreeNew2(b *v2.BeaconBlockBodyCapella) (MerkleTreeNode, err error) {
//	hh := ssz.DefaultHasherPool.Get()
//	defer ssz.DefaultHasherPool.Put(hh)
//	//indx := hh.Index()
//	// Field (0) 'RandaoReveal'
//	if size := len(b.RandaoReveal); size != 96 {
//		err = ssz.ErrBytesLengthFn("--.RandaoReveal", size, 96)
//		return
//	}
//	hh.PutBytes(b.RandaoReveal)
//
//	// Field (1) 'Eth1Data'
//	if err = b.Eth1Data.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (2) 'Graffiti'
//	if size := len(b.Graffiti); size != 32 {
//		err = ssz.ErrBytesLengthFn("--.Graffiti", size, 32)
//		return
//	}
//	hh.PutBytes(b.Graffiti)
//
//	// Field (3) 'ProposerSlashings'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.ProposerSlashings))
//		if num > 16 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.ProposerSlashings {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 16)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 16)
//		}
//	}
//
//	// Field (4) 'AttesterSlashings'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.AttesterSlashings))
//		if num > 2 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.AttesterSlashings {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 2)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 2)
//		}
//	}
//
//	// Field (5) 'Attestations'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.Attestations))
//		if num > 128 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.Attestations {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 128)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 128)
//		}
//	}
//
//	// Field (6) 'Deposits'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.Deposits))
//		if num > 16 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.Deposits {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 16)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 16)
//		}
//	}
//
//	// Field (7) 'VoluntaryExits'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.VoluntaryExits))
//		if num > 16 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.VoluntaryExits {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 16)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 16)
//		}
//	}
//
//	// Field (8) 'SyncAggregate'
//	if err = b.SyncAggregate.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (9) 'ExecutionPayload'
//	if err = b.ExecutionPayload.HashTreeRootWith(hh); err != nil {
//		return
//	}
//
//	// Field (10) 'BlsToExecutionChanges'
//	{
//		subIndx := hh.Index()
//		num := uint64(len(b.BlsToExecutionChanges))
//		if num > 16 {
//			err = ssz.ErrIncorrectListSize
//			return
//		}
//		for _, elem := range b.BlsToExecutionChanges {
//			if err = elem.HashTreeRootWith(hh); err != nil {
//				return
//			}
//		}
//		if ssz.EnableVectorizedHTR {
//			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 16)
//		} else {
//			hh.MerkleizeWithMixin(subIndx, num, 16)
//		}
//	}
//
//	return create(hh.GetBuf(), BeaconBlockBodyTreeDepth), nil
//}

func BeaconBlockBodyMerkleTreeNew(b *v2.BeaconBlockBodyCapella) (MerkleTreeNode, error) {
	leaves := make([][32]byte, 11)
	// field 0
	if hashRoot, err := BytesHashTreeRoot(b.RandaoReveal, 96, "RandaoReveal"); err != nil {
		return nil, err
	} else {
		leaves[0] = hashRoot
	}
	// field 1
	if hashRoot, err := b.Eth1Data.HashTreeRoot(); err != nil {
		return nil, err
	} else {
		leaves[1] = hashRoot
	}

	// field 2
	if hashRoot, err := BytesHashTreeRoot(b.Graffiti, 32, "Graffiti"); err != nil {
		return nil, err
	} else {
		leaves[2] = hashRoot
	}

	// field 3
	hrs := make([]ssz.HashRoot, len(b.ProposerSlashings))
	for i, v := range b.ProposerSlashings {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[3] = hashRoot
	}

	// field 4
	hrs = make([]ssz.HashRoot, len(b.AttesterSlashings))
	for i, v := range b.AttesterSlashings {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 2); err != nil {
		return nil, err
	} else {
		leaves[4] = hashRoot
	}

	// field 5
	hrs = make([]ssz.HashRoot, len(b.Attestations))
	for i, v := range b.Attestations {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 128); err != nil {
		return nil, err
	} else {
		leaves[5] = hashRoot
	}

	// field 6
	hrs = make([]ssz.HashRoot, len(b.Deposits))
	for i, v := range b.Deposits {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[6] = hashRoot
	}

	// field 7
	hrs = make([]ssz.HashRoot, len(b.VoluntaryExits))
	for i, v := range b.VoluntaryExits {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[7] = hashRoot
	}

	// field 8
	if hashRoot, err := b.SyncAggregate.HashTreeRoot(); err != nil {
		return nil, err
	} else {
		leaves[8] = hashRoot
	}

	// field 9
	if hashRoot, err := b.ExecutionPayload.HashTreeRoot(); err != nil {
		return nil, err
	} else {
		leaves[9] = hashRoot
	}

	// field 10
	hrs = make([]ssz.HashRoot, len(b.BlsToExecutionChanges))
	for i, v := range b.BlsToExecutionChanges {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[10] = hashRoot
	}
	return create(leaves, BeaconBlockBodyTreeDepth), nil
}

func ExecutionPayloadMerkleTreeNew(b *v11.ExecutionPayloadCapella) (MerkleTreeNode, error) {

	leaves := make([][32]byte, 15)
	// field 0
	if hashRoot, err := BytesHashTreeRoot(b.ParentHash, 32, "ParentHash"); err != nil {
		return nil, err
	} else {
		leaves[0] = hashRoot
	}

	// field 1
	if hashRoot, err := BytesHashTreeRoot(b.FeeRecipient, 20, "FeeRecipient"); err != nil {
		return nil, err
	} else {
		leaves[1] = hashRoot
	}

	// field 2
	if hashRoot, err := BytesHashTreeRoot(b.StateRoot, 32, "StateRoot"); err != nil {
		return nil, err
	} else {
		leaves[2] = hashRoot
	}

	// field 3
	if hashRoot, err := BytesHashTreeRoot(b.ReceiptsRoot, 32, "ReceiptsRoot"); err != nil {
		return nil, err
	} else {
		leaves[3] = hashRoot
	}

	// field 4
	if hashRoot, err := BytesHashTreeRoot(b.LogsBloom, 256, "LogsBloom"); err != nil {
		return nil, err
	} else {
		leaves[4] = hashRoot
	}

	// field 5
	if hashRoot, err := BytesHashTreeRoot(b.PrevRandao, 32, "PrevRandao"); err != nil {
		return nil, err
	} else {
		leaves[5] = hashRoot
	}

	// field 6
	if hashRoot, err := Uint64HashTreeRoot(b.BlockNumber); err != nil {
		return nil, err
	} else {
		leaves[6] = hashRoot
	}

	// field 7
	if hashRoot, err := Uint64HashTreeRoot(b.GasLimit); err != nil {
		return nil, err
	} else {
		leaves[7] = hashRoot
	}

	// field 8
	if hashRoot, err := Uint64HashTreeRoot(b.GasUsed); err != nil {
		return nil, err
	} else {
		leaves[8] = hashRoot
	}

	// field 9
	if hashRoot, err := Uint64HashTreeRoot(b.Timestamp); err != nil {
		return nil, err
	} else {
		leaves[9] = hashRoot
	}

	// field 10
	if hashRoot, err := specialFieldExtraDataHashTreeRoot(b.ExtraData); err != nil {
		return nil, err
	} else {
		leaves[10] = hashRoot
	}

	// field 11
	if hashRoot, err := BytesHashTreeRoot(b.BaseFeePerGas, 32, "BaseFeePerGas"); err != nil {
		return nil, err
	} else {
		leaves[11] = hashRoot
	}

	// field 12
	if hashRoot, err := BytesHashTreeRoot(b.BlockHash, 32, "BlockHash"); err != nil {
		return nil, err
	} else {
		leaves[12] = hashRoot
	}

	// field 13
	if hashRoot, err := specialFieldTransactionsHashTreeRoot(b.Transactions); err != nil {
		return nil, err
	} else {
		leaves[13] = hashRoot
	}

	// field 14
	hrs := make([]ssz.HashRoot, len(b.Withdrawals))
	for i, v := range b.Withdrawals {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[14] = hashRoot
	}
	return create(leaves, BeaconBlockBodyTreeDepth), nil
}
