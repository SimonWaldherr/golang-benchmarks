// Package split benchmarks various string splitting methods.
package split

import (
	"strings"
	"testing"
)

func BenchmarkSplitMethods(b *testing.B) {
	b.Run("Strings.Split", func(b *testing.B) {
		str := "one,two,three,four,five"
		for n := 0; n < b.N; n++ {
			_ = strings.Split(str, ",")
		}
	})

	b.Run("Strings.SplitN", func(b *testing.B) {
		str := "one,two,three,four,five"
		for n := 0; n < b.N; n++ {
			_ = strings.SplitN(str, ",", 3)
		}
	})

	b.Run("Strings.FieldsFunc", func(b *testing.B) {
		str := "one,two,three,four,five"
		for n := 0; n < b.N; n++ {
			_ = strings.FieldsFunc(str, func(r rune) bool { return r == ',' })
		}
	})
}
