package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

func openSQLite(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn)
}

// prepareSchema creates the users and orders tables used by benchmarks.
func prepareSchema(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE users (id INT, name TEXT, email TEXT, active BOOL)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE orders (id INT, user_id INT, amount FLOAT, status TEXT, meta JSON)`); err != nil {
		return err
	}
	return nil
}

func benchInsert(b *testing.B, open func(string) (*sql.DB, error), dsn string) {
	db, err := open(dsn)
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer db.Close()

	if err := prepareSchema(db); err != nil {
		b.Fatalf("schema: %v", err)
	}

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

func benchSelectJoin(b *testing.B, open func(string) (*sql.DB, error), dsn string, rows int) {
	db, err := open(dsn)
	if err != nil {
		b.Fatalf("open: %v", err)
	}
	defer db.Close()

	if err := prepareSchema(db); err != nil {
		b.Fatalf("schema: %v", err)
	}

	// populate a modest dataset
	tx, err := db.Begin()
	if err != nil {
		b.Fatalf("begin: %v", err)
	}
	ustmt, _ := tx.Prepare(`INSERT INTO users (id, name, email, active) VALUES (?, ?, ?, ?)`)
	ostmt, _ := tx.Prepare(`INSERT INTO orders (id, user_id, amount, status, meta) VALUES (?, ?, ?, ?, ?)`)
	for i := 1; i <= rows; i++ {
		if _, err := ustmt.Exec(i, fmt.Sprintf("user%d", i), nil, true); err != nil {
			b.Fatalf("populate user: %v", err)
		}
		// attach a couple orders for each user
		for j := 0; j < 2; j++ {
			id := i*10 + j
			if _, err := ostmt.Exec(id, i, float64(id)*1.5, "PAID", `{"device":"web"}`); err != nil {
				b.Fatalf("populate order: %v", err)
			}
		}
	}
	ustmt.Close()
	ostmt.Close()
	if err := tx.Commit(); err != nil {
		b.Fatalf("commit populate: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`
            SELECT u.name AS user, SUM(o.amount) AS total, COUNT(*) AS cnt
            FROM users u
            LEFT JOIN orders o ON u.id = o.user_id AND o.status = 'PAID'
            GROUP BY u.name
        `)
		if err != nil {
			b.Fatalf("query: %v", err)
		}
		for rows.Next() {
			var name string
			var total float64
			var cnt int
			if err := rows.Scan(&name, &total, &cnt); err != nil {
				b.Fatalf("scan: %v", err)
			}
		}
		rows.Close()
	}
}

func Benchmark_Insert_SQLite(b *testing.B) {
	// modernc.org/sqlite uses driver name "sqlite" and accepts a file: DSN
	benchInsert(b, openSQLite, "file::memory:?mode=memory&cache=shared")
}

func Benchmark_SelectJoin_SQLite(b *testing.B) {
	benchSelectJoin(b, openSQLite, "file::memory:?mode=memory&cache=shared", 200)
}
