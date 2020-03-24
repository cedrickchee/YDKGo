// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go test -run none -bench . -benchtime 3s -benchmem

// Basic benchmark test.
package basic

import (
	"fmt"
	"testing"
)

var gs string

// BenchmarkSprint tests the performance of using Sprint.
func BenchmarkSprint(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

	gs = s
}

// BenchmarkSprint tests the performance of using Sprintf.
func BenchmarkSprintf(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}

	gs = s
}

// Output:
// $ go test -run none -bench . -benchtime 3s -benchmem
// goos: linux
// goarch: amd64
// pkg: github.com/cedrickchee/ultimate-go/testing/benchmarks/basic
// BenchmarkSprint-4       33619288               103 ns/op               5 B/op          1 allocs/op
// BenchmarkSprintf-4      43777660                81.0 ns/op             5 B/op          1 allocs/op
// PASS
// ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/basic     9.914s
