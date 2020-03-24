// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to use a read/write mutex to define critical
// sections of code that needs synchronous access.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// data is a slice that will be shared.
var data []string

// rwMutex is used to define a critical section of code.
var rwMutex sync.RWMutex

// Number of reads occurring at ay given time.
var readCount int64

// init is called prior to main.
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a writer goroutine.
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			writer(i)
		}
		wg.Done()
	}()

	// Create eight reader goroutines.
	for i := 0; i < 8; i++ {
		go func(id int) {
			for {
				reader(id)
			}
		}(i)
	}

	// Wait for the write goroutine to finish.
	wg.Wait()
	fmt.Println("Program Complete")
}

// writer adds a new string to the slice in random intervals.
func writer(i int) {

	// Only allow one goroutine to read/write to the slice at a time.
	rwMutex.Lock()
	{
		// Capture the current read count.
		// Keep this safe though we can due without this call.
		rc := atomic.LoadInt64(&readCount)

		// Perform some work since we have a full lock.
		fmt.Printf("****> : Performing Write : RCount[%d]\n", rc)
		data = append(data, fmt.Sprintf("String: %d", i))
	}
	rwMutex.Unlock()
	// Release the lock.
}

// reader wakes up and iterates over the data slice.
func reader(id int) {

	// Any goroutine can read when no write operation is taking place.
	rwMutex.RLock()
	{
		// Increment the read count value by 1.
		rc := atomic.AddInt64(&readCount, 1)

		// Perform some read work and display values.
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Printf("%d : Performing Read : Length[%d] RCount[%d]\n", id, len(data), rc)

		// Decrement the read count value by 1.
		atomic.AddInt64(&readCount, -1)
	}
	rwMutex.RUnlock()
	// Release the read lock.
}

// Outputs:
// ... ... snipped ... ...
// ... ... snipped ... ...
// 1 : Performing Read : Length[0] RCount[8]
// 2 : Performing Read : Length[0] RCount[8]
// 2 : Performing Read : Length[0] RCount[8]
// 7 : Performing Read : Length[0] RCount[8]
// 1 : Performing Read : Length[0] RCount[8]
// 0 : Performing Read : Length[0] RCount[8]
// 3 : Performing Read : Length[0] RCount[8]
// ... ... snipped ... ...
// ... ... snipped ... ...
// 0 : Performing Read : Length[8] RCount[8]
// 7 : Performing Read : Length[8] RCount[8]
// 5 : Performing Read : Length[8] RCount[8]
// 4 : Performing Read : Length[8] RCount[8]
// 2 : Performing Read : Length[8] RCount[8]
// 6 : Performing Read : Length[8] RCount[8]
// 0 : Performing Read : Length[8] RCount[8]
// 1 : Performing Read : Length[8] RCount[8]
// 3 : Performing Read : Length[8] RCount[8]
// ****> : Performing Write : RCount[0]
// 0 : Performing Read : Length[9] RCount[3]
// 0 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[6]
// 2 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[2]
// 5 : Performing Read : Length[9] RCount[4]
// 7 : Performing Read : Length[9] RCount[1]
// 3 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[7]
// 1 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[5]
// 3 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 0 : Performing Read : Length[9] RCount[8]
// 4 : Performing Read : Length[9] RCount[8]
// 6 : Performing Read : Length[9] RCount[8]
// 1 : Performing Read : Length[9] RCount[8]
// 3 : Performing Read : Length[9] RCount[8]
// 2 : Performing Read : Length[9] RCount[8]
// 5 : Performing Read : Length[9] RCount[8]
// 7 : Performing Read : Length[9] RCount[8]
// ****> : Performing Write : RCount[0]
// Program Complete
