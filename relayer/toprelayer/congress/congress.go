package congress

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/sha3"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
	maxValidators      = 21   // Max validators allowed to seal.
)

// Congress proof-of-stake-authority protocol constants.
var (
	extraVanity = 32                     // Fixed number of extra-data prefix bytes reserved for validator vanity
	extraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for validator seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")

	// errUnauthorizedValidator is returned if a header is signed by a non-authorized entity.
	errUnauthorizedValidator = errors.New("unauthorized validator")

	// errRecentlySigned is returned if a header is signed by an authorized entity
	// that already signed a header recently, thus is temporarily not allowed to.
	errRecentlySigned = errors.New("recently signed")
)

// TODO: add db
type Congress struct {
	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	client *ethclient.Client
}

// New creates a Congress proof-of-stake-authority consensus engine with the initial
// validators set to the ones provided by the user.
func New(client *ethclient.Client) *Congress {
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Congress{
		recents:    recents,
		signatures: signatures,
		client:     client,
	}
}

func (c *Congress) Init(height uint64) error {
	var baseHeight uint64
	if height < Epoch {
		baseHeight = 0
	} else {
		if height%Epoch >= ValidatorNum {
			baseHeight = height / Epoch * Epoch
		} else {
			baseHeight = (height/Epoch - 1) * Epoch
		}
	}

	logger.Info("initing congress snapshot from %v to %v", baseHeight, height)
	// init baseheight
	{
		header, err := c.client.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(baseHeight))
		if err != nil {
			logger.Error(err)
			return err
		}
		hash := header.Hash()

		validators := make([]common.Address, (len(header.Extra)-extraVanity-extraSeal)/common.AddressLength)
		for i := 0; i < len(validators); i++ {
			copy(validators[i][:], header.Extra[extraVanity+i*common.AddressLength:])
		}
		snap := newSnapshot(c.signatures, baseHeight, hash, validators)
		c.recents.Add(snap.Hash, snap)
	}

	for i := baseHeight + 1; i <= height; i++ {
		header, err := c.client.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(i))
		if err != nil {
			logger.Error(err)
			return err
		}
		snap, err := c.GetLastSnap(header.Number.Uint64()-1, header.ParentHash)
		if err != nil {
			logger.Error(err)
			return err
		}
		err = c.Apply(snap, header)
		if err != nil {
			logger.Error(err)
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

func (c *Congress) GetLastSnap(number uint64, hash common.Hash) (*Snapshot, error) {
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := c.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// TODO: db
		if number%checkpointInterval == 0 {
		}
		if number == 0 || (number%Epoch == 0 && len(headers) >= int(maxValidators)) {
			checkpoint, err := c.client.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(number))
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			if checkpoint == nil {
				logger.Error(err)
				return nil, fmt.Errorf("header is nil")
			}
			hash := checkpoint.Hash()
			validators := make([]common.Address, (len(checkpoint.Extra)-extraVanity-extraSeal)/common.AddressLength)
			for i := 0; i < len(validators); i++ {
				copy(validators[i][:], checkpoint.Extra[extraVanity+i*common.AddressLength:])
			}
			snap = newSnapshot(c.signatures, number, hash, validators)
			// TODO:db
			break
		}
		h, err := c.client.HeaderByHash(context.Background(), hash)
		if err != nil {
			logger.Error(err)
			return nil, fmt.Errorf("HeaderByHash error")
		}
		headers = append(headers, h)
		number, hash = number-1, h.ParentHash
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}
	snap, err := snap.apply(headers)
	if err != nil {
		return nil, err
	}
	c.recents.Add(snap.Hash, snap)
	logger.Debug(snap)
	// TODO:db
	return snap, err
}

func (c *Congress) GetLastSnapBytes(header *types.Header) ([]byte, error) {
	snap, err := c.GetLastSnap(header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	bytes, err := encodeSnapshot(header, snap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return bytes, nil
}

func (c *Congress) Apply(snap *Snapshot, header *types.Header) error {
	var headers []*types.Header
	headers = append(headers, header)
	snap, err := snap.apply(headers)
	if err != nil {
		return err
	}
	c.recents.Add(snap.Hash, snap)
	return nil
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var validator common.Address
	copy(validator[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, validator)
	return validator, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-crypto.SignatureLength], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}
