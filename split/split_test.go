// Package split benchmarks various string splitting methods.
package split

import (
	"strings"
	"testing"
)

var splitResult []string

func BenchmarkSplitMethods(b *testing.B) {
	input := "one,two,three,four,five"
	comma := func(r rune) bool { return r == ',' }

	b.Run("Strings.Split", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			splitResult = strings.Split(input, ",")
		}
	})

	b.Run("Strings.SplitN", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			splitResult = strings.SplitN(input, ",", 3)
		}
	})

	b.Run("Strings.CutLoop", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			splitResult = splitWithCut(input)
		}
	})

	b.Run("Strings.FieldsFunc", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			splitResult = strings.FieldsFunc(input, comma)
		}
	})
}

func splitWithCut(s string) []string {
	parts := make([]string, 0, strings.Count(s, ",")+1)
	for {
		before, after, found := strings.Cut(s, ",")
		parts = append(parts, before)
		if !found {
			return parts
		}
		s = after
	}
}
