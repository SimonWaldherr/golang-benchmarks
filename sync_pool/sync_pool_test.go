package sync_pool

import (
	"sync"
	"testing"
)

type Item struct {
	buf [256]byte
}

func BenchmarkNewPerOp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &Item{}
	}
}

func BenchmarkSyncPool(b *testing.B) {
	var pool = sync.Pool{New: func() interface{} { return &Item{} }}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := pool.Get().(*Item)
		pool.Put(it)
	}
}
