package ethereum

import (
	"fmt"
	"strings"
	"toprelayer/relayer/toprelayer/ethtypes"
	"toprelayer/rpc/ethereum/light_client"

	fieldparams "github.com/OffchainLabs/prysm/v6/config/fieldparams"
	"github.com/OffchainLabs/prysm/v6/consensus-types/interfaces"
	"github.com/OffchainLabs/prysm/v6/consensus-types/primitives"
	eth "github.com/OffchainLabs/prysm/v6/proto/prysm/v1alpha1"
	ssz "github.com/prysmaticlabs/fastssz"
	"github.com/wonderivan/logger"
)

const (
	ONE_EPOCH_IN_SLOTS = 32
	HEADER_BATCH_SIZE  = 128

	SLOTS_PER_EPOCH   = 32
	EPOCHS_PER_PERIOD = 256

	ERROR_NO_BLOCK_FOR_SLOT = "not find requested block"
)

const (
	ExecutionPayloadTreeDepth uint64 = 4

	L1BeaconBlockBodyTreeExecutionPayloadIndex uint64 = 9

	L2ExecutionPayloadTreeExecutionBlockIndex uint64 = 12
	L2ExecutionPayloadProofSize               uint64 = ExecutionPayloadTreeDepth
)

func GetPeriodForSlot(slot primitives.Slot) uint64 {
	return uint64(slot) / (SLOTS_PER_EPOCH * EPOCHS_PER_PERIOD)
}

func GetEpochForSlot(slot primitives.Slot) primitives.Epoch {
	return primitives.Epoch(uint64(slot) / SLOTS_PER_EPOCH)
}

func GetPeriodForEpoch(epoch primitives.Epoch) uint64 {
	return uint64(epoch) / EPOCHS_PER_PERIOD
}

func epochInPeriodForPeriod(period uint64) primitives.Epoch {
	batch := period * EPOCHS_PER_PERIOD / 154
	return primitives.Epoch((batch+1)*154 - (period * EPOCHS_PER_PERIOD))
}

func GetFinalizedSlotForPeriod(period uint64) primitives.Slot {
	epoch := epochInPeriodForPeriod(period)
	return primitives.Slot(period*EPOCHS_PER_PERIOD*SLOTS_PER_EPOCH + uint64(epoch)*ONE_EPOCH_IN_SLOTS)
}

func IsErrorNoBlockForSlot(err error) bool {
	return strings.Contains(err.Error(), ERROR_NO_BLOCK_FOR_SLOT)
}

func getBeforeSlotInSamePeriod(finalizedSlot primitives.Slot) (primitives.Slot, error) {
	slot := finalizedSlot - 3*ONE_EPOCH_IN_SLOTS

	if GetPeriodForSlot(slot) != GetPeriodForSlot(finalizedSlot) {
		return slot, fmt.Errorf("not an available slot:%d,it should be bigger", finalizedSlot)
	}
	return slot, nil
}

func getAttestationSlot(lastFinalizedSlotOnTop primitives.Slot) primitives.Slot {
	nextFinalizedSlot := lastFinalizedSlotOnTop + ONE_EPOCH_IN_SLOTS
	return nextFinalizedSlot + 2*ONE_EPOCH_IN_SLOTS
}

func BytesHashTreeRoot(data []byte, lenLimit int, remark string) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)

	if size := len(data); size != lenLimit {
		return [32]byte{}, ssz.ErrBytesLengthFn("--."+remark, size, lenLimit)
	}
	hh.PutBytes(data)
	root, err := hh.HashRoot()
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

	hh.MerkleizeWithMixin(subIdx, num, lenLimit)

	return nil
}

func VecObjectHashTreeRoot(data []ssz.HashRoot, lenLimit uint64) ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)

	if err := vecObjectHashTreeRootWith(hh, data, lenLimit); err != nil {
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
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
	hh.MerkleizeWithMixin(elemIdx, byteLen, (32+31)/32)
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
			hh.MerkleizeWithMixin(elemIdx, byteLen, (1073741824+31)/32)
		}
	}
	hh.MerkleizeWithMixin(subIdx, num, 1048576)
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func specialFieldBlobKzgCommitmentsHashTreeRoot(kzgCommitments [][]byte) ([32]byte, error) {
	numItems := uint64(len(kzgCommitments))
	if numItems > 4096 {
		return [32]byte{}, ssz.ErrListTooBigFn("--.BlobKzgCommitments", int(numItems), 4096)
	}

	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)

	subIndx := hh.Index()
	for _, i := range kzgCommitments {
		if len(i) != 48 {
			return [32]byte{}, ssz.ErrBytesLength
		}
		hh.PutBytes(i)
	}

	hh.MerkleizeWithMixin(subIndx, numItems, 4096)

	return hh.HashRoot()
}

