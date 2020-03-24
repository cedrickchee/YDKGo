// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go test -run none -bench . -benchtime 3s -benchmem

// Write three benchmark tests for converting an integer into a string. First using the
// fmt.Sprintf function, then the strconv.FormatInt function and then strconv.Itoa.
// Identify which function performs the best.
package main

import (
	"fmt"
	"strconv"
	"testing"
)

var gs string

// BenchmarkSprintf provides performance numbers for the fmt.Sprintf function.
func BenchmarkSprintf(b *testing.B) {
	number := 6
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("%d", number)
	}

	gs = s
}

// BenchmarkFormat provides performance numbers for the strconv.FormatInt function.
func BenchmarkFormat(b *testing.B) {
	number := int64(6)
	var s string

	for i := 0; i < b.N; i++ {
		s = strconv.FormatInt(number, 10)
	}

	gs = s
}

// BenchmarkItoa provides performance numbers for the strconv.Itoa function.
func BenchmarkItoa(b *testing.B) {
	number := 6
	var s string

	for i := 0; i < b.N; i++ {
		s = strconv.Itoa(number)
	}

	gs = s
}

// Output:
// $ go test -run none -bench . -benchtime 3s -benchmem
// goos: linux
// goarch: amd64
// pkg: github.com/cedrickchee/ultimate-go/testing/benchmarks/exercises/exercise1
// BenchmarkSprintf-4      33091200               104 ns/op               8 B/op          1 allocs/op
// BenchmarkFormat-4       896546001                3.93 ns/op            0 B/op          0 allocs/op
// BenchmarkItoa-4         856785422                3.98 ns/op            0 B/op          0 allocs/op
// PASS
// ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/exercises/exercise1       11.313s
