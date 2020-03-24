// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to takes slices of slices to create different
// views of and make changes to the underlying array.
package main

import "fmt"

func main() {

	// Create a slice with a length of 5 elements and a capacity of 8.
	slice1 := make([]string, 5, 8)
	slice1[0] = "Apple"
	slice1[1] = "Orange"
	slice1[2] = "Banana"
	slice1[3] = "Grape"
	slice1[4] = "Plum"

	inspectSlice(slice1)

	// Take a slice of slice1. We want just indexes 2 and 3.
	// Parameters are [starting_index : (starting_index + length)]
	slice2 := slice1[2:4]
	inspectSlice(slice2)

	fmt.Println("*************************")

	// Change the value of the index 0 of slice2.
	slice2[0] = "CHANGED"

	// Display the change across all existing slices.
	inspectSlice(slice1)
	inspectSlice(slice2)

	fmt.Println("*************************")

	// Make a new slice big enough to hold elements of slice 1 and copy the
	// values over using the builtin copy function.
	slice3 := make([]string, len(slice1))
	copy(slice3, slice1)
	inspectSlice(slice3)
}

// inspectSlice exposes the slice header for review.
func inspectSlice(slice []string) {
	fmt.Printf("Length[%d] Capacity[%d]\n", len(slice), cap(slice))
	for i, s := range slice {
		fmt.Printf("[%d] %p %s\n",
			i,
			&slice[i],
			s)
	}
}

// Outputs:
// Length[5] Capacity[8]
// [0] 0xc00007c000 Apple
// [1] 0xc00007c010 Orange
// [2] 0xc00007c020 Banana
// [3] 0xc00007c030 Grape
// [4] 0xc00007c040 Plum
// Length[2] Capacity[6]
// [0] 0xc00007c020 Banana
// [1] 0xc00007c030 Grape
// *************************
// Length[5] Capacity[8]
// [0] 0xc00007c000 Apple
// [1] 0xc00007c010 Orange
// [2] 0xc00007c020 CHANGED
// [3] 0xc00007c030 Grape
// [4] 0xc00007c040 Plum
// Length[2] Capacity[6]
// [0] 0xc00007c020 CHANGED
// [1] 0xc00007c030 Grape
// *************************
// Length[5] Capacity[5]
// [0] 0xc000080000 Apple
// [1] 0xc000080010 Orange
// [2] 0xc000080020 CHANGED
// [3] 0xc000080030 Grape
// [4] 0xc000080040 Plum
