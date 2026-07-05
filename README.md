# Go Benchmarks

[![test](https://github.com/SimonWaldherr/golang-benchmarks/actions/workflows/audit.yml/badge.svg?branch=master&event=push)](https://github.com/SimonWaldherr/golang-benchmarks/actions/workflows/audit.yml) 
[![DOI](https://zenodo.org/badge/154216722.svg)](https://zenodo.org/badge/latestdoi/154216722) 
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonWaldherr/golang-benchmarks)](https://goreportcard.com/report/github.com/SimonWaldherr/golang-benchmarks) 
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  

In programming in general, and in Golang in particular, many roads lead to Rome.
From time to time I ask myself which of these ways is the fastest. 
In Golang there is a wonderful solution, with `go test -bench` you can measure the speed very easily and quickly.
In order for you to benefit from it too, I will publish such benchmarks in this repository in the future.

## ToC

* [base64](https://github.com/SimonWaldherr/golang-benchmarks#base64)
* [between](https://github.com/SimonWaldherr/golang-benchmarks#between)
* [caseinsensitivecompare](https://github.com/SimonWaldherr/golang-benchmarks#caseinsensitivecompare)
* [concat](https://github.com/SimonWaldherr/golang-benchmarks#concat)
* [contains](https://github.com/SimonWaldherr/golang-benchmarks#contains)
* [concurrency_counter](https://github.com/SimonWaldherr/golang-benchmarks#concurrency_counter)
* [embed](https://github.com/SimonWaldherr/golang-benchmarks#embed)
* [floodfill](https://github.com/SimonWaldherr/golang-benchmarks#floodfill)
* [foreach](https://github.com/SimonWaldherr/golang-benchmarks#foreach)
* [hash](https://github.com/SimonWaldherr/golang-benchmarks#hash)
* [index](https://github.com/SimonWaldherr/golang-benchmarks#index)
* [json](https://github.com/SimonWaldherr/golang-benchmarks#json)
* [math](https://github.com/SimonWaldherr/golang-benchmarks#math)
* [parse](https://github.com/SimonWaldherr/golang-benchmarks#parse)
* [random](https://github.com/SimonWaldherr/golang-benchmarks#random)
* [regexp](https://github.com/SimonWaldherr/golang-benchmarks#regexp)
* [sql](https://github.com/SimonWaldherr/golang-benchmarks#sql)
* [template](https://github.com/SimonWaldherr/golang-benchmarks#template)
* [trim](https://github.com/SimonWaldherr/golang-benchmarks#trim)

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

Golang Version: [go version go1.26.4 darwin/arm64](https://tip.golang.org/doc/go1.26)  
Hardware Spec: [Apple MacBook Pro 16-Inch M2 Max 2023](https://support.apple.com/kb/SP890) [(?)](https://everymac.com/systems/apple/macbook_pro/specs/macbook-pro-m2-max-12-core-cpu-30-core-gpu-16-2023-specs.html) [(buy)](https://amzn.to/3K80lP4)  

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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/base64
cpu: Apple M2 Max
BenchmarkBase64decode-12    	22767247	        52.63 ns/op	      32 B/op	       2 allocs/op
BenchmarkBase64regex-12     	  134428	      9086 ns/op	   21875 B/op	     198 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/base64	2.859s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/between
cpu: Apple M2 Max
BenchmarkNumberRegEx-12      	  184062	      6191 ns/op	   16849 B/op	     142 allocs/op
BenchmarkFulltextRegEx-12    	  240684	      5190 ns/op	   12068 B/op	     104 allocs/op
BenchmarkNumberParse-12      	31407274	        38.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkFulltextParse-12    	 2468192	       480.1 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/between	5.728s
```

### caseinsensitivecompare

```go
package caseinsensitivecompare

import (
	"strings"
	"testing"
)

var caseInsensitiveCompareResult bool

func BenchmarkEqualFold(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.EqualFold("abc", "ABC")
		caseInsensitiveCompareResult = strings.EqualFold("ABC", "ABC")
		caseInsensitiveCompareResult = strings.EqualFold("1aBcD", "1AbCd")
	}
}

func BenchmarkToUpper(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.ToUpper("abc") == strings.ToUpper("ABC")
		caseInsensitiveCompareResult = strings.ToUpper("ABC") == strings.ToUpper("ABC")
		caseInsensitiveCompareResult = strings.ToUpper("1aBcD") == strings.ToUpper("1AbCd")
	}
}

func BenchmarkToLower(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.ToLower("abc") == strings.ToLower("ABC")
		caseInsensitiveCompareResult = strings.ToLower("ABC") == strings.ToLower("ABC")
		caseInsensitiveCompareResult = strings.ToLower("1aBcD") == strings.ToLower("1AbCd")
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/caseinsensitivecompare
cpu: Apple M2 Max
BenchmarkEqualFold-12    	81751046	        14.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkToUpper-12      	10225863	       117.2 ns/op	      24 B/op	       3 allocs/op
BenchmarkToLower-12      	 8540114	       141.1 ns/op	      40 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/caseinsensitivecompare	4.833s
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

const concatParts = 64

var concatResult string

func BenchmarkConcatString(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < concatParts; i++ {
			str += "x"
		}
		concatResult = str
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		for i := 0; i < concatParts; i++ {
			buffer.WriteString("x")
		}
		concatResult = buffer.String()
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var builder strings.Builder
		for i := 0; i < concatParts; i++ {
			builder.WriteString("x")
		}
		concatResult = builder.String()
	}
}

func BenchmarkConcatBuilderGrow(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var builder strings.Builder
		builder.Grow(concatParts)
		for i := 0; i < concatParts; i++ {
			builder.WriteString("x")
		}
		concatResult = builder.String()
	}
}

func BenchmarkConcat(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var str string
			for i := 0; i < concatParts; i++ {
				str += "x"
			}
			concatResult = str
		}
	})

	b.Run("Buffer", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var buffer bytes.Buffer
			for i := 0; i < concatParts; i++ {
				buffer.WriteString("x")
			}
			concatResult = buffer.String()
		}
	})

	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var builder strings.Builder
			for i := 0; i < concatParts; i++ {
				builder.WriteString("x")
			}
			concatResult = builder.String()
		}
	})

	b.Run("BuilderGrow", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var builder strings.Builder
			builder.Grow(concatParts)
			for i := 0; i < concatParts; i++ {
				builder.WriteString("x")
			}
			concatResult = builder.String()
		}
	})
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/concat
cpu: Apple M2 Max
BenchmarkConcatString-12         	  872364	      1336 ns/op	    2400 B/op	      63 allocs/op
BenchmarkConcatBuffer-12         	 4814239	       249.2 ns/op	     128 B/op	       2 allocs/op
BenchmarkConcatBuilder-12        	 7447760	       159.9 ns/op	     120 B/op	       4 allocs/op
BenchmarkConcatBuilderGrow-12    	 9211669	       129.8 ns/op	      64 B/op	       1 allocs/op
BenchmarkConcat/String-12        	  888828	      1386 ns/op	    2400 B/op	      63 allocs/op
BenchmarkConcat/Buffer-12        	 4823826	       250.3 ns/op	     128 B/op	       2 allocs/op
BenchmarkConcat/Builder-12       	 7349016	       164.6 ns/op	     120 B/op	       4 allocs/op
BenchmarkConcat/BuilderGrow-12   	 8999282	       131.3 ns/op	      64 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/concat	11.016s
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

// BenchmarkContainsMethods benchmarks different methods to check substring presence.
func BenchmarkContainsMethods(b *testing.B) {
	b.Run("Strings.Contains", func(b *testing.B) {
		str := "Lorem Ipsum"
		substr := "em Ip"
		for n := 0; n < b.N; n++ {
			_ = strings.Contains(str, substr)
			_ = strings.Contains(str, "Dolor")
		}
	})

	b.Run("Bytes.Contains", func(b *testing.B) {
		str := []byte("Lorem Ipsum")
		substr := []byte("em Ip")
		for n := 0; n < b.N; n++ {
			_ = bytes.Contains(str, substr)
			_ = bytes.Contains(str, []byte("Dolor"))
		}
	})

	b.Run("RegexMatchString", func(b *testing.B) {
		re := regexp.MustCompile(`em Ip`)
		for n := 0; n < b.N; n++ {
			_ = re.MatchString("Lorem Ipsum")
			_ = re.MatchString("Dolor")
		}
	})

	b.Run("RegexMatch", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = regexp.MatchString(`em Ip`, "Lorem Ipsum")
			_, _ = regexp.MatchString(`em Ip`, "Dolor")
		}
	})
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/contains
cpu: Apple M2 Max
BenchmarkContains-12            	207514189	         5.730 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsNot-12         	120990061	         9.995 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytes-12       	191096216	         6.421 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytesNot-12    	100000000	        10.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatch-12        	26216328	        46.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatchNot-12     	42367768	        27.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatch-12               	 1839638	       682.7 ns/op	    1396 B/op	      17 allocs/op
BenchmarkMatchNot-12            	 1842693	       658.3 ns/op	    1397 B/op	      17 allocs/op
BenchmarkContainsMethods/Strings.Contains-12         	69026180	        17.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/Bytes.Contains-12           	67496623	        17.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatchString-12         	17844970	        69.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatch-12               	  856844	      1344 ns/op	    2795 B/op	      34 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/contains	18.382s
```

### concurrency_counter

```go
package concurrency_counter

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// Benchmarks comparing common concurrency counter patterns.
// Variants:
// - sync.Mutex (write)
// - sync.RWMutex (write + read-heavy)
// - atomic.AddInt64
// - channels (various buffer sizes)
// - worker-pool (fixed workers consuming jobs)

func BenchmarkMutexParallel(b *testing.B) {
	var mu sync.Mutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkRWMutexWriteParallel(b *testing.B) {
	var mu sync.RWMutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkRWMutexReadParallel(b *testing.B) {
	var mu sync.RWMutex
	var counter int64

	// populate counter
	counter = 42

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = counter
			mu.RUnlock()
		}
	})
}

func BenchmarkAtomicParallel(b *testing.B) {
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&counter, 1)
		}
	})
}

func BenchmarkChannelBufferedSizes(b *testing.B) {
	sizes := []int{0, 1, 16, 1024}

	for _, sz := range sizes {
		b.Run(fmt.Sprintf("buf=%d", sz), func(b *testing.B) {
			ch := make(chan struct{}, sz)
			var wg sync.WaitGroup
			var counter int64

			wg.Add(1)
			go func() {
				for range ch {
					atomic.AddInt64(&counter, 1)
				}
				wg.Done()
			}()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					ch <- struct{}{}
				}
			})

			close(ch)
			wg.Wait()
		})
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	workers := runtime.GOMAXPROCS(0)

	b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
		jobs := make(chan struct{}, 1024)
		var wg sync.WaitGroup
		var counter int64

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				for range jobs {
					atomic.AddInt64(&counter, 1)
				}
				wg.Done()
			}()
		}

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				jobs <- struct{}{}
			}
		})

		close(jobs)
		wg.Wait()
	})
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/concurrency_counter
cpu: Apple M2 Max
BenchmarkMutexParallel-12           	 9617564	       128.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutexWriteParallel-12    	12490365	        93.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutexReadParallel-12     	10118828	       116.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAtomicParallel-12          	19465969	        61.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=0-12         	 3474572	       349.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=1-12         	 4407970	       270.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=16-12        	 8073873	       170.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=1024-12      	10008231	       124.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool/workers=12-12              	20238248	        67.43 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/concurrency_counter	12.872s
```

### embed

```go
package embed

