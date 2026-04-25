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
