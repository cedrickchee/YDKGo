// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Declare a nil slice of integers. Create a loop that appends 10 values to the
// slice. Iterate over the slice and display each value.
//
// Declare a slice of five strings and initialize the slice with string literal
// values. Display all the elements. Take a slice of index one and two
// and display the index position and value of each element in the new slice.
package main

import "fmt"

// Add imports.

func main() {

	// Declare a nil slice of integers.
	var a []int

	// Append numbers to the slice.
	for i := 0; i < 10; i++ {
		a = append(a, i)
	}

	// Display each value in the slice.
	for _, v := range a {
		fmt.Println(v)
	}

	// Declare a slice of strings and populate the slice with names.
	names := []string{"John", "Jane", "David", "Emmy"}

	// Display each index position and slice value.
	fmt.Println("*****************************")
	for i, v := range names {
		fmt.Printf("%d - %s\n", i, v)
	}

	// Take a slice of index 1 and 2 of the slice of strings.
	family := names[:2]

	// Display each index position and slice values for the new slice.
	fmt.Println("*****************************")
	for i, v := range family {
		fmt.Printf("%d - %s\n", i, v)
	}
}
