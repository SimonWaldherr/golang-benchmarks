@echo off

setlocal
set "benchs=base64 between concat contains foreach hash index json math parse random regexp"

> README.md (
echo.# Go Benchmarks
echo.
echo.In programming in general, and in Golang in particular, many roads lead to Rome.
echo.From time to time I ask myself which of these ways is the fastest. 
echo.In Golang there is a wonderful solution, with \`go test -bench\` you can measure the speed very easily and quickly.
echo.In order for you to benefit from it too, I will publish such benchmarks in this repository in the future.
echo.
echo.## ToC
echo.
)

for %%i in (%benchs%) do (
  echo * [%%i](https://github.com/SimonWaldherr/golang-benchmarks#%%i^) >> README.md
)

>> README.md (
echo. 
echo.## Golang?
echo. 
echo.I published another repository where I show some Golang examples.
echo.If you\'re interested in new programming languages, you should definitely take a look at Golang:
echo. 
echo.* [Golang examples](https://github.com/SimonWaldherr/golang-examples^)
echo.* [tour.golang.org](https://tour.golang.org/^)
echo.* [Go by example](https://gobyexample.com/^)
echo.* [Golang Book](http://www.golang-book.com/^)
echo.* [Go-Learn](https://github.com/skippednote/Go-Learn^)
echo. 
echo.## Is it any good?
echo. 
echo.[Yes](https://news.ycombinator.com/item?id=3067434^)
echo.
echo.## Benchmark Results
echo.
)

go fmt ./...

for /f "delims=;" %%i in ('go version') do set GO_VERSION=%%i

echo.Golang Version: %GO_VERSION% >> README.md
echo. >> README.md

for %%i in (%benchs%) do (

@cd %%i
@gofmt -s -w -l -e .
    
echo.### %%i
echo.
echo.```go
type *_test.go
echo.```
echo.
echo.```
echo.$ go test -bench . -benchmem
go test -bench . -benchmem
echo.```
echo.

@cd ..
) >> README.win.md
