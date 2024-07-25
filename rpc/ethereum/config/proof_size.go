package config

import "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"

type ProofSize struct {
	BeaconBlockBodyTreeDepth                   uint64
	L1BeaconBlockBodyTreeExecutionPayloadIndex uint64
	L2ExecutionPayloadTreeExecutionBlockIndex  uint64
	L1BeaconBlockBodyProofSize                 uint64
	L2ExecutionPayloadProofSize                uint64
	ExecutionProofSize                         uint64
}

// NewProofSize returns a new ProofSize.
func NewProofSize(epoch primitives.Epoch) *ProofSize {

	return &ProofSize{}
}
