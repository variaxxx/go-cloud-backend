package auth_refresh_token

import (
	"crypto/sha256"
	"encoding/hex"
)

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