import (
	_ "embed"
	"io/ioutil"
	"os"
	"testing"
)

//go:embed example.txt
var embeddedFile []byte

func BenchmarkEmbed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Access the embedded file
		_ = embeddedFile
	}
}

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Read the file from disk
		data, err := os.ReadFile("example.txt")
		if err != nil {
			b.Fatalf("failed to read file: %v", err)
		}
		_ = data
	}
}

func BenchmarkIoutilReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Read the file using ioutil
		data, err := ioutil.ReadFile("example.txt")
		if err != nil {
			b.Fatalf("failed to read file: %v", err)
		}
		_ = data
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/embed
cpu: Apple M2 Max
BenchmarkEmbed-12             	1000000000	         0.3051 ns/op	       0 B/op	       0 allocs/op
BenchmarkReadFile-12          	  110800	     11140 ns/op	     840 B/op	       5 allocs/op
BenchmarkIoutilReadFile-12    	  107013	     11134 ns/op	     840 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/embed	3.297s
```

### floodfill

```go
// Package floodfill benchmarks various flood fill implementations.
package main

import (
	"testing"
)

// Rekursive Implementierung
func floodFillRecursive(image [][]int, sr int, sc int, newColor int) [][]int {
	oldColor := image[sr][sc]
	if oldColor == newColor {
		return image
	}
	floodFillRecurse(image, sr, sc, oldColor, newColor)
	return image
}

