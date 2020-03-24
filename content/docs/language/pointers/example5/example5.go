// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how stacks grow/change.
package main

// Number of elements to grow each stack frame.
// Run with 1 and then with 1024
const size = 1024

// main is the entry point for the application.
func main() {
	s := "HELLO"
	stackCopy(&s, 0, [size]int{})
}

// stackCopy recursively runs increasing the size
// of the stack.
//go:noinline
func stackCopy(s *string, c int, a [size]int) {
	println(c, s, *s)

	c++
	if c == 10 {
		return
	}

	stackCopy(s, c, a)
}

// Outputs - run with 1:
// 0 0xc00003a740 HELLO
// 1 0xc00003a740 HELLO
// 2 0xc00003a740 HELLO
// 3 0xc00003a740 HELLO
// 4 0xc00003a740 HELLO
// 5 0xc00003a740 HELLO
// 6 0xc00003a740 HELLO
// 7 0xc00003a740 HELLO
// 8 0xc00003a740 HELLO
// 9 0xc00003a740 HELLO

// Outputs - run with 1024:
// 0 0xc000091f40 HELLO
// 1 0xc000091f40 HELLO
// 2 0xc0000a1f40 HELLO
// 3 0xc0000a1f40 HELLO
// 4 0xc0000a1f40 HELLO
// 5 0xc0000a1f40 HELLO
// 6 0xc0000c1f40 HELLO
// 7 0xc0000c1f40 HELLO
// 8 0xc0000c1f40 HELLO
// 9 0xc0000c1f40 HELLO
