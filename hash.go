package twofa

import (
	"crypto/sha1"
	"hash"
)

type HashProvider func() hash.Hash

type IHashFunc interface {
	Hash(msg []byte) ([]byte, error)
}

type HashFunc struct {
	Provider HashProvider
}

func (h *HashFunc) Hash(msg []byte) ([]byte, error) {
	provider := h.Provider()
	if _, err := provider.Write(msg); err != nil {
		return nil, err
	}

	return provider.Sum(nil), nil
}

func DefaultHashFunc() IHashFunc {
	hashFunc := new(HashFunc)
	hashFunc.Provider = sha1.New
	return hashFunc
}