func floodFillRecurse(image [][]int, sr int, sc int, oldColor int, newColor int) {
	if sr < 0 || sr >= len(image) || sc < 0 || sc >= len(image[0]) || image[sr][sc] != oldColor {
		return
	}
	image[sr][sc] = newColor

	floodFillRecurse(image, sr+1, sc, oldColor, newColor)
	floodFillRecurse(image, sr-1, sc, oldColor, newColor)
	floodFillRecurse(image, sr, sc+1, oldColor, newColor)
	floodFillRecurse(image, sr, sc-1, oldColor, newColor)
}

// Iterative Implementierung mit Stack (DFS)
func floodFillDFS(image [][]int, sr int, sc int, newColor int) [][]int {
	oldColor := image[sr][sc]
	if oldColor == newColor {
		return image
	}

	stack := [][]int{{sr, sc}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		r, c := current[0], current[1]
		if r < 0 || r >= len(image) || c < 0 || c >= len(image[0]) || image[r][c] != oldColor {
			continue
		}

		image[r][c] = newColor

		stack = append(stack, []int{r + 1, c})
		stack = append(stack, []int{r - 1, c})
		stack = append(stack, []int{r, c + 1})
		stack = append(stack, []int{r, c - 1})
	}
	return image
}

// Iterative Implementierung mit Queue (BFS)
func floodFillBFS(image [][]int, sr int, sc int, newColor int) [][]int {
	oldColor := image[sr][sc]
	if oldColor == newColor {
		return image
	}

	queue := [][]int{{sr, sc}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		r, c := current[0], current[1]
		if r < 0 || r >= len(image) || c < 0 || c >= len(image[0]) || image[r][c] != oldColor {
			continue
		}

		image[r][c] = newColor

		queue = append(queue, []int{r + 1, c})
		queue = append(queue, []int{r - 1, c})
		queue = append(queue, []int{r, c + 1})
		queue = append(queue, []int{r, c - 1})
	}
	return image
}

// Iterative Implementierung mit Stack (4-Wege-Verbindung)
func floodFillStack4Way(image [][]int, sr int, sc int, newColor int) [][]int {
	oldColor := image[sr][sc]
	if oldColor == newColor {
		return image
	}

	stack := [][]int{{sr, sc}}
	directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		r, c := current[0], current[1]
		if r < 0 || r >= len(image) || c < 0 || c >= len(image[0]) || image[r][c] != oldColor {
			continue
		}

		image[r][c] = newColor

		for _, dir := range directions {
			stack = append(stack, []int{r + dir[0], c + dir[1]})
		}
	}
	return image
}

// Komplexes Beispielbild
var complexImage = [][]int{
	{0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 1, 0},
	{0, 1, 0, 0, 1, 0},
	{0, 1, 1, 1, 1, 0},
	{0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 1, 1},
}

// Benchmark für die rekursive Implementierung
func BenchmarkFloodFillRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		image := make([][]int, len(complexImage))
		for j := range complexImage {
			image[j] = make([]int, len(complexImage[j]))
			copy(image[j], complexImage[j])
		}
		floodFillRecursive(image, 1, 1, 2)
	}
}

// Benchmark für die iterative Implementierung mit Stack (DFS)
func BenchmarkFloodFillDFS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		image := make([][]int, len(complexImage))
		for j := range complexImage {
			image[j] = make([]int, len(complexImage[j]))
			copy(image[j], complexImage[j])
		}
		floodFillDFS(image, 1, 1, 2)
	}
}

// Benchmark für die iterative Implementierung mit Queue (BFS)
func BenchmarkFloodFillBFS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		image := make([][]int, len(complexImage))
		for j := range complexImage {
			image[j] = make([]int, len(complexImage[j]))
			copy(image[j], complexImage[j])
		}
		floodFillBFS(image, 1, 1, 2)
	}
}

