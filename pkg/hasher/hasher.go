package hasher

import (
	"crypto/sha256"
	"fmt"
)

type PasswordHash interface {
	Hash(password string) string
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
