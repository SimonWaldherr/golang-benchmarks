package hash_crypto

import (
	"crypto/sha256"
	"testing"

	"golang.org/x/crypto/blake2b"
)

var sample = []byte("the quick brown fox jumps over the lazy dog")

func BenchmarkSHA256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := sha256.Sum256(sample)
		_ = h
	}
}

func BenchmarkBLAKE2b(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = blake2b.Sum256(sample)
	}
}