// Benchmark für die iterative Implementierung mit Stack (4-Wege-Verbindung)
func BenchmarkFloodFillStack4Way(b *testing.B) {
	for i := 0; i < b.N; i++ {
		image := make([][]int, len(complexImage))
		for j := range complexImage {
			image[j] = make([]int, len(complexImage[j]))
			copy(image[j], complexImage[j])
		}
		floodFillStack4Way(image, 1, 1, 2)
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/floodfill
cpu: Apple M2 Max
BenchmarkFloodFillRecursive-12    	 5952084	       194.4 ns/op	     432 B/op	       7 allocs/op
BenchmarkFloodFillDFS-12          	 1405957	       924.8 ns/op	    1744 B/op	      48 allocs/op
BenchmarkFloodFillBFS-12          	 1000000	      1089 ns/op	    2704 B/op	      53 allocs/op
BenchmarkFloodFillStack4Way-12    	 1260918	       883.8 ns/op	    1744 B/op	      48 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/floodfill	7.030s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/foreach
cpu: Apple M2 Max
BenchmarkForMap-12           	61557022	        18.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeMap-12         	27512103	        44.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSlice-12       	438105837	         2.756 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSliceKey-12    	428354636	         2.768 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/foreach	6.414s
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
	"hash/crc64"
	"hash/fnv"
	"math/rand"
	"testing"

	"github.com/jzelinskie/whirlpool"
	"github.com/reusee/mmh3"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

func benchmarkHashAlgo(b *testing.B, h hash.Hash) {
	data := make([]byte, 2048)
	rand.Read(data)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		h.Reset()
		h.Write(data)
		_ = h.Sum(nil)
	}
}

func benchmarkBCryptHashAlgo(b *testing.B, cost int) {
	data := make([]byte, 2048)
	rand.Read(data)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bcrypt.GenerateFromPassword(data, cost)
	}
}

func BenchmarkAdler32(b *testing.B) {
	benchmarkHashAlgo(b, adler32.New())
}

func BenchmarkBCryptCost4(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 4)
}

func BenchmarkBCryptCost10(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 10)
}

func BenchmarkBCryptCost16(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 16)
}

/*
func BenchmarkBCryptCost22(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 22)
}

func BenchmarkBCryptCost28(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 28)
}

func BenchmarkBCryptCost31(b *testing.B) {
	benchmarkBCryptHashAlgo(b, 31)
}
*/

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

func BenchmarkCRC64ISO(b *testing.B) {
	benchmarkHashAlgo(b, crc64.New(crc64.MakeTable(crc64.ISO)))
}

func BenchmarkCRC64ECMA(b *testing.B) {
	benchmarkHashAlgo(b, crc64.New(crc64.MakeTable(crc64.ECMA)))
}

func BenchmarkFnv32(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New32())
}

func BenchmarkFnv32a(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New32a())
}

func BenchmarkFnv64(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New64())
}

func BenchmarkFnv64a(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New64a())
}

func BenchmarkFnv128(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New128())
}

func BenchmarkFnv128a(b *testing.B) {
	benchmarkHashAlgo(b, fnv.New128a())
}

func BenchmarkMD4(b *testing.B) {
	benchmarkHashAlgo(b, md4.New())
}

func BenchmarkMD5(b *testing.B) {
	benchmarkHashAlgo(b, md5.New())
}

func BenchmarkSHA1(b *testing.B) {
	benchmarkHashAlgo(b, sha1.New())
}

func BenchmarkSHA224(b *testing.B) {
	benchmarkHashAlgo(b, sha256.New224())
}

func BenchmarkSHA256(b *testing.B) {
	benchmarkHashAlgo(b, sha256.New())
}

func BenchmarkSHA384(b *testing.B) {
	benchmarkHashAlgo(b, sha512.New384())
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

func BenchmarkRIPEMD160(b *testing.B) {
	benchmarkHashAlgo(b, ripemd160.New())
}

func BenchmarkWhirlpool(b *testing.B) {
	benchmarkHashAlgo(b, whirlpool.New())
}

func BenchmarkSHA256Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		data := make([]byte, 2048)
		rand.Read(data)
		for pb.Next() {
			h := sha256.New()
			h.Write(data)
			h.Sum(nil)
		}
	})
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/hash
cpu: Apple M2 Max
BenchmarkAdler32-12           	 1750003	       685.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkBCryptCost4-12       	610751516	         1.975 ns/op	       0 B/op	       0 allocs/op
BenchmarkBCryptCost10-12      	606473218	         1.997 ns/op	       0 B/op	       0 allocs/op
BenchmarkBCryptCost16-12      	606074126	         1.988 ns/op	       0 B/op	       0 allocs/op
BenchmarkBlake2b256-12        	  459516	      2626 ns/op	      32 B/op	       1 allocs/op
BenchmarkBlake2b512-12        	  464535	      2548 ns/op	      64 B/op	       1 allocs/op
BenchmarkBlake3256-12         	  411738	      2953 ns/op	      32 B/op	       1 allocs/op
BenchmarkMMH3-12              	 3586282	       333.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkCRC32-12             	 5003482	       243.4 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ISO-12          	 1000000	      1153 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ECMA-12         	 1000000	      1160 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32-12             	  498943	      2438 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32a-12            	  498172	      2465 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64-12             	  495234	      2455 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64a-12            	  493644	      2469 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv128-12            	  185625	      6515 ns/op	      24 B/op	       2 allocs/op
BenchmarkFnv128a-12           	  161979	      7609 ns/op	      24 B/op	       2 allocs/op
BenchmarkMD4-12               	  396938	      2950 ns/op	      16 B/op	       1 allocs/op
BenchmarkMD5-12               	  390572	      3073 ns/op	      16 B/op	       1 allocs/op
BenchmarkSHA1-12              	 1401793	       854.4 ns/op	      24 B/op	       1 allocs/op
BenchmarkSHA224-12            	 1356103	       871.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA256-12            	 1369776	       866.4 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA384-12            	  776698	      1523 ns/op	      48 B/op	       1 allocs/op
BenchmarkSHA512-12            	  788835	      1499 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA3256-12           	  474246	      2541 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA3512-12           	  263305	      4578 ns/op	      64 B/op	       1 allocs/op
BenchmarkRIPEMD160-12         	  199186	      6223 ns/op	      24 B/op	       1 allocs/op
BenchmarkWhirlpool-12         	   47346	     25365 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA256Parallel-12    	12143214	       106.7 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/hash	40.551s
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

