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

Golang Version: [go version go1.26.2 darwin/arm64](https://tip.golang.org/doc/go1.26)  
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
BenchmarkBase64decode-12    	22126984	        53.87 ns/op	      32 B/op	       2 allocs/op
BenchmarkBase64regex-12     	  127070	      9335 ns/op	   21887 B/op	     198 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/base64	3.783s
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
BenchmarkNumberRegEx-12      	  166154	      6475 ns/op	   16850 B/op	     142 allocs/op
BenchmarkFulltextRegEx-12    	  215836	      5202 ns/op	   12072 B/op	     104 allocs/op
BenchmarkNumberParse-12      	31613541	        37.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkFulltextParse-12    	 2502544	       484.3 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/between	5.488s
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
BenchmarkEqualFold-12    	66977056	        16.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkToUpper-12      	10253871	       119.9 ns/op	      24 B/op	       3 allocs/op
BenchmarkToLower-12      	 8311576	       144.2 ns/op	      40 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/caseinsensitivecompare	4.027s
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
BenchmarkConcatString-12         	  849830	      1357 ns/op	    2400 B/op	      63 allocs/op
BenchmarkConcatBuffer-12         	 4819561	       252.1 ns/op	     128 B/op	       2 allocs/op
BenchmarkConcatBuilder-12        	 7375850	       162.5 ns/op	     120 B/op	       4 allocs/op
BenchmarkConcatBuilderGrow-12    	 9110570	       129.7 ns/op	      64 B/op	       1 allocs/op
BenchmarkConcat/String-12        	  883183	      1350 ns/op	    2400 B/op	      63 allocs/op
BenchmarkConcat/Buffer-12        	 4850577	       270.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkConcat/Builder-12       	 7153880	       171.8 ns/op	     120 B/op	       4 allocs/op
BenchmarkConcat/BuilderGrow-12   	 8977294	       133.5 ns/op	      64 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/concat	11.011s
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
BenchmarkContains-12            	241628170	         4.904 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsNot-12         	209303389	         5.720 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytes-12       	219988339	         5.504 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytesNot-12    	186626149	         6.828 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatch-12        	25491687	        46.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatchNot-12     	48798174	        25.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatch-12               	 1672022	       720.1 ns/op	    1397 B/op	      17 allocs/op
BenchmarkMatchNot-12            	 1777761	       672.1 ns/op	    1397 B/op	      17 allocs/op
BenchmarkContainsMethods/Strings.Contains-12         	100000000	        10.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/Bytes.Contains-12           	100000000	        12.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatchString-12         	17138757	        67.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatch-12               	  836452	      1381 ns/op	    2796 B/op	      34 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/contains	19.378s
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
BenchmarkMutexParallel-12           	 9550502	       120.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutexWriteParallel-12    	12136197	        99.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkRWMutexReadParallel-12     	12322478	       110.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAtomicParallel-12          	19968507	        55.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=0-12         	 3619378	       334.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=1-12         	 4133910	       291.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=16-12        	 6987964	       169.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkChannelBufferedSizes/buf=1024-12      	 9380547	       129.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool/workers=12-12              	17303614	        65.79 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/concurrency_counter	13.349s
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
BenchmarkEmbed-12             	1000000000	         0.3050 ns/op	       0 B/op	       0 allocs/op
BenchmarkReadFile-12          	  108997	     11371 ns/op	     840 B/op	       5 allocs/op
BenchmarkIoutilReadFile-12    	  108512	     11071 ns/op	     840 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/embed	3.204s
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
BenchmarkFloodFillRecursive-12    	 5497546	       201.2 ns/op	     432 B/op	       7 allocs/op
BenchmarkFloodFillDFS-12          	 1388280	       851.0 ns/op	    1744 B/op	      48 allocs/op
BenchmarkFloodFillBFS-12          	 1000000	      1075 ns/op	    2704 B/op	      53 allocs/op
BenchmarkFloodFillStack4Way-12    	 1286047	       912.8 ns/op	    1744 B/op	      48 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/floodfill	6.834s
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
BenchmarkForMap-12           	55885108	        19.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeMap-12         	28007853	        43.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSlice-12       	441927368	         2.706 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSliceKey-12    	444138961	         2.715 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/foreach	5.492s
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
BenchmarkAdler32-12           	 1726119	       683.7 ns/op	       8 B/op	       1 allocs/op
BenchmarkBCryptCost4-12       	    1225	    985500 ns/op	    8166 B/op	       9 allocs/op
BenchmarkBCryptCost10-12      	      19	  60300417 ns/op	    8257 B/op	       9 allocs/op
BenchmarkBCryptCost16-12      	       1	3922344291 ns/op	    9976 B/op	      11 allocs/op
BenchmarkBlake2b256-12        	  451138	      2677 ns/op	      32 B/op	       1 allocs/op
BenchmarkBlake2b512-12        	  459603	      2635 ns/op	      64 B/op	       1 allocs/op
BenchmarkBlake3256-12         	  408350	      3005 ns/op	      32 B/op	       1 allocs/op
BenchmarkMMH3-12              	 3655274	       331.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkCRC32-12             	 5002388	       244.1 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ISO-12          	 1000000	      1162 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ECMA-12         	 1000000	      1152 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32-12             	  491220	      2459 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32a-12            	  489003	      2451 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64-12             	  493273	      2440 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64a-12            	  495297	      2516 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv128-12            	  181533	      6573 ns/op	      24 B/op	       2 allocs/op
BenchmarkFnv128a-12           	  159890	      7530 ns/op	      24 B/op	       2 allocs/op
BenchmarkMD4-12               	  266040	      4532 ns/op	      16 B/op	       1 allocs/op
BenchmarkMD5-12               	  391636	      3090 ns/op	      16 B/op	       1 allocs/op
BenchmarkSHA1-12              	 1400571	       898.2 ns/op	      24 B/op	       1 allocs/op
BenchmarkSHA224-12            	 1221715	       868.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA256-12            	 1378735	       870.7 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA384-12            	  757268	      1531 ns/op	      48 B/op	       1 allocs/op
BenchmarkSHA512-12            	  796516	      1514 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA3256-12           	  208836	      5742 ns/op	     480 B/op	       2 allocs/op
BenchmarkSHA3512-12           	  119353	     10065 ns/op	     576 B/op	       3 allocs/op
BenchmarkRIPEMD160-12         	  189320	      6148 ns/op	      24 B/op	       1 allocs/op
BenchmarkWhirlpool-12         	   47458	     25466 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA256Parallel-12    	11993842	       100.0 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/hash	43.831s
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
BenchmarkMapStringKeys-12     	25698018	        45.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntKeys-12        	57162136	        18.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapStringIndex-12    	22459677	        49.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntIndex-12       	35164170	        29.37 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/index	6.343s
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
BenchmarkJsonMarshal-12      	 1645453	       720.8 ns/op	     480 B/op	       5 allocs/op
BenchmarkJsonUnmarshal-12    	  314257	      3880 ns/op	    1816 B/op	      27 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/json	3.384s
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
BenchmarkMathInt8-12           	1000000000	         0.3092 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt32-12          	1000000000	         0.3076 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt64-12          	1000000000	         0.3080 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt32-12    	290558541	         4.220 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt64-12    	287433439	         4.130 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathMutexInt-12       	145752284	         8.404 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat32-12        	1000000000	         0.3076 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat64-12        	1000000000	         0.3069 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/math	7.229s
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
BenchmarkParseBool-12     	522835948	         2.159 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt-12      	100000000	        10.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat-12    	18935383	        63.53 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/parse	4.112s
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
BenchmarkMathRand-12            	170200674	         6.958 ns/op	       0 B/op	       0 allocs/op
BenchmarkCryptoRand-12          	10248249	       115.9 ns/op	      48 B/op	       3 allocs/op
BenchmarkMathRandBytes-12       	18543044	        64.42 ns/op	      32 B/op	       1 allocs/op
BenchmarkCryptoRandBytes-12     	 4600580	       261.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkMathRandString-12      	 9909012	       119.8 ns/op	     128 B/op	       3 allocs/op
BenchmarkCryptoRandString-12    	 3888849	       311.7 ns/op	     128 B/op	       3 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/random	9.080s
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
BenchmarkMatchString-12            	  266162	      4548 ns/op	   10214 B/op	      86 allocs/op
BenchmarkMatchStringCompiled-12    	 3757484	       319.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchStringGolibs-12      	 3659856	       326.3 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/regexp	5.612s
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

	ustmt := prepareBenchmarkStmt(b, tx, `INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
	ostmt := prepareBenchmarkStmt(b, tx, `INSERT INTO orders (id, user_id, amount, status, meta) VALUES (?, ?, ?, ?, ?)`)

	for i := 1; i <= users; i++ {
		if _, err := ustmt.Exec(i, fmt.Sprintf("user%d", i), fmt.Sprintf("u%d@example.com", i), i%3 != 0); err != nil {
			b.Fatalf("populate user: %v", err)
		}
		for j := 0; j < ordersPerUser; j++ {
			id := i*100 + j
			status := "PAID"
			if j%4 == 0 {
				status = "PENDING"
			}
			if _, err := ostmt.Exec(id, i, float64(id)*1.5, status, `{"device":"web"}`); err != nil {
				b.Fatalf("populate order: %v", err)
			}
		}
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
	stmt := prepareBenchmarkStmt(b, tx, `INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)

	names := make([]string, b.N)
	emails := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		names[i] = fmt.Sprintf("user%d", i)
		emails[i] = fmt.Sprintf("u%d@example.com", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := stmt.Exec(i, names[i], emails[i], i%2 == 0); err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()

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
	b.Run("SelectJoin", func(b *testing.B) {
		benchSelectJoin(b, open, dsn("select_join"), 200)
	})
}

func sqliteDSN(name string) string {
	return "file:" + name + "?mode=memory&cache=shared"
}

func Benchmark_SQLite(b *testing.B) {
	runSQLBenchmarks(b, openSQLite, func(name string) string {
		return sqliteDSN("sqlite_" + name)
	})
}

func Benchmark_Insert_SQLite(b *testing.B) {
	benchInsertTxPerRow(b, openSQLite, sqliteDSN("insert_legacy"))
}

func Benchmark_SelectJoin_SQLite(b *testing.B) {
	benchSelectJoin(b, openSQLite, sqliteDSN("select_join_legacy"), 200)
}
//go:build tinysql
// +build tinysql

package main

import (
	"database/sql"
	"testing"

	_ "github.com/SimonWaldherr/tinySQL/internal/driver"
)

func Benchmark_Insert_TinySQL(b *testing.B) {
	benchInsert(b, func(dsn string) (*sql.DB, error) { return sql.Open("tinysql", dsn) }, "mem://?tenant=bench")
}

func Benchmark_SelectJoin_TinySQL(b *testing.B) {
	benchSelectJoin(b, func(dsn string) (*sql.DB, error) { return sql.Open("tinysql", dsn) }, "mem://?tenant=bench_select", 200)
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/sql
cpu: Apple M2 Max
Benchmark_SQLite/InsertTxPerRow-12         	  165094	      7298 ns/op	     993 B/op	      27 allocs/op
Benchmark_SQLite/InsertPrepared-12         	  849130	      1559 ns/op	     304 B/op	      10 allocs/op
Benchmark_SQLite/TransactionBatch100-12    	    7782	    173344 ns/op	   31031 B/op	    1111 allocs/op
Benchmark_SQLite/SelectPoint-12            	  697264	      1689 ns/op	     554 B/op	      20 allocs/op
Benchmark_SQLite/SelectRange25-12          	  101010	     11794 ns/op	    2221 B/op	     180 allocs/op
Benchmark_SQLite/UpdatePrepared-12         	  865459	      1381 ns/op	     162 B/op	       6 allocs/op
Benchmark_SQLite/SelectJoin-12             	    3432	    353008 ns/op	   16424 B/op	    1412 allocs/op
Benchmark_Insert_SQLite-12                 	  166222	      7450 ns/op	     992 B/op	      26 allocs/op
Benchmark_SelectJoin_SQLite-12             	    3484	    367422 ns/op	   16424 B/op	    1412 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/sql	12.119s
```

### sql (tinysql)

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

	ustmt := prepareBenchmarkStmt(b, tx, `INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
	ostmt := prepareBenchmarkStmt(b, tx, `INSERT INTO orders (id, user_id, amount, status, meta) VALUES (?, ?, ?, ?, ?)`)

	for i := 1; i <= users; i++ {
		if _, err := ustmt.Exec(i, fmt.Sprintf("user%d", i), fmt.Sprintf("u%d@example.com", i), i%3 != 0); err != nil {
			b.Fatalf("populate user: %v", err)
		}
		for j := 0; j < ordersPerUser; j++ {
			id := i*100 + j
			status := "PAID"
			if j%4 == 0 {
				status = "PENDING"
			}
			if _, err := ostmt.Exec(id, i, float64(id)*1.5, status, `{"device":"web"}`); err != nil {
				b.Fatalf("populate order: %v", err)
			}
		}
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
	stmt := prepareBenchmarkStmt(b, tx, `INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)

	names := make([]string, b.N)
	emails := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		names[i] = fmt.Sprintf("user%d", i)
		emails[i] = fmt.Sprintf("u%d@example.com", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := stmt.Exec(i, names[i], emails[i], i%2 == 0); err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()

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
	b.Run("SelectJoin", func(b *testing.B) {
		benchSelectJoin(b, open, dsn("select_join"), 200)
	})
}

func sqliteDSN(name string) string {
	return "file:" + name + "?mode=memory&cache=shared"
}

func Benchmark_SQLite(b *testing.B) {
	runSQLBenchmarks(b, openSQLite, func(name string) string {
		return sqliteDSN("sqlite_" + name)
	})
}

func Benchmark_Insert_SQLite(b *testing.B) {
	benchInsertTxPerRow(b, openSQLite, sqliteDSN("insert_legacy"))
}

func Benchmark_SelectJoin_SQLite(b *testing.B) {
	benchSelectJoin(b, openSQLite, sqliteDSN("select_join_legacy"), 200)
}
//go:build tinysql
// +build tinysql

package main

import (
	"database/sql"
	"testing"

	_ "github.com/SimonWaldherr/tinySQL/internal/driver"
)

func Benchmark_Insert_TinySQL(b *testing.B) {
	benchInsertTxPerRow(b, openTinySQL, "mem://?tenant=bench_insert")
}

func Benchmark_SelectJoin_TinySQL(b *testing.B) {
	benchSelectJoin(b, openTinySQL, "mem://?tenant=bench_select", 200)
}

func Benchmark_TinySQL(b *testing.B) {
	runSQLBenchmarks(b, openTinySQL, func(name string) string {
		return "mem://?tenant=bench_" + name
	})
}

func openTinySQL(dsn string) (*sql.DB, error) {
	return sql.Open("tinysql", dsn)
}
```

```
$ go test -bench . -benchmem -tags tinysql
# github.com/SimonWaldherr/golang-benchmarks/sql
bench_tinysql_test.go:10:2: no required module provides package github.com/SimonWaldherr/tinySQL/internal/driver; to add it:
	go get github.com/SimonWaldherr/tinySQL/internal/driver
FAIL	github.com/SimonWaldherr/golang-benchmarks/sql [setup failed]
(tinysql benchmarks skipped or failed)
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
BenchmarkTextTemplate-12    	 3247826	       376.4 ns/op	     272 B/op	       5 allocs/op
BenchmarkHTMLTemplate-12    	 1000000	      1083 ns/op	     496 B/op	      15 allocs/op
BenchmarkRegExp-12          	 2886160	       391.3 ns/op	     298 B/op	       9 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/template	4.627s
```