func ExecutionPayloadMerkleTreeShanghai(executionData interfaces.ExecutionData) (MerkleTreeNode, error) {
	leaves := make([][32]byte, 15)

	// Field (0) 'ParentHash'
	parentHash := executionData.ParentHash()
	if hashRoot, err := BytesHashTreeRoot(parentHash, 32, "ParentHash"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(parentHash) error ", err)
		return nil, err
	} else {
		leaves[0] = hashRoot
	}

	// Field (1) 'FeeRecipient'
	feeRecipient := executionData.FeeRecipient()
	if hashRoot, err := BytesHashTreeRoot(feeRecipient, 20, "FeeRecipient"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(feeRecipient) error ", err)
		return nil, err
	} else {
		leaves[1] = hashRoot
	}

	// Field (2) 'StateRoot'
	stateRoot := executionData.StateRoot()
	if hashRoot, err := BytesHashTreeRoot(stateRoot, 32, "StateRoot"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(stateRoot) error ", err)
		return nil, err
	} else {
		leaves[2] = hashRoot
	}

	// Field (3) 'ReceiptsRoot'
	receiptsRoot := executionData.ReceiptsRoot()
	if hashRoot, err := BytesHashTreeRoot(receiptsRoot, 32, "ReceiptsRoot"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(receiptsRoot) error ", err)
		return nil, err
	} else {
		leaves[3] = hashRoot
	}

	// Field (4) 'LogsBloom'
	logsBloom := executionData.LogsBloom()
	if hashRoot, err := BytesHashTreeRoot(logsBloom, 256, "LogsBloom"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(logsBloom) error ", err)
		return nil, err
	} else {
		leaves[4] = hashRoot
	}

	// Field (5) 'PrevRandao'
	prevRandao := executionData.PrevRandao()
	if hashRoot, err := BytesHashTreeRoot(prevRandao, 32, "PrevRandao"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(prevRandao) error ", err)
		return nil, err
	} else {
		leaves[5] = hashRoot
	}

	// Field (6) 'BlockNumber'
	if hashRoot, err := Uint64HashTreeRoot(executionData.BlockNumber()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew Uint64HashTreeRoot(executionData.BlockNumber()) error ", err)
		return nil, err
	} else {
		leaves[6] = hashRoot
	}

	// Field (7) 'GasLimit'
	if hashRoot, err := Uint64HashTreeRoot(executionData.GasLimit()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew Uint64HashTreeRoot(executionData.GasLimit()) error ", err)
		return nil, err
	} else {
		leaves[7] = hashRoot
	}

	// Field (8) 'GasUsed'
	if hashRoot, err := Uint64HashTreeRoot(executionData.GasUsed()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew Uint64HashTreeRoot(executionData.GasUsed()) error ", err)
		return nil, err
	} else {
		leaves[8] = hashRoot
	}

	// Field (9) 'Timestamp'
	if hashRoot, err := Uint64HashTreeRoot(executionData.Timestamp()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew Uint64HashTreeRoot(executionData.Timestamp()) error ", err)
		return nil, err
	} else {
		leaves[9] = hashRoot
	}

	// Field (10) 'ExtraData'
	if hashRoot, err := specialFieldExtraDataHashTreeRoot(executionData.ExtraData()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew specialFieldExtraDataHashTreeRoot(executionData.ExtraData() error ", err)
		return nil, err
	} else {
		leaves[10] = hashRoot
	}

	// Field (11) 'BaseFeePerGas'
	baseFeePerGas := executionData.BaseFeePerGas()
	if hashRoot, err := BytesHashTreeRoot(baseFeePerGas, len(baseFeePerGas), "BaseFeePerGas"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(baseFeePerGas) error ", err)
		return nil, err
	} else {
		leaves[11] = hashRoot
	}

	// Field (12) 'BlockHash'
	blockHash := executionData.BlockHash()
	if hashRoot, err := BytesHashTreeRoot(blockHash, len(blockHash), "BlockHash"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(blockHash) error ", err)
		return nil, err
	} else {
		leaves[12] = hashRoot
	}

	// Field (13) 'Transactions'
	transactions, err := executionData.Transactions()
	leaves[13] = [32]byte{}
	if err != nil {
		logger.Error("ExecutionPayloadMerkleTreeNew BytesHashTreeRoot(blockHash) error ", err)
	} else {
		if hashRoot, err := specialFieldTransactionsHashTreeRoot(transactions); err != nil {
			logger.Error("ExecutionPayloadMerkleTreeNew specialFieldTransactionsHashTreeRoot(transactions) error ", err)
		} else {
			leaves[13] = hashRoot
		}
	}

	// Field (14) 'Withdrawals'
	leaves[14] = [32]byte{0}
	withdrawals, err := executionData.Withdrawals()
	hrs := make([]ssz.HashRoot, len(withdrawals))
	for i, v := range withdrawals {
		hrs[i] = v
	}
	if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err != nil {
		return nil, err
	} else {
		leaves[14] = hashRoot
	}
	return create(leaves, ExecutionPayloadTreeDepth), nil
}

func ExecutionPayloadMerkleTreeCancun(executionData interfaces.ExecutionData) (MerkleTreeNode, error) {
	var depth = ExecutionPayloadTreeDepth
	leaves := make([][32]byte, 15)

	// Field (0) 'ParentHash'
	parentHash := executionData.ParentHash()
	if hashRoot, err := BytesHashTreeRoot(parentHash, 32, "ParentHash"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(parentHash) error ", err)
		return nil, err
	} else {
		leaves[0] = hashRoot
	}

	// Field (1) 'FeeRecipient'
	feeRecipient := executionData.FeeRecipient()
	if hashRoot, err := BytesHashTreeRoot(feeRecipient, 20, "FeeRecipient"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(feeRecipient) error ", err)
		return nil, err
	} else {
		leaves[1] = hashRoot
	}

	// Field (2) 'StateRoot'
	stateRoot := executionData.StateRoot()
	if hashRoot, err := BytesHashTreeRoot(stateRoot, 32, "StateRoot"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(stateRoot) error ", err)
		return nil, err
	} else {
		leaves[2] = hashRoot
	}

	// Field (3) 'ReceiptsRoot'
	receiptsRoot := executionData.ReceiptsRoot()
	if hashRoot, err := BytesHashTreeRoot(receiptsRoot, 32, "ReceiptsRoot"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(receiptsRoot) error ", err)
		return nil, err
	} else {
		leaves[3] = hashRoot
	}

	// Field (4) 'LogsBloom'
	logsBloom := executionData.LogsBloom()
	if hashRoot, err := BytesHashTreeRoot(logsBloom, 256, "LogsBloom"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(logsBloom) error ", err)
		return nil, err
	} else {
		leaves[4] = hashRoot
	}

	// Field (5) 'PrevRandao'
	prevRandao := executionData.PrevRandao()
	if hashRoot, err := BytesHashTreeRoot(prevRandao, 32, "PrevRandao"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(prevRandao) error ", err)
		return nil, err
	} else {
		leaves[5] = hashRoot
	}

	// Field (6) 'BlockNumber'
	if hashRoot, err := Uint64HashTreeRoot(executionData.BlockNumber()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(executionData.BlockNumber()) error ", err)
		return nil, err
	} else {
		leaves[6] = hashRoot
	}

	// Field (7) 'GasLimit'
	if hashRoot, err := Uint64HashTreeRoot(executionData.GasLimit()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(executionData.GasLimit()) error ", err)
		return nil, err
	} else {
		leaves[7] = hashRoot
	}

	// Field (8) 'GasUsed'
	if hashRoot, err := Uint64HashTreeRoot(executionData.GasUsed()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(executionData.GasUsed()) error ", err)
		return nil, err
	} else {
		leaves[8] = hashRoot
	}

	// Field (9) 'Timestamp'
	if hashRoot, err := Uint64HashTreeRoot(executionData.Timestamp()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(executionData.Timestamp()) error ", err)
		return nil, err
	} else {
		leaves[9] = hashRoot
	}

	// Field (10) 'ExtraData'
	if hashRoot, err := specialFieldExtraDataHashTreeRoot(executionData.ExtraData()); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun specialFieldExtraDataHashTreeRoot(executionData.ExtraData() error ", err)
		return nil, err
	} else {
		leaves[10] = hashRoot
	}

	// Field (11) 'BaseFeePerGas'
	baseFeePerGas := executionData.BaseFeePerGas()
	if hashRoot, err := BytesHashTreeRoot(baseFeePerGas, len(baseFeePerGas), "BaseFeePerGas"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(baseFeePerGas) error ", err)
		return nil, err
	} else {
		leaves[11] = hashRoot
	}

	// Field (12) 'BlockHash'
	blockHash := executionData.BlockHash()
	if hashRoot, err := BytesHashTreeRoot(blockHash, len(blockHash), "BlockHash"); err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun BytesHashTreeRoot(blockHash) error ", err)
		return nil, err
	} else {
		leaves[12] = hashRoot
	}

	// Field (13) 'Transactions'
	transactions, err := executionData.Transactions()
	leaves[13] = [32]byte{}
	if err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun executionData.Transactions() error ", err)
	} else {
		if hashRoot, err := specialFieldTransactionsHashTreeRoot(transactions); err != nil {
			logger.Error("ExecutionPayloadMerkleTreeCancun specialFieldTransactionsHashTreeRoot(transactions) error ", err)
		} else {
			leaves[13] = hashRoot
		}
	}

	// Field (14) 'Withdrawals'
	leaves[14] = [32]byte{0}
	withdrawals, err := executionData.Withdrawals()
	if err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun executionData.Withdrawals() error: ", err)
	} else {
		hrs := make([]ssz.HashRoot, len(withdrawals))
		for i, v := range withdrawals {
			hrs[i] = v
		}
		if hashRoot, err := VecObjectHashTreeRoot(hrs, 16); err == nil {
			leaves[14] = hashRoot
		}
	}

	// Field (15) 'BlobGasUsed'
	blobGasUsed, err := executionData.BlobGasUsed()
	if err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun executionData.BlobGasUsed() error: ", err)
	} else {
		if hashRoot, err := Uint64HashTreeRoot(blobGasUsed); err != nil {
			logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(blobGasUsed) error: ", err)
		} else {
			depth += 1
			leaves = append(leaves, hashRoot)
		}
	}

	// Field (16) 'ExcessBlobGas'
	excessBlobGas, err := executionData.ExcessBlobGas()
	if err != nil {
		logger.Error("ExecutionPayloadMerkleTreeCancun executionData.ExcessBlobGas() error: ", err)
	} else {
		if hashRoot, err := Uint64HashTreeRoot(excessBlobGas); err != nil {
			logger.Error("ExecutionPayloadMerkleTreeCancun Uint64HashTreeRoot(excessBlobGas) error: ", err)
		} else {
			leaves = append(leaves, hashRoot)
		}
	}

	return create(leaves, depth), nil
}