func BenchmarkMapStringIndex(b *testing.B) {
	i := 0

	for n := 0; n < b.N; n++ {
		_ = ms[ks[i]]

		i++
		if i >= NumItems {
			i = 0
		}
	}
}

func BenchmarkMapIntIndex(b *testing.B) {
	i := 0

	for n := 0; n < b.N; n++ {
		_ = mi[ki[i]]

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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/index
cpu: Apple M2 Max
BenchmarkMapStringKeys-12     	21227931	        47.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntKeys-12        	70636360	        22.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapStringIndex-12    	22286070	        48.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntIndex-12       	54642250	        22.10 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/index	7.888s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/json
cpu: Apple M2 Max
BenchmarkJsonMarshal-12      	 1683763	       718.5 ns/op	     480 B/op	       5 allocs/op
BenchmarkJsonUnmarshal-12    	  330643	      3680 ns/op	    1816 B/op	      27 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/json	3.277s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/math
cpu: Apple M2 Max
BenchmarkMathInt8-12           	1000000000	         0.3055 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt32-12          	1000000000	         0.3098 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt64-12          	1000000000	         0.3050 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt32-12    	288814456	         4.183 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt64-12    	271699471	         4.161 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathMutexInt-12       	147039955	         8.163 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat32-12        	1000000000	         0.3008 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat64-12        	1000000000	         0.3009 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/math	7.288s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/parse
cpu: Apple M2 Max
BenchmarkParseBool-12     	553994132	         2.157 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt-12      	100000000	        10.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat-12    	18276117	        63.55 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/parse	4.082s
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
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/random
cpu: Apple M2 Max
BenchmarkMathRand-12            	168725125	         7.316 ns/op	       0 B/op	       0 allocs/op
BenchmarkCryptoRand-12          	10485511	       113.9 ns/op	      48 B/op	       3 allocs/op
BenchmarkMathRandBytes-12       	18394468	        66.23 ns/op	      32 B/op	       1 allocs/op
BenchmarkCryptoRandBytes-12     	 4607526	       260.0 ns/op	      32 B/op	       1 allocs/op
BenchmarkMathRandString-12      	10100998	       123.7 ns/op	     128 B/op	       3 allocs/op
BenchmarkCryptoRandString-12    	 3905157	       301.1 ns/op	     128 B/op	       3 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/random	8.915s
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/regexp
cpu: Apple M2 Max
BenchmarkMatchString-12            	  293149	      4311 ns/op	   10217 B/op	      86 allocs/op
BenchmarkMatchStringCompiled-12    	 3783412	       321.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchStringGolibs-12      	 3521288	       325.9 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/regexp	4.640s
```

### sql

```go
package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

var sqlResultSink int64

func openSQLite(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn)
}

// prepareSchema creates the users and orders tables used by benchmarks.
func prepareSchema(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE users (id INT PRIMARY KEY, name TEXT, email TEXT, active BOOL)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE orders (id INT PRIMARY KEY, user_id INT, amount FLOAT, status TEXT, meta JSON)`); err != nil {
		return err
	}
	return nil
}

func benchmarkDB(b *testing.B, open func(string) (*sql.DB, error), dsn string) *sql.DB {
	db, err := open(dsn)
	if err != nil {
		b.Fatalf("open: %v", err)
	}

	if err := prepareSchema(db); err != nil {
		db.Close()
		b.Fatalf("schema: %v", err)
	}
	b.Cleanup(func() {
		if err := db.Close(); err != nil {
			b.Fatalf("close: %v", err)
		}
	})

	return db
}

func prepareBenchmarkStmt(b *testing.B, tx interface {
	Prepare(string) (*sql.Stmt, error)
}, query string) *sql.Stmt {
	stmt, err := tx.Prepare(query)
	if err != nil {
		b.Fatalf("prepare: %v", err)
	}
	b.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			b.Fatalf("close stmt: %v", err)
		}
	})

	return stmt
}

func populateUsersAndOrders(b *testing.B, db *sql.DB, users, ordersPerUser int) {
	tx, err := db.Begin()
	if err != nil {
		b.Fatalf("begin populate: %v", err)
	}

	ustmt, err := tx.Prepare(`INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		b.Fatalf("prepare users: %v", err)
	}
	ostmt, err := tx.Prepare(`INSERT INTO orders (id, user_id, amount, status, meta) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		ustmt.Close()
		tx.Rollback()
		b.Fatalf("prepare orders: %v", err)
	}

	for i := 1; i <= users; i++ {
		if _, err := ustmt.Exec(i, fmt.Sprintf("user%d", i), fmt.Sprintf("u%d@example.com", i), i%3 != 0); err != nil {
			ustmt.Close()
			ostmt.Close()
			tx.Rollback()
			b.Fatalf("populate user: %v", err)
		}
		for j := 0; j < ordersPerUser; j++ {
			id := i*100 + j
			status := "PAID"
			if j%4 == 0 {
				status = "PENDING"
			}
			if _, err := ostmt.Exec(id, i, float64(id)*1.5, status, `{"device":"web"}`); err != nil {
				ustmt.Close()
				ostmt.Close()
				tx.Rollback()
				b.Fatalf("populate order: %v", err)
			}
		}
	}

	if err := ustmt.Close(); err != nil {
		ostmt.Close()
		tx.Rollback()
		b.Fatalf("close users stmt: %v", err)
	}
	if err := ostmt.Close(); err != nil {
		tx.Rollback()
		b.Fatalf("close orders stmt: %v", err)
	}
	if err := tx.Commit(); err != nil {
		b.Fatalf("commit populate: %v", err)
	}
}

func benchInsertTxPerRow(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
	db := benchmarkDB(b, open, dsn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := db.Begin()
		if err != nil {
			b.Fatalf("begin: %v", err)
		}
		stmt, err := tx.Prepare(`INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
		if err != nil {
			b.Fatalf("prepare: %v", err)
		}
		if _, err := stmt.Exec(i, fmt.Sprintf("user%d", i), fmt.Sprintf("u%d@example.com", i), i%2 == 0); err != nil {
			b.Fatalf("exec: %v", err)
		}
		stmt.Close()
		if err := tx.Commit(); err != nil {
			b.Fatalf("commit: %v", err)
		}
	}
}

