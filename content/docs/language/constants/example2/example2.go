// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how constants do have a parallel type system.
package main

import "fmt"

const (
	// Max integer value on 64 bit architecture.
	maxInt = 9223372036854775807

	// Much larger value than int64.
	bigger = 9223372036854775808543522345

	// Will NOT compile
	// Compiler: "constant 9223372036854775808543522345 overflows int64"
	// biggerInt int64 = 9223372036854775808543522345
)

func main() {
	fmt.Println("Will Compile")
}

// Outputs:
// Will Compile
