// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show the basic concept of using a pointer
// to share data.
package main

func main() {

	// Declare variable of type int with a value of 10.
	count := 10

	// Display the "value of" and "address of" count.
	println("count:\tValue Of[", count, "]\t\t\tAddr Of[", &count, "]")

	// Pass the "address of" count.
	// This is still considered pass by value, not by reference because the address itself is a value.
	increment(&count)

	// count is updated.
	println("count:\tValue Of[", count, "]\t\t\tAddr Of[", &count, "]")
}

// increment declares count as a pointer variable whose value is
// always an address and points to values of type int.
// The * here is not an operator. It is part of the type name.
// Every type that is declared, whether you declare or it is predeclared, you get for free a pointer.
//go:noinline
func increment(inc *int) {

	// Increment the "value of" count that the "pointer points to".
	// The * is an operator. It tells us the value of the pointer points to.
	*inc++

	println("inc:\tValue Of[", inc, "]\tAddr Of[", &inc, "]\tValue Points To [", *inc, "]")
}

// Outputs:
// count:  Value Of[ 10 ]                  Addr Of[ 0xc00003a748 ]
// inc:    Value Of[ 0xc00003a748 ]        Addr Of[ 0xc00003a738 ] Value Points To [ 11 ]
// count:  Value Of[ 11 ]                  Addr Of[ 0xc00003a748 ]