func benchInsertPrepared(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
	db := benchmarkDB(b, open, dsn)
	tx, err := db.Begin()
	if err != nil {
		b.Fatalf("begin: %v", err)
	}
	stmt, err := tx.Prepare(`INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		b.Fatalf("prepare: %v", err)
	}

	const valuePoolSize = 1024
	names := make([]string, valuePoolSize)
	emails := make([]string, valuePoolSize)
	for i := range names {
		names[i] = fmt.Sprintf("user%d", i)
		emails[i] = fmt.Sprintf("u%d@example.com", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valueIndex := i % valuePoolSize
		if _, err := stmt.Exec(i, names[valueIndex], emails[valueIndex], i%2 == 0); err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()

	if err := stmt.Close(); err != nil {
		tx.Rollback()
		b.Fatalf("close stmt: %v", err)
	}
	if err := tx.Commit(); err != nil {
		b.Fatalf("commit: %v", err)
	}
}

func benchSelectPoint(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 2)
	stmt := prepareBenchmarkStmt(b, db, `SELECT email, active FROM users WHERE id = ?`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var email string
		var active bool
		id := i%users + 1
		if err := stmt.QueryRow(id).Scan(&email, &active); err != nil {
			b.Fatalf("query row: %v", err)
		}
		if active {
			sqlResultSink += int64(len(email))
		}
	}
}

func benchSelectRange(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 2)
	stmt := prepareBenchmarkStmt(b, db, `SELECT id, name FROM users WHERE id BETWEEN ? AND ? ORDER BY id`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := i%users + 1
		end := start + 24
		if end > users {
			end = users
		}
		rows, err := stmt.Query(start, end)
		if err != nil {
			b.Fatalf("query: %v", err)
		}

		var count int64
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				rows.Close()
				b.Fatalf("scan: %v", err)
			}
			count += int64(id + len(name))
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatalf("rows: %v", err)
		}
		if err := rows.Close(); err != nil {
			b.Fatalf("close rows: %v", err)
		}
		sqlResultSink += count
	}
}

func benchUpdatePrepared(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 1)
	stmt := prepareBenchmarkStmt(b, db, `UPDATE users SET active = ? WHERE id = ?`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i%users + 1
		result, err := stmt.Exec(i%2 == 0, id)
		if err != nil {
			b.Fatalf("exec update: %v", err)
		}
		affected, err := result.RowsAffected()
		if err == nil {
			sqlResultSink += affected
		}
	}
}

func benchDeletePrepared(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 3)
	stmt := prepareBenchmarkStmt(b, db, `DELETE FROM orders WHERE id = ?`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := i%users + 1
		orderID := userID*100 + i%3
		result, err := stmt.Exec(orderID)
		if err != nil {
			b.Fatalf("exec delete: %v", err)
		}
		affected, err := result.RowsAffected()
		if err == nil {
			sqlResultSink += affected
		}
	}
}

func benchAggregateByStatus(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 6)
	stmt := prepareBenchmarkStmt(b, db, `
		SELECT status, COUNT(id), SUM(amount)
		FROM orders
		GROUP BY status
		ORDER BY status
	`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query()
		if err != nil {
			b.Fatalf("query aggregate: %v", err)
		}

		var total int64
		for rows.Next() {
			var status string
			var count int
			var amount sql.NullFloat64
			if err := rows.Scan(&status, &count, &amount); err != nil {
				rows.Close()
				b.Fatalf("scan aggregate: %v", err)
			}
			total += int64(len(status) + count)
			if amount.Valid {
				total += int64(amount.Float64)
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatalf("rows aggregate: %v", err)
		}
		if err := rows.Close(); err != nil {
			b.Fatalf("close rows aggregate: %v", err)
		}
		sqlResultSink += total
	}
}

func benchSelectOrderedOrders(b *testing.B, open func(string) (*sql.DB, error), dsn string, users int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, users, 5)
	stmt := prepareBenchmarkStmt(b, db, `
		SELECT id, user_id, amount, status
		FROM orders
		WHERE status = ?
		ORDER BY amount
	`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query("PAID")
		if err != nil {
			b.Fatalf("query ordered orders: %v", err)
		}

		var total int64
		for rows.Next() {
			var id int
			var userID int
			var amount float64
			var status string
			if err := rows.Scan(&id, &userID, &amount, &status); err != nil {
				rows.Close()
				b.Fatalf("scan ordered orders: %v", err)
			}
			total += int64(id + userID + len(status))
			total += int64(amount)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatalf("rows ordered orders: %v", err)
		}
		if err := rows.Close(); err != nil {
			b.Fatalf("close rows ordered orders: %v", err)
		}
		sqlResultSink += total
	}
}

func benchSelectJoin(b *testing.B, open func(string) (*sql.DB, error), dsn string, rows int) {
	db := benchmarkDB(b, open, dsn)
	populateUsersAndOrders(b, db, rows, 2)

	query := strings.TrimSpace(`
		SELECT u.name AS user, SUM(o.amount) AS total, COUNT(o.id) AS cnt
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id AND o.status = 'PAID'
		GROUP BY u.name
		ORDER BY u.name
	`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(query)
		if err != nil {
			b.Fatalf("query: %v", err)
		}

		var totalCount int64
		for rows.Next() {
			var name string
			var total sql.NullFloat64
			var cnt int
			if err := rows.Scan(&name, &total, &cnt); err != nil {
				rows.Close()
				b.Fatalf("scan: %v", err)
			}
			totalCount += int64(len(name) + cnt)
			if total.Valid {
				totalCount += int64(total.Float64)
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatalf("rows: %v", err)
		}
		if err := rows.Close(); err != nil {
			b.Fatalf("close rows: %v", err)
		}
		sqlResultSink += totalCount
	}
}

func benchTransactionBatch(b *testing.B, open func(string) (*sql.DB, error), dsn string, batchSize int) {
	db := benchmarkDB(b, open, dsn)

	names := make([]string, batchSize)
	emails := make([]string, batchSize)
	for i := 0; i < batchSize; i++ {
		names[i] = fmt.Sprintf("user%d", i)
		emails[i] = fmt.Sprintf("u%d@example.com", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, err := db.Begin()
		if err != nil {
			b.Fatalf("begin: %v", err)
		}
		stmt, err := tx.Prepare(`INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
		if err != nil {
			tx.Rollback()
			b.Fatalf("prepare: %v", err)
		}
		baseID := i * batchSize
		for j := 0; j < batchSize; j++ {
			if _, err := stmt.Exec(baseID+j, names[j], emails[j], j%2 == 0); err != nil {
				stmt.Close()
				tx.Rollback()
				b.Fatalf("exec: %v", err)
			}
		}
		if err := stmt.Close(); err != nil {
			tx.Rollback()
			b.Fatalf("close stmt: %v", err)
		}
		if err := tx.Commit(); err != nil {
			b.Fatalf("commit: %v", err)
		}
	}
}

