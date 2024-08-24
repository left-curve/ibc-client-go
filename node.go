package grug

import "crypto/sha256"

var (
	HashPrefixInternalNode byte     = 0
	HashPrefixLeafNode     byte     = 1
	ZeroHash               [32]byte = [32]byte{}
)

// Node is a node is Grug's Merkle tree. It's an enum between either an internal
// or a leaf node.
type Node struct {
	Internal *InternalNode `json:"internal,omitempty"`
	Leaf     *LeafNode     `json:"leaf,omitempty"`
}

// InternalNode is an internal node in Grug's Merkle tree.
type InternalNode struct {
	LeftHash  *Hash `json:"left_hash,omitempty"`
	RightHash *Hash `json:"right_hash,omitempty"`
}

// LeafNode is a leaf node in Grug's Merkle tree.
type LeafNode struct {
	KeyHash   Hash `json:"key_hash"`
	ValueHash Hash `json:"value_hash"`
}

// Valiate performs basic validation on the node.
func (node Node) Validate() error {
	// One and only one between `Internal` and `Leaf` should be nil.
	if (node.Internal == nil) == (node.Leaf == nil) {
		return ErrMalformedProof
	}

	return nil
}

// Hash produces the hash of the internal node.
func (internalNode InternalNode) Hash() Hash {
	hasher := sha256.New()
	hasher.Write([]byte{HashPrefixInternalNode})
	if internalNode.LeftHash != nil {
		hasher.Write(internalNode.LeftHash[:])
	} else {
		hasher.Write(ZeroHash[:])
	}
	if internalNode.RightHash != nil {
		hasher.Write(internalNode.RightHash[:])
	} else {
		hasher.Write(ZeroHash[:])
	}
	return Hash(hasher.Sum(nil))
}

// Hash produces the hash of the leaf node.
func (leafNode LeafNode) Hash() Hash {
	hasher := sha256.New()
	hasher.Write([]byte{HashPrefixLeafNode})
	hasher.Write(leafNode.KeyHash[:])
	hasher.Write(leafNode.ValueHash[:])
	return Hash(hasher.Sum(nil))
}
