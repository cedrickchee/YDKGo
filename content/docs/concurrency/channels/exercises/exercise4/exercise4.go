// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Write a program that creates a fixed set of workers to generate random
// numbers. Discard any number divisible by 2. Continue receiving until 100
// numbers are received. Tell the workers to shut down before terminating.
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

func main() {
	// Create the channel for sharing results.
	ch := make(chan int)

	// Create a channel "shutdown" to tell goroutines when to terminate.
	shutdown := make(chan struct{})

	// Define the size of the worker pool. Use runtime.NumCPU to size the pool based on number of processors.
	poolSize := runtime.NumCPU()

	// Create a sync.WaitGroup to monitor the Goroutine pool. Add the count.
	var wg sync.WaitGroup
	wg.Add(poolSize)

	// Create a fixed size pool of goroutines to generate random numbers.
	for g := 0; g < poolSize; g++ {
		go func(id int) {
			// Start an infinite loop.
			for {
				// Generate a random number up to 1000.
				n := rand.Intn(1000)

				// Use a select to either send the number or receive the shutdown signal.
				select {
				case ch <- n:
					// In one case send the random number.
					fmt.Printf("Worker %d sent %d\n", id, n)
				case <-shutdown:
					// In another case receive from the shutdown channel.
					fmt.Printf("Worker %d shutting down\n", id)
					wg.Done()
					return
				}
			}
		}(g)
	}

	// Create a slice to hold the random numbers.
	var nums []int

	// Receive from the values channel with range.
	for n := range ch {
		// continue the loop if the value was even.
		if n%2 == 0 {
			fmt.Println("Discarding", n)
			continue
		}

		// Store the odd number.
		fmt.Println("Keeping", n)
		nums = append(nums, n)

		// break the loop once we have 100 results.
		if len(nums) == 100 {
			break
		}
	}

	// Send the shutdown signal by closing the shutdown channel.
	fmt.Println("Receiver sending shutdown signal")
	close(shutdown)

	// Wait for the Goroutines to finish.
	wg.Wait()

	// Print the values in our slice.
	fmt.Printf("Result count: %d\n", len(nums))
	fmt.Println(nums)
}

