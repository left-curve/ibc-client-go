package grug

import errorsmod "cosmossdk.io/errors"

var (
	ErrIncorrectHashLength = errorsmod.Register(ClientType, 2, "incorrect hash length")
	ErrMalformedProof      = errorsmod.Register(ClientType, 3, "malformed proof")
	ErrIncorrectProofType  = errorsmod.Register(ClientType, 4, "incorrect proof type")
	ErrUnexpectedChild     = errorsmod.Register(ClientType, 5, "invalid non-membership proof: expecting node to not have a child, but it has one")
	ErrNotCommonPrefix     = errorsmod.Register(ClientType, 6, "invalid non-membership proof: node does't have a common bit prefix with the key")
	ErrRootHashMismatch    = errorsmod.Register(ClientType, 7, "computed root hash doesn't match the actual")
)
