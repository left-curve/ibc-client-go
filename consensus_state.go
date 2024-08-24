package grug

import (
	"time"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

// Assert that `ConsensusState` implements the ICS-02 `ConsensusState` interface.
var _ exported.ConsensusState = (*ConsensusState)(nil)

// ConsensusState is the client state of Grug.
//
// This is a wrapper of the 07-tendermint consensus state, with with a few
// methods substituted by the Grug ones.
type ConsensusState struct{ *tendermint.ConsensusState }

// NewConsensusState creates a new Grug consensus state instance.
func NewConsensusState(
	timestamp time.Time, root commitmenttypes.MerkleRoot, nextValsHash tmbytes.HexBytes,
) *ConsensusState {
	return &ConsensusState{
		tendermint.NewConsensusState(timestamp, root, nextValsHash),
	}
}

// ClientType implements the ICS-02 `ConsensusState` interface.
func (ConsensusState) ClientType() string {
	return ClientType
}
