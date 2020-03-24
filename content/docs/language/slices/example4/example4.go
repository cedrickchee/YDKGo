// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to grow a slice using the built-in function append
// and how append grows the capacity of the underlying array.
package main

import "fmt"

func main() {

	// Declare a nil slice of strings.
	var data []string

	// Capture the capacity of the slice.
	lastCap := cap(data)

	// Append ~100k strings to the slice.
	for record := 1; record <= 1e5; record++ {

		// Use the built-in function append to add to the slice.
		value := fmt.Sprintf("Rec: %d", record)
		data = append(data, value)

		// When the capacity of the slice changes, display the changes.
		if lastCap != cap(data) {

			// Calculate the percent of change.
			capChg := float64(cap(data)-lastCap) / float64(lastCap) * 100

			// Save the new values for capacity.
			lastCap = cap(data)

			// Display the results.
			fmt.Printf("Addr[%p]\tIndex[%d]\t\tCap[%d - %2.f%%]\n",
				&data[0],
				record,
				cap(data),
				capChg)
		}
	}
}

// Outputs:
// Addr[0xc000010200]      Index[1]                Cap[1 - +Inf%]
// Addr[0xc00000c080]      Index[2]                Cap[2 - 100%]
// Addr[0xc000064080]      Index[3]                Cap[4 - 100%]
// Addr[0xc00007e000]      Index[5]                Cap[8 - 100%]
// Addr[0xc000080000]      Index[9]                Cap[16 - 100%]
// Addr[0xc00007c200]      Index[17]               Cap[32 - 100%]
// Addr[0xc000082000]      Index[33]               Cap[64 - 100%]
// Addr[0xc000084000]      Index[65]               Cap[128 - 100%]
// Addr[0xc000079000]      Index[129]              Cap[256 - 100%]
// Addr[0xc000086000]      Index[257]              Cap[512 - 100%]
// Addr[0xc00008a000]      Index[513]              Cap[1024 - 100%]
// Addr[0xc000090000]      Index[1025]             Cap[1280 - 25%]
// Addr[0xc00009a000]      Index[1281]             Cap[1704 - 33%]
// Addr[0xc0000b2000]      Index[1705]             Cap[2560 - 50%]
// Addr[0xc0000c0000]      Index[2561]             Cap[3584 - 40%]
// Addr[0xc0000d4000]      Index[3585]             Cap[4608 - 29%]
// Addr[0xc0000ec000]      Index[4609]             Cap[6144 - 33%]
// Addr[0xc00010e000]      Index[6145]             Cap[7680 - 25%]
// Addr[0xc000134000]      Index[7681]             Cap[9728 - 27%]
// Addr[0xc000166000]      Index[9729]             Cap[12288 - 26%]
// Addr[0xc0001a6000]      Index[12289]            Cap[15360 - 25%]
// Addr[0xc0001f4000]      Index[15361]            Cap[19456 - 27%]
// Addr[0xc000258000]      Index[19457]            Cap[24576 - 26%]
// Addr[0xc0002d6000]      Index[24577]            Cap[30720 - 25%]
// Addr[0xc000372000]      Index[30721]            Cap[38400 - 25%]
// Addr[0xc000434000]      Index[38401]            Cap[48128 - 25%]
// Addr[0xc00053c000]      Index[48129]            Cap[60416 - 26%]
// Addr[0xc000628000]      Index[60417]            Cap[75776 - 25%]
// Addr[0xc000750000]      Index[75777]            Cap[94720 - 25%]
// Addr[0xc0008c2000]      Index[94721]            Cap[118784 - 25%]