func beaconBlockHeaderConvert(header *eth.BeaconBlockHeader) *light_client.BeaconBlockHeader {
	return &light_client.BeaconBlockHeader{
		Slot:          header.Slot,
		ProposerIndex: header.ProposerIndex,
		ParentRoot:    [32]byte(header.ParentRoot),
		StateRoot:     [32]byte(header.StateRoot),
		BodyRoot:      [32]byte(header.BodyRoot),
	}
}

func convertEth2LightClientUpdate(lcu *ethtypes.LightClientUpdate) *light_client.LightClientUpdate {
	var executionHashBranch = make([][fieldparams.RootLength]byte, len(lcu.FinalizedUpdate.HeaderUpdate.ExecutionHashBranch))
	for i, v := range lcu.FinalizedUpdate.HeaderUpdate.ExecutionHashBranch {
		executionHashBranch[i] = v
	}

	ret := &light_client.LightClientUpdate{
		AttestedBeaconHeader: beaconBlockHeaderConvert(lcu.AttestedBeaconHeader),
		SyncAggregate: &light_client.SyncAggregate{
			SyncCommitteeBits:      [fieldparams.SyncAggregateSyncCommitteeBytesLength]byte(lcu.SyncAggregate.SyncCommitteeBits),
			SyncCommitteeSignature: [fieldparams.BLSSignatureLength]byte(lcu.SyncAggregate.SyncCommitteeSignature),
		},
		SignatureSlot: primitives.Slot(lcu.SignatureSlot),
		FinalityUpdate: &light_client.FinalizedHeaderUpdate{
			HeaderUpdate: &light_client.HeaderUpdate{
				BeaconHeader:        beaconBlockHeaderConvert(lcu.FinalizedUpdate.HeaderUpdate.BeaconHeader),
				ExecutionBlockHash:  lcu.FinalizedUpdate.HeaderUpdate.ExecutionBlockHash,
				ExecutionHashBranch: executionHashBranch,
			},
			FinalityBranch: lcu.FinalizedUpdate.FinalityBranch,
		},
	}
	if lcu.NextSyncCommitteeUpdate != nil {
		ret.NextSyncCommitteeUpdate = &light_client.SyncCommitteeUpdate{
			NextSyncCommittee:       lcu.NextSyncCommitteeUpdate.NextSyncCommittee,
			NextSyncCommitteeBranch: lcu.NextSyncCommitteeUpdate.NextSyncCommitteeBranch,
		}
	}
	return ret
}
