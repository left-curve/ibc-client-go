package grug

import "github.com/cosmos/ibc-go/v8/modules/core/exported"

// Assert that `Path` implements the `exported.Path` interface.
var _ exported.Path = (*Path)(nil)

// Path is the Merkle path for Grug proofs. It's basically just a byte array.
type Path struct {
	Bytes []byte
}

// Empty implements the `exported.Path` interface.
func (path Path) Empty() bool {
	return len(path.Bytes) == 0
}
