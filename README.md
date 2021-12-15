# Go Benchmarks

[![Go Report Card](https://goreportcard.com/badge/github.com/SimonWaldherr/golang-benchmarks)](https://goreportcard.com/report/github.com/SimonWaldherr/golang-benchmarks)  
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  

In programming in general, and in Golang in particular, many roads lead to Rome.
From time to time I ask myself which of these ways is the fastest. 
In Golang there is a wonderful solution, with `go test -bench` you can measure the speed very easily and quickly.
In order for you to benefit from it too, I will publish such benchmarks in this repository in the future.

## ToC

* [base64](https://github.com/SimonWaldherr/golang-benchmarks#base64)
* [between](https://github.com/SimonWaldherr/golang-benchmarks#between)
* [concat](https://github.com/SimonWaldherr/golang-benchmarks#concat)
* [contains](https://github.com/SimonWaldherr/golang-benchmarks#contains)
* [foreach](https://github.com/SimonWaldherr/golang-benchmarks#foreach)
* [hash](https://github.com/SimonWaldherr/golang-benchmarks#hash)
* [index](https://github.com/SimonWaldherr/golang-benchmarks#index)
* [json](https://github.com/SimonWaldherr/golang-benchmarks#json)
* [math](https://github.com/SimonWaldherr/golang-benchmarks#math)
* [parse](https://github.com/SimonWaldherr/golang-benchmarks#parse)
* [random](https://github.com/SimonWaldherr/golang-benchmarks#random)
* [regexp](https://github.com/SimonWaldherr/golang-benchmarks#regexp)

## Golang?

I published another repository where I show some Golang examples.
If you\'re interested in new programming languages, you should definitely take a look at Golang:

* [Golang examples](https://github.com/SimonWaldherr/golang-examples)
* [tour.golang.org](https://tour.golang.org/)
* [Go by example](https://gobyexample.com/)
* [Golang Book](http://www.golang-book.com/)
* [Go-Learn](https://github.com/skippednote/Go-Learn)

## Is it any good?

[Yes](https://news.ycombinator.com/item?id=3067434)

## Benchmark Results

Golang Version: go version go1.18beta1 darwin/amd64

### base64

```go
// Package base64 benchmarks some base64 functions.
// On all tested systems it's faster to decode a
// base64 encoded string instead of a check via regex.
package base64

import (
	"encoding/base64"
	"regexp"
	"testing"
)

func base64decode(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func base64regex(s string) bool {
	matched, _ := regexp.MatchString(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$`, s)
	return matched
}

func BenchmarkBase64decode(b *testing.B) {
	isNotBase64 := `Invalid string`
	isBase64 := `VmFsaWQgc3RyaW5nCg==`

	for n := 0; n < b.N; n++ {
		base64decode(isNotBase64)
		base64decode(isBase64)
	}
}

func BenchmarkBase64regex(b *testing.B) {
	isNotBase64 := `Invalid string`
	isBase64 := `VmFsaWQgc3RyaW5nCg==`

	for n := 0; n < b.N; n++ {
		base64regex(isNotBase64)
		base64regex(isBase64)
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkBase64decode-8   	 8433190	       121.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkBase64regex-8    	   54193	     20814 ns/op	   21382 B/op	     198 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/base64	3.648s
```

### between

```go
// Package between compares the performance of checking
// if a number is between two other numbers via regex
// and by parsing the number as integers.
package between

import (
	"regexp"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/ranger"
	"testing"
)

func BenchmarkNumberRegEx(b *testing.B) {
	re := ranger.Compile(89, 1001)
	re = "^(" + re + ")$"
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		matched, err := regexp.MatchString(re, "404")
		if !matched || err != nil {
			b.Log("Error in Benchmark")
		}

		matched, err = regexp.MatchString(re, "2000")
		if matched || err != nil {
			b.Log("Error in Benchmark")
		}
	}
}

func BenchmarkFulltextRegEx(b *testing.B) {
	re := ranger.Compile(89, 1001)
	re = " (" + re + ") "
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		matched, err := regexp.MatchString(re, "lorem ipsum 404 dolor sit")
		if !matched || err != nil {
			b.Log("Error in Benchmark")
		}

		matched, err = regexp.MatchString(re, "lorem ipsum 2000 dolor sit")
		if matched || err != nil {
			b.Log("Error in Benchmark")
		}
	}
}

func BenchmarkNumberParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i1 := as.Int("404")
		i2 := as.Int("2000")

		if i1 < 89 || i1 > 1001 {
			b.Log("Error in Benchmark")
		}

		if !(i2 < 89 || i2 > 1001) {
			b.Log("Error in Benchmark")
		}
	}
}

func BenchmarkFulltextParse(b *testing.B) {
	re := regexp.MustCompile("[0-9]+")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		i1 := as.Int(re.FindString("lorem ipsum 404 dolor sit"))
		i2 := as.Int(re.FindString("lorem ipsum 2000 dolor sit"))

		if i1 < 89 || i1 > 1001 {
			b.Log("Error in Benchmark")
		}

		if !(i2 < 89 || i2 > 1001) {
			b.Log("Error in Benchmark")
		}
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkNumberRegEx-8     	   66166	     15380 ns/op	   16102 B/op	     142 allocs/op
BenchmarkFulltextRegEx-8   	   88393	     12268 ns/op	   11598 B/op	     104 allocs/op
BenchmarkNumberParse-8     	20069037	        61.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkFulltextParse-8   	 1336287	       913.5 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/between	6.021s
```

### concat

```go
// Package concat benchmarks the performance of
// various string concatenation methods.
// Instead of just concatenating a string to another string
// it is also possible (and much faster) to use
// a bytes buffer.
package concat

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkConcatString(b *testing.B) {
	var str string
	for n := 0; n < b.N; n++ {
		str += "x"
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")

	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	var builder strings.Builder
	for n := 0; n < b.N; n++ {
		builder.WriteString("x")
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkConcatString-8    	 1000000	     58661 ns/op	  503994 B/op	       1 allocs/op
BenchmarkConcatBuffer-8    	200141751	         5.710 ns/op	       2 B/op	       0 allocs/op
BenchmarkConcatBuilder-8   	517038657	         3.592 ns/op	       5 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/concat	62.792s
```

### contains

```go
// Package contains tests various ways of checking
// if a string is contained in another string.
package contains

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

// strings.Contains
func contains() bool {
	return strings.Contains("Lorem Ipsum", "em Ip")
}

func containsNot() bool {
	return strings.Contains("Lorem Ipsum", "Dolor")
}

func TestContains(t *testing.T) {
	if contains() == false {
		t.Error("ERROR: contains")
	}
	if containsNot() == true {
		t.Error("ERROR: contains not")
	}
}

func BenchmarkContains(b *testing.B) {
	for n := 0; n < b.N; n++ {
		contains()
	}
}

func BenchmarkContainsNot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsNot()
	}
}

// bytes.Contains
func containsBytes() bool {
	return bytes.Contains([]byte("Lorem Ipsum"), []byte("em Ip"))
}

func containsBytesNot() bool {
	return bytes.Contains([]byte("Lorem Ipsum"), []byte("Dolor"))
}

func TestContainsBytes(t *testing.T) {
	if containsBytes() == false {
		t.Error("ERROR: bytes contains")
	}
	if containsBytesNot() == true {
		t.Error("ERROR: bytes contains not")
	}
}

func BenchmarkContainsBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsBytes()
	}
}

func BenchmarkContainsBytesNot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsBytesNot()
	}
}

// regexp.MustCompile + regexp.MatchString
func compileMatch(re *regexp.Regexp) bool {
	matched := re.MatchString("Lorem Ipsum")
	return matched
}

func compileMatchNot(re *regexp.Regexp) bool {
	matched := re.MatchString("Lorem Ipsum")
	return matched
}

func TestCompileMatch(t *testing.T) {
	re1 := regexp.MustCompile("em Ip")
	re2 := regexp.MustCompile("Dolor")
	if compileMatch(re1) == false {
		t.Error("ERROR: compile match")
	}
	if compileMatchNot(re2) == true {
		t.Error("ERROR: compile match not")
	}
}

func BenchmarkCompileMatch(b *testing.B) {
	re := regexp.MustCompile("em Ip")
	for n := 0; n < b.N; n++ {
		compileMatch(re)
	}
}

func BenchmarkCompileMatchNot(b *testing.B) {
	re := regexp.MustCompile("Dolor")
	for n := 0; n < b.N; n++ {
		compileMatchNot(re)
	}
}

// regexp.MatchString
func match() bool {
	matched, _ := regexp.MatchString("em Ip", "Lorem Ipsum")
	return matched
}

func matchNot() bool {
	matched, _ := regexp.MatchString("Dolor", "Lorem Ipsum")
	return matched
}

func TestMatch(t *testing.T) {
	if match() == false {
		t.Error("ERROR: match")
	}
	if matchNot() == true {
		t.Error("ERROR: match not")
	}
}

func BenchmarkMatch(b *testing.B) {
	for n := 0; n < b.N; n++ {
		match()
	}
}

func BenchmarkMatchNot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		matchNot()
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkContains-8           	126274000	         9.490 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsNot-8        	130826602	         8.798 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytes-8      	100000000	        11.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytesNot-8   	127874654	         9.298 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatch-8       	13490476	        88.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatchNot-8    	27309816	        45.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatch-8              	  727933	      1472 ns/op	    1363 B/op	      17 allocs/op
BenchmarkMatchNot-8           	  819375	      1397 ns/op	    1362 B/op	      17 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/contains	12.569s
```

### foreach

```go
// Package foreach benchmarks ranging over slices and maps.
package foreach

import (
	"testing"
)

var amap map[int]string
var aslice []string

func init() {
	amap = map[int]string{
		0: "lorem",
		1: "ipsum",
		2: "dolor",
		3: "sit",
		4: "amet",
	}

	aslice = []string{
		"lorem",
		"ipsum",
		"dolor",
		"sit",
		"amet",
	}
}

func forMap() {
	for i := 0; i < len(amap); i++ {
		_ = amap[i]
	}
}

func rangeMap() {
	for _, v := range amap {
		_ = v
	}
}

func rangeSlice() {
	for _, v := range aslice {
		_ = v
	}
}

func rangeSliceKey() {
	for k := range aslice {
		_ = aslice[k]
	}
}

func BenchmarkForMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		forMap()
	}
}

func BenchmarkRangeMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rangeMap()
	}
}

func BenchmarkRangeSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rangeSlice()
	}
}

func BenchmarkRangeSliceKey(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rangeSliceKey()
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkForMap-8          	46463888	        23.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeMap-8        	17785903	        66.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSlice-8      	392926434	         3.082 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSliceKey-8   	384636010	         3.073 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/foreach	5.567s
```

### hash

```go
// Package hash benchmarks various hashing algorithms.
// Especially with hashing algorithms, faster is not always better.
// One should always decide on the basis of the respective requirements.
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"math/rand"
	"testing"

	"github.com/jzelinskie/whirlpool"
	"github.com/reusee/mmh3"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
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

func BenchmarkBlake2b256(b *testing.B) {
	h, err := blake2b.New256(nil)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkHashAlgo(b, h)
}

func BenchmarkBlake2b512(b *testing.B) {
	h, err := blake2b.New512(nil)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkHashAlgo(b, h)
}

func BenchmarkBlake3256(b *testing.B) {
	benchmarkHashAlgo(b, blake3.New())
}

func BenchmarkMMH3(b *testing.B) {
	benchmarkHashAlgo(b, mmh3.New128())
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
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkAdler32-8      	 1354020	       822.2 ns/op	       8 B/op	       1 allocs/op
BenchmarkBlake2b256-8   	  527420	      2313 ns/op	      32 B/op	       1 allocs/op
BenchmarkBlake2b512-8   	  519346	      2281 ns/op	      64 B/op	       1 allocs/op
BenchmarkBlake3256-8    	  414732	      3161 ns/op	      64 B/op	       2 allocs/op
BenchmarkMMH3-8         	 2679139	       447.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkCRC32-8        	 8325501	       142.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv128-8       	  237967	      4991 ns/op	      16 B/op	       1 allocs/op
BenchmarkMD5-8          	  383752	      3095 ns/op	      16 B/op	       1 allocs/op
BenchmarkSHA1-8         	  513009	      2596 ns/op	      24 B/op	       1 allocs/op
BenchmarkSHA256-8       	  220134	      5392 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA512-8       	  316483	      3805 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA3256-8      	  178647	      6817 ns/op	     512 B/op	       3 allocs/op
BenchmarkSHA3512-8      	   95163	     12360 ns/op	     576 B/op	       3 allocs/op
BenchmarkWhirlpool-8    	   23505	     50425 ns/op	      64 B/op	       1 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/hash	19.579s
```

### index

```go
// Package index benchmarks access on maps with various data types as keys.
package index

import (
	"math/rand"
	"strconv"
	"testing"
)

var NumItems int = 1000000

var ms map[string]string
var ks []string

var mi map[int]string
var ki []int

func initMapStringIndex() {
	ms = make(map[string]string)
	ks = make([]string, 0)

	for i := 0; i < NumItems; i++ {
		key := strconv.Itoa(rand.Intn(NumItems))
		ms[key] = "value" + strconv.Itoa(i)
		ks = append(ks, key)
	}
}

func initMapIntIndex() {
	mi = make(map[int]string)
	ki = make([]int, 0)

	for i := 0; i < NumItems; i++ {
		key := rand.Intn(NumItems)
		mi[key] = "value" + strconv.Itoa(i)
		ki = append(ki, key)
	}
}

func init() {
	initMapStringIndex()
	initMapIntIndex()
}

func BenchmarkMapStringKeys(b *testing.B) {
	i := 0

	for n := 0; n < b.N; n++ {
		if _, ok := ms[ks[i]]; ok {
		}

		i++
		if i >= NumItems {
			i = 0
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	i := 0

	for n := 0; n < b.N; n++ {
		if _, ok := mi[ki[i]]; ok {
		}

		i++
		if i >= NumItems {
			i = 0
		}
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkMapStringKeys-8   	10225266	       121.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntKeys-8      	15159615	        80.97 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/index	4.303s
```

### json

```go
package json

import (
	"encoding/json"
	"math"
	"math/big"
	"testing"
	"time"
)

type Data struct {
	String   string
	Time     time.Time
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Boolean  bool
	Float32  float32
	Float64  float64
	BigInt   big.Int
	BigFloat big.Float
}

func BenchmarkJsonMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var d = Data{
			String:   "",
			Time:     time.Now(),
			Int:      math.MaxInt32,
			Int8:     math.MaxInt8,
			Int16:    math.MaxInt16,
			Int32:    math.MaxInt32,
			Int64:    math.MaxInt64,
			Boolean:  false,
			Float32:  math.MaxFloat32,
			Float64:  math.MaxFloat64,
			BigInt:   *big.NewInt(math.MaxInt64),
			BigFloat: *big.NewFloat(math.MaxFloat64),
		}

		_, err := json.Marshal(d)
		if err != nil {
			b.Error(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	str := `
{
  "String": "",
  "Time": "2019-10-30T16:41:29.853426+07:00",
  "Int": 2147483647,
  "Int8": 127,
  "Int16": 32767,
  "Int32": 2147483647,
  "Int64": 9223372036854775807,
  "Boolean": false,
  "Float32": 3.4028235e+38,
  "Float64": 1.7976931348623157e+308,
  "BigInt": 9999999999999999999,
  "BigFloat": "2.7976931348623157e+308"
}
`

	for n := 0; n < b.N; n++ {
		var d Data
		err := json.Unmarshal([]byte(str), &d)
		if err != nil {
			b.Error(err)
			b.Fail()
			return
		}
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkJsonMarshal-8     	  600850	      1726 ns/op	     480 B/op	       5 allocs/op
BenchmarkJsonUnmarshal-8   	  174606	      7467 ns/op	    2136 B/op	      39 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/json	2.599s
```

### math

```go
// Package math compares the speed of various mathematical operations.
package math

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkMathInt8(b *testing.B) {
	var intVal int8
	for n := 0; n < b.N; n++ {
		intVal = intVal + 2
	}
}

func BenchmarkMathInt32(b *testing.B) {
	var intVal int32
	for n := 0; n < b.N; n++ {
		intVal = intVal + 2
	}
}

func BenchmarkMathInt64(b *testing.B) {
	var intVal int64
	for n := 0; n < b.N; n++ {
		intVal = intVal + 2
	}
}

func BenchmarkMathAtomicInt32(b *testing.B) {
	var intVal int32
	for n := 0; n < b.N; n++ {
		atomic.AddInt32(&intVal, 2)
	}
}

func BenchmarkMathAtomicInt64(b *testing.B) {
	var intVal int64
	for n := 0; n < b.N; n++ {
		atomic.AddInt64(&intVal, 2)
	}
}

type IntMutex struct {
	v   int64
	mux sync.Mutex
}

func BenchmarkMathMutexInt(b *testing.B) {
	var m IntMutex
	for n := 0; n < b.N; n++ {
		m.mux.Lock()
		m.v = m.v + 2
		m.mux.Unlock()
	}
}

func BenchmarkMathFloat32(b *testing.B) {
	var floatVal float32
	for n := 0; n < b.N; n++ {
		floatVal = floatVal + 2
	}
}

func BenchmarkMathFloat64(b *testing.B) {
	var floatVal float64
	for n := 0; n < b.N; n++ {
		floatVal = floatVal + 2
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkMathInt8-8          	1000000000	         0.3538 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt32-8         	1000000000	         0.3365 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt64-8         	1000000000	         0.3147 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt32-8   	207475148	         5.865 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt64-8   	227856810	         5.302 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathMutexInt-8      	81229758	        14.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat32-8       	1000000000	         0.2867 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat64-8       	1000000000	         0.2930 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/math	6.696s
```

### parse

```go
// Package parse benchmarks parsing.
package parse

import (
	"strconv"
	"testing"
)

func BenchmarkParseBool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := strconv.ParseBool("true")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkParseInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := strconv.ParseInt("1337", 10, 64)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkParseFloat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := strconv.ParseFloat("3.141592653589793238462643383", 64)
		if err != nil {
			panic(err)
		}
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkParseBool-8    	1000000000	         0.4824 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt-8     	82117677	        14.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat-8   	14214274	        90.59 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/parse	3.279s
```

### random

```go
// Package random compares math/rand with crypto/rand.
// math/rand is much faster than crypto/rand, but it
// returns only a pseudo random number.
package random

import (
	crand "crypto/rand"
	"encoding/base64"
	"math/big"
	mrand "math/rand"
	"testing"
)

func BenchmarkMathRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mrand.Int63n(0xFFFF)
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := crand.Int(crand.Reader, big.NewInt(0xFFFF))
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCryptoRandString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := GenerateRandomString(32)
		if err != nil {
			panic(err)
		}
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := mrand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkMathRand-8           	45873254	        25.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkCryptoRand-8         	 1229484	       987.9 ns/op	      56 B/op	       4 allocs/op
BenchmarkCryptoRandString-8   	 7062478	       166.2 ns/op	     128 B/op	       3 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/random	4.904s
```

### regexp

```go
// Package regexp benchmarks the performance of a pre-compiled regexp match
// a non-pre-compiled match and JIT-cached-compilation via golibs: https://simonwaldherr.de/go/golibs
package regexp

import (
	"regexp"
	"testing"

	"simonwaldherr.de/go/golibs/regex"
)

var regexpStr string = `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,9}$`

func BenchmarkMatchString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := regexp.MatchString(regexpStr, "john.doe@example.tld")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMatchStringCompiled(b *testing.B) {
	r, err := regexp.Compile(regexpStr)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.MatchString("john.doe@example.tld")
	}
}

func BenchmarkMatchStringGolibs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := regex.MatchString("john.doe@example.tld", regexpStr)
		if err != nil {
			panic(err)
		}
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-6920HQ CPU @ 2.90GHz
BenchmarkMatchString-8           	  129441	      8488 ns/op	    9948 B/op	      86 allocs/op
BenchmarkMatchStringCompiled-8   	 2205708	       543.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchStringGolibs-8     	 2125749	       567.5 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	_/Users/simonwaldherr/git/golang-benchmarks/regexp	4.908s
```

