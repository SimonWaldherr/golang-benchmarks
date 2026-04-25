// Package random compares math/rand with crypto/rand.
// math/rand is much faster than crypto/rand, but it
// returns only a pseudo random number.
package random

import (
	crand "crypto/rand"
	"encoding/base64"
	"io"
	"math/big"
	mrand "math/rand"
	"testing"
)

var (
	randomIntResult    int64
	randomBigIntResult *big.Int
	randomBytesResult  []byte
	randomStringResult string
)

type mathRandReader struct{}

func (mathRandReader) Read(p []byte) (int, error) {
	return mrand.Read(p)
}

func BenchmarkMathRand(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		randomIntResult = mrand.Int63n(0xFFFF)
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	limit := big.NewInt(0xFFFF)
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		r, err := crand.Int(crand.Reader, limit)
		if err != nil {
			panic(err)
		}
		randomBigIntResult = r
	}
}

func BenchmarkMathRandBytes(b *testing.B) {
	reader := mathRandReader{}
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		r, err := GenerateRandomBytes(reader, 32)
		if err != nil {
			panic(err)
		}
		randomBytesResult = r
	}
}

func BenchmarkCryptoRandBytes(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r, err := GenerateRandomBytes(crand.Reader, 32)
		if err != nil {
			panic(err)
		}
		randomBytesResult = r
	}
}

func BenchmarkMathRandString(b *testing.B) {
	reader := mathRandReader{}
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		r, err := GenerateRandomString(reader, 32)
		if err != nil {
			panic(err)
		}
		randomStringResult = r
	}
}

func BenchmarkCryptoRandString(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r, err := GenerateRandomString(crand.Reader, 32)
		if err != nil {
			panic(err)
		}
		randomStringResult = r
	}
}

func GenerateRandomBytes(r io.Reader, n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GenerateRandomString(r io.Reader, s int) (string, error) {
	b, err := GenerateRandomBytes(r, s)
	return base64.URLEncoding.EncodeToString(b), err
}
