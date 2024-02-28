package config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"math"
)

const (
	VersionLength = 4
	SlotsPerEpoch = 32
)

type Version [VersionLength]byte
type Slot uint64
type Epoch uint64

func (s Slot) Epoch() Epoch {
	return Epoch(s / SlotsPerEpoch)
}

func (s Slot) ToPrysmSlot() primitives.Slot {
	return primitives.Slot(s)
}

func FromPrysmSlot(s primitives.Slot) Slot {
	return Slot(s)
}

func (e Epoch) ToPrysmEpoch() primitives.Epoch {
	return primitives.Epoch(e)
}

func FromPrysmEpoch(e primitives.Epoch) Epoch {
	return Epoch(e)
}

type NetworkConfig struct {
	GenesisValidatorsRoot common.Hash

	GenesisForkVersion Version
	GenesisForkEpoch   Epoch

	AltairForkVersion Version
	AltairForkEpoch   Epoch

	BellatrixForkVersion Version
	BellatrixForkEpoch   Epoch

	CapellaForkVersion Version
	CapellaForkEpoch   Epoch

	DenebForkVersion Version
	DenebForkEpoch   Epoch
}

func NewNetworkConfig(networkId NetworkId) *NetworkConfig {
	switch networkId {
	case MAINNET:
		return &NetworkConfig{
			GenesisValidatorsRoot: common.Hash{},
			GenesisForkVersion:    Version{0x00, 0x00, 0x00, 0x00},
			GenesisForkEpoch:      0,
			AltairForkVersion:     Version{0x01, 0x00, 0x00, 0x00},
			AltairForkEpoch:       74240,
			BellatrixForkVersion:  Version{0x02, 0x00, 0x00, 0x00},
			BellatrixForkEpoch:    144896,
			CapellaForkVersion:    Version{0x03, 0x00, 0x00, 0x00},
			CapellaForkEpoch:      194048,
			DenebForkVersion:      Version{0x04, 0x00, 0x00, 0x00},
			DenebForkEpoch:        Epoch(math.MaxUint64),
		}
	case SEPOLIA:
		return &NetworkConfig{
			GenesisValidatorsRoot: common.Hash{},
			GenesisForkVersion:    Version{0x90, 0x00, 0x00, 0x69},
			GenesisForkEpoch:      0,
			AltairForkVersion:     Version{0x90, 0x00, 0x00, 0x70},
			AltairForkEpoch:       50,
			BellatrixForkVersion:  Version{0x90, 0x00, 0x00, 0x71},
			BellatrixForkEpoch:    100,
			CapellaForkVersion:    Version{0x90, 0x00, 0x00, 0x72},
			CapellaForkEpoch:      56832,
			DenebForkVersion:      Version{0x90, 0x00, 0x00, 0x73},
			DenebForkEpoch:        132608,
		}
	default:
		return nil
	}
}

func (nc *NetworkConfig) ComputeForkVersion(epoch Epoch) Version {
	if epoch >= nc.DenebForkEpoch {
		return nc.DenebForkVersion
	}
	if epoch >= nc.CapellaForkEpoch {
		return nc.CapellaForkVersion
	}
	if epoch >= nc.BellatrixForkEpoch {
		return nc.BellatrixForkVersion
	}
	if epoch >= nc.AltairForkEpoch {
		return nc.AltairForkVersion
	}
	return nc.GenesisForkVersion
}

func (nc *NetworkConfig) ComputeForkVersionBySlot(slot Slot) Version {
	return nc.ComputeForkVersion(slot.Epoch())
}

func (nc *NetworkConfig) ComputeProofSize(epoch Epoch) *ProofSize {
	if epoch >= nc.DenebForkEpoch {
		return &ProofSize{
			BeaconBlockBodyTreeDepth:                   4,
			L1BeaconBlockBodyTreeExecutionPayloadIndex: 9,
			L2ExecutionPayloadTreeExecutionBlockIndex:  12,
			L1BeaconBlockBodyProofSize:                 4,
			L2ExecutionPayloadProofSize:                5,
			ExecutionProofSize:                         9,
		}
	}

	return &ProofSize{
		BeaconBlockBodyTreeDepth:                   4,
		L1BeaconBlockBodyTreeExecutionPayloadIndex: 9,
		L2ExecutionPayloadTreeExecutionBlockIndex:  12,
		L1BeaconBlockBodyProofSize:                 4,
		L2ExecutionPayloadProofSize:                4,
		ExecutionProofSize:                         8,
	}
}

func (nc *NetworkConfig) ComputeProofSizeBySlot(slot Slot) *ProofSize {
	return nc.ComputeProofSize(slot.Epoch())
}
