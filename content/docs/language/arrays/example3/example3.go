// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how the behavior of the for range and
// how memory for an array is contiguous.
package main

import "fmt"

func main() {

	// Declare an array of 5 strings initialized with values.
	friends := [5]string{"Annie", "Betty", "Charley", "Doug", "Edward"}

	// Iterate over the array displaying the value and
	// address of each element.
	for i, v := range friends {
		fmt.Printf("Value[%s]\tAddress[%p] IndexAddr[%p]\n", v, &v, &friends[i])
	}
}

// Outputs:
// Value[Annie]    Address[0xc00009a040] IndexAddr[0xc00009c000]
// Value[Betty]    Address[0xc00009a040] IndexAddr[0xc00009c010]
// Value[Charley]  Address[0xc00009a040] IndexAddr[0xc00009c020]
// Value[Doug]     Address[0xc00009a040] IndexAddr[0xc00009c030]
// Value[Edward]   Address[0xc00009a040] IndexAddr[0xc00009c040]
