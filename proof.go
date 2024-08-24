package grug

// Proof is Grug's Merkle proof. It's an enum between either a membership or a
// non-membership proof.
type Proof struct {
	Membership    *MembershipProof    `json:"membership,omitempty"`
	NonMembership *NonMembershipProof `json:"non_membership,omitempty"`
}

// MembershipProof is the variant of `Proof` for a key that exists.
type MembershipProof struct {
	SiblingHashes []*Hash `json:"sibling_hashes"`
}

// NonMembershipProof is the variant of `Proof` for a key that doesn't exist.
type NonMembershipProof struct {
	Node          Node    `json:"proof_node"`
	SiblingHashes []*Hash `json:"sibling_hashes"`
}

// Validate performs basic validation on the Merkle proof.
func (proof Proof) Validate() error {
	// One and only one between `Membership` and `NonMembership` should be nil.
	if (proof.Membership == nil) == (proof.NonMembership == nil) {
		return ErrMalformedProof
	}

	// In case the proof is a non-membership proof, the node must be valid.
	if proof.NonMembership != nil {
		proof.NonMembership.Node.Validate()
	}

	return nil
}

// VerifyMembership verifies that a key-value pair exists in the Merkle tree of
// the given root hash.
func (proof Proof) VerifyMembership(rootHash, keyHash, valueHash Hash) error {
	// The proof must be a membership proof.
	if proof.Membership == nil {
		return ErrIncorrectProofType
	}

	return proof.Membership.VerifyMembership(rootHash, keyHash, valueHash)
}

// VerifyNonMembership verfies that a key doesn't exist in the Merkle tree of
// the given hash.
func (proof Proof) VerifyNonMembership(rootHash, keyHash Hash) error {
	// The proof must be a non-membership proof.
	if proof.NonMembership == nil {
		return ErrIncorrectProofType
	}

	return proof.NonMembership.VerifyNonMembership(rootHash, keyHash)
}

// VerifyMembership verifies that a key-value pair exists in the Merkle tree of
// the given root hash.
func (memberProof MembershipProof) VerifyMembership(rootHash, keyHash, valueHash Hash) error {
	bits := BitArray(keyHash[:])
	currentHash := LeafNode{KeyHash: keyHash, ValueHash: valueHash}.Hash()

	return computeAndCompareRootHash(rootHash, bits, memberProof.SiblingHashes, currentHash)
}

// VerifyNonMembership verfies that a key doesn't exist in the Merkle tree of
// the given hash.
func (nonMemberProof NonMembershipProof) VerifyNonMembership(rootHash, keyHash Hash) error {
	bits := BitArray(keyHash[:])

	var currentHash Hash
	if nonMemberProof.Node.Internal != nil {
		// If the node is an internal node, we check the bit at the depth:
		// - if the bit is 0, the node must not have a left child;
		// - if the bit is 1, the node must not have a right child.
		if bits.BitAtIndex(len(nonMemberProof.SiblingHashes)) == 0 {
			if nonMemberProof.Node.Internal.LeftHash != nil {
				return ErrUnexpectedChild
			}
		} else {
			if nonMemberProof.Node.Internal.RightHash != nil {
				return ErrUnexpectedChild
			}
		}

		currentHash = nonMemberProof.Node.Internal.Hash()
	} else {
		// If the node is a leaf, it's bit path must share a common prefix with the
		// key we want to prove not exist.
		existBits := BitArray(nonMemberProof.Node.Leaf.KeyHash[:])
		for i := 0; i < len(nonMemberProof.SiblingHashes); i++ {
			if existBits.BitAtIndex(i) != bits.BitAtIndex(i) {
				return ErrNotCommonPrefix
			}
		}

		currentHash = nonMemberProof.Node.Leaf.Hash()
	}

	return computeAndCompareRootHash(rootHash, bits, nonMemberProof.SiblingHashes, currentHash)
}

func computeAndCompareRootHash(rootHash Hash, bits BitArray, siblingHashes []*Hash, currentHash Hash) error {
	for depth := len(siblingHashes) - 1; depth >= 0; depth-- {
		bit := bits.BitAtIndex(depth)
		siblingHash := siblingHashes[depth]

		if bit == 0 {
			currentHash = InternalNode{LeftHash: &currentHash, RightHash: siblingHash}.Hash()
		} else {
			currentHash = InternalNode{LeftHash: siblingHash, RightHash: &currentHash}.Hash()
		}
	}

	if currentHash != rootHash {
		return ErrRootHashMismatch
	}

	return nil
}
