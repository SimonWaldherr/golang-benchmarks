package caseinsensitivecompare

import (
	"strings"
	"testing"
)

var caseInsensitiveCompareResult bool

func BenchmarkEqualFold(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.EqualFold("abc", "ABC")
		caseInsensitiveCompareResult = strings.EqualFold("ABC", "ABC")
		caseInsensitiveCompareResult = strings.EqualFold("1aBcD", "1AbCd")
	}
}

func BenchmarkToUpper(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.ToUpper("abc") == strings.ToUpper("ABC")
		caseInsensitiveCompareResult = strings.ToUpper("ABC") == strings.ToUpper("ABC")
		caseInsensitiveCompareResult = strings.ToUpper("1aBcD") == strings.ToUpper("1AbCd")
	}
}

func BenchmarkToLower(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		caseInsensitiveCompareResult = strings.ToLower("abc") == strings.ToLower("ABC")
		caseInsensitiveCompareResult = strings.ToLower("ABC") == strings.ToLower("ABC")
		caseInsensitiveCompareResult = strings.ToLower("1aBcD") == strings.ToLower("1AbCd")
	}
}
