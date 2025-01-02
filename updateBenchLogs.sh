#!/bin/bash

declare -a benchs=(base64 between caseinsensitivecompare concat contains embed floodfill foreach hash index json math parse random regexp template trim)

cat > README.md <<- EOM
# Go Benchmarks

[![DOI](https://zenodo.org/badge/154216722.svg)](https://zenodo.org/badge/latestdoi/154216722) 
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonWaldherr/golang-benchmarks)](https://goreportcard.com/report/github.com/SimonWaldherr/golang-benchmarks) 
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  

In programming in general, and in Golang in particular, many roads lead to Rome.
From time to time I ask myself which of these ways is the fastest. 
In Golang there is a wonderful solution, with \`go test -bench\` you can measure the speed very easily and quickly.
In order for you to benefit from it too, I will publish such benchmarks in this repository in the future.

## ToC

EOM

for i in "${benchs[@]}"
do
    echo "* [$i](https://github.com/SimonWaldherr/golang-benchmarks#$i)" >> README.md
done

cat >> README.md <<- EOM

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

EOM

go fmt ./...

echo -n "Golang Version: "  >> README.md
go version                  >> README.md
echo ""                     >> README.md

for i in "${benchs[@]}"
do
    cd $i
    gofmt -s -w -l -e .
    echo "### $i"                       >> ../README.md
    echo                                >> ../README.md
    echo "\`\`\`go"                     >> ../README.md
    cat *_test.go                       >> ../README.md
    echo "\`\`\`"                       >> ../README.md
    echo                                >> ../README.md
    echo "\`\`\`"                       >> ../README.md
    echo "$ go test -bench . -benchmem" >> ../README.md
    go test -bench . -benchmem          >> ../README.md
    echo "\`\`\`"                       >> ../README.md
    echo                                >> ../README.md
    cd ..
done
