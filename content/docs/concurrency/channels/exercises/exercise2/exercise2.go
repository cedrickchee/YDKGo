// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Write a program that uses a fan out pattern to generate 100 random numbers
// concurrently. Have each goroutine generate a single random number and return
// that number to the main goroutine over a buffered channel. Set the size of
// the buffered channel so no send ever blocks. Don't allocate more capacity
// than you need. Have the main goroutine store each random number it receives
// in a slice. Print the slice values then terminate the program.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Declare constant for number of goroutines.
const grs = 100

func init() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Create the buffered channel with room for
	// each goroutine to be created.
	ch := make(chan int, grs)

	// Iterate and launch each goroutine.
	for g := 0; g < grs; g++ {
		// Create an anonymous function for each goroutine that
		// generates a random number and sends it on the channel.
		go func() {
			ch <- rand.Intn(500)
		}()
	}

	// Create a variable to be used to track received messages.
	// Set the value to the number of goroutines created.
	wait := grs

	// Iterate receiving each value until they are all received.
	// Store them in a slice of ints.
	var nums []int
	for wait > 0 {
		nums = append(nums, <-ch)
		wait--
	}

	// Print the values in our slice.
	fmt.Println(nums)
}

// Outputs:
// [409 22 279 402 120 442 304 457 199 402 330 40 43 299 167 0 57 271 19 351 479
// 2 35 378 81 403 324 346 233 172 206 425 394 275 410 8 245 123 91 458 313 389
// 143 242 236 342 396 172 30 430 301 466 491 379 63 278 479 52 114 253 71 11
// 381 469 327 439 278 274 448 264 267 150 11 370 237 474 251 18 74 320 402 380
// 249 469 472 38 158 447 399 482 197 326 286 194 76 239 0 132 409 312]
