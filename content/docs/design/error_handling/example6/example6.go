// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how wrapping errors work.
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// AppError represents a custom error type.
type AppError struct {
	State int
}

// Error implements the error interface.
func (c *AppError) Error() string {
	return fmt.Sprintf("App Error, State: %d", c.State)
}

func main() {

	// Make the function call and validate the error.
	if err := firstCall(10); err != nil {

		// Use type as context to determine cause.
		switch v := errors.Cause(err).(type) {
		case *AppError:

			// We got our custom error type.
			fmt.Println("Custom App Error:", v.State)

		default:

			// We did not get any specific error type.
			fmt.Println("Default Error")
		}

		// Display the stack trace for the error.
		fmt.Println("\nStack Trace\n********************************")
		fmt.Printf("%+v\n", err)
		fmt.Println("\nNo Trace\n********************************")
		fmt.Printf("%v\n", err)
	}
}

// firstCall makes a call to a second function and wraps any error.
func firstCall(i int) error {
	if err := secondCall(i); err != nil {
		return errors.Wrapf(err, "firstCall->secondCall(%d)", i)
	}
	return nil
}

// secondCall makes a call to a third function and wraps any error.
func secondCall(i int) error {
	if err := thirdCall(); err != nil {
		return errors.Wrap(err, "secondCall->thirdCall()")
	}
	return nil
}

// thirdCall create an error value we will validate.
func thirdCall() error {
	return &AppError{99}
}

// Outputs:
// Custom App Error: 99

// Stack Trace
// ********************************
// App Error, State: 99
// secondCall->thirdCall()
// main.secondCall
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/design/error_handling/example6/example6.go:62
// main.firstCall
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/design/error_handling/example6/example6.go:53
// main.main
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/design/error_handling/example6/example6.go:28
// runtime.main
//         /usr/local/go/src/runtime/proc.go:203
// runtime.goexit
//         /usr/local/go/src/runtime/asm_amd64.s:1357
// firstCall->secondCall(10)
// main.firstCall
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/design/error_handling/example6/example6.go:54
// main.main
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/design/error_handling/example6/example6.go:28
// runtime.main
//         /usr/local/go/src/runtime/proc.go:203
// runtime.goexit
//         /usr/local/go/src/runtime/asm_amd64.s:1357

// No Trace
// ********************************
// firstCall->secondCall(10): secondCall->thirdCall(): App Error, State: 99
