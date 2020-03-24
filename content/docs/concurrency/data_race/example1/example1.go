// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go build -race

// Sample program to show how to create race conditions in
// our programs. We don't want to do this.
package main

import (
	"fmt"
	"sync"
)

// counter is a variable incremented by all goroutines.
var counter int

func main() {

	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {

				// Capture the value of Counter.
				value := counter

				// Increment our local value of Counter.
				value++

				fmt.Println(value)

				// Store the value back into Counter.
				counter = value
			}

			wg.Done()
		}()
	}

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

/*
==================
WARNING: DATA RACE
Read at 0x0000005fa138 by goroutine 8:
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:35 +0x47

Previous write at 0x0000005fa138 by goroutine 7:
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:41 +0x63

Goroutine 8 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:31 +0xab

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:31 +0xab
==================
Final Counter: 4
Found 1 data race(s)
*/
