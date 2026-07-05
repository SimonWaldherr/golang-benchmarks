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
