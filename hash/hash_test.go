// Package hash benchmarks various hashing algorithms.
// Especially with hashing algorithms, faster is not always better.
// One should always decide on the basis of the respective requirements.
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/sha3"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"math/rand"
	"testing"
)

func benchmarkHashAlgo(b *testing.B, h hash.Hash) {
	data := make([]byte, 2048)
	rand.Read(data)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		h.Write(data)
		h.Sum(nil)
	}
}

func BenchmarkAdler32(b *testing.B) {
	benchmarkHashAlgo(b, adler32.New())
}

func BenchmarkCRC32(b *testing.B) {
	benchmarkHashAlgo(b, crc32.NewIEEE())
}

func BenchmarkFnv128(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New128())
}

func BenchmarkMD5(b *testing.B) {
	benchmarkHashAlgo(b, md5.New())
}

func BenchmarkSHA1(b *testing.B) {
	benchmarkHashAlgo(b, sha1.New())
}

func BenchmarkSHA256(b *testing.B) {
	benchmarkHashAlgo(b, sha256.New())
}

func BenchmarkSHA512(b *testing.B) {
	benchmarkHashAlgo(b, sha512.New())
}

func BenchmarkSHA3256(b *testing.B) {
	benchmarkHashAlgo(b, sha3.New256())
}

func BenchmarkSHA3512(b *testing.B) {
	benchmarkHashAlgo(b, sha3.New512())
}

func BenchmarkWhirlpool(b *testing.B) {
	benchmarkHashAlgo(b, whirlpool.New())
}