func runSQLBenchmarks(b *testing.B, open func(string) (*sql.DB, error), dsn func(string) string) {
	b.Run("InsertTxPerRow", func(b *testing.B) {
		benchInsertTxPerRow(b, open, dsn("insert_tx"))
	})
	b.Run("InsertPrepared", func(b *testing.B) {
		benchInsertPrepared(b, open, dsn("insert_prepared"))
	})
	b.Run("TransactionBatch100", func(b *testing.B) {
		benchTransactionBatch(b, open, dsn("batch_100"), 100)
	})
	b.Run("SelectPoint", func(b *testing.B) {
		benchSelectPoint(b, open, dsn("select_point"), 1000)
	})
	b.Run("SelectRange25", func(b *testing.B) {
		benchSelectRange(b, open, dsn("select_range"), 1000)
	})
	b.Run("UpdatePrepared", func(b *testing.B) {
		benchUpdatePrepared(b, open, dsn("update_prepared"), 1000)
	})
	b.Run("DeletePrepared", func(b *testing.B) {
		benchDeletePrepared(b, open, dsn("delete_prepared"), 1000)
	})
	b.Run("AggregateByStatus", func(b *testing.B) {
		benchAggregateByStatus(b, open, dsn("aggregate_status"), 1000)
	})
	b.Run("SelectOrderedOrders", func(b *testing.B) {
		benchSelectOrderedOrders(b, open, dsn("select_ordered_orders"), 500)
	})
	b.Run("SelectJoin", func(b *testing.B) {
		benchSelectJoin(b, open, dsn("select_join"), 200)
	})
}

func sqliteDSN(name string) string {
	return "file:" + name + "?mode=memory&cache=shared"
}

func Benchmark_Insert_SQLite(b *testing.B) {
	benchInsertTxPerRow(b, openSQLite, sqliteDSN("insert_legacy"))
}

func Benchmark_SelectJoin_SQLite(b *testing.B) {
	benchSelectJoin(b, openSQLite, sqliteDSN("select_join_legacy"), 200)
}
package main

import (
	"database/sql"
	"testing"

	tsqldriver "github.com/SimonWaldherr/tinySQL/driver"
)

func Benchmark_Insert_TinySQL(b *testing.B) {
	benchInsertTxPerRow(b, openTinySQL, "mem://?tenant=bench_insert")
}

func Benchmark_SelectJoin_TinySQL(b *testing.B) {
	benchSelectJoin(b, openTinySQL, "mem://?tenant=bench_select", 200)
}

func Benchmark_SQLCompare(b *testing.B) {
	runComparedSQLBenchmark(b, "InsertTxPerRow", benchInsertTxPerRow)
	runComparedSQLBenchmark(b, "InsertPrepared", benchInsertPrepared)
	runComparedSQLBenchmark(b, "TransactionBatch100", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchTransactionBatch(b, open, dsn, 100)
	})
	runComparedSQLBenchmark(b, "SelectPoint", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchSelectPoint(b, open, dsn, 1000)
	})
	runComparedSQLBenchmark(b, "SelectRange25", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchSelectRange(b, open, dsn, 1000)
	})
	runComparedSQLBenchmark(b, "UpdatePrepared", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchUpdatePrepared(b, open, dsn, 1000)
	})
	runComparedSQLBenchmark(b, "DeletePrepared", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchDeletePrepared(b, open, dsn, 1000)
	})
	runComparedSQLBenchmark(b, "AggregateByStatus", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchAggregateByStatus(b, open, dsn, 1000)
	})
	runComparedSQLBenchmark(b, "SelectOrderedOrders", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchSelectOrderedOrders(b, open, dsn, 500)
	})
	runComparedSQLBenchmark(b, "SelectJoin", func(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
		benchSelectJoin(b, open, dsn, 200)
	})
}

