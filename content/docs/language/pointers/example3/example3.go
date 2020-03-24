// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show the basic concept of using a pointer
// to share data.
package main

import "fmt"

// user represents a user in the system.
type user struct {
	name   string
	email  string
	logins int
}

func main() {

	// Declare and initialize a variable named john of type user.
	john := user{
		name:  "John",
		email: "john@foo.bar",
	}

	//** We don't need to include all the fields when specifying field
	// names with a struct literal.

	// Pass the "address of" the john value.
	display(&john)

	// Pass the "address of" the logins field from within the john value.
	increment(&john.logins)

	// Pass the "address of" the john value.
	display(&john)
}

// increment declares logins as a pointer variable whose value is
// always an address and points to values of type int.
func increment(logins *int) {
	*logins++
	fmt.Printf("&logins[%p] logins[%p] *logins[%d]\n\n", &logins, logins, *logins)
}

// display declares u as user pointer variable whose value is always an address
// and points to values of type user.
func display(u *user) {
	fmt.Printf("%p\t%+v\n", u, *u)
	fmt.Printf("Name: %q Email: %q Logins: %d\n\n", u.name, u.email, u.logins)
}

// Outputs:
// 0xc00007a150    {name:John email:john@foo.bar logins:0}
// Name: "John" Email: "john@foo.bar" Logins: 0
//
// &logins[0xc00000e030] logins[0xc00007a170] *logins[1]
//
// 0xc00007a150    {name:John email:john@foo.bar logins:1}
// Name: "John" Email: "john@foo.bar" Logins: 1
