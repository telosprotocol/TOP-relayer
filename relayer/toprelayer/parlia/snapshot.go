// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package parlia

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/wonderivan/logger"
)

var (
	Epoch        uint64 = 200
	ValidatorNum uint64 = 21
)

// Snapshot is the state of the validatorSet at a given point.
type Snapshot struct {
	sigCache *lru.ARCCache // Cache of recent block signatures to speed up ecrecover

	Number           uint64                      `json:"number"`     // Block number where the snapshot was created
	Hash             common.Hash                 `json:"hash"`       // Block hash where the snapshot was created
	Validators       map[common.Address]struct{} `json:"validators"` // Set of authorized validators at this moment
	LastValidators   map[common.Address]struct{} `json:"last_validators"`
	Recents          map[uint64]common.Address   `json:"recents"`            // Set of recent validators for spam protections
	RecentForkHashes map[uint64]string           `json:"recent_fork_hashes"` // Set of recent forkHash
}

// newSnapshot creates a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use it for
// the genesis block.
func newSnapshot(
	sigCache *lru.ARCCache,
	number uint64,
	hash common.Hash,
	lastValidators []common.Address,
	validators []common.Address,
) *Snapshot {
	snap := &Snapshot{
		sigCache:         sigCache,
		Number:           number,
		Hash:             hash,
		Recents:          make(map[uint64]common.Address),
		RecentForkHashes: make(map[uint64]string),
		Validators:       make(map[common.Address]struct{}),
		LastValidators:   make(map[common.Address]struct{}),
	}
	for _, v := range lastValidators {
		snap.LastValidators[v] = struct{}{}
	}
	for _, v := range validators {
		snap.Validators[v] = struct{}{}
	}
	return snap
}

// validatorsAscending implements the sort interface to allow sorting a list of addresses
type validatorsAscending []common.Address

func (s validatorsAscending) Len() int           { return len(s) }
func (s validatorsAscending) Less(i, j int) bool { return bytes.Compare(s[i][:], s[j][:]) < 0 }
func (s validatorsAscending) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(sigCache *lru.ARCCache, db ethdb.Database, hash common.Hash) (*Snapshot, error) {
	blob, err := db.Get(append([]byte("parlia-"), hash[:]...))
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal(blob, snap); err != nil {
		return nil, err
	}
	snap.sigCache = sigCache

	return snap, nil
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db ethdb.Database) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return db.Put(append([]byte("parlia-"), s.Hash[:]...), blob)
}

// copy creates a deep copy of the snapshot
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		sigCache:         s.sigCache,
		Number:           s.Number,
		Hash:             s.Hash,
		LastValidators:   make(map[common.Address]struct{}),
		Validators:       make(map[common.Address]struct{}),
		Recents:          make(map[uint64]common.Address),
		RecentForkHashes: make(map[uint64]string),
	}

	for v := range s.LastValidators {
		cpy.LastValidators[v] = struct{}{}
	}
	for v := range s.Validators {
		cpy.Validators[v] = struct{}{}
	}
	for block, v := range s.Recents {
		cpy.Recents[block] = v
	}
	for block, id := range s.RecentForkHashes {
		cpy.RecentForkHashes[block] = id
	}
	return cpy
}

func (s *Snapshot) isMajorityFork(forkHash string) bool {
	ally := 0
	for _, h := range s.RecentForkHashes {
		if h == forkHash {
			ally++
		}
	}
	return ally > len(s.RecentForkHashes)/2
}

func (s *Snapshot) apply(headers []*types.Header, chainId *big.Int) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errOutOfRangeChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errOutOfRangeChain
	}
	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	for _, header := range headers {
		number := header.Number.Uint64()
		// Delete the oldest validator from the recent list to allow it signing again
		if limit := uint64(len(snap.Validators)/2 + 1); number >= limit {
			delete(snap.Recents, number-limit)
		}
		// Resolve the authorization key and check against signers
		validator, err := ecrecover(header, s.sigCache, chainId)
		if err != nil {
			return nil, err
		}
		_, ok1 := snap.Validators[validator]
		_, ok2 := snap.LastValidators[validator]
		if !ok1 && !ok2 {
			return nil, errUnauthorizedValidator
		}
		for _, recent := range snap.Recents {
			if recent == validator {
				return nil, errRecentlySigned
			}
		}
		snap.Recents[number] = validator
		// change validator set
		if number > 0 && (number%Epoch == 0) {
			checkpointHeader := header
			if checkpointHeader == nil {
				return nil, consensus.ErrUnknownAncestor
			}

			validatorBytes := checkpointHeader.Extra[extraVanity : len(checkpointHeader.Extra)-extraSeal]
			// get validators from headers and use that for new validator set
			newValArr, err := ParseValidators(validatorBytes)
			if err != nil {
				return nil, err
			}
			newVals := make(map[common.Address]struct{}, len(newValArr))
			for _, val := range newValArr {
				newVals[val] = struct{}{}
			}
			oldLimit := len(snap.Validators)/2 + 1
			newLimit := len(newVals)/2 + 1
			if newLimit < oldLimit {
				for i := 0; i < oldLimit-newLimit; i++ {
					delete(snap.Recents, number-uint64(newLimit)-uint64(i))
				}
			}
			snap.LastValidators = snap.Validators
			snap.Validators = newVals
		}
	}
	snap.Number += uint64(len(headers))
	snap.Hash = headers[len(headers)-1].Hash()
	return snap, nil
}

// validators retrieves the list of validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, len(s.Validators))
	for v := range s.Validators {
		validators = append(validators, v)
	}
	sort.Sort(validatorsAscending(validators))
	return validators
}

func ParseValidators(validatorsBytes []byte) ([]common.Address, error) {
	if len(validatorsBytes)%validatorBytesLength != 0 {
		return nil, errors.New("invalid validators bytes")
	}
	n := len(validatorsBytes) / validatorBytesLength
	result := make([]common.Address, n)
	for i := 0; i < n; i++ {
		address := make([]byte, validatorBytesLength)
		copy(address, validatorsBytes[i*validatorBytesLength:(i+1)*validatorBytesLength])
		result[i] = common.BytesToAddress(address)
	}
	return result, nil
}

type SnapshotOut struct {
	Header        []byte
	ValidatorsNum uint64
	Validators    [][]byte
	RecentsNum    uint64
	Recents       [][]byte
}

func encodeSnapshot(header *types.Header, snap *Snapshot) ([]byte, error) {
	out := new(SnapshotOut)

	headerRlp, err := rlp.EncodeToBytes(header)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	out.Header = headerRlp
	out.ValidatorsNum = uint64(len(snap.Validators))
	for k := range snap.Validators {
		out.Validators = append(out.Validators, k.Bytes())
	}
	out.RecentsNum = uint64(len(snap.Recents))
	for k, v := range snap.Recents {
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(k))
		out.Recents = append(out.Recents, buf)
		out.Recents = append(out.Recents, v.Bytes())
	}
	bytes, err := rlp.EncodeToBytes(out)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return bytes, nil
}
