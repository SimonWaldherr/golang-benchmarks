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

Golang Version: [go version go1.25.3 darwin/arm64](https://tip.golang.org/doc/go1.25)  
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
BenchmarkBase64decode-12    	23728497	        49.39 ns/op	      32 B/op	       2 allocs/op
BenchmarkBase64regex-12     	  135303	      8744 ns/op	   21901 B/op	     198 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/base64	3.719s
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
BenchmarkNumberRegEx-12      	  147270	      8500 ns/op	   16881 B/op	     142 allocs/op
BenchmarkFulltextRegEx-12    	  142051	     11951 ns/op	   12055 B/op	     104 allocs/op
BenchmarkNumberParse-12      	28720615	        45.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkFulltextParse-12    	 1860694	       692.4 ns/op	      32 B/op	       2 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/between	6.805s
```

### caseinsensitivecompare

```go
package trim

import (
	"strings"
	"testing"
)

func BenchmarkEqualFold(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = strings.EqualFold("abc", "ABC")
		_ = strings.EqualFold("ABC", "ABC")
		_ = strings.EqualFold("1aBcD", "1AbCd")
	}
}

func BenchmarkToUpper(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = strings.ToUpper("abc") == strings.ToUpper("ABC")
		_ = strings.ToUpper("ABC") == strings.ToUpper("ABC")
		_ = strings.ToUpper("1aBcD") == strings.ToUpper("1AbCd")
	}
}