func runComparedSQLBenchmark(
	b *testing.B,
	name string,
	bench func(*testing.B, func(string) (*sql.DB, error), string),
) {
	b.Run(name, func(b *testing.B) {
		b.Run("SQLite", func(b *testing.B) {
			bench(b, openSQLite, sqliteDSN("compare_"+name))
		})
		b.Run("TinySQL", func(b *testing.B) {
			bench(b, openTinySQL, "mem://?tenant=compare_"+name)
		})
	})
}

func openTinySQL(dsn string) (*sql.DB, error) {
	return sql.Open(tsqldriver.DriverName, dsn)
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/sql
cpu: Apple M2 Max
Benchmark_Insert_SQLite-12         	  162249	      9078 ns/op	     993 B/op	      27 allocs/op
Benchmark_SelectJoin_SQLite-12     	    2610	    454634 ns/op	   14928 B/op	    1215 allocs/op
Benchmark_Insert_TinySQL-12        	   10000	    675432 ns/op	 1104907 B/op	   14886 allocs/op
Benchmark_SelectJoin_TinySQL-12    	      22	  45697669 ns/op	99678560 B/op	  324824 allocs/op
Benchmark_SQLCompare/InsertTxPerRow/SQLite-12  	  138482	      8293 ns/op	     991 B/op	      26 allocs/op
Benchmark_SQLCompare/InsertTxPerRow/TinySQL-12 	   10000	    562790 ns/op	 1105589 B/op	   14887 allocs/op
Benchmark_SQLCompare/InsertPrepared/SQLite-12  	  693192	      1826 ns/op	     304 B/op	      10 allocs/op
Benchmark_SQLCompare/InsertPrepared/TinySQL-12 	  307714	      3806 ns/op	    2337 B/op	      74 allocs/op
Benchmark_SQLCompare/TransactionBatch100/SQLite-12         	    5790	    200223 ns/op	   31031 B/op	    1111 allocs/op
Benchmark_SQLCompare/TransactionBatch100/TinySQL-12        	    1155	   7893682 ns/op	12418972 B/op	  180570 allocs/op
Benchmark_SQLCompare/SelectPoint/SQLite-12                 	  445681	      2249 ns/op	     635 B/op	      23 allocs/op
Benchmark_SQLCompare/SelectPoint/TinySQL-12                	  120090	     10460 ns/op	    3503 B/op	      66 allocs/op
Benchmark_SQLCompare/SelectRange25/SQLite-12               	   81038	     13597 ns/op	    2113 B/op	     159 allocs/op
Benchmark_SQLCompare/SelectRange25/TinySQL-12              	   48229	     23466 ns/op	   17896 B/op	     245 allocs/op
Benchmark_SQLCompare/UpdatePrepared/SQLite-12              	  653145	      1577 ns/op	     162 B/op	       6 allocs/op
Benchmark_SQLCompare/UpdatePrepared/TinySQL-12             	    8412	    140919 ns/op	    2589 B/op	      59 allocs/op
Benchmark_SQLCompare/DeletePrepared/SQLite-12              	 1338363	       898.5 ns/op	     135 B/op	       6 allocs/op
Benchmark_SQLCompare/DeletePrepared/TinySQL-12             	    8547	    138182 ns/op	   15994 B/op	      50 allocs/op
Benchmark_SQLCompare/AggregateByStatus/SQLite-12           	     422	   2764831 ns/op	     744 B/op	      32 allocs/op
Benchmark_SQLCompare/AggregateByStatus/TinySQL-12          	     271	   4352741 ns/op	 4541232 B/op	   32140 allocs/op
Benchmark_SQLCompare/SelectOrderedOrders/SQLite-12         	     772	   1537046 ns/op	  138520 B/op	   13954 allocs/op
Benchmark_SQLCompare/SelectOrderedOrders/TinySQL-12        	    1194	    843736 ns/op	  887911 B/op	   16017 allocs/op
Benchmark_SQLCompare/SelectJoin/SQLite-12                  	    2436	    453353 ns/op	   14928 B/op	    1215 allocs/op
Benchmark_SQLCompare/SelectJoin/TinySQL-12                 	      14	  85874595 ns/op	99678608 B/op	  324826 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/sql	942.457s
```

### template

```go
// Package template benchmarks the performance of different templating methods
package template

import (
	"bytes"
	htmltemplate "html/template"
	"regexp"
	"testing"
	texttemplate "text/template"
)

// Define a struct to hold the data for the templates
type Data struct {
	Name    string
	Address string
}

// Prepare the templates and data
var (
	data = Data{Name: "John Doe", Address: "123 Elm St"}

	textTplString = "Name: {{.Name}}, Address: {{.Address}}"
	htmlTplString = "<div>Name: {{.Name}}</div><div>Address: {{.Address}}</div>"
	regExpString  = "Name: {{NAME}}, Address: {{ADDRESS}}"
)

// Benchmark for text/template
func BenchmarkTextTemplate(b *testing.B) {
	tpl, _ := texttemplate.New("text").Parse(textTplString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		tpl.Execute(&buf, data)
	}
}

// Benchmark for html/template
func BenchmarkHTMLTemplate(b *testing.B) {
	tpl, _ := htmltemplate.New("html").Parse(htmlTplString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		tpl.Execute(&buf, data)
	}
}

// Benchmark for replacing placeholders using regexp
func BenchmarkRegExp(b *testing.B) {
	rName := regexp.MustCompile(`{{NAME}}`)
	rAddress := regexp.MustCompile(`{{ADDRESS}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := rName.ReplaceAllString(regExpString, data.Name)
		result = rAddress.ReplaceAllString(result, data.Address)
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/template
cpu: Apple M2 Max
BenchmarkTextTemplate-12    	 3342513	       330.1 ns/op	     272 B/op	       5 allocs/op
BenchmarkHTMLTemplate-12    	 1215890	       996.8 ns/op	     496 B/op	      15 allocs/op
BenchmarkRegExp-12          	 3088400	       388.4 ns/op	     297 B/op	       9 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/template	5.660s
```

