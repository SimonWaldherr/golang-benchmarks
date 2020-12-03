# Go Benchmarks

In programming in general, and in Golang in particular, many roads lead to Rome.
From time to time I ask myself which of these ways is the fastest. 
In Golang there is a wonderful solution, with \`go test -bench\` you can measure the speed very easily and quickly.
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

Golang Version: go version devel +a6755fc0de Sat Nov 7 07:33:23 2020 +0000 windows/amd64 
 
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
goos: windows
goarch: amd64
pkg: github.com/SimonWaldherr/golang-benchmarks/base64
cpu: AMD Ryzen 5 3500U with Radeon Vega Mobile Gfx  
BenchmarkBase64decode-8   	 5367226	       245.6 ns/op	      32 B/op	       2 allocs/op
BenchmarkBase64regex-8    	   30693	     37875 ns/op	   21444 B/op	     198 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/base64	3.170s
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
goos: windows
goarch: amd64
pkg: github.com/SimonWaldherr/golang-benchmarks/between
cpu: AMD Ryzen 5 3500U with Radeon Vega Mobile Gfx  
BenchmarkNumberRegEx-8     	   45007	     26404 ns/op	   16139 B/op	     142 allocs/op
BenchmarkFulltextRegEx-8   	   58605	     21743 ns/op	   11640 B/op	     104 allocs/op
BenchmarkNumberParse-8     	12766242	        95.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkFulltextParse-8   	  800138	      1296 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/between	5.382s
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
goos: windows
goarch: amd64
pkg: github.com/SimonWaldherr/golang-benchmarks/concat
cpu: AMD Ryzen 5 3500U with Radeon Vega Mobile Gfx  
BenchmarkConcatString-8    	  654996	     64654 ns/op	  331436 B/op	       1 allocs/op
BenchmarkConcatBuffer-8    	96764048	        12.40 ns/op	       2 B/op	       0 allocs/op
BenchmarkConcatBuilder-8   	