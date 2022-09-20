package utils

import (
	"crypto/rand"
)

func NewEntropy(bitSize int) []byte {
	entropy := make([]byte, bitSize/8)
	_, _ = rand.Read(entropy)
	return entropy
}
