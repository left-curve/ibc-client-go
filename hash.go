package grug

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

// Assert that `Hash` implements both the `json.Marshaler` and `json.Unmarshaler`
// interfaces.
var (
	_ json.Marshaler   = (*Hash)(nil)
	_ json.Unmarshaler = (*Hash)(nil)
)

// Hash represents a byte array of exactly 32 bytes.
type Hash [32]byte

// doSha256 performs the SHA2-256 hash.
func doSha256(preimage []byte) Hash {
	hasher := sha256.New()
	hasher.Write(preimage)
	return Hash(hasher.Sum(nil))
}

// MarshalJSON implements the `json.Marshaler` interface.
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(h[:]))
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	if len(decoded) != 32 {
		return ErrIncorrectHashLength
	}

	copy(h[:], decoded)
	return nil
}
