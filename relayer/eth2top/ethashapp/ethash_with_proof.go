package ethashapp

import (
	"math/big"
	"os"
	"time"

	"toprelayer/relayer/eth2top/ethash"
	"toprelayer/relayer/eth2top/ethashproof"
	"toprelayer/relayer/eth2top/mtree"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const (
	BLOCKS_PER_EPOCH uint64 = 30000
)

var futureEpoch uint64 = 0
var futureEpochProcessing bool = false

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
	epoch := h / BLOCKS_PER_EPOCH
	cache, err := ethashproof.LoadCache(int(epoch))
	if err != nil {
		if futureEpochProcessing {
			logger.Debug("waiting for futureEpochProcessing...")
			totalTime := 0
			for {
				time.Sleep(time.Duration(5) * time.Second)
				totalTime += 5
				if !futureEpochProcessing || totalTime > 3600 {
					break
				}
			}
			cache, err = ethashproof.LoadCache(int(epoch))
		}

		if err != nil {
			logger.Info("epoch %v cache is missing, calculate dataset merkle tree to create the cache first...", epoch)
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
	}

	if futureEpoch != epoch+1 {
		if !futureEpochProcessing && !ethashproof.ExistCache(int(epoch+1)) {
			futureEpochProcessing = true
			logger.Info("future epoch %v cache is missing, calculate dataset merkle tree to create the cache first...", epoch+1)
			go func() {
				_, e := ethashproof.CalculateDatasetMerkleRoot(epoch+1, true)
				if e != nil || !ethashproof.ExistCache(int(epoch+1)) {
					logger.Error("Creating cache failed: ", err)
				} else {
					futureEpoch = epoch + 1
				}
				futureEpochProcessing = false
			}()
		}
	}

	// Remove outdated epoch
	if epoch > 1 {
		outdatedEpoch := epoch - 2
		err = os.Remove(ethash.PathToDAG(outdatedEpoch, ethash.DefaultDir))
		if err != nil {
			if os.IsNotExist(err) {
			} else {
				logger.Error("Remove DAG: ", err)
			}
		}

		err = os.Remove(ethashproof.PathToCache(outdatedEpoch))
		if err != nil {
			if os.IsNotExist(err) {
			} else {
				logger.Error("Remove cache error: ", err)
			}
		}
	}

	indices := ethash.Instance.GetVerificationIndices(
		h,
		ethash.Instance.SealHash(header),
		header.Nonce.Uint64(),
	)
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
