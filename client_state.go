package grug

import (
	"encoding/json"
	"time"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibcerrors "github.com/cosmos/ibc-go/v8/modules/core/errors"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ics23 "github.com/cosmos/ics23/go"
)

// ClientType indicates the client is a Grug chain.
//
// TODO: Grug doesn't have an ICS number, so putting "xx" for now.
const ClientType = "xx-grug"

var (
	// Assert that `ClientState` implements the ICS-02 `ClientState` interface.
	_ exported.ClientState = (*ClientState)(nil)

	// placeholderProofSpec is a placeholder for the ICS-23 proof spec for use in
	// the 07-tendermint client state.
	//
	// The Grug client doesn't actually make use of this. It's ignored in the
	// actual Merkle verification logic.
	placeholderProofSpecs = []*ics23.ProofSpec{}
)

// ClientState is the client state of Grug.
//
// This is a wrapper of the 07-tendermint client state, with with a few methods
// substituted by the Grug ones.
type ClientState struct{ *tendermint.ClientState }

// NewClientState creates a new Grug client state instance.
func NewClientState(
	chainID string, trustLevel tendermint.Fraction,
	trustingPeriod, unbondingPeriod, maxClockDrift time.Duration,
	latestHeight clienttypes.Height, upgradePath []string,
) *ClientState {
	return &ClientState{
		tendermint.NewClientState(
			chainID, trustLevel, trustingPeriod, unbondingPeriod, maxClockDrift,
			latestHeight, placeholderProofSpecs, upgradePath,
		),
	}
}

// ClientType implements the ICS-02 `ClientState` interface.
func (ClientState) ClientType() string {
	return ClientType
}

// VerifyMembership implements the ICS-02 `ClientState` interface.
func (cs ClientState) VerifyMembership(
	ctx sdk.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
	value []byte,
) error {
	if cs.GetLatestHeight().LT(height) {
		return errorsmod.Wrapf(
			ibcerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}

	if err := verifyDelayPeriodPassed(ctx, clientStore, height, delayTimePeriod, delayBlockPeriod); err != nil {
		return err
	}

	// Here's where we deviate from 07-tendermint: instead of unmarshalling the
	// proof into an ICS-23 proof in Protobuf, we unmarshal into the Grug proof in
	// JSON.
	var merkleProof Proof
	if err := json.Unmarshal(proof, &merkleProof); err != nil {
		return errorsmod.Wrap(
			commitmenttypes.ErrInvalidProof,
			"failed to unmarshal proof into Grug proof",
		)
	}

	// Also, the `path` is simply a `[]byte`, instead of `commitmenttypes.MerklePath`.
	merklePath, ok := path.(Path)
	if !ok {
		return errorsmod.Wrapf(
			ibcerrors.ErrInvalidType,
			"expected %T, got %T", Path{}, path,
		)
	}

	consensusState, found := tendermint.GetConsensusState(clientStore, cdc, height)
	if !found {
		return errorsmod.Wrap(
			clienttypes.ErrConsensusStateNotFound,
			"please ensure the proof was constructed against a height that exists on the client",
		)
	}

	return merkleProof.VerifyMembership(Hash(consensusState.GetRoot().GetHash()), doSha256(merklePath), doSha256(value))
}

// VerifyNonMembership implements the ICS-02 `ClientState` interface.
func (cs ClientState) VerifyNonMembership(
	ctx sdk.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
) error {
	if cs.GetLatestHeight().LT(height) {
		return errorsmod.Wrapf(
			ibcerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}

	if err := verifyDelayPeriodPassed(ctx, clientStore, height, delayTimePeriod, delayBlockPeriod); err != nil {
		return err
	}

	// Here's where we deviate from 07-tendermint: instead of unmarshalling the
	// proof into an ICS-23 proof in Protobuf, we unmarshal into the Grug proof in
	// JSON.
	var merkleProof Proof
	if err := json.Unmarshal(proof, &merkleProof); err != nil {
		return errorsmod.Wrap(
			commitmenttypes.ErrInvalidProof,
			"failed to unmarshal proof into Grug proof",
		)
	}

	// Also, the `path` is simply a `[]byte`, instead of `commitmenttypes.MerklePath`.
	merklePath, ok := path.(Path)
	if !ok {
		return errorsmod.Wrapf(
			ibcerrors.ErrInvalidType,
			"expected %T, got %T", Path{}, path,
		)
	}

	consensusState, found := tendermint.GetConsensusState(clientStore, cdc, height)
	if !found {
		return errorsmod.Wrap(
			clienttypes.ErrConsensusStateNotFound,
			"please ensure the proof was constructed against a height that exists on the client",
		)
	}

	return merkleProof.VerifyNonMembership(Hash(consensusState.GetRoot().GetHash()), doSha256(merklePath))
}

// verifyDelayPeriodPassed is copied without change from 07-tendermint.
//
// Unfortunately this function is private, so we have to copy the code instead
// of simply importing it.
//
// TODO: submit a PR to ibc-go to make this function public.
func verifyDelayPeriodPassed(ctx sdk.Context, store storetypes.KVStore, proofHeight exported.Height, delayTimePeriod, delayBlockPeriod uint64) error {
	if delayTimePeriod != 0 {
		// check that executing chain's timestamp has passed consensusState's processed time + delay time period
		processedTime, ok := tendermint.GetProcessedTime(store, proofHeight)
		if !ok {
			return errorsmod.Wrapf(tendermint.ErrProcessedTimeNotFound, "processed time not found for height: %s", proofHeight)
		}

		currentTimestamp := uint64(ctx.BlockTime().UnixNano())
		validTime := processedTime + delayTimePeriod

		// NOTE: delay time period is inclusive, so if currentTimestamp is validTime, then we return no error
		if currentTimestamp < validTime {
			return errorsmod.Wrapf(tendermint.ErrDelayPeriodNotPassed, "cannot verify packet until time: %d, current time: %d",
				validTime, currentTimestamp)
		}

	}

	if delayBlockPeriod != 0 {
		// check that executing chain's height has passed consensusState's processed height + delay block period
		processedHeight, ok := tendermint.GetProcessedHeight(store, proofHeight)
		if !ok {
			return errorsmod.Wrapf(tendermint.ErrProcessedHeightNotFound, "processed height not found for height: %s", proofHeight)
		}

		currentHeight := clienttypes.GetSelfHeight(ctx)
		validHeight := clienttypes.NewHeight(processedHeight.GetRevisionNumber(), processedHeight.GetRevisionHeight()+delayBlockPeriod)

		// NOTE: delay block period is inclusive, so if currentHeight is validHeight, then we return no error
		if currentHeight.LT(validHeight) {
			return errorsmod.Wrapf(tendermint.ErrDelayPeriodNotPassed, "cannot verify packet until height: %s, current height: %s",
				validHeight, currentHeight)
		}
	}

	return nil
}
