package ethashapp

import (
	"math/big"
	"os"

	"toprelayer/relayer/eth2top/ethash"
	"toprelayer/relayer/eth2top/ethashproof"
	"toprelayer/relayer/eth2top/mtree"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	// METHOD_GETBRIDGESTATE = "getCurrentBlockHeight"
	// SYNCHEADERS           = "syncBlockHeader"

	// SUCCESSDELAY int64 = 15 //mainnet 120
	// FATALTIMEOUT int64 = 24 //hours
	// FORKDELAY    int64 = 5  //mainnet 10000 seconds
	// ERRDELAY     int64 = 10
	// CONFIRMDELAY int64 = 5

	BLOCKS_PER_EPOCH uint64 = 30000
	// BLOCKS_TO_END_OF_EPOCH uint64 = 5000
)

type Output struct {
	HeaderRLP    string   `json:"header_rlp"`
	MerkleRoot   string   `json:"merkle_root"`
	Elements     []string `json:"elements"`
	MerkleProofs []string `json:"merkle_proofs"`
	ProofLength  uint64   `json:"proof_length"`
}

func zeroPad(bytes []byte, num int) []byte {
	for {
		if len(bytes) < num {
			bytes = append([]byte{uint8(0)}, bytes...)
		} else {
			break
		}
	}
	return bytes
}

func EthashWithProofs(h uint64, header *types.Header) (Output, error) {
	// var header *types.Header
	// if err := rlp.DecodeBytes(rlpheader, &header); err != nil {
	// 	logger.Error("RLP decoding of header failed: ", err)
	// 	return Output{}, err
	// }
	epoch := h / BLOCKS_PER_EPOCH
	cache, err := ethashproof.LoadCache(int(epoch))
	if err != nil {
		logger.Info("Cache is missing, calculate dataset merkle tree to create the cache first...")
		_, err = ethashproof.CalculateDatasetMerkleRoot(epoch, true)
		if err != nil {
			logger.Error("Creating cache failed: ", err)
			return Output{}, err
		}
		cache, err = ethashproof.LoadCache(int(epoch))
		if err != nil {
			logger.Error("Getting cache failed after trying to create it, abort: ", err)
			return Output{}, err
		}
	}

	// Remove outdated epoch
	if epoch > 1 {
		outdatedEpoch := epoch - 2
		err = os.Remove(ethash.PathToDAG(outdatedEpoch, ethash.DefaultDir))
		if err != nil {
			if os.IsNotExist(err) {
				logger.Info("DAG for previous epoch does not exist, nothing to remove: ", outdatedEpoch)
			} else {
				logger.Error("Remove DAG: ", err)
			}
		}

		err = os.Remove(ethashproof.PathToCache(outdatedEpoch))
		if err != nil {
			if os.IsNotExist(err) {
				logger.Info("Cache for previous epoch does not exist, nothing to remove: ", outdatedEpoch)
			} else {
				logger.Error("Remove cache error: ", err)
			}
		}
	}

	logger.Debug("SealHash: ", ethash.Instance.SealHash(header))
	indices := ethash.Instance.GetVerificationIndices(
		h,
		ethash.Instance.SealHash(header),
		header.Nonce.Uint64(),
	)
	logger.Debug("Proof length: ", cache.ProofLength)
	bytes, err := rlp.EncodeToBytes(header)
	if err != nil {
		logger.Error("RLP decoding of header failed: ", err)
		return Output{}, err
	}
	output := Output{
		HeaderRLP:    string(bytes),
		MerkleRoot:   string(cache.RootHash.Bytes()),
		Elements:     []string{},
		MerkleProofs: []string{},
		ProofLength:  cache.ProofLength,
	}
	for _, index := range indices {
		element, proof, err := ethashproof.CalculateProof(h, index, cache)
		if err != nil {
			logger.Error("calculating the proofs failed for index: %d, error: %s", index, err)
			return Output{}, err
		}
		es := element.ToUint256Array()
		for i := 0; i < len(es); i += 2 {
			eBytes := zeroPad(es[i].Bytes(), 32)
			eBytes = append(eBytes, zeroPad(es[i+1].Bytes(), 32)...)
			output.Elements = append(output.Elements, string(eBytes))
		}
		allProofs := []*big.Int{}
		for _, be := range mtree.HashesToBranchesArray(proof) {
			allProofs = append(allProofs, be.Big())
		}
		for _, pr := range allProofs {
			output.MerkleProofs = append(output.MerkleProofs, string(zeroPad(pr.Bytes(), 16)))
		}
	}

	return output, nil
}
