// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Declare three variables that are initialized to their zero value and three
// declared with a literal value. Declare variables of type string, int and
// bool. Display the values of those variables.
//
// Declare a new variable of type float32 and initialize the variable by
// converting the literal value of Pi (3.14).
package main

import "fmt"

// main is the entry point for the application.
func main() {

	// Declare variables that are set to their zero value.
	var name string
	var level int

	// Display the value of those variables.
	fmt.Println(name)
	fmt.Println(level)

	// Declare variables and initialize.
	// Using the short variable declaration operator.
	count := 6
	pi := 3.14159

	// Display the value of those variables.
	fmt.Println(count)
	fmt.Printf("value: %f, type: %T\n", pi, pi)

	// Perform a type conversion.
	pie := float32(pi) // convert float64 to float32

	// Display the value of that variable.
	fmt.Println(pie)
}
