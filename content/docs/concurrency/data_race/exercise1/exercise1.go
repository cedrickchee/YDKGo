// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Fix the race condition in this program.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// numbers maintains a set of random numbers.
var numbers []int

// mutex will help protect the slice.
var mutex sync.Mutex

// init is called prior to main.
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Number of goroutines to use.
	const grs = 3

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	// Create three goroutines to generate random numbers.
	for i := 0; i < grs; i++ {
		go func() {
			random(10)
			wg.Done()
		}()
	}

	// Wait for all the goroutines to finish.
	wg.Wait()

	// Display the set of random numbers.
	for i, number := range numbers {
		fmt.Println(i, number)
	}
}

// random generates random numbers and stores them into a slice.
func random(amount int) {

	// Generate as many random numbers as specified.
	for i := 0; i < amount; i++ {
		n := rand.Intn(100)

		mutex.Lock()
		{
			numbers = append(numbers, n)
		}
		mutex.Unlock()
	}
}

/*
==================
WARNING: DATA RACE
Read at 0x0000005f9200 by goroutine 9:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0x92
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Previous write at 0x0000005f9200 by goroutine 7:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0xf6
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Goroutine 9 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc
==================
==================
WARNING: DATA RACE
Read at 0x0000005f9200 by goroutine 8:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0x92
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Previous write at 0x0000005f9200 by goroutine 7:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0xf6
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Goroutine 8 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc
==================
==================
WARNING: DATA RACE
Read at 0x00c0000ac048 by goroutine 8:
  runtime.growslice()
      /usr/local/go/src/runtime/slice.go:76 +0x0
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0x168
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Previous write at 0x00c0000ac048 by goroutine 7:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0xd3
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Goroutine 8 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc
==================
==================
WARNING: DATA RACE
Read at 0x00c0000ac048 by goroutine 9:
  runtime.growslice()
      /usr/local/go/src/runtime/slice.go:76 +0x0
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0x168
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Previous write at 0x00c0000ac048 by goroutine 7:
  main.random()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:55 +0xd3
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:35 +0x37

Goroutine 9 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/exercise1/exercise1.go:34 +0xbc
==================
0 79
1 1
2 1
3 60
4 67
5 91
6 5
7 86
8 46
9 52
10 51
11 39
12 78
13 21
14 94
15 59
16 72
17 96
18 11
19 35
Found 4 data race(s)
*/