func BenchmarkToLower(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = strings.ToLower("abc") == strings.ToLower("ABC")
		_ = strings.ToLower("ABC") == strings.ToLower("ABC")
		_ = strings.ToLower("1aBcD") == strings.ToLower("1AbCd")
	}
}
```

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/caseinsensitivecompare
cpu: Apple M2 Max
BenchmarkEqualFold-12    	63957488	        18.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkToUpper-12      	 7216710	       165.3 ns/op	      24 B/op	       3 allocs/op
BenchmarkToLower-12      	 5655952	       234.7 ns/op	      40 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/caseinsensitivecompare	5.248s
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

func BenchmarkConcat(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		var str string
		for n := 0; n < b.N; n++ {
			str += "x"
		}
	})

	b.Run("Buffer", func(b *testing.B) {
		var buffer bytes.Buffer
		for n := 0; n < b.N; n++ {
			buffer.WriteString("x")
		}
	})

	b.Run("Builder", func(b *testing.B) {
		var builder strings.Builder
		for n := 0; n < b.N; n++ {
			builder.WriteString("x")
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
BenchmarkConcatString-12     	 1000000	     32225 ns/op	  503994 B/op	       1 allocs/op
BenchmarkConcatBuffer-12     	348260655	         3.649 ns/op	       3 B/op	       0 allocs/op
BenchmarkConcatBuilder-12    	569981112	         2.272 ns/op	       5 B/op	       0 allocs/op
BenchmarkConcat/String-12    	 1000000	     35984 ns/op	  503994 B/op	       1 allocs/op
BenchmarkConcat/Buffer-12    	347776692	         3.509 ns/op	       3 B/op	       0 allocs/op
BenchmarkConcat/Builder-12   	561156622	         2.273 ns/op	       5 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/concat	74.748s
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
BenchmarkContains-12            	255654132	         4.447 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsNot-12         	211925643	         5.653 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytes-12       	229659852	         5.126 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsBytesNot-12    	185971905	         6.666 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatch-12        	24490878	        61.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompileMatchNot-12     	50348594	        25.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatch-12               	 1716830	       719.9 ns/op	    1401 B/op	      17 allocs/op
BenchmarkMatchNot-12            	 1806837	       683.6 ns/op	    1401 B/op	      17 allocs/op
BenchmarkContainsMethods/Strings.Contains-12         	100000000	        10.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/Bytes.Contains-12           	100000000	        10.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatchString-12         	17377903	        68.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsMethods/RegexMatch-12               	  938626	      1333 ns/op	    2801 B/op	      34 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/contains	18.847s
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
BenchmarkEmbed-12             	1000000000	         0.3024 ns/op	       0 B/op	       0 allocs/op
BenchmarkReadFile-12          	  106364	     10251 ns/op	     840 B/op	       5 allocs/op
BenchmarkIoutilReadFile-12    	  122562	      9942 ns/op	     840 B/op	       5 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/embed	3.178s
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
BenchmarkFloodFillRecursive-12    	 5839147	       179.5 ns/op	     432 B/op	       7 allocs/op
BenchmarkFloodFillDFS-12          	 1521726	       804.6 ns/op	    1744 B/op	      48 allocs/op
BenchmarkFloodFillBFS-12          	 1231497	      1028 ns/op	    2704 B/op	      53 allocs/op
BenchmarkFloodFillStack4Way-12    	 1434441	       876.0 ns/op	    1744 B/op	      48 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/floodfill	7.954s
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
BenchmarkForMap-12           	64167015	        18.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeMap-12         	28646996	        42.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSlice-12       	444200887	         2.656 ns/op	       0 B/op	       0 allocs/op
BenchmarkRangeSliceKey-12    	452416611	         2.669 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/foreach	6.291s
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
BenchmarkAdler32-12           	 1762033	       675.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkBCryptCost4-12       	    1208	    997697 ns/op	    8166 B/op	       9 allocs/op
BenchmarkBCryptCost10-12      	      19	  59828660 ns/op	    8268 B/op	       9 allocs/op
BenchmarkBCryptCost16-12      	       1	3796500708 ns/op	   10184 B/op	      13 allocs/op
BenchmarkBlake2b256-12        	  478257	      2558 ns/op	      32 B/op	       1 allocs/op
BenchmarkBlake2b512-12        	  475470	      2553 ns/op	      64 B/op	       1 allocs/op
BenchmarkBlake3256-12         	  411357	      2926 ns/op	      32 B/op	       1 allocs/op
BenchmarkMMH3-12              	 3683704	       325.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkCRC32-12             	 5036356	       234.4 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ISO-12          	 1000000	      1128 ns/op	       8 B/op	       1 allocs/op
BenchmarkCRC64ECMA-12         	 1000000	      1140 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32-12             	  477223	      2446 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv32a-12            	  489124	      2403 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64-12             	  512338	      2375 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv64a-12            	  517750	      2402 ns/op	       8 B/op	       1 allocs/op
BenchmarkFnv128-12            	  191343	      6312 ns/op	      24 B/op	       2 allocs/op
BenchmarkFnv128a-12           	  162915	      7244 ns/op	      24 B/op	       2 allocs/op
BenchmarkMD4-12               	  261038	      4424 ns/op	      24 B/op	       2 allocs/op
BenchmarkMD5-12               	  402916	      2959 ns/op	      16 B/op	       1 allocs/op
BenchmarkSHA1-12              	 1435572	       845.1 ns/op	      24 B/op	       1 allocs/op
BenchmarkSHA224-12            	 1417270	       865.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA256-12            	 1406242	       863.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkSHA384-12            	  803680	      1495 ns/op	      48 B/op	       1 allocs/op
BenchmarkSHA512-12            	  809352	      1491 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA3256-12           	  214844	      5605 ns/op	     480 B/op	       2 allocs/op
BenchmarkSHA3512-12           	  122396	      9833 ns/op	     576 B/op	       3 allocs/op
BenchmarkRIPEMD160-12         	  194943	      6258 ns/op	      24 B/op	       1 allocs/op
BenchmarkWhirlpool-12         	   44792	     26769 ns/op	      64 B/op	       1 allocs/op
BenchmarkSHA256Parallel-12    	11694115	       103.5 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/hash	42.425s
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
BenchmarkMapStringKeys-12     	22595022	        50.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntKeys-12        	56522802	        21.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapStringIndex-12    	21632437	        51.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkMapIntIndex-12       	52751019	        26.77 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/index	7.016s
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
BenchmarkJsonMarshal-12      	 1710452	       674.3 ns/op	     480 B/op	       5 allocs/op
BenchmarkJsonUnmarshal-12    	  341144	      3583 ns/op	    1816 B/op	      27 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/json	3.184s
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
BenchmarkMathInt8-12           	1000000000	         0.3041 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt32-12          	1000000000	         0.3003 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathInt64-12          	1000000000	         0.2972 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt32-12    	295297237	         4.033 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathAtomicInt64-12    	293039517	         4.090 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathMutexInt-12       	147570950	         7.976 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat32-12        	1000000000	         0.2979 ns/op	       0 B/op	       0 allocs/op
BenchmarkMathFloat64-12        	1000000000	         0.2951 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/math	7.198s
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
BenchmarkParseBool-12     	1000000000	         0.3027 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseInt-12      	122210564	         9.882 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFloat-12    	19550142	        60.96 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/parse	4.216s
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

func BenchmarkCryptoRandBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := GenerateRandomBytes(32)
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
goarch: arm64
pkg: github.com/SimonWaldherr/golang-benchmarks/random
cpu: Apple M2 Max
BenchmarkMathRand-12            	168541000	         6.786 ns/op	       0 B/op	       0 allocs/op
BenchmarkCryptoRand-12          	 9774474	       121.5 ns/op	      48 B/op	       3 allocs/op
BenchmarkCryptoRandString-12    	10994502	       123.6 ns/op	     128 B/op	       3 allocs/op
BenchmarkCryptoRandBytes-12     	17583844	        63.77 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/random	5.840s
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
BenchmarkMatchString-12            	  297730	      4326 ns/op	   10228 B/op	      86 allocs/op
BenchmarkMatchStringCompiled-12    	 3917988	       303.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchStringGolibs-12      	 3859442	       310.7 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/regexp	5.498s
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
BenchmarkTextTemplate-12    	 3585704	       339.8 ns/op	     272 B/op	       5 allocs/op
BenchmarkHTMLTemplate-12    	 1000000	      1076 ns/op	     496 B/op	      15 allocs/op
BenchmarkRegExp-12          	 2882409	       411.6 ns/op	     298 B/op	       9 allocs/op
PASS
ok  	github.com/SimonWaldherr/golang-benchmarks/template	4.589s
```

### trim

```go
```

```
$ go test -bench . -benchmem
FAIL	. [setup failed]
```

