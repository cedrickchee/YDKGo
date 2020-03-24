// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Declare an array of 5 strings with each element initialized to its zero value.
//
// Declare a second array of 5 strings and initialize this array with literal string
// values. Assign the second array to the first and display the results of the first array.
// Display the string value and address of each element.
package main

import "fmt"

func main() {

	// Declare an array of 5 strings set to its zero value.
	var fruits [5]string

	// Declare an array of 5 strings and pre-populate it with names.
	pets := [5]string{"cat", "dog", "rabbit", "fish", "hamster"}

	// Assign the populated array to the array of zero values.
	fruits = pets

	// Iterate over the first array declared.
	// Display the string value and address of each element.
	for i, v := range fruits {
		fmt.Printf("Value[%s]\tAddress[%p]\tIndexAddr[%p]\n", v, &v, &fruits[i])
	}
}

// Outputs:
// Value[cat]      Address[0xc00007c1e0]   IndexAddr[0xc0000a6000]
// Value[dog]      Address[0xc00007c1e0]   IndexAddr[0xc0000a6010]
// Value[rabbit]   Address[0xc00007c1e0]   IndexAddr[0xc0000a6020]
// Value[fish]     Address[0xc00007c1e0]   IndexAddr[0xc0000a6030]
// Value[hamster]  Address[0xc00007c1e0]   IndexAddr[0xc0000a6040]
