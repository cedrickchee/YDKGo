// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Declare a struct that represents a baseball player. Include name, atBats and hits.
// Declare a method that calculates a player's batting average. The formula is hits / atBats.
// Declare a slice of this type and initialize the slice with several players. Iterate over
// the slice displaying the players name and batting average.
package main

import "fmt"

// Add imports.

// Declare a struct that represents a ball player.
// Include fields called name, atBats and hits.
type player struct {
	name   string
	atBats int
	hits   int
}

// Declare a method that calculates the batting average for a player.
func (p *player) average() float64 {
	avg := float64(p.hits) / float64(p.atBats)
	return avg
}

func main() {

	// Create a slice of players and populate each player
	// with field values.
	p := []player{
		{"john", 9, 1},
		{"david", 15, 5},
	}

	// Display the batting average for each player in the slice.
	for i := range p {
		fmt.Printf("%s, avg: .%.f\n", p[i].name, p[i].average()*1000)

		// fmt.Printf("%s: addr[%p], indexaddr[%p]\n", p[i].name, &v, &p[i])
	}
}
