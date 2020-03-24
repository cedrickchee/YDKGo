// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how the for range has both value and pointer semantics.
package main

import "fmt"

func main() {

	// Using the value semantic form of the for range.
	friends := []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	for _, v := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", v)
	}

	// Using the pointer semantic form of the for range.
	friends = []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	for i := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", friends[i])
	}
}

// Outputs:
// v[Annie]
// v[Betty]
// v[Charley]
// v[Doug]
// v[Edward]
// v[Annie]
// v[Betty]
// panic: runtime error: index out of range [2] with length 2

// goroutine 1 [running]:
// main.main()
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/slices/example8/example8.go:24 +0x1fd
// exit status 2