/* Outputs:
Worker 0 sent 81
Keeping 81
Keeping 887
Keeping 847
Keeping 59
Keeping 81
Worker 0 sent 81
Worker 0 sent 318
Discarding 318
Keeping 425
Worker 0 sent 425
Worker 0 sent 540
Discarding 540
Discarding 456
Worker 0 sent 456
Worker 0 sent 300
Discarding 300
Discarding 694
Worker 0 sent 694
Worker 0 sent 511
Keeping 511
Discarding 162
Worker 0 sent 162
Worker 0 sent 89
Keeping 89
Discarding 728
Worker 0 sent 728
Worker 0 sent 274
Worker 1 sent 847
Worker 2 sent 59
Discarding 274
Keeping 211
Keeping 445
Keeping 237
Worker 2 sent 237
Worker 2 sent 106
Discarding 106
Keeping 495
Worker 1 sent 445
Worker 1 sent 466
Worker 2 sent 495
Discarding 466
Discarding 528
Discarding 258
Worker 3 sent 887
Worker 3 sent 47
Keeping 47
Keeping 947
Worker 3 sent 947
Worker 3 sent 287
Worker 1 sent 528
Worker 2 sent 258
Worker 0 sent 211
Keeping 287
Discarding 888
Discarding 790
Keeping 15
Keeping 541
Worker 1 sent 790
Worker 2 sent 15
Worker 1 sent 408
Discarding 408
Keeping 387
Keeping 831
Worker 1 sent 831
Worker 1 sent 429
Keeping 429
Discarding 356
Worker 1 sent 356
Worker 1 sent 737
Keeping 737
Keeping 631
Worker 3 sent 888
Worker 3 sent 485
Keeping 485
Discarding 26
Worker 3 sent 26
Worker 3 sent 413
Keeping 413
Discarding 90
Worker 3 sent 90
Worker 3 sent 194
Discarding 194
Keeping 563
Worker 3 sent 563
Worker 3 sent 433
Keeping 433
Keeping 147
Worker 3 sent 147
Worker 3 sent 78
Discarding 78
Discarding 324
Worker 3 sent 324
Worker 3 sent 159
Keeping 159
Keeping 353
Worker 3 sent 353
Worker 3 sent 957
Keeping 957
Keeping 721
Worker 3 sent 721
Worker 3 sent 189
Keeping 189
Keeping 199
Worker 3 sent 199
Worker 3 sent 0
Discarding 0
Keeping 705
Worker 3 sent 705
Worker 3 sent 888
Discarding 888
Discarding 538
Worker 3 sent 538
Worker 3 sent 703
Keeping 703
Keeping 355
Worker 3 sent 355
Worker 3 sent 451
Keeping 451
Discarding 510
Worker 3 sent 510
Worker 3 sent 605
Keeping 605
Discarding 156
Worker 3 sent 156
Worker 3 sent 266
Discarding 266
Discarding 828
Worker 3 sent 828
Worker 3 sent 561
Keeping 561
Discarding 202
Worker 3 sent 202
Worker 3 sent 783
Keeping 783
Discarding 746
Worker 3 sent 746
Worker 3 sent 563
Keeping 563
Discarding 376
Worker 3 sent 376
Worker 3 sent 2
Discarding 2
Discarding 718
Worker 3 sent 718
Worker 3 sent 447
Keeping 447
Discarding 94
Worker 3 sent 94
Worker 3 sent 577
Keeping 577
Keeping 463
Worker 3 sent 463
Worker 3 sent 996
Discarding 996
Discarding 420
Worker 3 sent 420
Worker 3 sent 623
Keeping 623
Keeping 953
Worker 3 sent 953
Worker 3 sent 137
Keeping 137
Keeping 133
Worker 3 sent 133
Worker 3 sent 241
Worker 1 sent 631
Worker 0 sent 541
Keeping 241
Keeping 59
Keeping 33
Keeping 643
Worker 0 sent 643
Worker 0 sent 891
Keeping 891
Discarding 2
Worker 0 sent 2
Worker 0 sent 878
Discarding 878
Discarding 336
Worker 0 sent 336
Worker 0 sent 546
Discarding 546
Keeping 107
Worker 0 sent 107
Worker 0 sent 940
Discarding 940
Keeping 503
Worker 0 sent 503
Worker 0 sent 552
Discarding 552
Keeping 843
Worker 0 sent 843
Worker 0 sent 205
Keeping 205
Discarding 598
Worker 0 sent 598
Worker 0 sent 425
Keeping 425
Keeping 351
Worker 0 sent 351
Worker 0 sent 515
Keeping 515
Keeping 757
Worker 0 sent 757
Worker 0 sent 687
Keeping 687
Discarding 10
Worker 0 sent 10
Worker 0 sent 410
Discarding 410
Keeping 285
Worker 0 sent 285
Worker 0 sent 590
Discarding 590
Discarding 632
Worker 0 sent 632
Worker 0 sent 98
Discarding 98
Keeping 553
Worker 0 sent 553
Worker 0 sent 591
Keeping 591
Discarding 582
Worker 0 sent 582
Worker 0 sent 384
Discarding 384
Keeping 297
Worker 0 sent 297
Worker 0 sent 267
Keeping 267
Keeping 137
Worker 0 sent 137
Worker 0 sent 271
Keeping 271
Discarding 894
Worker 0 sent 894
Worker 0 sent 726
Discarding 726
Discarding 802
Worker 0 sent 802
Worker 0 sent 981
Keeping 981
Keeping 79
Worker 0 sent 79
Worker 0 sent 66
Discarding 66
Discarding 270
Worker 0 sent 270
Worker 0 sent 493
Keeping 493
Discarding 86
Worker 0 sent 86
Worker 0 sent 819
Keeping 819
Keeping 981
Worker 0 sent 981
Worker 0 sent 52
Discarding 52
Keeping 175
Worker 0 sent 175
Worker 0 sent 885
Keeping 885
Discarding 710
Worker 0 sent 710
Worker 0 sent 387
Keeping 387
Keeping 749
Worker 0 sent 749
Worker 0 sent 528
Discarding 528
Discarding 818
Worker 0 sent 818
Worker 0 sent 384
Discarding 384
Keeping 903
Worker 0 sent 903
Worker 3 sent 59
Worker 0 sent 224
Discarding 224
Keeping 547
Discarding 612
Worker 0 sent 612
Worker 0 sent 532
Discarding 532
Discarding 616
Worker 0 sent 616
Worker 0 sent 839
Keeping 839
Discarding 540
Worker 0 sent 540
Worker 0 sent 786
Worker 1 sent 33
Worker 3 sent 547
Worker 2 sent 387
Discarding 786
Keeping 51
Discarding 76
Discarding 640
Keeping 351
Worker 2 sent 351
Worker 2 sent 844
Discarding 844
Discarding 364
Worker 2 sent 364
Worker 0 sent 51
Worker 3 sent 640
Worker 2 sent 305
Keeping 305
Keeping 183
Worker 0 sent 183
Keeping 801
Worker 1 sent 76
Discarding 90
Discarding 602
Discarding 258
Worker 1 sent 258
Worker 1 sent 767
Keeping 767
Keeping 231
Worker 1 sent 231
Worker 1 sent 578
Discarding 578
Discarding 154
Worker 1 sent 154
Worker 1 sent 822
Worker 3 sent 801
Worker 0 sent 602
Worker 2 sent 90
Discarding 822
Keeping 223
Discarding 342
Discarding 208
Keeping 743
Worker 2 sent 743
Worker 2 sent 968
Discarding 968
Discarding 166
Worker 2 sent 166
Worker 2 sent 710
Discarding 710
Keeping 535
Worker 2 sent 535
Worker 2 sent 440
Discarding 440
Discarding 904
Worker 2 sent 904
Worker 2 sent 162
Discarding 162
Keeping 657
Worker 2 sent 657
Worker 2 sent 415
Keeping 415
Keeping 371
Worker 2 sent 371
Worker 2 sent 39
Keeping 39
Discarding 430
Worker 2 sent 430
Worker 2 sent 513
Keeping 513
Discarding 700
Worker 2 sent 700
Worker 2 sent 359
Keeping 359
Discarding 720
Worker 2 sent 720
Worker 2 sent 783
Keeping 783
Discarding 870
Worker 2 sent 870
Worker 2 sent 984
Discarding 984
Keeping 247
Worker 2 sent 247
Worker 2 sent 10
Discarding 10
Keeping 565
Receiver sending shutdown signal
Worker 2 sent 565
Worker 2 shutting down
Worker 1 sent 223
Worker 1 shutting down
Worker 3 sent 342
Worker 3 shutting down
Worker 0 sent 208
Worker 0 shutting down
Result count: 100
[81 887 847 59 81 425 511 89 211 445 237 495 47 947 287 15 541 387 831 429 737
631 485 413 563 433 147 159 353 957 721 189 199 705 703 355 451 605 561 783 563
447 577 463 623 953 137 133 241 59 33 643 891 107 503 843 205 425 351 515 757
687 285 553 591 297 267 137 271 981 79 493 819 981 175 885 387 749 903 547 839
51 351 305 183 801 767 231 223 743 535 657 415 371 39 513 359 783 247 565]
*/
