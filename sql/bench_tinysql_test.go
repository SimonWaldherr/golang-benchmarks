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
