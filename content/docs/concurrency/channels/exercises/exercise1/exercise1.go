// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Write a program where two goroutines pass an integer back and forth
// ten times. Display when each goroutine receives the integer. Increment
// the integer with each pass. Once the integer equals ten, terminate
// the program cleanly.
package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create an unbuffered channel.
	ch := make(chan int)

	// Create the WaitGroup and add a count
	// of two, one for each goroutine.
	var wg sync.WaitGroup
	wg.Add(2)

	// Launch the goroutine and handle Done.
	go func() {
		goroutine("John", ch)
		wg.Done()
	}()

	// Launch the goroutine and handle Done.
	go func() {
		goroutine("David", ch)
		wg.Done()
	}()

	// Send a value to start the counting.
	ch <- 1

	// Wait for the program to finish.
	wg.Wait()
}

// goroutine simulates sharing a value.
func goroutine(name string, ch chan int) {
	for {
		// Wait for the value to be sent.
		value, ok := <-ch
		if !ok {
			// If the channel was closed, return.
			fmt.Printf("Goroutine %s down\n", name)
			return
		}

		// Display the value.
		fmt.Printf("Goroutine %s inc %d\n", name, value)

		// Terminate when the value is 10.
		if value == 10 {
			close(ch)
			fmt.Printf("Goroutine %s down\n", name)
			return
		}

		// Increment the value and send it
		// over the channel.
		ch <- (value + 1)
	}
}

// Outputs:
// Goroutine David inc 1
// Goroutine John inc 2
// Goroutine David inc 3
// Goroutine John inc 4
// Goroutine David inc 5
// Goroutine John inc 6
// Goroutine David inc 7
// Goroutine John inc 8
// Goroutine David inc 9
// Goroutine John inc 10
// Goroutine John down
// Goroutine David down
