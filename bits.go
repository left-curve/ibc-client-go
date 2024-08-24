package grug

// Bits represents an array of bits, compacted into an array of bytes.
type BitArray []byte

// BitAtIndex returns the bit at the i-th position.
func (bits BitArray) BitAtIndex(i int) byte {
	quotient := i / 8
	remainder := i % 8

	return (bits[quotient] >> remainder) & 1
}
