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
