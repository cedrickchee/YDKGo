// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show that you cannot take the address
// of an element in a map.
package main

// player represents someone playing our game.
type player struct {
	name  string
	score int
}

func main() {

	// Declare a map with initial values using a map literal.
	players := map[string]player{
		"anna":  {"Anna", 42},
		"jacob": {"Jacob", 21},
	}

	// Trying to take the address of a map element fails.
	anna := &players["anna"]
	anna.score++

	// Compiler error: "./example4.go:25:10: cannot take the address of players["anna"]"

	// Instead take the element, modify it, and put it back.
	player := players["anna"]
	player.score++
	players["anna"] = player
}
