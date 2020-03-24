// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to recover from panics.
package main

import (
	"fmt"
	"runtime"
)

func main() {

	// Call the testPanic function to run the test.
	if err := testPanic(); err != nil {
		fmt.Println("Error:", err)
	}
}

// testPanic simulates a function that encounters a panic to
// test our catchPanic function.
func testPanic() (err error) {
	// Schedule the catchPanic function to be called when
	// the testPanic function returns.
	defer catchPanic(&err)

	fmt.Println("Start test")

	// Mimic a traditional error from a function.
	err = mimicError("1")

	// Trying to dereference a nil pointer will cause the
	// runtime to panic.
	var p *int
	*p = 10

	fmt.Println("End test")

	return err
}

// catchPanic catches panics and processes the error.
func catchPanic(err *error) {
	// Check if a panic occurred.
	if r := recover(); r != nil {
		fmt.Println("PANIC deferred")

		// Capture the stack trace.
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		fmt.Println("Stack trace:", string(buf))

		// If the caller wants the error back provide it.
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func mimicError(key string) error {
	return nil
}

// Outputs:
// Start test
// PANIC deferred
// Stack trace: goroutine 1 [running]:
// main.catchPanic(0xc00006eef0)
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/functions/advanced/examp
// le1/example1.go:52 +0x10d
// panic(0x4a9380, 0x55b8d0)
//         /usr/local/go/src/runtime/panic.go:679 +0x1b2
// main.testPanic(0x0, 0x0)
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/functions/advanced/examp
// le1/example1.go:37 +0xe2
// main.main()
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/functions/advanced/examp
// le1/example1.go:17 +0x26

// Error: runtime error: invalid memory address or nil pointer dereference
