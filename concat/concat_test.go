// Package concat benchmarks the performance of
// various string concatenation methods.
// Instead of just concatenating a string to another string
// it is also possible (and much faster) to use
// a bytes buffer.
package concat

import (
	"bytes"
	"strings"
	"testing"
)

const concatParts = 64

var concatResult string

func BenchmarkConcatString(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < concatParts; i++ {
			str += "x"
		}
		concatResult = str
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var buffer bytes.Buffer
		for i := 0; i < concatParts; i++ {
			buffer.WriteString("x")
		}
		concatResult = buffer.String()
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var builder strings.Builder
		for i := 0; i < concatParts; i++ {
			builder.WriteString("x")
		}
		concatResult = builder.String()
	}
}

func BenchmarkConcatBuilderGrow(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var builder strings.Builder
		builder.Grow(concatParts)
		for i := 0; i < concatParts; i++ {
			builder.WriteString("x")
		}
		concatResult = builder.String()
	}
}

func BenchmarkConcat(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var str string
			for i := 0; i < concatParts; i++ {
				str += "x"
			}
			concatResult = str
		}
	})

	b.Run("Buffer", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var buffer bytes.Buffer
			for i := 0; i < concatParts; i++ {
				buffer.WriteString("x")
			}
			concatResult = buffer.String()
		}
	})

	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var builder strings.Builder
			for i := 0; i < concatParts; i++ {
				builder.WriteString("x")
			}
			concatResult = builder.String()
		}
	})

	b.Run("BuilderGrow", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var builder strings.Builder
			builder.Grow(concatParts)
			for i := 0; i < concatParts; i++ {
				builder.WriteString("x")
			}
			concatResult = builder.String()
		}
	})
}
