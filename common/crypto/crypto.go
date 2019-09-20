package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type Hash []byte

func (h Hash) String() string {
	return hex.EncodeToString(h)
}

func Sum(data []byte) Hash {
	h := sha256.Sum256(data)
	return h[:]
}

func RandBytes(size int) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}
