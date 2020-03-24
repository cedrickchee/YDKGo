// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Declare a struct type to maintain information about a user (name, email and age).
// Create a value of this type, initialize with values and display each field.
//
// Declare and initialize an anonymous struct type with the same three fields. Display the value.
package main

import "fmt"

// Add imports.

// user represents a user in the system
type user struct {
	name  string
	age   int
	email string
}

func main() {

	// Declare variable of type user and init using a struct literal.
	john := user{
		name:  "John",
		age:   30,
		email: "john@foo.bar",
	}

	// Display the field values.
	fmt.Println("Name", john.name)
	fmt.Println("Age", john.age)
	fmt.Println("Email", john.email)

	// Declare a variable using an anonymous struct.
	david := struct {
		name  string
		age   int
		email string
	}{
		name:  "David",
		age:   32,
		email: "david@foo.bar",
	}

	// Display the field values.
	fmt.Println("Name", david.name)
	fmt.Println("Age", david.age)
	fmt.Println("Email", david.email)
}
