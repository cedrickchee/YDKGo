// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go test -run none -bench . -benchtime 3s

// Tests to show how Data Oriented Design matters.
package caching

import "testing"

var fa int

// Capture the time it takes to perform a link list traversal.
func BenchmarkLinkListTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = LinkedListTraverse()
	}

	fa = a
}

// Capture the time it takes to perform a column traversal.
func BenchmarkColumnTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = ColumnTraverse()
	}

	fa = a
}

// Capture the time it takes to perform a row traversal.
func BenchmarkRowTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = RowTraverse()
	}

	fa = a
}

// Benchmark run results:
// $ go test -run none -bench . -benchtime 3s
// Elements in the link list 16777216
// Elements in the matrix 16777216
// goos: linux
// goarch: amd64
// pkg: github.com/cedrickchee/ultimate-go/language/arrays/benchmarks/caching
// BenchmarkLinkListTraverse-4           93          36718229 ns/op
// BenchmarkColumnTraverse-4             19         184652302 ns/op
// BenchmarkRowTraverse-4               165          21787691 ns/op
// PASS
// ok      github.com/cedrickchee/ultimate-go/language/arrays/benchmarks/caching   21.069s
