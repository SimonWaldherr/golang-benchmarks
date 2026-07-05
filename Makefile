GO ?= go
BENCH ?= .
BENCHTIME ?= 1s
PKG ?= ./...

.PHONY: all fmt test bench bench-sql update-readme update-deps tidy clean

all: fmt test

fmt:
	$(GO) fmt ./...
	gofmt -s -w $$(find . -name '*.go' -not -path './.git/*')

test:
	$(GO) test $(PKG)

bench:
	$(GO) test -bench $(BENCH) -benchmem -benchtime $(BENCHTIME) $(PKG)

bench-sql:
	$(GO) test -bench $(BENCH) -benchmem -benchtime $(BENCHTIME) ./sql

update-readme:
	./updateBenchLogs.sh

update-deps:
	$(GO) get -u all
	$(GO) get github.com/SimonWaldherr/tinySQL@latest github.com/zeebo/blake3@latest golang.org/x/crypto@latest modernc.org/sqlite@latest simonwaldherr.de/go/golibs@latest simonwaldherr.de/go/ranger@latest
	$(GO) mod tidy

tidy:
	$(GO) mod tidy

clean:
	$(GO) clean -testcache
