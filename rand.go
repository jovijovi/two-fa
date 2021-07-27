package twofa

import (
	"crypto/rand"
)

const (
	randSize = 32
)

func newRand() ([]byte, error) {
	b := make([]byte, randSize)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
